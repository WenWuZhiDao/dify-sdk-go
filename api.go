package dify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const (
	Chat    = "Chat"
	Dataset = "Dataset"
)

type API struct {
	c          *Client
	chatSecret string
	dataSecret string
}

func (api *API) WithChatSecret(secret string) *API {
	api.chatSecret = secret
	return api
}

func (api *API) getChatSecret() string {
	if api.chatSecret != "" {
		return api.chatSecret
	}
	return api.c.getChatAPISecret()
}

func (api *API) WithDatasetSecret(secret string) *API {
	api.dataSecret = secret
	return api
}

func (api *API) getDatasetSecret() string {
	if api.dataSecret != "" {
		return api.dataSecret
	}
	return api.c.getDatasetAPISecret()
}

func (api *API) createBaseRequest(ctx context.Context, method, apiUrl string, body interface{}, apiType string) (*http.Request, error) {
	var b io.Reader
	if body != nil {
		// 断言body是否是io.Reader类型
		if reader, ok := body.(io.Reader); ok {
			b = reader
		} else {
			reqBytes, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			b = bytes.NewBuffer(reqBytes)
		}
	} else {
		b = http.NoBody
	}
	req, err := http.NewRequestWithContext(ctx, method, api.c.getHost()+apiUrl, b)
	if err != nil {
		return nil, err
	}
	var token string
	switch apiType {
	case Chat:
		token = api.getChatSecret()
	case Dataset:
		token = api.getDatasetSecret()
	default:
		token = api.getChatSecret()
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	return req, nil
}

func (api *API) CreateFormFileRequest(ctx context.Context, method, apiUrl string, params *UploadFileRequest, apiType string) (*http.Request, error) {
	// 打开文件
	file, err := os.Open(params.FilePath)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		return nil, err
	}
	defer file.Close()

	// 创建表单数据
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文件字段
	part, err := writer.CreateFormFile("file", filepath.Base(params.FilePath))
	if err != nil {
		fmt.Printf("Failed to create form file: %v\n", err)
		return nil, err
	}

	// 将文件内容写入表单
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Printf("Failed to copy file to form: %v\n", err)
		return nil, err
	}

	// 添加用户标识字段
	writer.WriteField("user", params.User)

	// 关闭 writer
	writer.Close()
	var token string
	switch apiType {
	case Chat:
		token = api.getChatSecret()
	case Dataset:
		token = api.getDatasetSecret()
	default:
		token = api.getChatSecret()
	}

	req, err := http.NewRequestWithContext(ctx, method, api.c.getHost()+apiUrl, body)
	if err != nil {
		return nil, err
	}
	// 设置请求头
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Cache-Control", "no-cache")

	return req, nil
}
