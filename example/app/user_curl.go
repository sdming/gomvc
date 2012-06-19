package main

import (
	"fmt"
	"strings"
	"github.com/sdming/gomvc/example/lib"
	"github.com/sdming/gomvc/example/user"
)

var (
	host    string = "http://localhost:8080"
	verbose bool   = true
)

func Test(url, method, accept, body, expect string) {
	fmt.Println("")
	fmt.Println(method, url, accept)
	b, err := lib.Curl(host+url, method, accept, body)
	if err != nil {
		fmt.Println("fail", err)
		return
	}

	res := fmt.Sprintf("%s", b)

	if !strings.Contains(res, expect) {
		fmt.Println("fail", "want", expect, "reponse", res)
		return
	}
	fmt.Println("pass")
	if verbose {
		fmt.Println(res)
	}
}

func main() {

	Test("/user/xml/1000", "GET", "", "", user.GetById(1000).Xml())
	Test("/user/json/1000", "GET", "", "", user.GetById(1000).Json())
	Test("/user/user/1000", "GET", "", "", user.GetById(1000).String())
	Test("/user/string/1000", "GET", "", "", user.GetById(1000).String())
	Test("/user/int/1000", "GET", "", "", "1000")
	Test("/user/user/1000", "GET", "application/json", "", user.GetById(1000).Json())
	Test("/user/user/1000", "GET", "application/xml", "", user.GetById(1000).Xml())
	Test("/user/user/1000", "GET", "text/plain", "", user.GetById(1000).String())
	Test("/user/user/error", "GET", "", "", "")
	Test("/user/slice/10", "GET", "application/json", "", lib.Json(user.Take(10)))
	Test("/user/search/123456_10_20", "GET", "application/json", "", lib.Json(user.Search("123456", 10, 20)))
	Test("/user/struct?Id=1000&Name=hello&Zipcode=000000&Age=18", "GET", "", "", user.User{Id: 1000, Name: "hello", Age: 18, ZipCode: "000000"}.String())
	Test("/user/content", "GET", "", "", "request /user/content")
	Test("/user/post", "POST", "", `{"Id":1000,"Name":"Name 1000","Age":18,"ZipCode":"000000"}`, "Id:1000; Name:Name 1000; Age:18; ZipCode:000000")
	Test("/user/Put/1000", "PUT", "", `{"Id":1000,"Name":"Name 1000","Age":18,"ZipCode":"000000"}`, "PUT 1000 Id:1000; Name:Name 1000; Age:18; ZipCode:000000")

}
