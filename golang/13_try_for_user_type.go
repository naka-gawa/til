package main

import "fmt"

type User struct {
	userid string
	gameid int
	gamept int
}

func main() {
	u := User{"tmnakagawa", 10, 100}
	fmt.Println(u)
}
