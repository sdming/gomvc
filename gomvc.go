//Copyright 

//go mvc web framework

//

package gomvc

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"path"
	"runtime/debug"
	"strings"
	"time"
	//"io/ioutil"
	"os"
	//"time"
)

// http verbs
const (
	HttpVerbsGet     = "GET"     //http get
	HttpVerbsPost    = "POST"    //http post 
	HttpVerbsPut     = "PUT"     //http put
	HttpVerbsDelete  = "DELETE"  //http delete
	HttpVerbsHead    = "HEAD"    //http head
	HttpVerbsTrace   = "TRACE"   //http trace
	HttpVerbsConnect = "CONNECT" //http connect
	HttpVerbsOptions = "OPTIONS" //http options
)

//content type
const (
	ContentTypeJson = "application/json" //json
	ContentTypeXml  = "application/xml"  //xml
)

//gomvc system messages
const (
	MsgServerTimeout     = "server timeout"        //server execute timeout
	MsgServerInternalErr = "server internal error" //5xx, server internal error
	MsgNotFound          = "404 not found"         //404
)

//gomvc const var
const (
	MvcFilterName    = "mvc"    //default mvc filter name
	MvcDefaultAction = "Index"  // default action name
	MvcActionName    = "action" // default action name
)

//gomvc errors
var (
	ErrViewNotFound   = errors.New("can not find view")    // can not find a view
	ErrInvalidFilters = errors.New("invalid http filters") // invalid http filter
	ErrInvalidRoute   = errors.New("invalid route")        //
)

//gomvc system result
var (
	ResultVoid = &VoidResult{} //action does have return value
)

//global 
var (
	Logger        *log.Logger         //server logger
	EnableDebug   bool        = false //enable debug or not
	EnableProfile bool        = false //enable http package profile or not
)

// http error
type HttpError struct {
	Message string
}

func NewHttpError(format string, args ...interface{}) HttpError {
	return HttpError{fmt.Sprintf(format, args...)}
}

func (self HttpError) Error() string {
	return self.Message
}

//http server
type HttpServer struct {
	Config   WebConfig    //config of server	
	Listener net.Listener //listener
	Filters  []FilterItem //http filters
}

func DefaultServer() *HttpServer {
	return NewHttpServer(WebConfig{})
}

//create a http server 
func NewHttpServer(config WebConfig) *HttpServer {
	s := &HttpServer{Config: config, Filters: []FilterItem{}}
	s.appendFilter("staic", "/", "GET", true, newStaticFilter(s))
	s.appendFilter(MvcFilterName, "/", "*", true, newMvcFilter(s))
	s.appendFilter("render", "/", "*", false, newRenderFilter(s))
	return s
}

//append a filter 
func (s *HttpServer) appendFilter(name, path, method string, p bool, filter HttpFilter) {
	f := FilterItem{Name: name, PathPrefix: path, Method: method, PassOnResult: p, Handler: filter}
	s.Filters = append(s.Filters, f)
}

//insert http filter 
func (s *HttpServer) AddFiler(index int, f FilterItem) {
	s.Filters = append(s.Filters[:index], append([]FilterItem{f}, s.Filters[index:]...)...)
}

//remove filter 
func (s *HttpServer) RemoveFiler(index int) {
	//todo:
}

//add route
func (s *HttpServer) Route(name string, pattern string, method string, controller interface{}) error {
	for _, r := range s.Filters {
		if r.Name == MvcFilterName {
			if mvc, ok := r.Handler.(*MvcFilter); ok {
				return mvc.Route(name, pattern, method, controller)
			}

		}
	}
	return ErrInvalidRoute
}

