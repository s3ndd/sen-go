package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Response struct {
	HTTPResponse *http.Response
	body         []byte
}

func (r *Response) Unmarshal(value interface{}) error {
	body, err := r.Body()
	if err != nil {
		return err
	}
	if len(body) == 0 {
		body = []byte("{}")
	}
	err = json.Unmarshal(body, &value)
	if err != nil {
		return UnmarshalError{
			Body: body,
			Err:  err,
		}
	}
	return nil
}

func (r *Response) StatusCode() int {
	return r.HTTPResponse.StatusCode
}

func (r *Response) Response() *http.Response {
	return r.HTTPResponse
}

func (r *Response) Body() ([]byte, error) {
	if r.body == nil {
		defer r.HTTPResponse.Body.Close()
		body, err := ioutil.ReadAll(r.HTTPResponse.Body)
		if err != nil {
			return nil, err
		}
		r.body = body
	}
	return r.body, nil
}

type UnmarshalError struct {
	Body []byte
	Err  error
}

func (e UnmarshalError) Error() string {
	return fmt.Sprintf("Failed to unmarshal the response. Error: %s, \nResponse: %s", e.Err, e.Body)
}
