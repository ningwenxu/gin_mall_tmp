package service

import (
	"gin_mall_tmp/conf"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strconv"
)

func UploadAvatarToLocalStatic(file multipart.File, userId uint, userName string) (filePath string, err error) {
	bId := strconv.Itoa(int(userId)) //路径拼接
	basePath := "." + conf.AvatarPath + "user" + bId + "/"
	if !DirExistOrNot(basePath) {
		CreateDir(basePath)
	}
	avatarPath := basePath + userName + ".jpg" //todo:把file的后缀提取出来
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(avatarPath, content, 0666)
	if err != nil {
		return
	}
	return "user" + bId + "/" + userName + ".jpg", nil
}

func UploadProductToLocalStatic(file multipart.File, userId uint, productName string) (filePath string, err error) {
	bId := strconv.Itoa(int(userId)) //路径拼接
	basePath := "." + conf.ProductPath + "boss" + bId + "/"
	if !DirExistOrNot(basePath) {
		CreateDir(basePath)
	}
	productPath := basePath + productName + ".jpg" //todo:把file的后缀提取出来
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(productPath, content, 0666)
	if err != nil {
		return
	}
	return "boss" + bId + "/" + productName + ".jpg", nil
}

// 判断文件夹路径是否存在
func DirExistOrNot(fileAddr string) bool {
	s, err := os.Stat(fileAddr)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 创建文件夹
func CreateDir(dirName string) bool {
	err := os.MkdirAll(dirName, 755)
	if err != nil {
		return false
	}
	return true
}
