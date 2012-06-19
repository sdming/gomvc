//Copyright 

//go mvc web framework

//http result

package gomvc

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

//execute error result
func (ctx *HttpContext) ExeError(err string) {
	e := ErrorResult{Data: err}
	e.Execute(ctx)
}

//http result interface
type HttpResult interface {
	Execute(ctx *HttpContext)
}

//static file 
type FileResult struct {
	Data string //file path
}

//render file
func (file *FileResult) Execute(ctx *HttpContext) {
	http.ServeFile(ctx.Resonse, ctx.Request, file.Data)
}

//http raw content
type ContentResult struct {
	ModifyTime  time.Time   //modify time
	ContentType string      //content type
	Data        interface{} //content
}

//render file
func (c *ContentResult) Execute(ctx *HttpContext) {
	if r, ok := c.Data.(io.ReadSeeker); ok {
		http.ServeContent(ctx.Resonse, ctx.Request, c.ContentType, c.ModifyTime, r)
		return
	}
	fmt.Fprintln(ctx.Resonse, c.Data)
	//http.ServeContent(ctx.Resonse, ctx.Request, c.Name, c.ModifyTime, c.Data)
}

//javascript ContentType = "application/x-javascript";
type JavaScriptResult struct {
	Data interface{} //data
}

//render javascript
func (j *JavaScriptResult) Execute(ctx *HttpContext) {
	//TODO
}

//json ContentType = "application/json"
type JsonResult struct {
	Data interface{} //data
}

//render json
func (j *JsonResult) Execute(ctx *HttpContext) {
	b, err := json.Marshal(j.Data)
	if err != nil {
		ctx.ExeError(err.Error())
		return
	}

	ctx.ContentType(ContentTypeJson)
	ctx.Write(b)
}

//xml ContentType = "application/xml"
type XmlResult struct {
	Data interface{} //data
}

//render xml
func (x *XmlResult) Execute(ctx *HttpContext) {
	b, err := xml.Marshal(x.Data)
	if err != nil {
		ctx.ExeError(err.Error())
		return
	}

	ctx.ContentType(ContentTypeXml)
	ctx.Write(b)

}

//jsonp
type JsonpResult struct {
	Data interface{} //data
}

//render jsonp
func (j *JsonpResult) Execute(ctx *HttpContext) {
	//TODO
}

//view
type ViewResult struct {
}

//Partial View
type PartialViewResult struct {
}

//redirect 
type RedirectResult struct {
	Permanent bool //permanent redirect or not
	Path      string
}

//http 404 TODO
type NotFoundResult struct {
}

//TODO: render 404.xxx view page
func (r *NotFoundResult) Execute(ctx *HttpContext) {
	http.Error(ctx.Resonse, MsgNotFound, http.StatusNotFound)
}

//Error
type ErrorResult struct {
	Data string
}

//TODO: render error.xxx view page
func (r *ErrorResult) Execute(ctx *HttpContext) {
	http.Error(ctx.Resonse, r.Data, http.StatusInternalServerError)
}

//VoidResult
type VoidResult struct {
}

func (r *VoidResult) Execute(ctx *HttpContext) {
	ctx.Resonse.Write([]byte(``))
}
