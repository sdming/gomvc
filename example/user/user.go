package user

import (
	"fmt"
	"strconv"
	"encoding/json"
	"encoding/xml"
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

//Xml
func (u User) Xml() string {
	b, err := xml.Marshal(u)
	if err != nil {
		return ""
	}
	return string(b)
}

//Json
func (u User) Json() string {
	b, err := json.Marshal(u)
	if err != nil {
		return ""
	}
	return string(b)
}

//repository 
func GetById(id int) User {
	return User{Id: id, Name: "Name " + strconv.Itoa(id), Age: 18, ZipCode: "000000"}
}

//Take
func Take(count int) []User {
	users := make([]User, count, count)

	for i := 0; i < count; i++ {
		users[i] = GetById(i)
	}
	return users
}

//search
func Search(zipcode string, ageFrom, ageTo int) []User {
	users := make([]User, 2, 2)
	users[0] = User{Id: 1, Name: "Name " + strconv.Itoa(1), Age: ageFrom, ZipCode: zipcode}
	users[1] = User{Id: 2, Name: "Name " + strconv.Itoa(2), Age: ageTo, ZipCode: zipcode}
	return users
}
