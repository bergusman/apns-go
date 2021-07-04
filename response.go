package apns

import (
	"encoding/json"
	"io"
	"net/http"
)

const StatusSent = http.StatusOK

type Response struct {
	Id        string // header: apns-id
	Status    int    // header: :status
	Reason    string
	Timestamp int64
}

func ParseResponse(r *http.Response) (*Response, error) {
	defer r.Body.Close()
	res := &Response{
		Id:     r.Header.Get("apns-id"),
		Status: r.StatusCode,
	}
	if err := json.NewDecoder(r.Body).Decode(res); err != nil && err != io.EOF {
		return nil, err
	}
	return res, nil
}
