package service

import (
	"context"
	"errors"
	"net/url"

	"github.com/alex8obrien/flow/internal/domain"
	"github.com/alex8obrien/flow/internal/telemetry"
)

type LinkService struct {
	repo   LinkRepository
	tracer telemetry.Tracer
}

func NewLinkService(repo LinkRepository, tracer telemetry.Tracer) *LinkService {
	return &LinkService{repo: repo, tracer: tracer}
}

func (srv *LinkService) Health(ctx context.Context) domain.Readyz {
	ctx, span := srv.tracer.Start(ctx, "service.links.health")
	defer span.End()
	span.AddEvent("service.links.health.request_received")

	// Define responses, assumes a negative state
	r := domain.NewReadyz()
	s := telemetry.StatusUnset

	// Ping Database
	err := srv.repo.Ping(ctx)
	if err != nil {
		span.RecordError(err)
		span.AddEvent("service.links.health.ping_failed")
	} else {
		r.DBChecked()
	}

	if r.Overall() == domain.OK {
		s = telemetry.StatusSuccess
	} else {
		s = telemetry.StatusError
	}

	span.AddEvent("service.links.health.completed")
	span.SetStatus(s, "health check succeeded")
	return *r
}

func (srv *LinkService) Create(
	ctx context.Context,
	destinationURL, displayName, description string,
) (*domain.Link, error) {
	ctx, span := srv.tracer.Start(ctx, "service.links.create")
	defer span.End()
	span.AddEvent("service.links.create.request_received")

	// Build Link
	link, err := domain.NewLink(destinationURL, displayName, description)
	if err != nil {
		span.AddEvent("service.links.create.validation_failed")
		span.RecordError(err)
		span.SetStatus(telemetry.StatusError, err.Error())
		return nil, err
	}

	link, err = srv.repo.Insert(ctx, link)
	if err != nil {
		span.RecordError(err)
		span.AddEvent("service.links.create.repository_insert_failed")
		span.SetStatus(telemetry.StatusError, "repository insert failed")
		return nil, err
	}

	span.SetAttributes(telemetry.String("link.id", link.ID.String()))
	span.SetStatus(telemetry.StatusSuccess, "link created")
	span.AddEvent("service.links.create.completed")
	return link, nil
}

func (srv *LinkService) Redirect(ctx context.Context, shortCode string) (*url.URL, error) {
	ctx, span := srv.tracer.Start(ctx, "service.links.redirect")
	defer span.End()
	span.AddEvent("service.links.redirect.request_received")

	if len(shortCode) != domain.ShortCodeMaximumLength {
		span.SetStatus(telemetry.StatusError, "invalid id length")
		span.AddEvent("service.links.redirect.invalid_id")
		return nil, errors.New("invalid id")
	}

	link, err := srv.repo.FindByShortCode(ctx, shortCode)
	if err != nil {
		span.RecordError(err)
		span.AddEvent("service.links.redirect.repository_find_failed")
		span.SetStatus(telemetry.StatusError, "repository find failed")
		return nil, err
	}

	// Build URL
	url, err := url.Parse(link.DestinationURL)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(telemetry.String("link.id", link.ID.String()))
	span.SetStatus(telemetry.StatusSuccess, "link created")
	span.AddEvent("service.links.create.completed")
	return url, nil
}
