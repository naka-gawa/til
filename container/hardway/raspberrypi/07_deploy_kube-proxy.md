## Summary

ここでやること

- binary 取得
- kube-controller-manager config

## binary 取得

- binary 取得

気にするのは cpu

```
admin@k8s-master-1:~ $ uname -a
Linux k8s-master-1 5.4.83-v7+ #1379 SMP Mon Dec 14 13:08:57 GMT 2020 armv7l GNU/Linux
admin@k8s-master-1:~ $
```

raspi 3 は armv7l なので arm のバイナリを選択すること

```
wget -q --show-progress --https-only --timestamping https://storage.googleapis.com/kubernetes-release/release/v1.18.6/bin/linux/arm/kube-proxy && \
  chmod +x kube-proxy && \
  sudo mv kube-proxy /usr/local/bin/
```

kubeconfig を配置する

```
sudo mkdir -p /var/lib/kube-proxy
sudo mv kube-proxy.kubeconfig /var/lib/kube-proxy/kubeconfig
```

kube-apiVersion の service ファイルを作成し起動する

```
POD_NETWORK=10.2.0.0/16

cat <<EOF | sudo tee /var/lib/kube-proxy/kube-proxy-config.yaml
kind: KubeProxyConfiguration
apiVersion: kubeproxy.config.k8s.io/v1alpha1
clientConnection:
  kubeconfig: "/var/lib/kube-proxy/kubeconfig"
mode: "iptables"
clusterCIDR: "${POD_NETWORK}"
EOF


cat <<EOF | sudo tee /etc/systemd/system/kube-proxy.service
[Unit]
Description=Kubernetes Kube Proxy
Documentation=https://github.com/kubernetes/kubernetes

[Service]
ExecStart=/usr/local/bin/kube-proxy \\
  --config=/var/lib/kube-proxy/kube-proxy-config.yaml
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload && \
  sudo systemctl enable --now kube-proxy.service
```

## verification

- systemctl status kube-proxy の起動状態を確認する

```
admin@k8s-worker-3:~/k8s $ systemctl status kube-proxy.service
● kube-proxy.service - Kubernetes Kube Proxy
   Loaded: loaded (/etc/systemd/system/kube-proxy.service; enabled; vendor preset: enabled)
   Active: active (running) since Wed 2021-03-24 19:33:10 GMT; 7s ago
     Docs: https://github.com/kubernetes/kubernetes
 Main PID: 17807 (kube-proxy)
    Tasks: 10 (limit: 2063)
   Memory: 6.3M
   CGroup: /system.slice/kube-proxy.service
```

ref.[how-to-create-cluster-logical-hardway](https://github.com/CyberAgentHack/home-kubernetes-2020/tree/master/how-to-create-cluster-logical-hardway)

