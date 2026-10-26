package main

import (
	"log"

	"github.com/sleticalboy/golearn/apps/dy-wss/core"
)

func main() {
	url := "https://live.douyin.com/534955839450"
	liveInfo := core.GetLiveInfo(url)
	log.Printf("liveInfo: %+v", liveInfo)
	core.WssServerStart(liveInfo)
}
