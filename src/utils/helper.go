package utils

import (
	"encoding/base64"
)

func DecodeBase64(base64_encode_data map[string]string) (map[string]string, error) {
	var decoded_map map[string]string = make(map[string]string)

	for k, v := range base64_encode_data {
		data, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return nil, err
		}
		decoded_map[k] = string(data)
	}
	return decoded_map, nil
}
