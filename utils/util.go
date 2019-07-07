package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

const (
	RECODE_OK           = "0"
	RECODE_DBERR        = "4001"
	RECODE_NODATA       = "4002"
	RECODE_DATAEXIST    = "4003"
	RECODE_DATAERR      = "4004"
	RECODE_SESSIONERR   = "4101"
	RECODE_LOGINERR     = "4102"
	RECODE_PARAMERR     = "4103"
	RECODE_USERERR      = "4104"
	RECODE_HASHERR      = "4105"
	RECODE_PWDERR       = "4106"
	RECODE_EXISTSERR    = "4201"
	RECODE_IPCERR       = "4202"
	RECODE_THIRDERR     = "4301"
	RECODE_IOERR        = "4302"
	RECODE_SERVERERR    = "4500"
	RECODE_UNKNOWERR    = "4501"
	RECODE_MINTERR      = "4502"
	RECODE_ETHERERR     = "4503"
	RECODE_ERC20ERR     = "4504"
	RECODE_ERC20POORERR = "4505"
)

var recodeText = map[string]string{
	RECODE_OK:           "成功",
	RECODE_DBERR:        "数据库操作错误",
	RECODE_NODATA:       "无数据",
	RECODE_DATAEXIST:    "数据已存在",
	RECODE_DATAERR:      "数据错误",
	RECODE_SESSIONERR:   "用户未登录",
	RECODE_LOGINERR:     "用户登录失败",
	RECODE_PARAMERR:     "参数错误",
	RECODE_USERERR:      "用户不存在或密码错误",
	RECODE_HASHERR:      "HASH错误",
	RECODE_PWDERR:       "密码错误",
	RECODE_EXISTSERR:    "重复上传错误",
	RECODE_IPCERR:       "IPC错误",
	RECODE_THIRDERR:     "与以太坊交互失败",
	RECODE_IOERR:        "文件读写错误",
	RECODE_SERVERERR:    "内部错误",
	RECODE_UNKNOWERR:    "未知错误",
	RECODE_MINTERR:      "挖矿错误",
	RECODE_ETHERERR:     "以太操作错误",
	RECODE_ERC20ERR:     "ERC20转账错误",
	RECODE_ERC20POORERR: "ERC20余额不足",
}

// RecodeText ...
func RecodeText(code string) string {
	str, ok := recodeText[code]
	if ok {
		return str
	}
	return recodeText[RECODE_UNKNOWERR]
}

// Resp ...
type Resp struct {
	Errno  string      `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

// ResponseData resp数据响应
func ResponseData(c echo.Context, resp *Resp) {
	resp.ErrMsg = RecodeText(resp.Errno)
	c.JSON(http.StatusOK, resp) // 响应给终端
}

// GetFileName 读取dir目录下文件名带address的文件
/*
address: d559d8a8990d0ca570737e3d32c89409c125f0c0
dirname: ../../../data/keystore
fileName: UTC--2019-06-24T13-16-07.992079300Z--d559d8a8990d0ca570737e3d32c89409c125f0c0
*/
func GetFileName(address, dirname string) (string, error) {

	data, err := ioutil.ReadDir(dirname)
	if err != nil {
		fmt.Println("read dir err", err)
		return "", err
	}
	for _, v := range data {
		if strings.Index(v.Name(), address) > 0 {
			//代表找到文件
			return v.Name(), nil
		}
	}
	return "", nil
}
