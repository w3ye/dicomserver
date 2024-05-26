package service

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"os"

	"github.com/google/uuid"
	"github.com/suyashkumar/dicom"
	dicomTag "github.com/suyashkumar/dicom/pkg/tag"
)

type FileService struct {
	filePath string
	Encoder  interface {
		Encode(w io.Writer, m image.Image) error
	}
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

	// validate the file
	_, err = dicom.Parse(file, fileHeader.Size, nil)
	if err != nil {
		return "", err
	}
	// reset the file pointer to the beginning of the file
	file.Seek(0, io.SeekStart)

	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(file); err != nil {
		return "", err
	}

	id := uuid.New().String()

	if err := os.WriteFile(
		fmt.Sprintf("%s/%s", f.filePath, id),
		buf.Bytes(),
		0666,
	); err != nil {
		// check if the file exists, if there's no error, the file exists
		if _, err := os.Stat(fmt.Sprintf("%s/%s.dicom", f.filePath, id)); err == nil {
			// remove the file
			os.Remove(fmt.Sprintf("%s/%s.dicom", f.filePath, id))
		}
		return "", err
	}
	return id, nil
}

type GetDicomHeaderAttributeResponse struct {
	Tag   string `json:"tag"`
	Name  string `json:"name"`
	VR    string `json:"vr"`
	Value string `json:"value"`
}

func (f FileService) GetDicomHeaders(id string, query string) (*GetDicomHeaderAttributeResponse, error) {
	tagQuery := queryTag(query)
	tag, err := tagQuery.convertQueryTagToDicomTag()
	if err != nil {
		return nil, err
	}

	tagInfo, err := dicomTag.Find(tag)
	if err != nil {
		return nil, err
	}

	dataset, err := dicom.ParseFile(fmt.Sprintf("%s/%s", f.filePath, id), nil)
	if err != nil {
		return nil, err
	}

	element, err := dataset.FindElementByTag(tag)
	if err != nil {
		return nil, err
	}

	response := &GetDicomHeaderAttributeResponse{
		Tag:   tagInfo.Tag.String(),
		Name:  tagInfo.Name,
		VR:    tagInfo.VR,
		Value: element.Value.String(),
	}

	return response, nil
}

func (f *FileService) SetImageEncoder(fileType string) {
	switch fileType {
	case "png":
		f.Encoder = PNGEncoder{}
	case "jpg":
		f.Encoder = JPEGEncoder{}
	}
}

func (f FileService) GetDicomImage(id string) ([]byte, error) {
	if f.Encoder == nil {
		return nil, errors.New("image encoder not set")
	}
	dataset, err := dicom.ParseFile(fmt.Sprintf("%s/%s", f.filePath, id), nil)
	if err != nil {
		return nil, err
	}

	element, err := dataset.FindElementByTag(dicomTag.PixelData)
	if err != nil {
		return nil, err
	}

	pixelDataInfo := dicom.MustGetPixelDataInfo(element.Value)

	buf := bytes.NewBuffer(nil)
	for _, frame := range pixelDataInfo.Frames {
		img, err := frame.GetImage()
		if err != nil {
			return nil, err
		}
		f.Encoder.Encode(buf, img)
	}
	return buf.Bytes(), nil
}
