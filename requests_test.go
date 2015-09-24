package requests

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type httpBin struct {
	Args    struct{} `json:"args"`
	Data    string   `json:"data"`
	Files   struct{} `json:"files"`
	Form    struct{} `json:"form"`
	Headers struct {
		Accept_Encoding string `json:"Accept-Encoding"`
		Content_Length  string `json:"Content-Length"`
		Content_Type    string `json:"Content-Type"`
		Host            string `json:"Host"`
		User_Agent      string `json:"User-Agent"`
	} `json:"headers"`
	JSON struct {
		Hello string `json:"hello"`
	} `json:"json"`
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

type echo struct {
	Hello string `json:"hello"`
}

func TestParse2Struct(t *testing.T) {
	var bin httpBin
	testJSON := `
	{
		"args": {},
		"data": "{\"hello\":\"world\"}",
		"files": {},
		"form": {},
		"headers": {
		"Accept-Encoding": "gzip",
		"Content-Length": "17",
		"Content-Type": "application/json; charset=utf-8",
		"Host": "httpbin.org",
		"User-Agent": "beegoServer"
		},
		"json": {
			"hello": "world"
		},
		"origin": "118.244.254.30",
		"url": "http://httpbin.org/post"
	}
	`
	err := Parse2Struct("json", []byte(testJSON), &bin)
	if err != nil {
		t.Fatal(err)
	}
	if bin.JSON.Hello != "world" {
		t.Errorf("want world, got %s", bin.JSON.Hello)
	}
}

func TestParse2Bytes(t *testing.T) {
	var bin httpBin
	bin.Headers.Content_Length = "application/json; charset=utf-8"
	b, err := Parse2Bytes("json", bin)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "application/json; charset=utf-8") {
		t.Errorf("want true, got %t", strings.Contains(string(b), "application/json; charset=utf-8"))
	}
}

func TestConvertResponseToBytes(t *testing.T) {
	resp, err := http.Get("https://api.github.com")
	if err != nil {
		t.Fatal(err)
	}
	b1, err := ConvertResponseToBytes(resp)
	if err != nil {
		t.Fatal(err)
	}
	b2, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(b1) != string(b2) {
		t.Errorf("want true, got %t", string(b1) == string(b2))
	}
}

func TestGet(t *testing.T) {
	var bin httpBin
	// var url = "http://www.baidu.com"
	var url = "http://httpbin.org/get"
	result, err := Get(url, nil, &bin)
	if err != nil {
		m, _ := ConvertResponseToBytes(result)
		t.Log(string(m))
		t.Fatal(err)
	}
	if len(strings.Split(bin.Origin, ".")) != 4 {
		t.Errorf("want 4, got %d", len(strings.Split(bin.Origin, ".")))
	}
}

func TestPost(t *testing.T) {
	var bin httpBin
	request := &echo{Hello: "world"}
	var url = "http://httpbin.org/post"

	result, err := Post(url, request, nil, &bin)
	if err != nil {
		m, _ := ConvertResponseToBytes(result)
		t.Log(string(m))
		t.Fatal(err)
	}
	if bin.JSON.Hello != request.Hello {
		t.Errorf("want %s, got %+v", request.Hello, bin.JSON.Hello)
	}
}

func TestDelete(t *testing.T) {
	var bin httpBin
	request := &echo{Hello: "world"}
	var url = "http://httpbin.org/delete"
	result, err := Delete(url, request, nil, &bin)
	if err != nil {
		m, _ := ConvertResponseToBytes(result)
		t.Log(string(m))
		t.Fatal(err)
	}
	if bin.JSON.Hello != request.Hello {
		t.Errorf("want %s, got %+v", request.Hello, bin.JSON.Hello)
	}
}

func TestPut(t *testing.T) {
	var bin httpBin
	request := &echo{Hello: "world"}
	var url = "http://httpbin.org/put"
	result, err := Put(url, request, nil, &bin)
	if err != nil {
		m, _ := ConvertResponseToBytes(result)
		t.Log(string(m))
		t.Fatal(err)
	}
	if bin.JSON.Hello != request.Hello {
		t.Errorf("want %s, got %+v", request.Hello, bin.JSON.Hello)
	}
}
