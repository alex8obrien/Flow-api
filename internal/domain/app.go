package domain

import "encoding/json"

// Version sets the main application Version.
const Version string = "0.1.0"

type Status string

const (
	OK   Status = "ok"
	DOWN Status = "down"
)

type Readyz struct {
	version string
	overall Status
	db      Status
}

func (r Readyz) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Version string `json:"version"`
		Overall Status `json:"overall"`
		DB      Status `json:"db"`
	}{
		Version: r.version,
		Overall: r.overall,
		DB:      r.db,
	})
}

func NewReadyz() *Readyz {
	return &Readyz{
		version: Version,
		overall: DOWN,
		db:      DOWN,
	}
}

func (r *Readyz) DBChecked() {
	r.db = OK
	r.determineOverallStatus()
}

func (r *Readyz) determineOverallStatus() {
	// Assumes an OK Status
	s := OK

	// Iterate through all dependencies
	if r.db == DOWN {
		s = DOWN
	}

	// Set overall result
	r.overall = s
}

func (r *Readyz) Overall() Status {
	return r.overall
}
