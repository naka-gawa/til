## Summary

Go 言語で作られてる sops を CLI 上で実行するのではなく

パッケージを直接扱う場合のサンプルコードを調べた

## Descrypt

```go
package main

import (
    "fmt"
    "os"

    "github.com/getsops/sops/v3/decrypt" // sopsの復号化機能を提供するパッケージ
)

func main() {
    encryptedFilePath := "test.yaml"

    decryptedData, err := decrypt.File(encryptedFilePath, "yaml")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error decrypting file: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Decrypted Data:", string(decryptedData))
}
```

## Encrypt

```go
package main

import (
	"fmt"
	"os"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/aes" // AES暗号化ライブラリ
	"github.com/getsops/sops/v3/cmd/sops/common" // SOPSの共通機能
	"github.com/getsops/sops/v3/keyservice" // キーサービスクライアント
	"github.com/getsops/sops/v3/kms" // AWS KMSに関する機能
	"github.com/getsops/sops/v3/stores/yaml" // YAMLファイルの読み書き
)

func main() {
	inputPath := "test-non-enc.yaml"
	outputPath := "test-enc.yaml"

	// AWS KMSのARN（リソース識別子）
	masterKeyArn := "arn:aws:kms:ap-northeast-1:xxxxxxxxxxxx:key/hogehoge"

	// ファイルの暗号化を実行
	err := EncryptFile(inputPath, outputPath, masterKeyArn)
	if err != nil {
		// エラーが発生した場合、標準エラー出力にメッセージを出力して終了
		fmt.Fprintf(os.Stderr, "Error encrypting file: %v\n", err)
		os.Exit(1)
	}

	// 暗号化成功メッセージ
	fmt.Println("File encrypted successfully")
}

// EncryptFile はファイルを暗号化する関数
func EncryptFile(inputPath, outputPath, masterKeyArn string) error {
	// 入力ファイルを読み込み
	fileBytes, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("error reading file: %s", err)
	}

	// YAMLファイル用のストアインスタンスを作成
	store := &yaml.Store{}
	// YAMLファイルからプレーンテキストデータを読み込む
	branches, err := store.LoadPlainFile(fileBytes)
	if err != nil {
		return fmt.Errorf("error loading plain file: %s", err)
	}

	// AWS KMSを使ってマスターキーを作成
	kmsMasterKey := kms.NewMasterKey(masterKeyArn, "", make(map[string]*string))

	// SOPSデータ構造を構築
	tree := sops.Tree{
		Branches: branches,
		Metadata: sops.Metadata{
			KeyGroups: []sops.KeyGroup{
				{
					kmsMasterKey, // マスターキーをキーグループに追加
				},
			},
			Version: "3.7.1", // SOPSのバージョン
		},
	}

	// データキーの生成とマスターキーでの暗号化
	dataKey, errs := tree.GenerateDataKeyWithKeyServices([]keyservice.KeyServiceClient{keyservice.NewLocalClient()})
	if len(errs) > 0 {
		return fmt.Errorf("error generating data key: %v", errs[0])
	}

	// ツリーの暗号化
	err = common.EncryptTree(common.EncryptTreeOpts{
		DataKey: dataKey,
		Tree:    &tree,
		Cipher:  aes.NewCipher(), // AES暗号化を使用
	})
	if err != nil {
		return fmt.Errorf("error encrypting tree: %s", err)
	}

	// 暗号化されたデータをファイルに出力
	encryptedBytes, err := store.EmitEncryptedFile(tree)
	if err != nil {
		return fmt.Errorf("could not emit encrypted file: %s", err)
	}

	// 出力ファイルへの書き込み
	err = os.WriteFile(outputPath, encryptedBytes, 0644)
	if err != nil {
		return fmt.Errorf("error writing encrypted file: %s", err)
	}

	return nil
}
```

上記のケースは KMS を利用して暗号化するケース
ただし、Encrypt は unstable なので利用は推奨されない点注意が必要

