//generated by lazy
//author: seanlan
/**
//表结构
type BaseStruct struct {
	Package    string    //包名
	StructName string    //结构名
	TableName  string    //表名
	Members    []*Member //成员
}

// 字段结构
type Member struct {
	Name          string    //字段名
	Type          string    //字段类型
	NewType       string    //新字段类型
	ColumnName    string    //字段列名
	ColumnComment string    //字段注释
	ModelType     string    //字段模型类型
	JSONTag       string    //json tag
	GORMTag       string    //gorm tag
	NewTag        string    //新tag
}
*/

package sqlmodel

const TableNameChainToken = "chain_token"

var ChainTokenColumns = struct {
	ID              FieldBase
	ChainSymbol     FieldBase
	ContractAddress FieldBase
	Name            FieldBase
	Symbol          FieldBase
	Decimals        FieldBase
	Threshold       FieldBase
	GasFee          FieldBase
	ArrangeSwitch   FieldBase
}{
	ID:              FieldBase{"id", "chain_token.id"},
	ChainSymbol:     FieldBase{"chain_symbol", "chain_token.chain_symbol"},
	ContractAddress: FieldBase{"contract_address", "chain_token.contract_address"},
	Name:            FieldBase{"name", "chain_token.name"},
	Symbol:          FieldBase{"symbol", "chain_token.symbol"},
	Decimals:        FieldBase{"decimals", "chain_token.decimals"},
	Threshold:       FieldBase{"threshold", "chain_token.threshold"},
	GasFee:          FieldBase{"gas_fee", "chain_token.gas_fee"},
	ArrangeSwitch:   FieldBase{"arrange_switch", "chain_token.arrange_switch"},
}

type ChainToken struct {
	ID              int64   `json:"id" gorm:"column:id;type:bigint;primaryKey;autoIncrement"`              //
	ChainSymbol     string  `json:"chain_symbol" gorm:"column:chain_symbol;type:varchar;not null"`         //链的符号
	ContractAddress string  `json:"contract_address" gorm:"column:contract_address;type:varchar;not null"` //代币合约地址，如果是空表示是主币
	Name            string  `json:"name" gorm:"column:name;type:varchar;not null"`                         //币种名称
	Symbol          string  `json:"symbol" gorm:"column:symbol;type:varchar;not null"`                     //币种符号
	Decimals        int32   `json:"decimals" gorm:"column:decimals;type:int;not null"`                     //小数位
	Threshold       float64 `json:"threshold" gorm:"column:threshold;type:decimal;not null"`               //零钱整理阀值
	GasFee          float64 `json:"gas_fee" gorm:"column:gas_fee;type:decimal;not null"`                   //Gas费用
	ArrangeSwitch   int32   `json:"arrange_switch" gorm:"column:arrange_switch;type:tinyint;not null"`     //零钱整理开关
}

// TableName ChainToken's table name
func (*ChainToken) TableName() string {
	return TableNameChainToken
}
