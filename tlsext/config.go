package tlsext

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"strings"
)

// TlsClientConfig represents the standard client TLS config.
type TlsClientConfig struct {
	CaCerts            []string `json:"ca_certs" toml:"ca_certs"`
	Cert               string   `json:"cert" toml:"cert"`
	CertKey            string   `json:"cert_key" toml:"cert_key"`
	InsecureSkipVerify bool     `json:"insecure_skip_verify" toml:"insecure_skip_verify"`
	ServerName         string   `json:"server_name" toml:"server_name"`
}

// TLSConfig returns a tls.Config, may be nil without error if TLS is not configured.
func (this *TlsClientConfig) TlsConfig() (*tls.Config, error) {
	// This check returns a nil (aka, "use the default")
	// tls.Config if no field is set that would have an effect on
	// a TLS connection. That is, any of:
	//     * client certificate settings,
	//     * peer certificate authorities,
	//     * disabled security, or
	//     * an SNI server name.
	if len(this.CaCerts) == 0 && this.CertKey == "" && this.Cert == "" && !this.InsecureSkipVerify && this.ServerName == "" {
		return nil, nil
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: this.InsecureSkipVerify,
		Renegotiation:      tls.RenegotiateNever,
	}

	if len(this.CaCerts) != 0 {
		if pool, err := MakeCertPool(this.CaCerts); err != nil {
			return nil, err
		} else {
			tlsConfig.RootCAs = pool
		}
	}

	if this.Cert != "" && this.CertKey != "" {
		if err := loadCertificate(tlsConfig, this.Cert, this.CertKey); err != nil {
			return nil, err
		}
	}
	if this.ServerName != "" {
		tlsConfig.ServerName = this.ServerName
	}

	return tlsConfig, nil
}

// ServerConfig represents the standard server TLS config.
type ServerConfig struct {
	Cert           string   `json:"cert" toml:"cert"`
	CertKey        string   `json:"cert_key" toml:"cert_key"`
	AllowedCaCerts []string `json:"allowed_ca_certs" toml:"allowed_ca_certs"`
	CipherSuites   []string `json:"cipher_suites" toml:"cipher_suites"`
	TlsMaxVersion  string   `json:"tls_max_version" toml:"tls_max_version"`
	TlsMinVersion  string   `json:"tls_min_version" toml:"tls_min_version"`
}

// TLSConfig returns a tls.Config, may be nil without error if TLS is not
// configured.
func (this *ServerConfig) TlsConfig() (*tls.Config, error) {
	if this.Cert == "" && this.CertKey == "" && len(this.AllowedCaCerts) == 0 {
		return nil, nil
	}

	tlsConfig := &tls.Config{}

	if len(this.AllowedCaCerts) != 0 {
		if pool, err := MakeCertPool(this.AllowedCaCerts); err != nil {
			return nil, err
		} else {
			tlsConfig.ClientCAs = pool
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}
	}

	if this.Cert != "" && this.CertKey != "" {
		if err := loadCertificate(tlsConfig, this.Cert, this.CertKey); err != nil {
			return nil, err
		}
	}

	if len(this.CipherSuites) != 0 {
		if cipherSuites, err := ParseCiphers(this.CipherSuites); err != nil {
			return nil, fmt.Errorf("could not parse server cipher suites %s: %v", strings.Join(this.CipherSuites, ","), err)
		} else {
			tlsConfig.CipherSuites = cipherSuites
		}
	}

	if this.TlsMaxVersion != "" {
		if version, err := ParseTLSVersion(this.TlsMaxVersion); err != nil {
			return nil, fmt.Errorf("could not parse tls max version %q: %v", this.TlsMaxVersion, err)
		} else {
			tlsConfig.MaxVersion = version
		}
	}
	if this.TlsMinVersion != "" {
		if version, err := ParseTLSVersion(this.TlsMinVersion); err != nil {
			return nil, fmt.Errorf("could not parse tls min version %q: %v", this.TlsMinVersion, err)
		} else {
			tlsConfig.MinVersion = version
		}
	}

	if tlsConfig.MinVersion != 0 && tlsConfig.MaxVersion != 0 && tlsConfig.MinVersion > tlsConfig.MaxVersion {
		return nil, fmt.Errorf("tls min version %q can't be greater than tls max version %q", tlsConfig.MinVersion, tlsConfig.MaxVersion)
	}

	return tlsConfig, nil
}

func loadCertificate(config *tls.Config, certFile, keyFile string) error {
	if cert, err := tls.LoadX509KeyPair(certFile, keyFile); err != nil {
		return fmt.Errorf("could not load keypair %s:%s: %v\n", certFile, keyFile, err)
	} else {
		config.Certificates = []tls.Certificate{cert}
		config.BuildNameToCertificate()

		return nil
	}
}

func MakeCertPool(certFiles []string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	for _, certFile := range certFiles {
		if pem, err := ioutil.ReadFile(certFile); err != nil {
			return nil, fmt.Errorf("could not read certificate %q: %v", certFile, err)
		} else {
			if ok := pool.AppendCertsFromPEM(pem); !ok {
				return nil, fmt.Errorf("could not parse any PEM certificates %q: %v", certFile, err)
			}
		}
	}

	return pool, nil
}

func ParseCiphers(ciphers []string) ([]uint16, error) {
	var (
		supported = tls.CipherSuites()
		suites    []uint16
	)
	for _, cipher := range ciphers {
		for _, suite := range supported {
			if cipher == suite.Name {
				suites = append(suites, suite.ID)
			} else {
				return nil, fmt.Errorf("unsupported cipher %q", cipher)
			}
		}
	}

	return suites, nil
}

// "TLS10": tls.VersionTLS10
// "TLS11": tls.VersionTLS11
// "TLS12": tls.VersionTLS12
// "TLS13": tls.VersionTLS13
// "TLS30": tls.VersionTLS30
func ParseTLSVersion(version string) (uint16, error) {
	switch version {
	case "TLS10":
		return tls.VersionTLS10, nil
	case "TLS11":
		return tls.VersionTLS11, nil
	case "TLS12":
		return tls.VersionTLS12, nil
	case "TLS13":
		return tls.VersionTLS13, nil
	// case "TLS30":
	// 	return tls.VersionSSL30, nil
	default:
		return 0, fmt.Errorf("unsupported TLS version: %q", version)
	}
}
