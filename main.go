package main

import (
	"github.com/niwho/hellox/config"
	"github.com/niwho/hellox/http"
)

var COMMITVER string
var DATE string
var GOVERSION string

func main() {
	config.ParseCommandParams(COMMITVER, DATE, GOVERSION)
	http.InitHttp()
}
