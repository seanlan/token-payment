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

const TableNameAppChainAddress = "app_chain_address"

var AppChainAddressColumns = struct {
	Address     FieldBase
	AppKey      FieldBase
	ChainSymbol FieldBase
	CreateAt    FieldBase
	EncKey      FieldBase
	Hook        FieldBase
	ID          FieldBase
}{
	Address:     FieldBase{"`address`"},
	AppKey:      FieldBase{"`app_key`"},
	ChainSymbol: FieldBase{"`chain_symbol`"},
	CreateAt:    FieldBase{"`create_at`"},
	EncKey:      FieldBase{"`enc_key`"},
	Hook:        FieldBase{"`hook`"},
	ID:          FieldBase{"`id`"},
}

type AppChainAddress struct {
	Address     string `json:"address" gorm:"column:address;type:varchar(255);not null"`               //地址
	AppKey      string `json:"app_key" gorm:"column:app_key;type:varchar(200);not null"`               //app key
	ChainSymbol string `json:"chain_symbol" gorm:"column:chain_symbol;type:varchar(200);not null"`     //链的符号
	CreateAt    int64  `json:"create_at" gorm:"column:create_at;type:bigint;not null"`                 //创建时间
	EncKey      string `json:"enc_key" gorm:"column:enc_key;type:varchar(255);not null"`               //加密后的私钥
	Hook        string `json:"hook" gorm:"column:hook;type:varchar(512);not null"`                     //账号变动通知url
	ID          uint64 `json:"id" gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement:true"` //
}

// TableName AppChainAddress's table name
func (*AppChainAddress) TableName() string {
	return TableNameAppChainAddress
}