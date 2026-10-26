package core

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dop251/goja"
)

const liveUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0"

var (
	reRoomInfo = regexp.MustCompile(`room\\":\{.*\\"id_str\\":\\"(\d+)\\".*,\\"status\\":(\d+).*"title\\":\\"([^"]*)\\"`)
	reRoomID   = regexp.MustCompile(`roomId\\":\\"(\d+)\\"`)
	reOwner    = regexp.MustCompile(`owner\\":(.*?),\\"room_auth`)
	reM3U8Map  = regexp.MustCompile(`hls_pull_url_map\\":(\{.*?})`)
	reFLV      = regexp.MustCompile(`flv\\":\\"(.*?)\\"`)
)

// ParseLiveRoomURL ports dy-live parseLiveRoomUrl (signature via JS deferred).
func ParseLiveRoomURL(roomURL string) (*LiveInfo, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, roomURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("User-Agent", liveUserAgent)
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	html := string(body)

	ttwid := ""
	u, _ := url.Parse(roomURL)
	for _, c := range jar.Cookies(u) {
		if c.Name == "ttwid" {
			ttwid = c.Value
			break
		}
	}
	if ttwid == "" {
		return nil, fmt.Errorf("ttwid cookie missing")
	}

	roomTitle := ""
	if m := reRoomInfo.FindStringSubmatch(html); len(m) == 4 {
		roomTitle = m[3]
		log.Printf("房间标题: %s", roomTitle)
		if m[2] == "4" {
			return nil, fmt.Errorf("直播已结束")
		}
	}

	roomIDMatch := reRoomID.FindStringSubmatch(html)
	if len(roomIDMatch) < 2 {
		return nil, fmt.Errorf("roomId not found in page")
	}
	liveRoomID := roomIDMatch[1]

	ownerMatch := reOwner.FindStringSubmatch(html)
	if len(ownerMatch) < 2 {
		return nil, fmt.Errorf("owner info not found in page")
	}
	ownerJSON := strings.ReplaceAll(ownerMatch[1], `\"`, `"`)
	var userInfo json.RawMessage
	if err := json.Unmarshal([]byte(ownerJSON), &userInfo); err != nil {
		return nil, fmt.Errorf("parse owner json: %w", err)
	}
	log.Printf("主播账号信息: %s", ownerJSON)

	m3u8URL := ""
	if m := reM3U8Map.FindStringSubmatch(html); len(m) >= 2 {
		mapJSON := strings.ReplaceAll(m[1], `\"`, `"`)
		var m3u8s map[string]string
		if err := json.Unmarshal([]byte(mapJSON), &m3u8s); err != nil {
			return nil, fmt.Errorf("parse hls_pull_url_map: %w", err)
		}
		m3u8URL = strings.Replace(m3u8s["FULL_HD1"], "http", "https", 1)
		if m3u8URL == "" {
			m3u8URL = strings.Replace(m3u8s["HD1"], "http", "https", 1)
		}
		log.Printf("直播流m3u8链接地址是: %s", m3u8URL)
	}

	flvURL := ""
	if m := reFLV.FindStringSubmatch(html); len(m) >= 2 {
		flvURL = strings.ReplaceAll(m[1], `\"`, `"`)
		flvURL = strings.ReplaceAll(flvURL, `\u0026`, "&")
		if !strings.Contains(flvURL, "https") {
			flvURL = strings.Replace(flvURL, "http", "https", 1)
		}
		log.Printf("直播流FLV地址是: %s", flvURL)
	}

	userUniqueID := getUserUniqueID()
	const (
		versionCode       = 180800
		webcastSDKVersion = "1.0.14-beta.0"
	)
	sigParams := []kv{
		{"live_id", "1"},
		{"aid", "6383"},
		{"version_code", strconv.Itoa(versionCode)},
		{"webcast_sdk_version", webcastSDKVersion},
		{"room_id", liveRoomID},
		{"sub_room_id", ""},
		{"sub_channel_id", ""},
		{"did_rule", "3"},
		{"user_unique_id", userUniqueID},
		{"device_platform", "web"},
		{"device_type", ""},
		{"ac", ""},
		{"identity", "audience"},
	}
	signature := getSignature(getXMsStub(sigParams), liveUserAgent)

	webcastParams := []kv{
		{"room_id", liveRoomID},
		{"compress", "gzip"},
		{"version_code", strconv.Itoa(versionCode)},
		{"webcast_sdk_version", webcastSDKVersion},
		{"live_id", "1"},
		{"did_rule", "3"},
		{"user_unique_id", userUniqueID},
		{"identity", "audience"},
		{"signature", signature},
	}
	q := make([]string, 0, len(webcastParams))
	for _, p := range webcastParams {
		q = append(q, p.k+"="+p.v)
	}
	wssURL := "wss://webcast5-ws-web-lf.douyin.com/webcast/im/push/v2/?" + strings.Join(q, "&")
	wssURL, err = buildRequestURL(wssURL)
	if err != nil {
		return nil, err
	}

	info := &LiveInfo{
		RoomUrl:   roomURL,
		RoomTitle: roomTitle,
		RoomId:    liveRoomID,
		WssUrl:    wssURL,
		Ttwid:     ttwid,
		M3U8Url:   m3u8URL,
		FlvUrl:    flvURL,
		UserAgent: liveUserAgent,
	}
	_ = json.Unmarshal(userInfo, &info.UserInfo)
	return info, nil
}

type kv struct{ k, v string }

func getUserUniqueID() string {
	n := rand.Int63n(7999999999999999999-7300000000000000000) + 7300000000000000000
	return strconv.FormatInt(n, 10)
}

func getXMsStub(params []kv) string {
	parts := make([]string, 0, len(params))
	for _, p := range params {
		parts = append(parts, p.k+"="+p.v)
	}
	sum := md5.Sum([]byte(strings.Join(parts, ",")))
	return hex.EncodeToString(sum[:])
}

// getSignature runs embedded webmssdk.js via goja (same shape as dy-live).
func getSignature(xMsStub, userAgent string) string {
	vm := goja.New()
	jsDom := fmt.Sprintf(`
document = {};
window = {};
navigator = { userAgent: %s };
`, strconv.Quote(userAgent))
	if _, err := vm.RunString(jsDom + string(webmssdkJS)); err != nil {
		log.Printf("getSignature: load webmssdk failed: %v", err)
		return "00000000"
	}
	v, err := vm.RunString(fmt.Sprintf("get_sign(%s)", strconv.Quote(xMsStub)))
	if err != nil {
		log.Printf("getSignature: get_sign failed: %v", err)
		return "00000000"
	}
	sig, ok := v.Export().(string)
	if !ok || sig == "" {
		log.Printf("getSignature: unexpected result type=%T value=%v", v.Export(), v.Export())
		return "00000000"
	}
	return sig
}

func buildRequestURL(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("aid", "6383")
	q.Set("device_platform", "web")
	q.Set("browser_language", "zh-CN")
	q.Set("browser_platform", "Win32")
	browserName := strings.SplitN(liveUserAgent, "/", 2)[0]
	q.Set("browser_name", browserName)
	// match Python: USER_AGENT.split(browser_name)[-1][1:]
	rest := strings.Split(liveUserAgent, browserName)
	browserVersion := ""
	if len(rest) > 0 {
		tail := rest[len(rest)-1]
		if len(tail) > 0 {
			browserVersion = tail[1:]
		}
	}
	q.Set("browser_version", browserVersion)
	u.RawQuery = q.Encode()
	log.Printf("new_query_string: %s", u.RawQuery)
	return u.String(), nil
}
