package rest

// ReqData REST request and/or response data
type ReqData struct {
	State int `json:"state"`
}

// PinNumberData represents state of pin given by number
type PinNumberData struct {
	Number int `json:"number"`
	State  int `json:"state"`
}
