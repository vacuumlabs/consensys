package hash

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"vax/pkg/model/dto"
)

func EventPayload(payload dto.EventPayload) (string, error) {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return "", nil
	}

	h := sha256.New()
	h.Write(payloadJson)
	return base64.URLEncoding.EncodeToString(h.Sum(nil)), nil
}
