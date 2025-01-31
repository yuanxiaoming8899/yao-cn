package crypto

import (
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/md4"
)

func init() {
	crypto.RegisterHash(crypto.MD4, md4.New)
}

// HashTypes string
var HashTypes = map[string]crypto.Hash{
	"MD4":         crypto.MD4,
	"MD5":         crypto.MD5,
	"SHA1":        crypto.SHA1,
	"SHA224":      crypto.SHA224,
	"SHA256":      crypto.SHA256,
	"SHA384":      crypto.SHA384,
	"SHA512":      crypto.SHA512,
	"MD5SHA1":     crypto.MD5SHA1,
	"RIPEMD160":   crypto.RIPEMD160,
	"SHA3_224":    crypto.SHA3_224,
	"SHA3_256":    crypto.SHA3_256,
	"SHA3_384":    crypto.SHA3_384,
	"SHA3_512":    crypto.SHA3_512,
	"SHA512_224":  crypto.SHA512_224,
	"SHA512_256":  crypto.SHA512_256,
	"BLAKE2s_256": crypto.BLAKE2s_256,
	"BLAKE2b_256": crypto.BLAKE2b_256,
	"BLAKE2b_384": crypto.BLAKE2b_384,
	"BLAKE2b_512": crypto.BLAKE2b_512,
}

// Hash string
func Hash(hash crypto.Hash, value string) (string, error) {
	h := hash.New()
	_, err := h.Write([]byte(value))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// Hmac the Keyed-Hash Message Authentication Code (HMAC)
func Hmac(hash crypto.Hash, value string, key string, encoding ...string) (string, error) {
	mac := hmac.New(hash.New, []byte(key))
	_, err := mac.Write([]byte(value))
	if err != nil {
		return "", err
	}

	if len(encoding) > 0 && encoding[0] == "base64" {
		return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
	}

	return fmt.Sprintf("%x", mac.Sum(nil)), nil
}

// RSA2Sign RSA2 Sign
func RSA2Sign(prikey string, hash crypto.Hash, value string, encoding ...string) (string, error) {

	privateKey, err := parsePrivateKey(prikey)
	if err != nil {
		return "", err
	}

	h := hash.New()
	_, err = h.Write([]byte(value))
	if err != nil {
		return "", err
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, h.Sum(nil))
	if err != nil {
		return "", err
	}

	if len(encoding) > 0 && encoding[0] == "base64" {
		return base64.StdEncoding.EncodeToString(signature), nil
	}

	return hex.EncodeToString(signature), nil
}

// RSA2Verify RSA2 Verify
func RSA2Verify(pubkey string, hash crypto.Hash, value string, signatureString string, encoding ...string) (bool, error) {

	publicKey, err := parsePublicKey(pubkey)
	if err != nil {
		return false, err
	}

	h := hash.New()
	_, err = h.Write([]byte(value))
	if err != nil {
		return false, err
	}

	var signature []byte
	if len(encoding) > 0 && encoding[0] == "base64" {
		signature, err = base64.StdEncoding.DecodeString(signatureString)
		if err != nil {
			return false, err
		}
	} else {
		signature, err = hex.DecodeString(signatureString)
		if err != nil {
			return false, err
		}
	}

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, h.Sum(nil), []byte(signature))
	return err == nil, nil
}

func parsePrivateKey(privateKeyStr string) (*rsa.PrivateKey, error) {
	privateKeyStr = strings.TrimSpace(privateKeyStr)
	if !strings.HasPrefix(privateKeyStr, "-----BEGIN RSA PRIVATE KEY-----") {
		privateKeyStr = fmt.Sprintf("-----BEGIN RSA PRIVATE KEY-----\n%s\n-----END RSA PRIVATE KEY-----\n", privateKeyStr)
	}

	block, _ := pem.Decode([]byte(privateKeyStr))
	if block == nil {
		return nil, fmt.Errorf("cannot decode PEM block")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch key := key.(type) {
	case *rsa.PrivateKey:
		return key, nil
	default:
		return nil, errors.New("private key error")
	}

}

func parsePublicKey(publicKeyStr string) (*rsa.PublicKey, error) {

	publicKeyStr = strings.TrimSpace(publicKeyStr)
	if !strings.HasPrefix(publicKeyStr, "-----BEGIN RSA PUBLIC KEY-----") {
		publicKeyStr = fmt.Sprintf("-----BEGIN RSA PUBLIC KEY-----\n%s\n-----END RSA PUBLIC KEY-----\n", publicKeyStr)
	}

	block, _ := pem.Decode([]byte(publicKeyStr))
	if block == nil {
		return nil, fmt.Errorf("cannot decode PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("public key error")
	}
}
