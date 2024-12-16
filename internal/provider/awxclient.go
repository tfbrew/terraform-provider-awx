package provider

import (
	"net/http"
)

type AwxClient struct {
	client   *http.Client
	endpoint string
	token    string
}

func (c *AwxClient) DoSomething(varname string) string {
	return "hi"
}
