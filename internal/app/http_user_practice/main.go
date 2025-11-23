package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User 定义用户结构体
type User struct {
	ID       int    `json:"id" gorm:"primaryKey`
	Username string `json:"username" gorm:"not null;size:100"`
	Email    string `json:"email" gorm:"unique;not null;size:100"`
	Age      int    `json:"age" gorm:"check:age >= 0 AND age <= 150"`
	Phone    int    `json:"phone" gorm:"type:int"`
}

// APIResponse 统一响应结构
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

var DB *gorm.DB

// 初始化数据库连接
func initDB() error {
	dsn := "root:123456@tcp(localhost:3306)/db2?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	// 自动迁移表结构（创建或更新表)
	//return DB.AutoMigrate(&User{})
	return nil
}

// 验证邮箱格式
func isValidEmail(email string) bool {
	return true
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`) //a.@b.c
	return regex.MatchString(email)
}

// 通过post查询参数添加用户的处理函数
func addUserWithGetHandler(w http.ResponseWriter, r *http.Request) { //w响应，r请求
	w.Header().Set("Content-Type", "application/json") //响应的头设置对应的k，v    告诉客户端响应的数据格式是json

	// 只允许post请求
	if r.Method != http.MethodPost { //如果请求的method方法不是post
		w.WriteHeader(http.StatusMethodNotAllowed) //返回错误
		json.NewEncoder(w).Encode(APIResponse{     //响应的新代码按照结构体形式
			Success: false,
			Error:   "只支持post方法",
		})
		return
	}

	// 从查询参数中获取用户信息
	username := r.PostFormValue("username") //获取username，存到变量username
	email := r.PostFormValue("email")
	ageStr := r.PostFormValue("age")
	phoneStr := r.PostFormValue("phone")

	// 数据验证
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{ //响应的json数据是结构体的格式
			Success: false,
			Error:   "用户名不能为空 (参数: username)",
		})
		return
	}

	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   "邮箱不能为空 (参数: email)",
		})
		return
	}

	if !isValidEmail(email) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   "邮箱格式无效",
		})
		return
	}

	phone, err := strconv.Atoi(phoneStr)
	if err != nil || phone < 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   "手机号必须是整数 (参数: phone)",
		})
		return
	}

	if phone == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   "电话号不能为空 (参数: phoneStr)",
		})
		return
	}

	// 解析年龄参数
	age, err := strconv.Atoi(ageStr)
	if err != nil || age < 0 || age > 150 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   "年龄必须是0-150之间整数 (参数: age)",
		})
		return
	}

	// 构造用户对象
	newUser := User{ //其中包含自动生成的id
		Username: username,
		Email:    email,
		Age:      age,
		Phone:    phone,
	}
	// 插入数据库
	result := DB.Create(&newUser)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   "插入数据库失败: " + result.Error.Error(),
		})
		return
	}
	// 返回成功响应
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Message: "用户添加成功",
		Data:    newUser,
	})
}

// 根据id查询用户信息
func getUserHandler(w http.ResponseWriter, r *http.Request) { //w响应，r请求
	w.Header().Set("Content-Type", "application/json") //响应的头设置对应的k，v    告诉客户端响应的数据格式是json

	// 从查询参数中获取用户id
	idStr := r.URL.Query().Get("id") //获取id，存到变量idStr
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   "缺少参数 id",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   "id 必须是数字",
		})
		return
	}

	var user User
	result := DB.First(&user, id) //根据id查询
	if result.Error == gorm.ErrRecordNotFound {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   "用户不存在",
		})
		return
	}
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   "数据库查询失败: " + result.Error.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    user,
	})
}

// 更新用户名
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Query().Get("id")
	username := r.URL.Query().Get("username")

	if idStr == "" || username == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   "缺少参数 id 或 username",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   "id 必须是数字",
		})
		return
	}

	// 检查用户是否存在
	var user User
	result := DB.Select("id").First(&user, id)
	if result.Error == gorm.ErrRecordNotFound {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   "用户不存在",
		})
		return
	}

	// 更新
	result = DB.Model(&User{}).Where("id = ?", id).Update("username", username)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   "更新失败: " + result.Error.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Message: "用户名更新成功",
		Data: map[string]interface{}{
			"id":       id,
			"username": username,
		},
	})
}

func main() {
	// 初始化数据库
	if err := initDB(); err != nil {
		log.Fatal("数据库连接失败: ", err)
	}

	// 注册路由
	http.HandleFunc("/api/user/add", addUserWithGetHandler) //绑定路径和函数，当客户端请求路径为""时使用这个函数处理请求
	http.HandleFunc("/api/user", getUserHandler)            // 查询
	http.HandleFunc("/api/user/update", updateUserHandler)  // 修改
	// 启动服务器
	log.Println("服务器启动在 :8080 端口")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
