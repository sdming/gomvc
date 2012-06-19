//Copyright 

//go mvc web framework

//utils 

package gomvc

import (
    "errors"
    "path"
    "reflect"
    "strconv"
    "strings"
)

func ActionByMethod(method string) string {
    switch method {
    case HttpVerbsConnect, HttpVerbsOptions, HttpVerbsTrace:
        return "NotSupported"
    case HttpVerbsGet:
        return "Get"
    case HttpVerbsDelete:
        return "Delete"
    case HttpVerbsHead:
        return "Get"
    case HttpVerbsPost:
        return "Post"
    }

    return "NotSupported"
}

func IsWebPrintType(kind reflect.Kind) bool {
    switch kind {
    case reflect.Uintptr, reflect.Ptr, reflect.Chan, reflect.Func, reflect.UnsafePointer:
        return false
    }
    return true
}

func IsStructureType(kind reflect.Kind) bool {
    switch kind {
    case reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
        return true
    }
    return false
}

//is simple type or not
func IsSimpleType(kind reflect.Kind) bool {

    switch kind {
    case reflect.Bool:
        return true
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return true
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        return true
    case reflect.Float32, reflect.Float64:
        return true
    case reflect.String:
        return true
    case reflect.Complex64, reflect.Complex128:
        return true
    }
    return false
}

//get method by name
func PublicMethodByName(v reflect.Value, name string) (method reflect.Value, ok bool) {
    ok = false

    if name == "" {
        return
    }

    method = v.MethodByName(strings.Title(name))
    if method.Kind() != reflect.Func {
        return
    }
    return method, ok
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

func convertResult(accept string, value interface{}) HttpResult {

    var (
        v   interface{}
        rv  reflect.Value
        ok  bool
    )

    if rv, ok = value.(reflect.Value); ok {
        v = rv.Interface()
    } else {
        v = value
        rv = reflect.ValueOf(value)
    }

    kind := rv.Kind()

    switch {
    case IsSimpleType(kind):
        return &ContentResult{Data: v}
    case !IsWebPrintType(kind):
        return ResultVoid
    case IsStructureType(kind):
        switch {
        case strings.Index(accept, "application/json") > -1:
            return &JsonResult{Data: v}
        case strings.Index(accept, "application/xml") > -1:
            return &XmlResult{Data: v}
        case strings.Index(accept, "application/jsonp") > -1:
            return &JsonpResult{Data: v}
        case strings.Index(accept, "application/javascript") > -1:
            return &JavaScriptResult{Data: v}
        }
    }
    return &ContentResult{Data: v}
}

//extract args from request
func DecodeStruct(typ reflect.Type, value func(string) string) (reflect.Value, bool) {
    sv := reflect.New(typ)
    sv = reflect.Indirect(sv)

    //if typ.Kind() == reflect.Ptr {
    //typ = typ.Elem()  
    //if v.IsNil() { v.Set(reflect.New(t)) }; v = v.Elem()
    //}

    for i := 0; i < typ.NumField(); i++ {
        f := typ.Field(i)
        name := f.Name
        strv := value(name)
        if strv == "" {
            continue
        }

        fValue := sv.FieldByName(name)
        if !fValue.CanSet() {
            continue
        }

        if v, err := ReflectValue(strv, f.Type); err == nil {
            fValue.Set(v)
        } else {
            return sv, false
        }
    }

    return sv, true
}

//call action
func callAction(function reflect.Value, args []reflect.Value) (result []reflect.Value, err interface{}) {
    defer func() {
        if err = recover(); err != nil {
            result = nil
        }

    }()
    return function.Call(args), nil
}

func Json(a interface{}) *JsonResult {
    return &JsonResult{Data: a}
}

func Xml(a interface{}) *XmlResult {
    return &XmlResult{Data: a}
}

// Return the canonical path for p, eliminating . and .. elements.
func cleanPath(p string) string {
    if p == "" {
        return "/"
    }
    if p[0] != '/' {
        p = "/" + p
    }
    np := path.Clean(p)
    // path.Clean removes trailing slash except for root;
    // put the trailing slash back if necessary.
    if p[len(p)-1] == '/' && np != "/" {
        np += "/"
    }
    return np
}
