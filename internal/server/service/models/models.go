package models

type RequestArrayMetric []RequestMetric

type RequestMetric struct {
	Delta  *int64   `json:"delta,omitempty"`
	Value  *float64 `json:"value,omitempty"`
	Hash   *string  `json:"hash,omitempty"`
	MType  string   `param:"mType" json:"type"`
	MName  string   `param:"mName" json:"id"`
	MValue string   `param:"mValue"`
}
