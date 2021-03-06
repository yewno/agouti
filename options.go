package agouti

import (
	"net/http"
	"time"
)

type config struct {
	Proxy               string
	Timeout             time.Duration
	DesiredCapabilities Capabilities
	BrowserName         string
	RejectInvalidSSL    bool
	Debug               bool
	HTTPClient          *http.Client
}

// An Option specifies configuration for a new WebDriver or Page.
type Option func(*config)

// Browser provides an Option for specifying a browser.
func Browser(name string) Option {
	return func(c *config) {
		c.BrowserName = name
	}
}

func Proxy(address string) Option {
	return func(c *config) {
		c.Proxy = address
	}
}

// Timeout provides an Option for specifying a timeout in seconds.
func Timeout(seconds int) Option {
	return func(c *config) {
		c.Timeout = time.Duration(seconds) * time.Second
	}
}

// Desired provides an Option for specifying desired WebDriver Capabilities.
func Desired(capabilities Capabilities) Option {
	return func(c *config) {
		c.DesiredCapabilities = capabilities
	}
}

// RejectInvalidSSL is an Option specifying that the WebDriver should reject
// invalid SSL certificates. All WebDrivers should accept invalid SSL certificates
// by default. See: http://www.w3.org/TR/webdriver/#invalid-ssl-certificates
var RejectInvalidSSL Option = func(c *config) {
	c.RejectInvalidSSL = true
}

// Debug is an Option that connects the running WebDriver to stdout and stdin.
var Debug Option = func(c *config) {
	c.Debug = true
}

// HTTPClient provides an Option for specifying a *http.Client
func HTTPClient(client *http.Client) Option {
	return func(c *config) {
		c.HTTPClient = client
	}
}

func (c config) Merge(options []Option) *config {
	for _, option := range options {
		option(&c)
	}
	return &c
}

func (c *config) Capabilities() Capabilities {
	merged := Capabilities{"acceptSslCerts": true}
	for feature, value := range c.DesiredCapabilities {
		merged[feature] = value
	}
	if c.BrowserName != "" {
		merged.Browser(c.BrowserName)
	}
	if c.RejectInvalidSSL {
		merged.Without("acceptSslCerts")
	}
	return merged
}
