package dify

import (
	"context"
	"fmt"
	"net/http"
)

type UploadFileRequest struct {
	User     string `json:"user"`
	FilePath string `json:"filePath"`
	FileType string `json:"fileType"`
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

	httpReq, err := api.CreateFormFileRequest(ctx, http.MethodPost, "/v1/files/upload", req, Chat)
	if err != nil {
		return
	}
	err = api.c.sendJSONRequest(httpReq, &resp)
	if err != nil {
		fmt.Printf("Failed to read response: %v\n", err)
		return nil, err
	}
	return
}
