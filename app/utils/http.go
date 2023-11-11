package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

func ReadAndReplaceBody(r *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	return body, nil
}
