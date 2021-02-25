package seelesdk

import (
	gosdktypes "github.com/Seele-N/go-sdk/types"
)

// RestClient restful api client
type RestClient struct {
	config   gosdktypes.ClientConfig
	cdc      *gosdktypes.Codec
	appCodec gosdktypes.AppCodec
}

// NewRestClient new restful api client
func NewRestClient(config gosdktypes.ClientConfig) *RestClient {
	cdc, appCodec := gosdktypes.NewCodec()
	client := &RestClient{
		config:   config,
		cdc:      cdc,
		appCodec: appCodec,
	}
	return client
}
