package models

import "time"

type FileMetadata struct {
	FileName     string    `json:"file_name"`
	FilePath     string    `json:"file_path"`
	UploadTime   time.Time `json:"upload_time"`
	StorageClass string    `json:"storage_class"`
}
