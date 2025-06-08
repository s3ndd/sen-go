package client

import "time"

type Config struct {
	EndpointURL             string
	UseSecure               bool
	IgnoreCertificateErrors bool
	Timeout                 time.Duration
}
