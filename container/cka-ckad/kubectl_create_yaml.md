# Certification Tips!

## Manifest yaml fileの生成
### Pod
```
kubectl run nginx --restart=Never --image=nginx --dry-run -o yaml
```

### Deployment
```
kubectl run nginx --restart=Always --image=nginx --dry-run -o yaml
```

### Replicas 4 の Deployment
```
kubectl run nginx --restart=Always --replicas=4 --images=nginx --dry-run -o yaml
```
### Service ClusterIP
```
kubectl expose nginx --port=32689 --target-port=80 --dry-run -o yaml
```

### Service NodePort
```
Service ClusterIP で作ったリソースに対して, kubectl edit で ClusterIP -> NodePortに変える
NodePort はランダムに割り当てられるが、もし NodePort を固定したい場合は nodePort: xxxx を追加する
```
 - 補足
   - kubectl get svc で表示されるものについて
   - NodePort リソースを作成すると下記のように表記される
   - PORT の部分ってどう意味か雰囲気でしか分かっていなかったのでメモ
     - `port:nodePort[Protocol]` となる
     - port: nodePortと同じ
       - Source は見つかっていないが、 ClusterIP ではPort/TargetPortのみなので、port は コンテナと接続する面の Port、NodePort は hostのポート
         - ならNodePort は 1ホスト で共通だが、Port はコンテナに紐づくので重複可能では？
           - 動作はあっていた。 `32689:32689/TCP` と `32689:32690/TCP` という service object 作れた
     - nodePort: 各 Node の IP にて、記載の値でサービスを公開するためのポート
     - targetPort: 各コンテナで開いている Port、Service はこの値を使って接続してくる

```
😀 ❯❯❯ k get svc
NAME         TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)           AGE
blue         NodePort    10.100.155.121   <none>        32689:30811/TCP   23m
kubernetes   ClusterIP   10.100.0.1       <none>        443/TCP           55m
```
