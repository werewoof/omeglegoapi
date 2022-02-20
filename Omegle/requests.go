package omegle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (o *Omegle) PostRequest(url string, data interface{}) (string, error) {
	jsondata, err := json.Marshal(data)
	fmt.Println(string(jsondata))
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsondata))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "Application/json")
	resp, err := o.client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("recieved status code other than 200: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (o *Omegle) GetRequest(url string, data interface{}) (string, error) {
	jsondata, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(jsondata))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "Application/json")
	resp, err := o.client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("recieved status code other than 200: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
