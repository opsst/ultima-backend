package models

type Ingredient struct {
	Name    string        `json:"name,omitempty" `
	Rate    string        `json:"rate,omitempty" `
	Calling string        `json:"calling,omitempty"`
	Func    []interface{} `json:"func,omitempty"`
	Irr     string        `json:"irr,omitempty"`
	Come    string        `json:"come,omitempty"`
	Cosing  string        `json:"cosing,omitempty"`
	Quick   []interface{} `json:"quick,omitempty"`
	Detail  string        `json:"detail,omitempty"`
	Proof   []interface{} `json:"proof,omitempty"`
	Link    string        `json:"link,omitempty"`
	IsTryOn string        `json:"istryon,omitempty"`
}
