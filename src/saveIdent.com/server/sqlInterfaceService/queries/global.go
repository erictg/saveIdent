package queries

var Device 	*DeviceQueries
var User		*UserQueries

func Init(){
	Device, _ = NewDeviceQueries()
	User, _ = NewUserQueries()
}