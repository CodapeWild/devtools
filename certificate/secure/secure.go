package secure

import (
	"crypto"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"math/big"
)

type CryptoAlgorithm int

const (
	RSA CryptoAlgorithm = iota + 1
	ECDSA
	DSA
)

var (
	UnknowCryptoAlgorithm = errors.New("unknow cryptography algorithm")
)

func GenRSAPrivateKey(bits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, bits)
}

func GenECDSAPrivateKey(c elliptic.Curve) (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(c, rand.Reader)
}

func GenDSAPrivateKey(pSize dsa.ParameterSizes) (*dsa.PrivateKey, error) {
	param := &dsa.Parameters{}
	err := dsa.GenerateParameters(param, rand.Reader, pSize)
	if err != nil {
		return nil, err
	}

	privKey := &dsa.PrivateKey{
		PublicKey: dsa.PublicKey{
			Parameters: *param,
		},
	}
	if err = dsa.GenerateKey(privKey, rand.Reader); err != nil {
		return nil, err
	} else {
		return privKey, nil
	}
}

func Encrypt(pubKey *rsa.PublicKey, plain []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, pubKey, plain)
}

func Decrypt(privKey *rsa.PrivateKey, cipher []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privKey, cipher)
}

func MD5Digester(data []byte) func() (hash crypto.Hash, hashed []byte) {
	return func() (hash crypto.Hash, hashed []byte) {
		sbuf := md5.Sum(data)

		return crypto.MD5, sbuf[:]
	}
}

func SHA1Digester(data []byte) func() (hash crypto.Hash, hashed []byte) {
	return func() (hash crypto.Hash, hashed []byte) {
		sbuf := sha1.Sum(data)

		return crypto.SHA1, sbuf[:]
	}
}

func SHA256Digester(data []byte) func() (hash crypto.Hash, hashed []byte) {
	return func() (hash crypto.Hash, hashed []byte) {
		sbuf := sha256.Sum256(data)

		return crypto.SHA256, sbuf[:]
	}
}

func SHA512Digester(data []byte) func() (hash crypto.Hash, hashed []byte) {
	return func() (hash crypto.Hash, hashed []byte) {
		sbuf := sha512.Sum512(data)

		return crypto.SHA512, sbuf[:]
	}
}

// sign with RSA
func Sign(privKey *rsa.PrivateKey, digester func() (hash crypto.Hash, hashed []byte)) ([]byte, error) {
	hash, hashed := digester()

	return rsa.SignPKCS1v15(rand.Reader, privKey, hash, hashed)
}

// verify RSA signature
func Verify(pubKey *rsa.PublicKey, digester func() (hash crypto.Hash, hashed []byte), sig []byte) error {
	hash, hashed := digester()

	return rsa.VerifyPKCS1v15(pubKey, hash, hashed, sig)
}

// sign with ECDSA or DSA
func DigitalSign(privKey crypto.PrivateKey, algo CryptoAlgorithm, hash []byte) (r, s *big.Int, err error) {
	switch algo {
	case ECDSA:
		return ecdsa.Sign(rand.Reader, privKey.(*ecdsa.PrivateKey), hash)
	case DSA:
		return dsa.Sign(rand.Reader, privKey.(*dsa.PrivateKey), hash)
	default:
		return nil, nil, UnknowCryptoAlgorithm
	}
}

// verify ECDSA or DSA signature
func DigitalVerify(pubKey crypto.PublicKey, algo CryptoAlgorithm, hash []byte, r, s *big.Int) bool {
	switch algo {
	case ECDSA:
		return ecdsa.Verify(pubKey.(*ecdsa.PublicKey), hash, r, s)
	case DSA:
		return dsa.Verify(pubKey.(*dsa.PublicKey), hash, r, s)
	default:
		return false
	}
}
