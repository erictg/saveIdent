package query

import "saveIdent.com/server/deviceInputService/dto"

type ShardInfo struct {
	Total int		`json:"total"`
	Success int		`json:"successful"`
	Failed int		`json:"failed"`
}

type DocInfo struct {
	Index string				`json:"_index"`
	Type string					`json:"_type"`
	Id string					`json:"_id"`
	Score float32				`json:"_score"`
	Source dto.UpdateRequestDTO	`json:"_source"`
	Sort []string				`json:"sort"`
}

type HitsInfo struct {
	Total int			`json:"total"`
	MaxScore float32	`json:"max_score"`
	Hits []DocInfo		`json:"hits"`
}

type SearchInfoResponse struct {
	Took int			`json:"took"`
	TimeOut bool		`json:"timed_out"`
	Shards ShardInfo	`json:"_shards"`
	Hits HitsInfo		`json:"hits"`
}

type MatchUID struct {
	Uid string `json:"_uid"`
}

type MatchDeviceID struct {
	ID int	`json:"device_id"`
}

type MustGeoMatch struct {
	MatchAll interface{} 	`json:"match_all"`
}

type TermMatch struct {
	Status int 		`json:"status"`
}

type MustGeoStatusMatch struct {
	Term TermMatch		`json:"term"`
}

type GeoBox struct {

}

type FilterGeo struct {
	BoundingBox GeoBox	`json:"geo_bounding_box"`
}

type FilterGeoStatus struct {}

type BoolGeo struct {
	M MustGeoMatch		`json:"must"`
	F FilterGeo			`json:"filter"`
}

type BoolGeoStatus struct {
	M MustGeoStatusMatch	`json:"must"`
	F FilterGeoStatus		`json:"filter"`
}

type QueryDeviceId struct {
	Match MatchDeviceID		`json:"match"`
}

type QueryGeoShit struct {
	B BoolGeo	`json:"bool"`
}

type QueryGeoStatusShit struct {
	B BoolGeoStatus `json:"bool"`
}

type SearchDeviceIDRequest struct {
	Query QueryDeviceId		`json:"query"`
	Sort []MatchUID		`json:"sort"`
}

type SearchGeoQueryRequest struct {
	Query QueryGeoShit		`json:"query"`
	Sort []MatchUID			`json:"sort"`
}

type SearchGeoStatusRequest struct {
	Query QueryGeoStatusShit	`json:"query"`
	Sort []MatchUID				`json:"sort"`
}