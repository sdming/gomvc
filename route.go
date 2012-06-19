//Copyright 

//go mvc web framework

//route

package gomvc

import (
	"reflect"
	"regexp"
	"strings"
	"strconv"
	"io/ioutil"
	"encoding/json"
)

//Mvc Filter
type MvcFilter struct {
	Server *HttpServer //
	Routes []Route
}

//
func newMvcFilter(s *HttpServer) *MvcFilter {
	f := &MvcFilter{Server: s}
	return f
}

//mvc route
type Route struct {
	Name       string        //route name
	Pattern    string        //request path pattern
	Method     string        //http method
	Controller reflect.Value //controller	

	reg   *regexp.Regexp //regular expression 
	names []string       //subexp names

}

// add route
func (mvc *MvcFilter) Route(name string, pattern string, method string, controller interface{}) error {
	if pattern == "" {
		return NewHttpError("gomvc invalid pattern " + pattern)
	}
	if controller == nil {
		return NewHttpError("gomvc invalid controller ")
	}

	var r *regexp.Regexp
	var err error
	if r, err = regexp.Compile(pattern); err != nil {
		return err
	}

	route := Route{Name: name, Pattern: pattern, Method: method}
	route.reg = r
	route.names = r.SubexpNames()[1:]

	if fv, ok := controller.(reflect.Value); ok {
		route.Controller = fv
	} else {
		route.Controller = reflect.ValueOf(controller)
	}

	mvc.Routes = append(mvc.Routes, route)

	return nil
}

//match request
func (route *Route) match(ctx *HttpContext) (data map[string]string, ok bool) {
	reg := route.reg

	if ctx.Method != route.Method && route.Method != "*" && route.Method != "" && !(ctx.Method == "HEAD" && route.Method == "GET") {
		return nil, false
	}

	if !reg.MatchString(ctx.RequestPath) {
		return nil, false
	}

	match := reg.FindStringSubmatch(ctx.RequestPath)

	if len(match[0]) != len(ctx.RequestPath) {
		return nil, false
	}

	data = make(map[string]string)
	for i, m := range match[1:] {
		n := route.names[i]
		if n != "" {
			data[n] = m
		}
	}

	return data, true
}

//match route
func (mvc *MvcFilter) findRoute(ctx *HttpContext) (route Route, data map[string]string, ok bool) {
	ok = false

	if mvc.Routes == nil || len(mvc.Routes) == 0 {
		return
	}

	for _, r := range mvc.Routes {
		data, ok = r.match(ctx)
		if EnableDebug {
			Logger.Println("mvc filter route match test", ok, r.Name, r.Pattern, ctx.RequestPath)
		}
		if ok {
			route = r
			break
		}
	}

	return
}

//find action
func (route *Route) findAction(ctx *HttpContext) reflect.Value {

	action, ok := ctx.RouteData[MvcActionName]
	if !ok || action == "" {
		action = ActionByMethod(ctx.Method)
	}
	action = strings.Title(action)
	method := route.Controller.MethodByName(action)
	if !method.IsValid() {
		action = MvcDefaultAction
		method = route.Controller.MethodByName(action)
	}

	return method
}

//invoke action method
func (route *Route) invoke(ctx *HttpContext) (result interface{}) {

	defer func() {
		if err := recover(); err != nil {
			Logger.Println("go mvc filter invoke fail", err)
			result = &ErrorResult{Data: "go mvc filter invoke fail"}
		}
	}()

	method := route.findAction(ctx)
	if !method.IsValid() {
		return &ErrorResult{Data: "go mvc filter invalid action"}
	}
	methodType := method.Type()

	if EnableDebug {
		Logger.Println("go mvc filter dispatch method", methodType.Name())
	}

	args, ok := decodeArgs(ctx, method)
	if ok == false {
		return &ErrorResult{Data: "invalid method arg"}
	}

	ret, err := callAction(method, args)
	if err != nil {
		Logger.Println("mvc filter call action method fail", methodType.Name(), err)
		return &ErrorResult{Data: "call action method fail"}
	}
	if len(ret) == 0 {
		return ResultVoid
	}
	return ret[0].Interface()
}

