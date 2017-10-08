package models

type EmergencyServices struct {
	Id 			uint 		`json:"id"`
	Name 		string 		`json:"name"`
	CenterLat 	float32 	`json:"center_lat"`
	CenterLon 	float32 	`json:"center_lon"`
	Users		[]User		`json:"users"`
}
