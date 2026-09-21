package store

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"tickets-system-v1/internal/models"
)

var (
	Data       models.NumberSystem
	ConfigData models.Config
)

const storeFilePath = "data/store.json"
const configFilePath = "config.json"

func InitStore() {
	Data = models.NumberSystem{
		CurrentNumbers: []int{},
		WaitingNumbers: []int{},
	}

	if err := os.MkdirAll("data", 0755); err != nil {
		log.Printf("警告: dataディレクトリの作成に失敗しました: %v", err)
	}

	loadConfig()
	loadData()
}

func loadConfig() {
	file, err := os.Open(configFilePath)
	if err != nil {
		log.Printf("警告: config.jsonが見つかりません。デフォルトのカテゴリを使用します。")
		return
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("エラー: config.jsonの読み込みに失敗しました: %v", err)
		return
	}

	if err := json.Unmarshal(bytes, &ConfigData); err != nil {
		log.Printf("エラー: config.jsonの解析に失敗しました: %v", err)
	}
}

func loadData() {
	file, err := os.Open(storeFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("エラー: storeファイルの展開に失敗しました: %v", err)
		}
		return
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("エラー: storeファイルの読み込みに失敗しました: %v", err)
		return
	}

	if len(bytes) == 0 {
		return
	}

	if err := json.Unmarshal(bytes, &Data); err != nil {
		log.Printf("エラー: storeデータの解析に失敗しました: %v", err)
	} else {
		log.Printf("%s からデータを正常に読み込みました", storeFilePath)
	}
}

func SaveData() {
	Data.Mutex.RLock()
	defer Data.Mutex.RUnlock()

	bytes, err := json.MarshalIndent(Data, "", "  ")
	if err != nil {
		log.Printf("エラー: データのJSON変換に失敗しました: %v", err)
		return
	}

	tempFile := storeFilePath + ".tmp"
	if err := os.WriteFile(tempFile, bytes, 0644); err != nil {
		log.Printf("エラー: 一時ファイルへの書き込みに失敗しました: %v", err)
		return
	}

	if err := os.Rename(tempFile, storeFilePath); err != nil {
		if err := os.WriteFile(storeFilePath, bytes, 0644); err != nil {
			log.Printf("エラー: storeファイルへの書き込みに失敗しました: %v", err)
		} else {
			os.Remove(tempFile)
		}
	}
}
