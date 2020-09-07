# Terraform 概要

## what is terrafrom?
- インフラストラクチャ定義ツール
- 宣言的にクラウド上のリソースを定義することができるのが特徴

## 構成要素
- HCL
  - *H*ashicorp *C*onfiguration *L*anguageの略
  - 設定言語のこと
- Resource
  - リソースは重要な構成要素
  - このブロックでAWS、GCPなどのクラウドリソースを定義する
  - なので、クラウドリソース = Resourceと読み替えても良いかも
- Datasource
  - ReadonlyなResourceのこと
  - ROなためTerraform管理外だが、Terraformから参照したいデータのことを指す
- Provider
  - Resourceの作成更新削除に実行するプラグイン
- State
  - Terraformが認識するResouce状態
  - `*.tfstate`ファイルに保存される

## 主なコマンド
- terraform init
  - Terraformを初期化する上でまず最初に実行するコマンド
  - 作業領域の初期化
  - これを実行すると`*.tf`ファイルで定義しているプラグインのダウンロード処理が始まる
  - ダウンロードファイルは`.terraform`フォルダに保存される
- terraform plan
  - apply時にどのようなリソースが作成/更新/削除されるかをdry-runするコマンド
- terraform apply
  - `*.tf`ファイルを基に、リソースの生成を行うコマンド
  - リソースが生成されると、`terraform.state`というファイルにリソースに関連する情報が保存される
  - 1世代前のものがバックアップされ、`terraform.state.backup`に保存される
- terraform show
  - ステータス確認コマンド
  - 実際には`terraform.state`ファイルの状態をステータスを表示する
- terraform fmt
  - formatを自動で修正してくれるコマンド
- terraform destroy
  - Terraformプロジェクトを削除するコマンド

## reference
[Terraform職人入門: 日々の運用で学んだ知見を淡々とまとめる](https://qiita.com/minamijoyo/items/1f57c62bed781ab8f4d7)
