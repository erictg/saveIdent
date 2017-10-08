package queries


import (
	"database/sql"
	"saveIdent.com/server/sqlInterfaceService/dbUtils"
	"errors"
	sq "github.com/Masterminds/squirrel"
)

type RoleQueries struct{
	db *sql.DB
}

func NewRoleQueries() (*RoleQueries, error){
	holderDb, err := dbUtils.GetDB()
	if err != nil{
		//.Error(err, "NewRoleQueries", "failed to establish database connection")
		return nil, err
	}

	return &RoleQueries{db:holderDb}, nil
}

func NewRoleQueriesWithDb(passDb *sql.DB) (*RoleQueries, error){
	if passDb == nil{
		err := errors.New("passed in null database connection")
		//.Error(err, "NewRoleQueriesWithDb", "passed in null database connection")
		return nil, err
	}
	return &RoleQueries{db:passDb}, nil
}

func (p *RoleQueries) Create(name string) error{
	//validate name
	if len(name) < 1{
		err := errors.New("name is empty")
		//.Error(err, "Create_Role", "name is empty")
		return err
	}

	query := sq.Insert("roles").Columns("name").Values(name)
	//.InfoWithField("Create_Role", "querying", "query", query)

	result, err := query.RunWith(p.db).Exec()
	if err != nil{
		//.Error(err, "Create_Role", "query failed")
		return err
	}

	rows, qErr := result.RowsAffected()
	if qErr != nil{
		//.Error(qErr, "Create_Role", "info failed")
	}

	//rows should be one if it was added successfully
	if rows != 1{
		err := errors.New("name exists")
		//.Error(err, "Create_Role", "name exists")
		return err
	}

	return nil
}

func (p *RoleQueries) DeleteByName(name string) error{
	//validate name
	if len(name) < 1{
		err := errors.New("name is empty")
		//.Error(err, "DeleteByName_Role", "name is empty")
		return err
	}

	query := sq.Delete("roles").Where(sq.Eq{"name": name})
	//.InfoWithField("DeleteByName_Role", "querying", "query", query)

	result, err := p.db.Exec(query.ToSql())
	if err != nil{
		//.Error(err, "DeleteByName_Role", "query failed")
		return err
	}

	rows, qErr := result.RowsAffected()
	if qErr != nil{
		//.Error(qErr, "DeleteByName_Role", "info failed")
	}

	//rows should be one if it was added successfully
	if rows != 1{
		err := errors.New("not found")
		//.Error(err, "DeleteByName_Role", "not found")
		return err
	}

	//it worked
	return nil
}

func (p *RoleQueries) DeleteByID(id uint) error{
	query := sq.Delete("roles").Where(sq.Eq{"id": id})
	//.InfoWithField("DeleteById_Role", "querying", "query", query)

	result, err := p.db.Exec(query.ToSql())
	if err != nil{
		//.Error(err, "DeleteById_Role", "query failed")
		return err
	}

	rows, qErr := result.RowsAffected()
	if qErr != nil{
		//.Error(qErr, "DeleteById_Role", "info failed")
	}

	//rows should be one if it was added successfully
	if rows != 1{
		err := errors.New("not found")
		//.Error(err, "DeleteById_Role", "not found")
		return err
	}

	//it worked
	return nil
}

func (p *RoleQueries) Exists(name string) (bool, error){
	//validate name
	if len(name) < 1{
		err := errors.New("name is empty")
		//.Error(err, "Exists_Role", "name is empty")
		return false, err
	}

	query := sq.Select("name").From("roles").Where(sq.Eq{"name": name})
	//.InfoWithField("Exists_Role", "querying", "query", query)

	result, err := p.db.Exec(query.ToSql())
	if err != nil{
		//.Error(err, "Exists_Role", "query failed")
		return false, err
	}

	rows, qErr := result.RowsAffected()
	if qErr != nil{
		//.Error(qErr, "Exists_Role", "info failed")
	}

	//rows should be one if it was added successfully
	if rows != 1{
		err := errors.New("not found")
		//.Error(err, "Exists_Role", "not found")
		return false, err
	}

	//it worked
	return true, nil
}

func (p *RoleQueries) ExistsId(id uint) (bool, error){

	query := sq.Select("name").From("roles").Where(sq.Eq{"id": id})
	//.InfoWithField("Exists_Role", "querying", "query", query)

	result, err := p.db.Exec(query.ToSql())
	if err != nil{
		//.Error(err, "Exists_Role", "query failed")
		return false, err
	}

	rows, qErr := result.RowsAffected()
	if qErr != nil{
		//.Error(qErr, "Exists_Role", "info failed")
	}

	//rows should be one if it was added successfully
	if rows != 1{
		err := errors.New("not found")
		//.Error(err, "Exists_Role", "not found")
		return false, err
	}

	//it worked
	return true, nil
}

func (p *RoleQueries) GetId(name string) (uint, error){
	//validate name
	if len(name) < 1{
		err := errors.New("name is empty")
		//.Error(err, "Exists_Role", "name is empty")
		return 0, err
	}

	rows, err := sq.Select("id").From("roles").Where(sq.Eq{"name": name}).RunWith(p.db).Query()
	if err != nil{
		//.Error(err, "GetId", "query failed")
		return 0, err
	}
	defer rows.Close()

	id := uint(0)

	count := 0
	for rows.Next(){
		if count != 0{
			err = errors.New("too many columns returned")
			//.Error(err, "GetId", "too many columns returned")
			return 0, err
		}
		err := rows.Scan(&id)
		if err != nil{
			//.Error(err, "GetId", "failed to get id")
			return 0, err
		}
		count ++
	}

	return id, nil
}