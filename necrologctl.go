package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const defaultEndpoint = "http://localhost:8080/log"

type logRequest struct {
	Path  string `json:"path"`
	Level string `json:"level"`
	Msg   string `json:"msg"`
}

func main() {
	path := flag.String("path", "/var/log/uah_log/edotensei/system_info.log", "Log file path")
	level := flag.String("level", "info", "Log level (debug|info|warn|error)")
	msg := flag.String("msg", "", "Log message")
	endpoint := flag.String("endpoint", defaultEndpoint, "necrolog API endpoint")
	flag.Parse()

	if *msg == "" {
		fmt.Fprintln(os.Stderr, "[error] log message cannot be empty")
		os.Exit(1)
	}

	req := logRequest{
		Path:  *path,
		Level: *level,
		Msg:   *msg,
	}

	payload, err := json.Marshal(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[error] failed to marshal json: %v\n", err)
		os.Exit(1)
	}

	resp, err := http.Post(*endpoint, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Fprintf(os.Stderr, "[error] failed to send request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("[necrologctl] response: %s\n", string(body))
}
