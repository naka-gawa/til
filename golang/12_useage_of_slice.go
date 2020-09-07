package main

func main() {
	/*
		n1 := 19
		n2 := 86
		n3 := 1
		n4 := 12
	*/
	ns := []int{19, 86, 1, 12}
	sum := 0
	for _, v := range ns {
		sum = sum + v
	}
	println(sum)
}
