package main

import "strconv"

func main(){
	for i := 1; i <= 100; i = i + 1 {
		if i % 2 == 0 {
			println(strconv.Itoa(i) + "-偶数")
		} else {
			println(strconv.Itoa(i) + "-奇数")
		}
	}
}