package routes

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"go-echo/configs"
	"go-echo/dbs"
	"go-echo/eths"
	"go-echo/utils"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

// PageMaxPic ...
const (
	PageMaxPic    = 3
	defaultFormat = "2006-01-02 15:04:05 PM"
)

// Price ...
var (
	Price  int64
	mutex  sync.Mutex
	endBid bool
)

// PingHandler ...
func PingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

/*
17601329166@163.com
eilinge
eilinge
*/

var passwd string

// Register ...
func Register(c echo.Context) error {
	//1. 响应数据结构初始化
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	//2. 解析数据
	account := &dbs.Accounts{}

	/*
		将前端传过来的数据, 与dbs.Accounts进行数据绑定
		&dbs.Account{
			Email       `json:"email"`			name="email"
			IdentitiyID `json:"identity_id"`	name="identity_id"
			UserName 	`json:"username"`		name="username"
		}
	*/
	if err := c.Bind(account); err != nil { // 解析form表单
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}
	//3. 操作geth创建账户(account.IdentitiyId->pass)
	passwd = account.IdentitiyID
	address, err := eths.NewAcc(account.IdentitiyID, configs.Config.Eth.Connstr)
	if err != nil {
		// fmt.Println("failed to NewAcc: ", err)
		resp.Errno = utils.RECODE_IPCERR
		return err
	}
	//4. 操作Mysql插入数据
	pwd := fmt.Sprintf("%x", sha256.Sum256([]byte(account.IdentitiyID)))
	sql := fmt.Sprintf("insert into account(email, username, identity_id, address) values('%s', '%s', '%s', '%s')",
		account.Email, account.UserName, pwd, address)
	fmt.Println(sql)
	_, err = dbs.Create(sql)
	if err != nil {
		resp.Errno = utils.RECODE_DBERR
		return err
	}
	//5. session处理
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
	}
	sess.Values["address"] = address
	sess.Values["username"] = account.UserName
	sess.Save(c.Request(), c.Response())
	return nil
}

// Login ...
func Login(c echo.Context) error {
	//1. 响应数据结构初始化
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	//2. 解析数据
	account := &dbs.Accounts{}

	if err := c.Bind(account); err != nil { // 解析form表单
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}
	// fmt.Println("account.IdentitiyID: ", account.UserName, account.IdentitiyID)
	passwd = account.IdentitiyID

	//3. 操作Mysql查询数据
	pwd := fmt.Sprintf("%x", sha256.Sum256([]byte(account.IdentitiyID)))
	sql := fmt.Sprintf("select * from account where username='%s' and identity_id='%s';",
		account.UserName, pwd)
	// fmt.Println(sql)
	values, num, err := dbs.DBQuery(sql)
	if err != nil || num <= 0 {
		resp.Errno = utils.RECODE_DATAERR
		return err
	}
	row1 := values[0]
	//5. session处理
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
	}
	sess.Values["address"] = row1["address"]
	sess.Values["username"] = row1["username"]
	sess.Save(c.Request(), c.Response())
	return nil
}

// GetSession ....
func GetSession(c echo.Context) error {
	//1. 响应数据结构初始化
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)

	sess, err := session.Get("session", c)
	if err != nil {
		fmt.Println("failed to get session")
		resp.Errno = utils.RECODE_SESSIONERR
		return err
	}
	address := sess.Values["address"]
	// username := sess.Values["username"]
	if address == nil {
		fmt.Println("failed to get address")
		resp.Errno = utils.RECODE_SESSIONERR
		return err
	}
	return nil
}

// UpLoad ...
func UpLoad(c echo.Context) error {
	//1. 响应数据结构初始化
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	//2. 解析数据
	content := &dbs.Content{}

	price, _ := strconv.ParseInt(c.FormValue("price"), 10, 32)
	weight, _ := strconv.ParseInt(c.FormValue("weight"), 10, 32)

	content.Price = price
	content.Weight = weight

	h, err := c.FormFile("fileName") // 解析文件名
	if err != nil {
		fmt.Println("failed to FormFile: ", err)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}
	src, err := h.Open()
	defer src.Close()
	// 3. 打开一个目标文件用于保存
	content.Content = "static/photo/" + h.Filename
	dst, err := os.Create(content.Content)
	if err != nil {
		fmt.Println("failed to create file: ", err)
		resp.Errno = utils.RECODE_IOERR
		return err
	}
	defer dst.Close()

	// 4. get hash
	cData := make([]byte, h.Size)
	n, err := src.Read(cData)
	if err != nil || h.Size != int64(n) {
		resp.Errno = utils.RECODE_IOERR
		return err
	}
	content.ContentHash = fmt.Sprintf("%x", sha256.Sum256(cData))
	dst.Write(cData) // 图片存储

	content.Title = h.Filename
	// 5. write to dbs / 给上传图片页面, 添加weight, price, 并一起传入
	// content.AddContent()

	// 6. 操作以太坊
	sess, _ := session.Get("session", c)
	fromAddr, ok := sess.Values["address"].(string)
	if fromAddr == "" || !ok {
		resp.Errno = utils.RECODE_SESSIONERR
		return errors.New("no session")
	}
	// from, pass, hash, data string
	// fmt.Printf("price: %d, weight: %d\n", price, weight)
	// go func() {
	// 使用go func开启协程, 则当挖矿失败, 无法返回resp.Errno
	err = eths.Upload(fromAddr, passwd, content.ContentHash, name, age, record)
	if err != nil {
		resp.Errno = utils.RECODE_MINTERR
		return err
	}
	content.AddContent()
	// }()
	return nil
}
