package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"math/big"

	"github.com/adinandradrs/omni-service-sdk/pkg/domain"
	"github.com/sethvargo/go-password/password"
	"go.uber.org/zap"
)

func RandomOtp(d int) (string, error) {
	ns := "012345679"
	b := make([]byte, d)
	for i := 0; i < d; i++ {
		max := big.NewInt(int64(len(ns)))
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = ns[num.Int64()]
	}
	return string(b), nil
}

func RandomPassword(len int, d int, sym int, logger *zap.Logger) (res string, e *domain.TechnicalError) {
	res, err := password.Generate(len, d, sym, false, false)
	if err != nil {
		return res, Exception("failed to generate password", err, logger)
	}
	logger.Info("success generate password", zap.String("generated", res))
	return res, e
}

func Hash(key string) string {
	md5 := md5.New()
	md5.Write([]byte(key))
	return hex.EncodeToString(md5.Sum(nil))
}

func Encrypt(logger *zap.Logger, d string, h string) (res []byte, ex *domain.TechnicalError) {
	c, _ := aes.NewCipher([]byte(h))
	o, err := cipher.NewGCM(c)
	if err != nil {
		return res, Exception("failed to encrypt GCM data", err, logger)
	}
	nsz := make([]byte, o.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nsz); err != nil {
		return res, Exception("failed to encrypt I/O reader", err, logger)
	}
	return o.Seal(nsz, nsz, []byte(d), nil), ex
}

func Decrypt(data []byte, hash string, logger *zap.Logger) (res string, ex *domain.TechnicalError) {
	c, err := aes.NewCipher([]byte(hash))
	if err != nil {
		return res, Exception("failed to decrypt chiper", err, logger)
	}
	aead, err := cipher.NewGCM(c)
	if err != nil {
		return res, Exception("failed to decrypt init GCM", err, logger)
	}
	nsz := aead.NonceSize()
	n, cbytes := data[:nsz], data[nsz:]
	o, err := aead.Open(nil, n, cbytes, nil)
	if err != nil {
		return res, Exception("failed to decrypt nonce size parse logic", err, logger)
	}
	return string(o), ex
}
