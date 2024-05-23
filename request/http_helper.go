package request

import (
	"context"
	"encoding/json"

	"github.com/go-resty/resty/v2"

	"github.com/night-sword/kratos-kit/errors"
)

func DecodeJsonResponse[T any](jsonRsp []byte) (rsp *T, err error) {
	err = json.Unmarshal(jsonRsp, &rsp)
	if err != nil {
		err = errors.BadRequest(errors.RsnAccessRepoFail, "decode response fail").AddMetadata("rsp", string(jsonRsp))
	}
	return
}

func DecodeKratosJsonResponse[T any](response *resty.Response) (rsp *T, err error) {
	if response.StatusCode() != 200 {
		err = errors.FromHttpRsp(response.Body())
		return
	}

	err = json.Unmarshal(response.Body(), &rsp)
	if err != nil {
		err = errors.BadRequest(errors.RsnAccessRepoFail, "decode response fail").AddMetadata("rsp", string(response.Body()))
	}

	return
}

func HttpGet[Response any](ctx context.Context, client *resty.Client, api string) (rsp *Response, err error) {
	r, err := client.R().
		SetContext(ctx).
		Get(api)
	if err != nil {
		return
	}

	return DecodeKratosJsonResponse[Response](r)
}

func HttpPost[Request, Response any](ctx context.Context, client *resty.Client, api string, req Request) (rsp *Response, err error) {
	r, err := client.R().
		SetContext(ctx).
		SetBody(req).
		Post(api)
	if err != nil {
		return
	}

	return DecodeKratosJsonResponse[Response](r)
}
