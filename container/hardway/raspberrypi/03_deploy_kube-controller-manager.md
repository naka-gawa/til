## Summary

ここでやること

- binary 取得
- kube-controller-manager config

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
wget -q --show-progress --https-only --timestamping "https://storage.googleapis.com/kubernetes-release/release/v1.18.6/bin/linux/arm/kube-controller-manager" && \
  chmod +x kube-controller-manager && \
  sudo mv kube-controller-manager /usr/local/bin/
```

kubeconfig を配置する

```
sudo cp -ai kube-controller-manager.kubeconfig /var/lib/kubernetes/
```

kube-apiVersion の service ファイルを作成し起動する

```
POD_NETWORK=10.2.0.0/16
CLUSTER_IP_NETWORK=10.32.0.0/24

cat <<EOF | sudo tee /etc/systemd/system/kube-controller-manager.service
[Unit]
Description=Kubernetes Controller Manager
Documentation=https://github.com/kubernetes/kubernetes

[Service]
ExecStart=/usr/local/bin/kube-controller-manager \\
  --bind-address=0.0.0.0 \\
  --cluster-cidr=${POD_NETWORK} \\
  --cluster-name=kubernetes \\
  --cluster-signing-cert-file=/var/lib/kubernetes/ca.pem \\
  --cluster-signing-key-file=/var/lib/kubernetes/ca-key.pem \\
  --kubeconfig=/var/lib/kubernetes/kube-controller-manager.kubeconfig \\
  --leader-elect=true \\
  --root-ca-file=/var/lib/kubernetes/ca.pem \\
  --service-account-private-key-file=/var/lib/kubernetes/service-account-key.pem \\
  --service-cluster-ip-range=${CLUSTER_IP_NETWORK} \\
  --use-service-account-credentials=true \\
  --v=2
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload && \
  sudo systemctl enable --now kube-controller-manager.service
```

## verification

- systemctl status kube-controller-manager の起動状態を確認する

```
● kube-controller-manager.service - Kubernetes Controller Manager
   Loaded: loaded (/etc/systemd/system/kube-controller-manager.service; enabled; vendor preset: enabled)
   Active: active (running) since Tue 2021-03-23 20:21:12 GMT; 10s ago
     Docs: https://github.com/kubernetes/kubernetes
 Main PID: 14777 (kube-controller)
    Tasks: 10 (limit: 2063)
   Memory: 10.5M
   CGroup: /system.slice/kube-controller-manager.service
```

ref.[how-to-create-cluster-logical-hardway](https://github.com/CyberAgentHack/home-kubernetes-2020/tree/master/how-to-create-cluster-logical-hardway)

