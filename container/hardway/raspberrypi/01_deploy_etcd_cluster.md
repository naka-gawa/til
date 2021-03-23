## Summary

ここでやること

- etcd build
- etcd config
- verification

## etcd build (Masterのみ)

- build

arm の release は用意されていないのでローカルでetcdをbuildする

```
sudo apt install -y git && \
  git clone https://github.com/coreos/etcd.git && \
  export GOOS=linux && \
  export GOARCH=arm && \
  export GOARM=7 &&
  git clone --branch v3.4.0 https://github.com/coreos/etcd.git && \
  cd etcd && \
  go mod vendor && \
  ./build
```

bin フォルダ内に binary できるので scp で master node に持っていく

```
scp ~/k8s/etcd/bin/* admin@10.1.101.1:~/k8s
```

## etcd config

- 作成した証明書の配置

```
sudo mkdir -p /etc/etcd /var/lib/etcd && \
  sudo chmod 700 /var/lib/etcd && \
  sudo cp ca.pem kubernetes-key.pem kubernetes.pem /etc/etcd/
```

- systemd の etcd.services ファイルの作成

```
ETCD_NAME=etcd_cluster
INTERNAL_IP=10.1.101.1

cat <<EOF | sudo tee /etc/systemd/system/etcd.service
[Unit]
Description=etcd
Documentation=https://github.com/coreos

[Service]
Type=notify
ExecStart=/usr/local/bin/etcd \\
  --name ${ETCD_NAME} \\
  --cert-file=/etc/etcd/kubernetes.pem \\
  --key-file=/etc/etcd/kubernetes-key.pem \\
  --peer-cert-file=/etc/etcd/kubernetes.pem \\
  --peer-key-file=/etc/etcd/kubernetes-key.pem \\
  --trusted-ca-file=/etc/etcd/ca.pem \\
  --peer-trusted-ca-file=/etc/etcd/ca.pem \\
  --peer-client-cert-auth \\
  --client-cert-auth \\
  --initial-advertise-peer-urls https://${INTERNAL_IP}:2380 \\
  --listen-peer-urls https://${INTERNAL_IP}:2380 \\
  --listen-client-urls https://${INTERNAL_IP}:2379,https://127.0.0.1:2379 \\
  --advertise-client-urls https://${INTERNAL_IP}:2379 \\
  --initial-cluster-token etcd-initial-token \\
  --initial-cluster ${ETCD_NAME}=https://${INTERNAL_IP}:2380 \\
  --initial-cluster-state new \\
  --data-dir=/var/lib/etcd
Restart=on-failure
RestartSec=5
Environment=ETCD_UNSUPPORTED_ARCH=arm

[Install]
WantedBy=multi-user.target
EOF
```

- etcd プロセスの起動

```
sudo systemctl daemon-reload && \
  sudo systemctl enable --now etcd.service
```

## verification

- 起動状態の確認

```
sudo ETCDCTL_API=3 etcdctl member list \
  --endpoints=https://127.0.0.1:2379 \
  --cacert=/etc/etcd/ca.pem \
  --cert=/etc/etcd/kubernetes.pem \
  --key=/etc/etcd/kubernetes-key.pem
```

ref.[how-to-create-cluster-logical-hardway](https://github.com/CyberAgentHack/home-kubernetes-2020/tree/master/how-to-create-cluster-logical-hardway)

