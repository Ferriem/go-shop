package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Ferriem/go-web/Product/common"
	"github.com/Ferriem/go-web/Product/datamodels"
	"github.com/Ferriem/go-web/Product/encrypt"
)

func Auth(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("Exec Auth")
	//cookie
	err := CheckUserInfo(r)
	if err != nil {
		return err
	}
	return nil
}

func CheckUserInfo(r *http.Request) error {
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		return errors.New("fail to get uid cookie")
	}
	signCookie, err := r.Cookie("sign")
	if err != nil {
		return errors.New("fail to get sign cookie")
	}
	signByte, err := encrypt.DePwdCode(signCookie.Value)
	if err != nil {
		return errors.New("sign is invalid")
	}
	userInfo := datamodels.UserInfo{}
	err = json.Unmarshal(signByte, &userInfo)
	if err != nil {
		return err
	}
	userID, err := strconv.Atoi(uidCookie.Value)
	if userInfo.ID != int64(userID) {
		return errors.New("Dismatch")
	}
	return nil
}

func Check(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Exec check")
}

func main() {
	// filter
	filter := common.NewFilter()

	filter.RegisterFilterUri("/check", Auth)

	http.HandleFunc("/check", filter.Handle(Check))

	http.ListenAndServe(":8083", nil)
}
