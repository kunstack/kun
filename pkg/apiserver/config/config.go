/*
 * Copyright 2021 Aapeli.Smith<aapeli.nian@gmail.com>.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"fmt"
	"github.com/aapelismith/kun/pkg/auth"
	"github.com/aapelismith/kun/pkg/log"
	"github.com/aapelismith/kun/pkg/types"
	"github.com/spf13/pflag"
	"golang.org/x/exp/slices"
	"k8s.io/apimachinery/pkg/util/errors"
	"net/http"
	"os"
	"time"
)

var (
	// 'request'	Ask clients for a certificate, but allow even if there isn't one; do not verify it
	// 'require'	Require clients to present a certificate, but do not verify it
	// 'verify_if_given'	Ask clients for a certificate; allow even if there isn't one, but verify it if there is
	// 'require_and_verify'	Require clients to present a valid certificate that is verified
	certificateAuthMode = []string{"request", "require", "verify_if_given", "require_and_verify"}
)

type GlobalOptions struct {
	// CacheDefaultExpiration given default expiration duration, If the expiration duration
	// is less than one the items in the cache never expire
	CacheDefaultExpiration types.Duration `yaml:"cache_default_expiration,omitempty" json:"cache_default_expiration,omitempty"`

	// CacheCleanupInterval If the cleanup interval is less than one, expired items are not
	// deleted from the cache before calling c.DeleteExpired().
	CacheCleanupInterval types.Duration `yaml:"cache_cleanup_interval,omitempty" json:"cache_cleanup_interval,omitempty"`
}

type PeerOptions struct {
	// BindAddr address bound for the peer server
	BindAddr string `yaml:"bind_addr,omitempty" json:"bind_addr,omitempty"`

	// MemberURLS list of the member's URLs for this cluster
	MemberURLs []string `yaml:"member_urls,omitempty" json:"member_urls,omitempty"`

	// RootCAPool CA certificate used to validate client certificate.
	// Enables client certificate verification when specified.
	RootCAPool []string `yaml:"root_ca_pool,omitempty" json:"root_ca_pool,omitempty"`

	// CertificateFile Public key file of the peer server.
	CertificateFile string `yaml:"certificate_file,omitempty" json:"certificate_file,omitempty"`

	// CertificateKeyFile  Private key file of the peer server.
	CertificateKeyFile string `yaml:"certificate_key_file,omitempty" json:"certificate_key_file,omitempty"`

	// CertificateAuthMode is the mode for authenticating the client. Allowed values are:
	// 'request'	Ask clients for a certificate, but allow even if there isn't one; do not verify it
	// 'require'	Require clients to present a certificate, but do not verify it
	// 'verify_if_given'	Ask clients for a certificate; allow even if there isn't one, but verify it if there is
	// 'require_and_verify'	Require clients to present a valid certificate that is verified
	CertificateAuthMode string `yaml:"certificate_auth_mode,omitempty" json:"client_auth_mode,omitempty"`

	// ClientCertificateFile Path to the client server TLS cert file.
	ClientCertificateFile string `yaml:"client_certificate_file,omitempty" json:"client_certificate_file,omitempty"`

	// CertificateKeyFile  Private key file of the peer server.
	ClientCertificateKeyFile string `yaml:"client_certificate_key_file,omitempty" json:"client_certificate_key_file,omitempty"`

	// ClientInsecureSkipTLSVerify skip server certificate verification (CAUTION: this option should be enabled only for testing purposes)
	ClientInsecureSkipTLSVerify bool `yaml:"client_insecure_skip_tls_verify,omitempty" json:"client_insecure_skip_tls_verify,omitempty"`
}

// SetDefaults sets the default values.
func (o *PeerOptions) SetDefaults() {
	o.BindAddr = ":7070"
	o.CertificateAuthMode = "request"
	o.MemberURLs = []string{"http://localhost:7070"}
}

// AddFlags add backend related command line parameters
func (o *PeerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.BindAddr, "peer.bind-addr", o.BindAddr, "Peer server listening address")

	fs.StringArrayVar(&o.MemberURLs, "peer.member-urls", o.MemberURLs, "List of the member's URLs for this cluster")

	fs.StringArrayVar(&o.RootCAPool, "peer.root-ca-pool", o.RootCAPool, "Optional list of "+
		"base64-encoded DER-encoded CA certificates file to trust for peer server.")

	fs.StringVar(&o.CertificateFile, "peer.certificate-file", o.CertificateFile, "Path to the peer server TLS cert file.")

	fs.StringVar(&o.CertificateKeyFile, "peer.certificate-key-file", o.CertificateKeyFile, "Path to the peer server TLS key file.")

	fs.StringVar(&o.CertificateAuthMode, "peer.certificate-auth-mode", o.CertificateAuthMode, "The mode for "+
		"authenticating the client for peer server. Allowed values are: \n\t'request' Ask clients for a certificate, "+
		"but allow even if there isn't one; do not verify it, \n\t'require' Require clients to present a certificate, "+
		"but do not verify it \n\t 'verify_if_given' Ask clients for a certificate; allow even if there isn't one, but "+
		"verify it if there is\n\t 'require_and_verify' Require clients to present a valid certificate that is verified")

	fs.StringVar(&o.ClientCertificateFile, "peer.client-certificate-file", o.CertificateFile, "Path to the "+
		"client server TLS cert file.")

	fs.StringVar(&o.ClientCertificateKeyFile, "peer.client-certificate-key-file", o.ClientCertificateKeyFile,
		"Path to the client server TLS key file.")

	fs.BoolVar(&o.ClientInsecureSkipTLSVerify, "peer.client-insecure-skip-tls-verify", o.ClientInsecureSkipTLSVerify,
		"skip server certificate verification (CAUTION: this option should be enabled only for testing purposes)")
}

// Validate verify the configuration and return an error if correct
func (o *PeerOptions) Validate() error {
	if o.BindAddr == "" {
		return fmt.Errorf("bind_addr is required field")
	}

	for _, caFile := range o.RootCAPool {
		stat, err := os.Stat(caFile)
		if err != nil {
			return err
		}
		if !stat.Mode().IsRegular() {
			return fmt.Errorf("trusted_ca_file '%s' is not regular file", caFile)
		}
	}

	if o.CertificateFile != "" && o.CertificateKeyFile == "" {
		return fmt.Errorf("when setting the certificate_file, the certificate_key_file is a required field")
	}

	if o.CertificateFile == "" && o.CertificateKeyFile != "" {
		return fmt.Errorf("when setting the certificate_key_file, the certificate_file is a required field")
	}

	if o.CertificateFile != "" {
		stat, err := os.Stat(o.CertificateFile)
		if err != nil {
			return err
		}
		if !stat.Mode().IsRegular() {
			return fmt.Errorf("certificate_file '%s' is not regular file", o.CertificateFile)
		}
	}

	if o.CertificateKeyFile != "" {
		stat, err := os.Stat(o.CertificateKeyFile)
		if err != nil {
			return err
		}
		if !stat.Mode().IsRegular() {
			return fmt.Errorf("certificate_key_file '%s' is not regular file", o.CertificateKeyFile)
		}
	}

	if o.CertificateAuthMode != "" {
		if ok := slices.Contains(certificateAuthMode, o.CertificateAuthMode); !ok {
			return fmt.Errorf("%s is an unknown tls certificate auth mode", o.CertificateAuthMode)
		}
	}
	return nil
}

// FrontendOptions frontend related configuration
type FrontendOptions struct {
	// HttpBindAddr http listen address bind by the backend service
	HttpBindAddr string `yaml:"http_bind_addr,omitempty" json:"http_bind_addr,omitempty"`

	// HttpsBindAddr https listen address bind by the backend service
	HttpsBindAddr string `yaml:"https_bind_addr,omitempty" json:"https_bind_addr,omitempty"`

	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body. A zero or negative value means
	// there will be no timeout.
	ReadTimeout types.Duration `yaml:"read_timeout,omitempty" json:"read_timeout,omitempty"`

	// ReadHeaderTimeout is the amount of time allowed to read
	// request headers. The connection's read deadline is reset
	// after reading the headers and the Handler can decide what
	// is considered too slow for the body. If ReadHeaderTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, there is no timeout.
	ReadHeaderTimeout types.Duration `yaml:"read_header_timeout,omitempty" json:"read_header_timeout,omitempty"`

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read. Like ReadTimeout, it does not
	// let Handlers make decisions on a per-request basis.
	// A zero or negative value means there will be no timeout.
	WriteTimeout types.Duration `yaml:"write_timeout,omitempty" json:"write_timeout,omitempty"`

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when kee-alive are enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, there is no timeout.
	IdleTimeout types.Duration `yaml:"idle_timeout,omitempty" json:"idle_timeout,omitempty"`

	// MaxHeaderBytes controls the maximum number of bytes the
	// server will read parsing the request header's keys and
	// values, including the request line. It does not limit the
	// size of the request body.
	// If zero, DefaultMaxHeaderBytes is used.
	MaxHeaderBytes int `yaml:"max_header_bytes,omitempty" json:"max_header_bytes,omitempty"`

	// default public key of the frontend service
	DefaultCertificateFile string `yaml:"default_certificate_file,omitempty" json:"default_certificate_file,omitempty"`

	// default private key of the frontend service
	DefaultCertificateKeyFile string `yaml:"default_certificate_key_file,omitempty" json:"default_certificate_key_file,omitempty"`
}

// SetDefaults sets the default values.
func (o *FrontendOptions) SetDefaults() {
	o.HttpBindAddr = ":8080"
	o.HttpsBindAddr = ":8443"
	o.MaxHeaderBytes = http.DefaultMaxHeaderBytes
	o.IdleTimeout = types.Duration(time.Minute * 5)
	o.ReadTimeout = types.Duration(time.Minute * 5)
	o.ReadHeaderTimeout = types.Duration(time.Minute * 5)
}

// AddFlags add backend related command line parameters
func (o *FrontendOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.HttpBindAddr, "frontend.http-bind-addr", o.HttpBindAddr, "Frontend http entrypoint listening address")

	fs.StringVar(&o.HttpsBindAddr, "frontend.https-bind-addr", o.HttpsBindAddr, "Frontend https entrypoint listening address")

	fs.Var(&o.ReadTimeout, "frontend.read-timeout", "ReadTimeout is the maximum duration for reading the "+
		"entire request, including the body. A zero or negative value means there will be no timeout.")

	fs.Var(&o.ReadHeaderTimeout, "frontend.read-header-timeout", " ReadHeaderTimeout is the amount of time"+
		" allowed to read request headers. The connection's read deadline is reset after reading the headers and the "+
		"Handler can decide what is considered too slow for the body. If ReadHeaderTimeout is zero, the value of "+
		"ReadTimeout is used. If both are zero, there is no timeout.")

	fs.Var(&o.IdleTimeout, "frontend.idle-timeout", "IdleTimeout is the maximum amount of time to wait for"+
		" the next request when kee-alive are enabled. If IdleTimeout is zero, the value of ReadTimeout is used. "+
		"If both are zero, there is no timeout.")

	fs.Var(&o.WriteTimeout, "frontend.write-timeout", " WriteTimeout is the maximum duration before timing out"+
		" writes of the response. It is reset whenever a new request's header is read. Like ReadTimeout, it does not"+
		" let Handlers make decisions on a per-request basis. A zero or negative value means there will be no timeout.")

	fs.IntVar(&o.MaxHeaderBytes, "frontend.max-header-bytes", o.MaxHeaderBytes, "MaxHeaderBytes controls the "+
		"maximum number of bytes the frontend.will read parsing the request header's keys and  values, including the"+
		" request line. It does not limit the size of the request body. If zero, DefaultMaxHeaderBytes is used.")
}

// Validate verify the configuration and return an error if correct
func (o *FrontendOptions) Validate() error {
	if o.HttpBindAddr == "" && o.HttpsBindAddr == "" {
		return fmt.Errorf("one of http_address/https_address is required")
	}

	if o.HttpBindAddr == o.HttpsBindAddr {
		return fmt.Errorf("http_address/https_address cannot be the same value")
	}

	if o.DefaultCertificateFile != "" && o.DefaultCertificateKeyFile == "" {
		return fmt.Errorf("when setting the default_certificate_file, " +
			"the default_certificate_key_file is a required field")
	}

	if o.DefaultCertificateFile == "" && o.DefaultCertificateKeyFile != "" {
		return fmt.Errorf("when setting the default_certificate_key_file, " +
			"the default_certificate_file is a required field")
	}

	if o.DefaultCertificateFile != "" {
		stat, err := os.Stat(o.DefaultCertificateFile)
		if err != nil {
			return err
		}
		if !stat.Mode().IsRegular() {
			return fmt.Errorf("default_certificate_file '%s' is not regular file", o.DefaultCertificateFile)
		}
	}

	if o.DefaultCertificateKeyFile != "" {
		stat, err := os.Stat(o.DefaultCertificateKeyFile)
		if err != nil {
			return err
		}
		if !stat.Mode().IsRegular() {
			return fmt.Errorf("default_certificate_key_file '%s' is not regular file", o.DefaultCertificateKeyFile)
		}
	}
	return nil
}

// BackendOptions backend entrypoint  related configuration
type BackendOptions struct {
	// Listen address bind by the backend service
	BindAddr string `yaml:"bind_addr,omitempty" json:"bind_addr,omitempty"`

	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body. A zero or negative value means
	// there will be no timeout.
	ReadTimeout types.Duration `yaml:"read_timeout,omitempty" json:"read_timeout,omitempty"`

	// ReadHeaderTimeout is the amount of time allowed to read
	// request headers. The connection's read deadline is reset
	// after reading the headers and the Handler can decide what
	// is considered too slow for the body. If ReadHeaderTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, there is no timeout.
	ReadHeaderTimeout types.Duration `yaml:"read_header_timeout,omitempty" json:"read_header_timeout,omitempty"`

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read. Like ReadTimeout, it does not
	// let Handlers make decisions on a per-request basis.
	// A zero or negative value means there will be no timeout.
	WriteTimeout types.Duration `yaml:"write_timeout,omitempty" json:"write_timeout,omitempty"`

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when kee-alive are enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, there is no timeout.
	IdleTimeout types.Duration `yaml:"idle_timeout,omitempty" json:"idle_timeout,omitempty"`

	// MaxHeaderBytes controls the maximum number of bytes the
	// server will read parsing the request header's keys and
	// values, including the request line. It does not limit the
	// size of the request body.
	// If zero, DefaultMaxHeaderBytes is used.
	MaxHeaderBytes int `yaml:"max_header_bytes,omitempty" json:"max_header_bytes,omitempty"`

	// RootCAPool CA certificate used to validate client certificate.
	// Enables client certificate verification when specified.
	RootCAPool []string `yaml:"root_ca_pool,omitempty" json:"root_ca_pool,omitempty"`

	// ClientAuthMode is the mode for authenticating the client. Allowed values are:
	// 'request'	Ask clients for a certificate, but allow even if there isn't one; do not verify it
	// 'require'	Require clients to present a certificate, but do not verify it
	// 'verify_if_given'	Ask clients for a certificate; allow even if there isn't one, but verify it if there is
	// 'require_and_verify'	Require clients to present a valid certificate that is verified
	ClientAuthMode string `yaml:"client_auth_mode,omitempty" json:"client_auth_mode,omitempty"`

	// public key of the web service
	CertificateFile string `yaml:"certificate_file,omitempty" json:"certificate_file,omitempty"`

	// private key of the web service
	CertificateKeyFile string `yaml:"certificate_key_file,omitempty" json:"certificate_key_file,omitempty"`
}

// SetDefaults sets the default values.
func (o *BackendOptions) SetDefaults() {
	o.BindAddr = ":8090"
	o.ClientAuthMode = "request"
	o.MaxHeaderBytes = http.DefaultMaxHeaderBytes
	o.IdleTimeout = types.Duration(time.Minute * 5)
	o.ReadTimeout = types.Duration(time.Minute * 5)
	o.ReadHeaderTimeout = types.Duration(time.Minute * 5)
}

// AddFlags add backend related command line parameters
func (o *BackendOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.BindAddr, "server.bind-addr", o.BindAddr, "Server entrypoint listening address")

	fs.StringArrayVar(&o.RootCAPool, "server.root-ca-pool", o.RootCAPool, "Optional list of "+
		"base64-encoded DER-encoded CA certificates file to trust.")

	fs.StringVar(&o.ClientAuthMode, "server.client-auth-mode", o.ClientAuthMode, "The mode for authenticating the client")

	fs.StringVar(&o.CertificateFile, "server.certificate-file", o.CertificateFile, "Public key file of the web service.")

	fs.StringVar(&o.CertificateKeyFile, "server.certificate-key-file", o.CertificateKeyFile, "Private key file of the web service")

	fs.Var(&o.ReadTimeout, "server.read-timeout", "ReadTimeout is the maximum duration for reading the "+
		"entire request, including the body. A zero or negative value means there will be no timeout.")

	fs.Var(&o.ReadHeaderTimeout, "server.read-header-timeout", " ReadHeaderTimeout is the amount of time"+
		" allowed to read request headers. The connection's read deadline is reset after reading the headers and the "+
		"Handler can decide what is considered too slow for the body. If ReadHeaderTimeout is zero, the value of "+
		"ReadTimeout is used. If both are zero, there is no timeout.")

	fs.Var(&o.IdleTimeout, "server.idle-timeout", "IdleTimeout is the maximum amount of time to wait for"+
		" the next request when kee-alive are enabled. If IdleTimeout is zero, the value of ReadTimeout is used. "+
		"If both are zero, there is no timeout.")

	fs.Var(&o.WriteTimeout, "server.write-timeout", " WriteTimeout is the maximum duration before timing out"+
		" writes of the response. It is reset whenever a new request's header is read. Like ReadTimeout, it does not"+
		" let Handlers make decisions on a per-request basis. A zero or negative value means there will be no timeout.")

	fs.IntVar(&o.MaxHeaderBytes, "server.max-header-bytes", o.MaxHeaderBytes, "MaxHeaderBytes controls the "+
		"maximum number of bytes the server will read parsing the request header's keys and  values, including the"+
		" request line. It does not limit the size of the request body. If zero, DefaultMaxHeaderBytes is used.")
}

// Validate verify the configuration and return an error if correct
func (o *BackendOptions) Validate() error {
	if o.BindAddr == "" {
		return fmt.Errorf("bind_addr is required field")
	}

	for _, caFile := range o.RootCAPool {
		stat, err := os.Stat(caFile)
		if err != nil {
			return err
		}
		if !stat.Mode().IsRegular() {
			return fmt.Errorf("trusted_ca_file '%s' is not regular file", caFile)
		}
	}

	if o.CertificateFile != "" && o.CertificateKeyFile == "" {
		return fmt.Errorf("when setting the certificate_file, the certificate_key_file is a required field")
	}

	if o.CertificateFile == "" && o.CertificateKeyFile != "" {
		return fmt.Errorf("when setting the certificate_key_file, the certificate_file is a required field")
	}

	if o.CertificateFile != "" {
		stat, err := os.Stat(o.CertificateFile)
		if err != nil {
			return err
		}
		if !stat.Mode().IsRegular() {
			return fmt.Errorf("certificate_file '%s' is not regular file", o.CertificateFile)
		}
	}

	if o.CertificateKeyFile != "" {
		stat, err := os.Stat(o.CertificateKeyFile)
		if err != nil {
			return err
		}
		if !stat.Mode().IsRegular() {
			return fmt.Errorf("certificate_key_file '%s' is not regular file", o.CertificateKeyFile)
		}
	}

	if o.ClientAuthMode != "" {
		if ok := slices.Contains(certificateAuthMode, o.ClientAuthMode); !ok {
			return fmt.Errorf("%s is an unknown tls client auth mode", o.ClientAuthMode)
		}
	}
	return nil
}

// Configuration Profile contents
type Configuration struct {
	Log      *log.Options     `yaml:"log,omitempty" json:"log,omitempty"`
	Auth     *auth.Options    `yaml:"auth,omitempty" json:"auth,omitempty"`
	Peer     *PeerOptions     `yaml:"peer,omitempty" json:"peer,omitempty"`
	Frontend *FrontendOptions `yaml:"frontend,omitempty" json:"frontend,omitempty"`
	Backend  *BackendOptions  `yaml:"backend,omitempty" json:"backend,omitempty"`
}

// AddFlags   added the configuration  flag to the  specified pflag.FlagSet
func (c *Configuration) AddFlags(fs *pflag.FlagSet) {
	c.Log.AddFlags(fs)
	c.Peer.AddFlags(fs)
	c.Frontend.AddFlags(fs)
	c.Backend.AddFlags(fs)
}

// SetDefaults sets the default values.
func (c *Configuration) SetDefaults() {
	c.Log.SetDefaults()
	c.Peer.SetDefaults()
	c.Frontend.SetDefaults()
	c.Backend.SetDefaults()
}

// Validate Verify that the data in the Configuration meets the requirements
func (c *Configuration) Validate() error {
	errs := make([]error, 0)

	if err := c.Log.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("log: %w", err))
	}

	if err := c.Peer.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("peer: %w", err))
	}

	if err := c.Frontend.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("frontend: %w", err))
	}

	if err := c.Backend.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("backend: %w", err))
	}

	return errors.NewAggregate(errs)
}

// NewConfiguration create Configuration with `zero` value
func NewConfiguration() *Configuration {
	return &Configuration{
		Log:      log.NewOptions(),
		Auth:     auth.NewOptions(),
		Peer:     new(PeerOptions),
		Backend:  new(BackendOptions),
		Frontend: new(FrontendOptions),
	}
}
