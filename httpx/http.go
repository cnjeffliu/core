/*
 * @Author: Jeffrey Liu
 * @Date: 2021-07-19 12:00:32
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-12-14 17:40:38
 * @Description:
 */
package httpx

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type HttpxClient struct {
	cli *http.Client
}

type httpxOption func(r *http.Request)

const defaultTimeout = 5 * time.Second

var client *HttpxClient

func WithJSONContent() httpxOption {
	return func(r *http.Request) {
		r.Header.Add("Content-Type", "application/json; charset=utf-8")
	}
}

func WithHeaders(headers map[string]string) httpxOption {
	return func(r *http.Request) {
		if headers != nil {
			for key, val := range headers {
				r.Header.Add(key, val)
			}
		}
	}
}

func WithParams(params map[string]string) httpxOption {
	return func(r *http.Request) {
		q := r.URL.Query()
		if params != nil {
			for key, val := range params {
				q.Add(key, val)
			}
			r.URL.RawQuery = q.Encode()
		}
	}
}

func Client() *HttpxClient {
	if client == nil {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		cli := &http.Client{
			Transport: tr,
			Timeout:   defaultTimeout,
		}
		client = &HttpxClient{
			cli: cli,
		}
	}
	return client
}

func (h *HttpxClient) GetToMap(url string, opts ...httpxOption) (map[string]interface{}, error) {
	return h.GetToMapWithTimeOut(url, defaultTimeout, opts...)
}

func (h *HttpxClient) GetToMapWithTimeOut(url string, timeout time.Duration, opts ...httpxOption) (map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return nil, errors.New("new request is fail ")
	}

	for _, opt := range opts {
		opt(req)
	}

	h.cli.Timeout = timeout

	resp, error := h.cli.Do(req)
	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (h *HttpxClient) GetJSON(url string, opts ...httpxOption) ([]byte, error) {
	opts = append(opts, WithJSONContent())
	return h.Get(url, opts...)
}

func (h *HttpxClient) Get(url string, opts ...httpxOption) ([]byte, error) {
	return h.GetWithTimeOut(url, defaultTimeout, opts...)
}

func (h *HttpxClient) GetWithTimeOut(url string, timeout time.Duration, opts ...httpxOption) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return nil, errors.New("new request is fail ")
	}

	for _, opt := range opts {
		opt(req)
	}

	h.cli.Timeout = timeout

	resp, error := h.cli.Do(req)
	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (h *HttpxClient) PostToMap(url string, body interface{}, opts ...httpxOption) (map[string]interface{}, error) {
	return h.PostToMapWithTimeOut(url, body, defaultTimeout, opts...)
}

func (h *HttpxClient) PostToMapWithTimeOut(url string, body interface{}, timeout time.Duration, opts ...httpxOption) (map[string]interface{}, error) {
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
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Println(err)
		return nil, errors.New("new request is fail: %v \n")
	}

	for _, opt := range opts {
		opt(req)
	}

	h.cli.Timeout = timeout

	resp, error := h.cli.Do(req)
	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()

	var rawData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&rawData)
	if err != nil {
		return nil, err
	}
	return rawData, nil
}

func (h *HttpxClient) PostJSON(url string, body interface{}, opts ...httpxOption) ([]byte, error) {
	opts = append(opts, WithJSONContent())
	return h.Post(url, body, opts...)
}

func (h *HttpxClient) Post(url string, body interface{}, opts ...httpxOption) ([]byte, error) {
	return h.PostWithOptions(url, defaultTimeout, body, opts...)
}

func (h *HttpxClient) PostWithOptions(url string, timeout time.Duration, body interface{}, opts ...httpxOption) ([]byte, error) {
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
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Println(err)
		return nil, errors.New("new request is fail: %v \n")
	}

	for _, opt := range opts {
		opt(req)
	}

	h.cli.Timeout = timeout

	resp, error := h.cli.Do(req)
	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// GetIP returns request real ip.
func GetIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip
	}

	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}

	if net.ParseIP(ip) != nil {
		return ip
	}

	return ""
}
