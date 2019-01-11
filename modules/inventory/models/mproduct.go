package models
type (
	Mproduct struct {
		M_product_id	string		`gorm:"type:varchar(64);PRIMARY_KEY"`
		Name			string		`gorm:"type:varchar(255);"`
	}
    // transformedmProductrepresents a formatted product
	TransformedMproduct struct {
		M_product_id	string		`json:"sku"`
		Name			string		`json:"product_name"`
	}
)

//change table name
func (Mproduct) TableName() string {
    return "m_product"
}
