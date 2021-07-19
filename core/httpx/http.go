package httpx

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type HttpxClient struct {
	cli *http.Client
}

type httpxClientOption func(r *http.Request)

const defaultTimeout = 5 * time.Second

var client *HttpxClient

func NewHttpxCli() *HttpxClient {
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

func (h *HttpxClient) GetJSON(url string, opts ...httpxClientOption) ([]byte, error) {
	args := []httpxClientOption{WithJSONContent()}
	for _, v := range opts {
		args = append(args, v)
	}
	return h.Get(url, args...)
}

func (h *HttpxClient) Get(url string, opts ...httpxClientOption) ([]byte, error) {
	return h.GetWithOptions(url, defaultTimeout, opts...)
}

func (h *HttpxClient) GetWithOptions(url string, timeout time.Duration, opts ...httpxClientOption) ([]byte, error) {
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

func (h *HttpxClient) PostJSON(url string, body interface{}, opts ...httpxClientOption) ([]byte, error) {
	args := []httpxClientOption{WithJSONContent()}
	for _, v := range opts {
		args = append(args, v)
	}
	return h.Post(url, body, args...)
}

func (h *HttpxClient) Post(url string, body interface{}, opts ...httpxClientOption) ([]byte, error) {
	return h.PostWithOptions(url, defaultTimeout, body, opts...)
}

func (h *HttpxClient) PostWithOptions(url string, timeout time.Duration, body interface{}, opts ...httpxClientOption) ([]byte, error) {
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

func WithJSONContent() httpxClientOption {
	return func(r *http.Request) {
		r.Header.Add("Content-Type", "application/json; charset=utf-8")
	}
}

func WithHeaders(headers map[string]string) httpxClientOption {
	return func(r *http.Request) {
		if headers != nil {
			for key, val := range headers {
				r.Header.Add(key, val)
			}
		}
	}
}

func WithParams(params map[string]string) httpxClientOption {
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
