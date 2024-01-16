## Summary
社内 Go もくもく会
テスト駆動開発でGO言語を学びましょうの[整数編](https://andmorefine.gitbook.io/learn-go-with-tests/go-fundamentals/integers)
整数のテストを行う

### TDDプロセスについて
まずテストコードを書き、失敗するテストを用意する。これにより、なぜ失敗するのか？が明確になり、理解が進む。
テストコードを成功させるために、必要な最小限のコードを用意する。そこからリファクタリングを重ねることで
機能の安全性を担保させながらリファクタリングを進めることができる。

#### Step.1 テストコード

```golang
package integers

import (
    "testing"
)

func TestAdder(t *testing.T) {
    sum := Add(2, 2)
    expected := 4

    if sum != expected {
        t.Errorf("expected '%d' but got '%d'", expected, sum)
    }
}
```

```output
╰─ go test
# github.com/naka-gawa/tiltdd/tdd [github.com/naka-gawa/tiltdd/tdd.test]
./adder_test.go:6:12: undefined: Add
FAIL    github.com/naka-gawa/tiltdd/tdd [build failed]

```

`Add` が定義されてない

#### Step.2 最小コードを用意
テストが実行できる最小限のコードを用意する

```golang
package integers

// Add takes two integers and return the sum of them.
func Add(x, y int) int {
    return 0
}
```

```output
╰─ go test
--- FAIL: TestAdder (0.00s)
    adder_test.go:10: expected '4' but got '0'
FAIL
exit status 1
FAIL    github.com/naka-gawa/tiltdd/tdd 0.454s

```

test は実行できたが、期待値と違う

#### Step.3 合格する最小コードを用意


```golang
package integers

// Add takes two integers and return the sum of them.
func Add(x, y int) int {
    return x + y
}
```

```output
╰─ go test
PASS
ok      github.com/naka-gawa/tiltdd/tdd 0.449s
```

#### Step.4 リファクタリング


