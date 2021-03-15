## Summary
Golang 開発ツール

### gofmt
標準でついているフォーマッタ
ソースのインデントなどを整形してくれる

before
```
package main

import (
"fmt"
"math/rand"
)

func main() {
fmt.Println("My favorite number is", rand.Intn(10))
}
```

after
```
package main

import (
  "fmt"
  "math/rand"
)

func main() {
  fmt.Println("My favorite number is", rand.Intn(10))
}
```
