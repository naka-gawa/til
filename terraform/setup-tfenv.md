# Terraformのバージョン管理にtfenvを使うと楽になった

インストール自体は楽だが、バージョンを管理したい場合に煩雑になるなぁと感じていたのでtfenvを導入した

## 導入
- brewでインストール可能
```$ brew install tfenv
```

## 使い方
- `list-remote`で確認
```$ tfenv list-remote
0.12.13
0.12.12
0.12.11
0.12.10
0.12.9
0.12.8
0.12.7
0.12.6
0.12.5
-snip-
```

- `install`でインストール
```$ tfenv uninstall 0.12.13
[INFO] Installing Terraform v0.12.13
[INFO] Downloading release tarball from https://releases.hashicorp.com/terraform/0.12.13/terraform_0.12.13_darwin_amd64.zip
#######################################################################   99.9%
[INFO] Downloading SHA hash file from https://releases.hashicorp.com/terraform/0.12.13/terraform_0.12.13_SHA256SUMS
tfenv: tfenv-install: [WARN] No keybase install found, skipping OpenPGP signature verification
Archive:  tfenv_download.2Z0wVk/terraform_0.12.13_darwin_amd64.zip
  inflating: /usr/local/Cellar/tfenv/1.0.2/versions/0.12.13/terraform
  [INFO] Installation of terraform v0.12.13 successful
  [INFO] Switching to v0.12.13
  [INFO] Switching completed
```

- `uninstall`で削除
```$ tfenv uninstall 0.12.13
[INFO] Uninstall Terraform v0.12.13
[INFO] Terraform v0.12.13 is successfully uninstalled
```
## Tips
- 実行ログから解るようにOSに合ったバイナリをダウンロードし、`/usr/local/bin`に配置しているよう。

