package controller

import (
	"demo1/internal/app/mydemo/model"
	"demo1/internal/app/mydemo/service"
	"demo1/internal/app/mydemo/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 通过post查询参数添加用户的处理函数
func AddUserHandler(c *gin.Context) { //c

	// 从表单中获取用户信息
	username := c.PostForm("username") //获取username，存到变量username
	email := c.PostForm("email")
	ageStr := c.PostForm("age")
	phoneStr := c.PostForm("phone")

	// 数据验证
	if username == "" {
		MakeApiResponseError(c, CODE_PARAMS_ERROR)
		return
	}

	if email == "" {
		MakeApiResponseError(c, CODE_PARAMS_ERROR)
		return
	}

	if !utils.IsValidEmail(email) {
		MakeApiResponseError(c, CODE_PARAMS_ERROR)
		return
	}

	phone, err := strconv.Atoi(phoneStr)
	if err != nil || phone < 0 {
		service.Logger.Error("phoneStrAtoi err", zap.Error(err))
		MakeApiResponseError(c, CODE_PARAMS_ERROR)
		return
	}

	if phone == 0 {
		MakeApiResponseError(c, CODE_PARAMS_ERROR)
		return
	}

	// 解析年龄参数
	age, err := strconv.Atoi(ageStr)
	if err != nil || age < 0 || age > 150 {
		service.Logger.Error("ageStrAtoi err", zap.Error(err))
		MakeApiResponseError(c, CODE_PARAMS_ERROR)
		return
	}

	resp, err := http.Get("http://v.juhe.cn/historyWeather/weather?weather=sunny")
	if err != nil {
		service.Logger.Error("http.Get err", zap.Error(err))
		return
	}

	//在函数结束前关闭resp结构体的Body属性
	defer resp.Body.Close()
	//读resp的体Body，返回[]byte和错误
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		service.Logger.Error("io.ReadAll err", zap.Error(err))
		return
	}

	//将[]byte类型转换为string
	str := string(body)
	str = `{"code":0,"data":{"tianqi":"晴"}}`

	type wdata struct {
		Tianqi string `json:"tianqi"`
	}

	type Weather struct {
		Code int   `json:"code"`
		Data wdata `json:"data"`
	}

	bytestr := []byte(str)
	var w Weather

	err = json.Unmarshal(bytestr, &w)
	if err != nil {
		service.Logger.Error("wUnmarshal err", zap.Error(err))
		return
	}

	tianq := w.Data.Tianqi
	fmt.Println(tianq)

	// 构造用户对象
	newUser := &model.User{ //其中包含自动生成的id
		Username: username,
		Email:    email,
		Age:      age,
		Phone:    phone,
		//Tianqi:  tianq,
	}

	// 插入数据库
	result := service.CreateUser(newUser)
	if result.Error != nil {
		service.Logger.Error("CreateUser err", zap.Error(result.Error))
		MakeApiResponseErrorDefault(c)
		return
	}

	// 返回成功响应
	MakeApiResponseSuccess(c, CODE_SUCCESS)
}

// 根据id查询用户信息
func GetUserHandler(c *gin.Context) { //
	// 从查询参数中获取用户id
	idStr := c.Query("id") //获取id，存到变量idStr
	if idStr == "" {
		MakeApiResponseError(c, CODE_PARAMS_ERROR)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		service.Logger.Error("idAtoi err", zap.Error(err))
		MakeApiResponseError(c, CODE_PARAMS_ERROR)
		return
	}

	//查 Redis
	user, cacheErr := service.GetUserFromCache(id)
	if cacheErr != nil {
		service.Logger.Error("GetUserFromCache err", zap.Error(cacheErr))
		MakeApiResponseErrorDefault(c)
		return
	} else {

	}
	if user != nil {
		MakeApiResponseSuccess(c, user)
		return
	} else {
		//redis没有，查数据库
		var result *gorm.DB
		result, user = service.GerUserById(id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				service.Logger.Error("GerUserById err", zap.Error(result.Error))
				MakeApiResponseErrorDefault(c)
				return
			} else {
				service.Logger.Error("GerUserById err", zap.Error(result.Error))
				MakeApiResponseErrorDefault(c)
				return
			}

		}

		//写入redis
		err = service.SetUserToCache(user)
		if err != nil {
			service.Logger.Error("SetUserToCache err", zap.Error(err))
			MakeApiResponseErrorDefault(c)
			return
		}

		//响应成功
		MakeApiResponseSuccess(c, user)
	}
}

