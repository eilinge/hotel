package routes

import (
	"fmt"
	"go-hotel/configs"
	"go-hotel/eths"
	"go-hotel/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// Guest ...
type Guest struct {
	IDCard string `json:"idcard"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Reword bool   `json:"reword"`
}

// PingHandler ...
func PingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

// UpLoad ...
func UpLoad(c echo.Context) error {
	//1. 响应数据结构初始化
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	//2. 解析数据
	form, _ := c.FormParams()
	name := form["name"][0]
	reword, _ := strconv.Atoi(form["reword"][0])
	idcard := form["idcard"][0]

	var status bool
	if reword == 1 {
		status = true
	} else {
		status = false
	}

	age, _ := strconv.Atoi(form["age"][0])
	fmt.Printf("form name: %s, reword: %v, idcard: %s, age: %d\n", name, reword, idcard, age)

	// Upload(from, pass, hash, name string, age int64, record bool)
	go eths.Upload(configs.Config.Eth.Fundation, configs.Config.Eth.FundationPWD, idcard, name, int64(age), status)

	respData := make(map[string]interface{})
	respData["guest"] = Guest{IDCard: idcard, Name: name, Age: age, Reword: status}

	resp.Data = respData

	return nil
}

// Check ...
func Check(c echo.Context) error {
	//1. 响应数据结构初始化
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	//2. 解析数据
	idcard := c.QueryParam("idcard")
	fmt.Println("idcard: ", idcard)
	// go func() {
	age, status, _ := eths.CheckReward(idcard)
	fmt.Println("status: ", status)

	respData := make(map[string]interface{})
	respData["CriminalRecord"] = status
	respData["Age"] = age

	respData["checkin"] = true
	if age < 18 || status == true {
		respData["checkin"] = false
	}
	resp.Data = respData
	// }()
	return nil
}
