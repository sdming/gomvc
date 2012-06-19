package main

import (
	"fmt"
	"gomvc"
	"strconv"
	"math/rand"
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

func (u User) String() string {
	return fmt.Sprint("Id:", u.Id, "; Name:Name ", u.Id, "; Age:", u.Age, "; ZipCode:", u.ZipCode)
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

	for i := 0; i < 10; i++ {
		users = append(users, GetUser(i))
	}
	return users
}

//get http://localhost:8080/user/search/123456_10_20
//muti parameter
func (*UserController) Search(zipcode string, ageFrom, ageTo int) []User {
	users := make([]User, 0, 0)

	for i := 0; i < 10; i++ {
		age := rand.Intn(ageTo-ageFrom) + ageFrom
		u := User{Id: i, Name: "Name " + strconv.Itoa(i), Age: age, ZipCode: zipcode}
		users = append(users, u)
	}
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
	return fmt.Sprintln("request from ", ctx.RemoteAddr)
}

//PUT http://localhost:8080/user/Put/1000
//put all them together
func (*UserController) Put(id int, u User, ctx *gomvc.HttpContext) string {
	return fmt.Sprintf("put %v as %#v from %v", id, u, ctx.RemoteAddr)
}

//http://localhost:8080/user/error
//raise error
func (*UserController) Error() string {
	n := 0
	i := 100 / n

	return strconv.Itoa(i)
}

//user content
//return as a file

//start web server
func StartWeb() {

	server := gomvc.DefaultServer()
	server.Route("user", "^/user/(?P<action>[A-Za-z]+)(/(?P<p0>[0-9]+))*", "*", &UserController{})
	server.Route("user", "^/user/(?P<action>[A-Za-z]+)/(?P<p0>[0-9]+)_(?P<p1>[0-9]+)_(?P<p2>[0-9]+)", "GET", &UserController{})

	server.Start()

	//http://localhost:8080/user/detail/1000
}

func main() {
	StartWeb()
	//TestReflect()
}

type TestA struct {
	Id      int
	Name    string
	Age     int
	Zipcode string
}

func TestReflect() {
	t := TestA{}
	tv := reflect.ValueOf(t)

	fmt.Println("t", t)

	typ := tv.Type()
	fmt.Println("typ", typ)

	arg := reflect.New(typ)
	fmt.Println("arg", arg)

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)

		name := f.Name
		fmt.Println("filed", i, f.Type, name)

		strv := "100"
		fff := reflect.Indirect(arg).FieldByName(name)
		if fff.CanSet() {
			if fv, err := ReflectValue(strv, f.Type); err == nil {
				fmt.Println(name, "get value", fv)
				fff.Set(fv)
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println("can not set value ", name)
		}

	}
	fmt.Println("arg", arg.Interface())
}

//convert string to reflect.Value according to type
func ReflectValue(str string, typ reflect.Type) (val reflect.Value, err error) {

	switch typ.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(str)
		return reflect.ValueOf(b), err
	case reflect.Int:
		i, err := strconv.Atoi(str)
		return reflect.ValueOf(i), err
	case reflect.Int8:
		i64, err := strconv.ParseInt(str, 10, 8)
		return reflect.ValueOf(int8(i64)), err
	case reflect.Int16:
		i64, err := strconv.ParseInt(str, 10, 16)
		return reflect.ValueOf(int16(i64)), err
	case reflect.Int32:
		i64, err := strconv.ParseInt(str, 10, 32)
		return reflect.ValueOf(int32(i64)), err
	case reflect.Int64:
		i64, err := strconv.ParseInt(str, 10, 32)
		return reflect.ValueOf(i64), err
	case reflect.Uint:
		u64, err := strconv.ParseUint(str, 10, 32)
		return reflect.ValueOf(uint(u64)), err
	case reflect.Uint8:
		u64, err := strconv.ParseUint(str, 10, 8)
		return reflect.ValueOf(uint8(u64)), err
	case reflect.Uint16:
		u64, err := strconv.ParseUint(str, 10, 16)
		return reflect.ValueOf(uint16(u64)), err
	case reflect.Uint32:
		u64, err := strconv.ParseUint(str, 10, 32)
		return reflect.ValueOf(uint32(u64)), err
	case reflect.Uint64:
		u64, err := strconv.ParseUint(str, 10, 64)
		return reflect.ValueOf(u64), err
	case reflect.Float32:
		f64, err := strconv.ParseFloat(str, 32)
		return reflect.ValueOf(float32(f64)), err
	case reflect.Float64:
		f64, err := strconv.ParseFloat(str, 32)
		return reflect.ValueOf(f64), err
	case reflect.String:
		return reflect.ValueOf(str), nil
	case reflect.Complex64, reflect.Complex128:
		err = errors.New("todo")
	case reflect.Array, reflect.Slice:
		err = errors.New("todo")
	case reflect.Map:
		err = errors.New("todo")
	case reflect.Struct:
		err = errors.New("todo")
	case reflect.Uintptr, reflect.Ptr, reflect.Chan, reflect.Func, reflect.UnsafePointer:
		err = errors.New("do not support")
	default:
		err = errors.New("unknownType")
	}

	return reflect.Zero(typ), err

}
