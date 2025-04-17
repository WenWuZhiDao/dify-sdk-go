package dify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

type DocumentsRequest struct {
	DatasetID string `json:"dataset_id"`
	Keyword   string `json:"keyword,omitempty"`
	Page      int    `json:"page,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

type DocumentsDataResponse struct {
	ID                   string `json:"id"`
	URL                  string `json:"url"`
	Position             int    `json:"position"`
	DataSourceType       string `json:"data_source_type"`
	DataSourceInfo       any    `json:"data_source_info"`
	DatasetProcessRuleId any    `json:"dataset_process_rule_id"`
	Name                 string `json:"name"`
	CreatedFrom          string `json:"created_from"`
	CreatedBy            string `json:"created_by"`
	CreatedAt            int    `json:"created_at"`
	Tokens               int    `json:"tokens"`
	IndexingStatus       string `json:"indexing_status"`
	Error                any    `json:"error"`
	Enabled              bool   `json:"enabled"`
	DisabledAt           any    `json:"disabled_at"`
	DisabledBy           any    `json:"disabled_by"`
	Archived             bool   `json:"archived"`
}

type DocumentsResponse struct {
	Limit   int                     `json:"limit"`
	HasMore bool                    `json:"has_more"`
	Total   int                     `json:"total"`
	Page    int                     `json:"page"`
	Data    []DocumentsDataResponse `json:"data"`
}

func (api *API) Documents(ctx context.Context, req *DocumentsRequest) (resp *DocumentsResponse, err error) {
	httpReq, err := api.createBaseRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/datasets/%s/documents", req.DatasetID), nil, Dataset)
	if err != nil {
		return
	}
	query := httpReq.URL.Query()
	query.Set("keyword", req.Keyword)
	if req.Page > 0 {
		query.Set("page", strconv.FormatInt(int64(req.Page), 10))
	}
	if req.Limit > 0 {
		query.Set("limit", strconv.FormatInt(int64(req.Limit), 10))
	}
	httpReq.URL.RawQuery = query.Encode()

	err = api.c.sendJSONRequest(httpReq, &resp)
	return
}

// ------------------------------

type ProcessRule struct {
	Rules any    `json:"rules,omitempty"`
	Mode  string `json:"mode"`
}

type DocumentCreateByTextRequest struct {
	DatasetID         string      `json:"dataset_id"`
	Name              string      `json:"name"`
	Text              string      `json:"text"`
	IndexingTechnique string      `json:"indexing_technique,omitempty"`
	ProcessRule       ProcessRule `json:"process_rule"`
}

type DocumentCreateByTextResponse struct {
	Document struct {
		ID             string `json:"id"`
		Position       int    `json:"position"`
		DataSourceType string `json:"data_source_type"`
		DataSourceInfo struct {
			UploadFileId string `json:"upload_file_id"`
		} `json:"data_source_info"`
		DatasetProcessRuleId string `json:"dataset_process_rule_id"`
		Name                 string `json:"name"`
		CreatedFrom          string `json:"created_from"`
		CreatedBy            string `json:"created_by"`
		CreatedAt            int    `json:"created_at"`
		Tokens               int    `json:"tokens"`
		IndexingStatus       string `json:"indexing_status"`
		Error                any    `json:"error"`
		Enabled              bool   `json:"enabled"`
		DisabledAt           any    `json:"disabled_at"`
		DisabledBy           any    `json:"disabled_by"`
		Archived             bool   `json:"archived"`
		DisplayStatus        string `json:"display_status"`
		WordCount            int    `json:"word_count"`
		HitCount             int    `json:"hit_count"`
		DocForm              string `json:"doc_form"`
	} `json:"document"`
	Batch string `json:"batch"`
}

func (api *API) DocumentCreateByText(ctx context.Context, req *DocumentCreateByTextRequest) (resp *DocumentCreateByTextResponse, err error) {
	httpReq, err := api.createBaseRequest(ctx, http.MethodPost, fmt.Sprintf("/v1/datasets/%s/document/create_by_text", req.DatasetID), req, Dataset)
	if err != nil {
		return
	}
	err = api.c.sendJSONRequest(httpReq, &resp)
	return
}

// ------------------------------

type DatasetSegmentsDataResponse struct {
	ID            string   `json:"id"`
	Position      int      `json:"position"`
	DocumentId    string   `json:"document_id"`
	Keywords      []string `json:"keywords"`
	Content       string   `json:"content"`
	Answer        string   `json:"answer"`
	WordCount     int      `json:"word_count"`
	Tokens        int      `json:"tokens"`
	IndexNodeId   string   `json:"index_node_id"`
	IndexNodeHash string   `json:"index_node_hash"`
	HitCount      int      `json:"hit_count"`
	Enabled       bool     `json:"enabled"`
	DisabledAt    any      `json:"disabled_at"`
	DisabledBy    any      `json:"disabled_by"`
	Status        string   `json:"status"`
	CreatedBy     string   `json:"created_by"`
	CreatedAt     int      `json:"created_at"`
	IndexingAt    int      `json:"indexing_at"`
	CompletedAt   int      `json:"completed_at"`
	Error         any      `json:"error"`
	StoppedAt     any      `json:"stopped_at"`
}

type GetSegmentsByDocumentIdResponse struct {
	Data    []DatasetSegmentsDataResponse `json:"data"`
	DocForm string                        `json:"doc_form"`
}

type GetSegmentsByDocumentIdRequest struct {
	DatasetID  string `json:"dataset_id"`
	DocumentId string `json:"document_id"`
	Keyword    string `json:"keyword"`
	Status     string `json:"status"`
}

func (api *API) GetSegmentsByDocumentId(ctx context.Context, req *GetSegmentsByDocumentIdRequest) (resp *GetSegmentsByDocumentIdResponse, err error) {
	httpReq, err := api.createBaseRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/datasets/%s/documents/%s/segments", req.DatasetID, req.DocumentId), req, Dataset)
	if err != nil {
		return
	}
	err = api.c.sendJSONRequest(httpReq, &resp)
	return
}

// ------------------------------
type AddDocumentSegmentsSegmentsRequest struct {
	Content  string   `json:"content"`
	Answer   string   `json:"answer"`
	Keywords []string `json:"keywords"`
}
type AddDocumentSegmentsRequest struct {
	DatasetID  string                               `json:"dataset_id"`
	DocumentId string                               `json:"document_id"`
	Segments   []AddDocumentSegmentsSegmentsRequest `json:"segments"`
}

type AddDocumentSegmentsResponse struct {
	Data    []DatasetSegmentsDataResponse `json:"data"`
	DocForm string                        `json:"doc_form"`
}

func (api *API) AddSegmentsToDocument(ctx context.Context, req *AddDocumentSegmentsRequest) (resp *AddDocumentSegmentsResponse, err error) {
	httpReq, err := api.createBaseRequest(ctx, http.MethodPost, fmt.Sprintf("/v1/datasets/%s/documents/%s/segments", req.DatasetID, req.DocumentId), req, Dataset)
	if err != nil {
		return
	}
	err = api.c.sendJSONRequest(httpReq, &resp)
	return
}

// ------------------------------

type DeleteSegmentsRequest struct {
	DatasetID  string `json:"dataset_id"`
	DocumentId string `json:"document_id"`
	SegmentId  string `json:"segment_id"`
}

type DeleteResponse struct {
	Result string `json:"result"`
}

func (api *API) DeleteDatasetSegments(ctx context.Context, req *DeleteSegmentsRequest) (resp *DeleteResponse, err error) {
	httpReq, err := api.createBaseRequest(ctx, http.MethodDelete, fmt.Sprintf("/v1/datasets/%s/documents/%s/segments/%s", req.DatasetID, req.DocumentId, req.SegmentId), nil, Dataset)
	if err != nil {
		return
	}
	err = api.c.sendJSONRequest(httpReq, &resp)
	return
}

// ------------------------------
type DeleteDocumentRequest struct {
	DatasetID  string `json:"dataset_id"`
	DocumentId string `json:"document_id"`
}

func (api *API) DeleteDatasetDocument(ctx context.Context, req *DeleteDocumentRequest) (resp *DeleteResponse, err error) {
	httpReq, err := api.createBaseRequest(ctx, http.MethodDelete, fmt.Sprintf("/v1/datasets/%s/documents/%s", req.DatasetID, req.DocumentId), nil, Dataset)
	if err != nil {
		return
	}
	err = api.c.sendJSONRequest(httpReq, &resp)
	return
}

// ------------------------------

type DocumentUpdateByTextRequest struct {
	DocumentCreateByTextRequest
	DocumentID string `json:"document_id"`
}

type DocumentUpdateByTextResponse struct {
	DocumentCreateByTextResponse
}

func (api *API) DocumentUpdateByText(ctx context.Context, req *DocumentUpdateByTextRequest) (resp *DocumentUpdateByTextResponse, err error) {
	httpReq, err := api.createBaseRequest(ctx, http.MethodPost, fmt.Sprintf("/v1/datasets/%s/documents/%s/update_by_text", req.DatasetID, req.DocumentID), req, Dataset)
	if err != nil {
		return
	}
	err = api.c.sendJSONRequest(httpReq, &resp)
	return
}

// ------------------------------

type DocumentCreateByFileRequest struct {
	DatasetID string         `json:"dataset_id"`
	File      multipart.File `json:"file"`
	FileName  string         `json:"file_name"`
	Data      struct {
		OriginalDocumentID string      `json:"original_document_id,omitempty"`
		IndexingTechnique  string      `json:"indexing_technique,omitempty"`
		ProcessRule        ProcessRule `json:"process_rule"`
	} `json:"data"`
}

type DocumentCreateByFileResponse struct {
	Document struct {
		ID             string `json:"id"`
		Position       int    `json:"position"`
		DataSourceType string `json:"data_source_type"`
		DataSourceInfo struct {
			UploadFileId string `json:"upload_file_id"`
		} `json:"data_source_info"`
		DatasetProcessRuleId string `json:"dataset_process_rule_id"`
		Name                 string `json:"name"`
		CreatedFrom          string `json:"created_from"`
		CreatedBy            string `json:"created_by"`
		CreatedAt            int    `json:"created_at"`
		Tokens               int    `json:"tokens"`
		IndexingStatus       string `json:"indexing_status"`
		Error                any    `json:"error"`
		Enabled              bool   `json:"enabled"`
		DisabledAt           any    `json:"disabled_at"`
		DisabledBy           any    `json:"disabled_by"`
		Archived             bool   `json:"archived"`
		DisplayStatus        string `json:"display_status"`
		WordCount            int    `json:"word_count"`
		HitCount             int    `json:"hit_count"`
		DocForm              string `json:"doc_form"`
	} `json:"document"`
	Batch string `json:"batch"`
}

func (api *API) DocumentCreateByFile(ctx context.Context, req *DocumentCreateByFileRequest) (resp *DocumentCreateByFileResponse, err error) {
	reqData, err := json.Marshal(req.Data)
	if err != nil {
		return nil, fmt.Errorf("error json.Marshal form data: %v", err)
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", req.FileName)
	if err != nil {
		return nil, fmt.Errorf("error creating form file: %v", err)
	}
	_, err = io.Copy(part, req.File)
	if err != nil {
		return nil, fmt.Errorf("error copying file: %v", err)
	}
	err = writer.WriteField("data", string(reqData))
	if err != nil {
		return nil, fmt.Errorf("error writer.WriteField: %v", err)
	}
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing writer: %v", err)
	}
	httpReq, err := api.createBaseRequest(ctx, http.MethodPost, fmt.Sprintf("/v1/datasets/%s/document/create_by_file", req.DatasetID), body, Dataset)
	if err != nil {
		return
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	err = api.c.sendJSONRequest(httpReq, &resp)
	return
}

// ------------------------------

type DocumentsIndexingStatusRequest struct {
	DatasetID string `json:"dataset_id"`
	Batch     string `json:"batch"`
}

type DocumentsIndexingStatusDataResponse struct {
	ID                   string  `json:"id"`
	IndexingStatus       string  `json:"indexing_status"`
	ProcessingStartedAt  float64 `json:"processing_started_at"`
	ParsingCompletedAt   float64 `json:"parsing_completed_at"`
	CleaningCompletedAt  float64 `json:"cleaning_completed_at"`
	SplittingCompletedAt float64 `json:"splitting_completed_at"`
	CompletedAt          any     `json:"completed_at"`
	PausedAt             any     `json:"paused_at"`
	Error                any     `json:"error"`
	StoppedAt            any     `json:"stopped_at"`
	CompletedSegments    int     `json:"completed_segments"`
	TotalSegments        int     `json:"total_segments"`
}

type DocumentsIndexingStatusResponse struct {
	Data []DocumentsIndexingStatusDataResponse `json:"data"`
}

func (api *API) DocumentsIndexingStatus(ctx context.Context, req *DocumentsIndexingStatusRequest) (resp *DocumentsIndexingStatusResponse, err error) {
	httpReq, err := api.createBaseRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/datasets/%s/documents/%s/indexing-status", req.DatasetID, req.Batch), nil, Dataset)
	if err != nil {
		return
	}
	err = api.c.sendJSONRequest(httpReq, &resp)
	return
}

// ------------------------------

type DocumentDeleteRequest struct {
	DatasetID  string `json:"dataset_id"`
	DocumentID string `json:"document_id"`
}

type DocumentDeleteResponse struct {
	Result string `json:"result"`
}

func (api *API) DocumentDelete(ctx context.Context, req *DocumentDeleteRequest) (resp *DocumentDeleteResponse, err error) {
	httpReq, err := api.createBaseRequest(ctx, http.MethodDelete, fmt.Sprintf("/v1/datasets/%s/documents/%s", req.DatasetID, req.DocumentID), nil, Dataset)
	if err != nil {
		return
	}
	err = api.c.sendJSONRequest(httpReq, &resp)
	return
}
