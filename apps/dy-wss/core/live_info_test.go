package core

import (
	"net/url"
	"strings"
	"testing"
)

func TestGetXMsStub(t *testing.T) {
	params := []kv{
		{"live_id", "1"},
		{"aid", "6383"},
		{"version_code", "180800"},
		{"webcast_sdk_version", "1.0.14-beta.0"},
		{"room_id", "7519127586287504191"},
		{"sub_room_id", ""},
		{"sub_channel_id", ""},
		{"did_rule", "3"},
		{"user_unique_id", "7692890738690032830"},
		{"device_platform", "web"},
		{"device_type", ""},
		{"ac", ""},
		{"identity", "audience"},
	}
	got := getXMsStub(params)
	if len(got) != 32 {
		t.Fatalf("md5 hex length: got %d want 32, value=%q", len(got), got)
	}
	// stable for fixed input
	again := getXMsStub(params)
	if got != again {
		t.Fatalf("not deterministic: %q vs %q", got, again)
	}
}

func TestBuildRequestURL(t *testing.T) {
	raw := "wss://webcast5-ws-web-lf.douyin.com/webcast/im/push/v2/?room_id=1&compress=gzip&signature=deadbeef"
	got, err := buildRequestURL(raw)
	if err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(got)
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()
	for _, key := range []string{"aid", "device_platform", "browser_language", "browser_platform", "browser_name", "browser_version", "room_id", "signature"} {
		if q.Get(key) == "" {
			t.Fatalf("missing query key %q in %s", key, got)
		}
	}
	if q.Get("aid") != "6383" {
		t.Fatalf("aid=%q", q.Get("aid"))
	}
	if q.Get("browser_name") != "Mozilla" {
		t.Fatalf("browser_name=%q", q.Get("browser_name"))
	}
	if !strings.HasPrefix(q.Get("browser_version"), "5.0") {
		t.Fatalf("browser_version=%q", q.Get("browser_version"))
	}
}

func TestGetSignature_SyntheticStub(t *testing.T) {
	// synthetic X-MS-STUB (md5 hex); no network
	stub := "0123456789abcdef0123456789abcdef"
	sig := getSignature(stub, liveUserAgent)
	if sig == "" {
		t.Fatal("empty signature")
	}
	if sig == "00000000" {
		t.Fatalf("goja/webmssdk fallback stub; get_sign failed (check engine/SDK load)")
	}
	t.Logf("signature=%s (len=%d)", sig, len(sig))
}

func TestGetUserUniqueID(t *testing.T) {
	id := getUserUniqueID()
	if len(id) < 18 {
		t.Fatalf("unexpected id length: %q", id)
	}
}
