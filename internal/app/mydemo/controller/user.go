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
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "用户名不能为空 (参数: username)",
		})
		return
	}

	if email == "" {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "邮箱不能为空 (参数: email)",
		})
		return
	}

	if !utils.IsValidEmail(email) {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "邮箱格式无效",
		})
		return
	}

	phone, err := strconv.Atoi(phoneStr)
	if err != nil || phone < 0 {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "手机号必须是整数 (参数: phone)",
		})
		return
	}

	if phone == 0 {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "电话号不能为空 (参数: phoneStr)",
		})
		return
	}

	// 解析年龄参数
	age, err := strconv.Atoi(ageStr)
	if err != nil || age < 0 || age > 150 {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "年龄必须是0-150之间整数 (参数: age)",
		})
		return
	}

	resp, err := http.Get("http://v.juhe.cn/historyWeather/weather?weather=sunny")
	if err != nil {
		return
	}

	//在函数结束前关闭resp结构体的Body属性
	defer resp.Body.Close()
	//读resp的体Body，返回[]byte和错误
	body, err := io.ReadAll(resp.Body)
	if err != nil {
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
		panic(err)
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
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   "插入数据库失败: " + result.Error.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "用户添加成功",
		Data:    newUser,
	})
}

// 根据id查询用户信息
func GetUserHandler(c *gin.Context) { //
	// 从查询参数中获取用户id
	idStr := c.Query("id") //获取id，存到变量idStr
	if idStr == "" {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "缺少参数 id",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "id 必须是数字",
		})
		return
	}

	//查 Redis
	user, cacheErr := service.GetUserFromCache(id)
	if cacheErr != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Message: "从缓存获取失败",
		})
		return
	} else {

	}
	if user != nil {
		c.JSON(http.StatusOK, model.APIResponse{
			Success: true,
			Message: "从缓存获取",
			Data:    user,
		})

	} else {
		//redis没有，查数据库
		var result *gorm.DB
		result, user = service.GerUserById(id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, model.APIResponse{
					Success: false,
					Error:   "用户不存在",
				})
				return
			} else {
				c.JSON(http.StatusInternalServerError, model.APIResponse{
					Success: false,
					Error:   "数据库查询失败: " + result.Error.Error(),
				})
				return
			}

		}

		//写入redis
		err = service.SetUserToCache(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.APIResponse{
				Success: false,
				Error:   "写入redis失败 " + err.Error(),
			})
			return
		}

		//响应成功
		c.JSON(http.StatusOK, model.APIResponse{
			Success: true,
			Message: "新加入缓存",
			Data:    user,
		})
	}
}

// 更新用户名
func UpdateUserHandler(c *gin.Context) {

	idStr := c.PostForm("id")
	username := c.PostForm("username")

	if idStr == "" || username == "" {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "缺少参数 id 或 username",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		service.Logger.Error("格式错误", zap.String("idstr", idStr), zap.Error(err))
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "id 必须是数字",
		})
		return
	}

	// 检查用户是否存在
	result, _ := service.GerUserById(id)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, model.APIResponse{
			Success: false,
			Error:   "用户不存在",
		})
		return
	}

	// 更新
	result = service.UpdateUserById(id, username)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   "更新失败: " + result.Error.Error(),
		})
		return
	}

	//更新后删除 Redis 缓存，下次查询会重建
	err = service.DelRedisUser(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   "删除失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "用户名更新成功",
		Data: map[string]interface{}{
			"id":       id,
			"username": username,
		},
	})
}

// 获取所有用户并按年龄排序
func GetUsersHandlerAll(c *gin.Context) {
	order := c.Query("order")
	if order != "desc" && order != "asc" {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "order只能是desc或asc",
		})
		return
	}

	pagestr := c.Query("page")
	// if pagestr == "" {
	// 	c.JSON(http.StatusBadRequest, model.APIResponse{
	// 		Success: false,
	// 		Error:   "page不能为空",
	// 	})
	// 	return
	// }

	page, err := strconv.Atoi(pagestr)
	if err != nil {
		service.Logger.Error("格式错误", zap.String("pagestr", pagestr), zap.Error(err))
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "page 必须是数字",
		})
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
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   "缓存查询用户数量错误" + err.Error(),
		})
		return
	}

	if total != 0 {
		service.Logger.Debug("缓存的用户数量为", zap.String("total", strconv.Itoa(int(total))))
	}

	//先获取个数
	total, err = service.GetUserCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   "获取用户数量错误" + err.Error(),
		})
		return
	}

	//在缓存存count
	err = service.SetRedisUserCount(total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   "存缓存错误" + err.Error(),
		})
		return
	} else {
		service.Logger.Debug("缓存的用户数量为", zap.String("total", strconv.Itoa(int(total))))
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
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   "查询redis失败: " + err.Error(),
		})
		return
	} else {

	}

	if users != nil {
		c.JSON(http.StatusOK, model.APIResponse{
			Success: true,
			Message: "查询redis成功",
			Data: map[string]interface{}{
				"strSlice": users,
				"hasNext":  hasNext,
			},
		})
		return
	} else {

		users, err = service.GetUserByPage(page, pagesize, isasc)
		// page=100时，users，err怎么样，user cap 变为20，值[]
		if err != nil {
			// 记录错误日志
			service.Logger.Error("数据库查询失败", zap.Error(err))
			c.JSON(http.StatusInternalServerError, model.APIResponse{
				Success: false,
				Error:   "查询失败: " + err.Error(),
			})
			return
		}

		//添加缓存
		//在没有数据也要存缓存，page=100页，存的是[]
		err = service.SetRedisUserSlice(users, order, page, pagesize)
		if err != nil {
			service.Logger.Error("添加缓存失败", zap.Error(err))
			c.JSON(http.StatusInternalServerError, model.APIResponse{
				Success: false,
				Error:   "添加缓存失败: " + err.Error(),
			})
			return
		}

		service.Logger.Debug("users值", zap.Any("users", users))
		//service.Logger.Debug("user值", zap.Any("user", users[0]))
		strusers := fmt.Sprintf("users值%v", users)
		//struser := fmt.Sprintf("user值%+v", users[0])
		service.Logger.Debug("users值", zap.String("users", strusers))
		//service.Logger.Debug("user值", zap.String("user", struser))
		// 记录查到的用户数量
		service.Logger.Info("数据库查询成功", zap.Int("用户数量", len(users)))

		c.JSON(http.StatusOK, model.APIResponse{
			Success: true,
			Message: "新添加redis成功",
			Data: map[string]interface{}{
				"strSlice": users,   //缓存
				"hasNext":  hasNext, //判断是否有下一页
			},
		})
	}
}
