package ekycclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"salon_be/component/logger"

	"go.uber.org/zap"
)

type EKYCConfig struct {
	BaseURL     string
	TokenID     string
	TokenKey    string
	AccessToken string
}

type EKYCClient struct {
	config EKYCConfig
	client *http.Client
}

type UploadResponse struct {
	Message string       `json:"message"`
	Object  UploadObject `json:"object"`
}

type UploadObject struct {
	FileName     string `json:"fileName"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Hash         string `json:"hash"`
	FileType     string `json:"fileType"`
	UploadedDate string `json:"uploadedDate"`
	StorageType  string `json:"storageType"`
	TokenId      string `json:"tokenId"`
}

func NewEKYCClient() *EKYCClient {
	ekycConfig := EKYCConfig{
		BaseURL:     "https://api.idg.vnpt.vn",
		TokenID:     "your-token-id",
		TokenKey:    "your-token-key",
		AccessToken: "your-access-token",
	}

	return &EKYCClient{
		config: ekycConfig,
		client: &http.Client{},
	}
}

func (c *EKYCClient) makeRequest(ctx context.Context, method, endpoint string, body io.Reader, contentType string) (*http.Response, error) {
	url := c.config.BaseURL + endpoint
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		logger.AppLogger.Error(ctx, "failed to create request",
			zap.String("method", method),
			zap.String("endpoint", endpoint),
			zap.Error(err))
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Token-id", c.config.TokenID)
	req.Header.Set("Token-key", c.config.TokenKey)
	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
	req.Header.Set("mac-address", "TEST1")

	resp, err := c.client.Do(req)
	if err != nil {
		logger.AppLogger.Error(ctx, "failed to execute request",
			zap.String("method", method),
			zap.String("endpoint", endpoint),
			zap.Error(err))
		return nil, err
	}

	if resp.StatusCode >= 400 {
		logger.AppLogger.Error(ctx, "received error response",
			zap.String("method", method),
			zap.String("endpoint", endpoint),
			zap.Int("status_code", resp.StatusCode))
	}

	return resp, nil
}

func (c *EKYCClient) UploadFile(ctx context.Context, file *multipart.FileHeader) (*UploadResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file
	part, err := writer.CreateFormFile("file", file.Filename)
	if err != nil {
		logger.AppLogger.Error(ctx, "failed to create form file",
			zap.String("filename", file.Filename),
			zap.Error(err))
		return nil, err
	}

	fileContent, err := file.Open()
	if err != nil {
		logger.AppLogger.Error(ctx, "failed to open file",
			zap.String("filename", file.Filename),
			zap.Error(err))
		return nil, err
	}
	defer fileContent.Close()

	if _, err = io.Copy(part, fileContent); err != nil {
		logger.AppLogger.Error(ctx, "failed to copy file content",
			zap.String("filename", file.Filename),
			zap.Error(err))
		return nil, err
	}

	// Add other form fields
	writer.WriteField("title", "Document Upload")
	writer.WriteField("description", "eKYC Document")

	if err := writer.Close(); err != nil {
		logger.AppLogger.Error(ctx, "failed to close multipart writer", zap.Error(err))
		return nil, err
	}

	logger.AppLogger.Info(ctx, "uploading file to EKYC service",
		zap.String("filename", file.Filename))

	resp, err := c.makeRequest(ctx, "POST", "/file-service/v1/addFile", body, writer.FormDataContentType())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result UploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.AppLogger.Error(ctx, "failed to decode upload response", zap.Error(err))
		return nil, err
	}

	return &result, nil
}

type ClassifyResponse struct {
	Message string `json:"message"`
	Object  struct {
		Type int    `json:"type"`
		Name string `json:"name"`
	} `json:"object"`
}

func (c *EKYCClient) ClassifyDocument(ctx context.Context, hash, clientSession string) (*ClassifyResponse, error) {
	payload := map[string]string{
		"img_card":       hash,
		"client_session": clientSession,
		"token":          "verification",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		logger.AppLogger.Error(ctx, "failed to marshal classify request",
			zap.Error(err))
		return nil, err
	}

	resp, err := c.makeRequest(ctx, "POST", "/ai/v1/classify/id", bytes.NewBuffer(jsonData), "application/json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ClassifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.AppLogger.Error(ctx, "failed to decode classify response", zap.Error(err))
		return nil, err
	}

	return &result, nil
}

type LivenessResponse struct {
	Message string `json:"message"`
	Object  struct {
		LivenessStatus string `json:"liveness"`
		LivenessMsg    string `json:"liveness_msg"`
		FaceSwapping   bool   `json:"face_swapping"`
		FakeLiveness   bool   `json:"fake_liveness"`
		IsEyeOpen      string `json:"is_eye_open,omitempty"`
	} `json:"object"`
}

func (c *EKYCClient) ValidateDocument(ctx context.Context, hash, clientSession string) (*LivenessResponse, error) {
	payload := map[string]string{
		"img":            hash,
		"client_session": clientSession,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		logger.AppLogger.Error(ctx, "failed to marshal validate document request",
			zap.Error(err))
		return nil, err
	}

	resp, err := c.makeRequest(ctx, "POST", "/ai/v1/card/liveness", bytes.NewBuffer(jsonData), "application/json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result LivenessResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.AppLogger.Error(ctx, "failed to decode validate document response", zap.Error(err))
		return nil, err
	}

	return &result, nil
}

type DocumentInfo struct {
	Message string `json:"message"`
	Object  struct {
		Name           string          `json:"name"`
		CardType       string          `json:"card_type"`
		ID             string          `json:"id"`
		IDProbs        string          `json:"id_probs"`
		BirthDay       string          `json:"birth_day"`
		BirthDayProb   float64         `json:"birth_day_prob"`
		Nationality    string          `json:"nationality"`
		Gender         string          `json:"gender"`
		ValidDate      string          `json:"valid_date"`
		IssueDate      string          `json:"issue_date"`
		IssuePlace     string          `json:"issue_place"`
		OriginLocation string          `json:"origin_location"`
		RecentLocation string          `json:"recent_location"`
		PostCode       json.RawMessage `json:"post_code"`
		Tampering      json.RawMessage `json:"tampering"`
	} `json:"object"`
}

func (c *EKYCClient) ExtractDocumentInfo(ctx context.Context, frontHash, backHash string, clientSession string) (*DocumentInfo, error) {
	payload := map[string]interface{}{
		"img_front":         frontHash,
		"img_back":          backHash,
		"client_session":    clientSession,
		"type":              -1,
		"validate_postcode": true,
		"token":             "verification",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		logger.AppLogger.Error(ctx, "failed to marshal extract document info request",
			zap.Error(err))
		return nil, err
	}

	resp, err := c.makeRequest(ctx, "POST", "/ai/v1/ocr/id", bytes.NewBuffer(jsonData), "application/json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result DocumentInfo
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.AppLogger.Error(ctx, "failed to decode document info response", zap.Error(err))
		return nil, err
	}

	return &result, nil
}

type FaceVerification struct {
	Message string `json:"message"`
	Object  struct {
		Result string  `json:"result"`
		Msg    string  `json:"msg"`
		Prob   float64 `json:"prob"`
	} `json:"object"`
}

func (c *EKYCClient) VerifyFace(ctx context.Context, docHash, faceHash string, clientSession string) (*FaceVerification, error) {
	payload := map[string]string{
		"img_front":      docHash,
		"img_face":       faceHash,
		"client_session": clientSession,
		"token":          "verification",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		logger.AppLogger.Error(ctx, "failed to marshal verify face request",
			zap.Error(err))
		return nil, err
	}

	resp, err := c.makeRequest(ctx, "POST", "/ai/v1/face/compare", bytes.NewBuffer(jsonData), "application/json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result FaceVerification
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.AppLogger.Error(ctx, "failed to decode face verification response", zap.Error(err))
		return nil, err
	}

	return &result, nil
}
