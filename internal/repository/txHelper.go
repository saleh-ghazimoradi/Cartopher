package repository

import "gorm.io/gorm"

func exec(db, tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return db
}
