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

const TableNameChain = "chain"

var ChainColumns = struct {
	ID          FieldBase
	ChainSymbol FieldBase
	Name        FieldBase
	ChainID     FieldBase
	Currency    FieldBase
	ChainType   FieldBase
	Confirm     FieldBase
	Gas         FieldBase
	GasPrice    FieldBase
	LatestBlock FieldBase
	RebaseBlock FieldBase
	Concurrent  FieldBase
}{
	ID:          FieldBase{"`id`"},
	ChainSymbol: FieldBase{"`chain_symbol`"},
	Name:        FieldBase{"`name`"},
	ChainID:     FieldBase{"`chain_id`"},
	Currency:    FieldBase{"`currency`"},
	ChainType:   FieldBase{"`chain_type`"},
	Confirm:     FieldBase{"`confirm`"},
	Gas:         FieldBase{"`gas`"},
	GasPrice:    FieldBase{"`gas_price`"},
	LatestBlock: FieldBase{"`latest_block`"},
	RebaseBlock: FieldBase{"`rebase_block`"},
	Concurrent:  FieldBase{"`concurrent`"},
}

type Chain struct {
	ID          int64  `json:"id" gorm:"column:id;type:bigint;primaryKey;autoIncrement"`              //
	ChainSymbol string `json:"chain_symbol" gorm:"column:chain_symbol;type:varchar;not null"`         //链的符号
	Name        string `json:"name" gorm:"column:name;type:varchar;not null"`                         //链名称
	ChainID     int64  `json:"chain_id" gorm:"column:chain_id;type:bigint;not null"`                  //链ID
	Currency    string `json:"currency" gorm:"column:currency;type:varchar;not null"`                 //货币
	ChainType   string `json:"chain_type" gorm:"column:chain_type;type:varchar;not null;default:evm"` //链类型 默认evm
	Confirm     int32  `json:"confirm" gorm:"column:confirm;type:int;not null"`                       //确认区块数量
	Gas         int64  `json:"gas" gorm:"column:gas;type:bigint;not null"`                            //gas费用配置
	GasPrice    int64  `json:"gas_price" gorm:"column:gas_price;type:bigint;not null"`                //gas price 配置
	LatestBlock int64  `json:"latest_block" gorm:"column:latest_block;type:bigint;not null"`          //最新区块
	RebaseBlock int64  `json:"rebase_block" gorm:"column:rebase_block;type:bigint;not null"`          //重新构建区块
	Concurrent  int32  `json:"concurrent" gorm:"column:concurrent;type:int;not null"`                 //并发量
}

// TableName Chain's table name
func (*Chain) TableName() string {
	return TableNameChain
}
