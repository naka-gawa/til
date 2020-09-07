# Dockerfile Tips

## WORKDIR
### Overview
- ディレクトリがない場合は生成されてから作業ディレクトリになる
  - なので`mkdir`や`cd`が不要
  - 記述量が減り可読性が上がる?
  - こんな感じを記載を簡易化できる

```
ENV $APP_PATH = /app
RUN mkdir $APP_PATH && \
    cd $APP_PATH && \
    curl ...
```

```
ENV $APP_PATH = /app
WORKDIR $APP_PATH
RUN curl...
```

- [ref](https://qiita.com/DogFortune/items/05bf806ecbb256a823f8)

## CMD & ENTRYPOINT
### Overview




