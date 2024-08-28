package repo

import (
	"context"
	"database/sql"
	"fmt"
	"utility/internal/entities"

	"github.com/lib/pq"
	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/utils"
)

// LookupRepo is responsible for handling lookup-related data.
type LookupRepo struct {
	db *sql.DB
}

// LookupRepoImply is the interface defining the methods for working with lookup data.
type LookupRepoImply interface {
	GetLookupByIdList(ctx context.Context, idList entities.LookupIDs) (entities.LookupData, error)
	GetLookupByTypeName(ctx context.Context, name string, filter map[string]string) ([]entities.Lookup, error)
}

// NewLookupRepo creates a new LookupRepo instance.
func NewLookupRepo(db *sql.DB) LookupRepoImply {
	return &LookupRepo{db: db}
}

// GetLookupByTypeName retrieves lookup data based on specified parameters.
func (lookup *LookupRepo) GetLookupByTypeName(ctx context.Context, name string, filter map[string]string) ([]entities.Lookup, error) {

	var (
		log        = logger.Log().WithContext(ctx)
		lookupData []entities.Lookup
		args       []interface{}
	)
	query := `
		SELECT l.id, COALESCE(l.name, ''), COALESCE(l.value, ''), l.lookup_type_id  
		FROM lookup l
		JOIN lookup_type lt
		ON  lt.id = l.lookup_type_id 
		WHERE lt.name = LOWER($1)`

	args = append(args, name)

	if !utils.IsEmpty(filter["value"]) {
		query = fmt.Sprintf("%s AND l.value = LOWER($2)", query)
		args = append(args, filter["value"])
	}

	rows, err := lookup.db.QueryContext(ctx, query, args...)

	if err != nil {
		log.Errorf("[LookupRepo][GetLookupByTypeName], Error : %s", err.Error())
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var lookupInfo entities.Lookup
		err := rows.Scan(&lookupInfo.ID, &lookupInfo.Name, &lookupInfo.Value, &lookupInfo.LookupTypeId)

		if err != nil {
			log.Errorf("[LookupRepo][GetLookupByTypeName], Error : %s", err.Error())
			return nil, err
		}

		lookupData = append(lookupData, lookupInfo)
	}

	return lookupData, nil
}

// GetLookupByIdList retrieves lookup data based on array of lookup Ids
func (LookupRepo *LookupRepo) GetLookupByIdList(ctx context.Context, idList entities.LookupIDs) (entities.LookupData, error) {

	var (
		lookupData []entities.Lookup
		log        = logger.Log().WithContext(ctx)
		foundIds   = make(map[int64]bool)
		invalidIDs []int64
		output     entities.LookupData
		ids        = idList.ID
	)

	// SQL query to retrieve lookup data by ID
	getLookup := `
			SELECT id,
			COALESCE(name, '') AS name, 
			COALESCE(value, '') AS value,
			COALESCE(description, '') AS description,
			COALESCE(position, '0')AS position, 
			COALESCE(language_identifier, '' )AS language_identifier
			FROM lookup
			WHERE id = ANY($1)		
	`

	// Execute the query with the provided ID
	res, err := LookupRepo.db.QueryContext(ctx, getLookup, pq.Array(ids))

	if err != nil {
		log.Errorf("[LookupRepo][GetLookupByIdList], Error : %s", err.Error())
		return entities.LookupData{}, err
	}

	defer res.Close()
	// Scan the resulting row into the lookup entity

	for res.Next() {

		var lookup entities.Lookup
		err := res.Scan(
			&lookup.ID,
			&lookup.Name,
			&lookup.Value,
			&lookup.Description,
			&lookup.Position,
			&lookup.LanguageIdentifier,
		)

		if err != nil {
			log.Errorf("[LookupRepo][GetLookupByIdList], Error : %s", err.Error())
			return entities.LookupData{}, err
		}
		lookupData = append(lookupData, lookup)
		foundIds[lookup.ID] = true
	}

	// Handle errors from the scan operation
	if err != nil {
		if err == sql.ErrNoRows {
			log.Errorf("[LookupRepo][GetLookupByIdList], No Content, Error : %s", err.Error())
			return entities.LookupData{}, err
		}
		log.Errorf("[LookupRepo][GetLookupByIdList], .Error : %s", err.Error())
		return entities.LookupData{}, err
	}

	//To check missing IDs
	for _, id := range ids {
		if !foundIds[id] {
			invalidIDs = append(invalidIDs, id)
		}
	}
	//Remove redundancy in array
	invalidIDs = RemoveDuplicateValues(invalidIDs)

	output.Lookup = lookupData
	output.InvalidLookupIds = invalidIDs
	return output, nil
}

// To remove duplicate values in an array
func RemoveDuplicateValues(invalidIDs []int64) []int64 {
	keys := make(map[int64]bool)
	list := []int64{}

	for _, entry := range invalidIDs {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
