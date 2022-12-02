package service

import (
	"context"
	"gin_mall_tmp/conf"
	"gin_mall_tmp/dao"
	"gin_mall_tmp/model"
	"gin_mall_tmp/pkg/e"
	"gin_mall_tmp/pkg/util"
	"gin_mall_tmp/serializer"
	"gopkg.in/mail.v2"
	"mime/multipart"
	"strings"
	"time"
)

type UserService struct {
	NickName string `json:"nick_name" form:"nick_name"`
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
	Key      string `json:"key" form:"key"` // 密钥  前端验证
}
type SendEmailService struct {
	Email         string `json:"email" form:"email"`
	Password      string `json:"password" form:"password"`
	OperationType uint   `json:"operation_type" form:"operation_type"`
	//1.绑定邮箱 2.解绑邮箱  3 改密码

}

type ValidEmailService struct {
}
type ShowMoneyService struct {
	Key string `json:"key" form:"key"`
}

func (service *UserService) Register(ctx context.Context) serializer.Response {
	var user model.User
	code := e.Success
	if service.Key == "" || len(service.Key) != 16 {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  "密钥长度不足",
		}
	}
	// 10000 ---->  密文存储  对称加密
	util.Encrypt.SetKey(service.Key)

	userDao := dao.NewUserDao(ctx) //数据库连接
	_, exist, err := userDao.ExistOrNotByUserName(service.UserName)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	if exist {
		code = e.ErrorExistUser
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}

	}
	user = model.User{
		UserName: service.UserName,
		Email:    "",
		NickName: service.NickName,
		Status:   model.Active,
		Avatar:   "avatar.jpg",
		Money:    util.Encrypt.AesEncoding("10000"), //初始金额的加密
	}
	//密码的加密
	if err = user.SetPassword(service.Password); err != nil {
		code = e.ErrorFailEncryption
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}

	}
	//创建用户
	err = userDao.CreateUser(&user)
	if err != nil {
		code = e.Error

	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

func (service *UserService) Login(ctx context.Context) serializer.Response {
	var user *model.User
	code := e.Success
	userDao := dao.NewUserDao(ctx)
	//判断用户是否存在
	user, exist, err := userDao.ExistOrNotByUserName(service.UserName)
	if !exist || err != nil {
		code = e.ErrorExistUserNotFound
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "用户不存在，请先注册",
		}
	}
	//校验密码
	if user.CheckPassword(service.Password) == false {
		code = e.ErrorNotCompare
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "密码错误，请重新登录",
		}
	}
	//http 无状态 （认证 token）  token签发
	token, err := util.GenerateToken(user.ID, service.UserName, 0)
	if err != nil {
		code = e.ErrorAuthToken
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "token校验错误",
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.TokenData{User: serializer.BuildUser(user), Token: token},
		Msg:    e.GetMsg(code),
	}
}

// Update 用户修改信息
func (service UserService) Update(ctx context.Context, uId uint) serializer.Response {
	var user *model.User
	var err error
	code := e.Success
	//找到这个用户
	userDao := dao.NewUserDao(ctx)
	user, err = userDao.GetUserById(uId)
	// 修改昵称  nickname
	if service.NickName != "" {
		user.NickName = service.NickName

	}
	err = userDao.UpdateUserById(uId, user)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Error:  err.Error(),
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildUser(user),
	}
}

// Post 头像更新
func (service *UserService) Post(ctx context.Context, uId uint, file multipart.File, fileSize int64) serializer.Response {
	code := e.Success
	var user *model.User
	var err error
	userDao := dao.NewUserDao(ctx)
	user, err = userDao.GetUserById(uId)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Error:  err.Error(),
			Msg:    e.GetMsg(code),
		}
	}
	path, err := UploadAvatarToLocalStatic(file, uId, user.UserName)
	if err != nil {
		code = e.ErrorUploadFail
		return serializer.Response{
			Status: code,
			Error:  err.Error(),
			Msg:    e.GetMsg(code),
		}
	}
	user.Avatar = path
	err = userDao.UpdateUserById(uId, user)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Error:  err.Error(),
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildUser(user),
		Msg:    e.GetMsg(code),
	}
}

// 发送邮箱
func (service *SendEmailService) Send(ctx context.Context, uId uint) serializer.Response {
	code := e.Success
	var address string
	var notice *model.Notice //绑定邮箱  修改密码 模板通知
	token, err := util.GenerateEmailToken(uId, service.OperationType, service.Email, service.Password)
	if err != nil {
		code = e.ErrorAuthToken
		return serializer.Response{
			Status: code,
			Error:  err.Error(),
			Msg:    e.GetMsg(code),
		}
	}
	noticeDao := dao.NewNoticeDao(ctx)
	notice, err = noticeDao.GetNoticeById(service.OperationType)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Error:  err.Error(),
			Msg:    e.GetMsg(code),
		}
	}
	address = conf.ValidEmail + token //发送方
	print(address)
	mailStr := notice.Text
	mailText := strings.Replace(mailStr, "Email", address, -1)
	m := mail.NewMessage()
	m.SetHeader("From", conf.SmtpEmail)
	m.SetHeader("To", service.Email)
	m.SetHeader("Subject", "FanOne")
	m.SetBody("text/html", mailText)
	d := mail.NewDialer(conf.SmtpHost, 465, conf.SmtpEmail, conf.SmtpPass)
	d.StartTLSPolicy = mail.MandatoryStartTLS
	if err = d.DialAndSend(m); err != nil {
		code = e.ErrorSendEmail
		return serializer.Response{
			Status: code,

			Msg: e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,

		Msg: e.GetMsg(code),
	}
}

// 验证邮箱
func (service *ValidEmailService) Valid(ctx context.Context, token string) serializer.Response {
	var userId uint
	var email string
	var password string
	var operationType uint
	code := e.Success
	//验证token
	if token == "" {
		code = e.InvalidParams
	} else {
		claims, err := util.ParseEmailToken(token)
		if err != nil {
			code = e.ErrorAuthToken
		} else if time.Now().Unix() > claims.ExpiresAt {
			code = e.ErrorAuthCheckTokenTimeout
		} else {
			userId = claims.UserID
			email = claims.Email
			password = claims.Password
			operationType = claims.OperationType
		}
	}
	if code != e.Success {
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	//获取用户
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserById(userId)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	if operationType == 1 {
		//绑定邮箱
		user.Email = email
	} else if operationType == 2 {
		//解绑邮箱
		user.Email = ""
	} else if operationType == 3 {
		err = user.SetPassword(password)
		if err != nil {
			code = e.Error
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			}
		}
	}
	err = userDao.UpdateUserById(userId, user)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildUser(user),
	}
}

// 展示用户金额
func (service *ShowMoneyService) Show(ctx context.Context, uId uint) serializer.Response {
	code := e.Success
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserById(uId)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildMoney(user, service.Key),
		Msg:    e.GetMsg(code),
	}
}
