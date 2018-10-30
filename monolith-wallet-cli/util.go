package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

func HTTPGet(url string, v interface{}) error {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, v)
}
