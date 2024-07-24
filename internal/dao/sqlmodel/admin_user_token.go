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

const TableNameAdminUserToken = "admin_user_token"

var AdminUserTokenColumns = struct {
	ID              FieldBase
	UserID          FieldBase
	Token           FieldBase
	ExpireAt        FieldBase
	RefreshToken    FieldBase
	RefreshExpireAt FieldBase
	CreateAt        FieldBase
}{
	ID:              FieldBase{"id", "id"},
	UserID:          FieldBase{"user_id", "user_id"},
	Token:           FieldBase{"token", "token"},
	ExpireAt:        FieldBase{"expire_at", "expire_at"},
	RefreshToken:    FieldBase{"refresh_token", "refresh_token"},
	RefreshExpireAt: FieldBase{"refresh_expire_at", "refresh_expire_at"},
	CreateAt:        FieldBase{"create_at", "create_at"},
}

type AdminUserToken struct {
	ID              int64  `json:"id" gorm:"column:id;type:bigint;primaryKey;autoIncrement"`                         //
	UserID          int64  `json:"user_id" gorm:"column:user_id;type:bigint;not null"`                               //用户id
	Token           string `json:"token" gorm:"column:token;type:varchar;not null"`                                  //token
	ExpireAt        int64  `json:"expire_at" gorm:"column:expire_at;type:bigint;not null;default:0"`                 //token过期时间
	RefreshToken    string `json:"refresh_token" gorm:"column:refresh_token;type:varchar;not null"`                  //refresh token
	RefreshExpireAt int64  `json:"refresh_expire_at" gorm:"column:refresh_expire_at;type:bigint;not null;default:0"` //refresh token过期时间
	CreateAt        int64  `json:"create_at" gorm:"column:create_at;type:bigint;not null;default:0"`                 //创建时间
}

// TableName AdminUserToken's table name
func (*AdminUserToken) TableName() string {
	return TableNameAdminUserToken
}
