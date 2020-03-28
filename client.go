package backlog

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

const (
	apiVersion = "v2"
)

// ClinetError is a description of a Backlog API client error.
type ClinetError struct {
	msg string
}

func (e *ClinetError) Error() string {
	return e.msg
}

func newClientError(msg string) *ClinetError {
	return &ClinetError{msg: msg}
}

// Client is Backlog API client.
type Client struct {
	url        *url.URL
	httpClient *http.Client
	token      string

	Activity    *ActivityService
	Category    *CategoryService
	CustomField *CustomFieldService
	Issue       *IssueService
	Priority    *PriorityService
	Project     *ProjectService
	PullRequest *PullRequestService
	Resolution  *ResolutionService
	Space       *SpaceService
	Status      *StatusService
	User        *UserService
	Version     *VersionService
	Wiki        *WikiService
}

// Response represents Backlog API response.
// It wraps http.Response.
type response struct {
	*http.Response
	Error *APIResponseError
}

// Request wraps http.Request.
type request struct {
	*http.Request
}

// RequestParams wraps url.Values.
type requestParams struct {
	*url.Values
}

type clientGet func(spath string, params *requestParams) (*response, error)
type clientPost func(spath string, params *requestParams) (*response, error)
type clientPatch func(spath string, params *requestParams) (*response, error)
type clientDelete func(spath string, params *requestParams) (*response, error)
type clientUploade func(spath, fPath, fName string) (*response, error)

type clientMethod struct {
	Get     clientGet
	Post    clientPost
	Patch   clientPatch
	Delete  clientDelete
	Uploade clientUploade
}

// NewClient creates a new Backlog API Client.
func NewClient(baseURL, token string) (*Client, error) {
	if len(token) == 0 {
		return nil, newClientError("missing token")
	}

	parsedURL, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		url:        parsedURL,
		httpClient: http.DefaultClient,
		token:      token,
	}

	cm := &clientMethod{
		Get: func(spath string, params *requestParams) (*response, error) {
			return c.get(spath, params)
		},
		Post: func(spath string, params *requestParams) (*response, error) {
			return c.post(spath, params)
		},
		Patch: func(spath string, params *requestParams) (*response, error) {
			return c.patch(spath, params)
		},
		Delete: func(spath string, params *requestParams) (*response, error) {
			return c.delete(spath, params)
		},
		Uploade: func(spath, fPath, fName string) (*response, error) {
			return c.uploade(spath, fPath, fName)
		},
	}

	c.Activity = newActivityService(cm)
	c.Category = newCategoryService(cm)
	c.CustomField = newCustomFieldService(cm)
	c.Issue = newIssueService(cm)
	c.Priority = newPriorityService(cm)
	c.Project = newProjectService(cm)
	c.PullRequest = newPullRequestService(cm)
	c.Resolution = newResolutionService(cm)
	c.Space = newSpaceService(cm)
	c.Status = newStatusService(cm)
	c.User = newUserService(cm)
	c.Version = newVersionService(cm)
	c.Wiki = newWikiService(cm)

	return c, nil
}

// Creates new request.
func (c *Client) newReqest(method, spath string, params *requestParams, body io.Reader) (*request, error) {
	if spath == "" {
		return nil, errors.New("spath must not empty")
	}

	if params == nil {
		params = newRequestParams()
	}
	params.Set("apiKey", c.token)

	u := *c.url
	u.Path = path.Join(u.Path, "api", apiVersion, spath)
	u.RawQuery = params.Encode()

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	return &request{Request: req}, nil
}

// Do http request, and return Response.
func (c *Client) do(req *request) (*response, error) {
	resp, err := c.httpClient.Do(req.Request)
	if err != nil {
		return nil, err
	}

	r := newResponse(resp)

	return checkResponseError(r)
}

// Get method of http reqest.
// It creates new http reqest and do and return Response.
func (c *Client) get(spath string, params *requestParams) (*response, error) {
	req, err := c.newReqest(http.MethodGet, spath, params, nil)
	if err != nil {
		return nil, err
	}

	return c.do(req)
}

// Post method of http reqest.
// It creates new http reqest and do and return Response.
func (c *Client) post(spath string, params *requestParams) (*response, error) {
	if params == nil {
		params = newRequestParams()
	}
	req, err := c.newReqest(http.MethodPost, spath, nil, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.do(req)
}

// Patch method of http reqest.
// It creates new http reqest and do and return Response.
func (c *Client) patch(spath string, params *requestParams) (*response, error) {
	if params == nil {
		params = newRequestParams()
	}
	req, err := c.newReqest(http.MethodPatch, spath, nil, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.do(req)
}

// Delete method of http reqest.
// It creates new http reqest and do and return Response.
func (c *Client) delete(spath string, params *requestParams) (*response, error) {
	if params == nil {
		params = newRequestParams()
	}
	req, err := c.newReqest(http.MethodDelete, spath, nil, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.do(req)
}

// Uploade file method used http reqest.
// It creates new http reqest and do and return Response.
func (c *Client) uploade(spath, fPath, fName string) (*response, error) {
	if fPath == "" || fName == "" {
		return nil, newClientError("file's path and name is required")
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	f, err := os.Open(fPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fw, err := w.CreateFormFile("file", fName)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, f); err != nil {
		return nil, err
	}
	w.Close()

	req, err := c.newReqest(http.MethodPost, spath, nil, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	return c.do(req)
}

// Create new parameter for request.
func newRequestParams() *requestParams {
	return &requestParams{&url.Values{}}
}

// Creates new Response.
func newResponse(resp *http.Response) *response {
	r := &response{
		Response: resp,
		Error:    &APIResponseError{},
	}

	return r
}

// Check HTTP status code. If it has errors, return error.
func checkResponseError(r *response) (*response, error) {
	if sc := r.StatusCode; 200 <= sc && sc <= 299 {
		return r, nil
	}

	if r.Body == nil {
		return nil, newClientError("response body is empty")
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(r.Error); err != nil {
		return nil, err
	}

	return nil, r.Error
}
