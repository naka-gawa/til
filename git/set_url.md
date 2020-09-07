# git remote set-url

https://qiita.com/hirotatsuuu/items/d9abff1296469f9c861c

コントリビュートの練習として複数アカウントでpushしようとしたら403エラーがでた。

```
❯❯❯ git push -u origin master
remote: Permission to tmnkgwa4/til.git denied to xxxx.
fatal: unable to access 'https://github.com/tmnkgwa4/til.git/': The requested URL returned error: 403
```

原因はhttps://github.comのユーザ名がxxxxになっているからの模様。下記で対処した。

```
❯❯❯ git remote set-url https://tmnkgwa4@github.com/tmnkgwa4/til.git
❯❯❯ git push -u origin master
```
