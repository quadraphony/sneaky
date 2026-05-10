package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var allowed = map[string]bool{
	"sing-box": true,
	"ssh-socks": true,
	"wireguard": true,
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: validate <config.json>")
		os.Exit(2)
	}
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
		os.Exit(1)
	}
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		fmt.Fprintf(os.Stderr, "invalid json: %v\n", err)
		os.Exit(1)
	}
	a, ok := v["adapter"].(string)
	if !ok || !allowed[a] {
		fmt.Fprintf(os.Stderr, "invalid or missing adapter field\n")
		os.Exit(1)
	}
	if _, ok := v["config"].(map[string]interface{}); !ok {
		fmt.Fprintf(os.Stderr, "invalid or missing config field\n")
		os.Exit(1)
	}
	fmt.Println("config looks valid")
}
