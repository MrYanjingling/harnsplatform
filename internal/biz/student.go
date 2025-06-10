package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// var (
// 	// ErrUserNotFound is user not found.
// 	ErrStudentNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
// )

// Greeter is a Greeter model.
type Student struct {
	ID      string
	Name    string
	Sayname string
}

// GreeterRepo is a Greater repo.
type StudentRepo interface {
	Save(context.Context, *Student) (*Student, error)
	Update(context.Context, *Student) (*Student, error)
	FindByID(context.Context, int64) (*Student, error)
	ListByHello(context.Context, string) ([]*Student, error)
	ListAll(context.Context) ([]*Student, error)
}

// GreeterUsecase is a Greeter usecase.
type StudentUsecase struct {
	repo StudentRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewStudentUsecase(repo StudentRepo, logger log.Logger) *StudentUsecase {
	return &StudentUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *StudentUsecase) CreateStudent(ctx context.Context, g *Student) (*Student, error) {
	uc.log.WithContext(ctx).Infof("CreateStudent: %v", g.Name)
	return uc.repo.Save(ctx, g)
}
