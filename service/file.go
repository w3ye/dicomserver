package service

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/google/uuid"
)

type FileService struct {
	filePath string
}

func NewFileService() *FileService {
	// establish a local file path
	filepath := "store"
	// if filePath does not exist, create it
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		os.Mkdir(filepath, 0777)
	}

	return &FileService{
		filePath: filepath,
	}
}

func (f FileService) WriteFile(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(file); err != nil {
		return "", err
	}

	id := uuid.New().String()

	filePath := "store"
	if err := os.WriteFile(
		fmt.Sprintf("%s/%s", filePath, id),
		buf.Bytes(),
		0666,
	); err != nil {
		// check if the file exists, if there's no error, the file exists
		if _, err := os.Stat(fmt.Sprintf("%s/%s.dicom", filePath, id)); err == nil {
			// remove the file
			os.Remove(fmt.Sprintf("%s/%s.dicom", filePath, id))
		}
	}
	return id, nil
}

func (f FileService) GetFileHeaders() {}

func (f FileService) GetImage() {}
