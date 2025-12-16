package service

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/saleh-ghazimoradi/Cartopher/pkg/uploadProvider"
)

type UploadService interface {
	UploadProductImage(productId uint, file *multipart.FileHeader) (string, error)
}

type uploadService struct {
	provider uploadProvider.UploadProvider
}

func (u *uploadService) UploadProductImage(productId uint, file *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isValidImageExt(ext) {
		return "", fmt.Errorf("invalid file type: %s", ext)
	}

	newFileName := uuid.New().String()

	path := fmt.Sprintf("products/%d/%s%s", productId, newFileName, ext)

	return u.provider.UploadFile(file, path)
}

func isValidImageExt(ext string) bool {
	validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}

func NewUploadService(provider uploadProvider.UploadProvider) UploadService {
	return &uploadService{
		provider: provider,
	}
}
