package hagent

import "fmt"

type QueryParam struct {
	Key   string
	Value any
}

func (q QueryParam) Map() map[string]string {
	return map[string]string{
		q.Key: fmt.Sprintf("%v", q.Value),
	}
}

type QueryParams []QueryParam

func (q QueryParams) Map() map[string]string {
	queryParams := make(map[string]string)
	for _, queryParam := range q {
		queryParams[queryParam.Key] = fmt.Sprintf("%v", queryParam.Value)
	}
	return queryParams
}

func (q *QueryParams) Add(key string, value any) {
	*q = append(*q, QueryParam{Key: key, Value: value})
}
