package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var default_timeout = 5 * time.Second

func Get(url string, params map[string]string, headers map[string]string) ([]byte, error) {
	return GetWithTimeOut(url, default_timeout, params, headers)
}

func GetWithTimeOut(url string, timeout time.Duration, params map[string]string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return nil, errors.New("new request is fail ")
	}

	q := req.URL.Query()
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}

	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}

	client := &http.Client{Timeout: timeout}

	resp, error := client.Do(req)
	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func Post(url string, body map[string]string, params map[string]string, headers map[string]string) ([]byte, error) {
	return PostWithTimeOut(url, default_timeout, body, params, headers)
}

func PostWithTimeOut(url string, timeout time.Duration, body map[string]string, params map[string]string, headers map[string]string) ([]byte, error) {
	var bodyJson []byte
	var req *http.Request
	if body != nil {
		var err error
		bodyJson, err = json.Marshal(body)
		if err != nil {
			log.Println(err)
			return nil, errors.New("http post body to json failed")
		}
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Println(err)
		return nil, errors.New("new request is fail: %v \n")
	}
	req.Header.Set("Content-type", "application/json")

	q := req.URL.Query()
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}

	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}

	client := &http.Client{Timeout: timeout}
	resp, error := client.Do(req)
	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
