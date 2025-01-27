package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type AwxClient struct {
	client   *http.Client
	endpoint string
	auth     string
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
	httpReq.Header.Add("Authorization", c.auth)

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
	httpReq.Header.Add("Authorization", c.auth)

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
	httpReq.Header.Add("Authorization", c.auth)

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
	httpReq.Header.Add("Authorization", c.auth)

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
	httpReq.Header.Add("Authorization", c.auth)

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
	httpReq.Header.Add("Authorization", c.auth)

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
	httpReq.Header.Add("Authorization", c.auth)

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
