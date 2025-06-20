package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/donovanhide/eventsource"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
)

type Event struct {
	Event string
	Data  string
	ID    string
	Retry int
}

func processSSE(url string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "text/event-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	reader := bufio.NewReader(resp.Body)
	var event Event

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			log.Fatal("Error reading SSE stream:", err)
		}

		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			// 空行表示事件结束，处理完整事件
			handleEvent(event)
			event = Event{} // 重置事件
			continue
		}

		// 解析 SSE 字段
		if bytes.HasPrefix(line, []byte("event:")) {
			event.Event = string(bytes.TrimSpace(line[6:]))
		} else if bytes.HasPrefix(line, []byte("data:")) {
			event.Data += string(bytes.TrimSpace(line[5:])) + "\n"
		} else if bytes.HasPrefix(line, []byte("id:")) {
			event.ID = string(bytes.TrimSpace(line[3:]))
		} else if bytes.HasPrefix(line, []byte("retry:")) {
			fmt.Sscanf(string(line[6:]), "%d", &event.Retry)
		}
		// 忽略注释行和其他字段
	}
}

func handleEvent(event Event) {
	if event.Data == "" {
		return
	}
	// 去除 data 末尾的换行符
	event.Data = strings.TrimSuffix(event.Data, "\n")
	fmt.Printf("New Event [%s]\n", event.Event)
	fmt.Printf("ID: %s\n", event.ID)
	fmt.Printf("Data: %s\n", event.Data)
	if event.Retry != 0 {
		fmt.Printf("Retry: %dms\n", event.Retry)
	}
	var m map[string]any
	err := json.Unmarshal([]byte(event.Data), &m)
	fmt.Printf("Parsed JSON: %+v, err: %v\n", m, err)
	fmt.Println("-------------------")
}

func TestSse(t *testing.T) {
	processSSE("http://localhost:8099/api/v1/sse")
}

func TestSse2(t *testing.T) {
	stream, err := eventsource.Subscribe("http://localhost:8099/api/v1/sse", "")
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case ev := <-stream.Events:
			handleSSEEvent(ev)
		case err := <-stream.Errors:
			handleSSEError(err)
		case <-time.After(30 * time.Second):
			log.Println("No messages received for 30 seconds")
		}
	}
}

func handleSSEEvent(ev eventsource.Event) {
	fmt.Printf("Event: %s\n", ev.Event())
	fmt.Printf("ID: %s\n", ev.Id())
	fmt.Printf("Data: %s\n", ev.Data())
	// fmt.Printf("Timestamp: %v\n", ev.Timestamp())
	fmt.Println("-------------------")
}

func handleSSEError(err error) {
	log.Printf("Error: %v", err)
	// 这里可以实现重连逻辑
}
