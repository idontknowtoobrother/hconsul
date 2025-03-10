package hagent

import (
	"resty.dev/v3"
)

type Request struct {
	rcl *resty.Client
}

func (d *hAgent) NewRequest(
	service Service,
	id ServiceID,
) (*Request, error) {
	addr, err := d.DiscoveryServiceId(service, id)
	if err != nil {
		return nil, err
	}

	rcl := resty.New().SetBaseURL(addr)

	return &Request{
		rcl: rcl,
	}, nil
}

func (r *Request) R() *resty.Request {
	return r.rcl.R()
}
