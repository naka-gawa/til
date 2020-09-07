# nodeName
- Podを作成するときにnodeNameをstaticに指定することができる
- taintとか付いていなかったらこれを確認すると良い

```
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - image: nginx
    name: nginx
  nodeName: master
```
