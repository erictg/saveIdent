package queries

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"errors"
	"saveIdent.com/server/sqlInterfaceService/dbUtils"
	"saveIdent.com/common/models"
	"log"
)

type UserQueries struct{
	db *sql.DB
}

func NewUserQueries() (*UserQueries, error){
	holderDb, err := dbUtils.GetDB()
	if err != nil{
		//.Error(err, "NewUserQueries", "failed to establish database connection")
		return nil, err
	}

	return &UserQueries{db:holderDb}, nil
}

func NewUserQueriesWithDb(passDb *sql.DB) (*UserQueries, error){
	if passDb == nil{
		err := errors.New("passed in null database connection")
		//.Error(err, "NewUserQueriesWithDb", "passed in null database connection")
		return nil, err
	}
	return &UserQueries{db:passDb}, nil
}

func (q *UserQueries) Create(name, password, email, address string) error{
	query := sq.Insert("users").Columns("name", "password", "email", "address").Values(name, password, email, address)


	if result, err := query.RunWith(q.db).Exec(); err != nil{
		return err
	}else{
		if rows, err := result.RowsAffected(); rows != 1{
			//.Error(err, "DeleteUserById", "failed to delete user")
			return err
		}

		return nil
	}

}

func (q *UserQueries) GetById(id uint) (models.User, error){

	query := sq.Select("name", "address", "email").Where(sq.Eq{"id": id})

	result, err := query.RunWith(q.db).Query()
	if err != nil{
		log.Println(err)
		return models.User{}, err
	}

	defer result.Close()

	user := models.User{}

	for result.Next(){
		var name string = ""
		var address string = ""
		var email string = ""
		if err = result.Scan(&name, &address, &email); err != nil{
			return models.User{}, err
		}
		user.Name = name
		user.Address = address
		user.Email = email
	}

	return user, nil
}

//deletion queries
func (q *UserQueries) DeleteUserById(id uint) error{
	query := sq.Delete("users").Where(sq.Eq{"id":id})
	result, err := query.RunWith(q.db).Exec()
	if err != nil{
		//.Error(err, "DeleteUserById", "failed to delete user")
		return err
	}

	if rows, err := result.RowsAffected(); rows != 1{
		//.Error(err, "DeleteUserById", "failed to delete user")
		return err
	}

	return nil
}

func (q *UserQueries) DeleteUserByName(name string) error{
	if len(name) < 1{
		err := errors.New("invalid name")
		//.Error(err, "DeleteUserByName", "invalid name")
		return err
	}

	query, _, err := sq.Delete("users").Where(sq.Eq{"name": name}).ToSql()
	if err != nil{
		//.Error(err, "DeleteUserByName", "generation of query string failed")
		return err
	}
	return dbUtils.ExecuteLineCheckQuery(query, "DeleteUserByName", q.db)
}


func (q *UserQueries) UpdatePassword(password string, id uint) error {
	err := validateString("UpdatePassword", password)
	if err != nil{
		return err
	}

	query := createUpdateQuery("password", password, id)

	return dbUtils.ExecuteLineCheckQuery(query, "UpdatePassword", q.db)
}

//roles
func (q *UserQueries) AddUserToRoleById(roleId uint, userId uint) error{

	query, _, err := sq.Insert("user_roles_relation_table").Columns("user_id", "role_id").Values(userId, roleId).ToSql()
	if err != nil{
		//.Error(err, "AddUserToRoleById", "generation of query string failed")
		return err
	}

	return dbUtils.ExecuteLineCheckQuery(query, "AddUserToRole", q.db)
}

func (q *UserQueries) AddUserToRoleByName(role string, userId uint) error{
	roleQuery := RoleQueries{db:q.db}

	roleId, err := roleQuery.GetId(role)
	if err != nil{
		//.Error(err, "AddUserToRoleByName", "name not found")
		return err
	}

	return q.AddUserToRoleById(roleId, userId)
}

func (q *UserQueries) UserInRoleByName(name string, userId uint) error{
	roleQuery := RoleQueries{db:q.db}

	roleId, err := roleQuery.GetId(name)
	if err != nil{
		//.Error(err, "UserInRoleByName", "failed to find id")
		return err
	}

	query, _, err := sq.Select("user_id").From("user_roles_relation_table").Where(sq.Eq{"user_id": userId, "role_id": roleId}).ToSql()
	if err != nil{
		//.Error(err, "UserInRoleByName", "generation of query string failed")
		return err
	}
	err = dbUtils.ExecuteLineCheckQuery(query, "UserInRoleByName", q.db)
	if err != nil{
		//.Error(err, "UserInRoleByName", "failed to execute query")
		return err
	}

	return nil
}

