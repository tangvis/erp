package n11

import "github.com/go-resty/resty/v2"

type Client struct {
	C *resty.Client
}

func NewClient() *Client {
	c := resty.New()
	// Registering Request Middleware
	c.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		// Now you have access to Client and current Request object
		// manipulate it as per your need

		return nil // if its success otherwise return error
	})

	// Registering Response Middleware
	c.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		// Now you have access to Client and current Response object
		// manipulate it as per your need

		return nil // if its success otherwise return error
	})
	return &Client{
		C: c,
	}
}
