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

const TableNameAdminLog = "admin_log"

var AdminLogColumns = struct {
	ID        FieldBase
	UserID    FieldBase
	AssetID   FieldBase
	AssetType FieldBase
	Remarks   FieldBase
	CreateAt  FieldBase
}{
	ID:        FieldBase{"id", "admin_log.id"},
	UserID:    FieldBase{"user_id", "admin_log.user_id"},
	AssetID:   FieldBase{"asset_id", "admin_log.asset_id"},
	AssetType: FieldBase{"asset_type", "admin_log.asset_type"},
	Remarks:   FieldBase{"remarks", "admin_log.remarks"},
	CreateAt:  FieldBase{"create_at", "admin_log.create_at"},
}

type AdminLog struct {
	ID        int64  `json:"id" gorm:"column:id;type:bigint;primaryKey;autoIncrement"` //
	UserID    int64  `json:"user_id" gorm:"column:user_id;type:bigint;not null"`       //用户ID
	AssetID   int64  `json:"asset_id" gorm:"column:asset_id;type:bigint;not null"`     //资源ID
	AssetType int32  `json:"asset_type" gorm:"column:asset_type;type:int;not null"`    //资源类型
	Remarks   string `json:"remarks" gorm:"column:remarks;type:varchar;not null"`      //日志
	CreateAt  int64  `json:"create_at" gorm:"column:create_at;type:bigint;not null"`   //记录时间
}

// TableName AdminLog's table name
func (*AdminLog) TableName() string {
	return TableNameAdminLog
}
