package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type AwxClient struct {
	client   *http.Client
	endpoint string
	token    string
}

// A wrapper for http.NewRequestWithContext() that prepends tower endpoint to URL & sets authorization
// headers and then makes the actual http request.
func (c *AwxClient) MakeHTTPRequestToAPI(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	url = c.endpoint + url
	httpReq, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+c.token)

	//TODO Add if body != nil, set body of http request

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode != 200 && httpResp.StatusCode != 404 {
		return nil, fmt.Errorf("API call's HTTP status wasn't 200 or 404, was %vd", httpResp.StatusCode)
	}

	return httpResp, nil
}

func (c *AwxClient) AssocJobTemplChild(ctx context.Context, body ChildResult, url string) error {

	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+c.token)

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	if httpResp.StatusCode != 204 {
		err = fmt.Errorf("expected http code 204, got %d", httpResp.StatusCode)
		return err
	}

	return nil
}

func (c *AwxClient) AssocSuccessNode(ctx context.Context, body ChildAssocBody, url string) error {

	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+c.token)

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	if httpResp.StatusCode != 204 {
		err = fmt.Errorf("expected http code 204, got %d", httpResp.StatusCode)
		return err
	}

	return nil
}

func (c *AwxClient) DisassocJobTemplChild(ctx context.Context, body ChildDissasocBody, url string) error {

	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+c.token)

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	if httpResp.StatusCode != 204 {
		err = fmt.Errorf("expected http code 204, got %d", httpResp.StatusCode)
		return err
	}

	return nil
}

func (c *AwxClient) AssocJobTemplCredential(ctx context.Context, id int, body Result) error {
	url := c.endpoint + fmt.Sprintf("/api/v2/job_templates/%d/credentials/", id)

	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+c.token)

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	if httpResp.StatusCode != 204 {
		err = fmt.Errorf("expected http code 204, got %d", httpResp.StatusCode)
		return err
	}

	return nil
}

func (c *AwxClient) DisassocJobTemplCredential(ctx context.Context, id int, body DissasocBody) error {
	url := c.endpoint + fmt.Sprintf("/api/v2/job_templates/%d/credentials/", id)

	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+c.token)

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	if httpResp.StatusCode != 204 {
		err = fmt.Errorf("expected http code 204, got %d", httpResp.StatusCode)
		return err
	}

	return nil
}

func (c *AwxClient) AssocJobTemplLabel(ctx context.Context, id int, body LabelResult) error {
	url := c.endpoint + fmt.Sprintf("/api/v2/job_templates/%d/labels/", id)

	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+c.token)

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	if httpResp.StatusCode != 204 {
		err = fmt.Errorf("expected http code 204, got %d", httpResp.StatusCode)
		return err
	}

	return nil
}

func (c *AwxClient) DisassocJobTemplLabel(ctx context.Context, id int, body LabelDissasocBody) error {
	url := c.endpoint + fmt.Sprintf("/api/v2/job_templates/%d/labels/", id)

	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+c.token)

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	if httpResp.StatusCode != 204 {
		err = fmt.Errorf("expected http code 204, got %d", httpResp.StatusCode)
		return err
	}

	return nil
}

// func (c *AwxClient) GetByUrl(ctx context.Context, url string) *http.Response {

// 	return nil
// }
