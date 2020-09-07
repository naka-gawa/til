package main

import "strconv"

func main() {
	for i := 1; i < 100; i = i + 1 {
		n := i % 2
		switch n {
		case 0 :
			println(strconv.Itoa(i) + "-偶数")
		case 1 :
			println(strconv.Itoa(i) + "-奇数")
		}
		
	}
}