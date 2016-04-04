package requests

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"net"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	TLSConfig           *tls.Config
	DialTimeout         time.Duration
	DialKeepAlive       time.Duration
	TLSHandshakeTimeout time.Duration
	DiableKeepAlive     bool
	MaxIdelConnsPerHost int
}

func httpCall(action, url string, request interface{}, headers map[string]string, response interface{}, requestFuncs ...func(*Config)) (*http.Response, error) {
	var requestType string
	if headers == nil {
		headers = map[string]string{
			"Content-Type": "application/json; charset=utf-8",
		}
	}
	// 默认使用json
	if headers["Content-Type"] == "" {
		headers["Content-Type"] = "application/json; charset=utf-8"
	}

	debugf("httpCall default request header: %+v", headers)

	requestType = headers["Content-Type"]
	var data []byte
	if strings.Contains(requestType, "xml") {
		data, _ = xml.Marshal(request)
	} else {
		data, _ = json.Marshal(request)
	}

	debugf("httpCall method: %s, url %s", action, url)

	req, err := http.NewRequest(action, url, bufio.NewReader(bytes.NewBuffer(data)))
	if err != nil {
		return &http.Response{}, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	var config = &Config{
		TLSConfig:           &tls.Config{InsecureSkipVerify: true},
		DialTimeout:         10 * time.Second,
		DialKeepAlive:       30 * time.Second,
		TLSHandshakeTimeout: 10 & time.Second,
		MaxIdelConnsPerHost: 1,
		DiableKeepAlive:     false,
	}

	// 覆盖原有的方法
	for i := 0; i < len(requestFuncs); i++ {
		f := requestFuncs[i]
		f(config)
	}

	clt := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: config.TLSConfig,
			Proxy:           http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   config.DialTimeout,
				KeepAlive: config.DialKeepAlive,
			}).Dial,
			TLSHandshakeTimeout: config.TLSHandshakeTimeout,
			DisableKeepAlives:   config.DiableKeepAlive,
			MaxIdleConnsPerHost: config.MaxIdelConnsPerHost,
		},
	}
	// 发送请求
	resp, err := clt.Do(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	// 解析结果
	responseType := resp.Header.Get("Content-Type")
	if strings.Contains(responseType, "json") {
		err = json.NewDecoder(resp.Body).Decode(response)
		return resp, err
	} else if strings.Contains(responseType, "xml") {
		err = xml.NewDecoder(resp.Body).Decode(response)
		return resp, err
	}
	return resp, err
}

func Post(url string, request interface{}, headers map[string]string, response interface{}, requestFuncs ...func(*Config)) (*http.Response, error) {
	return httpCall("POST", url, request, headers, response, requestFuncs...)
}

func Put(url string, request interface{}, headers map[string]string, response interface{}, requestFuncs ...func(*Config)) (*http.Response, error) {
	return httpCall("PUT", url, request, headers, response, requestFuncs...)
}

func Delete(url string, request interface{}, headers map[string]string, response interface{}, requestFuncs ...func(*Config)) (*http.Response, error) {
	return httpCall("DELETE", url, request, headers, response, requestFuncs...)
}
func Get(url string, headers map[string]string, response interface{}, requestFuncs ...func(*Config)) (*http.Response, error) {
	return httpCall("GET", url, nil, headers, response, requestFuncs...)
}
