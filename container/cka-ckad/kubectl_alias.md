# よく使うalias

忘れたりするのでよく使うaliasをメモ

## さくっと設定できる系
```
😀 ❯❯❯ alias k='kubectl'

😀 ❯❯❯ alias ka='kubectl apply'

😀 ❯❯❯ alias kd='kubectl delete'

😀 ❯❯❯ alias kall='kubectl get all'
```

## 自作関数にしたほうが良い系
- kubectl get all -n [namespace] では全リソースを取得できないので、下記を使う
```
function kgetall {
  for i in $(kubectl get namespaces); do
    kubectl api-resources --verbs=list --namespaced -o name | xargs -n 1 kubectl get --show-kind --ignore-not-found -n $i
  done
}
```
