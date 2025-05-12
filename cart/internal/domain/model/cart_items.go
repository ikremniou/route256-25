package model

type EnrichedCartItemModel struct {
	SkuId int64
	Name  string
	Count uint32
	Price uint32
}

type AllCartItemsModel struct {
	Items []EnrichedCartItemModel
	Total uint32
}
