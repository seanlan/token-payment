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

const TableNameAdminUserPermission = "admin_user_permission"

var AdminUserPermissionColumns = struct {
	ID           FieldBase
	UserID       FieldBase
	PermissionID FieldBase
}{
	ID:           FieldBase{"id", "admin_user_permission.id"},
	UserID:       FieldBase{"user_id", "admin_user_permission.user_id"},
	PermissionID: FieldBase{"permission_id", "admin_user_permission.permission_id"},
}

type AdminUserPermission struct {
	ID           int64 `json:"id" gorm:"column:id;type:bigint;primaryKey;autoIncrement"`       //
	UserID       int64 `json:"user_id" gorm:"column:user_id;type:bigint;not null"`             //用户ID
	PermissionID int64 `json:"permission_id" gorm:"column:permission_id;type:bigint;not null"` //权限ID
}

// TableName AdminUserPermission's table name
func (*AdminUserPermission) TableName() string {
	return TableNameAdminUserPermission
}
