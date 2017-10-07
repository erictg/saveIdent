package dto

type UpdateRequestDTO struct{
	Lat float32		`json:"lat"`
	Lon float32		`json:"lon"`
	Status int		`json:"status"`
	Device_id int	`json:"device_id"`
}

type UpdateResponseDTO struct {
	IsAffected bool
}

