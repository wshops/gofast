package gofast

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type Header map[string]string

//type Body map[string]string

type Client struct {
	fastClient      *fasthttp.Client
	errorHandler    ErrorHandler
	requestEncoder  RequestEncoder
	responseDecoder ResponseDecoder
}

func New(config ...Config) *Client {
	cfg := configDefault(config...)

	return &Client{
		fastClient: &fasthttp.Client{
			Name:                     cfg.Name,
			NoDefaultUserAgentHeader: cfg.NoDefaultUserAgentHeader,
			ReadTimeout:              cfg.ReadTimeout,
			WriteTimeout:             cfg.WriteTimeout,
		},
		errorHandler:    cfg.ErrorHandler,
		requestEncoder:  cfg.RequestEncoder,
		responseDecoder: cfg.ResponseDecoder,
	}
}

func (c *Client) Get(uri string, out any, header Header) (error, int) {
	return c.do(uri, fasthttp.MethodGet, nil, out, header)
}

func (c *Client) Post(uri string, in, out any, header Header) (error, int) {
	return c.do(uri, fasthttp.MethodPost, in, out, header)
}

func (c *Client) PostForm(uri string, in map[string]string, out any, header Header) (error, int) {
	cfg := configDefault()
	cfg.RequestEncoder = URLEncoder
	return New(cfg).do(uri, fasthttp.MethodPost, in, out, header)
}

func (c *Client) Put(uri string, in, out any, header Header) (error, int) {
	return c.do(uri, fasthttp.MethodPut, in, out, header)
}

func (c *Client) Patch(uri string, in, out any, header Header) (error, int) {
	return c.do(uri, fasthttp.MethodPatch, in, out, header)
}

func (c *Client) Delete(uri string, in, out any, header Header) (error, int) {
	return c.do(uri, fasthttp.MethodDelete, in, out, header)
}

func (c *Client) do(uri string, method string, in, out any, header Header) (err error, code int) {
	req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()

	req.SetRequestURI(uri)
	req.Header.SetMethod(method)
	for k, v := range header {
		req.Header.Set(k, v)
	}

	if in != nil {
		if err := c.requestEncoder(req, in); err != nil {
			return fmt.Errorf("encode request: %w", err), 0
		}
	}

	if err := c.fastClient.Do(req, resp); err != nil {
		return fmt.Errorf("send request: %w", err), 0
	}

	//if resp.StatusCode() < 200 || resp.StatusCode() > 299 {
	//	return c.errorHandler(resp), resp.StatusCode()
	//}

	if out != nil {
		if err := c.responseDecoder(resp, out); err != nil {
			return fmt.Errorf("decode response: %w", err), 0
		}
	}

	return nil, resp.StatusCode()
}
