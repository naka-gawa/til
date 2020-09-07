# namespaceを用いたServiceDiscoveryを検証する
- `[service name].[namespace name].svc.cluster.local` が利用できる
- ちなみに同一 namespace であれば`[service name]` が利用できる

## 動作確認
下記で環境を用意
```
$ k run nginx --restart=Never --image=nginx
$ k create ns test
$ k run nginx1 --restart=Never --image=nginx --labels='app=nginx' -n test
$ k run nginx2 --restart=Never --image=nginx --labels='app=nginx' -n test
$ k apply -f -n test - << EOF
apiVersion: v1
kind: Service
metadata:
  labels:
    app: nginx
  name: nginx
spec:
  ports:
  - name: "80"
    port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: nginx
  type: ClusterIP
EOF
$ k exec nginx1 -n test -- cp /etc/hostname /usr/share/nginx/html/index.html
$ k exec nginx2 -n test -- cp /etc/hostname /usr/share/nginx/html/index.html
```

こんな感じのリソースが出来上がる
- namespace test
  - Pod : nginx
- namespace : test
  - Pod : nginx
  - Service : nginx

これで default namespace から test namespace へアクセスする際にDNSを利用してみる

- namespaceマタギの場合
```
$ k exec -it nginx /bin/bash
--snip--
root@nginx:/# curl nginx.test.svc.cluster.local
nginx1
root@nginx:/# curl nginx.test.svc.cluster.local
nginx2
root@nginx:/#
```

- 同一namespaceの場合
```
$ k exec -it nginx1 -n test /bin/bash
--snip--
root@nginx1:/# curl nginx
nginx2
root@nginx1:/# curl nginx
nginx1
root@nginx1:/#
```
