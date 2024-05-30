package service

import (
	"context"

	"go-admin/internal/modules/sys/model"
	"go-admin/internal/modules/sys/repo"
	"go-admin/pkg/util"
)

// Logger management
type Logger struct {
	LoggerRepo *repo.Logger
}

// Query loggers from the data access object based on the provided parameters and options.
func (a *Logger) Query(ctx context.Context, params model.LoggerQueryParam) (*model.LoggerQueryResult, error) {
	params.Pagination = true

	result, err := a.LoggerRepo.Query(ctx, params, model.LoggerQueryOptions{
		QueryOptions: util.QueryOptions{
			OrderFields: []util.OrderByParam{
				{Field: "created_at", Direction: util.DESC},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
