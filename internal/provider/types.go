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

type JTLabelsAPIRead struct {
	Count        int           `json:"count"`
	LabelResults []LabelResult `json:"results"`
}

type LabelResult struct {
	Id int `json:"id"`
}

type LabelDissasocBody struct {
	Id           int  `json:"id"`
	Disassociate bool `json:"disassociate"`
}