func (q *UserQueries) UserInRoleById(roleId uint, userId uint) error{

	//("select user_roles_relation_table.user_id from user_roles_relation_table where user_id = %v and role_id = %v;",
	query, _, err := sq.Select("user_id").From("user_roles_relation_table").Where(sq.Eq{"user_id": userId, "role_id": roleId}).ToSql()
	if err != nil{
		//.Error(err, "UserInRoleById", "generation of query string failed")
		return err
	}

	err = dbUtils.ExecuteLineCheckQuery(query, "UserInRoleByName", q.db)
	if err != nil{
		//.Error(err, "UserInRoleByName", "failed to execute query")
		return err
	}

	return nil
}

func (q *UserQueries) RemoveUserFromRoleById(roleId uint, userId uint) error{

	query, _, err := sq.Delete("user_roles_relation_table").Where(sq.Eq{"role_id": roleId, "user_id": userId}).ToSql()
	if err != nil{
		//.Error(err, "RemoveUserFromRoleById", "generation of query string failed")
		return err
	}
	return dbUtils.ExecuteLineCheckQuery(query, "RemoveUserFromRoleById", q.db)
}

func (q *UserQueries) RemoveUserFromRoleByName(role string, userId uint) error{
	roleQuery := RoleQueries{db:q.db}

	roleId, err := roleQuery.GetId(role)
	if err != nil{
		//.Error(err, "RemoveUserFromRoleByName", "remove role")
		return err
	}

	return q.RemoveUserFromRoleById(roleId, userId)
}

func (q *UserQueries) GetUserRolesById(id uint) ([]models.Role, error){
	subQuery := "(select role_id from user_roles_relation_table where user_roles_relation_table.user_id = " + string(id) + ") as q on role_id = roles.id"
	query := sq.Select("roles.id", "roles.name").From("roles").LeftJoin(subQuery)

	rows, err := query.RunWith(q.db).Query()
	if err != nil{
		//.Error(err, "GetUserRolesById", "query fucked up")
		return nil, err
	}
	defer rows.Close()

	roles := []models.Role{}

	for rows.Next(){
		var name string = ""
		var id uint = 0
		err = rows.Scan(&id, &name)
		if err != nil{
			//.Error(err, "GetUserRolesById", "row scan failed")
			return nil, err
		}
		roles = append(roles, models.Role{Name:name, Id:id})
	}

	return roles, nil
}

//permission
func (q *UserQueries) AddUserToPermissionById(permissionId uint, userId uint) error{

	query, _, err := sq.Insert("user_permission_relation_table").Columns("user_id", "permission_id").Values(userId, permissionId).ToSql()
	if err != nil{
		//.Error(err, "AddUserToPermissionById", "generation of query string failed")
		return err
	}
	return dbUtils.ExecuteLineCheckQuery(query, "AddUserToPermissionById", q.db)
}

func (q *UserQueries) AddUserToPermissionByName(permission string, userId uint) error{
	permissionQuery := PermissionQueries{db:q.db}

	permissionId, err := permissionQuery.GetId(permission)
	if err != nil{
		//.Error(err, "AddUserToPermissionByName", "name not found")
		return err
	}

	return q.AddUserToRoleById(permissionId, userId)
}

func (q *UserQueries) UserInPermissionByName(name string, userId uint) error{
	permissionQuery := PermissionQueries{db:q.db}

	permissionId, err := permissionQuery.GetId(name)
	if err != nil{
		//.Error(err, "UserInPermissionByName", "name not found")
		return err
	}

	//query := fmt.Spruintf("select user_permission_relation_table.user_id from user_permission_relation_table where user_id = %v and permission_id = %v;",
	//	userId, permissionId)
	query, _, err := sq.Select("user_id").From("user_permission_relation_table").Where(sq.Eq{"user_id": userId, "permission_id": permissionId}).ToSql()
	if err != nil{
		//.Error(err, "UserInPermissionByName", "generation of query string failed")
		return err
	}
	err = dbUtils.ExecuteLineCheckQuery(query, "AddUserToPermissionByName", q.db)
	if err != nil{
		//.Error(err, "UserInPermissionByName", "failed to execute query")
		return err
	}

	return nil
}

func (q *UserQueries) UserInPermissionById(permissionId uint, userId uint) error{

	query, _, err := sq.Select("user_id").From("user_permission_relation_table").Where(sq.Eq{"user_id": userId, "permission_id": permissionId}).ToSql()
	if err != nil{
		//.Error(err, "UserInPermissionById", "generation of query string failed")
		return err
	}
	err = dbUtils.ExecuteLineCheckQuery(query, "UserInPermissionById", q.db)
	if err != nil{
		//.Error(err, "UserInPermissionById", "failed to execute query")
		return err
	}

	return nil
}

func (q *UserQueries) RemoveUserFromPermissionById(permissionId uint, userId uint) error{

	query, _, err := sq.Delete("user_permission_relation_table").Where(sq.Eq{"permission_id": permissionId, "user_id": userId}).ToSql()

	if err != nil{
		//.Error(err, "RemoveUserFromPermissionById", "generation of query string failed")
		return err
	}
	return dbUtils.ExecuteLineCheckQuery(query, "RemoveUserFromPermissionById", q.db)
}

