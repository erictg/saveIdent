package dto

type UpdateRequestDTO struct{
	Lat float32
	Lon float32
	Status int
	Device_id int
}

type UpdateResponseDTO struct {
	IsAffected bool
}

