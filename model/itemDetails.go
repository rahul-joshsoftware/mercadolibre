package model

type ItemRequest struct{}

type ItemResponse struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type CallerInfo struct {
	LastRequestId int   `json:"lastrequestid"`
	FailureCount  int   `json:"failurecount"`
	ErrorCount    int   `json:"errorcount"`
	FailureIds    []int `json:"failureids"`
}
type ResponseChan struct {
	RequestID int
	ItemData  ItemResponse
	Error     string
}
