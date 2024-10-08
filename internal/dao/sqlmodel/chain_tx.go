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

const TableNameChainTx = "chain_tx"

var ChainTxColumns = struct {
	ID                FieldBase
	ApplicationID     FieldBase
	ChainSymbol       FieldBase
	BlockNumber       FieldBase
	BlockHash         FieldBase
	TxHash            FieldBase
	FromAddress       FieldBase
	ToAddress         FieldBase
	ContractAddress   FieldBase
	Symbol            FieldBase
	Value             FieldBase
	TokenID           FieldBase
	TxIndex           FieldBase
	BatchIndex        FieldBase
	Confirm           FieldBase
	Confirmed         FieldBase
	Removed           FieldBase
	TransferType      FieldBase
	Arranged          FieldBase
	CreateAt          FieldBase
	SerialNo          FieldBase
	NotifySuccess     FieldBase
	NotifyFailedTimes FieldBase
	NotifyNextTime    FieldBase
}{
	ID:                FieldBase{"id", "chain_tx.id"},
	ApplicationID:     FieldBase{"application_id", "chain_tx.application_id"},
	ChainSymbol:       FieldBase{"chain_symbol", "chain_tx.chain_symbol"},
	BlockNumber:       FieldBase{"block_number", "chain_tx.block_number"},
	BlockHash:         FieldBase{"block_hash", "chain_tx.block_hash"},
	TxHash:            FieldBase{"tx_hash", "chain_tx.tx_hash"},
	FromAddress:       FieldBase{"from_address", "chain_tx.from_address"},
	ToAddress:         FieldBase{"to_address", "chain_tx.to_address"},
	ContractAddress:   FieldBase{"contract_address", "chain_tx.contract_address"},
	Symbol:            FieldBase{"symbol", "chain_tx.symbol"},
	Value:             FieldBase{"value", "chain_tx.value"},
	TokenID:           FieldBase{"token_id", "chain_tx.token_id"},
	TxIndex:           FieldBase{"tx_index", "chain_tx.tx_index"},
	BatchIndex:        FieldBase{"batch_index", "chain_tx.batch_index"},
	Confirm:           FieldBase{"confirm", "chain_tx.confirm"},
	Confirmed:         FieldBase{"confirmed", "chain_tx.confirmed"},
	Removed:           FieldBase{"removed", "chain_tx.removed"},
	TransferType:      FieldBase{"transfer_type", "chain_tx.transfer_type"},
	Arranged:          FieldBase{"arranged", "chain_tx.arranged"},
	CreateAt:          FieldBase{"create_at", "chain_tx.create_at"},
	SerialNo:          FieldBase{"serial_no", "chain_tx.serial_no"},
	NotifySuccess:     FieldBase{"notify_success", "chain_tx.notify_success"},
	NotifyFailedTimes: FieldBase{"notify_failed_times", "chain_tx.notify_failed_times"},
	NotifyNextTime:    FieldBase{"notify_next_time", "chain_tx.notify_next_time"},
}

type ChainTx struct {
	ID                int64   `json:"id" gorm:"column:id;type:bigint;primaryKey;autoIncrement"`                //
	ApplicationID     int64   `json:"application_id" gorm:"column:application_id;type:bigint;not null"`        //应用ID
	ChainSymbol       string  `json:"chain_symbol" gorm:"column:chain_symbol;type:varchar;not null"`           //链的符号
	BlockNumber       int64   `json:"block_number" gorm:"column:block_number;type:bigint;not null"`            //区块高度
	BlockHash         string  `json:"block_hash" gorm:"column:block_hash;type:varchar;not null"`               //区块hash值
	TxHash            string  `json:"tx_hash" gorm:"column:tx_hash;type:varchar;not null"`                     //交易hash值
	FromAddress       string  `json:"from_address" gorm:"column:from_address;type:varchar;not null"`           //支付地址
	ToAddress         string  `json:"to_address" gorm:"column:to_address;type:varchar;not null"`               //收款地址
	ContractAddress   string  `json:"contract_address" gorm:"column:contract_address;type:varchar;not null"`   //代币合约地址，如果是空表示是主币
	Symbol            string  `json:"symbol" gorm:"column:symbol;type:varchar;not null"`                       //代币符号
	Value             float64 `json:"value" gorm:"column:value;type:decimal;not null"`                         //数量
	TokenID           int64   `json:"token_id" gorm:"column:token_id;type:bigint;not null"`                    //tokenid （NFT）
	TxIndex           int64   `json:"tx_index" gorm:"column:tx_index;type:bigint;not null"`                    //交易序号
	BatchIndex        int64   `json:"batch_index" gorm:"column:batch_index;type:bigint;not null"`              //交易批次号
	Confirm           int32   `json:"confirm" gorm:"column:confirm;type:int;not null"`                         //确认次数
	Confirmed         int32   `json:"confirmed" gorm:"column:confirmed;type:int;not null"`                     //确认过的
	Removed           int32   `json:"removed" gorm:"column:removed;type:int;not null"`                         //是否已移除
	TransferType      int32   `json:"transfer_type" gorm:"column:transfer_type;type:int;not null"`             //交易类型 1到账 2提币
	Arranged          int32   `json:"arranged" gorm:"column:arranged;type:int;not null"`                       //是否整理过
	CreateAt          int64   `json:"create_at" gorm:"column:create_at;type:bigint;not null"`                  //交易时间
	SerialNo          string  `json:"serial_no" gorm:"column:serial_no;type:varchar;not null"`                 //订单序列号
	NotifySuccess     int32   `json:"notify_success" gorm:"column:notify_success;type:int;not null"`           //是否通知成功
	NotifyFailedTimes int32   `json:"notify_failed_times" gorm:"column:notify_failed_times;type:int;not null"` //通知失败次数
	NotifyNextTime    int64   `json:"notify_next_time" gorm:"column:notify_next_time;type:bigint;not null"`    //下次通知时间
}

// TableName ChainTx's table name
func (*ChainTx) TableName() string {
	return TableNameChainTx
}
