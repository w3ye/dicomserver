package service

import (
	"dicomserver/repositories"
	"image"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestFileService_GetDicomHeaders(t *testing.T) {
	type fields struct {
		repo    fileServiceRepo
		Encoder interface {
			Encode(w io.Writer, m image.Image) error
		}
	}
	type args struct {
		id    string
		query string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *GetDicomHeaderAttributeResponse
		wantErr bool
	}{
		{
			name: "should return the correct dicom header attribute when a file exists",
			fields: fields{
				repo: repositories.NewLocalFileRepository("../test_file"),
			},
			args: args{
				id:    "07fbab94-a2b6-45fe-84b6-515bb65354e1",
				query: "(0002,0000)",
			},
			want: &GetDicomHeaderAttributeResponse{
				Tag:   "(0002,0000)",
				Name:  "FileMetaInformationGroupLength",
				VR:    "UL",
				Value: "[186]",
			},
			wantErr: false,
		},
		{
			name: "should return an error if the file does not exist",
			fields: fields{
				repo: repositories.NewLocalFileRepository("../test_file"),
			},
			args: args{
				id:    "1234",
				query: "(0002,0000)",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return an error if the tag is incorrect",
			fields: fields{
				repo: repositories.NewLocalFileRepository("../test_file"),
			},
			args: args{
				id:    "07fbab94-a2b6-45fe-84b6-515bb65354e1",
				query: "(0002,abcd)",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileService{
				repo:    tt.fields.repo,
				Encoder: tt.fields.Encoder,
			}
			got, err := f.GetDicomHeaders(tt.args.id, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileService.GetDicomHeaders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FileService.GetDicomHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileService_GetDicomImage(t *testing.T) {
	type fields struct {
		repo    fileServiceRepo
		Encoder interface {
			Encode(w io.Writer, m image.Image) error
		}
	}
	type args struct {
		id string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		emptyBytes bool
		wantErr    bool
	}{
		{
			name: "should not error when an existing file contains an image",
			fields: fields{
				repo:    repositories.NewLocalFileRepository("../test_file"),
				Encoder: PNGEncoder{},
			},
			args: args{
				id: "07fbab94-a2b6-45fe-84b6-515bb65354e1",
			},
			emptyBytes: false,
			wantErr:    false,
		},
		{
			name: "should return an error if there's no encoder",
			fields: fields{
				repo: repositories.NewLocalFileRepository("../test_file"),
			},
			args: args{
				id: "07fbab94-a2b6-45fe-84b6-515bb65354e1",
			},
			emptyBytes: true,
			wantErr:    true,
		},
		{
			name: "should return an error if the file path is incorrect",
			fields: fields{
				repo:    repositories.NewLocalFileRepository("../test_file2"),
				Encoder: PNGEncoder{},
			},
			args: args{
				id: "07fbab94-a2b6-45fe-84b6-515bb65354e1",
			},
			emptyBytes: true,
			wantErr:    true,
		},
		{
			name: "should return an error if the id incorrect",
			fields: fields{
				repo:    repositories.NewLocalFileRepository("../test_file"),
				Encoder: PNGEncoder{},
			},
			args: args{
				id: "07fbab94-a2b6-45fe-84b6-515bb65354e3",
			},
			emptyBytes: true,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileService{
				repo:    tt.fields.repo,
				Encoder: tt.fields.Encoder,
			}
			got, err := f.GetDicomImage(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileService.GetDicomImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.emptyBytes && len(got) != 0 || !tt.emptyBytes && len(got) == 0 {
				t.Errorf("FileService.GetDicomImage() = %v, want %v", got, tt.emptyBytes)
			}
		})
	}
	if err := os.Remove("../test_file2"); err != nil {
		t.Error("Failed to remove test_file2")
	}
}
