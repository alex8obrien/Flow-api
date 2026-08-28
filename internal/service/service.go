package service

import (
	"context"

	"github.com/alex8obrien/echo-v5-template/internal/domain"
	"github.com/alex8obrien/echo-v5-template/internal/telemetry"
)

type UserService struct {
	repo   UserRepository
	tracer telemetry.Tracer
}

func NewUserService(repo UserRepository, tracer telemetry.Tracer) *UserService {
	return &UserService{repo: repo, tracer: tracer}
}

func (srv *UserService) Register(ctx context.Context, name, email string) (*domain.User, error) {
	ctx, span := srv.tracer.Start(ctx, "user.register", telemetry.String("user.email_domain", email))
	defer span.End()
	span.AddEvent("user.register.request_received")

	if email == "" {
		span.AddEvent("user.register.validation_failed", telemetry.String("field", "email"))
		span.RecordError(domain.ErrInvalidEmail)
		span.SetStatus(telemetry.StatusError, domain.ErrInvalidEmail.Error())
		return nil, domain.ErrInvalidEmail
	}

	u := &domain.User{Name: name, Email: email}
	span.AddEvent("user.register.repository_insert_start")

	user, err := srv.repo.Insert(ctx, u)
	if err != nil {
		span.RecordError(err)
		span.AddEvent("user.register.repository_insert_failed")
		span.SetStatus(telemetry.StatusError, "repository insert failed")
		return nil, err
	}

	span.SetAttributes(telemetry.String("user.id", user.ID.String()))
	span.SetStatus(telemetry.StatusSuccess, "user registered")
	span.AddEvent("user.register.completed")
	return user, nil
}
