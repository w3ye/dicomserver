package service

import (
	"bytes"
	"errors"
	"image"
	"io"
	"mime/multipart"

	"github.com/nfnt/resize"
	"github.com/suyashkumar/dicom"
	dicomTag "github.com/suyashkumar/dicom/pkg/tag"
)

type fileServiceRepo interface {
	Write(bytes []byte) (string, error)
	GetDataset(id string) (dicom.Dataset, error)
}

type FileService struct {
	repo    fileServiceRepo
	Encoder interface {
		Encode(w io.Writer, m image.Image) error
	}
}

func NewFileService(repo fileServiceRepo) *FileService {
	return &FileService{
		repo: repo,
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

	id, err := f.repo.Write(buf.Bytes())
	if err != nil {
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

	dataset, err := f.repo.GetDataset(id)
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

	dataset, err := f.repo.GetDataset(id)
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
		// scale the image to the max size
		width := uint(img.Bounds().Max.Y)
		height := uint(img.Bounds().Max.Y)
		scaledImage := resize.Resize(width, height, img, resize.Lanczos3)
		f.Encoder.Encode(buf, scaledImage)
	}
	return buf.Bytes(), nil
}
