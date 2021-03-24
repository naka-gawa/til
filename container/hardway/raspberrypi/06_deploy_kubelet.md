## Summary

ここでやること

- 証明書とフォルダ作成
- 各種 binary 取得
  - Low Level Runtime
  - High Level Runtime
  - kubelet
  - cni plugin
- kubelet の service ファイル作成と起動

## 証明書とフォルダ作成

```
sudo mkdir -p \
  /etc/cni/net.d \
  /opt/cni/bin \
  /var/lib/kubelet \
  /var/lib/kubernetes \
  /etc/containerd

sudo cp -ai ${HOSTNAME}-key.pem ${HOSTNAME}.pem /var/lib/kubelet/
sudo cp -ai ${HOSTNAME}.kubeconfig /var/lib/kubelet/kubeconfig
```

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
  https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.18.0/crictl-v1.18.0-linux-arm.tar.gz \
  https://github.com/containernetworking/plugins/releases/download/v0.8.6/cni-plugins-linux-arm-v0.8.6.tgz \
  https://storage.googleapis.com/kubernetes-release/release/v1.18.6/bin/linux/arm/kubelet

tar -xvf crictl-v1.18.0-linux-arm.tar.gz
sudo tar -xvf cni-plugins-linux-arm-v0.8.6.tgz -C /opt/cni/bin/
chmod +x crictl kubelet
sudo mv crictl kubelet /usr/local/bin/
```

containerd をインストールする

```
curl -fsSL https://download.docker.com/linux/raspbian/gpg | sudo apt-key add -
cat <<EOF | sudo tee /etc/apt/sources.list
deb http://raspbian.raspberrypi.org/raspbian/ buster main contrib non-free rpi
# Uncomment line below then 'apt-get update' to enable 'apt-get source'
#deb-src http://raspbian.raspberrypi.org/raspbian/ buster main contrib non-free rpi
deb https://download.docker.com/linux/raspbian buster stable
EOF

sudo apt update -y && \
  sudo apt install -y containerd.io && \
  sudo systemctl enable --now containerd.service
containerd --version
```

Pod Network を定義する

```
POD_CIDR=10.2.1.0/24

cat <<EOF | sudo tee /etc/cni/net.d/10-bridge.conf
{
    "cniVersion": "0.3.1",
    "name": "bridge",
    "type": "bridge",
    "bridge": "cnio0",
    "isGateway": true,
    "ipMasq": true,
    "ipam": {
        "type": "host-local",
        "ranges": [
          [{"subnet": "${POD_CIDR}"}]
        ],
        "routes": [{"dst": "0.0.0.0/0"}]
    }
}
EOF

cat <<EOF | sudo tee /etc/cni/net.d/99-loopback.conf
{
    "cniVersion": "0.3.1",
    "name": "lo",
    "type": "loopback"
}
EOF
```

containerd の設定をする

```
cat << EOF | sudo tee /etc/containerd/config.toml
[plugins]
  [plugins.cri.containerd]
    snapshotter = "overlayfs"
    [plugins.cri.containerd.default_runtime]
      runtime_type = "io.containerd.runtime.v1.linux"
      runtime_engine = "/usr/sbin/runc"
      runtime_root = ""
EOF
```

kubelet の service ファイルを作成し起動する

```
POD_CIDR=10.2.1.0/24

cat <<EOF | sudo tee /var/lib/kubelet/kubelet-config.yaml
kind: KubeletConfiguration
apiVersion: kubelet.config.k8s.io/v1beta1
authentication:
  anonymous:
    enabled: false
  webhook:
    enabled: true
  x509:
    clientCAFile: "/var/lib/kubernetes/ca.pem"
authorization:
  mode: Webhook
clusterDomain: "cluster.local"
clusterDNS:
  - "10.32.0.10"
podCIDR: "${POD_CIDR}"
resolvConf: "/run/systemd/resolve/resolv.conf"
runtimeRequestTimeout: "15m"
tlsCertFile: "/var/lib/kubelet/${HOSTNAME}.pem"
tlsPrivateKeyFile: "/var/lib/kubelet/${HOSTNAME}-key.pem"
EOF


cat <<EOF | sudo tee /etc/systemd/system/kubelet.service
[Unit]
Description=Kubernetes Kubelet
Documentation=https://github.com/kubernetes/kubernetes
After=containerd.service
Requires=containerd.service

[Service]
ExecStart=/usr/local/bin/kubelet \\
  --config=/var/lib/kubelet/kubelet-config.yaml \\
  --container-runtime=remote \\
  --container-runtime-endpoint=unix:///var/run/containerd/containerd.sock \\
  --image-pull-progress-deadline=2m \\
  --kubeconfig=/var/lib/kubelet/kubeconfig \\
  --network-plugin=cni \\
  --register-node=true \\
  --v=2
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload && \
  sudo systemctl enable --now kubelet.service
```

## verification

- systemctl status kubelet の起動状態を確認する

```
systemctl status kubelet.service
● kubelet.service - Kubernetes Kubelet
   Loaded: loaded (/etc/systemd/system/kubelet.service; enabled; vendor preset: enabled)
   Active: active (running) since Wed 2021-03-24 19:16:43 GMT; 4s ago
     Docs: https://github.com/kubernetes/kubernetes
 Main PID: 16334 (kubelet)
    Tasks: 14 (limit: 2063)
   Memory: 9.9M
   CGroup: /system.slice/kubelet.service
```

ref.[how-to-create-cluster-logical-hardway](https://github.com/CyberAgentHack/home-kubernetes-2020/tree/master/how-to-create-cluster-logical-hardway)

