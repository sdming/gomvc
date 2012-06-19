package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"bytes"
	"log"
)

func do(url, method, accept, body string) {
	fmt.Println(method, url, accept)

	b, err := curl(url, method, accept, body)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf(" %s \n\n", b)
	}
}

func curl(url, method, accept, body string) ([]byte, error) {

	buf := bytes.NewBuffer([]byte(body))
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		log.Println("http.NewRequest", err)
		return nil, err
	}

	if accept != "" {
		req.Header.Set("Accept", accept)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println("client.do error", err)
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("ioutil.ReadAll", err)
		return nil, err
	}

	res.Body.Close()

	return b, nil

}

func main() {
	do("http://localhost:8080/user/xml/1000", "GET", "", "")
	do("http://localhost:8080/user/json/1000", "GET", "", "")
	do("http://localhost:8080/user/user/1000", "GET", "", "")
	do("http://localhost:8080/user/string/1000", "GET", "", "")
	do("http://localhost:8080/user/int/1000", "GET", "", "")
	do("http://localhost:8080/user/user/1000", "GET", "application/json", "")
	do("http://localhost:8080/user/user/1000", "GET", "application/xml", "")
	do("http://localhost:8080/user/user/1000", "GET", "", "")
	do("http://localhost:8080/user/slice", "GET", "", "")
	do("http://localhost:8080/user/search/123456_10_20", "GET", "application/json", "")
	do("http://localhost:8080/user/error", "GET", "", "")
	do("http://localhost:8080/user/struct?Id=1000&Name=hello&Zipcode=000000&Age=18", "GET", "", "")
}
