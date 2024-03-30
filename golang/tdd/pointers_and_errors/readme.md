## Summary
社内 Go もくもく会
テスト駆動開発でGO言語を学びましょうの[ポインタとエラー編](https://andmorefine.gitbook.io/learn-go-with-tests/go-fundamentals/pointers-and-errors)

## step1

Bitcoinを預金するWallet構造体を作成

- test code

```go
package main

import "testing"

func TestWallet(t *testing.T) {

	wallet := Wallet{}

	wallet.Deposit(10)

	got := wallet.Balance()
	want := 10

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
```

- code

```go
package main

type Wallet struct {
	balance int
}

func (w *Wallet) Deposit(amount int) {
	w.balance += amount
}

func (w *Wallet) Balance() int {
	return w.balance
}
```

### 調査

最初、`func (w Wallet) Deposit(amount int) {}` としたが、下記のようにエラーになってた

```
ineffective assignment to field Wallet.balance (SA4005)
```

go では、関数やメソッドの引数はデフォルトで値渡し（コピー）つまり、その値のコピーが作成され、関数やメソッド内でそのコピーが変更されても、元の値には影響を与えない。

しかし、ポインタを使用すると、これが参照渡しつまり、そのポインタが指す元の値（ここではWalletのインスタンス）が直接変更される。

Deposit method のレシーバーを `*Wallet` にすると、Deposit method は Wallet の値を変更できるようになる。

値を表示してる Balance と違い Deposit は構造体で定義された先を変更する必要がある。

一方 Balance は値を表示されてるだけなので値渡しでも動作上は問題ない。しかし、一貫性を保つため、参照渡しとすることが多い。

### リファクタリング

wallet に対して、直接 int を定義するのではなく、Bitcoin を定義してどの値なのか説明させたい

- test code

```go
package main

import "testing"

func TestWallet(t *testing.T) {

	wallet := Wallet{}

	wallet.Deposit(Bitcoin(10))

	got := wallet.Balance()

	want := Bitcoin(10)

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
```

- code

```go
package main

import "fmt"

type Bitcoin int
type Wallet struct {
	balance Bitcoin
}

type Stringer interface {
	Stringer() string
}

func (w *Wallet) Deposit(amount Bitcoin) {
	w.balance += amount
}

func (w Wallet) Balance() Bitcoin {
	return w.balance
}

func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}
```

### なぜ、String() を追加したら、print で string を指定することで String() が実行されるのか

Go では、特定の型に対して `String()` メソッドを定義すると、その型の値を文字列としてフォーマットする際にそのメソッドが自動的に呼び出される。

Go の fmt パッケージが `Stringer` インターフェースを提供してて、このインターフェースを満たす任意の型 `（String()メソッドを持つ型）` は、fmtパッケージの関数（例えばPrintfやErrorfなど）によって自動的に文字列に変換される。

`Bitcoin` 型の値を `fmt.Errorf` や `fmt.Printf` などの関数で文字列としてフォーマットすると、`String()` メソッドが自動的に呼び出され、その結果が使用される。

これにより、Bitcoin型の値は常に"<値> BTC"の形式で表示される。

したがって、`t.Errorf("got %s want %s", got, want)` というコードでは、`got` と `want` が `Bitcoin` 型である場合、それぞれの `String()` メソッドが呼び出されてエラーメッセージが生成される。

## step2

Bitcoin を消費される `withDraw` 機能を追加する

- test code

```go
package main

import "testing"

func TestWallet(t *testing.T) {

	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}

		wallet.Deposit(Bitcoin(10))

		got := wallet.Balance()

		want := Bitcoin(10)

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("Withdraw", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(20)}

		wallet.WithDraw(Bitcoin(10))

		got := wallet.Balance()

		want := Bitcoin(10)

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
}
```

- code

```go
package main

import "fmt"

type Bitcoin int
type Wallet struct {
	balance Bitcoin
}

type Stringer interface {
	Stringer() string
}

func (w *Wallet) Deposit(amount Bitcoin) {
	w.balance += amount
}

func (w *Wallet) WithDraw(amount Bitcoin) {
	w.balance -= amount
}

func (w *Wallet) Balance() Bitcoin {
	return w.balance
}

func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}
```

### リファクタリング

重複があるので修正

- test code

```go
package main

import "testing"

func TestWallet(t *testing.T) {

	assertBalance := func(t *testing.T, wallet Wallet, want Bitcoin) {
		t.Helper()
		got := wallet.Balance()

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	}

	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("Withdraw", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(20)}
		wallet.WithDraw(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})
}
```

- code

```go
package main

import "fmt"

type Bitcoin int
type Wallet struct {
	balance Bitcoin
}

type Stringer interface {
	Stringer() string
}

func (w *Wallet) Deposit(amount Bitcoin) {
	w.balance += amount
}

func (w *Wallet) WithDraw(amount Bitcoin) {
	w.balance -= amount
}

func (w *Wallet) Balance() Bitcoin {
	return w.balance
}

func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}
```

## step3

残高以上に `WithDraw` された場合エラーを返すようにしたい

- test code

```go
package main

import "testing"

func TestWallet(t *testing.T) {

	assertBalance := func(t *testing.T, wallet Wallet, want Bitcoin) {
		t.Helper()
		got := wallet.Balance()

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	}

	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("Withdraw", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(20)}
		wallet.WithDraw(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("Withdraw insufficient funds", func(t *testing.T) {
		startingBalance := Bitcoin(20)
		wallet := Wallet{startingBalance}
		err := wallet.WithDraw(Bitcoin(100))

		assertBalance(t, wallet, startingBalance)

		if err == nil {
			t.Error("wanted an error but didn't get one")
		}
	})
}
```

- code

```go
package main

import (
	"errors"
	"fmt"
)

type Bitcoin int
type Wallet struct {
	balance Bitcoin
}

type Stringer interface {
	Stringer() string
}

func (w *Wallet) Deposit(amount Bitcoin) {
	w.balance += amount
}

func (w *Wallet) WithDraw(amount Bitcoin) error {

	if amount > w.balance {
		return errors.New("oh no")
	}

	w.balance -= amount
	return nil
}

func (w *Wallet) Balance() Bitcoin {
	return w.balance
}

func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}
```

### リファクタリング

テストを読みやすくするために、エラーチェック用のクイックテストヘルパーを作成する

- test code

```go
package main

import "testing"

func TestWallet(t *testing.T) {

	assertBalance := func(t *testing.T, wallet Wallet, want Bitcoin) {
		t.Helper()
		got := wallet.Balance()

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	}
	assertError := func(t *testing.T, err error) {
		t.Helper()
		if err == nil {
			t.Error("wanted an error but didn't get one")
		}
	}

	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("Withdraw", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(20)}
		wallet.WithDraw(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("Withdraw insufficient funds", func(t *testing.T) {
		startingBalance := Bitcoin(20)
		wallet := Wallet{startingBalance}
		err := wallet.WithDraw(Bitcoin(100))

		assertBalance(t, wallet, startingBalance)
		assertError(t, err)

	})
}
```

- code

```go
package main

import (
	"errors"
	"fmt"
)

type Bitcoin int
type Wallet struct {
	balance Bitcoin
}

type Stringer interface {
	Stringer() string
}

func (w *Wallet) Deposit(amount Bitcoin) {
	w.balance += amount
}

func (w *Wallet) WithDraw(amount Bitcoin) error {

	if amount > w.balance {
		return errors.New("oh no")
	}

	w.balance -= amount
	return nil
}

func (w *Wallet) Balance() Bitcoin {
	return w.balance
}

func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}
```

## step4

比較する stringのヘルパーを更新する

- test code

```go
package main

import "testing"

func TestWallet(t *testing.T) {

	assertBalance := func(t *testing.T, wallet Wallet, want Bitcoin) {
		t.Helper()
		got := wallet.Balance()

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	}
	assertError := func(t *testing.T, got error, want string) {
		t.Helper()
		if got == nil {
			t.Fatal("didn't get an error but wanted one")
		}

		if got.Error() != want {
			t.Errorf("got %q, want %q", got, want)
		}
	}
	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("Withdraw", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(20)}
		wallet.WithDraw(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("Withdraw insufficient funds", func(t *testing.T) {
		startingBalance := Bitcoin(20)
		wallet := Wallet{startingBalance}
		err := wallet.WithDraw(Bitcoin(100))

		assertBalance(t, wallet, startingBalance)
		assertError(t, err, "cannot withdraw, insufficient funds")

	})
}
```

- code

```go
package main

import (
	"errors"
	"fmt"
)

type Bitcoin int
type Wallet struct {
	balance Bitcoin
}

type Stringer interface {
	Stringer() string
}

func (w *Wallet) Deposit(amount Bitcoin) {
	w.balance += amount
}

func (w *Wallet) WithDraw(amount Bitcoin) error {

	if amount > w.balance {
		return errors.New("cannot withdraw, insufficient funds")
	}

	w.balance -= amount
	return nil
}

func (w *Wallet) Balance() Bitcoin {
	return w.balance
}

func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}
```

### リファクタリング

テストコードと Withdrawコードの両方でエラーメッセージが重複してるので解消させる

- test code

```go
package main

import "testing"

func TestWallet(t *testing.T) {

	assertBalance := func(t *testing.T, wallet Wallet, want Bitcoin) {
		t.Helper()
		got := wallet.Balance()

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	}
	assertError := func(t *testing.T, got error, want string) {
		t.Helper()
		if got == nil {
			t.Fatal("didn't get an error but wanted one")
		}

		if got.Error() != want {
			t.Errorf("got %q, want %q", got, want)
		}
	}
	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("Withdraw", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(20)}
		wallet.WithDraw(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("Withdraw insufficient funds", func(t *testing.T) {
		startingBalance := Bitcoin(20)
		wallet := Wallet{startingBalance}
		err := wallet.WithDraw(Bitcoin(100))

		assertBalance(t, wallet, startingBalance)
		assertError(t, err, ErrInsufficientFunds.Error())

	})
}
```

- code

```go
package main

import (
	"errors"
	"fmt"
)

type Bitcoin int
type Wallet struct {
	balance Bitcoin
}

type Stringer interface {
	Stringer() string
}

var ErrInsufficientFunds = errors.New("cannot withdraw, insufficient funds")

func (w *Wallet) Deposit(amount Bitcoin) {
	w.balance += amount
}

func (w *Wallet) WithDraw(amount Bitcoin) error {

	if amount > w.balance {
		return ErrInsufficientFunds
	}

	w.balance -= amount
	return nil
}

func (w *Wallet) Balance() Bitcoin {
	return w.balance
}

func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}
```

## extend errcheck

下記エラーチェックが引っかかるので修正する

```
╭─    ~/.github.com/til/golang/tdd/pointers_and_errors  on   master !3 ?1                                                     ✔  at 02:28:33 
╰─ errcheck .
wallet_test.go:33:18:   wallet.WithDraw(Bitcoin(10))
```

- test code

```go
package main

import "testing"

func assertBalance(t *testing.T, wallet Wallet, want Bitcoin) {
	t.Helper()
	got := wallet.Balance()

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func assertError(t *testing.T, got error, want string) {
	t.Helper()
	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got.Error() != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertNoError(t *testing.T, got error) {
	t.Helper()
	if got != nil {
		t.Fatal("got an error but didn't want one")
	}
}

func TestWallet(t *testing.T) {
	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("Withdraw", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(20)}
		err := wallet.WithDraw(Bitcoin(10))

		assertBalance(t, wallet, Bitcoin(10))
		assertNoError(t, err)
	})

	t.Run("Withdraw insufficient funds", func(t *testing.T) {
		startingBalance := Bitcoin(20)
		wallet := Wallet{startingBalance}
		err := wallet.WithDraw(Bitcoin(100))

		assertBalance(t, wallet, startingBalance)
		assertError(t, err, ErrInsufficientFunds.Error())
	})
}
```

- code

```go
package main

import (
	"errors"
	"fmt"
)

type Bitcoin int
type Wallet struct {
	balance Bitcoin
}

type Stringer interface {
	Stringer() string
}

var ErrInsufficientFunds = errors.New("cannot withdraw, insufficient funds")

func (w *Wallet) Deposit(amount Bitcoin) {
	w.balance += amount
}

func (w *Wallet) WithDraw(amount Bitcoin) error {

	if amount > w.balance {
		return ErrInsufficientFunds
	}

	w.balance -= amount
	return nil
}

func (w *Wallet) Balance() Bitcoin {
	return w.balance
}

func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}
```

## 全体通しての補足
- ラムダ関数（無名関数）のメリットとデメリット
  - メリット:
    - テスト関数内でのみ使用されるヘルパー関数を定義する際に便利。これにより、その関数がテスト関数の外部で誤って使用されることを防ぐことができる。
    - ラムダ関数はその定義が行われたスコープ内の変数にアクセスできるため、テスト関数内の特定の変数を再利用することが容易になる。
  - デメリット:
    - ラムダ関数はその定義が行われたテスト関数内でのみ使用できる。したがって、同じロジックを複数のテスト関数で再利用する必要がある場合、それぞれのテスト関数でラムダ関数を定義する必要がある。
- 独立したfuncのメリットとデメリット
  - メリット:
    - 独立したfuncはパッケージ全体で再利用可能。同じロジックを複数のテスト関数で再利用する必要がある場合、それを一度だけ定義して複数のテスト関数から呼び出すことができる。
    - 独立したfuncはテストレポートで独自の名前を持つため、テストが失敗した場合のデバッグが容易になる。
  - デメリット:
    - 独立したfuncはその定義が行われたスコープ外の変数に直接アクセスできないので、必要なすべての値を引数として渡す必要がある。
    - テスト関数内でのみ使用されるヘルパー関数を定義する場合、その関数がテスト関数の外部で誤って使用される可能性がある。