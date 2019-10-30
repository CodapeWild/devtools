package secure

import (
	"bytes"
	"crypto/dsa"
	"crypto/elliptic"
	"log"
	"os"
	"testing"
)

var data = []byte("hello,tnt")

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func TestRSA(t *testing.T) {
	privKey, err := GenRSAPrivateKey(2048)
	if err != nil {
		log.Fatalln(err.Error())
	}

	cipher, err := Encrypt(&privKey.PublicKey, data)
	if err != nil {
		log.Fatalln(err.Error())
	}

	plain, err := Decrypt(privKey, cipher)
	if err != nil {
		log.Fatalln(err.Error())
	} else if bytes.Compare(data, plain) != 0 {
		log.Fatalln("encrypt and decrypt using RSA failed")
	} else {
		log.Println("encrypt and decrypt using RSA success")
	}

	sig, err := Sign(privKey, MD5Digester(data))
	if err != nil {
		log.Fatalln(err.Error())
	}

	if err = Verify(&privKey.PublicKey, MD5Digester(data), sig); err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("sign and verify using RSA success")
}

func TestECDSA(t *testing.T) {
	privKey, err := GenECDSAPrivateKey(elliptic.P256())
	if err != nil {
		log.Fatalln(err.Error())
	}

	r, s, err := DigitalSign(privKey, ECDSA, data)
	if err != nil {
		log.Fatalln(err.Error())
	}

	if DigitalVerify(&privKey.PublicKey, ECDSA, data, r, s) {
		log.Println("sign and verify using ECDSA success")
	} else {
		log.Fatalln("signature verify failed")
	}
}

func TestDSA(t *testing.T) {
	privKey, err := GenDSAPrivateKey(dsa.L2048N256)
	if err != nil {
		log.Fatalln(err.Error())
	}

	r, s, err := DigitalSign(privKey, DSA, data)
	if err != nil {
		log.Fatalln(err.Error())
	}

	if DigitalVerify(&privKey.PublicKey, DSA, data, r, s) {
		log.Println("sign and verify using DSA success")
	} else {
		log.Fatalln("signature verify failed")
	}
}
