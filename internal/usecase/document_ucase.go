package usecase

import "github.com/V1taly5/APIYDisk/internal/infrastructure/repository"

type DocumentUseCase struct {
	repository repository.YandexDiskAPI
}

func NewDocumentUseCase(disk repository.YandexDiskAPI) *DocumentUseCase {
	return &DocumentUseCase{disk}
}

func (doc *DocumentUseCase) UploadDocument(imgUrl, uploadPath string) (map[string]interface{}, error) {
	document, err := doc.repository.UploadFileLink(imgUrl, uploadPath)
	if err != nil {
		return nil, err
	}
	return document, nil
}
