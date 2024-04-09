package crypto

import "crypto/sha256"

func SHA256(data []byte) ([]byte, error) {
	h := sha256.New()
	_, err := h.Write(data)
	if err != nil {
		return nil, err
	}
	r := h.Sum(nil)
	return r, nil
}
