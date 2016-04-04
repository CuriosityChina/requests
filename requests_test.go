package requests

import (
	"crypto/tls"
	"strings"
	"testing"
	"time"
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

type httpBinXML struct {
	Slideshow struct {
		Title  string `xml:"title,attr"`
		Date   string `xml:"date,attr"`
		Author string `xml:"author,attr"`
		Slide  []struct {
			Type  string   `xml:"type,attr"`
			Title string   `xml:"title"`
			Item  []string `xml:"item"`
		} `xml:"slide"`
	} `xml:"slideshow"`
}

type echo struct {
	Hello string `json:"hello"`
}

func TestGet(t *testing.T) {
	var bin httpBin
	var url = "http://httpbin.org/get"
	_, err := Get(url, nil, &bin)
	if err != nil {
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

	_, err := Post(url, request, nil, &bin)
	if err != nil {
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
	_, err := Delete(url, request, nil, &bin)
	if err != nil {
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
	_, err := Put(url, request, nil, &bin)
	if err != nil {
		t.Fatal(err)
	}
	if bin.JSON.Hello != request.Hello {
		t.Errorf("want %s, got %+v", request.Hello, bin.JSON.Hello)
	}
}

func TestRequestFuncs(t *testing.T) {
	// 设置超时时间
	var timeout = func(clt *Config) {
		clt.DialTimeout = 10 * time.Second
	}

	var insecure = func(clt *Config) {
		clt.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	var bin httpBin
	request := &echo{Hello: "world"}
	var url = "https://httpbin.org/put"
	_, err := Put(url, request, nil, &bin, timeout, insecure)
	if err != nil {
		t.Fatal(err)
	}
	if bin.JSON.Hello != request.Hello {
		t.Errorf("want %s, got %+v", request.Hello, bin.JSON.Hello)
	}
}

func BenchmarkRequestGetParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var bin httpBin
			var url = "http://httpbin.org/get"
			_, err := Get(url, nil, &bin)
			if err != nil {
				b.Fatal(err)
			}
			if len(strings.Split(bin.Origin, ".")) != 4 {
				b.Errorf("want 4, got %d", len(strings.Split(bin.Origin, ".")))
			}
		}
	})
}
