package main

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Request struct {
	Method  string
	Url     string
	Body    string
	Headers map[string]string
}

type Response struct {
	Code    int
	Reason  string
	Body    string
	Headers map[string]string
}

func doRequest(req *Request) *Response {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	request, err := http.NewRequest(req.Method, req.Url, strings.NewReader(req.Body))
	if err != nil {
		Error(err.Error())
		os.Exit(-1)
	}
	for k, v := range req.Headers {
		request.Header.Set(k, v)
	}
	if _, ok := req.Headers["User-Agent"]; !ok {
		request.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; "+
			"Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) "+
			"Chrome/92.0.4515.131 Mobile Safari/537.36")
	}
	if _, ok := req.Headers["Content-Type"]; !ok {
		request.Header.Set("Content-Type", "application/json")
	}
	resp, err := client.Do(request)
	if err != nil {
		Error(err.Error())
		os.Exit(-1)
	}
	response := &Response{
		Code:   resp.StatusCode,
		Reason: resp.Status,
	}
	respHeader := make(map[string]string)
	for k, v := range resp.Header {
		respHeader[k] = strings.Join(v, "; ")
	}
	response.Headers = respHeader
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Error(err.Error())
		os.Exit(-1)
	}
	response.Body = string(respBody)
	return response
}
