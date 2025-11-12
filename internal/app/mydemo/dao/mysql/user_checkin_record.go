package mysql

import (
	"demo1/internal/app/mydemo/model"
	"time"

	"gorm.io/gorm"
)

// 查询有没有打卡
func GetUserCheckinRecord(uid int, cid int, date int) (userCheckinRecord *model.UserCheckinRecord, err error) {
	err = DB.Where("uid = ? AND cid = ? AND date >= ?", uid, cid, date).First(&userCheckinRecord).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound { //没查到数据返回空
			return nil, nil
		}

		return nil, err
	}

	return userCheckinRecord, nil
}

// 根据uid查询
func GetUserCheckinRecordByUidDate(uid int, date int) (userCheckinRecordByuid []model.UserCheckinRecord, err error) {
	err = DB.Where("uid=? and date=? and status=?", uid, date, model.JoinStatusNormal).Find(&userCheckinRecordByuid).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound { //没查到数据返回空
			return nil, nil
		}

		return nil, err
	}

	return userCheckinRecordByuid, nil
}

// in 查询
func GetUserCheckinRecordInCheckinId(uid int, date int, cidSlice []int) (userCheckinRecordByuid []model.UserCheckinRecord, err error) {
	err = DB.Where("uid=? and date=? and status=? and cid IN (?)", uid, date, model.JoinStatusNormal, cidSlice).Find(&userCheckinRecordByuid).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound { //没查到数据返回空
			return nil, nil
		}

		return nil, err
	}

	return
}

// 添加打卡
func AddUserCheckinRecord(userCheckinRecord *model.UserCheckinRecord) (err error) {
	err = DB.Create(&userCheckinRecord).Error
	return
}

// 获取今天的用户打卡的时间
func GetCreateTimeFromUserCheckinRecord(uid int, cid int, date int) (createTime time.Time, err error) {
	err = DB.Model(&model.UserCheckinRecord{}).Where("uid=? and cid=? and date=?", uid, cid, date).Select("create_at").Scan(&createTime).Error
	return
}

// 获取用户该打卡在今天打卡的名次
func GetUserCheckinRecordByCount(cid int, date int, createTime time.Time) (rank int, err error) {
	var count int64
	err = DB.Model(&model.UserCheckinRecord{}).Where("cid=? and date=? and create_at<?", cid, date, createTime).Count(&count).Error
	if err != nil {
		return
	}

	rank = int(count) + 1
	return
}

// 打卡列表
func GetUserCheckinRecordList(uid int, cid int, isasc bool) (userCheckinRecordList []model.UserCheckinRecord, err error) {
	byOrder := "date desc"
	if isasc {
		byOrder = "date asc"
	}

	err = DB.Order(byOrder).Where("uid = ? AND cid = ?", uid, cid).Find(&userCheckinRecordList).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound { //没查到数据返回空
			return nil, nil
		}

		return nil, err
	}

	return userCheckinRecordList, nil
}

// 更新打卡
func UpdateUserCheckinRecord(uid int, cid int, date int) (err error) {
	now := time.Now()
	err = DB.Model(&model.UserCheckinRecord{}).Where("uid = ? AND cid = ? AND date = ?", uid, cid, date).Update("update_at", now).Error
	return
}
