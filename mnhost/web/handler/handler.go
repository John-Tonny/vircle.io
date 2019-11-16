package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"vircle.io/mnhost/model"

	vps "vircle.io/mnhost/web/proto/vps"

	"github.com/micro/go-micro/client"
)

func WebCall(w http.ResponseWriter, r *http.Request) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Println("start web")
	// call the backend service
	vpsClient := vps.NewVpsService("go.micro.srv.vps", client.DefaultClient)
	rsp, err := vpsClient.Call(context.TODO(), &vps.Request{
		Name: request["name"].(string),
	})
	if err != nil {
		fmt.Print("error 500")
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Println("start")
	//房屋切片信息
	node_list := []models.Node{}
	json.Unmarshal(rsp.Mix, &node_list)
	fmt.Println("1")
	//将房屋切片信息转换成map切片返回给前端
	var nodes1 []interface{}
	for _, nodeinfo := range node_list {
		fmt.Println(nodeinfo.Id)
		nodes1 = append(nodes1, nodeinfo.To_node_info())
	}
	data_map := make(map[string]interface{})
	data_map["nodes"] = nodes1

	fmt.Println(data_map["nodes"])

	// we want to augment the response
	response := map[string]interface{}{
		"msg":  rsp.Errno,
		"err":  rsp.Errmsg,
		"test": request["name"],
		"mmm":  data_map,
		"ref":  time.Now().UnixNano(),
	}

	w.Header().Set("Content-Type", "application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
