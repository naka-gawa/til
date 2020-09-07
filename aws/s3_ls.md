# aws s3 バケットの一覧を取得するまで

## やること
- 適切なIAMを割り当てる
  - 強めな権限が必要な場合は`AmazonS3FullAccess`
  - RO権限の場合は`AmazonS3ReadOnlyAccess`
  - 権限を絞りたい場合は新たにポリシーを作成する
    - Resourceのbucketを指定する
    - Actionで権限を指定する
```{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "s3:*",
            "Resource": [
                "arn:aws:s3:::hogehoge-bucket",
                "arn:aws:s3:::hogehoge-bucket/*"
            ]
        },
        {
            "Action": [
                "s3:ListAllMyBuckets"
            ],
            "Effect": "Allow",
            "Resource": "arn:aws:s3:::*"
        }
    ]
}
```

## Reference
- [S3のアクションとか](https://docs.aws.amazon.com/ja_jp/IAM/latest/UserGuide/list_amazons3.html#amazons3-actions-as-permissions)
