package pemfile

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"devtools/certificate/secure"
	"devtools/comerr"
	"devtools/file"
	"encoding/pem"
	"errors"
	"io"
	"math/big"
	"net"
	"net/url"
	"os"
	"path"
	"time"
)

type KeyFileType string

const (
	kf_pkcs1_private     KeyFileType = "RSA PRIVATE KEY"
	kf_pkcs1_public      KeyFileType = "RSA PUBLIC KEY"
	kf_pkcs8_private     KeyFileType = "PRIVATE KEY"
	kf_pkcs8_public      KeyFileType = "PUBLIC KEY"
	kf_ec_private        KeyFileType = "EC PRIVATE KEY"
	kf_encrypted_private KeyFileType = "ENCRYPTED PRIVATE KEY"
	kf_certificate       KeyFileType = "CERTIFICATE"
)

var (
	KeyTypeConvertFailed = errors.New("cryptography key type conversion failed")
	UnknowKeyFileType    = errors.New("unknow cryptography key file type")
)

func MarshalKey(key interface{}, kf KeyFileType) ([]byte, error) {
	var (
		keyBuf []byte
		err    error
	)
	switch kf {
	case kf_pkcs1_private:
		if rsaPrivKey, ok := key.(*rsa.PrivateKey); !ok {
			err = KeyTypeConvertFailed
		} else {
			keyBuf = x509.MarshalPKCS1PrivateKey(rsaPrivKey)
		}
	case kf_pkcs8_private:
		keyBuf, err = x509.MarshalPKCS8PrivateKey(key)
	case kf_ec_private:
		if ecdsaPrivKey, ok := key.(*ecdsa.PrivateKey); !ok {
			err = KeyTypeConvertFailed
		} else {
			keyBuf, err = x509.MarshalECPrivateKey(ecdsaPrivKey)
		}
	case kf_pkcs8_public:
		keyBuf, err = x509.MarshalPKIXPublicKey(key)
	case kf_pkcs1_public:
		if rsaPubKey, ok := key.(*rsa.PublicKey); !ok {
			err = KeyTypeConvertFailed
		} else {
			keyBuf = x509.MarshalPKCS1PublicKey(rsaPubKey)
		}
	default:
		err = UnknowKeyFileType
	}

	return keyBuf, err
}

func KeyBlock(key interface{}, kf KeyFileType) (*pem.Block, error) {
	if buf, err := MarshalKey(key, kf); err != nil {
		return nil, err
	} else {
		return &pem.Block{
			Type:  string(kf),
			Bytes: buf,
		}, nil
	}
}

func EncKeyBlock(key interface{}, kf KeyFileType, pswd string, alg x509.PEMCipher) (*pem.Block, error) {
	if buf, err := MarshalKey(key, kf); err != nil {
		return nil, err
	} else {
		return x509.EncryptPEMBlock(rand.Reader, string(kf_encrypted_private), buf, []byte(pswd), alg)
	}
}

func WritePemTo(out io.Writer, bc *pem.Block) error {
	if out == nil || bc == nil {
		return comerr.ParamInvalid
	}

	return pem.Encode(out, bc)
}

type CertTempConfig struct {
	Country           []string `json:"country"`
	City              []string `json:"city"`
	Province          []string `json:"province"`
	Section           []string `json:"section"`
	Company           []string `json:"company"`
	Street            []string `json:"street"`
	PostCode          []string `json:"post_code"`
	CommonName        string   `json:"common_name"`
	AlternateDNSNames []string `json:"alternate_dns_names"`
	AlternateEmails   []string `json:"alternate_emails"`
	AlternateIPs      []string `json:"alternate_ips"`
	AlternateURLs     []string `json:"alternate_urls"`
	NotBefore         string   `json:"not_before"`
	NotAfter          string   `json:"not_after"`
}

