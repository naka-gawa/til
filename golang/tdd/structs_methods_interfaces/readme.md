## Summary
社内 Go もくもく会
テスト駆動開発でGO言語を学びましょうの[構造体、メソッド、インターフェース編](https://andmorefine.gitbook.io/learn-go-with-tests/go-fundamentals/structs-methods-and-interfaces)

## 長方形の高さと幅から周囲を算出したい
### テストが実行できるコード

- test code
```golang
package main

import "testing"

func TestPerimeter(t *testing.T) {
	got := Perimeter(10.0, 10.0)
	want := 40.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}
```

```golang
package main

func Perimeter(w float64, h float64) float64 {
	return 0
}
```

### テストを成功させるには

```golang
package main

func Perimeter(w float64, h float64) float64 {
	return 2 * (w + h)
}
```

```
╰─ go test
PASS
ok      structs_methods_interfaces      0.469s
```

## 長方形の面積を返す `Area(w, h float64)` を実装する
- 愚直に実装するならこんなかんじ

```golang
func TestArea(t *testing.T) {
	got := Area(12.0, 6.0)
	want := 72.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}
```

```golang
func Area(w float64, h float64) float64 {
	return w * h
}
```

### リファクタリング
- 高さと幅しか渡されてないので、三角形という可能性もあるが、三角形の場合は算出ロジックが違うので間違った値を返してしまう
- `struct` を使って構造体という名前をつけられるデータを保存するフィールド作成する

```golang
package main

import "testing"

type Rectangle struct {
	Width  float64
	Height float64
}

func TestPerimeter(t *testing.T) {
	rectangle := Rectangle{10.0, 10.0}
	got := Perimeter(rectangle)
	want := 40.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}

func TestArea(t *testing.T) {
	rectangle := Rectangle{12.0, 6.0}
	got := Area(rectangle)
	want := 72.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}
```

```golang
package main

func Perimeter(r Rectangle) float64 {
	return 2 * (r.Width + r.Height)
}

func Area(r Rectangle) float64 {
	return r.Width * r.Height
}
```

## 円の面積を算出する関数を実装する
- 同じように円の構造体を作る
  - ただし、メソッドを利用する

```golang
package main

import "testing"

func TestPerimeter(t *testing.T) {
	rectangle := Rectangle{10.0, 10.0}
	got := Perimeter(rectangle)
	want := 40.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}

func TestArea(t *testing.T) {

	t.Run("rectangles", func(t *testing.T) {
		rectangle := Rectangle{12, 6}
		got := rectangle.Area()
		want := 72.0

		if got != want {
			t.Errorf("got %.2f want %.2f", got, want)
		}
	})

	t.Run("circles", func(t *testing.T) {
		circle := Circle{10}
		got := circle.Area()
		want := 314.1592653589793

		if got != want {
			t.Errorf("got %.2f want %.2f", got, want)
		}
	})
}
```

```golang
package main

import "math"

func Perimeter(r Rectangle) float64 {
	return 2 * (r.Width + r.Height)
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

type Circle struct {
	Radius float64
}

func (r Circle) Area() float64 {
	return math.Pi * r.Radius * r.Radius
}
```

### メソッドってなによ？
- メソッド = **レシーバーをもつ関数**
  - `t.Errorf` を呼び出してきたが、これは `t(testing.T)` のインスタンスでメソッド `Errorf` を呼び出している
- メソッド宣言は、識別子（メソッド名）をメソットに紐づけ、メソッドをレシーバーの基本タイプに関連付ける
  - よくわからないのでコードを見ながら飲み込む
    - 引用元: https://graff-it-i.com/2022/05/22/go-method-next/

- メソッドを利用しないケース
  - UserShow 関数は構造体を受け取って、格納されてる名前と年齢を表示させてる

```golang
type User struct {
	Name string
	Age  int
}

// 引数はUserとなる
func UserShow(user User) {
	fmt.Println(user.Name, user.Age)
}

func main() {
	tom := User{
		Name: "Tom",
		Age:  20,
	}
  // 初期化した構造体を引数とする
	UserShow(tom)
  // Tom 20
}
```
- メソッドを利用するケース
  - 構造体に紐づく関数を実行している

```golang
type User struct {
	Name string
	Age  int
}

// メソッドを定義
// 構造体Userを指定している
func (u User) Show() {
　// 構造体Userの値を扱うことができる
	fmt.Println(u.Name, u.Age)
}

func main() {
	tom := User{
		Name: "Tom",
		Age:  20,
	}
  // Showメソッドを呼び出す
	tom.Show()
}
```

- メソッドを扱うメリットとしては、構造体と関数を分離することができる

```golang
type User struct {
	Name string
	Age  int
}

// メソッドを使わないと無限に func の名前が長くなり何をやっているのかわからなくなりそう
func UserShow(user User) {
    ~snip~
}

func UserDelete(user User) {
    ~snip~
}

func UserCreate(user User) {
    ~snip~
}

// メソッドを使うと対象に対する命令という関係性にスッキリする
func (u User) Show() {
    ~snip~
}

func (u User) Create() {
    ~snip~
}

func (u User) Delete() {
    ~snip~
}
```

### リファクタリング
- やりたいことは図形を識別できる構造体を取得して、`Area()` メソッドを呼び出して結果を確認することなので、インターフェイスを使って意図しない図形を渡された場合は fail するようにしたい
  - インターフェイスを使用して、必要なもののみを宣言する
  - 渡したタイプがインターフェイスと一致する場合コンパイルされる
    - `circle` や `rectangle` は `Area()` メソッドがあるのでインターフェイスは要件を満たす
    - `triangle` は `Area()` メソッドを持たないので要件を満たさずコンパイル失敗する
      - そういう意味でもメソッドは有用かもしれない

```golang
package main

import "testing"

type Shape interface {
	Area() float64
}

func TestPerimeter(t *testing.T) {
	rectangle := Rectangle{10.0, 10.0}
	got := Perimeter(rectangle)
	want := 40.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}

func TestArea(t *testing.T) {

	checkArea := func(t *testing.T, shape Shape, want float64) {
		t.Helper()
		got := shape.Area()
		if got != want {
			t.Errorf("%#v got %.2f want %.2f", shape, got, want)
		}
	}
	t.Run("rectangles", func(t *testing.T) {
		rectangle := Rectangle{12, 6}
		checkArea(t, rectangle, 72.0)
	})

	t.Run("circles", func(t *testing.T) {
		circle := Circle{10}
		checkArea(t, circle, 314.1592653589793)
	})
}
```

```golang
package main

import "math"

func Perimeter(r Rectangle) float64 {
	return 2 * (r.Width + r.Height)
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

type Circle struct {
	Radius float64
}

func (r Circle) Area() float64 {
	return math.Pi * r.Radius * r.Radius
}
```

```
╰─ go test -v

=== RUN   TestPerimeter
--- PASS: TestPerimeter (0.00s)
=== RUN   TestArea
=== RUN   TestArea/rectangles
=== RUN   TestArea/circles
--- PASS: TestArea (0.00s)
    --- PASS: TestArea/rectangles (0.00s)
    --- PASS: TestArea/circles (0.00s)
PASS
ok      structs_methods_interfaces      0.155s
```

### テーブル駆動テストによるリファクタリング
- 同じような方法でテストできるテストケースのリストを作成するのに役立つ

```golang
package main

import "testing"

type Shape interface {
	Area() float64
}

func TestPerimeter(t *testing.T) {
	rectangle := Rectangle{10.0, 10.0}
	got := Perimeter(rectangle)
	want := 40.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}

func TestArea(t *testing.T) {

	areaTests := []struct {
		shape Shape
		want  float64
	}{
		{Rectangle{12, 6}, 72.0},
		{Circle{10}, 314.1592653589793},
	}

	for _, tt := range areaTests {
		got := tt.shape.Area()
		if got != tt.want {
			if got != tt.want {
				t.Errorf("got %.2f want %.2f", got, tt.want)
			}
		}
	}
}
```

```
╰─ go test -v
=== RUN   TestPerimeter
--- PASS: TestPerimeter (0.00s)
=== RUN   TestArea
--- PASS: TestArea (0.00s)
PASS
ok      structs_methods_interfaces      0.656s
```

## 三角形もサポートする

```golang
package main

import "testing"

type Shape interface {
	Area() float64
}

func TestPerimeter(t *testing.T) {
	rectangle := Rectangle{10.0, 10.0}
	got := Perimeter(rectangle)
	want := 40.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}

func TestArea(t *testing.T) {

	areaTests := []struct {
		shape Shape
		want  float64
	}{
		{Rectangle{12, 6}, 72.0},
		{Circle{10}, 314.1592653589793},
		{Triangle{12, 6}, 36.0},
	}

	for _, tt := range areaTests {
		got := tt.shape.Area()
		if got != tt.want {
			if got != tt.want {
				t.Errorf("got %.2f want %.2f", got, tt.want)
			}
		}
	}
}
```

```golang
package main

import "math"

func Perimeter(r Rectangle) float64 {
	return 2 * (r.Width + r.Height)
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

type Circle struct {
	Radius float64
}

func (r Circle) Area() float64 {
	return math.Pi * r.Radius * r.Radius
}

type Triangle struct {
	Width  float64
	Height float64
}

func (t Triangle) Area() float64 {
	return 0.5 * t.Width * t.Height
}
```

```
package main

import "math"

func Perimeter(r Rectangle) float64 {
	return 2 * (r.Width + r.Height)
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

type Circle struct {
	Radius float64
}

func (r Circle) Area() float64 {
	return math.Pi * r.Radius * r.Radius
}

type Triangle struct {
	Width  float64
	Height float64
}

func (t Triangle) Area() float64 {
	return 0.5 * t.Width * t.Height
}
```

### リファクタリング
- テストケースがそれぞれ何を意味しているのか不明なので、明示したい
  - `t.Run` を wrap することでどこでコケたのかがわかるようになる

```golang
func TestArea(t *testing.T) {

	areaTests := []struct {
		name  string
		shape Shape
		want  float64
	}{
		{name: "rectangle", shape: Rectangle{Width: 12, Height: 6}, want: 72.0},
		{name: "circle", shape: Circle{Radius: 10}, want: 314.1592653589793},
		{name: "triangle", shape: Triangle{Width: 12, Height: 6}, want: 36.0},
	}

	for _, tt := range areaTests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.shape.Area()
			if got != tt.want {
				t.Errorf("got %.2f want %.2f", got, tt.want)
			}
		})
	}
}
```

```
OK pattern
╰─ go test -v
=== RUN   TestPerimeter
--- PASS: TestPerimeter (0.00s)
=== RUN   TestArea
=== RUN   TestArea/rectangle
=== RUN   TestArea/circle
=== RUN   TestArea/triangle
--- PASS: TestArea (0.00s)
    --- PASS: TestArea/rectangle (0.00s)
    --- PASS: TestArea/circle (0.00s)
    --- PASS: TestArea/triangle (0.00s) w
PASS
ok      structs_methods_interfaces      0.481s

NG pattern
╰─ go test -v
=== RUN   TestPerimeter
--- PASS: TestPerimeter (0.00s)
=== RUN   TestArea
=== RUN   TestArea/rectangle
    perimeter_test.go:35: got 72.00 want 70.00
=== RUN   TestArea/circle
=== RUN   TestArea/triangle
--- FAIL: TestArea (0.00s)
    --- FAIL: TestArea/rectangle (0.00s)
    --- PASS: TestArea/circle (0.00s)
    --- PASS: TestArea/triangle (0.00s)
FAIL
exit status 1
FAIL    structs_methods_interfaces      0.482s
```