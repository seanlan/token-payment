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

const TableNameChainSendTx = "chain_send_tx"

var ChainSendTxColumns = struct {
	ID                  FieldBase
	ApplicationID       FieldBase
	SerialNo            FieldBase
	ChainSymbol         FieldBase
	ContractAddress     FieldBase
	Symbol              FieldBase
	FromAddress         FieldBase
	ToAddress           FieldBase
	Value               FieldBase
	GasPrice            FieldBase
	TokenID             FieldBase
	TxHash              FieldBase
	Nonce               FieldBase
	Hook                FieldBase
	TransferType        FieldBase
	CreateAt            FieldBase
	TransferAt          FieldBase
	TransferSuccess     FieldBase
	TransferFailedTimes FieldBase
	TransferNextTime    FieldBase
	Received            FieldBase
	ReceiveAt           FieldBase
}{
	ID:                  FieldBase{"id", "id"},
	ApplicationID:       FieldBase{"application_id", "application_id"},
	SerialNo:            FieldBase{"serial_no", "serial_no"},
	ChainSymbol:         FieldBase{"chain_symbol", "chain_symbol"},
	ContractAddress:     FieldBase{"contract_address", "contract_address"},
	Symbol:              FieldBase{"symbol", "symbol"},
	FromAddress:         FieldBase{"from_address", "from_address"},
	ToAddress:           FieldBase{"to_address", "to_address"},
	Value:               FieldBase{"value", "value"},
	GasPrice:            FieldBase{"gas_price", "gas_price"},
	TokenID:             FieldBase{"token_id", "token_id"},
	TxHash:              FieldBase{"tx_hash", "tx_hash"},
	Nonce:               FieldBase{"nonce", "nonce"},
	Hook:                FieldBase{"hook", "hook"},
	TransferType:        FieldBase{"transfer_type", "transfer_type"},
	CreateAt:            FieldBase{"create_at", "create_at"},
	TransferAt:          FieldBase{"transfer_at", "transfer_at"},
	TransferSuccess:     FieldBase{"transfer_success", "transfer_success"},
	TransferFailedTimes: FieldBase{"transfer_failed_times", "transfer_failed_times"},
	TransferNextTime:    FieldBase{"transfer_next_time", "transfer_next_time"},
	Received:            FieldBase{"received", "received"},
	ReceiveAt:           FieldBase{"receive_at", "receive_at"},
}

type ChainSendTx struct {
	ID                  int64   `json:"id" gorm:"column:id;type:bigint;primaryKey;autoIncrement"`                    //
	ApplicationID       int64   `json:"application_id" gorm:"column:application_id;type:bigint;not null"`            //应用ID
	SerialNo            string  `json:"serial_no" gorm:"column:serial_no;type:varchar;not null"`                     //订单序列号
	ChainSymbol         string  `json:"chain_symbol" gorm:"column:chain_symbol;type:varchar;not null"`               //链的符号
	ContractAddress     string  `json:"contract_address" gorm:"column:contract_address;type:varchar;not null"`       //代币合约地址，如果是空表示是主币
	Symbol              string  `json:"symbol" gorm:"column:symbol;type:varchar;not null"`                           //代币符号
	FromAddress         string  `json:"from_address" gorm:"column:from_address;type:varchar;not null"`               //发送地址
	ToAddress           string  `json:"to_address" gorm:"column:to_address;type:varchar;not null"`                   //收款地址
	Value               float64 `json:"value" gorm:"column:value;type:decimal;not null"`                             //数量
	GasPrice            int64   `json:"gas_price" gorm:"column:gas_price;type:bigint;not null"`                      //gas费用
	TokenID             int64   `json:"token_id" gorm:"column:token_id;type:bigint;not null"`                        //tokenid （NFT）
	TxHash              string  `json:"tx_hash" gorm:"column:tx_hash;type:varchar;not null"`                         //交易hash值
	Nonce               int64   `json:"nonce" gorm:"column:nonce;type:bigint;not null;default:-1"`                   //交易nonce
	Hook                string  `json:"hook" gorm:"column:hook;type:varchar;not null"`                               //到账变动通知url
	TransferType        int32   `json:"transfer_type" gorm:"column:transfer_type;type:int;not null"`                 //交易类型 1到账 2提币 3整理费用 4零钱整理
	CreateAt            int64   `json:"create_at" gorm:"column:create_at;type:bigint;not null"`                      //申请时间
	TransferAt          int64   `json:"transfer_at" gorm:"column:transfer_at;type:bigint;not null"`                  //转账时间
	TransferSuccess     int32   `json:"transfer_success" gorm:"column:transfer_success;type:tinyint;not null"`       //是否转账成功
	TransferFailedTimes int32   `json:"transfer_failed_times" gorm:"column:transfer_failed_times;type:int;not null"` //转账失败次数
	TransferNextTime    int64   `json:"transfer_next_time" gorm:"column:transfer_next_time;type:bigint;not null"`    //下次转账时间
	Received            int32   `json:"received" gorm:"column:received;type:tinyint;not null"`                       //是否到账
	ReceiveAt           int64   `json:"receive_at" gorm:"column:receive_at;type:bigint;not null"`                    //到账时间
}

// TableName ChainSendTx's table name
func (*ChainSendTx) TableName() string {
	return TableNameChainSendTx
}