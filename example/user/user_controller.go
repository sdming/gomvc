package user

import (
	"fmt"
	"strconv"
	"github.com/sdming/gomvc"
)

//controller
type UserController struct {
}

//get http://localhost:8080/user/user/1000 (application/json; application/xml;text/plain;text/html)
//return struct, marshal to json/xml/text accoring to "Accept"
func (*UserController) User(id int) User {
	return GetById(id)
}

//get http://localhost:8080/user/string/1000
//return string
func (*UserController) String(id int) string {
	return GetById(id).String()
}

//get http://localhost:8080/user/int/1000
//return int
func (*UserController) Int(id int) int {
	return GetById(id).Id
}

//get http://localhost:8080/user/json/1000
//return json
func (*UserController) Json(id int) gomvc.HttpResult {
	u := GetById(id)
	return gomvc.Json(u)
}

//get http://localhost:8080/user/xml/1000
//return xml
func (*UserController) Xml(id int) gomvc.HttpResult {
	u := GetById(id)
	return gomvc.Xml(u)
}

//get http://localhost:8080/user/Slice/10
//return slice
func (*UserController) Slice(count int) []User {
	users := Take(count)
	return users
}

//get http://localhost:8080/user/search/123456_10_20
//muti parameter
func (*UserController) Search(zipcode string, ageFrom, ageTo int) []User {
	users := Search(zipcode, ageFrom, ageTo)
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

//GET http://localhost:8080/user/content
//access http context
func (*UserController) Content(ctx *gomvc.HttpContext) string {
	return fmt.Sprintln("request", ctx.RequestPath)
}

//PUT http://localhost:8080/user/Put/1000
//put all them together
func (*UserController) Put(id int, u User, ctx *gomvc.HttpContext) string {
	return fmt.Sprintf("%v %v %v", ctx.Method, id, u.String())
}

//http://localhost:8080/user/error
//raise error
func (*UserController) Error() string {
	n := 0
	i := 100 / n
	return strconv.Itoa(i)
}
