# gorm-normalize

> 让gorm使用过程中更加规范, 而不是一个项目里面各种风格都有

1. 区别于直接使用gorm, 使其更加规范
   1. 结构上的规范
   2. 查询上的规范
   3. 事务上的规范

2. 区别于使用gorm gen, 使其更加灵活
   1. 不用生成代码, 直接使用
   2. 减少代码量, 使其更加简洁


## 1. Model

* user.go

```go
package model

type User struct {
	BaseModel
	Name  string  `gorm:"column:name;type:varchar(20);not null;default:'';comment:用户名" json:"name"`
	Files []*File `gorm:"foreignKey:UserID" json:"files"`
}

const (
	userTableName = "users"
)

const (
	// PreloadFiles 预加载关联文件
	PreloadFiles = "Files"
)

func (User) TableName() string {
	return userTableName
}
```

* file.go

```go
package model

type FileType int8

type File struct {
	BaseModel
	UserID   uint     `gorm:"column:user_id;type:int(10) unsigned;not null;default:0;comment:用户ID" json:"user_id"`
	Name     string   `gorm:"column:name;type:varchar(20);not null;default:'';comment:文件名" json:"name"`
	Url      string   `gorm:"column:url;type:varchar(255);not null;default:'';comment:文件地址" json:"url"`
	FileType FileType `gorm:"column:file_type;type:tinyint(1);not null;default:0;comment:文件类型" json:"file_type"`
	Ext      string   `gorm:"column:ext;type:varchar(10);not null;default:'';comment:文件后缀" json:"ext"`
	Size     int64    `gorm:"column:size;type:bigint(20) unsigned;not null;default:0;comment:文件大小" json:"size"`
}

const (
	fileTableName = "files"
)

const (
	// FileTypeImage 图片
	FileTypeImage FileType = iota + 1
	// FileTypeVideo 视频
	FileTypeVideo
	// FileTypeAudio 音频
	FileTypeAudio
	// FileTypeDocument 文档
	FileTypeDocument
	// FileTypeOther 其他
	FileTypeOther
)

func (File) TableName() string {
	return fileTableName
}
```

## 2. Data

* user.go

```go
type (
	// User ...
	User struct {
		query.IAction[model.User]

		PreloadFilesKey string
	}
)

// NewUser ...
func NewUser() *User {
	return &User{
		IAction:         query.NewAction(query.WithDB[model.User](conn.GetMysqlDB())),
		PreloadFilesKey: model.PreloadFiles,
	}
}

// PreloadFiles 预加载关联文件
func (l *User) PreloadFiles(scops ...model.Scopemethod) model.Scopemethod {
	return func(db *gorm.DB) *gorm.DB {
		if len(scops) == 0 {
			return db.Preload(l.PreloadFilesKey)
		}
		// add your code here
		return db.Preload(l.PreloadFilesKey, func(db *gorm.DB) *gorm.DB {
			return db.Scopes(scops...)
		})
	}
}
```

## service

* user.go

```go
type (
	// DetailReq ...
	DetailReq struct {
		// add request params
		ID int `uri:"id"`
	}

	// DetailResp ...
	DetailResp struct {
		// add response params
		ID        uint   `json:"id"`
		Name      string `json:"name"`
		CreatedAt int64  `json:"created_at"`
		UpdateAt  int64  `json:"update_at"`

		Files []*model.File `json:"files"`
	}
)

// GetDetail ...
func (l *User) GetDetail(ctx context.Context, req *DetailReq) (*DetailResp, error) {
	userData := dataUser.NewUser()

	first, err := userData.WithContext(ctx).First(model.WhereID(req.ID), userData.PreloadFiles())
	if err != nil {
		ginplus.Logger().Error("get user detail failed", zap.Any("req", req), zap.Error(err))
		return nil, err
	}

	// add your code here
	return &DetailResp{
		ID:        first.ID,
		Name:      first.Name,
		CreatedAt: first.CreatedAt.Unix(),
		UpdateAt:  first.UpdatedAt.Unix(),
		Files:     first.Files,
	}, nil
}
```