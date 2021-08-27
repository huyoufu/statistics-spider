package main

type RegionType int32

const (
	RegionType_Country RegionType = 0
	//RegionType_Province 省份
	RegionType_Province RegionType = 1
	//RegionType_City 市
	RegionType_City RegionType = 2
	//RegionType_County 县
	RegionType_County RegionType = 3
	//RegionType_Town 乡
	RegionType_Town RegionType = 4
	//RegionType_Village 村
	RegionType_Village RegionType = 5
)

type Region struct {
	t        RegionType
	Name     string `json:"name"`
	Code     string `json:"code"`
	Url      string `json:"url"`
	Children []*Region `json:"children"`
}

func NewRegion(t RegionType, name string, code string, url string) *Region {
	return &Region{t: t, Name: name, Code: code, Url: url}
}




func (r *Region) add(region *Region) {
	r.Children = append(r.Children, region)
}
