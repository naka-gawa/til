# kubectl を使った label の付け方
- yaml ならこんな感じでラベルをつけられる
```
apiVersion: apps/v1
kind: Pod
metadata
  name: test
  labels:
    app: testlabel
spec:
  containers:
  - image: nginx
    name: nginx
```

- でも kubectl を使ったほうが楽ちん
- リソース作成時に一緒にラベリングするとき
```
kubectl run nginx --image=nginx --restart=Never --labels='app=testlabel'
```

- 既存リソースにラベリングするとき
```
k label node node01 color=blue
node/node01 labeled
```

- 確認
```
k get node --show-labels
NAME      STATUS    ROLES     AGE       VERSION   LABELS
master    Ready     master    30m       v1.11.3   beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,kubernetes.io/hostname=master,node-role.kubernetes.io/master=
node01    Ready     <none>    30m       v1.11.3   beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,color=blue,kubernetes.io/hostname=node01
```
