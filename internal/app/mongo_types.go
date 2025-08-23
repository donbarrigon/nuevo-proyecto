package app

// GeoPoint es un punto geografico
type GeoPoint struct {
	Type        string    `bson:"type"        json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"`
}

func NewGeoPoint(latitude, longitude float64) *GeoPoint {
	return &GeoPoint{
		Type:        "Point",
		Coordinates: []float64{longitude, latitude},
	}
}
