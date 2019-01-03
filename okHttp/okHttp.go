package okHttp

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	MaxRetry   = 2
	MaxTimeOut = 200 //毫秒
)

//BuildQuery http_build_query
func BuildQuery(params map[string]string) string {
	if len(params) <= 0 {
		return ""
	}
	var query bytes.Buffer
	for k, v := range params {
		query.WriteString(fmt.Sprintf("%s=%s&", k, strings.TrimSpace(v)))
	}
	return strings.TrimRight(query.String(), "&")
}

func urlSlave(reqURL string) ([]string, error) {
	scheme := ""
	if strings.HasPrefix(reqURL, "http://") {
		scheme = "http://"
	} else if strings.HasPrefix(reqURL, "https://") {
		scheme = "https://"
	} else {
		return nil, errors.New("request url error")
	}
	reqURL = strings.TrimPrefix(reqURL, scheme)
	sps := strings.SplitN(reqURL, "/", 2)
	hostStr := sps[0]
	uri := ""
	if len(sps) == 2 {
		uri = sps[1]
	}

	hosts := strings.Split(hostStr, "|")

	var urls = make([]string, len(hosts))
	for k, host := range hosts {
		var url string
		if uri == "" {
			url = scheme + host
		} else {
			url = scheme + host + "/" + uri
		}
		urls[k] = url
	}

	return urls, nil
}

//Get get
func Get(url string, timeout int, retry int, header map[string]string) ([]byte, map[string]string, error) {
	var result []byte
	var err error
	var ext map[string]string
	if timeout <= 0 {
		timeout = MaxTimeOut
	}
	if retry < 0 {
		retry = MaxRetry
	}
	urls, err := urlSlave(url)
	if err != nil {
		return nil, nil, err
	}

	ln := len(urls)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < retry+1; i++ {
		for j := 0; j < ln; j++ {
			index := rand.Intn(ln)
			u := urls[index]
			result, ext, err = Request(1, u, "", header, timeout)
			if err != nil {
				continue
			} else {
				goto END
			}
		}
	}

END:
	if err != nil {
		return nil, ext, err
	}
	return result, ext, nil
}

//Post post
func Post(url string, params string, timeout int, retry int, header map[string]string) ([]byte, map[string]string, error) {
	var result []byte
	var err error
	var ext map[string]string
	if timeout <= 0 {
		timeout = MaxTimeOut
	}
	if retry < 0 {
		retry = MaxRetry
	}
	urls, err := urlSlave(url)
	if err != nil {
		return nil, nil, err
	}

	ln := len(urls)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < retry+1; i++ {
		for j := 0; j < ln; j++ {
			index := rand.Intn(ln)
			u := urls[index]
			result, ext, err = Request(2, u, params, header, timeout)
			if err != nil {
				continue
			} else {
				goto END
			}
		}
	}

END:
	if err != nil {
		return nil, ext, err
	}
	return result, ext, nil
}

/**
* method GET 1, POST 2
**/
func Request(method int, url string, params string, header map[string]string, timeout int) ([]byte, map[string]string, error) {
	var start = time.Now()
	var req *http.Request
	var err error
	var ext = make(map[string]string)
	if method == 1 {
		req, err = http.NewRequest("GET", url, nil)
	} else {
		req, err = http.NewRequest("POST", url, strings.NewReader(params))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if err != nil {
		return nil, ext, err
	}
	if header != nil {
		for k, v := range header {
			if k == "Host" {
				req.Host = v
			} else {
				req.Header.Set(k, v)
			}
		}
	}
	if timeout <= 0 {
		timeout = MaxTimeOut
	}
	client := &http.Client{
		Timeout: time.Millisecond * time.Duration(timeout),
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, ext, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, ext, err
	}

	ext["status"] = strconv.Itoa(resp.StatusCode)
	ext["totalTime"] = strconv.Itoa(int(time.Now().Sub(start)/1000000)) + "ms"
	return body, ext, nil
}

/**
* get real ip
**/

func RealIP(r *http.Request) string {
	ra := r.RemoteAddr
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		ra = strings.Split(ip, ", ")[0]
	} else if ip := r.Header.Get("X-Real-IP"); ip != "" {
		ra = ip
	} else {
		ra, _, _ = net.SplitHostPort(ra)
	}
	return ra
}
