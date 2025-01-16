package utils

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kolo/xmlrpc"
)

type OmoIntegrate struct {
	server   string
	db       string
	username string
	password string
	uid      int
}

func NewOmoIntegrate(server, db, username, password string) *OmoIntegrate {
	return &OmoIntegrate{
		server:   server,
		db:       db,
		username: username,
		password: password,
	}
}

func (o *OmoIntegrate) Login() (bool, error) {
	rpcUrl := fmt.Sprintf("https://%s/xmlrpc/2/common", o.server)
	log.Printf("尝试连接到服务器: %s", rpcUrl)

	client, err := xmlrpc.NewClient(rpcUrl, &http.Transport{})
	if err != nil {
		log.Printf("创建 XML-RPC 客户端失败: %v", err)
		return false, err
	}
	defer client.Close()

	var result interface{}
	args := []interface{}{o.db, o.username, o.password, map[string]interface{}{}}
	err = client.Call("authenticate", args, &result)

	if err != nil {
		log.Printf("调用 authenticate 失败: %v", err)
		return false, err
	}

	log.Printf("服务器返回结果类型: %T, 值: %v", result, result)

	switch v := result.(type) {
	case int:
		o.uid = v
	case int64:
		o.uid = int(v)
	case float64:
		o.uid = int(v)
	default:
		log.Printf("未知的返回类型: %T", v)
		return false, fmt.Errorf("unexpected return type: %T", v)
	}

	if o.uid != 0 {
		log.Printf("登录成功，UID: %d", o.uid)
		return true, nil
	}

	log.Println("登录失败，UID 为 0")
	return false, nil
}
