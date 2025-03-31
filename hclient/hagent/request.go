package hagent

import (
	"fmt"

	"resty.dev/v3"
)

type Request struct {
	rcl *resty.Client
}

func (d *hAgent) NewRequest(
	schema string,
	service Service,
	id ServiceID,
) (*Request, error) {
	addr, err := d.DiscoveryServiceId(service, id)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s://%s", schema, addr)
	rcl := resty.New().SetBaseURL(url)

	return &Request{
		rcl: rcl,
	}, nil
}

func (r *Request) R() *resty.Request {
	return r.rcl.R()
}
