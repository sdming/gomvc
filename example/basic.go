package main

import (
	"fmt"
	"gomvc"
	"strconv"
	"reflect"
	"errors"
)

//model
type User struct {
	Id      int
	Name    string
	Age     int
	ZipCode string
}

//Stringer
func (u User) String() string {
	return fmt.Sprint("Id:", u.Id, "; Name:", u.Name, "; Age:", u.Age, "; ZipCode:", u.ZipCode)
}

//controller
type UserController struct {
}

func GetUser(id int) User {
	return User{Id: id, Name: "Name " + strconv.Itoa(id), Age: 18, ZipCode: "000000"}
}

//get http://localhost:8080/user/user/1000 (application/json; application/xml;text/plain;text/html)
//return struct, marshal to json/xml/text accoring to "Accept"
func (*UserController) User(id int) User {
	return GetUser(id)
}

//get http://localhost:8080/user/string/1000
//return string
func (*UserController) String(id int) string {
	return GetUser(id).String()
}

//get http://localhost:8080/user/int/1000
//return int
func (*UserController) Int(id int) int {
	return GetUser(id).Id
}

//get http://localhost:8080/user/json/1000
//return json
func (*UserController) Json(id int) gomvc.HttpResult {
	u := GetUser(id)
	return gomvc.Json(u)
}

//get http://localhost:8080/user/xml/1000
//return xml
func (*UserController) Xml(id int) gomvc.HttpResult {
	u := GetUser(id)
	return gomvc.Xml(u)
}

//get http://localhost:8080/user/slice
//return slice
func (*UserController) Slice() []User {
	users := make([]User, 0, 0)

	for i := 0; i < 2; i++ {
		users = append(users, GetUser(i))
	}
	return users
}

//get http://localhost:8080/user/search/123456_10_20
//muti parameter
func (*UserController) Search(zipcode string, ageFrom, ageTo int) []User {
	users := make([]User, 0, 0)
	users = append(users, User{Id: 1, Name: "Name " + strconv.Itoa(1), Age: ageFrom, ZipCode: zipcode})
	users = append(users, User{Id: 2, Name: "Name " + strconv.Itoa(2), Age: ageTo, ZipCode: zipcode})
	return users
}

//get http://localhost:8080/user/struct?Id=1000&Name=hello&Zipcode=000000&Age=18
//struct as parameter
//validation(TODO):required,pattern,type,max,min,range,maxLength,minLength, rangeLength,
//  number,date,time,zipcode,alphanumeric,lettersonly,email,url,greaterThan,lessThan
func (*UserController) Struct(p struct {
	Id      int    `required, min:0`
	Name    string `required, rangeLength:1-10`
	Age     int    `default:1, range:0-99`
	Zipcode string `pattern:[0-9]+`
},) User {
	return User{Id: p.Id, Name: p.Name, Age: p.Age, ZipCode: p.Zipcode}
}

//post http://localhost:8080/user/post
//unmarshal parameter from json
func (*UserController) Post(u User) string {
	return u.String()
}

//PUT http://localhost:8080/user/content
//access http context
func (*UserController) Content(ctx *gomvc.HttpContext) string {
	return fmt.Sprintln("request ", ctx.RequestPath)
}

//PUT http://localhost:8080/user/Put/1000
//put all them together
func (*UserController) Put(id int, u User, ctx *gomvc.HttpContext) string {
	return fmt.Sprintf("%v %v %#v from %v", tx.Method, id, u)
}

//http://localhost:8080/user/error
//raise error
func (*UserController) Error() string {
	n := 0
	i := 100 / n
	return strconv.Itoa(i)
}

//start web server
func StartWeb() {

	server := gomvc.DefaultServer()
	server.Route("user", "^/user/(?P<action>[A-Za-z]+)(/(?P<p0>[0-9]+))*", "*", &UserController{})
	server.Route("user", "^/user/(?P<action>[A-Za-z]+)/(?P<p0>[0-9]+)_(?P<p1>[0-9]+)_(?P<p2>[0-9]+)", "GET", &UserController{})
	server.Start()

}

func main() {
	StartWeb()
}
