package models

type Query_struct struct{
	Catalog_id uint `json:"catalog_id"`
	Catalog_name string `json:"catalog_name"`
	Catalog_type string `json:"catalog_type"`
	Item_id uint	`json:"item_id"`
	Item_img string	`json:"item_img"`
	Item_link string `json:"item_link"`

}
type Id_s struct{
	Id uint 
}
type Custom_Response struct{
	Catalogs []Custom_Catalog
	Count int
}
type Custom_Catalog struct{
	Name string `json:"name"`
	Id uint `json:"id"`
	Items []Custom_Item `json:"items"`
}
type Custom_Item struct{
	Link string `json:"link"`
	Img string `json:"img"`
}