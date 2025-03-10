package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	//json.Marshal() でGoの構造体をJSONに変換し、json.Unmarshal() でJSONをGoの構造体に変換するらしい
	"os"
	//ファイルに関する操作ができる
	// STEP 5-1: uncomment this line
	// _ "github.com/mattn/go-sqlite3"
)

var (
	errImageNotFound = errors.New("image not found")
	errItemsNotFound = errors.New("item not found")
)

type Item struct {
	ID   int    `db:"id" json:"-"`
	Name string `db:"name" json:"name"`
	Category string `db:"category" json:"category"`
	ImageName string `json:"image_name"`
}

// Please run `go generate ./...` to generate the mock implementation
// ItemRepository is an interface to manage items.
//
//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -package=${GOPACKAGE} -destination=./mock_$GOFILE
type ItemRepository interface {
	Insert(ctx context.Context, item *Item) error
	List(ctx context.Context) ([]*Item, error)
	Select(ctx context.Context, id int) (*Item, error)
}

// itemRepository is an implementation of ItemRepository
type itemRepository struct {
	// fileName is the path to the JSON file storing items.
	fileName string
}

// NewItemRepository creates a new itemRepository.
func NewItemRepository() ItemRepository {
	return &itemRepository{fileName: "items.json"}
}

// Insert inserts an item into the repository.
func (i *itemRepository) Insert(ctx context.Context, item *Item) error {
	// STEP 4-1: add an implementation to store an item
	const fileName = "items.json"
	var data struct{
		Items [] *Item `json:"items"`
	}

	// 既存のファイルを開く (なければ作る)
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close() // 終わったら閉じる

	// データを読み込む
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		// JSONファイルが空の場合は無視
		if err.Error() != "EOF" {
			return err
		}
	}

	// 新しいアイテムを追加
	data.Items = append(data.Items, item)

	// ファイルを開き直して書き込み (上書き)
	file, err = os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(&data); err != nil {
		return err
	}

	return nil
}

func (i *itemRepository) List(ctx context.Context) ([]*Item, error) {
	var data struct {
		Items []*Item `json:"items"`
	}
	dataBytes, err := os.ReadFile(i.fileName)
	if err != nil {
		if(os.IsNotExist(err)) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read files", err)
	}
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return nil, fmt.Errorf("failed to Unmarshal JSON", err)
	}
	return data.Items, nil
}


func (i *itemRepository) Select(ctx context.Context, id int) (*Item, error) {
	if id <= 0 {
		return nil, errItemsNotFound
	} //id-1番目なので

	items, err := i.List(ctx)
	if err != nil {
		return nil, err
	}

	if len(items) < id {
		return nil, errItemsNotFound
	}

	return items[id-1], nil
}


// StoreImage stores an image and returns an error if any.
// This package doesn't have a related interface for simplicity.
func StoreImage(fileName string, image []byte) error {
	// STEP 4-4: add an implementation to store an image
	if err := os.WriteFile(fileName, image, 0644); err != nil {
		return fmt.Errorf("failed to write image file: %s", err)
	}
	return nil
}
