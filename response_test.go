package apns

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestParseResponse200(t *testing.T) {
	r := &http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Apns-Id": []string{"EC1BF194-B3B2-424A-89A9-5A918A6E6B5D"},
		},
		Body: io.NopCloser(strings.NewReader("")),
	}
	res, err := ParseResponse(r)
	if err != nil {
		t.Errorf("err not nil: %q", err)
	}
	if res != nil {
		if res.Id != "EC1BF194-B3B2-424A-89A9-5A918A6E6B5D" {
			t.Errorf("res.Id: %v want: EC1BF194-B3B2-424A-89A9-5A918A6E6B5D", res.Id)
		}
		if res.Status != Status200 {
			t.Errorf("res.Status: %v want: 200", res.Status)
		}
	} else {
		t.Error("res is nil")
	}
}

func TestParseResponseNot200(t *testing.T) {
	r := &http.Response{
		StatusCode: 410,
		Header: http.Header{
			"Apns-Id": []string{"EC1BF194-B3B2-424A-89A9-5A918A6E6B5D"},
		},
		Body: io.NopCloser(strings.NewReader(`{"reason":"Unregistered","timestamp":1629000000000}`)),
	}
	res, err := ParseResponse(r)
	if err != nil {
		t.Errorf("err not nil: %q", err)
	}
	if res != nil {
		if res.Id != "EC1BF194-B3B2-424A-89A9-5A918A6E6B5D" {
			t.Errorf("res.Id: %v want: EC1BF194-B3B2-424A-89A9-5A918A6E6B5D", res.Id)
		}
		if res.Status != 410 {
			t.Errorf("res.Status: %v want: 410", res.Status)
		}
		if res.Timestamp != 1629000000000 {
			t.Errorf("res.Timestamp: %v want: 1629000000000", res.Timestamp)
		}
	} else {
		t.Error("res is nil")
	}
}

func TestParseResponseErrors(t *testing.T) {
	r := &http.Response{
		StatusCode: 410,
		Header: http.Header{
			"Apns-Id": []string{"EC1BF194-B3B2-424A-89A9-5A918A6E6B5D"},
		},
		Body: io.NopCloser(strings.NewReader("Hello")),
	}
	_, err := ParseResponse(r)
	if err == nil {
		t.Error("err must be not nil")
	}
}
