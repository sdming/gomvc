package main

import (
	"fmt"
	"github.com/sdming/gomvc"
	"github.com/sdming/gomvc/example/user"
)

func main() {
	fmt.Println("user web server is starting...")
	server := gomvc.DefaultServer()
	controller := &user.UserController{}
	server.Route("user", "^/user/(?P<action>[A-Za-z]+)(/(?P<p0>[0-9]+))*", "*", controller)
	server.Route("user", "^/user/(?P<action>[A-Za-z]+)/(?P<p0>[0-9]+)_(?P<p1>[0-9]+)_(?P<p2>[0-9]+)", "GET", controller)
	server.Start()
}
