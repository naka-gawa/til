# Taint と Toleration
- Taint
  - Taintは汚れという意味
- Toleration
  - Tolerationは容認という意味
- Taintが識別子、Tolerationが識別子でSchedulingするような動き

## 動作確認
- 現状確認
```
😀 ❯❯❯ k get node
NAME                 STATUS   ROLES    AGE    VERSION
kind-control-plane   Ready    master   171m   v1.16.3
kind-worker          Ready    <none>   171m   v1.16.3
kind-worker2         Ready    <none>   171m   v1.16.3
```

- kind-workerにstage=prod、kind-worker2にstage=devのtaintをつける
```
k taint nodes kind-worker stage=prod:NoScheduling
k taint nodes kind-worker2 stage=dev:NoScheduling
```

- 確認
```
😀 ❯❯❯ k describe node | grep Taint
Taints:             node-role.kubernetes.io/master:NoSchedule
Taints:             stage=prod:NoSchedule
Taints:             stage=dev:NoSchedule
```
- toleration Podを作成して動作を確認する
```
$ k apply -f - << EOF
apiVersion: v1
kind: Pod
metadata:
  name: nginx-prod
spec:
  containers:
  - name: nginx-prod
    image: nginx
  tolerations:
  - key: "stage"
    operator: "Equal"
    value: "prod"
    effect: "NoSchedule"
EOF
$ k apply -f - << EOF
apiVersion: v1
kind: Pod
metadata:
  name: nginx-dev
spec:
  containers:
  - name: nginx-dev
    image: nginx
  tolerations:
  - key: "stage"
    operator: "Equal"
    value: "dev"
    effect: "NoSchedule"
EOF
```

- 確認
  - nginx-prodはkind-workerでSchedulingされていれば良し
  - nginx-devはkind-worker2でSchedulingされていれば良し
  - taint nodes xxx stage=hogehoge:NoScheduleのSyntax
    - stage=hogehogeのPodをScheduleする:それ以外はNoSchedule と覚えることにする
```
😀 ❯❯❯ k get pod -o wide
NAME         READY   STATUS    RESTARTS   AGE     IP           NODE           NOMINATED NODE   READINESS GATES
nginx-dev    1/1     Running   0          10s     10.244.2.6   kind-worker2   <none>           <none>
nginx-prod   1/1     Running   0          2m30s   10.244.1.4   kind-worker    <none>           <none>
```
