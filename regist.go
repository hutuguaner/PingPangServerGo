package main

import(
	"database/sql"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type RegistRequestParams struct{
	Username string `json:"username"`
	Pwd string `json:"pwd"`
}

func regist(w http.ResponseWriter,r *http.Request){
	body,_:=ioutil.ReadAll(r.Body)
	bodyStr:=string(body)
	fmt.Println(bodyStr)
	var requestParams RegistRequestParams
	err:=json.Unmarshal(body,&requestParams)
	if err != nil{
		fmt.Println(err)
		return
	}
	
	if !hasInitDB{
		initDB()
	}

	var count int64 = 0
	rows,err:=DB.Query("select * from user where nickname=?",requestParams.Username)
	defer func(){
		if rows!=nil{
			rows.Close()
		}
	}()
	if err!=nil&&err!=sql.ErrNoRows{
		fmt.Println(err)
		return
	}
	
	for rows.Next(){
		count++
	}

	ret:=new(Ret)
	if count>0{
		//説明該賬號已经存在 ，不能再注册了
		ret.Code = 1
		ret.Msg = "当前账号已被使用"
	}else {
		//当前注册账号不存在，可以注册
		rows,err:=DB.Query("insert into user(nickname,pwd)values(?,?)",requestParams.Username,requestParams.Pwd)
		if err!=nil{
			fmt.Println(err)
			return
		}
		defer func(){
			if rows!=nil{
				rows.Close()
			}
		}()

		ret.Code = 0
		ret.Msg = "success"
	}
	responseJson,err:=json.Marshal(ret)
	if err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w,string(responseJson))

}