package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/koltyakov/gosip"
)

// HTTPClient HTTP methods helper
type HTTPClient struct {
	sp *gosip.SPClient
}

// RequestConfig struct
type RequestConfig struct {
	Headers map[string]string
}

// HeadersPresets : SP REST OData headers presets
var HeadersPresets = struct {
	Verbose         *RequestConfig
	Minimalmetadata *RequestConfig
	Nometadata      *RequestConfig
}{
	Verbose: &RequestConfig{
		Headers: map[string]string{
			"Accept":       "application/json;odata=verbose",
			"Content-Type": "application/json;odata=verbose;charset=utf-8",
		},
	},
	Minimalmetadata: &RequestConfig{
		Headers: map[string]string{
			"Accept":       "application/json;odata=minimalmetadata",
			"Content-Type": "application/json;odata=minimalmetadata;charset=utf-8",
		},
	},
	Nometadata: &RequestConfig{
		Headers: map[string]string{
			"Accept":       "application/json;odata=nometadata",
			"Content-Type": "application/json;odata=nometadata;charset=utf-8",
		},
	},
}

// NewHTTPClient creates an instance of httpClient
func NewHTTPClient(spClient *gosip.SPClient) *HTTPClient {
	return &HTTPClient{sp: spClient}
}

// Get - generic GET request wrapper
func (ctx *HTTPClient) Get(endpoint string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create a request: %v", err)
	}

	// Default headers
	req.Header.Set("Accept", "application/json;odata=verbose") // default to SP2013 for backwards compatibility

	// Apply custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := ctx.sp.Execute(req)
	if err != nil {
		return nil, fmt.Errorf("unable to request api: %v", err)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// Post - generic POST request wrapper
func (ctx *HTTPClient) Post(endpoint string, body []byte, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("unable to create a request: %v", err)
	}

	// Default headers
	req.Header.Set("Accept", "application/json;odata=verbose") // default to SP2013 for backwards compatibility
	req.Header.Set("Content-Type", "application/json;odata=verbose;charset=utf-8")

	// Apply custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := ctx.sp.Execute(req)
	if err != nil {
		return nil, fmt.Errorf("unable to request api: %v", err)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// Delete - generic DELETE request wrapper
func (ctx *HTTPClient) Delete(endpoint string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create a request: %v", err)
	}

	// Default headers
	req.Header.Set("Accept", "application/json;odata=verbose") // default to SP2013 for backwards compatibility
	req.Header.Set("Content-Type", "application/json;odata=verbose;charset=utf-8")
	req.Header.Add("X-Http-Method", "DELETE")
	req.Header.Add("If-Match", "*")

	// Apply custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := ctx.sp.Execute(req)
	if err != nil {
		return nil, fmt.Errorf("unable to request api: %v", err)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// Update - generic MERGE request wrapper
func (ctx *HTTPClient) Update(endpoint string, body []byte, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("unable to create a request: %v", err)
	}

	// Default headers
	req.Header.Set("Accept", "application/json;odata=verbose") // default to SP2013 for backwards compatibility
	req.Header.Set("Content-Type", "application/json;odata=verbose;charset=utf-8")
	req.Header.Add("X-Http-Method", "MERGE")
	req.Header.Add("If-Match", "*")

	// Apply custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := ctx.sp.Execute(req)
	if err != nil {
		return nil, fmt.Errorf("unable to request api: %v", err)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// ProcessQuery - CSOM requests helper
func (ctx *HTTPClient) ProcessQuery(endpoint string, body []byte) ([]byte, error) {
	if strings.Index(strings.ToLower(endpoint), strings.ToLower("/_vti_bin/client.svc/ProcessQuery")) == -1 {
		endpoint = fmt.Sprintf("%s/_vti_bin/client.svc/ProcessQuery", strings.Split(endpoint, "/_api")[0])
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("unable to create a request: %v", err)
	}

	// CSOM headers
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", `text/xml;charset="UTF-8"`)
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	resp, err := ctx.sp.Execute(req)
	if err != nil {
		return nil, fmt.Errorf("unable to request api: %v", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := []*struct {
		SchemaVersion  string `json:"SchemaVersion"`
		LibraryVersion string `json:"LibraryVersion"`
		ErrorInfo      *struct {
			ErrorMessage  string `json:"ErrorMessage"`
			ErrorValue    string `json:"ErrorValue"`
			ErrorCode     int    `json:"ErrorCode"`
			ErrorTypeName string `json:"ErrorTypeName"`
		} `json:"ErrorInfo"`
		TraceCorrelationID string `json:"TraceCorrelationId"`
	}{}

	if err := json.Unmarshal(data, &res); err != nil {
		return data, err
	}

	if res[0].ErrorInfo != nil {
		return data, fmt.Errorf(
			"%s (Code: %d, %s, Correlation ID: %s)",
			res[0].ErrorInfo.ErrorMessage,
			res[0].ErrorInfo.ErrorCode,
			res[0].ErrorInfo.ErrorTypeName,
			res[0].TraceCorrelationID,
		)
	}

	return data, nil
}