// 更新用户名
func UpdateUserHandler(c *gin.Context) {

	idStr := c.PostForm("id")
	username := c.PostForm("username")

	if idStr == "" || username == "" {
		MakeApiResponseError(c, CODE_PARAMS_ERROR)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		service.Logger.Error("idAtoi err", zap.Error(err))
		MakeApiResponseError(c, CODE_PARAMS_ERROR)
		return
	}

	// 检查用户是否存在
	result, _ := service.GerUserById(id)
	if result.Error == gorm.ErrRecordNotFound {
		service.Logger.Error("GerUserById err", zap.Error(result.Error))
		MakeApiResponseErrorDefault(c)
		return
	}

	// 更新
	result = service.UpdateUserById(id, username)
	if result.Error != nil {
		service.Logger.Error("UpdateUserById err", zap.Error(err))
		MakeApiResponseErrorDefault(c)
		return
	}

	//更新后删除 Redis 缓存，下次查询会重建
	err = service.DelRedisUser(idStr)
	if err != nil {
		service.Logger.Error("DelRedisUser err", zap.Error(err))
		MakeApiResponseErrorDefault(c)
		return
	}

	MakeApiResponseSuccess(c, map[string]interface{}{
		"id":       id,
		"username": username,
	})
}

// 获取所有用户并按年龄排序
func GetUsersHandlerAll(c *gin.Context) {
	order := c.Query("order")
	if order != "desc" && order != "asc" {
		MakeApiResponseError(c, CODE_PARAMS_ERROR)
		return
	}

	pagestr := c.Query("page")
	page, err := strconv.Atoi(pagestr)
	if err != nil {
		service.Logger.Error("Atoipage err", zap.Error(err))
		MakeApiResponseError(c, CODE_PARAMS_ERROR)
		return
	}

	pagesize := 3
	isasc := true
	if order == "desc" {
		isasc = false
	}

	//先在缓存查用户个数
	total, err := service.GetRedisUserCount()
	if err != nil {
		service.Logger.Error("GetRedisUserCount err", zap.Error(err))
		MakeApiResponseErrorDefault(c)
		return
	}

	if total != 0 {
		service.Logger.Debug("GetRedisUserCount", zap.String("total", strconv.Itoa(int(total))))
	}

	//先获取个数
	total, err = service.GetUserCount()
	if err != nil {
		service.Logger.Error("GetUserCount err", zap.Error(err))
		MakeApiResponseErrorDefault(c)
		return
	}

	//在缓存存count
	err = service.SetRedisUserCount(total)
	if err != nil {
		service.Logger.Error("SetRedisUserCount err", zap.Error(err))
		MakeApiResponseErrorDefault(c)
		return
	} else {
		service.Logger.Debug("SetRedisUserCount", zap.String("total", strconv.Itoa(int(total))))
	}

	var hasNext bool
	if int(total)/pagesize > page {
		hasNext = true
	} else {
		hasNext = false
	}

	//查询redis列表
	//查出的数据要为slice，如果没有数据，应该存"[]"
	//users := []model.User{}
	users, err := service.GetRedisUserSlice(order, page, pagesize)
	if err != nil {
		service.Logger.Error("GetRedisUserSlice err", zap.Error(err))
		MakeApiResponseErrorDefault(c)
		return
	} else {

	}

	if users != nil {
		MakeApiResponseSuccess(c, map[string]interface{}{
			"strSlice": users,
			"hasNext":  hasNext,
		})
		return
	} else {

		users, err = service.GetUserByPage(page, pagesize, isasc)
		// page=100时，users，err怎么样，user cap 变为20，值[]
		if err != nil {
			// 记录错误日志
			service.Logger.Error("GetUserByPage err", zap.Error(err))
			MakeApiResponseErrorDefault(c)
			return
		}

		//添加缓存
		//在没有数据也要存缓存，page=100页，存的是[]
		err = service.SetRedisUserSlice(users, order, page, pagesize)
		if err != nil {
			service.Logger.Error("SetRedisUserSlice err", zap.Error(err))
			MakeApiResponseErrorDefault(c)
			return
		}

		MakeApiResponseSuccess(c, map[string]interface{}{
			"strSlice": users,   //缓存
			"hasNext":  hasNext, //判断是否有下一页
		})
	}
}
