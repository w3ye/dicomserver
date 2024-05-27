package repositories

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/suyashkumar/dicom"
)

type LocalFileRepository struct {
	filePath string
}

func NewLocalFileRepository(filePath string) *LocalFileRepository {
	// if filePath does not exist, create it
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.Mkdir(filePath, 0777)
	}
	return &LocalFileRepository{
		filePath: filePath,
	}
}

func (r *LocalFileRepository) Write(bytes []byte) (string, error) {
	id := uuid.New().String()
	if err := os.WriteFile(
		fmt.Sprintf("%s/%s", r.filePath, id),
		bytes,
		0666,
	); err != nil {
		// check if the file exists, if there's no error, the file exists
		if _, err := os.Stat(fmt.Sprintf("%s/%s.dicom", r.filePath, id)); err == nil {
			// remove the file
			os.Remove(fmt.Sprintf("%s/%s.dicom", r.filePath, id))
		}
		return "", err
	}
	return id, nil
}

func (r *LocalFileRepository) GetDataset(id string) (dicom.Dataset, error) {
	dataset, err := dicom.ParseFile(fmt.Sprintf("%s/%s", r.filePath, id), nil)
	if err != nil {
		return dicom.Dataset{}, err
	}
	return dataset, nil
}
