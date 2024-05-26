package service

import "testing"

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
