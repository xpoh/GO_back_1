package main

type Employee struct {
	Name   string  `json:"name" xml:"name" :"name"`
	Age    int     `json:"age" xml:"age" :"age"`
	Salary float32 `json:"salary" xml:"salary" :"salary"`
}
