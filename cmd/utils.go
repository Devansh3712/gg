package cmd

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

func ReadObject(hash string) ([]byte, error) {
	dir, file := hash[:2], hash[2:]
	path := fmt.Sprintf(".git/objects/%s/%s", dir, file)
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func WriteObject(hash string, content []byte) error {
	dir, file := hash[:2], hash[2:]
	dirPath := fmt.Sprintf(".git/objects/%s", dir)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.Mkdir(dirPath, 0755); err != nil {
			return err
		}
	}
	filePath := fmt.Sprintf("%s/%s", dirPath, file)
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return err
	}
	return nil
}

func HashObject(objectType string, content []byte) []byte {
	data := fmt.Sprintf("%s %d\x00%s", objectType, len(content), content)
	hash := sha1.New()
	hash.Write([]byte(data))
	return hash.Sum(nil)
}

func ZLibCompress(content []byte) ([]byte, error) {
	compressed := new(bytes.Buffer)
	writer := zlib.NewWriter(compressed)
	if _, err := writer.Write(content); err != nil {
		return nil, err
	}
	writer.Close()
	return compressed.Bytes(), nil
}

func ZLibDecompress(content []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(content)
	reader, err := zlib.NewReader(buffer)
	if err != nil {
		return nil, err
	}
	reader.Close()

	decompressed, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return decompressed, nil
}
