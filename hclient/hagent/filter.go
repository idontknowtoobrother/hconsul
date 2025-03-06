package hagent

import "fmt"

type Filter struct {
	str string

	// private available
	ID         string
	Service    string
	APIVersion string
}

func NewFilterService(service Service) *Filter {
	return &Filter{
		str: "",
	}
}

func (f *Filter) And() string {
	f.str += " and "
	return f.str
}

func (f *Filter) Or() string {
	f.str += " or "
	return f.str
}

func (f *Filter) Write(str string) string {
	f.str += str
	return f.str
}

func (f *Filter) Build() string {
	if f.Service != "" {
		f.Write(fmt.Sprintf("ServiceName == \"%s\"", f.Service))
	}

	if f.ID != "" {
		f.And()
		f.Write(fmt.Sprintf("ServiceMeta.id == \"%s\"", f.ID))
	}

	if f.APIVersion != "" {
		f.And()
		f.Write(fmt.Sprintf("ServiceMeta.api_version == \"%s\"", f.APIVersion))
	}

	return f.str
}

func (f *Filter) String() string {
	if f.str == "" {
		f.Build()
	}
	return f.str
}
