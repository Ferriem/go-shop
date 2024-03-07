package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/Ferriem/go-web/Product/common"
	"github.com/Ferriem/go-web/Product/datamodels"
	"github.com/Ferriem/go-web/Product/encrypt"
	"github.com/Ferriem/go-web/RabbitMQ"
)

var hostArray = []string{
	"127.0.0.1",
	"127.0.0.1",
}

var localHost = ""

var GetOneIp = "localhost"

var GetOnePort = "8084"

var port = "8083"

var hashConsistent *common.Consistent

var rabbitMQValidate *RabbitMQ.RabbitMQ

type AccessControl struct {
	sourcesArray map[int]interface{}
	sync.RWMutex
}

var accessControl = &AccessControl{sourcesArray: make(map[int]interface{})}

func (m *AccessControl) GetNewRecord(uid int) interface{} {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	data := m.sourcesArray[uid]
	return data
}

func (m *AccessControl) SetNewRecord(uid int) {
	m.RWMutex.Lock()
	m.RWMutex.Unlock()
	m.sourcesArray[uid] = "ferriem"
}

func (m *AccessControl) GetDistributedRight(req *http.Request) bool {
	uid, err := req.Cookie("uid")
	if err != nil {
		return false
	}

	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}

	if hostRequest == localHost {
		return m.GetDataFromMap(uid.Value)
	} else {
		return m.GetDataFromOtherMap(hostRequest, req)
	}
}

func (m *AccessControl) GetDataFromMap(uid string) bool {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}
	data := m.GetNewRecord(uidInt)

	if data != nil {
		return true
	}

	return false
}

func (m *AccessControl) GetDataFromOtherMap(host string, req *http.Request) bool {
	hostUrl := "http://" + host + ":" + port + "/checkRight"
	response, body, err := GetCurl(hostUrl, req)
	if err != nil {
		return false
	}

	if response.StatusCode == 200 {
		if string(body) == "true" {
			return true
		} else {
			return false
		}
	}
	return false
}

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

func GetCurl(hostUrl string, req *http.Request) (response *http.Response, body []byte, err error) {
	uidPre, err := req.Cookie("uid")
	if err != nil {
		return
	}
	sign, err := req.Cookie("sign")
	if err != nil {
		return
	}
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", "http://"+hostUrl+":"+port+"/access", nil)
	if err != nil {
		return
	}
	cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: sign.Value, Path: "/"}
	reqest.AddCookie(cookieUid)
	reqest.AddCookie(cookieSign)

	response, err = client.Do(reqest)
	defer response.Body.Close()
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(response.Body)
	return
}

func CheckRight(w http.ResponseWriter, r *http.Request) {
	right := accessControl.GetDistributedRight(r)
	if !right {
		w.Write([]byte("false"))
		return
	}
	w.Write([]byte("true"))
	return
}

func Check(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Exec check")
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil && len(queryForm["productID"]) <= 0 && len(queryForm["productID"][0]) <= 0 {
		w.Write([]byte("false"))
		return
	}
	productString := queryForm["productID"][0]
	fmt.Println(productString)
	userCookie, err := r.Cookie("uid")
	if err != nil {
		w.Write([]byte("false"))
	}

	right := accessControl.GetDistributedRight(r)
	if right == false {
		w.Write([]byte("false"))
		return
	}

	hostUrl := "http://" + GetOneIp + ":" + GetOnePort + "/getOne"
	response, body, err := GetCurl(hostUrl, r)

	if err != nil {
		w.Write([]byte("fail"))
		return
	}

	if response.StatusCode == 200 {
		if string(body) == "true" {
			productID, err := strconv.ParseInt(productString, 10, 64)
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			userID, err := strconv.ParseInt(userCookie.Value, 10, 64)
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			message := datamodels.NewMessage(userID, productID)
			byteMessage, err := json.Marshal(message)
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			err = rabbitMQValidate.PublishSimple(string(byteMessage))
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			w.Write([]byte("true"))
			return
		}
	}
	w.Write([]byte("false"))
	return
}

func main() {
	//consistent hash
	hashConsistent = common.NewConsistent()
	//add node
	for _, v := range hostArray {
		hashConsistent.Add(v)
	}

	localIp, err := common.GetIntranceIP()
	if err != nil {
		fmt.Println(err)
	}
	localHost = localIp
	fmt.Println(localHost)

	rabbitMQValidate = RabbitMQ.NewRabbitMQSimple("product")

	defer rabbitMQValidate.Destory()

	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("./fronted/web/htmlProductShow"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./fronted/web/assets"))))

	// filter
	filter := common.NewFilter()

	filter.RegisterFilterUri("/check", Auth)
	filter.RegisterFilterUri("/ckechRight", Auth)

	http.HandleFunc("/check", filter.Handle(Check))
	http.HandleFunc("/checkRight", filter.Handle(CheckRight))

	http.ListenAndServe(":8083", nil)
}
