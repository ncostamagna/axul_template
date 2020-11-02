package contacts

import (
	"context"

	"github.com/ncostamagna/rerrors"
)

//Contact model
type Template struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Type     string `gorm:"size:20" json:"type"`
	Template string `gorm:"type:text" json:"template"`
}

//Repository is a Repository handler interface
type Repository interface {
	Create(ctx context.Context, template *Template) error
	Update(ctx context.Context, template *Template, templateValues Template) error
	GetAll(ctx context.Context, template *[]Template) error
	Get(ctx context.Context, template *Template, id uint) error
	Delete(ctx context.Context, template *Template, id uint) error
}

//Service interface
type Service interface {
	Create(ctx context.Context, template *Template) rerrors.RestErr
	Update(ctx context.Context) (*Template, rerrors.RestErr)
	Get(ctx context.Context) (Template, rerrors.RestErr)
}
