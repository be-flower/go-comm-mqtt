package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	/*复用 http 链接*/
	httpClient *http.Client
)

const (
	MaxIdleConnections int = 20
	/*超时时间30s*/
	RequestTimeout int = 30
)

// init HTTPClient
func init() {
	httpClient = createHTTPClient()
}

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
			DisableKeepAlives:   false,
		},
		Timeout: time.Duration(RequestTimeout) * time.Second,
	}
	return client
}

func newHttpRequest(method string, url string, headers map[string]string, data []byte) (*http.Request, error) {
	var body io.Reader
	if len(data) > 0 {
		body = bytes.NewBuffer(data)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Errorf("newHttpRequest:NewRequest error! err = %v", err)
		return nil, err
	}

	if len(headers) > 0 {
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		if _, ok := headers["Content-Type"]; !ok {
			req.Header.Set("Content-Type", "application/json")
		}
	} else {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func DoHttpRequest(method string, url string, headers map[string]string, data []byte) (res []byte, err error) {
	fmt.Sprintf("DoHttpRequest:method = %v, url = %v, headers =%v, data = %v", method, url, headers, string(data))
	request, err := newHttpRequest(method, url, headers, data)
	if err != nil {
		fmt.Errorf("DoHttpRequest:newHttpRequest error! err = %v", err)
		return nil, err
	}

	resp, err := httpClient.Do(request)
	if err != nil {
		fmt.Errorf("DoHttpRequest:Do error! err = %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	// 判断返回状态
	if resp.StatusCode == http.StatusOK {
		// 读取返回的数据
		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			// 读取异常,返回错误
			fmt.Errorf("DoHttpRequest:send message (%v) ReadAll return error = %v", string(data), err)
			return nil, err
		}
		//fmt.Debugf("DoHttpRequest:data = %v", string(data))

		// 将收到的数据与状态返回
		return data, nil
	} else {
		fmt.Errorf("DoHttpRequest:send message (%v) ReadAll statuscode = %v,url = %v,", string(data), resp.StatusCode, url)

		// 读取返回的数据
		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			// 读取异常,返回错误
			fmt.Errorf("DoHttpRequest:send message (%v) ReadAll return error = %v", string(data), err)
			return nil, err
		}
		// 将收到的数据与状态返回
		return data, nil
	}
}
