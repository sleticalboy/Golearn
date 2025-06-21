package http

import (
	"bytes"
	"fmt"
	"net/http"
)

var (
	handlers = map[string]func(w http.ResponseWriter, r *http.Request){
		"/api/v1/hello_world": func(response http.ResponseWriter, request *http.Request) {
			dumpRequest(request)
			if request.Method == "GET" {
				response.WriteHeader(200)
				str := "Hello Go(GET) Http World!\n"
				_, _ = response.Write(bytes.NewBufferString(str).Bytes())
			} else if request.Method == "POST" {
				body := make([]byte, request.ContentLength)
				_, _ = request.Body.Read(body)
				fmt.Printf("request body is %s\n", string(body))
				str := "Hello Go(POST) Http World!\n"
				_, _ = response.Write(bytes.NewBufferString(str).Bytes())
				fmt.Printf("response body is %s\n", str)
			}
		},
	}
)

func dumpRequest(request *http.Request) {
	fmt.Println("request line >>>>>>>")
	fmt.Printf("%s %s %s from %s\n", request.Method, request.URL, request.Proto, request.RemoteAddr)
	fmt.Println("request headers >>>>>>>")
	for k, v := range request.Header {
		fmt.Printf("%s: %s\n", k, v)
	}
	fmt.Println("request body >>>>>>>")
	fmt.Printf("%v\n", request.Body)
}

func httpServerRun() {
	// 注册 api handler
	for api, handler := range handlers {
		http.HandleFunc(api, handler)
	}
	// 创建服务器
	server := &http.Server{
		Addr: "127.0.0.1:8099",
		Handler: nil,
	}
	// 启动服务器
	err := server.ListenAndServe()
	// err = http.ListenAndServe(server.Addr, nil)
	if err != nil {
		fmt.Printf("Listen http://127.0.0.1.8099 failed: %e\n", err)
	}
}
