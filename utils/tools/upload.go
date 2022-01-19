package tools

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"go-bot/setting"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

func GetImageFullUrl(name string) string {
	return setting.RuntimeSetting.ImagePrefixUrl + "/" + GetImagePath() + name
}

func GetImageName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = EncodeMD5(fileName)

	return fileName + ext
}

func GetImagePath() string {
	return setting.RuntimeSetting.ImageSavePath
}

func GetImageFullPath() string {
	return setting.RuntimeSetting.RuntimeRootPath + GetImagePath()
}

func CheckImageExt(filename string) bool {
	ext := GetExt(filename)
	for _, allowExt := range setting.RuntimeSetting.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}

	return false
}

func CheckImageSize(f multipart.File) bool {
	size, err := GetSize(f)
	if err != nil {
		return false
	}

	return size <= setting.RuntimeSetting.ImageMaxSize*1024*1024
}

func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}
	err = IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	perm := CheckPermission(src)
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	return nil
}

func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}
