package queries

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"errors"
	"saveIdent.com/server/sqlInterfaceService/dbUtils"
	//"saveIdent.com/common/models"
	"log"
)

type DeviceQueries struct{
	db *sql.DB
}

func NewDeviceQueries() (*DeviceQueries, error){
	holderDb, err := dbUtils.GetDB()
	if err != nil{
		//.Error(err, "NewUserQueries", "failed to establish database connection")
		return nil, err
	}

	return &DeviceQueries{db:holderDb}, nil
}

func NewDeviceQueriesWithDb(passDb *sql.DB) (*DeviceQueries, error){
	if passDb == nil{
		err := errors.New("passed in null database connection")
		//.Error(err, "NewUserQueriesWithDb", "passed in null database connection")
		return nil, err
	}
	return &DeviceQueries{db:passDb}, nil
}

func (q *DeviceQueries) Create() error{
	query := sq.Insert("device")
	s, _, _ :=query.ToSql()
	return dbUtils.ExecuteLineCheckQuery(s, "", q.db)
}

func (q *DeviceQueries) CreateReceiver() error{
	query := sq.Insert("device").Columns("isReceiver").Values(true)
	s, _, _ :=query.ToSql();
	return dbUtils.ExecuteLineCheckQuery(s, "", q.db)
}

func (q *DeviceQueries) AssociateWithUser(deviceId uint, userId uint) error{
	query := sq.Update("device").Set("user_id", userId).Where(sq.Eq{"id":deviceId})
	s, _, _ := query.ToSql()
	return dbUtils.ExecuteLineCheckQuery(s, "", q.db)
}

func (q *DeviceQueries) DissociateWithUser(deviceId uint, userId uint) error{
	query := sq.Update("device").Set("user_id", "").Where(sq.Eq{"id":deviceId})
	s, _, _ := query.ToSql()
	return dbUtils.ExecuteLineCheckQuery(s, "", q.db)
}

func (q *DeviceQueries) GetUserIdForDevice(deviceId uint) (uint, error){
	var id uint = 0
	query := sq.Select("user_id").From("device").Where(sq.Eq{"id": deviceId})
	result, err := query.Query()
	if err != nil{
		log.Println("error occured during query")
		log.Println(err)
		return 0, err
	}

	for result.Next(){
		if err := result.Scan(&id); err != nil{
			log.Println(err)
			log.Println("something went wrong decoding query")
			return 0, err
		}
	}
	return id, nil
}

func (q *DeviceQueries) GetDeviceIdForUser(userId uint) (uint, error){
	var id uint = 0
	query := sq.Select("id").From("device").Where(sq.Eq{"user_id": userId})
	result, err := query.Query()
	if err != nil{
		log.Println("error occured during query")
		log.Println(err)
		return 0, err
	}

	for result.Next(){
		if err := result.Scan(&id); err != nil{
			log.Println(err)
			log.Println("something went wrong decoding query")
			return 0, err
		}
	}
	return id, nil
}