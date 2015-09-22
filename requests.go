package requests

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/astaxie/beego/httplib"
	"github.com/facebookgo/stackerr"
	"io/ioutil"
	"net/http"
	"strings"
)

func Parse2Struct(contentType string, data []byte, response interface{}) error {
	if strings.Contains(contentType, "xml") {
		return stackerr.Wrap(json.Unmarshal(data, &response))
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
	return ioutil.ReadAll(r.Body)
}

func Post(url string, request interface{}, headers map[string]string, response ...interface{}) (*http.Response, error) {
	return _x("POST", url, request, headers, response)
}

func Put(url string, request interface{}, headers map[string]string, response ...interface{}) (*http.Response, error) {
	return _x("PUT", url, request, headers, response)
}

func Delete(url string, request interface{}, headers map[string]string, response ...interface{}) (*http.Response, error) {
	return _x("DELETE", url, request, headers, response)
}
func Get(url string, headers map[string]string, response ...interface{}) (*http.Response, error) {
	return _x("GET", url, nil, headers, response...)
}

func _x(action, url string, request interface{}, headers map[string]string, response ...interface{}) (*http.Response, error) {
	var req *httplib.BeegoHttpRequest
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

	for k, v := range headers {
		req.Header(k, v)
	}

	resp, err := req.SendOut()
	if err != nil {
		return nil, stackerr.Wrap(err)
	}

	data, err := req.Bytes()
	if err != nil {
		return resp, stackerr.Wrap(err)
	}

	if len(response) == 0 {
		return resp, errors.New("LENGTH IS 0")
	}

	// 转成struct
	if action != "GET" {
		result, err := Parse2Bytes(headers["Content-Type"], request)
		if err != nil {
			return resp, stackerr.Wrap(err)
		}
		req.Body(result)
	}

	err = Parse2Struct(resp.Header.Get("Content-Type"), data, response[0])
	if err != nil {
		return resp, stackerr.Wrap(err)
	}
	return resp, nil

}
