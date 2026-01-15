package httpclient

import "crypto/tls"

type NewClientOption func(*Client)

func WithConfig(config *Config) NewClientOption {
	return func(client *Client) {
		client.Config = config
	}
}

func WithTLSClientConfig(tlsConfig *tls.Config) NewClientOption {
	return func(client *Client) {
		client.TLSClientConfig = tlsConfig
	}
}
