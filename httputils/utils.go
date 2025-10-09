package httputils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"errors"
)

func SendJSON(w http.ResponseWriter, status int, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(obj)
	if err != nil {
		return errors.Join(err, fmt.Errorf("Error encoding json response: %v", obj))
	}
	w.WriteHeader(status)
	_, err = w.Write(b)
	return err
}

func DecodeJSON(r *http.Request, dst interface{}) error {
	// ボディサイズ制限（1MB）
	r.Body = http.MaxBytesReader(nil, r.Body, 1<<10)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // ✅ デフォルトで有効

	err := decoder.Decode(dst)
	if err != nil {
		return err
	}

	// 余分なデータがないかチェック
	if decoder.More() {
		return errors.New("body must only contain a single JSON object")
	}

	return nil
}
