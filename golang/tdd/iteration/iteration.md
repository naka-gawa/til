# Summary
社内 Go もくもく会
テスト駆動開発でGO言語を学びましょうの[反復、繰り返し編](https://andmorefine.gitbook.io/learn-go-with-tests/go-fundamentals/iteration)

## first code
- `go mod init iteration` と `go mod tidy` を忘れない

test code
```golang
package iteration

import "testing"

func TestRepeat(t *testing.T) {
    repeated := Repeat("a")
    expected := "aaaaa"

    if repeated != expected {
        t.Errorf("expected %q but got %q", expected, repeated)
    }
}
```

- test を実行するための最低限のコード

```
package iteration

func Repeat(charactor string) string {
    return ""
}
```

- test 結果
  - 期待値と違うと言われてるので修正

```
╰─ go test
--- FAIL: TestRepeat (0.00s)
    repeat_test.go:10: expected "aaaaa" but got ""
FAIL
exit status 1
FAIL    iteration       0.555s

```

## second commit
- 修正後
```
package iteration

func Repeat(charactor string) string {
    var repeated string
    for i := 0; i < 5; i++ {
        repeated = repeated + charactor
    }
    return repeated
}
```

- test 結果
  - test が通る
```
╰─ go test
PASS
ok      iteration       0.484s

```

## refactoring
- `+=` 代入演算子を導入する

```golang
package iteration

const repeatCount = 5

func Repeat(charactor string) string {
    var repeated string
    for i := 0; i < repeatCount; i++ {
        repeated += charactor
    }
    return repeated
}
```

- test 結果

```
╰─ go test
PASS
ok      iteration       0.496s

```

## exercise

- テストを変更して、発信者の文字が繰り返される回数を指定子、コードを修正できるようにする
  - `repeat_test.go` を修正

```golang
package iteration

import "testing"

func TestRepeat(t *testing.T) {
    repeated := Repeat("a", 6)
    expected := "aaaaaa"

    if repeated != expected {
        t.Errorf("expected %q but got %q", expected, repeated)
    }
}

func BenchmarkRepeat(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Repeat("a", 6)
    }
}
```



