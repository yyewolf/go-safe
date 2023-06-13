package main

type File struct {
	Sum string `json:"s"`
}

var database map[string]*File
