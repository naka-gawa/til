# Pull Request Modify

https://teratail.com/questions/122843

PullReqestを投げたが修正依頼が入ったため、手順を調べた
通常はこんな感じのフローになる

```
1. git fork
2. git checkout -b develop
3. code modify
4. git add -> git commit -> git push
5. PR
6. reviewをしてもらう
7. LGTMをもらったならばmergeする
```

今回は6で修正依頼が入った
その場合は 1-4をもう一度実施する
ポイントはgit push のリポジトリを同名称にする必要がある
そうすることで自動的にPRの対象となっているremote branchのcommit履歴が積み重なる

```
1. git fork
2. git checkout -b develop
3. code modify
4. git add -> git commit -> git push
5. PR
6. reviewをしてもらう
7. review内容を修正し再度push
```
