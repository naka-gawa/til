package main
func main() {
	/* 
	avg では少数が入ることが想定される
	なので avg は float とする
	ただし、計算元となる sum や ３ は int型なので
	float にキャストしてあげる必要がある
    */
	var sum int
	var avg float32
	sum = 5 + 6 + 3
	avg = float32(sum) / float32(3)
	// avg := sum / 3
	if avg > 4.5 {
		println("good")
	}
}