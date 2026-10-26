package main

import (
	"log"
	"os"

	"github.com/sleticalboy/golearn/apps/dy-wss/core"
)

func main() {
	url := "https://live.douyin.com/80017709309"
	if len(os.Args) > 1 {
		url = os.Args[1]
	}
	liveInfo := core.GetLiveInfo(url)
	log.Printf("liveInfo room_id=%s title=%s ttwid=%s", liveInfo.RoomId, liveInfo.RoomTitle, liveInfo.Ttwid)
	log.Printf("wss_url=%s", liveInfo.WssUrl)
	core.WssServerStart(liveInfo)
}
