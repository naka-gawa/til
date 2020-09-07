package main

func main() {
	p := struct {
		name string
		age  int
	}{
		name: "Gopher",
		age:  38,
	}
	println(p.name)
	println(p.age)
}
