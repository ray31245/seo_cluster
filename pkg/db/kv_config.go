package db

import (
	"errors"
	"fmt"

	"github.com/ray31245/seo_cluster/pkg/db/model"

	"gorm.io/gorm"
)

type KVConfigDAO struct {
	db *gorm.DB
}

func (d *DB) NewKVConfigDAO() (*KVConfigDAO, error) {
	err := d.db.AutoMigrate(&model.KVConfig{})
	if err != nil {
		return nil, fmt.Errorf("NewKVConfigDAO: %w", err)
	}

	return &KVConfigDAO{db: d.db}, nil
}

// UpsertByKey Upsert a key value pair
func (d *KVConfigDAO) UpsertByKey(key string, value string) error {
	kv, err := d.GetByKey(key)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	kv.Value = value
	kv.Key = key

	err = d.db.Save(&kv).Error
	if err != nil {
		return err
	}

	return nil
}

// GetByKey gets a key value pair by key
func (d *KVConfigDAO) GetByKey(key string) (model.KVConfig, error) {
	var kv model.KVConfig

	err := d.db.Where("key = ?", key).First(&kv).Error
	if err != nil {
		return model.KVConfig{}, err
	}

	return kv, nil
}

// GetByKeyWithDefault gets a key value pair by key with default value
func (d *KVConfigDAO) GetByKeyWithDefault(key string, defaultValue string) (model.KVConfig, error) {
	kv, err := d.GetByKey(key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			kv.Key = key
			kv.Value = defaultValue

			return kv, nil
		}

		return model.KVConfig{}, err
	}

	return kv, nil
}

// ListKVConfigs lists all key value pairs
func (d *KVConfigDAO) ListKVConfigs() ([]model.KVConfig, error) {
	var kvs []model.KVConfig
	err := d.db.Find(&kvs).Error

	return kvs, err
}

// DeleteByKey deletes a key value pair by key
func (d *KVConfigDAO) DeleteByKey(key string) error {
	err := d.db.Where("key = ?", key).Delete(&model.KVConfig{}).Error
	if err != nil {
		return err
	}

	return nil
}