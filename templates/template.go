package templates

import (
	"context"
	"time"

	"github.com/ncostamagna/rerrors"
	"gorm.io/gorm"
)

//Template model
type Template struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Type      string         `gorm:"size:20;unique;not null" json:"type"`
	Template  string         `gorm:"type:text;not null" json:"template"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

//Repository is a Repository handler interface
type Repository interface {
	Create(ctx context.Context, template *Template) error
	Update(ctx context.Context, template *Template, templateValues Template) error
	GetAll(ctx context.Context, templates *[]Template) error
	Get(ctx context.Context, template *Template, id uint) error
	Delete(ctx context.Context, template *Template, id uint) error
}

//Service interface
type Service interface {
	Create(ctx context.Context, template *Template) rerrors.RestErr
	Update(ctx context.Context) (*Template, rerrors.RestErr)
	Get(ctx context.Context) (Template, rerrors.RestErr)
	GetAll(ctx context.Context, templates *[]Template) rerrors.RestErr
}
