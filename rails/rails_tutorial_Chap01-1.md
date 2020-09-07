# Rails Tutorial Chapter 1-1
- Rails チュートリアル の備忘録

## Environment
- Docker CE
```
😀 ❯❯❯ docker version`
Client: Docker Engine - Community
 Version:           19.03.5
 API version:       1.40
 Go version:        go1.12.12
 Git commit:        633a0ea
 Built:             Wed Nov 13 07:22:34 2019
 OS/Arch:           darwin/amd64
 Experimental:      false

Server: Docker Engine - Community
 Engine:
  Version:          19.03.5
  API version:      1.40 (minimum version 1.12)
  Go version:       go1.12.12
  Git commit:       633a0ea
  Built:            Wed Nov 13 07:29:19 2019
  OS/Arch:          linux/amd64
  Experimental:     false
 containerd:
  Version:          v1.2.10
  GitCommit:        b34a5c8af56e510852c35414db4c1f4fa6172339
 runc:
  Version:          1.0.0-rc8+dev
  GitCommit:        3e425f80a8c931f88e6d94a8c831b9d5aa481657
 docker-init:
  Version:          0.18.0
  GitCommit:        fec3683
```

- docker-compose
```
😀 ❯❯❯ docker-compose version
docker-compose version 1.25.2, build 698e2846
docker-py version: 4.1.0
CPython version: 3.7.5
OpenSSL version: OpenSSL 1.1.1d  10 Sep 2019
```

- Cloud9 を使うチュートリアルになっているが、docker をフル活用する

### How to create environment for development.
#### Container Settings
- [ref](https://docs.docker.com/compose/rails/)
- まず Dockerfile を用意する
  - Quick Start では postgresql を利用する形になっているので、mysql に変更を加えている
```
FROM ruby:2.5-alpine
LABEL maintainer="Tomoaki Nakagawa <tmnkgwa4@gmail.com>"

RUN apk add --no-cache nodejs mysql-dev build-base libxml2-dev libxslt-dev
WORKDIR /myapp
COPY Gemfile /myapp/Gemfile
COPY Gemfile.lock /myapp/Gemfile.lock
RUN bundle install
COPY . /myapp

# Add a script to be executed every time the container starts.
COPY entrypoint.sh /usr/bin/
RUN chmod +x /usr/bin/entrypoint.sh
ENTRYPOINT ["entrypoint.sh"]
EXPOSE 3000

# Start the main process.
CMD ["rails", "server", "-b", "0.0.0.0"]
```

- 次に Gemfile と Gemfile.lock を作成
  - Gemfile.lock は Empty でおｋ
```
😀 ❯❯❯ cat Gemfile
source 'https://rubygems.org'
gem 'rails', '~>5'

😀 ❯❯❯ cat Gemfile.lock

😀 ❯❯❯
```

- そして、Entrypoint.sh を作成
  - `server.pid` を削除する必要があるらしい（調査中）
```
😀 ❯❯❯ cat entrypoint.sh
#!/bin/bash
set -e

# Remove a potentially pre-existing server.pid for Rails.
rm -f /myapp/tmp/pids/server.pid

# Then exec the container's main process (what's set as CMD in the Dockerfile).
exec "$@"

😀 ❯❯❯
```

- 最後に`docker-compose.yaml`を作成する
  - db user は root, pw は password
  - database name は mysql
```
version: '3'
services:
  db:
    image: mysql:5.7
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: root
    ports:
      - "3306:3306"
  web:
    build: .
    command: bash -c "rm -f tmp/pids/server.pid && bundle exec rails s -p 3000 -b '0.0.0.0'"
    volumes:
      - .:/myapp
    ports:
      - "3000:3000"
    depends_on:
      - db
```

#### building environment
- Building the project
  - まず最初にプロジェクトを作成する
    - やっていることは Dockerfile の内容を Build して rails new でプロジェクトを current directory に作成している
  - `--build-skip` でイメージのビルドを飛ばす
  - `-G` オプションで git init を skip している
```
😀 ❯❯❯ docker-compose run web rails new . --force --database=mysql --skip-bundle -G

Creating network "develop-environment_default" with the default driver

[snip]

```

- すると、rails project が作成されるので `config/database.yaml` を編集する
```
default: &default
  adapter: mysql2
  username: root
  password: password
  host: db
  pool: 5

development:
  <<: *default
  database: dev_db

test:
  <<: *default
  database: test_db

production:
  <<: *default
  database: prd_db
```

- 改めて イメージ を ビルドする
```
😀 ❯❯❯ docker-compose build

[snip]

Successfully tagged develop-environment_web:latest

😀 ❯❯❯
```
- DB と Web を Up する
```
😀 ❯❯❯ docker-compose up -d

[snip]

