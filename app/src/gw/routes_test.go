package main

import (
	"bytes"
	"net/http"
	"os"
	"testing"

	"github.com/karlpokus/routest/v2"
	"gw/internal/remote"
)

func TestRoutes(t *testing.T) {
	os.Setenv("ELASTIC_APM_ACTIVE", "false")
	routest.Test(t, func() http.Handler {
		c := &remote.Mock{
			remote.Response{
				"/user": &http.Response{
					Body:       &remote.Body{},
					StatusCode: 200,
				},
				"/posts": &http.Response{
					Body:       remote.NewBody("posts"),
					StatusCode: 200,
				},
				"/post/post1": &http.Response{
					Body:       remote.NewBody("post1"),
					StatusCode: 200,
				},
				"/post": &http.Response{
					StatusCode: 200,
				},
			},
		}
		return Routes(c, "nologs")
	}, []routest.Data{
		{
			Name:         "ping",
			Method:       "GET",
			Path:         "/ping",
			Status:       200,
			ResponseBody: []byte("pong"),
		},
		{
			Name:         "get posts no basic auth",
			Method:       "GET",
			Path:         "/posts",
			Status:       400,
			ResponseBody: []byte("missing auth"),
		},
		{
			Name:   "get posts",
			Method: "GET",
			Path:   "/posts",
			RequestHeader: http.Header{
				"Authorization": []string{"Basic Yml4YTpmbHVmZgo="},
			},
			Status:       200,
			ResponseBody: []byte("posts"),
		},
		{
			Name:   "get post",
			Method: "GET",
			Path:   "/post/post1",
			RequestHeader: http.Header{
				"Authorization": []string{"Basic Yml4YTpmbHVmZgo="},
			},
			Status:       200,
			ResponseBody: []byte("post1"),
		},
		{
			Name:        "add post",
			Method:      "POST",
			Path:        "/post",
			RequestBody: bytes.NewBuffer([]byte("data")),
			RequestHeader: http.Header{
				"Authorization": []string{"Basic Yml4YTpmbHVmZgo="},
			},
			Status:       200,
			ResponseBody: []byte("data"),
		},
	})
}
