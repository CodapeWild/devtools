package pemfile

import (
	"crypto/rsa"
	"devtools/certificate/secure"
	"log"
	"os"
	"testing"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func TestCertificate(t *testing.T) {
	issuerTemp, err := SetCertTemp("./config.json")
	if err != nil {
		log.Fatalln(err.Error())
	}
	caCrt, _, caPriv, err := GenCertificatePem("../key/ca/ca.crt", issuerTemp, nil, nil)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, _, _, err = GenCertificatePem("../key/issuer/my.crt", issuerTemp, caCrt, caPriv)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func TestPemKeyParse(t *testing.T) {
	key, kt, err := ParsePemKeyFile("../key/issuer/priv.key")
	if err != nil {
		log.Fatalln(err.Error())
	}
	if kt == rsa_private_key {
		rsaPriv := key.(*rsa.PrivateKey)
		cipher, err := secure.Encrypt(&rsaPriv.PublicKey, []byte("hello, tnt"))
		if err != nil {
			log.Fatalln(err.Error())
		}
		plain, err := secure.Decrypt(rsaPriv, cipher)
		if err != nil {
			log.Fatalln(err.Error())
		}
		log.Println(string(plain))
	}
}

func TestCrtParse(t *testing.T) {
	crts, err := ParsePemCrtFile("../key/ca.crt")
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(len(crts))
	for _, v := range crts {
		log.Println(v.PublicKeyAlgorithm)
	}
}
