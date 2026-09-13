package jiosaavn

import (
	"crypto/des"
	"encoding/base64"
	"strings"
)

func DecryptMediaURL(encryptedURL string) string {
	key := []byte("38346591")
	block, err := des.NewCipher(key)
	if err != nil {
		return ""
	}

	decoded, err := base64.StdEncoding.DecodeString(encryptedURL)
	if err != nil {
		return ""
	}

	decrypted := make([]byte, len(decoded))
	
	// DES ECB mode decryption
	bs := block.BlockSize()
	for i := 0; i < len(decoded); i += bs {
		block.Decrypt(decrypted[i:i+bs], decoded[i:i+bs])
	}
	
	if len(decrypted) == 0 {
		return ""
	}
	
	// PKCS5 unpadding
	paddingLen := int(decrypted[len(decrypted)-1])
	if paddingLen > len(decrypted) || paddingLen == 0 {
		paddingLen = 0
	}
	url := string(decrypted[:len(decrypted)-paddingLen])
	
	url = strings.Replace(url, "_96.mp4", "_320.mp4", 1)
	url = strings.Replace(url, "_96_p.mp4", "_320.mp4", 1)
	return url
}
