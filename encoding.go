package gofast

import (
	"github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
)

type RequestEncoder func(req *fasthttp.Request, in any) error

type ResponseDecoder func(resp *fasthttp.Response, out any) error

var JSONEncoder = func(req *fasthttp.Request, in any) error {
	req.Header.SetContentType("application/json")
	return json.NewEncoder(req.BodyWriter()).Encode(in)
}

var JSONDecoder = func(resp *fasthttp.Response, out any) error {
	if err := json.Unmarshal(resp.Body(), out); err != nil {
		//log.Printf("[gofast] response decode failed - code: %v, body: %v", resp.StatusCode(), string(resp.Body()))
		return err
	}
	return nil
}

var URLEncoder = func(req *fasthttp.Request, in any) error {
	args := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(args)

	for k, v := range in.(map[string]string) {
		args.Set(k, v)
	}
	if _, err := args.WriteTo(req.BodyWriter()); err != nil {
		return err
	}
	req.Header.SetContentType("application/x-www-form-urlencoded")
	return nil
}

var TextDecoder = func(resp *fasthttp.Response, out any) error {
	s := out.(*string)
	*s = string(resp.Body())
	return nil
}
