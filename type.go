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
	Level    RegionType `json:"level"`
	Name     string     `json:"name"`
	Code     string     `json:"code"`
	Type     string     `json:"type,omitempty"`
	Url      string     `json:"-"`
	Children []*Region  `json:"children"`
}

func NewRegion(t RegionType, name string, code string, class string, url string) *Region {
	return &Region{Level: t, Name: name, Code: code, Type: class, Url: url}
}

func (r *Region) add(region *Region) {
	r.Children = append(r.Children, region)
}
