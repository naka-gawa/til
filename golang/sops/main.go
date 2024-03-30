package main

import (
	"fmt"
	"os"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/aes"
	"github.com/getsops/sops/v3/cmd/sops/common"
	"github.com/getsops/sops/v3/keyservice"
	"github.com/getsops/sops/v3/kms"
	"github.com/getsops/sops/v3/stores/yaml"
)

func main() {
	inputPath := "test-non-enc.yaml"
	outputPath := "test-enc.yaml"

	// AWS KMSのARN
	masterKeyArn := "arn:aws:kms:ap-northeast-1:063621714126:key/9201601c-bf4e-4144-8d0c-3519c3ff541f"

	err := EncryptFile(inputPath, outputPath, masterKeyArn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encrypting file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("File encrypted successfully")
}

func EncryptFile(inputPath, outputPath, masterKeyArn string) error {
	fileBytes, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("error reading file: %s", err)
	}

	store := &yaml.Store{} // YAML ストアのインスタンスを作成
	branches, err := store.LoadPlainFile(fileBytes)
	if err != nil {
		return fmt.Errorf("error loading plain file: %s", err)
	}

	// AWS KMSマスターキーを作成
	kmsMasterKey := kms.NewMasterKey(masterKeyArn, "", make(map[string]*string))

	tree := sops.Tree{
		Branches: branches,
		Metadata: sops.Metadata{
			KeyGroups: []sops.KeyGroup{
				{
					kmsMasterKey,
				},
			},
			Version: "3.7.1", // SOPSのバージョン情報を設定
		},
	}

	dataKey, errs := tree.GenerateDataKeyWithKeyServices([]keyservice.KeyServiceClient{keyservice.NewLocalClient()})
	if len(errs) > 0 {
		return fmt.Errorf("error generating data key: %v", errs[0])
	}

	err = common.EncryptTree(common.EncryptTreeOpts{
		DataKey: dataKey,
		Tree:    &tree,
		Cipher:  aes.NewCipher(),
	})
	if err != nil {
		return fmt.Errorf("error encrypting tree: %s", err)
	}

	encryptedBytes, err := store.EmitEncryptedFile(tree)
	if err != nil {
		return fmt.Errorf("could not emit encrypted file: %s", err)
	}

	err = os.WriteFile(outputPath, encryptedBytes, 0644)
	if err != nil {
		return fmt.Errorf("error writing encrypted file: %s", err)
	}

	return nil
}
