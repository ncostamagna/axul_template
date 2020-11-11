package templates

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/jinzhu/gorm"
)

type repo struct {
	db     *gorm.DB
	logger log.Logger
}

//NewRepo is a repositories handler
func NewRepo(db *gorm.DB, logger log.Logger) Repository {
	return &repo{
		db:     db,
		logger: log.With(logger, "repo", "sql"),
	}
}

func (repo *repo) Create(ctx context.Context, template *Template) error {

	logger := log.With(repo.logger, "method", "Create")

	result := repo.db.Create(&template)

	if result.Error != nil {
		_ = level.Error(logger).Log("err", result.Error)
		return result.Error
	}

	_ = logger.Log("RowAffected", result.RowsAffected)
	_ = logger.Log("ID", template.ID)

	return nil
}

func (repo *repo) GetAll(ctx context.Context, templates *[]Template) error {

	logger := log.With(repo.logger, "method", "GetAll")

	result := repo.db.Find(&templates)

	if result.Error != nil {
		_ = level.Error(logger).Log("err", result.Error)
		return result.Error
	}

	_ = logger.Log("RowAffected", result.RowsAffected)

	return nil
}

func (repo *repo) Get(ctx context.Context, template *Template, id uint) error {

	logger := log.With(repo.logger, "method", "GetAll")

	//result := repo.db.Find(&templates)
	result := repo.db.Find(&template, &Template{ID: id})
	if result.Error != nil {
		_ = level.Error(logger).Log("err", result.Error)
		return result.Error
	}

	_ = logger.Log("RowAffected", result.RowsAffected)

	return nil
}

func (repo *repo) Update(ctx context.Context, template *Template, templateValues Template) error {

	return nil
}

func (repo *repo) Delete(ctx context.Context, template *Template, id uint) error {

	return nil
}
