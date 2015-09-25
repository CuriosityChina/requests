#### **Requests**
-----
Reuqests is a library for working for http request. It currently supports
  > - HTTP request (GET, POST, PUT, POST)
  > - Automatic Marshal and Unmarshal data
  > - Support library internal debug logging



#### **Code Examples**
-----
``` go
package main

import (
	"log"

	"git.curio.im/golibs/requests"
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

func main() {
	var bin httpBin
	var url = "http://httpbin.org/get"
	result, err := requests.Get(url, nil, &bin)
    defer result.Body.Close()
	if err != nil {
		log.Println(err)
	}
	log.Printf("IP: %s", bin.Origin)
	log.Println(result.Status)
}

```

#### **Known issues**
-----
 - Unsupports Post Form
 - Unsupoorts non-utf8 encoding XML
