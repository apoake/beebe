package controller

import  macaron "gopkg.in/macaron.v1"

var m *macaron.Macaron

func init() {
	m = macaron.Classic()
}

func Macaron() *macaron.Macaron {
	return m
}