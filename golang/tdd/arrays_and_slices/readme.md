## Summary
社内 Go もくもく会
テスト駆動開発でGO言語を学びましょうの[arrays and slices編](https://andmorefine.gitbook.io/learn-go-with-tests/go-fundamentals/arrays-and-slices)

## Array
### 配列の要素数分繰り返すためのテストコードを書く

```golang
package main

import "testing"

func TestSum(t *Testing.T) {
    numbers := [5]int{1,2,3,4,5}
    got := Sum(numbers)
    want := 15

    if got := want {
        t.Error("got %d given %d, %v", got, want, numbers)
    }
}
```

### 最小コード
- 配列を受け取るが何が何でも 0 を返すようなもの

```golang
package main

func Sum(numbers [5]int) int {
    return 0
}
```
- 想定どおり、test は fail する

```
─ go test
--- FAIL: TestSum (0.00s)
    sum_test.go:11: got 0 want 15 given, [1 2 3 4 5]
FAIL
exit status 1
FAIL    arrays_and_slices       0.500s
```

- 成功させるためには、引数を受け取りその分　`for` で loop して加算する必要がある
　- *いつも忘れるのでメモ* for 文は, 初期値, 継続条件, 終了時動作である

```golang
package main

func Sum(numbers [5]int) int {
    for i := 0; i < 5; i++ {
        sum += numbers[i]
    }
    return sum
}
```

### リファクタリング
- for の条件を range で要素数を取得するようにする
  - range で取得すると、添字と値が取得できる
  - 今回は添字は何も使わないので `_` で廃棄している

```golang
func Sum(numbers [5]int) int {
    for _, number := range numbers {
        sum += number
    }
    return sum
}
```

- もちろんテスト結果は変わらない

```
╰─ go test
PASS
ok      arrays_and_slices       0.473s
```

## Slice
### スライスの要素数分繰り返すためのテストコードを書く

```golang
package main

import "testing"

func TestSum(t *testing.T) {

	t.Run("collection of 5 numbers", func(t *testing.T) {
			numbers := [5]int{1, 2, 3, 4, 5}

			got := Sum(numbers)
			want := 15

			if got != want {
					t.Errorf("got %d want %d given, %v", got, want, numbers)
			}
	})

	t.Run("collection of any size", func(t *testing.T) {
			numbers := []int{1, 2, 3}

			got := Sum(numbers)
			want := 6

			if got != want {
					t.Errorf("got %d want %d given, %v", got, want, numbers)
			}
	})
}
```

- test は fail する
  - 最初の  numbers では、配列を渡しているが、関数側では配列ではなく、スライスを受け取るようになっているため

```
╰─ go test
# arrays_and_slices [arrays_and_slices.test]
./sum_test.go:10:14: cannot use numbers (variable of type [5]int) as []int value in argument to Sum
FAIL    arrays_and_slices [build failed]
```

- 下記のように修正

```golang
package main

import "testing"

func TestSum(t *testing.T) {

	t.Run("collection of 5 numbers", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}

		got := Sum(numbers)
		want := 15

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

	t.Run("collection of any size", func(t *testing.T) {
		numbers := []int{1, 2, 3}

		got := Sum(numbers)
		want := 6

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})
}
```
- test も pass になった

```
╰─ go test
PASS
ok      arrays_and_slices       0.518s
```

## slice sum
### テストコード

``` golang
func TestSumAll(t *testing.T) {

    got := SumAll([]int{1, 2, 3}, []int{0, 9, 18})
    want := []int{6, 27}

    if got != want {
        t.Errorf("got %v want %v", got, want)
    }
}
```

```golang
func SumAll(numbersToSum ...[]int) (sums []int) {
    return
}
```

- やりたいことはスライスの中を加算して、新たなスライスを作成すること
  - Go では統合演算子をスライスで使うことができない
    - 繰り返し処理で頑張ることもできるが、reflect.DeepEqual を使うと同じ変数であるかわかるのでコレを使う

```golang
func TestSumAll(t *testing.T) {

    got := SumAll([]int{1, 2, 3}, []int{0, 9, 18})
    want := []int{6, 27}

    if !reflect.DeepEqual(got, want) {
        t.Errorf("got %v want %v", got, want)
    }
}
```

#### reflect.DeepEqual が安全ではない理由

- 仮に want を `string` に変えてしまったとしても、コンパイルされてしまい、実行されちゃう
  - [Time型が入った構造体などをDeepEqualするとfalseになる可能性がある](https://qiita.com/weloan/items/fb1cda8407022a52316f)
  - [go-cmp](https://github.com/google/go-cmp)というライブラリがあるらしい

```golang
func TestSumAll(t *testing.T) {
	got := SumAll([]int{1, 2, 3}, []int{0, 9, 18})
	want := "bob"
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
```

```
╰─ go test
--- FAIL: TestSumAll (0.00s)
    sum_test.go:37: got [6 27] want bob
FAIL
exit status 1
FAIL    arrays_and_slices       0.382s
```

### テストが通るようにする

```golang
func SumAll(numbersToSum ...[]int) []int {
	lengthOfNumbers := len(numbersToSum)
	sums := make([]int, lengthOfNumbers)

	for i, numbers := range numbersToSum {
		sum[i] = Sum(numbers)
	}
	return sums
}
```

```
╰─ go test
PASS
ok      arrays_and_slices       0.283s
```

### リファクタリング
- sums の容量は上限があるため、もし仮に容量を超える値を入れようとするとランタイムエラーになる
  - 特に厳密に容量を決めたいわけではないので柔軟に格納してほしいので append する

```golang
package main

func SumAll(numbersToSum ...[]int) []int {
	var sums []int

	for _, numbers := range numbersToSum {
		sums = append(sums, Sum(numbers))
	}
	return sums
}
```

```
╰─ go test
PASS
ok      arrays_and_slices       0.150s
```

## SumAllTails の実装
- スライスの最初のものを除くアイテムを加算する
  - ex. `[]int{1,2,3}` `[]int{4,5,6}` と SumAllTails しようとすると下記になる
　　- `[]int{2,3}` `[]int{5,6}` = `[]int(5,11)`
  - そんなテストコードを書く

```golang
func TestSumAllTails(t *testing.T) {
	got := SumAllTails([]int{1, 2, 3}, []int{4, 5, 6})
	want := []int{5, 6}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
```

- テストが実行できるコードはこれ
  - これだと、最初を除くという要件が満たせない

```golang
func SumAllTails(numbersToSum ...[]int) []int {
	var sums []int

	for _, numbers := range numbersToSum {
		sums = append(sums, Sum(numbers))
	}
	return sums
}
```

```
╰─ go test
--- FAIL: TestSumAllTails (0.00s)
    sum_test.go:45: got [6 15] want [5 6]
FAIL
exit status 1
FAIL    arrays_and_slices       0.653s
```

### 最低限テストがパスするコード

```golang
func SumAllTails(numbersToSum ...[]int) []int {
	var sums []int

	for _, numbers := range numbersToSum {
		numbers = numbers[1:]
		sums = append(sums, Sum(numbers))
	}
	return sums
}
```

```
╰─ go test
PASS
ok      arrays_and_slices       0.140s
```

### リファクタリング
- 無し

## 検証
下記を確認したいのでテストコードを書く

- 空のスライスを関数に渡すとどうなるか
- 空のスライスの「末尾」とは何ですか？
- Goに myEmptySlice[1:] からすべての要素をキャプチャするように指示するとどうなるか?

```golang
func TestSumAllTailsPoc(t *testing.T) {
	t.Run("make the sums of some slices", func(t *testing.T) {
		got := SumAllTails([]int{1, 2}, []int{0, 9})
		want := []int{2, 9}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("safely sum empty slices", func(t *testing.T) {
		got := SumAllTailsPoc([]int{}, []int{3, 4, 5})
		want := []int{0, 9}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
```

```golang
func SumAllTailsPoc(numbersToSum ...[]int) []int {
	var sums []int

	for _, numbers := range numbersToSum {
		if len(numbers) == 0 {
			sums = append(sums, 0)
		} else {
			tail := numbers[1:]
			sums = append(sums, Sum(tail))
		}
	}
	return sums
}
```

- 明確に要件は書いてなかったが、おそらくアクセスできない要素があった場合どのような挙動になるのかを把握する目的と推測
  - テストを pass させるため、スライスの容量を確認し、0 だったら、0 を合計にいれるようにした

### リファクタリング
- テストコードの中で reflect の部分 (アサーション=どんな型が入ってるのかのチェック) がほとんど同じなので、下記のように修正できる
  - Helper関数は単にどのファイルでどの行がエラー発生したのかがわかりやすくなるためのもの
    - ref. [Helper関数](https://fresopiya.com/2022/07/29/go-helper/)

```golang
func TestSumAllTailsPoc(t *testing.T) {

	checksum := func(t *testing.T, got, want []int) {
		t.Helper()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	}
	t.Run("make the sums of some slices", func(t *testing.T) {
		got := SumAllTails([]int{1, 2}, []int{0, 9})
		want := []int{2, 9}
		checksum(t, got, want)
	})
	t.Run("safely sum empty slices", func(t *testing.T) {
		got := SumAllTailsPoc([]int{}, []int{3, 4, 5})
		want := []int{0, 9}
		checksum(t, got, want)
	})
}
```

- このリファクタリングによって、checksum(t, got ,"12345") のように `[]int` じゃないテストケースが追加されることを防いでくれる

```golang
func TestSumAllTailsPoc(t *testing.T) {

	checksum := func(t *testing.T, got, want []int) {
		t.Helper()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	}
	t.Run("make the sums of some slices", func(t *testing.T) {
		got := SumAllTails([]int{1, 2}, []int{0, 9})
		want := []int{2, 9}
		checksum(t, got, want)
	})
	t.Run("safely sum empty slices", func(t *testing.T) {
		got := SumAllTailsPoc([]int{}, []int{3, 4, 5})
		want := []int{0, 9}
		checksum(t, got, want)
	})
	t.Run("fail to sum empty slices", func(t *testing.T) {
		got := SumAllTailsPoc([]int{}, []int{})
		want := "12345"
		checksum(t, got, want)
	})
}
```

```
╰─ go test
# arrays_and_slices [arrays_and_slices.test]
./sum_test.go:70:20: cannot use want (variable of type string) as []int value in argument to checksum
FAIL    arrays_and_slices [build failed]
```