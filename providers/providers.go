package providers

import (
	"errors"
	"fmt"
)

type Provider string

const (
	Cloudfront Provider = "cloudfront"
	Cloudflare Provider = "cloudflare"
)

var List = map[Provider]Provider{
	Cloudfront: Cloudfront,
	Cloudflare: Cloudflare,
}

func (p Provider) String() string {
	if v, ok := List[p]; ok {
		return string(v)
	}
	return ""
}

func (p *Provider) Validate() error {
	if _, ok := List[*p]; !ok {
		return errors.New(fmt.Sprint("invalid value", *p))
	}
	return nil
}
