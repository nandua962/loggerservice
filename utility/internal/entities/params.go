package entities

// Params represents a set of query parameters commonly used for paginating and searching data.
type Params struct {
	Name   string `form:"name"`
	Sort   string `form:"sort"`
	Order  string `form:"order"`
	Code   string `form:"code"`
	Id     string `form:"id"`
	IdList string `form:"ids"`
	Iso    string `form:"iso"`
}

type LangParams struct {
	Params Params
	Status string `form:"status"`
}

type Pagination struct {
	Page  int32 `form:"page,omitempty"`
	Limit int32 `form:"limit,omitempty"`
}

type LookupParams struct {
	Name         string `form:"name"`
	LookupTypeId string `form:"lookup_type_id"`
	Sort         string `form:"sort"`
	Order        string `form:"order"`
	Code         string `form:"code"`
}