//start server TODO: run server with goroutine
func (s *HttpServer) Start() error {
	var err error

	if Logger == nil {
		Logger = log.New(os.Stdout, "gomvc", log.Ldate|log.Ltime)
	}
	Logger.Println("gomvc http server start")

	err = s.Config.Check()
	if err != nil {
		Logger.Fatalln("gomvc check config error:", err)
		return err
	}
	Logger.Println("gomvc http server config:")
	Logger.Println(s.Config)

	if s.Filters == nil || len(s.Filters) == 0 {
		Logger.Println("gomvc http server invalid http filters")
		return ErrInvalidFilters
	}

	mux := http.NewServeMux()
	if s.Config.EnableProfile {
		Logger.Println("handle http profile on /debug/pprof")
		mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	}

	if s.Config.Timeout > 0 {
		http.TimeoutHandler(s, time.Duration(s.Config.Timeout)*time.Second, MsgServerTimeout)
	} else {
		mux.Handle("/", s)
	}

	l, err := net.Listen("tcp", s.Config.Address)
	if err != nil {
		Logger.Fatalln("gomvc http server listen error:", err)
		return err
	}

	//TODO: run server with go func(){...}()
	s.Listener = l
	err = http.Serve(s.Listener, mux)
	if err != nil {
		Logger.Fatalln("gomvc http server start error:", err)
		return err
	}

	return nil
}

// close server  TODO: send signal 
func (s *HttpServer) Close() error {

	Logger.Println("gomvc http server is closing")

	if s.Listener != nil {
		s.Listener.Close()
	}

	Logger.Println("gomvc http server closed")
	return nil
}

//build context
func (s *HttpServer) buildContext(w http.ResponseWriter, r *http.Request) *HttpContext {
	_ = r.ParseForm()
	return &HttpContext{
		Resonse:     w,
		Request:     r,
		Method:      r.Method,
		URL:         r.URL,
		RemoteAddr:  r.RemoteAddr,
		RequestPath: cleanPath(strings.TrimSpace(r.URL.Path))}
}

func handleError(err *error) {
	if x := recover(); x != nil {
		var buf bytes.Buffer
		fmt.Fprintln(&buf, "gomvc http server panic handle", x)
		buf.Write(debug.Stack())
		Logger.Fatalln(buf.String())
		*err = NewHttpError("internal error :{0}", x)
	}
}

//execute filter handle
func (s *HttpServer) exeFilter(ctx *HttpContext, f *FilterItem) (err error) {

	if EnableDebug {
		Logger.Println("filter begin:", f.Name)
	}

	defer func() {
		handleError(&err)
	}()

	if f.Match(ctx) {
		if EnableDebug {
			Logger.Println("execute filter:", f.Name)
		}

		f.Handler.Execute(ctx)
	} else {
		if EnableDebug {
			Logger.Println("filter pass:", f.Name)
		}
	}

	return nil
}

//handle http request
func (s *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := s.buildContext(w, r)

	if EnableDebug {
		Logger.Println("request", ctx.Method, ctx.URL, ctx.RemoteAddr)
		Logger.Printf("context:%+v", ctx)
	}
	defer func() {
		if EnableDebug {
			Logger.Println("request end", ctx.Method, ctx.URL, ctx.RemoteAddr)
		}
	}()

	defer func() {
		if x := recover(); x != nil {
			Logger.Fatalln("gomvc ServeHTTP internal error", x)
			s.InternalError(ctx)
		}
	}()

	for _, f := range s.Filters {
		e := s.exeFilter(ctx, &f)
		if e != nil {
			ctx.LastError = e
		}
	}

	if ctx.Result == nil {
		Logger.Println("http result is nil", ctx.Method, ctx.URL, ctx.RemoteAddr)
		s.InternalError(ctx)
	}
}

//server internal error (StatusInternalServerError)
func (s *HttpServer) InternalError(ctx *HttpContext) {
	http.Error(ctx.Resonse, MsgServerInternalErr, http.StatusInternalServerError)
}

//map physical path	 
func (s *HttpServer) MapPath(p string) string {

	f := path.Join(s.Config.PublicPath(), p)
	info, err := os.Stat(f)
	if err != nil || info.IsDir() {
		return ""
	}
	return f
}
