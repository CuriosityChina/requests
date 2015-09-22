package requests

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	// var x struct {
	// 	Code int    `json:"code"`
	// 	Msg  string `json:"msg"`
	// }
	var x interface{}
	// var url = "http://www.baidu.com"
	var url = "http://192.168.9.5:8081/v1/wx/1f6b2a478d4d4678276944b0a607058b/auth/accesstoken"
	result, err := Get(url, nil, &x)
	fmt.Println(x, err)
	m, _ := ConvertResponseToBytes(result)
	fmt.Println(string(m))
}

func TestPost(t *testing.T) {

}

func TestDelete(t *testing.T) {

}

func TestPut(t *testing.T) {

}
