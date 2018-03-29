package repositories

import (
	"database/sql"
	"github.com/zibilal/repoman/persistence"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"errors"
)

type StoresType struct {
}

type StoreOwnersType struct {
	Id           sql.NullInt64
	Fullname     sql.NullString
	Gender       sql.NullString
	Address      sql.NullString
	Email        sql.NullString
	ProvinceId   sql.NullInt64
	CityId       sql.NullInt64
	KecamatanId  sql.NullInt64
	KelurahanId  sql.NullInt64
	PostCode     sql.NullString
	PhotoKtp     sql.NullString
	PhotoWithKtp sql.NullString
	PhotoProfile sql.NullString
	HasStore     sql.NullInt64
}

type StoresOwnersRepo struct {
	savedQuery map[string]string
}

type StoresOwnerResult struct {
	LastInsertedId int64
	RowsAffected int64
}

func NewStoresOwnersRepo() *StoresOwnersRepo {
	repo := new(StoresOwnersRepo)
	repo.savedQuery = make(map[string]string)
	return repo
}

func (s *StoresOwnersRepo) AddQuery(name, query string) {
	s.savedQuery[name]=query
}

func (s *StoresOwnersRepo) GetQuery(name string) string {
	return s.savedQuery[name]
}

func (s *StoresOwnersRepo) Find(dbContext persistence.DatabaseContext, queryName string, dest interface{},  data ...interface{}) (interface{}, error) {
	query, found := s.savedQuery[queryName]

	if !found {
		return nil, fmt.Errorf("could not find any query named %s", queryName)
	}

	if len(data) <= 1 {
		return nil, fmt.Errorf("please provide the destination object")
	}

	db := dbContext.Db().(*sqlx.DB)
	err := db.Select(&dest, query, data...)

	if err != nil {
		return nil, err
	}

	return "OK", nil
}

func (s *StoresOwnersRepo) Update(dbContext persistence.DatabaseContext, queryName string, data ...interface{}) (interface{}, error) {
	query, found := s.savedQuery[queryName]

	if !found {
		return nil, fmt.Errorf("could not find any query named %s", queryName)
	}

	if len(data) <= 1 {
		return nil, fmt.Errorf("please provide the destination object")
	}

	db, found := dbContext.Db().(*sql.DB)
	if !found {
		return nil, fmt.Errorf("expected context Db of type *sql.DB got %T", dbContext.Db())
	}

	if dbContext.IsTransaction() {
		tx, err := db.Begin()
		if err != nil {
			return nil, err
		}
		result, err := tx.Exec(query, data...)
		if err != nil {
			return nil, err
		}
		lastInsertedId, err := result.LastInsertId()
		rowsAffected, err := result.RowsAffected()

		log.Println("Rows affected:", rowsAffected, "LastInsertedId:", lastInsertedId, "Error:", err)

		return StoresOwnerResult{
			LastInsertedId: lastInsertedId,
			RowsAffected: rowsAffected,
		}, nil

	} else {
		result, err := db.Exec(query, data ...)
		if err != nil {
			return nil, err
		}
		tmp1, err := result.LastInsertId()
		tmp2, err := result.RowsAffected()
		log.Println("Rows affected:", tmp1, "LastInsertedId:", tmp2, "Error:", err)
		return StoresOwnerResult{
			tmp1,
			tmp2,
		}, nil
	}


}

func (s *StoresOwnersRepo) Create(dbContext persistence.DatabaseContext, queryName string, data ...interface{}) (interface{}, error) {
	query, found := s.savedQuery[queryName]

	if !found {
		return nil, fmt.Errorf("could not find query named %s", queryName)
	}

	if len(data) < 0 || len(data) > 1 {
		return nil, errors.New("create only accept one struct data")
	}


}

func (s *StoresOwnersRepo) Delete(dbContext persistence.DatabaseContext, queryName string, data ...interface{}) (interface{}, error) {

}
