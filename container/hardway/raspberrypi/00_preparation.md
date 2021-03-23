## Summary

ここでやること

- 環境値の決定
- raspios のインストール
- 初期セットアップ
- 証明書と kubeconfig の生成

## 環境値の決定
- hostname と node ip address

| Hostname | Node IP |
| ---- | ---- |
|  k8s-master-1  |  10.1.101.1  |
|  k8s-worker-2  |  10.1.101.2  |
|  k8s-worker-3  |  10.1.101.3  |

- Network Range

| Address Range | Cidr | Gateway |
| ---- | ---- | ---- |
| Node Network | 10.1.101.0/24 | 10.1.101.254 |
| Pod Network (k8s-master-1) | 10.2.1.0/24 | - |
| Pod Network (k8s-worker-2) | 10.2.2.0/24 | - |
| Pod Network (k8s-worker-3) | 10.2.3.0/24 | - |
| Cluster Network | 10.32.0.0/24 | - |

## raspios のインストール
- latest イメージの download (Release date: January 11th 2021)

[Raspberry Pi OS Lite](https://downloads.raspberrypi.org/raspios_lite_armhf/images/raspios_lite_armhf-2021-01-12/2021-01-11-raspios-buster-armhf-lite.zip)

- sd カードへのコピー

```
sudo diskutil umountDisk /dev/disk2 && \
  sudo dd bs=1m if=2021-01-11-raspios-buster-armhf-lite.img of=/dev/rdisk2 conv=sync && \
  sleep 10 && \
  touch /Volumes/boot/ssh && \
  sed -i -e 's/.$/h cgroup_enable=cpuset cgroup_memory=1 cgroup_enable=memory' /Volumes/boot/cmdline.txt && \
  sudo rm -rf /Volumes/boot/cmdline.txt-e && \
  sudo diskutil umountDisk /dev/disk2
```

address は node 毎に変更する

## 初期セットアップ

- Initial Setup

hostname と ip address は node 毎に読み替える

```
sudo adduser admin && \
  sudo gpasswd -a admin sudo && \
  sudo update-alternatives --set iptables /usr/sbin/iptables-legacy && \
  sudo dphys-swapfile swapoff && \
  sudo systemctl stop dphys-swapfile && \
  sudo systemctl disable dphys-swapfile && \
  sudo dpkg-reconfigure locales

export NODEIP=10.1.101.1
export HOSTNAME=k8s-master-1

sudo sh -c "echo interface eth0 >> /etc/dhcpcd.conf"
sudo sh -c "echo static ip_address=${NODEIP}/24 >> /etc/dhcpcd.conf"
sudo sh -c "echo static routers=10.1.101.254 >> /etc/dhcpcd.conf"
sudo sh -c "echo static domain_name_servers=8.8.8.8 >> /etc/dhcpcd.conf"

sudo sh -c "echo 127.0.0.1   ${HOSTNAME} >> /etc/hosts"
sudo sh -c "echo 10.1.101.1  k8s-master-1 >> /etc/hosts"
sudo sh -c "echo 10.1.101.2  k8s-worker-2 >> /etc/hosts"
sudo sh -c "echo 10.1.101.3  k8s-worker-3 >> /etc/hosts"

sudo sh -c "echo ${HOSTNAME} > /etc/hostname"

sudo shutdown -h now
```

- Package Install (全 node 共通)

```
sudo apt update -y && \
  sudo apt full-upgrade -y && \
  sudo apt-get --purge remove vim-common vim-tiny -y && \
  sudo apt-get -y install vim autoconf automake libtool curl unzip gcc make libseccomp2 libseccomp-dev btrfs-progs libbtrfs-dev runc golang-cfssl
```

一緒に runc もインストールしちゃう

- Golang Setup

```
wget https://golang.org/dl/go1.15.8.linux-armv6l.tar.gz && \
  sudo tar -C /usr/local -xzf go1.15.8.linux-armv6l.tar.gz && \
  echo 'PATH="$PATH:/usr/local/go/bin"' | tee -a $HOME/.profile && \
  source $HOME/.profile && \
  echo "export GOPATH=$(go env GOPATH)" | tee -a $HOME/.profile && \
  source $HOME/.profile
```

## kubectl のインストール(Masterのみ)

```
wget https://storage.googleapis.com/kubernetes-release/release/v1.18.6/bin/linux/arm/kubectl && \
  chmod +x kubectl && \
  sudo mv kubectl /usr/local/bin

mkdir k8s && cd k8s
wget https://raw.githubusercontent.com/naka-gawa/til/master/container/hardway/raspberrypi/generate-cert.sh
wget https://raw.githubusercontent.com/naka-gawa/til/master/container/hardway/raspberrypi/generate-kubeconfig.sh
chmod u+x *.sh
```

## 証明書と kubeconfig の生成
CyberAgentHackで公開頂いているものを少し改変し流用した

- generate ssl

ca の雛形ファイル、有効期限とprofile 毎の証明書の用途を定義する
```ca-config.json
{
  "signing": {
    "default": {
      "expiry": "8760h"
    },
    "profiles": {
      "kubernetes": {
        "usages": ["signing", "key encipherment", "server auth", "client auth"],
        "expiry": "8760h"
      }
    }
  }
}
```

csr `Certificate Signing Request` 認証局(CA)に対し、SSLサーバ証明書への署名を申請するファイル
公開鍵とその所有者情報、及び申請者が対応する秘密鍵を持っていることを示すために申請者の署名が記載されている
```
{
  "CN": "Kubernetes",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "JP",
      "L": "Setagaya",
      "O": "Kubernetes",
      "OU": "HOME",
      "ST": "Tokyo"
    }
  ]
}
```

上記の2種ファイルを使って ssl 証明書を発行していく
```
cfssl gencert \
  -ca=ca.pem \
  -ca-key=ca-key.pem \
  -config=ca-config.json \
  -hostname=${NODE3_HOSTNAME},${NODE3_ADDRESS} \
  -profile=kubernetes \
  ${NODE3_HOSTNAME}-csr.json | cfssljson -bare ${NODE3_HOSTNAME}
```

- kubeconfig
```
kubectl config set-cluster kubernetes-the-hard-way \
  --certificate-authority=${CERT_DIR}/ca.pem \
  --embed-certs=true \
  --server=https://${MASTER_ADDRESS}:6443 \
  --kubeconfig=kube-proxy.kubeconfig
```

- できた証明書とkubeconfig を各 worker に送る

```
ssh admin@10.1.101.2 mkdir k8s
scp ~/k8s/cert/k8s-worker-2*.pem admin@10.1.101.2:~/k8s
scp ~/k8s/kubeconfig/k8s-worker-2*.kubeconfig admin@10.1.101.2:~/k8s
scp ~/k8s/kubeconfig/kube-proxy.kubeconfig admin@10.1.101.2:~/k8s
```

```
ssh admin@10.1.101.3 mkdir k8s
scp ~/k8s/cert/k8s-worker-3*.pem admin@10.1.101.3:~/k8s
scp ~/k8s/kubeconfig/k8s-worker-3*.kubeconfig admin@10.1.101.3:~/k8s
scp ~/k8s/kubeconfig/kube-proxy.kubeconfig admin@10.1.101.3:~/k8s
```

ref.[how-to-create-cluster-logical-hardway](https://github.com/CyberAgentHack/home-kubernetes-2020/tree/master/how-to-create-cluster-logical-hardway)
