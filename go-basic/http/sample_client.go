package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func httpClientRun(method string) {
	requestUrl := "http://127.0.0.1:8099/api/v1/hello_world"
	jsonStr := []byte(`{"hello": "world"}`)
	var (
		response *http.Response
		err      error
	)
	if method == "POST" {
		response, err = http.Post(requestUrl, "application/json", bytes.NewBuffer(jsonStr))
	} else {
		response, err = http.Get(requestUrl)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if err != nil {
		fmt.Println(err)
		return
	}
	// 记得处理 EOF
	// body := make([]byte, response.ContentLength)
	// _, err = response.Body.Read(body)
	// if err != nil && err != io.EOF {
	// 	fmt.Println(err)
	// 	return
	// }
	// 使用系统 API 时，系统内部会自动消费掉 EOF
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	fmt.Println(string(body))
}
