package http

import (
	"fmt"
	"testing"
)

func TestServer(t *testing.T) {
	go httpServerRun()
	fmt.Printf("ok!")
}

func TestClient(t *testing.T) {
	httpClientRun("GET")
	httpClientRun("POST")
}
