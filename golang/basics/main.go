package main

import (
	"fmt"
	"time"
)

func main() {
	a := 10
	if a == 10 {
		fmt.Println("a is 10")
	} else {
		fmt.Println("a is not 10")
	}
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}
	var i int
	for {
		if i > 2 {
			break
		}
		fmt.Println(i)
		i++
		time.Sleep(1 * time.Second)
	}

loop:
	for i := 0; i < 10; i++ {
		switch i {
		case 5:
			fmt.Printf("i is %d continue\n", i)
			continue
		case 6:
			fmt.Printf("i is %d continue\n", i)
			continue
		case 7:
			break loop
		default:
			fmt.Printf("i is %d\n", i)
		}
	}
	items := []item{
		{price: 10},
		{price: 20},
		{price: 30},
	}
	for i := range items {
		items[i].price += 100
	}
	fmt.Printf("items: %+v\n", items)
}

type item struct {
	price float64
}
