package templates

import (
	"context"
	"time"

	"github.com/ncostamagna/streetflow/slack"

	"github.com/ncostamagna/rerrors"

	"github.com/go-kit/kit/log"
)

type service struct {
	repo      Repository
	slackTran slack.SlackBuilder
	logger    log.Logger
}

type updateCb func(uint, time.Time) error

//NewService is a service handler
func NewService(repo Repository, logger log.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

//Create service
func (s service) Create(ctx context.Context, template *Template) rerrors.RestErr {

	err := s.repo.Create(ctx, template)

	if err != nil {
		return rerrors.NewInternalServerError(err)
	}

	return nil
}

func (s service) Update(ctx context.Context) (*Template, rerrors.RestErr) {

	template := Template{}

	return &template, nil
}

func (s service) Delete(ctx context.Context) (*Template, rerrors.RestErr) {

	template := Template{}

	return &template, nil
}

func (s service) Get(ctx context.Context, template *Template, id uint) rerrors.RestErr {

	if err := s.repo.Get(ctx, template, id); err != nil {
		return rerrors.NewNotFoundError(err)
	}

	return nil
}

func (s service) GetAll(ctx context.Context, templates *[]Template) rerrors.RestErr {

	if err := s.repo.GetAll(ctx, templates); err != nil {
		return rerrors.NewInternalServerError(err)
	}

	return nil
}
