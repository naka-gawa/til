## Summary
社内 Go もくもく会
テスト駆動開発でGO言語を学びましょうの[maps編](https://andmorefine.gitbook.io/learn-go-with-tests/go-fundamentals/maps)

## step 1

key でアイテムを保存し、すばやく検索する方法を確認する

マップを使用すると、辞書と同じようにアイテムを保存できる

keyは単語、valueは定義、そして、独自の辞書を構築するよりも、マップについて学ぶより良い方法は何か？

- test code

```go
package main

import (
	"testing"
)

func TestDictionary(t *testing.T) {
	dictionary := map[string]string{"test": "this is just a test"}

	got := Search(dictionary, "test")
	want := "this is just a test"

	if got != want {
		t.Errorf("got %q want %q given, %q", got, want, "test")
	}
}
```

- code

```go
package main

func Search(dictionary map[string]string, word string) string {
	return dictionary[word]
}
```

### リファクタリング

実装をより一般的なものにするために、 assertStringsヘルパーを作成することに。

- test code

```go
package main

import (
	"testing"
)

func TestDictionary(t *testing.T) {
	dictionary := Dictionary{"test": "this is just a test"}

	got := dictionary.Search("test")
	want := "this is just a test"

	assertStrings(t, got, want)
}

func assertStrings(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
```

- code

```go
package main

type Dictionary map[string]string

func (d Dictionary) Search(word string) string {
	return d[word]
}
```