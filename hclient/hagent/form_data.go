package hagent

import "fmt"

type FormData struct {
	Key   string
	Value any
}

func (f FormData) Map() map[string]string {
	return map[string]string{
		f.Key: fmt.Sprintf("%v", f.Value),
	}
}

type FormDataList []FormData

func (f FormDataList) Map() map[string]string {
	result := make(map[string]string)
	for _, formData := range f {
		result[formData.Key] = fmt.Sprintf("%v", formData.Value)
	}
	return result
}

func (f *FormDataList) Add(key string, value any) {
	*f = append(*f, FormData{Key: key, Value: value})
}
