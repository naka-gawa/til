## Summary

ここでやること

- binary 取得
- kube-apiserver config

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
wget -q --show-progress --https-only --timestamping \
  "https://storage.googleapis.com/kubernetes-release/release/v1.18.6/bin/linux/arm/kube-apiserver" && \
  chmod +x kube-apiserver && \
  sudo mv kube-apiserver /usr/local/bin/
```

今回は etcd のデータを暗号化するため、kube-apiserver に `--encryption-provider-config` で設定ファイルを渡す
まずはその設定ファイルを作成する

```
ENCRYPTION_KEY=$(head -c 32 /dev/urandom | base64)

cat > encryption-config.yaml <<EOF
kind: EncryptionConfig
apiVersion: v1
resources:
  - resources:
      - secrets
    providers:
      - aescbc:
          keys:
            - name: key1
              secret: ${ENCRYPTION_KEY}
      - identity: {}
EOF
```

証明書 `*.pem` を配置する

```
sudo mkdir -p /etc/kubernetes/config /var/lib/kubernetes/ && \
  sudo cp -ai ca.pem ca-key.pem kubernetes-key.pem kubernetes.pem \
  service-account-key.pem service-account.pem /var/lib/kubernetes/ && \
  sudo cp -ai encryption-config.yaml /var/lib/kubernetes/
```

kube-apiVersion の service ファイルを作成し起動する

```
INTERNAL_IP=10.1.101.1
CLUSTER_IP_NETWORK=10.32.0.0/24

cat <<EOF | sudo tee /etc/systemd/system/kube-apiserver.service
[Unit]
Description=Kubernetes API Server
Documentation=https://github.com/kubernetes/kubernetes

[Service]
ExecStart=/usr/local/bin/kube-apiserver \\
  --advertise-address=${INTERNAL_IP} \\
  --allow-privileged=true \\
  --apiserver-count=3 \\
  --audit-log-maxage=30 \\
  --audit-log-maxbackup=3 \\
  --audit-log-maxsize=100 \\
  --audit-log-path=/var/log/audit.log \\
  --authorization-mode=Node,RBAC \\
  --bind-address=0.0.0.0 \\
  --client-ca-file=/var/lib/kubernetes/ca.pem \\
  --enable-admission-plugins=NamespaceLifecycle,NodeRestriction,LimitRanger,ServiceAccount,DefaultStorageClass,ResourceQuota \\
  --etcd-cafile=/var/lib/kubernetes/ca.pem \\
  --etcd-certfile=/var/lib/kubernetes/kubernetes.pem \\
  --etcd-keyfile=/var/lib/kubernetes/kubernetes-key.pem \\
  --etcd-servers=https://${INTERNAL_IP}:2379 \\
  --event-ttl=1h \\
  --encryption-provider-config=/var/lib/kubernetes/encryption-config.yaml \\
  --kubelet-certificate-authority=/var/lib/kubernetes/ca.pem \\
  --kubelet-client-certificate=/var/lib/kubernetes/kubernetes.pem \\
  --kubelet-client-key=/var/lib/kubernetes/kubernetes-key.pem \\
  --kubelet-https=true \\
  --runtime-config='api/all=true' \\
  --service-account-key-file=/var/lib/kubernetes/service-account.pem \\
  --service-cluster-ip-range=${CLUSTER_IP_NETWORK} \\
  --service-node-port-range=30000-32767 \\
  --tls-cert-file=/var/lib/kubernetes/kubernetes.pem \\
  --tls-private-key-file=/var/lib/kubernetes/kubernetes-key.pem \\
  --v=2
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload && \
  sudo systemctl enable --now kube-apiserver.service
```

## verification

- systemctl status kube-apiserver の起動状態を確認する

```
admin@k8s-master-1:~/k8s $ sudo systemctl status kube-apiserver.service
● kube-apiserver.service - Kubernetes API Server
   Loaded: loaded (/etc/systemd/system/kube-apiserver.service; enabled; vendor preset: enabled)
   Active: active (running) since Tue 2021-03-23 20:10:04 GMT; 21s ago
     Docs: https://github.com/kubernetes/kubernetes
 Main PID: 14677 (kube-apiserver)
    Tasks: 14 (limit: 2063)
   Memory: 16.3M
   CGroup: /system.slice/kube-apiserver.service
-snip-
```

ref.[how-to-create-cluster-logical-hardway](https://github.com/CyberAgentHack/home-kubernetes-2020/tree/master/how-to-create-cluster-logical-hardway)

