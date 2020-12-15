# Terraform Custom Provider

下記記事のやってみたログ
[Terraform Provider実装 入門](https://febc-yamamoto.hatenablog.jp/entry/terraform-custom-provider-01)

## Provider とは ？
- プロバイダは2種類ある
  - HashiCorp が配布しているやつ
  - サードパーティにより配布されているやつ
- プロバイダの実行ファイルは下記ルール
  - HashiCorpが配布しているプロバイダーはterraform init で自動的にインストールされる
  - サードパーティが配布しているプロバイダは以下のディレクトリから検索される
    - カレントディレクトリ
    - Terraform 本体と同じディレクトリ
    - ホームディレクトリ配下の .terraform.d/plugins/ディレクトリ
  - 1.11 だと

- HashiCorp配布プロバイダ
```
❯ echo "provider arukas{}" > test.tf

❯ tf init

Initializing the backend...

Initializing provider plugins...
- Checking for available provider plugins...
- Downloading plugin for provider "arukas" (terraform-providers/arukas) 1.1.0...

The following providers do not have any version constraints in configuration,
so the latest version was installed.

To prevent automatic upgrades to new major versions that may contain breaking
changes, it is recommended to add version = "..." constraints to the
corresponding provider blocks in configuration, with the constraint strings
suggested below.

* provider.arukas: version = "~> 1.1"

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.

❯ tree -a .
.
├── .terraform
│   └── plugins
│       └── darwin_amd64
│           ├── lock.json
│           └── terraform-provider-arukas_v1.1.0_x4
├── .terraform-version
└── test.tf

3 directories, 4 files
```

- サードパーティ配布プロバイダ
```
❯ echo "provider aws{}" > test.tf

❯ echo "0.12.0" > .terraform-version

❯ tf init

Initializing the backend...

Initializing provider plugins...
- Checking for available provider plugins...
- Downloading plugin for provider "aws" (hashicorp/aws) 3.5.0...

The following providers do not have any version constraints in configuration,
so the latest version was installed.

To prevent automatic upgrades to new major versions that may contain breaking
changes, it is recommended to add version = "..." constraints to the
corresponding provider blocks in configuration, with the constraint strings
suggested below.

* provider.aws: version = "~> 3.5"

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.

❯ tree -a .
.
├── .terraform
│   └── plugins
│       └── darwin_amd64
│           ├── lock.json
│           └── terraform-provider-aws_v3.5.0_x5
├── .terraform-version
└── test.tf

```

- terraform-provider-aws は thirdparty ながら plugins に自動配備される




