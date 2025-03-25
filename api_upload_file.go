package dify

import (
	"bytes"
	"context"
	"net/http"
)

type UploadFileRequest struct {
	User  string       `form:"user"`
	Files bytes.Buffer `form:"files"`
}

type UploadFileResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	Extension string `json:"extension"`
	MimeType  string `json:"mime_type"`
	CreatedAt int    `json:"created_at"`
	CreatedBy string `json:"created_by"`
}

/*
 * 上传文件
 */
func (api *API) UploadFile(ctx context.Context, req *UploadFileRequest) (resp *UploadFileResponse, err error) {

	httpReq, err := api.createBaseRequest(ctx, http.MethodPost, "/files/upload", req, Chat)
	if err != nil {
		return
	}
	err = api.c.sendJSONRequest(httpReq, &resp)
	return
}
