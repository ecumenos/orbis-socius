package accounts

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/ecumenos/orbis-socius/cmd/api/repo"
	"github.com/ecumenos/orbis-socius/models"
	"github.com/ecumenos/orbis-socius/pkg/toolkit/errorsutils"
	"github.com/jackc/pgx/v4"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewService),
)

type Service struct {
	db repo.Driver
}

func NewService(db repo.Driver) *Service {
	return &Service{db: db}
}

func (s *Service) createAccount(ctx context.Context, uniqueName, domain, displayName string, civitas int64) (*models.Account, error) {
	a, err := s.getAccountByUniqueName(ctx, uniqueName)
	if err != nil {
		return nil, fmt.Errorf("issues with getting account by unique name (%v), err = %w", uniqueName, err)
	}
	if a != nil {
		return nil, fmt.Errorf("account with the unique name exists (unique name = %v)", uniqueName)
	}

	id, err := s.getSnowflakeID(ctx, civitas)
	if err != nil {
		return nil, err
	}
	createdAt := time.Now()
	updatedAt := time.Now()

	query := fmt.Sprintf(`INSERT INTO public.accounts
    (id, created_at, updated_at, unique_name, domain, civitas, display_name, tombstoned)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`)
	params := []interface{}{id, createdAt, updatedAt, uniqueName, domain, civitas, displayName, false}
	if _, err := s.db.QueryRow(ctx, query, params...); err != nil {
		return nil, err
	}

	return &models.Account{
		ID:          id,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		UniqueName:  uniqueName,
		Domain:      domain,
		Civitas:     civitas,
		DisplayName: displayName,
		Tombstoned:  false,
	}, nil
}

func (s *Service) getSnowflakeID(ctx context.Context, civitas int64) (int64, error) {
	node, err := snowflake.NewNode(civitas)
	if err != nil {
		return 0, fmt.Errorf("node creation err = %w", err)
	}

	for i := 0; i < 10; i++ {
		id := node.Generate().Int64()
		a, err := s.getAccountByID(ctx, id)
		if err != nil || a != nil {
			continue
		}
		return id, nil
	}

	return 0, errors.New("can not generate unique id for 10 times of try")
}

func (s *Service) getAccountByID(ctx context.Context, id int64) (*models.Account, error) {
	q := fmt.Sprintf(`
		SELECT
      id, created_at, updated_at, deleted_at, unique_name, domain, civitas, display_name, tombstoned
    FROM public.accounts
		WHERE id=$1;
	`)
	row, err := s.db.QueryRow(ctx, q, id)
	if err != nil {
		return nil, err
	}

	var a models.Account
	err = row.Scan(
		&a.ID,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.DeletedAt,
		&a.UniqueName,
		&a.Domain,
		&a.Civitas,
		&a.DisplayName,
		&a.Tombstoned,
	)
	if err == nil {
		return &a, nil
	}

	if errorsutils.Equals(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}

func (s *Service) getAccountByUniqueName(ctx context.Context, uniqueName string) (*models.Account, error) {
	q := fmt.Sprintf(`
		SELECT
      id, created_at, updated_at, deleted_at, unique_name, domain, civitas, display_name, tombstoned
    FROM public.accounts
		WHERE unique_name=$1;
	`)
	row, err := s.db.QueryRow(ctx, q, uniqueName)
	if err != nil {
		return nil, err
	}

	var a models.Account
	err = row.Scan(
		&a.ID,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.DeletedAt,
		&a.UniqueName,
		&a.Domain,
		&a.Civitas,
		&a.DisplayName,
		&a.Tombstoned,
	)
	if err == nil {
		return &a, nil
	}

	if errorsutils.Equals(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}
