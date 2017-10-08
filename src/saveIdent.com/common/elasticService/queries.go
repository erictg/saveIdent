package elasticService

import (
	"net/http"
	"fmt"
	"errors"
	"bytes"
	"encoding/json"
)

const (
	SEARCH = "/_search"
	JSON = "application/json"
)

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
	Source UpdateESDocType		`json:"_source"`
	Sort []string				`json:"sort"`
}

type HitsInfo struct {
	Total int			`json:"total"`
	MaxScore float32	`json:"max_score"`
	Hits []DocInfo		`json:"hits"`
}

type SearchResponse struct {
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

type Bound struct {
	Lat float32		`json:"lat"`
	Lng float32		`json:"lon"`
}

type BoxBounds struct {
	UpperLeft	Bound	`json:"top_left"`
	BottomRight	Bound	`json:"bottom_right"`
}

type GeoBox struct {
	Bounds BoxBounds 	`json:"geo.location"`
}

type FilterGeo struct {
	BoundingBox GeoBox	`json:"geo_bounding_box"`
}

type BoolGeo struct {
	M MustGeoMatch		`json:"must"`
	F FilterGeo			`json:"filter"`
}

type BoolGeoStatus struct {
	M MustGeoStatusMatch	`json:"must"`
	F FilterGeo				`json:"filter"`
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

type SearchGeoRequest struct {
	Query QueryGeoShit		`json:"query"`
	Sort []MatchUID			`json:"sort"`
}

type SearchGeoStatusRequest struct {
	Query QueryGeoStatusShit	`json:"query"`
	Sort []MatchUID				`json:"sort"`
}

func (db *ElasticSearchDB) SearchDeviceId(deviceId int) (SearchResponse, error) {
	var searchResponse SearchResponse

	// Construct search request
	searchRequest := SearchDeviceIDRequest{QueryDeviceId{MatchDeviceID{deviceId}}, nil}

	// Send request
	resp, err := sendJson(db, &searchRequest)
	defer resp.Body.Close()

	// If it failed give up
	if err != nil {
		if db.errLogger != nil {
			db.errLogger.Error(err, "Failed doing sumthing in search device id")
		}
		return searchResponse, err
	}

	// If it didn't fail get result
	decoder := json.NewDecoder(resp.Body)
	return searchResponse, decoder.Decode(&searchResponse)
}

func (db *ElasticSearchDB) SearchGeo(upperLeft, bottomRight Bound) (SearchResponse, error) {
	var searchResponse SearchResponse

	// Construct search request
	searchRequest := SearchGeoRequest{}

	return searchResponse, nil
}

func (db *ElasticSearchDB) SearchGeoStatus(upperLeft, bottomRight Bound, stat int) (SearchResponse,  error) {
	var searchResponse SearchResponse

	return searchResponse, nil
}

func sendJson(db *ElasticSearchDB, searchRequest interface{}) (*http.Response, error) {
	// Get client from fishie pool
	client, ok := db.clientPool.GetFishie().(http.Client)
	if !ok {
		if db.errLogger != nil {
			db.errLogger.Warn(1, "No free fishie")
		}
		fmt.Println("No free fishie")
		return &http.Response{}, errors.New("no free fishie")
	}

	defer db.clientPool.PutFishieBack(client)

	// Now shove that json request into a buffer
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)

	err := encoder.Encode(searchRequest)
	if err != nil {
		if db.errLogger != nil {
			db.errLogger.Error(err, "Failed enccoding json into buffer")
		}
		return &http.Response{}, err
	}

	// Send request
	return client.Post(db.dbIp + UPDATE + SEARCH, JSON, &b)
}