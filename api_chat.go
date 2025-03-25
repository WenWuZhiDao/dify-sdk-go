package dify

import (
	"context"
	"net/http"
)

type ChatMessageRequest struct {
	Inputs         map[string]interface{}   `json:"inputs"`
	Query          string                   `json:"query"`
	ResponseMode   string                   `json:"response_mode"`
	ConversationID string                   `json:"conversation_id,omitempty"`
	User           string                   `json:"user"`
	Files          []*ChatMessageFilesParam `json:"files"`
}

type ChatMessageResponse struct {
	ID             string `json:"id"`
	Answer         string `json:"answer"`
	ConversationID string `json:"conversation_id"`
	CreatedAt      int    `json:"created_at"`
}

/* Create chat message
 * Create a new conversation message or continue an existing dialogue.
 */
func (api *API) ChatMessages(ctx context.Context, req *ChatMessageRequest) (resp *ChatMessageResponse, err error) {
	req.ResponseMode = "blocking"

	httpReq, err := api.createBaseRequest(ctx, http.MethodPost, "/v1/chat-messages", req, Chat)
	if err != nil {
		return
	}
	err = api.c.sendJSONRequest(httpReq, &resp)
	return
}

type ChatMessageFilesParam struct {
	Type           string `json:"type"`
	TransferMethod string `json:"transfer_method"`
	Url            string `json:"url,omitempty"`
	UploadFileId   string `josn:"upload_file_id"`
}
