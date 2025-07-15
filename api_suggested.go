package dify

import (
	"context"
	"fmt"
	"net/http"
)

const (
	SuggestedResponseResultSuccess = "success"
)

type SuggestedRequest struct {
	MessageId string `json:"message_id"`
	User      string `json:"user"`
}

type SuggestedResponse struct {
	Result string   `json:"result"`
	Data   []string `json:"data"`
}

// 相似问题推荐
func (api *API) MessageSuggested(ctx context.Context, req *SuggestedRequest) (resp *SuggestedResponse, err error) {
	httpReq, err := api.createBaseRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/messages/%s/suggested?user=%s", req.MessageId, req.User), nil, Chat)
	if err != nil {
		return
	}
	err = api.c.sendJSONRequest(httpReq, &resp)
	return
}
