package queries

import (
	"database/sql"
	"saveIdent.com/server/sqlInterfaceService/dbUtils"
	"errors"
	sq "github.com/Masterminds/squirrel"
)

type PermissionQueries struct{
	db *sql.DB
}

func NewPermissionQueries() (*PermissionQueries, error){
	holderDb, err := dbUtils.GetDB()
	if err != nil{
		//logger.Error(err, "NewPermissionQueries", "failed to establish database connection")
		return nil, err
	}

	return &PermissionQueries{db:holderDb}, nil
}

func NewPermissionQueriesWithDb(passDb *sql.DB) (*PermissionQueries, error){
	if passDb == nil{
		err := errors.New("passed in null database connection")
		//logger.Error(err, "NewPermissionQueriesWithDb", "passed in null database connection")
		return nil, err
	}
	return &PermissionQueries{db:passDb}, nil
}

func (p *PermissionQueries) Create(name string) error{
	//validate name
	if len(name) < 1{
		err := errors.New("name is empty")
		//logger.Error(err, "Create_Permission", "name is empty")
		return err
	}

	query := sq.Insert("permissions").Columns("name").Values(name)

	//logger.InfoWithField("Create_Permission", "querying", "query", query)

	result, err := query.RunWith(p.db).Exec()
	if err != nil{
		//logger.Error(err, "Create_Permission", "query failed")
		return err
	}

	rows, qErr := result.RowsAffected()
	if qErr != nil{
		//logger.Error(qErr, "Create_Permission", "info failed")
	}

	//rows should be one if it was added successfully
	if rows != 1{
		err := errors.New("name exists")
		//logger.Error(err, "Create_Permission", "name exists")
		return err
	}

	return nil
}

func (p *PermissionQueries) DeleteByName(name string) error{
	//validate name
	if len(name) < 1{
		err := errors.New("name is empty")
		//logger.Error(err, "DeleteByName_Permission", "name is empty")
		return err
	}


	query := sq.Delete("permissions").Where(sq.Eq{"permissions.name": name})
	//logger.InfoWithField("DeleteByName_Permission", "querying", "query", query)

	result, err := query.RunWith(p.db).Exec()
	if err != nil{
		//logger.Error(err, "DeleteByName_Permission", "query failed")
		return err
	}

	rows, qErr := result.RowsAffected()
	if qErr != nil{
		//logger.Error(qErr, "DeleteByName_Permission", "info failed")
	}

	//rows should be one if it was added successfully
	if rows != 1{
		err := errors.New("not found")
		//logger.Error(err, "DeleteByName_Permission", "not found")
		return err
	}

	//it worked
	return nil
}

func (p *PermissionQueries) DeleteByID(id uint) error{
	query := sq.Delete("permissions").Where(sq.Eq{"permissions.id": id})
	//logger.InfoWithField("DeleteById_Permission", "querying", "query", query)

	result, err := query.RunWith(p.db).Exec()
	if err != nil{
		return err
	}

	rows, qErr := result.RowsAffected()
	if qErr != nil{
	}

	//rows should be one if it was added successfully
	if rows != 1{
		err := errors.New("not found")
		return err
	}

	//it worked
	return nil
}

func (p *PermissionQueries) Exists(name string) (bool, error){
	//validate name
	if len(name) < 1{
		err := errors.New("name is empty")
		return false, err
	}

	query := sq.Select("name").From("permissions").Where(sq.Eq{"name": name})

	result, err := query.RunWith(p.db).Exec()
	if err != nil{
		return false, err
	}

	rows, qErr := result.RowsAffected()
	if qErr != nil{
	}

	//rows should be one if it was added successfully
	if rows != 1{
		err := errors.New("not found")
		return false, err
	}

	//it worked
	return true, nil
}


func (p *PermissionQueries) ExistsId(id uint) (bool, error){

	query := sq.Select("name").From("permissions").Where(sq.Eq{"id": id})

	result, err := query.RunWith(p.db).Exec()
	if err != nil{
		return false, err
	}

	rows, qErr := result.RowsAffected()
	if qErr != nil{
	}

	//rows should be one if it was added successfully
	if rows != 1{
		err := errors.New("not found")
		return false, err
	}

	//it worked
	return true, nil
}

func (p *PermissionQueries) GetId(name string) (uint, error){
	//validate name
	if len(name) < 1{
		err := errors.New("name is empty")
		return 0, err
	}

	query := sq.Select("id").From("permissions").Where(sq.Eq{"name": name})
	rows, err := query.RunWith(p.db).Query()
	if err != nil{
		return 0, err
	}
	defer rows.Close()

	id := uint(0)

	count := 0
	for rows.Next(){
		if count != 0{
			err = errors.New("too many columns returned")
			return 0, err
		}
		err := rows.Scan(&id)
		if err != nil{
			return 0, err
		}
		count ++
	}

	return id, nil
}
