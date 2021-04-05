package main

import "errors"

var success error = errors.New("ok")
var failure error = errors.New("fail")
