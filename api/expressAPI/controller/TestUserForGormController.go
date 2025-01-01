package controller

import (
	"apiProject/api/expressAPI/service"
	"apiProject/api/expressAPI/types/domain"
	"apiProject/api/response"
	"apiProject/api/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"net/http"
)

// TestUserForGormController 测试用户控制器
type TestUserForGormController struct {
	testUserGormService service.TestUserForGormService
}

// TestUserForGormControllerInit 用户控制器初始化
func TestUserForGormControllerInit(u service.TestUserForGormService) *TestUserForGormController {
	return &TestUserForGormController{testUserGormService: u}
}

// RegisterRoutes 注册测试用户控制器请求路由
func (u *TestUserForGormController) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/userGorm/{id}", u.handlerGetUser).Methods(utils.GET)
	r.HandleFunc("/userGorm/create", u.handlerCreateUser).Methods(utils.POST)
	r.HandleFunc("/userGorm/batchCreate", u.handlerBatchCreateUser).Methods(utils.POST)
	r.HandleFunc("/userGorm/login", u.handlerUserLogin).Methods(utils.POST)
	r.HandleFunc("/userGorm/update", u.handlerUpdateUser).Methods(utils.PUT)
	r.HandleFunc("/userGorm/batchUpdate", u.handlerBatchUpdateUser).Methods(utils.PUT)
	r.HandleFunc("/userGorm/delete/{id}", u.handlerDeleteUser).Methods(utils.DELETE)
	r.HandleFunc("/userGorm/batchDelete", u.handlerBatchDeleteUser).Methods(utils.POST)
}

func (u *TestUserForGormController) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var queryId = vars["id"]
	if queryId == "" || cast.ToInt64(queryId) == 0 {
		response.WriteJson(w, response.FailMessageResp("测试用户ID不能为空"))
		return
	}
	t, err := u.testUserGormService.GetUserById(utils.ConvertToInt64(queryId))
	if err != nil {
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}
	response.WriteJson(w, response.OkDataResp(t))
}

// handlerCreateUser 处理用户新增
func (u *TestUserForGormController) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var user *domain.TestUser
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		zap.L().Sugar().Errorf("测试用户新增(gorm)参数解析错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("测试用户新增参数解析失败"))
		return
	}

	defer utils.CloseBodyError("测试用户新增(gorm)请求", w, r)

	createdUser, err := u.testUserGormService.CreateUser(user)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	// 将json格式化输出
	marshal, _ := json.MarshalIndent(createdUser, "", "    ")
	zap.L().Sugar().Infof("测试用户新增===\n%s", marshal)

	response.WriteJson(w, response.OkDataResp(createdUser))
}

// handlerBatchCreateUser 处理测试用户批量新增
func (u *TestUserForGormController) handlerBatchCreateUser(w http.ResponseWriter, r *http.Request) {
	var list []*domain.TestUser
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		zap.L().Sugar().Errorf("测试用户批量新增(gorm)参数解析错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("测试用户批量新增参数解析失败"))
		return
	}

	defer utils.CloseBodyError("测试用户批量新增(gorm)请求", w, r)

	if len(list) == 0 {
		response.FailMessageResp("批量创建参数验证失败")
		return
	}

	count, err := u.testUserGormService.BatchCreateUser(list)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	if count <= 0 {
		response.FailMessageResp("批量创建失败")
		return
	}

	response.WriteJson(w, response.OkMessageResp("批量创建成功"))
}

// handlerUserLogin 用户登录
func (u *TestUserForGormController) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	var user *domain.TestUser

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		zap.L().Sugar().Errorf("解析测试用户登录参数错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("解析测试用户登录参数失败"))
		return
	}
	defer utils.CloseBodyError("测试用户登录请求", w, r)

	if user.Username == "" {
		response.WriteJson(w, response.FailMessageResp("用户名不能为空"))
		return
	}

	if user.Password == "" {
		response.WriteJson(w, response.FailMessageResp("密码不能为空"))
		return
	}

	loginUser, err := u.testUserGormService.UserLogin(user)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkDataResp(loginUser))
}

// handlerUpdateUser 处理用户修改
func (u *TestUserForGormController) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	var user *domain.TestUser
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		zap.L().Sugar().Errorf("测试用户修改(gorm)参数解析错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("测试用户修改(gorm)参数解析失败"))
		return
	}

	defer utils.CloseBodyError("测试用户修改(gorm)请求", w, r)

	count, err := u.testUserGormService.UpdateUser(user)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	if count <= 0 {
		response.FailMessageResp("更新失败")
		return
	}

	response.WriteJson(w, response.OkMessageResp("更新成功"))
}

// handlerBatchUpdateUser 处理批量更新
func (u *TestUserForGormController) handlerBatchUpdateUser(w http.ResponseWriter, r *http.Request) {
	var list []*domain.TestUser
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		zap.L().Sugar().Errorf("测试用户批量修改(gorm)参数解析错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("测试用户批量修改(gorm)参数解析失败"))
		return
	}

	defer utils.CloseBodyError("测试用户批量修改(gorm)请求", w, r)

	if len(list) == 0 {
		response.FailMessageResp("批量修改参数验证失败")
		return
	}

	// 使用map进行ID去重
	uniqueMap := make(map[int64]*domain.TestUser)

	for _, user := range list {
		// 如果ID已经存在，则跳过
		if _, exists := uniqueMap[user.Id]; !exists {
			uniqueMap[user.Id] = user
		}
	}

	// 将map中的值转换为切片
	var uniqueList []*domain.TestUser
	for _, user := range uniqueMap {
		uniqueList = append(uniqueList, user)
	}

	list = uniqueList

	count, err := u.testUserGormService.BatchUpdateUser(list)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	if count <= 0 {
		response.FailMessageResp("批量更新失败")
		return
	}

	response.WriteJson(w, response.OkMessageResp("批量更新成功"))
}

func (u *TestUserForGormController) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var queryId = vars["id"]
	if queryId == "" || cast.ToInt64(queryId) == 0 {
		response.WriteJson(w, response.FailMessageResp("测试用户ID不能为空"))
		return
	}
	count, err := u.testUserGormService.DeleteUser(utils.ConvertToInt64(queryId))
	if err != nil {
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	if count <= 0 {
		response.FailMessageResp("删除失败")
		return
	}
	response.WriteJson(w, response.OkMessageResp("删除成功"))
}

func (u *TestUserForGormController) handlerBatchDeleteUser(w http.ResponseWriter, r *http.Request) {
	var ids []any
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		zap.L().Sugar().Errorf("测试用户批量删除(gorm)参数解析错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("测试用户批量删除(gorm)参数解析失败"))
		return
	}

	defer utils.CloseBodyError("测试用户批量删除(gorm)请求", w, r)

	if len(ids) == 0 {
		response.FailMessageResp("批量删除参数错误")
		return
	}

	count, err := u.testUserGormService.BatchDeleteUser(ids)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	if count <= 0 {
		response.FailMessageResp("批量删除失败")
		return
	}
	response.WriteJson(w, response.OkMessageResp("批量删除成功"))
}
