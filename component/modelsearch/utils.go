package modelsearch

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func applyOrderBy(db *gorm.DB, orderBy string) *gorm.DB {
	parts := strings.Split(orderBy, " ")
	if len(parts) == 2 {
		sortedField := parts[0]
		direction := strings.ToLower(parts[1])

		isDesc := direction == "desc"

		return db.Order(
			clause.OrderByColumn{
				Column: clause.Column{Name: sortedField},
				Desc:   isDesc,
			},
		)
	}

	return db
}
