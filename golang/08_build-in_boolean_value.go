package main

func main() {
	var a, b, c bool
	a = false
	b = false
	c = false
	if a && b || !c {
		println("true")
	} else {
		println("false")
	}
}
