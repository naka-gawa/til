## Summary
社内 Go もくもく会
テスト駆動開発でGO言語を学びましょうの[Hello, World編](https://andmorefine.gitbook.io/learn-go-with-tests/go-fundamentals/hello-world)

### Test コードを書く際のいくつかのルール

- `xxx_test.go` という名前が必要
- テスト関数は `Test` で始まる必要がある
- テスト関数は `t *testing.T` という引数のみを受け取る
- `*testing.T` 型を使うには他のファイルと同様に import で `testing` を呼び出してあげる必要がある

### Hello World を出力する関数のテストを書く

- helloworld.go
```golang
package main

import "fmt"

func Hello() string {
    return "hello world"
}

func main() {
    fmt.Println(Hello())
}
```

- helloworld_test.go
```golang
package main

import "testing"

func TestHello(t *testing.T) {
    got := Hello()
    want := "hello world"

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

- output
```shell
╰─ go test -v
=== RUN   TestHello
--- PASS: TestHello (0.00s)
PASS
ok      github.com/naka-gawa/tiltdd/tdd/helloworld      0.330s

```

#### `%q` ってなあに？
書式指定子で使われる `%q` について調べてみた

- integer の場合、シングルクォートで囲まれた文字列リテラルにフォーマットされる

```golang
package main

import "fmt"

func main() {
    integer := 12456
    fmt.Printf("出力は %q", integer)
}
// output: 出力は 'エ'
// https://0g0.org/unicode/30A7/
```

- string or slice 場合は、ダブルクォートで囲まれた文字列にフォーマットする。

```golang
package main

import "fmt"

func main() {
    s := "oange"
    slice := []string{"o", "a", "n", "g", "e"}
    fmt.Printf("出力は %q\n", s)
    fmt.Printf("出力は %q\n", slice)
}

/**
出力は "oange"
出力は ["o" "a" "n" "g" "e"]

**/
```

##### use case

[Goの書式指定子%qを絶対に忘れない](https://developer.so-tech.co.jp/entry/2022/08/31/110108) にある通り

文字列を出力する際に、空になる可能性があるなら%q を使うと良さそう。


```golang
u, err := url.Parse(r.Image.URL)
if err != nil {
    return nil, fmt.Errorf("画像URLの形式が正しくありません: %s : %w", r.Image.URL, err)
}

/**
画像URLの形式が正しくありません:  : hogehogeerror

クォートで囲われないのでパッと見分かりづらい（というよりも後者のほうがわかりやすい）
**/

u, err := url.Parse(r.Image.URL)
if err != nil {
    return nil, fmt.Errorf("画像URLの形式が正しくありません: %q : %w", r.Image.URL, err)
}

/**
画像URLの形式が正しくありません: "" : hogehogeerror

**/
```



整数型の場合は、Goシンタックスで安全にエスケープされたシングルクォートで囲まれた文字列リテラルにフォーマットする。
文字列型もしくはスライス型の場合は、Goシンタックスで安全にエスケープされたダブルクォートで囲まれた文字列にフォーマットする。
単一の整数コードポイントまたはルーン文字列（rune 型）の場合は、無効な Unicode コードポイントは strconv.QuoteRune のように、 Unicode 置換文字 U+FFFD に変更される。



### TDDプロセスについて
まずテストコードを書き、失敗するテストを用意する。これにより、なぜ失敗するのか？が明確になり、理解が進む。
テストコードを成功させるために、必要な最小限のコードを用意する。そこからリファクタリングを重ねることで
機能の安全性を担保させながらリファクタリングを進めることができる。


