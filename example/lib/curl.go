package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
	"encoding/xml"
)

//Xml
func Xml(value interface{}) string {
	b, err := xml.Marshal(value)
	if err != nil {
		return ""
	}
	return string(b)
}

//Json
func Json(value interface{}) string {
	b, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(b)
}

//String
func String(value interface{}) string {
	return fmt.Sprintln(value)
}

func DoCurl(url, method, accept, body string) {
	fmt.Println(method, url, accept)

	b, err := Curl(url, method, accept, body)
	if err != nil {
		fmt.Println("fail", method, url, err)
	} else {
		fmt.Printf(" %s \n\n", b)
	}
}

func Curl(url, method, accept, body string) ([]byte, error) {

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
