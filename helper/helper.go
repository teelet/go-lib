package helper

//
// 各种常用小方法，注意不要在这里放置有依赖的方法
//

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"strings"
)

// 单返回值的jsonEncode，修复golang marshal()之后 & 变 u0026 的问题
func JsonEncode(anyType interface{}) string {
	result := ""
	if str, err := json.Marshal(anyType); err == nil {
		result = strings.Replace(string(str), `\u0026`, "&", -1)
	}
	return result
}

func FileGetContents(file string) (string, error) {
	fc, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(fc), nil
}

func GzipEncode(in []byte) ([]byte, error) {
	var (
		buffer bytes.Buffer
		out    []byte
		err    error
	)
	writer := gzip.NewWriter(&buffer)
	_, err = writer.Write(in)
	if err != nil {
		writer.Close()
		return out, err
	}
	err = writer.Close()
	if err != nil {
		return out, err
	}

	return buffer.Bytes(), nil
}

func GzipDecode(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()

	return ioutil.ReadAll(reader)
}
