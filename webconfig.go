//Copyright 

//go mvc web framework

//mvc server config

package gomvc

import (
	//"fmt"
	"os"
	"path"
)

const (
	DefaultAddress string = "0.0.0.0:8080" //default listen address
	pubicDir              = "public"
)

//mvc server config
type WebConfig struct {
	ServerKey     string //TODO: used to session & cache
	Address       string //listen address, default is 0.0.0.0:8000
	RootDir       string //root dir, default is current work path
	EnableProfile bool   //enable http profile or not	
	Timeout       int    //server execute time out (in second)
	publicDir     string //static files path
}

//static files path
func (w *WebConfig) PublicPath() string {
	return path.Join(w.RootDir, "public")
}

//check config
func (w *WebConfig) Check() (err error) {
	if w.Address == "" {
		w.Address = DefaultAddress
	}

	if w.RootDir == "" {
		wd, _ := os.Getwd()
		w.RootDir = wd
	}

	f := path.Join(w.RootDir, pubicDir)
	info, err := os.Stat(f)
	if err == nil && info.IsDir() {
		w.publicDir = f
	} else {
		w.publicDir = w.RootDir
	}
	return nil
}

//TODO: create config from file
func (w WebConfig) FromFile(file string) (c *WebConfig, err error) {
	panic("not implemented")
}

//TODO: create config from json string
func (w WebConfig) FromJson(s string) (c *WebConfig, err error) {
	panic("not implemented")
}

//TODO: create config from xml
func (w WebConfig) FromXml(s string) (c *WebConfig, err error) {
	panic("not implemented")
}
