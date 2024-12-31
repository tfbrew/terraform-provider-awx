package provider

type JTChildAPIRead struct {
	Count   int           `json:"count"`
	Results []ChildResult `json:"results"`
}

type ChildResult struct {
	Id int `json:"id"`
}

type ChildDissasocBody struct {
	Id           int  `json:"id"`
	Disassociate bool `json:"disassociate"`
}