func (q *UserQueries) RemoveUserFromPermissionByName(permission string, userId uint) error{
	permissionQuery := PermissionQueries{db:q.db}

	permissionId, err := permissionQuery.GetId(permission)
	if err != nil{
		//.Error(err, "RemoveUserFromPermissionByName", "name not found")
		return err
	}

	return q.RemoveUserFromRoleById(permissionId, userId)
}

func (q *UserQueries) GetUserPermissionById(id uint) ([]models.Permission, error){
	subQuery := "(select permission_id from user_permission_relation_table where user_permission_relation_table.user_id = " + string(id) + ") as q on q.permission_id = permission.id"
	query := sq.Select("permission.id", "permission.name").From("permissions").LeftJoin(subQuery)

	rows, err := query.RunWith(q.db).Query()
	if err != nil{
		//.Error(err, "GetUserPermissionsById", "query fucked up")
		return nil, err
	}
	defer rows.Close()


	permissions := []models.Permission{}

	for rows.Next(){
		var name string = ""
		var id uint = 0
		err = rows.Scan(&id, &name)
		if err != nil{
			//.Error(err, "GetUserPermissionById", "row scan failed")
			return nil, err
		}
		permissions = append(permissions, models.Permission{Name:name, Id:id})
	}

	return permissions, nil
}
//Get the user
//- get complete profile
//- get profile without roles and permissions

func (q *UserQueries) GetUserWithoutRolesAndPermissionsById(id uint) (*models.User, error){
	query := sq.Select("id", "name", "email", "address").
		From("users").Where(sq.Eq{"id": id})

	rows, err := query.RunWith(q.db).Query()
	if err != nil{
		//.Error(err, "GetUserWithoutRolesAndPermissionsById", "query failed")
		return nil, err
	}
	defer rows.Close()

	user := new(models.User)
	for rows.Next(){
		rows.Scan(&user.Id, &user.Name, &user.Email, &user.Address)
		return nil, err
	}

	err = rows.Close()
	if err != nil{
		//.Error(err, "GetUserWithoutRolesAndPermissionsById", "failed to close rows")
		return nil, err
	}
	return user, nil
}

func (q *UserQueries) GetUserWithoutRolesAndPermissionByName(name string) (*models.User, error){
	query := sq.Select("id", "name", "email", "address").
		From("users").Where(sq.Eq{"name": name})

	rows, err := query.RunWith(q.db).Query()
	if err != nil{
		//.Error(err, "GetUserWithoutRolesAndPermissionsByName", "query failed")
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	user := new(models.User)
	for rows.Next(){
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Address)
		if err != nil{
			return nil, err
		}

	}


	if err != nil{
		//.Error(err, "GetUserWithoutRolesAndPermissionsByName", "failed to close rows")
		return nil, err
	}
	log.Println(user)
	return user, nil
}

func (q *UserQueries) GetUserById(id uint) (*models.User, error){
	user, err := q.GetUserWithoutRolesAndPermissionsById(id)
	if err != nil{
		//.Error(err, "GetUserById", "failed to retrieve user")
		return nil, err
	}

	roles, err := q.GetUserRolesById(id)
	if err != nil{
		//.Error(err, "GetUserById", "failed to get roles")
		return nil, err
	}

	user.Roles = roles

	permissions, err := q.GetUserPermissionById(id)
	if err != nil{
		//.Error(err, "GetUserById", "failed to get permissions")
		return nil, err
	}

	user.Permissions = permissions

	//.Info("GetUserById", "successfully retrieved user")
	return user, nil
}

func (q *UserQueries) GetUserByName(name string) (*models.User, error){
	user, err := q.GetUserWithoutRolesAndPermissionByName(name)
	if err != nil || user == nil{
		//.Error(err, "GetUserById", "failed to retrieve user")
		return nil, errors.New("something dun fucked up")
	}

	log.Println(user)

	id := user.Id

	roles, err := q.GetUserRolesById(id)
	if err != nil{
		//.Error(err, "GetUserById", "failed to get roles")
		return nil, err
	}

	user.Roles = roles

	permissions, err := q.GetUserPermissionById(id)
	if err != nil{
		//.Error(err, "GetUserById", "failed to get permissions")
		return nil, err
	}

	user.Permissions = permissions

	//.Info("GetUserById", "successfully retrieved user")
	return user, nil
}

//helper functions
func createUpdateQuery(column string, newString string, id uint) string{

	toReturn, _, err := sq.Update("users").Set(column, newString).Where(sq.Eq{"id": id}).ToSql()
	if err != nil{
		//.Error(err, "createUpdateQuery", "query gen got screwed up")
	}
	return toReturn
}

func validateString(callingFunction string, str string) error{
	if len(str) < 1{
		err := errors.New("invalid string")
		//.Error(err, callingFunction, "invalid string")
		return err
	}

	return nil
}