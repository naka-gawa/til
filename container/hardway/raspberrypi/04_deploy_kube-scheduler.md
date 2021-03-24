## Summary

ここでやること

- binary 取得
- kube-scheduler config

## binary 取得 (Masterのみ)

- binary 取得

気にするのは cpu

```
admin@k8s-master-1:~ $ uname -a
Linux k8s-master-1 5.4.83-v7+ #1379 SMP Mon Dec 14 13:08:57 GMT 2020 armv7l GNU/Linux
admin@k8s-master-1:~ $
```

raspi 3 は armv7l なので arm のバイナリを選択すること

```
wget -q --show-progress --https-only --timestamping "https://storage.googleapis.com/kubernetes-release/release/v1.18.6/bin/linux/arm/kube-scheduler" && \
  chmod +x kube-scheduler && \
  sudo mv kube-scheduler /usr/local/bin/
```

kubeconfig を配置する

```
sudo cp -ai kube-scheduler.kubeconfig /var/lib/kubernetes/
```

kube-apiVersion の service ファイルを作成し起動する

```
cat <<EOF | sudo tee /etc/kubernetes/config/kube-scheduler.yaml
apiVersion: kubescheduler.config.k8s.io/v1alpha1
kind: KubeSchedulerConfiguration
clientConnection:
  kubeconfig: "/var/lib/kubernetes/kube-scheduler.kubeconfig"
leaderElection:
  leaderElect: true
EOF

cat <<EOF | sudo tee /etc/systemd/system/kube-scheduler.service
[Unit]
Description=Kubernetes Scheduler
Documentation=https://github.com/kubernetes/kubernetes

[Service]
ExecStart=/usr/local/bin/kube-scheduler \\
  --config=/etc/kubernetes/config/kube-scheduler.yaml \\
  --v=2
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload && \
  sudo systemctl enable --now kube-scheduler.service
```

## verification

- systemctl status kube-scheduler.service の起動状態を確認する

```
>   sudo systemctl enable --now kube-scheduler.service
Created symlink /etc/systemd/system/multi-user.target.wants/kube-scheduler.service → /etc/systemd/system/kube-scheduler.service.
admin@k8s-master-1:~/k8s/kubeconfig $ systemctl status kube-scheduler.service
● kube-scheduler.service - Kubernetes Scheduler
   Loaded: loaded (/etc/systemd/system/kube-scheduler.service; enabled;
   Active: active (running) since Wed 2021-03-24 17:29:01 GMT; 11s ago
     Docs: https://github.com/kubernetes/kubernetes
 Main PID: 16616 (kube-scheduler)
    Tasks: 10 (limit: 2063)
   Memory: 7.3M
   CGroup: /system.slice/kube-scheduler.service
```

ref.[how-to-create-cluster-logical-hardway](https://github.com/CyberAgentHack/home-kubernetes-2020/tree/master/how-to-create-cluster-logical-hardway)

