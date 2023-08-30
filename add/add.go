package add

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func hashFile(path string) (string, error) {
    f, err := os.Open(path)
    if err != nil {
        return "", err
    }
    defer f.Close()

    h := sha1.New()
    if _, err := io.Copy(h, f); err != nil {
        return "", err
    }

    return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func saveFileToObjects(path string) (string, error) {
    // ファイルのハッシュ値を計算する
    hash, err := hashFile(path)
    if err != nil {
        return "", err
    }

    // ハッシュ値からオブジェクトファイル名を生成する
    objPath := fmt.Sprintf(".ggit/objects/%s/%s", hash[:2], hash[2:])

    // オブジェクトファイルが既に存在する場合は何もしない
    _, err = os.Stat(objPath)
    if err == nil {
        return hash, nil
    }

    // ファイルの内容をバイナリ形式でオブジェクトファイルに保存する
    data, err := os.ReadFile(path)
    if err != nil {
        return "", err
    }

    // オブジェクトファイルを保存するディレクトリが存在しない場合は作成する
    if err := os.MkdirAll(filepath.Dir(objPath), 0755); err != nil {
        return "", err
    }

    // blobとしてデータを圧縮
    compressedData, err := compressObject(data, "blob")
    if err != nil {
        return "", err
    }

    // オブジェクトファイルを保存する
    if err := os.WriteFile(objPath, compressedData, 0644); err != nil {
        return "", err
    }

    return hash, nil
}

// バイナリフォーマットでindexファイルに保存    
func addToIndex(hash string, file string) error {
	indexFile, err := os.OpenFile(".ggit/index", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer indexFile.Close()

	hashBytes := []byte(hash)
	fileBytes := []byte(file)
	
	err = binary.Write(indexFile, binary.LittleEndian, int64(len(hashBytes)))
	if err != nil {
		return err
	}

	err = binary.Write(indexFile, binary.LittleEndian, hashBytes)
	if err != nil {
		return err
	}

	err = binary.Write(indexFile, binary.LittleEndian, int64(len(fileBytes)))
	if err != nil {
		return err
	}

	err = binary.Write(indexFile, binary.LittleEndian, fileBytes)
	if err != nil {
		return err
	}

	return nil
}

// indexに再帰的にファイルを処理
func Add(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		files, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		for _, file := range files {
			err = Add(filepath.Join(path, file.Name()))
			if err != nil {
				return err
			}
		}
		return nil
	} else {
		hash, err := saveFileToObjects(path)
		if err != nil {
			return err
		}
		return addToIndex(hash, path)
	}
}

func compressObject(data []byte, objectType string) ([]byte, error) {
	var compressedData bytes.Buffer
	writer := zlib.NewWriter(&compressedData)
	header := fmt.Sprintf("%s %d\x00", objectType, len(data))
	writer.Write([]byte(header))
	writer.Write(data)
	writer.Close()
	return compressedData.Bytes(), nil
}