func SetCertTemp(confPath string) (*x509.Certificate, error) {
	certConf := &CertTempConfig{}
	err := file.ReadJsonFile(confPath, certConf)
	if err != nil {
		return nil, err
	}

	var ips []net.IP
	for _, v := range certConf.AlternateIPs {
		ips = append(ips, net.ParseIP(v))
	}

	var uris []*url.URL
	for _, v := range certConf.AlternateURLs {
		if u, err := url.Parse(v); err != nil {
			return nil, err
		} else {
			uris = append(uris, u)
		}
	}

	var serial *big.Int
	if serial, err = rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128)); err != nil {
		return nil, err
	}

	var notBefore, notAfter time.Time
	if notBefore, err = time.Parse("2006-01-02 15:04:05", certConf.NotBefore); err != nil {
		return nil, err
	}
	if notAfter, err = time.Parse("2006-01-02 15:04:05", certConf.NotAfter); err != nil {
		return nil, err
	}

	cert := &x509.Certificate{
		SerialNumber:          serial,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		BasicConstraintsValid: true,
		DNSNames:              certConf.AlternateDNSNames,
		IPAddresses:           ips,
		URIs:                  uris,
		Subject: pkix.Name{
			Country:            certConf.Country,
			Locality:           certConf.City,
			Province:           certConf.Province,
			OrganizationalUnit: certConf.Section,
			Organization:       certConf.Company,
			StreetAddress:      certConf.Street,
			PostalCode:         certConf.PostCode,
			CommonName:         certConf.CommonName,
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,
	}

	return cert, nil
}

func GenCertificatePem(filePath string, issuerTemp, signerCrt *x509.Certificate, signerPriv interface{}) (issuerCrt *x509.Certificate, issuerDer []byte, issuerPriv *rsa.PrivateKey, err error) {
	if issuerTemp == nil {
		return nil, nil, nil, comerr.ParamInvalid
	}

	issuerPriv, err = secure.GenRSAPrivateKey(2048)
	if err != nil {
		return nil, nil, nil, err
	}

	dir := path.Dir(filePath)
	if !file.IsDirExists(dir) {
		os.MkdirAll(dir, 0755)
	}

	// write private key pem file
	bc, err := KeyBlock(issuerPriv, kf_pkcs8_private)
	if err != nil {
		return nil, nil, nil, err
	}
	privPem, err := os.Create(dir + string(os.PathSeparator) + "priv.key")
	if err != nil {
		return nil, nil, nil, err
	}
	defer privPem.Close()
	WritePemTo(privPem, bc)
	// write public key pem file
	bc, err = KeyBlock(&issuerPriv.PublicKey, kf_pkcs8_public)
	if err != nil {
		return nil, nil, nil, err
	}
	pubPem, err := os.Create(dir + string(os.PathSeparator) + "pub.key")
	if err != nil {
		return nil, nil, nil, err
	}
	defer pubPem.Close()
	WritePemTo(pubPem, bc)

	// set certificate parameters
	issuerTemp.PublicKeyAlgorithm = x509.RSA
	sbuf := sha1.Sum(issuerPriv.N.Bytes())
	issuerTemp.SubjectKeyId = sbuf[:]
	if signerCrt == nil {
		issuerTemp.IsCA = true
		signerCrt = issuerTemp
		signerPriv = issuerPriv
	}

	if issuerDer, err = x509.CreateCertificate(rand.Reader, issuerTemp, signerCrt, &issuerPriv.PublicKey, signerPriv); err != nil {
		return nil, nil, nil, err
	}

	// write certificate pem file
	certPem, err := os.Create(filePath)
	if err != nil {
		return nil, nil, nil, err
	}
	if err = WritePemTo(certPem, &pem.Block{
		Type:  string(kf_certificate),
		Bytes: issuerDer,
	}); err != nil {
		return nil, nil, nil, err
	}

	issuerCrt, err = x509.ParseCertificate(issuerDer)

	return
}
