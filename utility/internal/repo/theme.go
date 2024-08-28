package repo

import (
	"context"
	"database/sql"
	"utility/internal/entities"

	"gitlab.com/tuneverse/toolkit/core/logger"
)

// ThemeRepo represents the repository for Theme-related operations.
type ThemeRepo struct {
	db *sql.DB
}

// ThemeRepoImply is an interface for the ThemeRepo.
type ThemeRepoImply interface {
	GetThemeByID(ctx context.Context, id string) (entities.Theme, error)
}

// NewThemeRepo creates a new instance of ThemeRepo.
func NewThemeRepo(db *sql.DB) ThemeRepoImply {
	return &ThemeRepo{db: db}
}

func (theme *ThemeRepo) GetThemeByID(ctx context.Context, id string) (entities.Theme, error) {

	var (
		res entities.Theme
		log = logger.Log().WithContext(ctx)
	)

	query := `
	SELECT 
		id,
		name,
		value,
		layout_id
	FROM theme
	WHERE id= $1
	`

	rows := theme.db.QueryRowContext(ctx, query, id)
	err := rows.Scan(
		&res.ID,
		&res.Name,
		&res.Value,
		&res.LayoutID,
	)
	if err != nil {
		log.Errorf("[ThemeRepo][GetThemeByID] Error while Scan, Error : %s", err.Error())
		return entities.Theme{}, err
	}
	return res, nil

}
