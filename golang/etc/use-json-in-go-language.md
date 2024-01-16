## Summary
Golang で json をゆるふわに使うためには interface を渡す

## struct として使う場合

こんな json がある場合
```
❯ cat sample.json
{
    "testKey1": "testValue1",
    "testKey2": "testValue2",
    "testKey3": "testValue3",
    "testKey4": "testValue4"
}
```

- 構造体を定義する
- ファイルやWebAPI等からJSON文字列を取得する
- json.UnmarshalでJSON文字列をデコードし、結果を構造体変数に放り込む
- 構造体変数のデータを処理をする

```
❯ cat before.go
package main

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "os"
)

type test struct {
        Test1 string `json:"testKey1"`
        Test2 string `json:"testKey2"`
        Test3 string `json:"testKey3"`
        Test4 string `json:"testKey4"`
}

func main() {
        jsonString, err := ioutil.ReadFile("./sample.json")
        if err != nil {
                os.Exit(1)
        }
        c := new(test)

        err = json.Unmarshal(jsonString, c)
        if err != nil {
                os.Exit(1)
        }
        fmt.Println(c)
}
```

となってめんどい

## interface を渡す場合

interface をポインタで渡すことで構造体を定義する必要がなくなる

```
❯ cat after.go
package main

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "os"
)

func main() {
        jsonString, err := ioutil.ReadFile("./sample.json")
        if err != nil {
                os.Exit(1)
        }
        var c interface{}

        err = json.Unmarshal(jsonString, &c)
        if err != nil {
                os.Exit(1)
        }
        fmt.Printf("%#v",c)
        fmt.Println(c.(map[string]interface{})["testKey1"])
}
```

取り出す場合は、type assertion して取り出せる
確認は `%#v` で確認できる