😀 ❯❯❯
```
- `localhost:3000` にアクセスして確認する
  - この状態だと mysql に database が存在しないので、rails から db:create を流す
```
😀 ❯❯❯ docker-compose run web rake db:create
Starting develop-environment_db_1 ... done
Created database 'dev_db'
Created database 'test_db'
```
- 改めて `localhost:3000` にアクセスして `Yay! You're on Rails!` が表示されれば OK

## Trouble Shoot
- docker-compose run -d で立ち上げてもコンテナが立ち上がらない
  - こんなログ
  - Gemfile の `tzinfo-data` platform を削除すればよい
    - `gem 'tzinfo-data', platforms: [:mingw, :mswin, :x64_mingw, :jruby]`
    - `gem 'tzinfo-data'
```
😀 ❯❯❯ docker logs develop-environment_web_1                             20-01-29 15:53:55
/usr/local/lib/ruby/site_ruby/2.5.0/bundler/spec_set.rb:43:in `block in for': Unable to find a spec satisfying tzinfo-data (>= 0) in the set. Perhaps the lockfile is corrupted? (Bundler::GemNotFound)
        from /usr/local/lib/ruby/site_ruby/2.5.0/bundler/spec_set.rb:25:in `loop'
        from /usr/local/lib/ruby/site_ruby/2.5.0/bundler/spec_set.rb:25:in `for'
        from /usr/local/lib/ruby/site_ruby/2.5.0/bundler/spec_set.rb:83:in `materialize'
        from /usr/local/lib/ruby/site_ruby/2.5.0/bundler/definition.rb:170:in `specs'
        from /usr/local/lib/ruby/site_ruby/2.5.0/bundler/definition.rb:237:in `specs_for'
        from /usr/local/lib/ruby/site_ruby/2.5.0/bundler/definition.rb:226:in `requested_specs'
        from /usr/local/lib/ruby/site_ruby/2.5.0/bundler/runtime.rb:108:in `block in definition_method'
        from /usr/local/lib/ruby/site_ruby/2.5.0/bundler/runtime.rb:20:in `setup'
        from /usr/local/lib/ruby/site_ruby/2.5.0/bundler.rb:107:in `setup'
        from /usr/local/lib/ruby/site_ruby/2.5.0/bundler/setup.rb:20:in `<top (required)>'
        from /usr/local/lib/ruby/site_ruby/2.5.0/rubygems/core_ext/kernel_require.rb:54:in `require'
        from /usr/local/lib/ruby/site_ruby/2.5.0/rubygems/core_ext/kernel_require.rb:54:in `require'
        from /myapp/config/boot.rb:3:in `<top (required)>'
        from bin/rails:3:in `require_relative'
        from bin/rails:3:in `<main>'
```

- Alpine Ruby Image に nokogiri を bundle install 使用とするとエラーでコケる
  - [ref](https://copo.jp/blog/2016/03/alpine-の-ruby-のイメージに-nokogiri-をインストール/)
  - 原因は必須パッケージが足らないから
    - build-base, libxml2-dev, libxslt-dev を追加することでbundle install が通る
```
😀 ❯❯❯ docker build -t test .
Sending build context to Docker daemon  4.608kB
Step 1/12 : FROM ruby:2.5-alpine
 ---> cbd297d70a23
Step 2/12 : RUN apk add --no-cache nodejs mysql-dev
 ---> Running in 7949a31ca313
fetch http://dl-cdn.alpinelinux.org/alpine/v3.11/main/x86_64/APKINDEX.tar.gz
fetch http://dl-cdn.alpinelinux.org/alpine/v3.11/community/x86_64/APKINDEX.tar.gz
(1/13) Installing openssl-dev (1.1.1d-r3)
(2/13) Installing mariadb-connector-c (3.1.6-r0)
(3/13) Installing mariadb-connector-c-dev (3.1.6-r0)
(4/13) Installing mariadb-common (10.4.10-r0)
(5/13) Installing libaio (0.3.112-r1)
(6/13) Installing xz-libs (5.2.4-r0)
(7/13) Installing pcre (8.43-r0)
(8/13) Installing mariadb-embedded (10.4.10-r0)
(9/13) Installing mariadb-dev (10.4.10-r0)
(10/13) Installing c-ares (1.15.0-r0)
(11/13) Installing nghttp2-libs (1.40.0-r0)
(12/13) Installing libuv (1.34.0-r0)
(13/13) Installing nodejs (12.14.0-r0)
Executing busybox-1.31.1-r9.trigger
OK: 83 MiB in 50 packages
Removing intermediate container 7949a31ca313
 ---> 96aefc0279a5
Step 3/12 : WORKDIR /myapp
 ---> Running in 2649334ed251
