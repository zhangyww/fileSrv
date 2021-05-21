/**
 * @Author: zhangyw
 * @Description:
 * @File:  redirect_test
 * @Date: 2021/5/20 11:54
 */

package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"testing"
	"time"
)

type MyRoundTriper struct {
}

func (this *MyRoundTriper) RoundTrip(request *http.Request) (*http.Response, error) {
	client := &http.Client{Timeout: 2 * time.Second}
	return client.Do(request)
}

func TestRedirect(t *testing.T) {
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		jsonBytes, _ := json.Marshal(req.URL)
		fmt.Println(string(jsonBytes))
		if len(req.URL.Path) < 3 {
			resp.Write([]byte("fail"))
			return
		}
		proxy := httputil.ReverseProxy{
			Director: func(request *http.Request) {
				request.URL.Scheme = "http"
				request.URL.Host = "172.16.16.212:9669"
				request.URL.Path = "/dir1/file2.txt"
			},
			//Transport: &MyRoundTriper{},
		}
		//newUrl, _ := url.Parse("http://172.16.16.212:9669/file1.txt")
		//proxy := httputil.NewSingleHostReverseProxy(newUrl)
		proxy.ServeHTTP(resp, req)
	})

	http.ListenAndServe("172.16.16.164:9669", nil)
}
