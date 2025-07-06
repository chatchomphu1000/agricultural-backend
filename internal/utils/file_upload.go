package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// AllowedImageTypes defines the allowed MIME types for image uploads
var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// MaxFileSize defines the maximum file size for uploads (5MB)
const MaxFileSize = 5 * 1024 * 1024

// UploadConfig contains configuration for file uploads
type UploadConfig struct {
	UploadDir    string
	MaxFileSize  int64
	AllowedTypes map[string]bool
}

// FileUploadResult contains information about an uploaded file
type FileUploadResult struct {
	ID       string
	Filename string
	FilePath string
	FileSize int64
	MimeType string
}

// NewUploadConfig creates a new upload configuration
func NewUploadConfig() *UploadConfig {
	return &UploadConfig{
		UploadDir:    "uploads/products",
		MaxFileSize:  MaxFileSize,
		AllowedTypes: AllowedImageTypes,
	}
}

// EnsureUploadDir creates the upload directory if it doesn't exist
func (uc *UploadConfig) EnsureUploadDir() error {
	return os.MkdirAll(uc.UploadDir, 0755)
}

// ValidateFile validates the uploaded file
func (uc *UploadConfig) ValidateFile(header *multipart.FileHeader) error {
	// Check file size
	if header.Size > uc.MaxFileSize {
		return fmt.Errorf("file size %d exceeds maximum allowed size %d", header.Size, uc.MaxFileSize)
	}

	// Check MIME type
	if !uc.AllowedTypes[header.Header.Get("Content-Type")] {
		return fmt.Errorf("file type %s is not allowed", header.Header.Get("Content-Type"))
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		return fmt.Errorf("file extension %s is not allowed", ext)
	}

	return nil
}

// SaveFile saves the uploaded file to disk
func (uc *UploadConfig) SaveFile(header *multipart.FileHeader) (*FileUploadResult, error) {
	// Validate file first
	if err := uc.ValidateFile(header); err != nil {
		return nil, err
	}

	// Ensure upload directory exists
	if err := uc.EnsureUploadDir(); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Generate unique filename
	fileID := uuid.New().String()
	ext := filepath.Ext(header.Filename)
	newFilename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), fileID, ext)
	filePath := filepath.Join(uc.UploadDir, newFilename)

	// Open uploaded file
	src, err := header.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	_, err = io.Copy(dst, src)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	return &FileUploadResult{
		ID:       fileID,
		Filename: header.Filename,
		FilePath: filePath,
		FileSize: header.Size,
		MimeType: header.Header.Get("Content-Type"),
	}, nil
}

// DeleteFile removes a file from disk
func (uc *UploadConfig) DeleteFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	// Only delete files within our upload directory for security
	if !strings.HasPrefix(filePath, uc.UploadDir) {
		return fmt.Errorf("file path is outside upload directory")
	}

	return os.Remove(filePath)
}

// GenerateImageURL generates a URL for serving the uploaded image
func GenerateImageURL(filePath string, baseURL string) string {
	if filePath == "" {
		return ""
	}

	// Convert file path to URL path
	urlPath := strings.ReplaceAll(filePath, "\\", "/")
	return fmt.Sprintf("%s/%s", strings.TrimSuffix(baseURL, "/"), strings.TrimPrefix(urlPath, "/"))
}
