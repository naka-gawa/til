# nodeAffinity
- Affinity = `相性`
- Podを特定のnodeへScheduleするための仕組み

## 動作検証
- 現状のnode labelを確認
```
😀 ❯❯❯ k get node --show-labels                                                                    20-01-17 12:54:55
NAME                   STATUS   ROLES    AGE     VERSION   LABELS
sample-control-plane   Ready    master   9m45s   v1.15.3   beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,kubernetes.io/arch=amd64,kubernetes.io/hostname=sample-control-plane,kubernetes.io/os=linux,node-role.kubernetes.io/master=
sample-worker          Ready    <none>   9m5s    v1.15.3   beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,kubernetes.io/arch=amd64,kubernetes.io/hostname=sample-worker,kubernetes.io/os=linux
sample-worker2         Ready    <none>   9m5s    v1.15.3   beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,kubernetes.io/arch=amd64,kubernetes.io/hostname=sample-worker2,kubernetes.io/os=linux
```
- sample-worker に label:color=blueを追加, sample-worker2 に label:color=redを追加する
```
k label node sample-worker color=blue
k label node sample-worker2 color=red
😀 ❯❯❯ k get node --show-labels
NAME                   STATUS   ROLES    AGE   VERSION   LABELS
sample-control-plane   Ready    master   13m   v1.15.3   beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,kubernetes.io/arch=amd64,kubernetes.io/hostname=sample-control-plane,kubernetes.io/os=linux,node-role.kubernetes.io/master=
sample-worker          Ready    <none>   13m   v1.15.3   beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,color=blue,kubernetes.io/arch=amd64,kubernetes.io/hostname=sample-worker,kubernetes.io/os=linux
sample-worker2         Ready    <none>   13m   v1.15.3   beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,color=red,kubernetes.io/arch=amd64,kubernetes.io/hostname=sample-worker2,kubernetes.io/os=linux
```

- deploy blue を作って, sample-workerにdeployする
```
k run blue --image=nginx --labels='color=blue'
😀 ❯❯❯ k get pod -o wide
NAME                    READY   STATUS    RESTARTS   AGE   IP           NODE             NOMINATED NODE   READINESS GATES
blue-6fcdf8f889-7hgjh   1/1     Running   0          16s   10.244.1.2   sample-worker2   <none>           <none>
```

- sample-worker2 にSchedulingされてしまっているので, nodeAffinityを追加する
```
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: color
            operator: In
            values:
            - blue
```

- sample-worker にSchedulingしていることを確認
```
😀 ❯❯❯ k get pod -o widei
NAME                    READY   STATUS        RESTARTS   AGE     IP           NODE             NOMINATED NODE   READINESS GATES
blue-6fcdf8f889-7hgjh   1/1     Terminating   0          4m39s   10.244.1.2   sample-worker2   <none>           <none>
blue-78f5b768bc-grfzd   1/1     Running       0          9s      10.244.2.4   sample-worker    <none>           <none>
```

