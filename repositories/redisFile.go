package repositories

import (
	"bytes"
	"dicomserver/db"

	"github.com/google/uuid"
	"github.com/suyashkumar/dicom"
)

type RedisFileRepository struct {
	*db.Redis
}

func NewRedisFileRepository(db *db.Redis) *RedisFileRepository {
	return &RedisFileRepository{
		db,
	}
}

func (r *RedisFileRepository) Write(bytes []byte) (string, error) {
	id := uuid.New().String()
	err := r.Client.Set(r.Ctx, id, bytes, 0).Err()
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *RedisFileRepository) GetDataset(id string) (dicom.Dataset, error) {
	b, err := r.Client.Get(r.Ctx, id).Bytes()
	if err != nil {
		return dicom.Dataset{}, err
	}
	buf := bytes.NewBuffer(b)

	dataset, err := dicom.Parse(buf, int64(len(b)), nil)
	if err != nil {
		return dicom.Dataset{}, err
	}
	return dataset, nil
}
