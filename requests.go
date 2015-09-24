package requests

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/astaxie/beego/httplib"
	"github.com/facebookgo/stackerr"
)

type loopReader struct {
	*bytes.Buffer
}

func (r loopReader) Close() error {
	return nil
}

func Parse2Struct(contentType string, data []byte, response interface{}) error {
	if strings.Contains(contentType, "xml") {
		return stackerr.Wrap(xml.Unmarshal(data, &response))
	}
	if strings.Contains(contentType, "text/") {
		switch response.(type) {
		case string:
			response = string(data)
		default:
			response = &data
		}
		return nil
	}
	return stackerr.Wrap(json.Unmarshal(data, &response))
}

func Parse2Bytes(contentType string, request interface{}) ([]byte, error) {
	if strings.Contains(contentType, "json") {
		result, x := json.Marshal(request)
		return result, stackerr.Wrap(x)
	}
	if strings.Contains(contentType, "xml") {
		result, x := xml.Marshal(request)
		return result, stackerr.Wrap(x)
	}
	return nil, stackerr.New("NOT FOUND")
}

func ConvertResponseToBytes(r *http.Response) ([]byte, error) {
	if r == nil {
		err := errors.New("http response address nil")
		return []byte{}, err
	}
	buf, err := ioutil.ReadAll(r.Body)
	origin := loopReader{bytes.NewBuffer(buf)}
	r.Body = origin
	return buf, err
}

func Post(url string, request interface{}, headers map[string]string, response interface{}) (*http.Response, error) {
	return httpCall("POST", url, request, headers, response)
}

func Put(url string, request interface{}, headers map[string]string, response interface{}) (*http.Response, error) {
	return httpCall("PUT", url, request, headers, response)
}

func Delete(url string, request interface{}, headers map[string]string, response interface{}) (*http.Response, error) {
	return httpCall("DELETE", url, request, headers, response)
}
func Get(url string, headers map[string]string, response interface{}) (*http.Response, error) {
	return httpCall("GET", url, nil, headers, response)
}

func httpCall(action, url string, request interface{}, headers map[string]string, response ...interface{}) (*http.Response, error) {
	var req *httplib.BeegoHttpRequest
	debugf("httpCall method: %s, url %s", action, url)
	debugf("httpCall default request header: %+v", headers)
	switch strings.ToUpper(action) {
	case "PUT":
		req = httplib.Put(url)
	case "POST":
		req = httplib.Post(url)
	case "DELETE":
		req = httplib.Delete(url)
	case "GET":
		req = httplib.Get(url)
	default:
		req = httplib.Post(url)
	}

	if headers == nil || headers["Content-Type"] == "" {
		headers = map[string]string{
			"Content-Type": "application/json; charset=utf-8",
		}
	}

	debugf("httpCall request header: %+v", headers)

	for k, v := range headers {
		req.Header(k, v)
	}

	// 转成struct
	if action != "GET" {
		result, err := Parse2Bytes(headers["Content-Type"], request)
		if err != nil {
			return nil, stackerr.Wrap(err)
		}
		debugf("httpCall send body: %s", string(result))
		req.Body(result)
	}

	resp, err := req.SendOut()
	if err != nil {
		return nil, stackerr.Wrap(err)
	}
	debugf("httpCall response header %+v", resp.Header)
	data, err := ConvertResponseToBytes(resp)
	if err != nil {
		return resp, stackerr.Wrap(err)
	}
	debugf("httpCall response body: %s", string(data))

	if len(data) == 0 {
		return resp, errors.New("LENGTH IS 0")
	}

	err = Parse2Struct(resp.Header.Get("Content-Type"), data, response[0])
	if err != nil {
		return resp, stackerr.Wrap(err)
	}
	debugf("httpCall response struct: %+v", response[0])
	return resp, nil
}
