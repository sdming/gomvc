//Copyright 

//go mvc web framework

//filter

package gomvc

import (
	"strings"
)

//http filter
type FilterItem struct {
	Name         string     //name of filter	
	PathPrefix   string     //path match
	Method       string     //http method
	PassOnResult bool       //pass if result exists
	Handler      HttpFilter //handler
}

func (f *FilterItem) Match(ctx *HttpContext) bool {
	if !strings.HasPrefix(ctx.RequestPath, f.PathPrefix) {
		return false
	}
	if f.PassOnResult && ctx.Result != nil {
		return false
	}
	if f.Method == "*" || f.Method == "" || strings.Contains(f.Method, ctx.Method) {
		return true
	}

	return false
}

//filter handler interface
type HttpFilter interface {
	Execute(ctx *HttpContext)
}

//wrap of func
type HandlerFunc func(ctx *HttpContext)

//wrap of func execute
func (f HandlerFunc) Execute(ctx *HttpContext) {
	f(ctx)
}

/*
//http filter base
type FilterBase struct {
	PassOnError  bool        //
	PassOnResult bool        //
	Server       *HttpServer //
}
*/

//static file fileter
type StaticFilter struct {
	Server *HttpServer //
}

//
func newStaticFilter(s *HttpServer) *StaticFilter {
	f := &StaticFilter{Server: s}
	return f
}

//static file fileter execute
func (filter *StaticFilter) Execute(ctx *HttpContext) {

	if EnableDebug {
		Logger.Println("static file filter, request path is", ctx.RequestPath)
	}

	ctx.PhysicalPath = filter.Server.MapPath(ctx.RequestPath)

	if EnableDebug {
		Logger.Println("static file filter, physical path is", ctx.PhysicalPath)
	}
	if ctx.PhysicalPath != "" {
		ctx.Result = &FileResult{Data: ctx.PhysicalPath}
	}
}

//render result filter
type RenderFilter struct {
	Server *HttpServer //
}

//
func newRenderFilter(s *HttpServer) *RenderFilter {
	f := &RenderFilter{Server: s}
	return f
}

//render result filter execute
func (filter *RenderFilter) Execute(ctx *HttpContext) {

	if EnableDebug {
		Logger.Println("render filter, result:", ctx.Result)
	}

	if ctx.Result == nil {
		ctx.Result = new(NotFoundResult)
	}

	ctx.Result.Execute(ctx)
}
