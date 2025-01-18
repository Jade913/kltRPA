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

func (o *OmoIntegrate) UpdateOmo(data []map[string]interface{}) ([]map[string]interface{}, error) {
	rpcUrl := fmt.Sprintf("https://%s/xmlrpc/2/object", o.server)
	log.Printf("开始更新 OMO 系统，RPC URL: %s", rpcUrl)

	client, err := xmlrpc.NewClient(rpcUrl, nil)
	if err != nil {
		log.Printf("创建 XML-RPC 客户端失败: %v", err)
		return nil, err
	}
	defer client.Close()

	results := make([]map[string]interface{}, len(data)) // 预分配切片大小

	for i, record := range data {
		var result interface{}
		args := []interface{}{
			o.db,
			o.uid,
			o.password,
			"kltcrm.import.customer",
			"create_customer",
			[]interface{}{i + 1, record},
		}

		err = client.Call("execute_kw", args, &result)
		if err != nil {
			log.Printf("调用 create_customer 失败: %v", err)
			return nil, err
		}

		// 将结果转换为map并存储
		if resultMap, ok := result.(map[string]interface{}); ok {
			results[i] = resultMap
			log.Printf("记录更新结果: %v", resultMap)
		}
	}

	log.Println("所有记录已成功更新到 OMO 系统")
	return results, nil
}
