package main

import(
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"bytes"
	"crypto/rand"
	"database/sql"
	"strings"
)

type LoginResponsStruct struct{
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data string `json:"data"`
}

func login(w http.ResponseWriter,r *http.Request){
	body,_:=ioutil.ReadAll(r.Body)
	bodyStr:=string(body)
	fmt.Println(bodyStr)

	var requestParams RegistRequestParams
	err:=json.Unmarshal(body,&requestParams)
	if err!=nil{
		fmt.Println(err)
		return
	}

	if !hasInitDB{
		initDB()
	}

	var count int64 = 0
	rows,err:=DB.Query("select * from user where nickname = ?",requestParams.Username)
	defer func(){
		if rows!=nil{
			rows.Close()
		}
	}()
	if err!=nil&&err!=sql.ErrNoRows{
		fmt.Println(err)
		return
	}
	user:=new(User)
	for rows.Next(){
		err=rows.Scan(&user.Id, &user.NickName, &user.Email, &user.Tel, &user.Sex, &user.AddrDes, &user.AddrLng, &user.AddrLat, &user.Pwd, &user.Token)
		count++
		if err!=nil{
			fmt.Println(err)
			return
		}
		break
	}
	loginResponsStruct:=new(LoginResponsStruct)
	if count>0&&user!=nil{
		//该用户存在，下一步认证密码是否正确
		if strings.EqualFold(requestParams.Pwd,user.Pwd){
			//登录成功
			loginResponsStruct.Code=0
			loginResponsStruct.Msg = "success"
			token :=createRandomString(20)
			rows,err:=DB.Query("update user set token=? where nickname=?",token,requestParams.Username)
			defer func(){
				if rows!=nil{
					rows.Close()
				}
			}()
			if err!=nil{
				fmt.Println(err)
				return
			}

			loginResponsStruct.Data = token
		}else{
			//密码输入有误
			loginResponsStruct.Code=1
			loginResponsStruct.Msg = "密码错误"
		}
	}else{
		//该用户不存在，尚未注册
		loginResponsStruct.Code = 1
		loginResponsStruct.Msg = "该用户尚未注册，请先注册"
	}

	loginResponsJson,err:=json.Marshal(loginResponsStruct)
	if err!=nil{
		fmt.Println(err)
		return
	}

	fmt.Fprintf(w,string(loginResponsJson))

}

//随机生成字符串
func createRandomString(len int)string{
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b:=bytes.NewBufferString(str)
	length:=b.Len()
	bigInt:=big.NewInt(int64(length))
	for i:=0;i<len;i++{
		randomInt,_:=rand.Int(rand.Reader,bigInt)
		container+=string(str[randomInt.Int64()])
	}
	return container
}