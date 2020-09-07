# Static Pod
## Overview
- StaticPod は api server を経由せず、特定のノード上の kubelet daemon によって直接管理される。
- Controler Plane によって管理される Pod とは違う。
- StaticPod 常に特定のノード上の一つの kubelet にバインドされる。
- kubelet は StaticPod ごとに api server に Podを自動作成しようとする。
  - DaemonSetの特定ノードのみのようなイメージ
- api server の管理外のため、kubectl で利用しての制御はできない。

## 試してガッテン
- Node 構成と Pod 状態
  - 1 Master / 2 Worker
  - Pod が無いことを確認
```
😀 ❯❯❯ k get node
NAME                   STATUS   ROLES    AGE     VERSION   INTERNAL-IP   EXTERNAL-IP   OS-IMAGE                                  KERNEL-VERSION     CONTAINER-RUNTIME
sample-control-plane   Ready    master   5m3s    v1.15.3   172.17.0.2    <none>        Ubuntu Disco Dingo (development branch)   4.9.184-linuxkit   containerd://1.2.6-0ubuntu1
sample-worker          Ready    <none>   4m17s   v1.15.3   172.17.0.3    <none>        Ubuntu Disco Dingo (development branch)   4.9.184-linuxkit   containerd://1.2.6-0ubuntu1
sample-worker2         Ready    <none>   4m17s   v1.15.3   172.17.0.4    <none>        Ubuntu Disco Dingo (development branch)   4.9.184-linuxkit   containerd://1.2.6-0ubuntu1

😀 ❯❯❯ k get pod
No resources found.

```

- `sample-worker` にログインし、StaticPodPathを確認する
  - `/etc/kubernetes/manifests`であることを確認
```
😀 ❯❯❯ docker exec -it $(docker ps -a | grep sample-worker2 | awk '{print $1}') /bin/bash 00:45:17  20-01-21 14:00:17  root@sample-worker2:/# grep staticPodPath /var/lib/kubelet/config.yaml
staticPodPath: /etc/kubernetes/manifests
root@sample-worker2:/#
```

- StaticPodを作成して動作を確認してみる
```
root@sample-worker2:/etc/kubernetes/manifests# cat <<EOF> static-web.yaml
apiVersion: v1
kind: Pod
metadata:
  name: static-web
  labels:
    role: myrole
spec:
  containers:
    - name: web
      image: nginx
      ports:
        - name: web
          containerPort: 80
          protocol: TCP
EOF
```

- 確認
  - `kubectl` を使わず `Pod` を `sample-worker2` ノードで起動できた
```
😀 ❯❯❯ k get pod -o wide
NAME                        READY   STATUS    RESTARTS   AGE   IP           NODE             NOMINATED NODE   READINESS GATES
static-web-sample-worker2   1/1     Running   0          41s   10.244.2.3   sample-worker2   <none>           <none>
```

## 補足事項
- `/etc/kubernetes/manifest` で無くとも良いが、変更後は kubeletの restart が必要
- `/var/lib/kubelet/config.json` に StaticPodPath が無いこともあり得る
  - 上記に StaticPodPath が存在するが、実際に Path がない場合はkubeletの再起動は不要
- リソース削除するには `api server` の管理下に居ないので、StaticPodファイルを削除すれば良い

