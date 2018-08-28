package utils

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
	proxy := os.Getenv("http_proxy")
	if proxy == "" {
		proxy = os.Getenv("HTTP_PROXY")
	}
	if proxy == "" {
		return nil, nil
	}

	proxyURL, err := url.Parse(proxy)
	if err != nil || !strings.HasPrefix(proxyURL.Scheme, "http") {
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
	userAgent  string
}

func (b *browser) performRequest(method, path string, body io.Reader, configure func(*http.Request)) (*HttpResponse, error) {
	url, err := url.Parse(path)
	if err == nil {
		return b.performRequestUrl(method, url, body, configure)
	} else {
		return nil, err
	}
}

func (b *browser) performRequestUrl(method string, url *url.URL, body io.Reader, configure func(*http.Request)) (res *HttpResponse, err error) {
	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", b.userAgent)

	if configure != nil {
		configure(req)
	}

	var bodyBackup io.ReadWriter
	if req.Body != nil {
		bodyBackup = &bytes.Buffer{}
		req.Body = ioutil.NopCloser(io.TeeReader(req.Body, bodyBackup))
	}

	httpResponse, err := b.httpClient.Do(req)
	if err != nil {
		return
	}

	res = &HttpResponse{httpResponse}

	return
}

func (b *browser) jsonRequest(method, path string, body interface{}, configure func(*http.Request)) (*HttpResponse, error) {
	json, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(json)

	return b.performRequest(method, path, buf, func(req *http.Request) {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		if configure != nil {
			configure(req)
		}
	})
}

func (b *browser) Get(path string) (*HttpResponse, error) {
	return b.performRequest("GET", path, nil, nil)
}

func (b *browser) GetFile(path string, mimeType string) (*HttpResponse, error) {
	return b.performRequest("GET", path, nil, func(req *http.Request) {
		req.Header.Set("Accept", mimeType)
	})
}

func (b *browser) Delete(path string) (*HttpResponse, error) {
	return b.performRequest("DELETE", path, nil, nil)
}

func (b *browser) PostJSON(path string, payload interface{}) (*HttpResponse, error) {
	return b.jsonRequest("POST", path, payload, nil)
}

func (b *browser) PatchJSON(path string, payload interface{}) (*HttpResponse, error) {
	return b.jsonRequest("PATCH", path, payload, nil)
}

func (b *browser) PostFile(path, filename string) (*HttpResponse, error) {
	stat, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return b.performRequest("POST", path, file, func(req *http.Request) {
		req.ContentLength = stat.Size()
		req.Header.Set("Content-Type", "application/octet-stream")
	})
}

type HttpResponse struct {
	*http.Response
}

func (res *HttpResponse) String() string {
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if nil != err {
		return err.Error()
	}

	return string(body)
}

func (res *HttpResponse) Unmarshal(dest interface{}) (err error) {
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	return json.Unmarshal(body, dest)
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

type HttpRequest struct {
	Url       string
	Payload   map[string]string
	UserAgent string
}

//获取http请求并将结果反序列化到一个接口上
func HttpGet(req HttpRequest, data interface{}) (string, error) {
	httpClient := newHttpClient()

	var client = browser{
		httpClient: httpClient,
		userAgent:  req.UserAgent,
	}

	var url string
	if nil == req.Payload {
		url = req.Url
	} else {
		query := Map2UrlQuery(req.Payload)
		url = req.Url + "?" + query
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}

	if err = resp.Unmarshal(data); err != nil {
		return "", err
	}

	return resp.String(), nil
}
