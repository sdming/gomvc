gomvc
=====

smart web &amp; lightweight  web framework in golang

overview
---
I am a new to golang, write your own code is the best way to learn a new programming language, so gomvc comes into the world.

But it's not a nother wheel, the goal of gomvc is to be a smart & fast web framework let you develop web application in KISS way.

my english is not fluent, so I try to keep this document straight and simple.

current version - 0.1

features
---
* develop http handle function in straight way
* data validation
* action method is a normal function, so you write unit test code easily
* support RESTful web api
* fully customizable http handler
* session, cookie, auth, bundle, validation...

road map
---
* 0.1 - basic functionality, make it work
* 0.2 - view enginer
* 0.3 - more result
* 0.4 - maybe validation
* 0.5 - session, cookie, make it better
* 0.6 - performance, make it faster
* 0.7 - bundle
* 0.8 - not planned
* 1.0 - release


quick start
---
```golang
//model
type User struct {
  Id      int
  Name    string
  Age     int
  ZipCode string
}

//repository 
func GetById(id int) User {
  return User{Id: id, Name: "Name " + strconv.Itoa(id), Age: 18, ZipCode: "000000"}
}

//controller
type UserController struct {
}

//get http://localhost:8080/user/user/1000 (application/json; application/xml;text/plain;text/html)
//return struct, marshal to json/xml/text accoring to "Accept"
func (*UserController) User(id int) User {
  return GetById(id)
}

//start server
server := gomvc.DefaultServer()
controller := &user.UserController{}
server.Route("user", "^/user/(?P<action>[A-Za-z]+)(/(?P<p0>[0-9]+))*", "*", controller)
server.Start()
```

example 
---
* basic example: please checkout code: example/user_controller.go
* api test: goto example folder, run user_server.go, and open link: http://localhost:8080/apitest.html
* demo: how to write a customer result(todo:)
* demo: how to write a filter(gzip)(todo:)
* demo: how to return file stream(todo:)

http handle process workflow
---
1. render phsical file if requested file exists
2. match route, call action method
3. rend result

http filter
---
* static filter：render a phsical file 
* render filter: write result to response stream
* mvc filter: match controller route, call action dynamic
* more...

http result
---
* ContentResult: html raw 
* JsonResult: "application/json"
* XmlResult: "application/xml"
* JavaScriptResult: "application/x-javascript"
* JsonpResult (TODO:)
* ViewResult (TODO:)
* PartialViewResult (TODO:)
* FileResult (TODO:)
* FileStreamResult (TODO:)
* RedirectResult (TODO:)
* NotFoundResult (TODO:render 404 template)
* ErrorResult (TODO:render error template)

view engine
---
find view template in fllowing priority:
1. views\controller_name\view_name
2. views\shared\view_name

invoke viw engine by file extension name
* .html - html/template 
* .mustache - mustche template


cache\gzip
---
nginx, haproxy, Varnish can provide awesome service


bundling & minification 
---
todo


change history
---
2012.6.18 init 


License
---------------------

About
---------------------
