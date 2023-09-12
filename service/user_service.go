package service

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gofly/dao"
	"gofly/global"
	"gofly/global/constants"
	"gofly/model"
	"gofly/service/dto"
	"gofly/utils"
	"strconv"
	"strings"
	"time"
)

var userService *UserService

type UserService struct {
	BaseService
	Dao *dao.UserDao
}

func NewUserService() *UserService {
	if userService == nil {
		userService = &UserService{
			Dao: dao.NewUserDao(),
		}
	}
	return userService
}
func (m *UserService) Login(iUserDTO dto.UserLoginDTO) (model.User, string, error) {
	var errResult error
	var token string
	iUser, err := m.Dao.GetUserByName(iUserDTO.Name)
	if err != nil || !utils.CompareHashAndPassword(iUser.Password, iUserDTO.Password) {
		errResult = errors.New("Invalid UserName Or Password")
	} else {
		//token, err = utils.GenerateToken(iUser.ID, iUser.Name)
		token, err = GenerateAndCacheLoginUserToken(iUser.ID, iUser.Name)
		if err != nil {
			errResult = errors.New(fmt.Sprintf("Generate Token Error: %s", err.Error()))
		}
	}
	return iUser, token, errResult
}
func (m *UserService) AddUser(iUserAddDTO *dto.UserAddDTO) error {
	if m.Dao.CheckUserNameExist(iUserAddDTO.Name) {
		return errors.New("user name exist")
	}
	return m.Dao.AddUser(iUserAddDTO)
}
func (m *UserService) GetUserById(ICommonIDDTO *dto.CommonIDDTO) (model.User, error) {
	return m.Dao.GetUserById(ICommonIDDTO.ID)
}
func (m *UserService) GetUserList(iUserListDTO *dto.UserListDTO) ([]model.User, int64, error) {
	return m.Dao.GetUserList(iUserListDTO)

}
func (m *UserService) UpdateUser(iUserUpdateDTO *dto.UserUpdateDTO) error {
	if iUserUpdateDTO.ID == 0 {
		return errors.New("Invalid User ID")
	}
	return m.Dao.UpdateUser(iUserUpdateDTO)
}
func (m *UserService) DeleteUserById(iCommonIDDTO *dto.CommonIDDTO) error {
	return m.Dao.DeleteUserById(iCommonIDDTO.ID)
}
func SetLoginUserTokenToRedis(uid uint, token string) error {
	return global.RedisClient.Set(strings.Replace(constants.LOGIN_USER_TOKEN_REDIS_KEY, "{ID}", strconv.Itoa(int(uid)), -1), token, viper.GetDuration("jwt.tokenExpire")*time.Minute)
}
func GenerateAndCacheLoginUserToken(nUserId uint, stUserName string) (string, error) {
	token, err := utils.GenerateToken(nUserId, stUserName)
	if err == nil {
		err = global.RedisClient.Set(strings.Replace(constants.LOGIN_USER_TOKEN_REDIS_KEY, "{ID}", strconv.Itoa(int(nUserId)), -1), token, viper.GetDuration("jwt.tokenExpire")*time.Minute)
	}
	return token, err
}
