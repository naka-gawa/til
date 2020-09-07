# Docker コマンド

## イメージ関連
### 指定イメージのpull
- 基本はこんな感じ
```
docker pull REPOSITORY[:TAG]
```

- 最初はローカルリポジトリを確認し、該当するイメージがなければリモートリポジトリを参照する動作
  - 何も指定しなければ、リモートリポジトリはDockerHubになる
  - ecr や gcr の場合は、こんな感じ

- ecr
  - [ref](https://docs.aws.amazon.com/ja_jp/AmazonECR/latest/userguide/docker-pull-ecr-image.html)
```
[aws_account_id].dkr.ecr.[region].amazonaws.com/[REPOSITORY][:TAG]
```

- gcr
  - [ref](https://cloud.google.com/container-registry/docs/overview?hl=ja)
  - hostnameは下記
    - gcr.io      : 米国をホストするが将来的にロケーションは変更になる可能性あり
    - us.gcr.io   : 米国ホスト
    - eu.gcr.io   : 欧州ホスト
    - asia.gcr.io : アジアホスト
```
[hostname]/[project-id]/[REPOSITORY][:TAG]
```

### イメージ一覧の取得
```
😀 ❯❯❯ docker images
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
aws-infra           latest              15a43d4fa30f        7 weeks ago         1.12GB
clusterops          latest              149b30a83791        8 weeks ago         539MB
myapp               latest              31c490c55532        8 weeks ago         5.55MB
<none>              <none>              80f7e3678503        8 weeks ago         5.55MB
memcached           latest              0860688a721c        2 months ago        82.2MB
redis               latest              dcf9ec9265e0        2 months ago        98.2MB
mysql               5.7                 1e4405fe1ea9        2 months ago        437MB
alpine              latest              965ea09ff2eb        3 months ago        5.55MB
kindest/node        v1.15.3             8ca0c8463ebe        5 months ago        1.45GB
ruby                2.6.2               8d6721e9290e        10 months ago       870MB
alpine              3.8                 dac705114996        10 months ago       4.41MB
fluent/fluentd      latest              9406ff63f205        13 months ago       38.3MB
```
## ビルド関連
### DockerfileからDockerイメージを構築もの
- 基本はこんな感じ
  - よく `docker build .` とするが、これはビルドコンテキストが Current Directoryであることを示している
    - ビルドコンテキストはファイル location set であり、PATH もしくは URLで指定される
      - URL では Git や tarboll , plaintext を 選べる
      - [ref](https://docs.docker.com/engine/reference/commandline/build/)
```
docker build [option] Path/URL
```

- ためしてガッテン
  - ビルドコンテキストは . なのでcurrent directory
  - ビルドコンテキストは 6.656kB で、それを tar で 固めて Docker daemon に送っている
    - なのでビルドコンテキストに不要な ファイルがあるとそれだけでイメージが大きくなる
    - .dockerignore で除外することもできる
      - ただし、.dockerignore 処理フローは下記のようになるのでファイル数が多くなればなるほどビルドパフォーマンスの劣化が起きる
        - .dockerignore をパースする
        - .dockerignore にかかれている文字から正規表現を作る
        - ディレクトリをトラバーサルする
          - 正規表現がディレクトリにマッチした場合、そのディレクトリをIgnore対象として、ディレクトリをスキップする
          - 正規表現がファイルにマッチした場合、そのファイルをIgnore対象とする
        - 上記を繰り返す
    - 根本的には不要なビルドコンテキストには置かない
- きれいなビルドの場合
```

😀 ❯❯❯  echo 'Hello World !!' > hello

😀 ❯❯❯  cat <<EOF> Dockerfile
FROM busybox:latest
COPY /hello /
RUN cat /hello
EOF

😀 ❯❯❯  docker build -t sampleapp:v1 .

Sending build context to Docker daemon  6.656kB
Step 1/3 : FROM busybox:latest
latest: Pulling from library/busybox
bdbbaa22dec6: Pull complete
Digest: sha256:6915be4043561d64e0ab0f8f098dc2ac48e077fe23f488ac24b665166898115a
Status: Downloaded newer image for busybox:latest
 ---> 6d5fcfe5ff17
Step 2/3 : COPY /hello /
 ---> 69d78bd44252
Step 3/3 : RUN cat /hello
 ---> Running in 5950e0a240af
Hello World !!
Removing intermediate container 5950e0a240af
 ---> 0da064e82d6d
Successfully built 0da064e82d6d
Successfully tagged sampleapp:v1

```

- 不要なファイルをおいた場合
  - 1G ファイル を含めた Current Context を Docker daemonに送っていることがわかる
```
😀 ❯❯❯ mkfile 1g test_1g_dummyfile

😀 ❯❯❯ docker build -t sampleapp:v2 .
Sending build context to Docker daemon  1.074GB
Step 1/3 : FROM busybox:latest
latest: Pulling from library/busybox
bdbbaa22dec6: Pull complete
Digest: sha256:6915be4043561d64e0ab0f8f098dc2ac48e077fe23f488ac24b665166898115a
Status: Downloaded newer image for busybox:latest
 ---> 6d5fcfe5ff17
Step 2/3 : COPY /hello /
 ---> 1fcf1022f4ba
Step 3/3 : RUN cat /hello
 ---> Running in 876615e52eef
Hello World !!
Removing intermediate container 876615e52eef
 ---> 1bd7af5c30ea
Successfully built 1bd7af5c30ea
Successfully tagged sampleapp:v2

```

## reference
- [.dockerignore アンチパターン](https://qiita.com/munisystem/items/b0f08b28e8cc26132212)
- [Dockerのビルドコンテキスト(build context)について確認したときのメモ](https://qiita.com/toshihirock/items/c85f3eb5f4752b15ca3d)
