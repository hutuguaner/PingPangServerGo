package main

import(
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"
	
)

type User struct {
	Id       int64  `db:"id" json:"id"`
	NickName string `db:"nickname" json:"nickname"`
	Email    string `db:"email" json:"email"`
	Tel      string `db:"tel" json:"tel"`
	Sex      string `db:"sex" json:"sex"`
	AddrDes  string `db:"addr_des" json:"addr_des"`
	AddrLng  string `db:"addr_lng" json:"addr_lng"`
	AddrLat  string `db:"addr_lat" json:"addr_lat"`
	Pwd      string `db:"pwd" json:"pwd"`
	Token    string `db:"token" json:"token"`
}

type Ret struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []User `json:"data"`
}

func getAllUser(w http.ResponseWriter, r *http.Request) {

	body,_:=ioutil.ReadAll(r.Body)
	bodyStr:=string(body)
	fmt.Println(bodyStr)

	if !hasInitDB{
		initDB()
	}

	user := new(User)
	rows, err := DB.Query("select * from user")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		fmt.Printf("query failed , err:%v", err)
		return
	}
	ret := new(Ret)
	ret.Code = 0
	ret.Msg = "success"

	for rows.Next() {
		err = rows.Scan(&user.Id, &user.NickName, &user.Email, &user.Tel, &user.Sex, &user.AddrDes, &user.AddrLng, &user.AddrLat, &user.Pwd, &user.Token)
		if err != nil {
			fmt.Printf("scan failed,err:%v", err)
			return
		}
		ret.Data = append(ret.Data, *user)
		//fmt.Print(*user)
	}
	jsonData, err := json.Marshal(ret)
	if err != nil {
		fmt.Printf("json marshal err:%v", err)
		return
	}
	fmt.Fprintf(w,string(jsonData))

	//fmt.Fprintf(w,"sfdfdsfsfsdfdfs")
}
