package service

import (
	"reflect"
	"testing"

	"github.com/suyashkumar/dicom/pkg/tag"
)

func Test_queryTag_extractGroupAndElement(t *testing.T) {
	tests := []struct {
		name    string
		q       queryTag
		want    uint16
		want1   uint16
		wantErr bool
	}{
		{
			name:    "valid tag",
			q:       queryTag("(0008,0005)"),
			want:    8,
			want1:   5,
			wantErr: false,
		},
		{
			name:    "invalid tag length",
			q:       queryTag("(0008,05)"),
			want:    0,
			want1:   0,
			wantErr: true,
		},
		{
			name:    "not hex decimals",
			q:       queryTag("(abcd,efgh)"),
			want:    0,
			want1:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.q.extractGroupAndElement()
			if (err != nil) != tt.wantErr {
				t.Errorf("queryTag.extractGroupAndElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("queryTag.extractGroupAndElement() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("queryTag.extractGroupAndElement() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_queryTag_validateFormat(t *testing.T) {
	tests := []struct {
		name    string
		q       queryTag
		wantErr bool
	}{
		{
			name:    "valid tag",
			q:       queryTag("(0008,0005)"),
			wantErr: false,
		},
		{
			name:    "invalid tag length",
			q:       queryTag("(0008,05)"),
			wantErr: true,
		},
		{
			name:    "not wrapped by paranthesis",
			q:       queryTag("0008,0005"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.q.validateFormat(); (err != nil) != tt.wantErr {
				t.Errorf("queryTag.validateFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_queryTag_convertQueryTagToDicomTag(t *testing.T) {
	tests := []struct {
		name    string
		q       queryTag
		want    tag.Tag
		wantErr bool
	}{
		{
			name:    "valid tag",
			q:       queryTag("(0008,0005)"),
			want:    tag.Tag{Group: 8, Element: 5},
			wantErr: false,
		},
		{
			name:    "invalid tag length",
			q:       queryTag("(0008,05)"),
			want:    tag.Tag{},
			wantErr: true,
		},
		{
			name:    "not hex decimals",
			q:       queryTag("(abcd,efgh)"),
			want:    tag.Tag{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.q.convertQueryTagToDicomTag()
			if (err != nil) != tt.wantErr {
				t.Errorf("queryTag.convertQueryTagToDicomTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queryTag.convertQueryTagToDicomTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
