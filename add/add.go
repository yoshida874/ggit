package add

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
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
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return "", err
    }

    // オブジェクトファイルを保存するディレクトリが存在しない場合は作成する
    if err := os.MkdirAll(filepath.Dir(objPath), 0755); err != nil {
        return "", err
    }

    // オブジェクトファイルを保存する
    if err := ioutil.WriteFile(objPath, data, 0644); err != nil {
        return "", err
    }

    return hash, nil
}

func Add(file string) error {
	// ファイルをオブジェクトファイルに保存する
	hash, err := saveFileToObjects(file)
	if err != nil {
		return err
	}

	// インデックスファイルにファイル名とハッシュ値を保存する
	indexFile, err := os.OpenFile(".ggit/index", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer indexFile.Close()

	if _, err := fmt.Fprintf(indexFile, "100644 %s 0\t%s\n", hash, file); err != nil {
		return err
	}

    return nil
}
