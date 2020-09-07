package main

import (
	"math/rand"
	"time"
)

func main() {
	/*
	完全なランダムな変数を作ることはできない
	なので math/rand に seed として EPOCH TIME を渡しランダム変数を作成する 
	*/
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(6)
	println(n)
	switch n {
	case 6:
		println("大吉")
	case 5, 4:
		println("中吉")
	case 3, 2:
		println("小吉")
	case 1:
		println("凶")
	case 0:
		println("お祈り")
	}
}