package testdata

import (
	"crypto/cipher"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
)

// --- DeprecatedPKCS1v15 ---

func checkPKCS1v15(pub *rsa.PublicKey, priv *rsa.PrivateKey) {
	msg := []byte("secret")
	key := make([]byte, 32)

	// Should trigger: deprecated PKCS#1 v1.5 encryption
	_, _ = rsa.EncryptPKCS1v15(rand.Reader, pub, msg)              // want: "rsa.EncryptPKCS1v15 is deprecated"
	_, _ = rsa.DecryptPKCS1v15(rand.Reader, priv, msg)             // want: "rsa.DecryptPKCS1v15 is deprecated"
	_ = rsa.DecryptPKCS1v15SessionKey(rand.Reader, priv, msg, key) // want: "rsa.DecryptPKCS1v15SessionKey is deprecated"
}

// --- WeakRSAKeySize ---

func checkWeakRSA() {
	// Should trigger: 1024-bit key
	_, _ = rsa.GenerateKey(rand.Reader, 1024) // want: "RSA 1024-bit keys are considered weak"

	// Should trigger: 512-bit key
	_, _ = rsa.GenerateKey(rand.Reader, 512) // want: "RSA keys smaller than 1024"

	// Should NOT trigger: adequate key size
	_, _ = rsa.GenerateKey(rand.Reader, 2048)
}

// --- DeprecatedCipherModes ---

func checkCipherModes(block cipher.Block, iv []byte) {
	// Should trigger: deprecated OFB
	_ = cipher.NewOFB(block, iv) // want: "cipher.NewOFB is deprecated"

	// Should trigger: deprecated CFB
	_ = cipher.NewCFBEncrypter(block, iv) // want: "cipher.NewCFBEncrypter is deprecated"
	_ = cipher.NewCFBDecrypter(block, iv) // want: "cipher.NewCFBDecrypter is deprecated"

	// Should NOT trigger: CTR (the recommended replacement)
	_ = cipher.NewCTR(block, iv)
}

// --- DeprecatedElliptic ---

func checkElliptic() {
	// Should trigger: deprecated elliptic.GenerateKey
	_, _, _, _ = elliptic.GenerateKey(elliptic.P256(), rand.Reader) // want: "elliptic.GenerateKey is deprecated"

	// Should trigger: deprecated elliptic.Marshal
	x, y := elliptic.P256().Params().Gx, elliptic.P256().Params().Gy
	_ = elliptic.Marshal(elliptic.P256(), x, y) // want: "elliptic.Marshal is deprecated"

	// Should trigger: deprecated elliptic.Unmarshal
	data := []byte{0x04}
	_, _ = elliptic.Unmarshal(elliptic.P256(), data) // want: "elliptic.Unmarshal is deprecated"

	// Should NOT trigger: getting a curve (not deprecated by our rules)
	_ = elliptic.P256()
}

// --- DeprecatedRSAMultiPrime ---

func checkMultiPrime() {
	// Should trigger: deprecated GenerateMultiPrimeKey
	_, _ = rsa.GenerateMultiPrimeKey(rand.Reader, 3, 2048) // want: "rsa.GenerateMultiPrimeKey is deprecated"
}
