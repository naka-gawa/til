# forkしたリポジトリを更新する

forkしたリポジトリを更新しようとしたらupstreamのリポジトリがアップデートされていたため困った
下記で対処した

```
0. (git switch master)
1. git remote add upstream [github upstream repo]
2. git fetch upstream
3. git merge upstream/master
4. git push -u origin master
```

# other knowlege
## git remote add
リモートリポジトリを名前付けして追加する
なので、今回のケースだと必ずしもupstreamである必要は無い

## git fetch vs git pull
fetchはリモートリポジトリからデータを取得するが、mergeはしない
対してpullはリモートリポジトリからデータを取得し、mergeする

# reference
https://techacademy.jp/magazine/10268
https://qiita.com/ota42y/items/e082d64f3f8b424e9b7d

