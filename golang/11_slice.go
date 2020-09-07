package main

import "fmt"

func main() {
	/* 下記よりスマート
	var array [10]int
	ns := array[:3]
	*/
	ns := make([]int, 10, 10)
	fmt.Println(ns)

	/* 下記よりスマート
	var array2 = [...]int{10, 20, 30, 40, 50}
	ms := array2[:]
	*/
	ms := []int{10, 20, 30, 40, 50}
	fmt.Println(ms)
}
