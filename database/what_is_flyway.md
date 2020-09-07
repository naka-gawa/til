

## OverView
- FlyWay は データベース マイグレーションツールである
  - この DB Migration とは、移行を意味するものではなく、DB 環境 (スキーマ+データ) を移行して、同一状態のDB環境を構築するもの

## 動作確認
- docker-compose.yaml を作成する
```
😀 ❯❯❯ cat <<EOF> docker-compose.yaml
version: '3.4'

x-template: &flyway-template
  image: boxfuse/flyway:latest
  volumes:
    - ./sql:/flyway/sql # マイグレーション用SQLファイルの格納先
    - ./conf:/flyway/conf # 設定ファイルの格納先
  depends_on:
    - db

services:
  flyway-clean:
    <<: *flyway-template
    command: clean

  flyway-migrate:
    <<: *flyway-template
    command: migrate

  flyway-info:
    <<: *flyway-template
    command: info

  db:
    image: postgres:latest
    environment:
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=pass
    ports:
      - 5432:5432
    volumes:
      - ./init:/docker-entrypoint-initdb.d # Create DataBase するinit用のSQLの格納先
    container_name: db
EOF
```

- `./conf` フォルダ配下に DB 接続情報ファイルを保存する
```
😀 ❯❯❯ cat <<EOF> ./conf/flyway.conf
flyway.url=jdbc:postgresql://db:5432/test_db
flyway.user=user
flyway.password=pass
EOF
```

- `./init` フォルダ配下に DB 初期設定 SQL ファイルを保存する
```
😀 ❯❯❯ cat <<EOF> ./init/CREATE_DB.sql
CREATE DATABASE test_db;
EOF
```

- `./migration` フォルダ配下に Migration 用 SQL ファイルを保存する
  - V1はCREATE TABLE
  - V2はINSERT DATA
  - ファイル名には[Syntax](https://flywaydb.org/documentation/migrations)があるので注意
```
😀 ❯❯❯ cat <<EOF> ./sql/V1.0__CREATE_TABLE.sql
CREATE TABLE IF NOT EXISTS test_table(test_id int, test_name varchar(255));
EOF

😀 ❯❯❯ cat <<EOF> ./sql/V2.0__INSERT_DATA.sql
INSERT INTO test_table(test_id, test_name) VALUES (1, 'test');
EOF
```

- まずは db コンテナを起動して、database に何もテーブルが存在しないことを確認
```
😀 ❯❯❯ docker-compose up -d db
Creating network "database_default" with the default driver
Creating db ... done

😀 ❯❯❯ docker ps -a
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                    NAMES
23f610fc4504        postgres:latest     "docker-entrypoint.s…"   4 seconds ago       Up 3 seconds        0.0.0.0:5432->5432/tcp   db

😀 ❯❯❯ docker exec -it db /bin/bash
root@23f610fc4504:/# psql -U user -d test_db
psql (12.1 (Debian 12.1-1.pgdg100+1))
Type "help" for help.

test_db=# \l
                             List of databases
   Name    | Owner | Encoding |  Collate   |   Ctype    | Access privileges
-----------+-------+----------+------------+------------+-------------------
 postgres  | user  | UTF8     | en_US.utf8 | en_US.utf8 |
 template0 | user  | UTF8     | en_US.utf8 | en_US.utf8 | =c/user          +
           |       |          |            |            | user=CTc/user
 template1 | user  | UTF8     | en_US.utf8 | en_US.utf8 | =c/user          +
           |       |          |            |            | user=CTc/user
 test_db   | user  | UTF8     | en_US.utf8 | en_US.utf8 |                  <<< test_dbが作成されている
 user      | user  | UTF8     | en_US.utf8 | en_US.utf8 |
(5 rows)

test_db=# \d
Did not find any relations.                                                <<< DBの中身は何も作られていない
test_db=#
```

- flyway migrate を実行して、先程作成した空のデータベースの中に table と data が入っていることを確認
```
😀 ❯❯❯ docker exec -it db /bin/bash
root@23f610fc4504:/# psql -U user -d test_db
psql (12.1 (Debian 12.1-1.pgdg100+1))
Type "help" for help.

test_db=# \l
                             List of databases
   Name    | Owner | Encoding |  Collate   |   Ctype    | Access privileges
-----------+-------+----------+------------+------------+-------------------
 postgres  | user  | UTF8     | en_US.utf8 | en_US.utf8 |
 template0 | user  | UTF8     | en_US.utf8 | en_US.utf8 | =c/user          +
           |       |          |            |            | user=CTc/user
 template1 | user  | UTF8     | en_US.utf8 | en_US.utf8 | =c/user          +
           |       |          |            |            | user=CTc/user
 test_db   | user  | UTF8     | en_US.utf8 | en_US.utf8 |
 user      | user  | UTF8     | en_US.utf8 | en_US.utf8 |
(5 rows)

test_db=# \d
               List of relations
 Schema |         Name          | Type  | Owner
--------+-----------------------+-------+-------
 public | flyway_schema_history | table | user
 public | test_table            | table | user  <<< yes !
(2 rows)

test_db=# select * from test_table;
 test_id | test_name
---------+-----------
       1 | test      <<< yes !!!
(1 row)

test_db=#
```

## Reference
- [Flyway + PostgreSQLのDBマイグレーション環境をDockerを使って構築する](https://qiita.com/supimen89/items/1008e633f6ac2028e1e9)

