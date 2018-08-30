package repository

import (
	"crypto/sha256"
	"github.com/worlvlhole/maladapt/internal/model"
	"github.com/worlvlhole/maladapt/pkg/storage/mongo"
)

type ScanRepository interface {
	FindByID(id string) (*model.ScanResult, error)
	FindBySHA256(hash [sha256.Size]byte) (*model.ScanResult, error)
	FindByMD5(hash string) (*model.ScanResult, error)
	Store(a *model.ScanResult) error
	Delete(id string) error
	DeleteWithMD5(hash string) error
	DeleteWithSHA256(hash string) error
	ApplyIndices() error
}

type ScanMongoRepository struct {
	mongoc *mongo.MongoClient
}

func NewScanMongoRepository(config mongo.Configuration) ScanRepository {
	return &ScanMongoRepository{
		mongo.NewMongoClient(config),
	}
}

func (s *ScanMongoRepository) FindByID(id string) (*model.ScanResult, error) {
	return &model.ScanResult{}, nil
}

func (s *ScanMongoRepository) FindBySHA256(hash [sha256.Size]byte) (*model.ScanResult, error) {
	return &model.ScanResult{}, nil
}

func (s *ScanMongoRepository) FindByMD5(hash string) (*model.ScanResult, error) {
	return &model.ScanResult{}, nil
}

func (s *ScanMongoRepository) Store(m *model.ScanResult) error {
	return nil
}

func (s *ScanMongoRepository) DeleteWithMD5(hash string) error {
	return nil
}

func (s *ScanMongoRepository) Delete(id string) error {
	return nil
}
func (s *ScanMongoRepository) DeleteWithSHA256(hash string) error {
	return nil
}
func (s *ScanMongoRepository) ApplyIndices() error {
	return nil
}
