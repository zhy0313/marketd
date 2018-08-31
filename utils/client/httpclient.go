package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

func newHttpClient() *http.Client {
	tr := &http.Transport{
		Proxy: proxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return &http.Client{
		Transport: tr,
	}
}

func cloneRequest(req *http.Request) *http.Request {
	dup := new(http.Request)
	*dup = *req
	dup.URL, _ = url.Parse(req.URL.String())
	dup.Header = make(http.Header)
	for k, s := range req.Header {
		dup.Header[k] = s
	}
	return dup
}

// An implementation of http.ProxyFromEnvironment that isn't broken
func proxyFromEnvironment(req *http.Request) (*url.URL, error) {
	proxy := os.Getenv("https_proxy")
	if proxy == "" {
		proxy = os.Getenv("HTTPS_PROXY")
	}
	if proxy == "" {
		return nil, nil
	}

	proxyURL, err := url.Parse(proxy)
	if err != nil || !strings.HasPrefix(proxyURL.Scheme, "https") {
		if proxyURL, err := url.Parse("http://" + proxy); err == nil {
			return proxyURL, nil
		}
	}

	if err != nil {
		return nil, fmt.Errorf("invalid proxy address %q: %v", proxy, err)
	}

	return proxyURL, nil
}

type browser struct {
	httpClient *http.Client
}

type HttpRequest struct {
	Url    string
	Body   interface{}
	Header http.Header
}

func (b *browser) performRequest(method, path string, body io.Reader, header http.Header) (*HttpResponse, error) {
	url, err := url.Parse(path)
	if err == nil {
		return b.performRequestUrl(method, url, body, header)
	} else {
		return nil, err
	}
}

func (b *browser) performRequestUrl(method string, url *url.URL, body io.Reader, header http.Header) (res *HttpResponse, err error) {
	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return
	}

	var bodyBackup io.ReadWriter
	if req.Body != nil {
		bodyBackup = &bytes.Buffer{}
		req.Body = ioutil.NopCloser(io.TeeReader(req.Body, bodyBackup))
	}

	req.Header = header

	httpResponse, err := b.httpClient.Do(req)
	if err != nil {
		return
	}

	res = &HttpResponse{httpResponse}

	return
}

func (b *browser) jsonRequest(method, path string, body interface{}, header http.Header) (*HttpResponse, error) {
	json, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(json)

	header.Set("Content-Type", "application/json; charset=utf-8")
	header.Add("Accept-Language", "zh-cn")

	return b.performRequest(method, path, buf, header)
}

func (b *browser) Get(req *HttpRequest) (*HttpResponse, error) {
	return b.performRequest("GET", req.Url, nil, req.Header)
}

func (b *browser) GetFile(path string, mimeType string, header http.Header) (*HttpResponse, error) {
	header.Set("Accept", mimeType)
	return b.performRequest("GET", path, nil, header)
}

func (b *browser) PostJSON(path string, payload interface{}, header http.Header) (*HttpResponse, error) {
	return b.jsonRequest("POST", path, payload, header)
}

type HttpResponse struct {
	*http.Response
}

func (res *HttpResponse) String() (string, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (res *HttpResponse) Link(name string) string {
	linkVal := res.Header.Get("Link")
	re := regexp.MustCompile(`<([^>]+)>; rel="([^"]+)"`)
	for _, match := range re.FindAllStringSubmatch(linkVal, -1) {
		if match[2] == name {
			return match[1]
		}
	}
	return ""
}

//获取http请求并将结果反序列化到一个接口上
func HttpGet(req *HttpRequest) (string, error) {
	httpClient := newHttpClient()

	var client = browser{
		httpClient: httpClient,
	}

	resp, err := client.Get(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	content, err := resp.String()
	if err != nil {
		return "", err
	}

	return content, nil
}
