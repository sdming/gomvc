//Copyright 

//go mvc web framework

//http content 

package gomvc

import (
	"net/http"
	"net/url"
)

//http context
type HttpContext struct {
	Request *http.Request //http request
	//Agent        string              //user agent
	Resonse      http.ResponseWriter //http response
	Method       string              //http method
	RequestPath  string              //request path
	PhysicalPath string              //physical Path
	URL          *url.URL            //request url
	RemoteAddr   string              //remote address

	//self fields
	RouteData map[string]string      //route data
	Result    HttpResult             //result
	LastError error                  //last error
	User      string                 //user name
	Variables map[string]interface{} //server variables
	Files     string                 //files uploaded, TODO:
	//Session, Cookie, Form, QueryString, Cache//TODO: 
}

func (ctx *HttpContext) Value(name string) string {
	v, ok := ctx.RouteData[name]
	if ok {
		return v
	}
	return ctx.Request.FormValue(name)
}

func (ctx *HttpContext) UserAgent() string {
	return ctx.Request.Header.Get("User-Agent")
}

func (ctx *HttpContext) SetHeader(key string, value string) {
	ctx.Resonse.Header().Set(key, value)
}

func (ctx *HttpContext) ContentType(ctype string) {
	ctx.Resonse.Header().Set("Content-Type", ctype)
}

func (ctx *HttpContext) Status(code int) {
	ctx.Resonse.WriteHeader(code)
}

func (ctx *HttpContext) Accept() string {
	return ctx.Request.Header.Get("Accept")
}

func (ctx *HttpContext) Write(b []byte) {
	ctx.Resonse.Write(b)
}