//extract args from request
func decodeArgs(ctx *HttpContext, method reflect.Value) (args []reflect.Value, ok bool) {

	methodType := method.Type()
	argsNumber := methodType.NumIn()

	args = make([]reflect.Value, argsNumber, argsNumber)

	for i := 0; i < argsNumber; i++ {
		in := methodType.In(i)

		if EnableDebug {
			Logger.Println("index:", i, "type=", in, "kind=", in.Kind())
		}

		switch in.Kind() {
		case reflect.Ptr:
			if strings.Contains(in.Elem().Name(), "HttpContext") {
				args[i] = reflect.ValueOf(ctx)
			} else {
				Logger.Printf("incoming %v", in)
			}
			//todo:
		case reflect.Struct:
			sv := reflect.New(in)
			if ctx.Method == HttpVerbsPost || ctx.Method == HttpVerbsPut {
				body, err := ioutil.ReadAll(ctx.Request.Body)
				if err != nil {
					Logger.Println("mvc filter read body err", err)
					return nil, false
				}
				err = json.Unmarshal(body, sv.Interface())
				if err != nil {
					Logger.Println("mvc filter unmarshal err", in, err, string(body))
					return nil, false
				}
				sv = reflect.Indirect(sv)
			} else {
				sv = reflect.Indirect(sv)
				for i := 0; i < in.NumField(); i++ {
					f := in.Field(i)
					strv := ctx.Value(f.Name)
					if strv == "" {
						continue
					}

					fValue := sv.FieldByName(f.Name)
					if !fValue.CanSet() {
						continue
					}

					if v, err := ReflectValue(strv, f.Type); err == nil {
						fValue.Set(v)
					} else {
						Logger.Println("mvc filter decode struct err", in, f, err, strv)
						return nil, false
					}
				}
			}
			args[i] = sv
		default:
			// simple type 
			name := "p" + strconv.Itoa(i)
			strv, ok := ctx.RouteData[name]
			if !ok {
				strv = ctx.Request.FormValue(name)
			}
			if v, err := ReflectValue(strv, in); err != nil {
				Logger.Println("mvc filter decode arg err", name, in, err, strv)
				return nil, false
			} else {
				args[i] = v
			}
		}
	}

	return args, true
}

func decodeFuncArgs(ctx *HttpContext, methodType reflect.Type) ([]reflect.Value, bool) {
	argsNumber := methodType.NumIn()
	args := make([]reflect.Value, argsNumber, argsNumber)

	for i := 0; i < argsNumber; i++ {
		in := methodType.In(i)
		name := "p" + strconv.Itoa(i)
		strv, ok := ctx.RouteData[name]
		if !ok {
			strv = ctx.Request.FormValue(name)
		}
		if EnableDebug {
			Logger.Println("in:", i, "name=", name, "kind=", in.Kind(), "value=", strv)
		}

		if v, err := ReflectValue(strv, in); err != nil {
			return nil, false
		} else {
			args[i] = v
		}
	}

	return args, true
}

//mvc filter execute
func (mvc *MvcFilter) Execute(ctx *HttpContext) {

	route, data, ok := mvc.findRoute(ctx)
	if !ok {
		return
	}
	ctx.RouteData = data

	if EnableDebug {
		Logger.Println("mvc filter route matched:", route.Name, data, route.Controller)
	}

	result := route.invoke(ctx)

	if r, ok := result.(HttpResult); ok {
		if EnableDebug {
			Logger.Println("mvc filter result is HttpResult")
		}
		ctx.Result = r
	} else {
		ctx.Result = convertResult(ctx.Accept(), result)
		if EnableDebug {
			Logger.Println("mvc filter convert result to ", reflect.TypeOf(ctx.Result))
		}
	}

	if EnableDebug {
		Logger.Println("mvc filter set context result to", ctx.Result)
	}
}
