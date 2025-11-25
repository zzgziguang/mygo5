package mysql

import (
	"demo1/internal/app/mydemo/model"

	"gorm.io/gorm"
)

func CreateCheckin(newcheckin *model.Checkin) (result *gorm.DB) {
	result = DB.Create(&newcheckin)
	return
}

// 获取全部打卡
func GetCheckinAll(date int) (checkins []model.Checkin, err error) {
	err = DB.Model(&model.Checkin{}).Where("end_time >= ?", date).Find(&checkins).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound { //没查到数据返回空
			return nil, nil
		}

		return nil, err
	}

	return checkins, nil
}

// 查询结束时间
func GetCheckinEndTimeByCid(cid int) (endTime int, err error) {
	err = DB.Select("end_time").Model(&model.Checkin{}).Where("id = ?", cid).Scan(&endTime).Error
	return
}

// 根据id排序
func GetCheckinOrderId(page int, pagesize int, isasc bool) (checkins []model.Checkin, err error) {
	offset := (page - 1) * pagesize
	var order string
	if isasc {
		order = "id asc"
	} else {
		order = "id desc"
	}

	err = DB.Order(order).Offset(offset).Limit(pagesize).Find(&checkins).Error
	return
}

// 查询获取全部id
func GetCheckinId() (checkinIds []int, err error) {
	err = DB.Model(&model.Checkin{}).Select("id").Find(&checkinIds).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound { //没查到数据返回空
			return nil, nil
		}

		return nil, err
	}

	return checkinIds, nil
}

// 根据参与人数排序
func GetCheckinOrderJoinnumber(page int, pagesize int, isasc bool) (checkins []model.Checkin, err error) {
	offset := (page - 1) * pagesize
	var order string
	if isasc {
		order = "join_num asc"
	} else {
		order = "join_num desc"
	}

	err = DB.Order(order).Offset(offset).Limit(pagesize).Find(&checkins).Error
	return
}

// 更新joinnumber
func UpdateCheckinJoinNum(cid int, checkin *model.Checkin) (err error) {
	err = DB.Model(&model.Checkin{}).Where("id = ?", cid).Update("join_num", checkin.JoinNum+1).Error
	return
}

// 更新weight
func UpdateCheckinWeight(newCheckin *model.Checkin, weight int) (err error) {
	err = DB.Model(&model.Checkin{}).Where("id = ?", newCheckin.Id).Update("weight = ?", weight).Error
	return
}

// 根据cid获取
func GetCheckinBycid(cid int) (checkin *model.Checkin, err error) {
	err = DB.Model(&model.Checkin{}).Where("id=?", cid).First(&checkin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound { //没查到数据返回空
			return nil, nil
		}

		return nil, err
	}

	return checkin, nil
}
