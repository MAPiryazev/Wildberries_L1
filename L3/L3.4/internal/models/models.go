package models

type ProcessingType string

const (
	ProcessingResize    ProcessingType = "resize"
	ProcessingThumbnail ProcessingType = "thumbnail"
	ProcessingWatermark ProcessingType = "watermark"
)

type ImageTaskMessage struct {
	ImageID         string         `json:"image_id"`
	OriginalBucket  string         `json:"original_bucket"`
	OriginalPath    string         `json:"original_path"`
	ProcessedBucket string         `json:"processed_bucket"`
	ProcessedPath   string         `json:"processed_path"`
	ProcessingType  ProcessingType `json:"processing_type"`
}

type ImageStatus string

const (
	StatusPending    ImageStatus = "pending"
	StatusProcessing ImageStatus = "processing"
	StatusReady      ImageStatus = "ready"
	StatusFailed     ImageStatus = "failed"
)

type ImageInfo struct {
	ImageID      string      `json:"image_id"`
	OriginalURL  string      `json:"original_url"`
	ProcessedURL string      `json:"processed_url,omitempty"`
	Status       ImageStatus `json:"status"`
	ErrorMessage string      `json:"error_message,omitempty"`
}
