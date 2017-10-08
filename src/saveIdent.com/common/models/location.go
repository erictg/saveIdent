package models


type Location struct {
	Lat			float32		`json:"lat"`
	Lon			float32		`json:"lon"`
	Status 		bool	 	`json:"status"`
	DeviceId	uint	 	`json:"device_id"`
}
