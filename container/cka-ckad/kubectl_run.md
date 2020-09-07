# run/create/expose の run

## kubectl run
### Pod を oneliner で作成する
```
😀 ❯❯❯ k run nginx --image=nginx --restart=Never
pod/nginx created

😀 ❯❯❯ k get pod
NAME    READY   STATUS    RESTARTS   AGE
nginx   1/1     Running   0          40s
```

### Pod に env を指定して作成する
```
😀 ❯❯❯ k run nginx --image=nginx --restart=Never --env=TESTENV=testvalue
pod/nginx created

😀 ❯❯❯ k exec nginx env | grep test
TESTENV=testvalue
```

## Pod と containerPort を指定して作成する
```
😀 ❯❯❯ k run nginx --restart=Never --image=nginx --port=8080
pod/nginx created

😀 ❯❯❯ k get pod -o json | grep Port
                                "containerPort": 8080,
```

## Pod と Resouce制限を加えて作成する
```
😀 ❯❯❯ k run nginx --restart=Never --image=nginx --requests='cpu=100m,memory=50Mi' --limits='cpu=200m,memory=100Mi'
pod/nginx created

😀 ❯❯❯ k get pod -o json | jq -r '.items[].spec.containers[].resources'
{
  "limits": {
    "cpu": "200m",
    "memory": "100Mi"
  },
  "requests": {
    "cpu": "100m",
    "memory": "50Mi"
  }
}
```
## Pod を 一時的に立てておいてリソース完了後に削除する
```
😀 ❯❯❯ k run alpine --restart=Never --image=alpine --rm -it --command -- ls -la
otal 64
drwxr-xr-x    1 root     root          4096 Jan 15 04:18 .
drwxr-xr-x    1 root     root          4096 Jan 15 04:18 ..
drwxr-xr-x    2 root     root          4096 Dec 24 15:04 bin
drwxr-xr-x    5 root     root           380 Jan 15 04:18 dev
drwxr-xr-x    1 root     root          4096 Jan 15 04:18 etc
drwxr-xr-x    2 root     root          4096 Dec 24 15:04 home
drwxr-xr-x    5 root     root          4096 Dec 24 15:04 lib
drwxr-xr-x    5 root     root          4096 Dec 24 15:04 media
drwxr-xr-x    2 root     root          4096 Dec 24 15:04 mnt
drwxr-xr-x    2 root     root          4096 Dec 24 15:04 opt
dr-xr-xr-x  271 root     root             0 Jan 15 04:18 proc
drwx------    2 root     root          4096 Dec 24 15:04 root
drwxr-xr-x    1 root     root          4096 Jan 15 04:18 run
drwxr-xr-x    2 root     root          4096 Dec 24 15:04 sbin
drwxr-xr-x    2 root     root          4096 Dec 24 15:04 srv
dr-xr-xr-x   13 root     root             0 Jan 15 04:18 sys
drwxrwxrwt    2 root     root          4096 Dec 24 15:04 tmp
drwxr-xr-x    7 root     root          4096 Dec 24 15:04 usr
drwxr-xr-x   12 root     root          4096 Dec 24 15:04 var
pod "alpine" deleted

😀 ❯❯❯ k get pod
No resources found.
```

## Reference
- [kubectl run/create/expose のススメ](https://qiita.com/sourjp/items/f0c8c8b4a2a494a80908)
