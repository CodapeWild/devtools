package pemfile

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

type KeyType int

const (
	nonsupport_key KeyType = iota
	rsa_private_key
	rsa_public_key
	ecdsa_private_key
	ecdsa_public_key
)

var (
	NoPemDataFound     = errors.New("no pem data found")
	UnsupportedKeyType = errors.New("unsupported key type, only RSA or ECDSA pem key file")
	UnsupportedKeyFile = errors.New("unsupported key file type, only PKCS#1 or PKCS#8 or EC PRIVATE KEY format")
	NoCrtDataFound     = errors.New("no certificate pem data found")
)

func ParsePemKeyFile(filePath string) (key interface{}, kt KeyType, err error) {
	kfBuf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, nonsupport_key, err
	}

	bc, _ := pem.Decode(kfBuf)
	if bc == nil {
		return nil, nonsupport_key, NoPemDataFound
	}

	var ok bool
	switch KeyFileType(bc.Type) {
	case kf_pkcs1_private:
		if key, err = x509.ParsePKCS1PrivateKey(bc.Bytes); err == nil {
			kt = rsa_private_key
		}
	case kf_pkcs8_private:
		if key, err = x509.ParsePKCS8PrivateKey(bc.Bytes); err == nil {
			if _, ok = key.(*rsa.PrivateKey); ok {
				kt = rsa_private_key
			} else if _, ok = key.(*ecdsa.PrivateKey); ok {
				kt = ecdsa_private_key
			} else {
				err = UnsupportedKeyType
			}
		}
	case kf_ec_private:
		if key, err = x509.ParseECPrivateKey(bc.Bytes); err == nil {
			kt = ecdsa_private_key
		}
	case kf_pkcs8_public:
		if key, err = x509.ParsePKIXPublicKey(bc.Bytes); err == nil {
			if _, ok = key.(*rsa.PublicKey); ok {
				kt = rsa_public_key
			} else if _, ok = key.(*ecdsa.PublicKey); ok {
				kt = ecdsa_public_key
			} else {
				err = UnsupportedKeyType
			}
		}
	case kf_pkcs1_public:
		if key, err = x509.ParsePKCS1PublicKey(bc.Bytes); err == nil {
			kt = rsa_public_key
		}
	default:
		err = UnsupportedKeyFile
	}

	return
}

func ParsePemCrtFile(filePath string) ([]*x509.Certificate, error) {
	crtBuf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var (
		crts []*x509.Certificate
		bc   *pem.Block
		rest []byte
	)
NEXT_CRT:
	bc, rest = pem.Decode(crtBuf)
	if bc != nil {
		if bc.Type == string(kf_certificate) {
			if crt, err := x509.ParseCertificate(bc.Bytes); err != nil {
				return nil, err
			} else {
				crts = append(crts, crt)
			}
		}
		if len(rest) != 0 {
			crtBuf = rest
			goto NEXT_CRT
		}
	}

	if len(crts) == 0 {
		return nil, NoCrtDataFound
	} else {
		return crts, nil
	}
}
