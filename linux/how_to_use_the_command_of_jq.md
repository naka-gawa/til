# jq の使い方

## Overview
- JSON から値を取り出したり、集計したり、整形して表示したりするコマンド

## 基本
- こんな JSON を例に取る
```
echo '{"item_list":[{"item_id":1,"name":"apple","price":50},{"item_id":2,"name":"orange","price":30}]}
```
- jq に読み込ませて、root から取得する
```
😀 ❯❯❯ echo '{"item_list":[{"item_id":1,"name":"apple","price":50},{"item_id":2,"name":"orange","price":30}]}' | jq .
{
  "item_list": [
    {
      "item_id": 1,
      "name": "apple",
      "price": 50
    },
    {
      "item_id": 2,
      "name": "orange",
      "price": 30
    }
  ]
}
```
- item_listだけを取得する
  - 配列の中身を取りたいは `[]` を追加する
  - 配列なので添字でアクセス可能
```
😱 ❯❯❯ echo '{"item_list":[{"item_id":1,"name":"apple","price":50},{"item_id":2,"name":"orange","price":30}]}' | jq -r '.item_list[]'
{
  "item_id": 1,
  "name": "apple",
  "price": 50
}
{
  "item_id": 2,
  "name": "orange",
  "price": 30
}
```

- 値だけを取り出したい場合はこんな感じ
  - 結果を別コマンドに渡したい場合は ダブルクォート `"` が邪魔なので、`-r` で消すことができる
```
😀 ❯❯❯ echo '{"item_list":[{"item_id":1,"name":"apple","price":50},{"item_id":2,"name":"orange","price":30}]}' | jq '.item_list[1].name'
"orange"

😀 ❯❯❯ echo '{"item_list":[{"item_id":1,"name":"apple","price":50},{"item_id":2,"name":"orange","price":30}]}' | jq -r '.item_list[1].name'
orange
```

# ちょっと複雑な使い方
- 正規表現を使った検索
  - `*ple` にマッチした文字列を抜き出したい
  - `select(検索対象 | 検索文字列)` みたいに使う
```
😱 ❯❯❯ echo '{"item_list":[{"item_id":1,"name":"apple","price":50},{"item_id":2,"name":"orange","price":30}]}' | jq '.item_list[] | select(.name | test("ple"))'
{
  "item_id": 1,
  "name": "apple",
  "price": 50
}

😀 ❯❯❯ echo '{"item_list":[{"item_id":1,"name":"apple","price":50},{"item_id":2,"name":"orange","price":30}]}' | jq '.item_list[].name | select(. | test("ple"))'
"apple"
```

- 整形
  - csv にしたい場合は pipe でつないで `,`区切りで配列にすればよい
```
😀 ❯❯❯ echo '{"item_list":[{"item_id":1,"name":"apple","price":50},{"item_id":2,"name":"orange","price":30}]}' | jq -r '.item_list[]'
{
  "item_id": 1,
  "name": "apple",
  "price": 50
}
{
  "item_id": 2,
  "name": "orange",
  "price": 30
}

😀 ❯❯❯ echo '{"item_list":[{"item_id":1,"name":"apple","price":50},{"item_id":2,"name":"orange","price":30}]}' | jq -r '.item_list[] | [.name,.price]'
[
  "apple",
  50
]
[
  "orange",
  30
]

😀 ❯❯❯ echo '{"item_list":[{"item_id":1,"name":"apple","price":50},{"item_id":2,"name":"orange","price":30}]}' | jq -cr '.item_list[] | [.name,.price]'
["apple",50]
["orange",30]
```
