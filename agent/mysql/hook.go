package mysql

import (
	"gorm.io/gorm"
	"time"
)

// BaseModel Define a base model that includes the fields ctime and mtime
type BaseModel struct {
	ID    uint64 `gorm:"primary_key"`
	Ctime int64
	Mtime int64
}

// BeforeCreate hook sets the ctime and mtime for new records
func (base *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	currentTime := time.Now().UnixNano() / int64(time.Millisecond)
	base.Ctime = currentTime
	base.Mtime = currentTime
	return nil
}

// BeforeUpdate hook sets the mtime for updated records
func (base *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	currentTime := time.Now().UnixNano() / int64(time.Millisecond)
	base.Mtime = currentTime
	return nil
}

// RegisterGlobalHooks Register global hooks for all models
func RegisterGlobalHooks(db *gorm.DB) {
	if err := db.Callback().Create().Before("gorm:before_create").Register("set_create_time", setCreateTime); err != nil {
		panic(err)
	}
	if err := db.Callback().Update().Before("gorm:before_update").Register("set_update_time", setUpdateTime); err != nil {
		panic(err)
	}
}

// setCreateTime sets the creation time for new records
func setCreateTime(db *gorm.DB) {
	if _, ok := db.Statement.Model.(*BaseModel); ok {
		currentTime := time.Now().UnixMilli()
		db.Statement.SetColumn("Ctime", currentTime)
		db.Statement.SetColumn("Mtime", currentTime)
	}
}

// setUpdateTime sets the update time for updated records
func setUpdateTime(db *gorm.DB) {
	if _, ok := db.Statement.Model.(*BaseModel); ok {
		currentTime := time.Now().UnixMilli()
		db.Statement.SetColumn("Mtime", currentTime)
	}
}
