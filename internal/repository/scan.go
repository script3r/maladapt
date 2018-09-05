package repository

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
	"github.com/worlvlhole/maladapt/internal/model"
	"github.com/worlvlhole/maladapt/pkg/storage/mongo"
)

const scanCollection = "scan"

type ScanRepository interface {
	FindByID(id string) (*model.ScanResult, error)
	FindBySHA256(hash [sha256.Size]byte) (*model.ScanResult, error)
	FindByMD5(hash [md5.Size]byte) (*model.ScanResult, error)
	Store(a *model.ScanResult) error
	Delete(id string) error
	DeleteWithMD5(hash [md5.Size]byte) error
	DeleteWithSHA256(hash [sha256.Size]byte) error
	ApplyIndices() error
}

type ScanMongoRepository struct {
	mongoc *mongo.Client
}

func NewScanMongoRepository(config mongo.Configuration) ScanRepository {
	return &ScanMongoRepository{
		mongo.NewClient(config),
	}
}

func (s *ScanMongoRepository) FindByID(id string) (*model.ScanResult, error) {
	return &model.ScanResult{}, nil
}

func (s *ScanMongoRepository) FindBySHA256(hash [sha256.Size]byte) (*model.ScanResult, error) {
	logger := log.WithFields(log.Fields{"func": "FindBySHA256"})

	c, ses, err := s.Collection()
	if err != nil {
		logger.Error(err)
		return nil, nil
	}

	defer ses.Close()

	res := model.ScanResult{}
	err = c.Find(&bson.M{"sha256": hash}).One(&res)
	if err != nil {
		if err == mgo.ErrNotFound {
			logger.Info("No item found")
			return nil, nil
		}
		return nil, err
	}

	return &res, nil
}

func (s *ScanMongoRepository) FindByMD5(hash [md5.Size]byte) (*model.ScanResult, error) {
	return &model.ScanResult{}, nil
}

func (s *ScanMongoRepository) Store(m *model.ScanResult) error {
	logger := log.WithFields(log.Fields{"func": "Store"})

	col, ses, err := s.Collection()
	if err != nil {
		logger.Error(err)
		return err
	}

	defer ses.Close()

	if !m.ID.Valid() {
		m.ID = bson.NewObjectId()
	}

	_, err = col.UpsertId(m.ID, m)
	return err
}

func (s *ScanMongoRepository) DeleteWithMD5(hash [md5.Size]byte) error {
	logger := log.WithFields(log.Fields{"func": "DeleteWithMD5"})

	col, ses, err := s.Collection()
	if err != nil {
		logger.Error(err)
		return err
	}

	defer ses.Close()

	if err := col.Remove(&bson.M{"md5": hash}); err != nil {
		if err == mgo.ErrNotFound {
			logger.WithFields(log.Fields{"md5": base64.RawURLEncoding.EncodeToString(hash[:])}).Info("md5 hash not found")
			return nil
		}
		logger.Error(err)
		return err
	}
	return nil
}

func (s *ScanMongoRepository) Delete(id string) error {
	logger := log.WithFields(log.Fields{"func": "Delete"})
	if !bson.IsObjectIdHex(id) {
		return errors.New("invalid hex bson ID")
	}

	col, ses, err := s.Collection()
	if err != nil {
		logger.Error(err)
		return err
	}

	defer ses.Close()

	return col.RemoveId(bson.ObjectIdHex(id))
}

func (s *ScanMongoRepository) DeleteWithSHA256(hash [sha256.Size]byte) error {
	logger := log.WithFields(log.Fields{"func": "DeleteWithSHA256"})

	col, ses, err := s.Collection()
	if err != nil {
		logger.Error(err)
		return err
	}

	defer ses.Close()

	if err := col.Remove(&bson.M{"sha256": hash}); err != nil {
		if err == mgo.ErrNotFound {
			logger.WithFields(log.Fields{"sha256": base64.RawURLEncoding.EncodeToString(hash[:])}).Info("sha256 hash not found")
			return nil
		}
		logger.Error(err)
		return err
	}
	return nil
}
func (s *ScanMongoRepository) ApplyIndices() error {
	return nil
}

func (s *ScanMongoRepository) Collection() (*mgo.Collection, *mgo.Session, error) {
	logger := log.WithFields(log.Fields{"func": "Collection"})

	ses, db, err := s.mongoc.Database()
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	return db.C(scanCollection), ses, nil
}
