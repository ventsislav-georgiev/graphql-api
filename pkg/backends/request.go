package backends

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/ventsislav-georgiev/graphql-api/pkg/backends/httpmethod"
)

type RequestOptions struct {
	Method  httpmethod.HttpMethod
	Url     string
	Params  map[string]string
	Data    map[string]interface{}
	Files   map[string]interface{}
	Header  http.Header
	Cookies []*http.Cookie
	Timeout *float32
}

type HttpClient interface {
	Request(requestOptions RequestOptions) (*http.Response, error)
}

type DefaultHttpClient struct{}

func (*DefaultHttpClient) Request(requestOptions RequestOptions) (*http.Response, error) {
	client := &http.Client{}

	url, err := url.Parse(requestOptions.Url)
	if err != nil {
		return nil, err
	}

	if requestOptions.Cookies != nil {
		client.Jar, err = cookiejar.New(nil)
		if err != nil {
			return nil, err
		}

		client.Jar.SetCookies(url, requestOptions.Cookies)
	}

	header := requestOptions.Header
	if header == nil {
		header = http.Header{}
	}

	header["User-Agent"] = []string{"APIMediation/0.0.1"}
	header["Content-Type"] = []string{"application/json"}

	req := &http.Request{
		Method:     requestOptions.Method.String(),
		URL:        url,
		Header:     header,
		ProtoMajor: 1,
		ProtoMinor: 1,
	}

	if requestOptions.Params != nil {
		query := url.Query()

		for k, v := range requestOptions.Params {
			query.Add(k, v)
		}

		url.RawQuery = query.Encode()
	}

	if requestOptions.Data != nil {
		dataBytes, err := json.Marshal(requestOptions.Data)
		if err != nil {
			return nil, err
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(dataBytes))
		req.ContentLength = int64(len(dataBytes))
		header["Content-Length"] = []string{strconv.Itoa(len(dataBytes))}
	}

	if requestOptions.Timeout != nil {
		client.Timeout = time.Duration(*requestOptions.Timeout)
	}

	if os.Getenv("reqdump") == "true" {
		dump, err := httputil.DumpRequest(req, true)
		if err != nil {
			return nil, err
		}
		fmt.Printf("\n## Request Dump ##\n%s\n## Request Dump ##\n", dump)
	} else {
		fmt.Printf("%s: %s\n", req.Method, req.URL)
	}

	resp, err := client.Do(req)

	if resp.StatusCode >= 300 {
		errMessage := resp.Status + " "

		if resp.Body != nil {
			defer resp.Body.Close()

			bytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}

			errMessage += string(bytes)
		}

		return nil, errors.New(errMessage)
	}

	return resp, err
}
