package mysql

import (
	"demo1/internal/app/mydemo/model"

	"gorm.io/gorm"
)

// 根据uid，cid查询表
func GetUserCheckinJoin(uid int, cid int) (userCheckinJoin *model.UserCheckinJoin, err error) {
	err = DB.Where("uid=? and cid =? and status=?", uid, cid, model.JoinStatusNormal).First(&userCheckinJoin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound { //没查到数据返回空
			return nil, nil
		}

		return nil, err
	}

	return userCheckinJoin, nil
}

// 根据uid查询表
func GetUserCheckinJoinByuid(uid int) (userCheckinJoinByuid []model.UserCheckinJoin, err error) {
	err = DB.Where("uid=? and status=?", uid, model.JoinStatusNormal).Find(&userCheckinJoinByuid).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound { //没查到数据返回空
			return nil, nil
		}

		return nil, err
	}

	return userCheckinJoinByuid, nil
}

// 根据uid，cid查询表
func GetUserCheckinJoinCount(uid int, cid int) (userCheckinJoin *model.UserCheckinJoin, err error) {
	err = DB.Where("uid=? and cid =? and status=?", uid, cid, model.JoinStatusNormal).First(&userCheckinJoin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound { //没查到数据返回空
			return nil, nil
		}

		return nil, err
	}

	return userCheckinJoin, nil
}

// 添加表
func AddUserCheckinJoin(newUserCheckinJoin *model.UserCheckinJoin) (err error) {
	err = DB.Create(&newUserCheckinJoin).Error
	return
}