Removing intermediate container 2649334ed251
 ---> 3c6504ffd9ad
Step 4/12 : COPY Gemfile /myapp/Gemfile
 ---> c5555fb16098
Step 5/12 : COPY Gemfile.lock /myapp/Gemfile.lock
 ---> 5538873d3748
Step 6/12 : RUN bundle install
 ---> Running in 1e8bb77cf57f
Fetching gem metadata from https://rubygems.org/.............
Fetching gem metadata from https://rubygems.org/.
Resolving dependencies...
Fetching rake 13.0.1
Installing rake 13.0.1
Fetching concurrent-ruby 1.1.5
Installing concurrent-ruby 1.1.5
Fetching i18n 1.8.2
Installing i18n 1.8.2
Fetching minitest 5.14.0
Installing minitest 5.14.0
Fetching thread_safe 0.3.6
Installing thread_safe 0.3.6
Fetching tzinfo 1.2.6
Installing tzinfo 1.2.6
Fetching activesupport 5.2.4.1
Installing activesupport 5.2.4.1
Fetching builder 3.2.4
Installing builder 3.2.4
Fetching erubi 1.9.0
Installing erubi 1.9.0
Fetching mini_portile2 2.4.0
Installing mini_portile2 2.4.0
Fetching nokogiri 1.10.7
Installing nokogiri 1.10.7 with native extensions
Gem::Ext::BuildError: ERROR: Failed to build gem native extension.

    current directory: /usr/local/bundle/gems/nokogiri-1.10.7/ext/nokogiri
/usr/local/bin/ruby -I /usr/local/lib/ruby/site_ruby/2.5.0 -r
./siteconf20200128-1-1ulffhe.rb extconf.rb
checking if the C compiler accepts ... *** extconf.rb failed ***
Could not create Makefile due to some reason, probably lack of necessary
libraries and/or headers.  Check the mkmf.log file for more details.  You may
need configuration options.

Provided configuration options:
        --with-opt-dir
        --without-opt-dir
        --with-opt-include
        --without-opt-include=${opt-dir}/include
        --with-opt-lib
        --without-opt-lib=${opt-dir}/lib
        --with-make-prog
        --without-make-prog
        --srcdir=.
        --curdir
        --ruby=/usr/local/bin/$(RUBY_BASE_NAME)
        --help
        --clean
/usr/local/lib/ruby/2.5.0/mkmf.rb:456:in `try_do': The compiler failed to
generate an executable file. (RuntimeError)
You have to install development tools first.
        from /usr/local/lib/ruby/2.5.0/mkmf.rb:574:in `block in try_compile'
        from /usr/local/lib/ruby/2.5.0/mkmf.rb:521:in `with_werror'
        from /usr/local/lib/ruby/2.5.0/mkmf.rb:574:in `try_compile'
        from extconf.rb:138:in `nokogiri_try_compile'
        from extconf.rb:162:in `block in add_cflags'
        from /usr/local/lib/ruby/2.5.0/mkmf.rb:632:in `with_cflags'
        from extconf.rb:161:in `add_cflags'
        from extconf.rb:416:in `<main>'

To see why this extension failed to compile, please check the mkmf.log which can
be found here:

  /usr/local/bundle/extensions/x86_64-linux/2.5.0/nokogiri-1.10.7/mkmf.log

extconf failed, exit code 1

Gem files will remain installed in /usr/local/bundle/gems/nokogiri-1.10.7 for
inspection.
Results logged to
/usr/local/bundle/extensions/x86_64-linux/2.5.0/nokogiri-1.10.7/gem_make.out

An error occurred while installing nokogiri (1.10.7), and Bundler cannot
continue.
Make sure that `gem install nokogiri -v '1.10.7' --source
'https://rubygems.org/'` succeeds before bundling.

In Gemfile:
  rails was resolved to 5.2.4.1, which depends on
    actioncable was resolved to 5.2.4.1, which depends on
      actionpack was resolved to 5.2.4.1, which depends on
        actionview was resolved to 5.2.4.1, which depends on
          rails-dom-testing was resolved to 2.0.3, which depends on
            nokogiri
The command '/bin/sh -c bundle install' returned a non-zero code: 5
```

- Alpine イメージを使ったときに `standard_init_linux.go:211: exec user process caused "no such file or directory"` エラーが出て Build がコケる
  - `entrypoint.sh` が原因
  - QuickStart は Ubuntu Image なので bash だが、容量を小さくしたいため、alpine を採用したらだめだった
    - alpine は `sh`
  - 修正する前がこちら
```
#!/bin/bash
set -e

# Remove a potentially pre-existing server.pid for Rails.
rm -f /myapp/tmp/pids/server.pid

# Then exec the container's main process (what's set as CMD in the Dockerfile).
exec "$@"
```
