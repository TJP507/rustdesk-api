package model

type Tag struct {
	IdModel
	Name         string                 `json:"name" gorm:"default:'';not null;"`
	UserId       uint                   `json:"user_id" gorm:"default:0;not null;index"`
	Color        uint                   `json:"color" gorm:"default:0;not null;"` // color is a Flutter color value ranging from 0x00000000 to 0xFFFFFFFF; the first two hex digits represent opacity and the remaining six represent the color, convertible to RGBA
	CollectionId uint                   `json:"collection_id" gorm:"default:0;not null;index"`
	Collection   *AddressBookCollection `json:"collection,omitempty"`
	TimeModel
}

type TagList struct {
	Tags []*Tag `json:"list"`
	Pagination
}
