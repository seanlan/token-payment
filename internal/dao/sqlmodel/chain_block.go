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

const TableNameChainBlock = "chain_block"

var ChainBlockColumns = struct {
	ID          FieldBase
	ChainSymbol FieldBase
	BlockNumber FieldBase
	BlockHash   FieldBase
	ParentHash  FieldBase
	Checked     FieldBase
}{
	ID:          FieldBase{"id", "chain_block.id"},
	ChainSymbol: FieldBase{"chain_symbol", "chain_block.chain_symbol"},
	BlockNumber: FieldBase{"block_number", "chain_block.block_number"},
	BlockHash:   FieldBase{"block_hash", "chain_block.block_hash"},
	ParentHash:  FieldBase{"parent_hash", "chain_block.parent_hash"},
	Checked:     FieldBase{"checked", "chain_block.checked"},
}

type ChainBlock struct {
	ID          int64  `json:"id" gorm:"column:id;type:bigint;primaryKey;autoIncrement"`      //
	ChainSymbol string `json:"chain_symbol" gorm:"column:chain_symbol;type:varchar;not null"` //链的符号
	BlockNumber int64  `json:"block_number" gorm:"column:block_number;type:bigint;not null"`  //区块高度
	BlockHash   string `json:"block_hash" gorm:"column:block_hash;type:varchar;not null"`     //区块hash值
	ParentHash  string `json:"parent_hash" gorm:"column:parent_hash;type:varchar;not null"`   //上一个区块hash值
	Checked     int32  `json:"checked" gorm:"column:checked;type:int;not null"`               //是否检测完成
}

// TableName ChainBlock's table name
func (*ChainBlock) TableName() string {
	return TableNameChainBlock
}
