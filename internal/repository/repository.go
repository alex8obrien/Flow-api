package repository

import (
	"context"
	"uuid"

	"github.com/alex8obrien/echo-v5-template/internal/domain"
	"github.com/alex8obrien/echo-v5-template/internal/telemetry"
)

type UserRepository struct {
	db     DB
	tracer telemetry.Tracer
}

func NewUserRepository(db DB, tracer telemetry.Tracer) UserRepository {
	return UserRepository{db: db, tracer: tracer}
}

func (r UserRepository) Insert(ctx context.Context, user *domain.User) (*domain.User, error) {
	ctx, span := r.tracer.Start(ctx, "repository.user.insert",
		telemetry.String("db.system", "postgresql"),
		telemetry.String("db.operation", "INSERT"),
		telemetry.String("db.sql.table", "users"),
	)
	defer span.End()
	span.AddEvent("repository.user.insert.query_start")

	var id uuid.UUID

	err := r.db.QueryRow(ctx, "INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id", user.Name, user.Email).
		Scan(&id)
	if err != nil {
		span.RecordError(err)
		span.AddEvent("repository.user.insert.query_failed")
		span.SetStatus(telemetry.StatusError, "insert query failed")
		return nil, err
	}

	user.ID = id
	span.SetAttributes(telemetry.String("user.id", user.ID.String()))
	span.SetStatus(telemetry.StatusSuccess, "insert completed")
	span.AddEvent("repository.user.insert.query_completed")

	return user, nil
}
