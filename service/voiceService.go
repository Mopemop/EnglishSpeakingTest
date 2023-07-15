package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"time"
)

type signatureOrigin struct {
	Host        string
	Date        string
	RequestLine string
}

type auth struct {
	ApiKey    string
	Algorithm string
	Headers   string
	Signature string
}

// 需要发送的信息
type paraMessage struct {
	Sub      string `json:"sub"`
	Ent      string `json:"ent"`
	Category string `json:"category"`
	Cmd      string `json:"cmd"`
	Text     string `json:"text"`
	Tte      string `json:"tte"`
	Ttp_skip bool   `json:"ttp_skip"`
	Aue      string `json:"aue"`
	Rstcd    string `json:"rstcd"`
}
type paraRequest struct {
	Common   CommonData  `json:"common"`
	Business paraMessage `json:"business"`
	Data     paraData    `json:"data"`
}

type CommonData struct {
	AppID string `json:"app_id"`
}

type paraData struct {
	Status int `json:"status"`
}

type contentData struct {
	Status int    `json:"status"`
	Data   string `json:"data"`
}

type contentBusiness struct {
	Cmd string `json:"cmd"`
	Aus int    `json:"aus"`
}

type contentRequest struct {
	Business contentBusiness `json:"business"`
	Data     contentData     `json:"data"`
}

type response struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Sid     string       `json:"sid"`
	Data    responseData `json:"data"`
}

type responseData struct {
	Status int    `json:"status"`
	Data   string `json:"data"`
}

func connectAPI() string {
	//创建连接字符
	signatureOrigin := signatureOrigin{
		Host:        "ise-api.xfyun.cn",
		Date:        time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123),
		RequestLine: "GET /v2/open-ise HTTP/1.1",
	}
	h := hmac.New(sha256.New, []byte("科大讯飞控制台的APISecret"))
	h.Write([]byte("host: " + signatureOrigin.Host + "\n" + "date: " + signatureOrigin.Date + "\n" + signatureOrigin.RequestLine))
	auth := auth{
		ApiKey:    "科大讯飞控制台的APIkey",
		Algorithm: "hmac-sha256",
		Signature: base64.StdEncoding.EncodeToString(h.Sum(nil)),
		Headers:   "host date request-line",
	}
	authorization := "api_key=" + "\"" + auth.ApiKey + "\"" + ", algorithm=\"hmac-sha256\", headers=\"host date request-line\", signature=\"" + auth.Signature + "\""
	authorization = base64.StdEncoding.EncodeToString([]byte(authorization))
	URL := "wss://ise-api.xfyun.cn/v2/open-ise?authorization=" + authorization + "&host=ise-api.xfyun.cn&date=" + url.QueryEscape(signatureOrigin.Date)
	// 连接WebSocket服务器
	return URL
}

// 设置第一次发送的信息
func setSendMassage() paraMessage {
	message := paraMessage{}
	message.Sub = "ise"
	message.Ent = "en_vip"
	message.Category = "read_chapter"
	message.Tte = "utf-8"
	//message.Aue = "lame"
	message.Ttp_skip = true
	message.Cmd = "ssb"
	message.Text = "\uFEFF[content]" + "\n" + "文本内容"
	message.Rstcd = "utf8"
	return message
}

func Communication(voice []byte) (bool, ReadChapter) {
	// 连接WebSocket服务器
	conn, http, err := websocket.DefaultDialer.Dial(connectAPI(), nil)
	fmt.Println(http.Status)
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	frameSize := 1280
	status := 0
	var message paraMessage
	message = setSendMassage()

	// 发送消息
	// 发送第一次消息，确定交互信息
	common := CommonData{"科大讯飞控制台的APPID"}
	data := paraData{status}
	req := paraRequest{common, message, data}
	jsonBytes, _ := json.Marshal(req)
	err = conn.WriteMessage(websocket.TextMessage, jsonBytes)
	if err != nil {
		fmt.Println("err1")
		log.Fatal(err)
		return false, ReadChapter{}
	}
	//后续发送音频数据
	status = 1
	//第一帧
	voiceData := base64.StdEncoding.EncodeToString(voice[0:frameSize])
	business := contentBusiness{"auw", 1}
	cData := contentData{status, voiceData}
	req1 := contentRequest{business, cData}
	jsonBytes, _ = json.Marshal(req1)
	err = conn.WriteMessage(websocket.TextMessage, jsonBytes)
	//中间帧
	index := frameSize
	for ; index+frameSize < len(voice); index = index + frameSize {
		voiceData = base64.StdEncoding.EncodeToString(voice[index : index+frameSize])
		business = contentBusiness{"auw", 2}
		cData = contentData{status, voiceData}
		req1 = contentRequest{business, cData}
		jsonBytes, _ = json.Marshal(req1)
		err = conn.WriteMessage(websocket.TextMessage, jsonBytes)
	}
	//最后帧
	status = 2
	voiceData = base64.StdEncoding.EncodeToString(voice[index:])
	business = contentBusiness{"auw", 4}
	cData = contentData{status, voiceData}
	req1 = contentRequest{business, cData}
	jsonBytes, _ = json.Marshal(req1)
	err = conn.WriteMessage(websocket.TextMessage, jsonBytes)
	response := response{}
	var success bool
	for true {
		success, response = readReturn(conn)
		if response.Data.Status == 2 {
			break
		} else if success == false {
			return success, ReadChapter{}
		}
	}
	byteResult, err := base64.StdEncoding.DecodeString(response.Data.Data)
	if err != nil {
		fmt.Println("err3")
		log.Fatal(err)
	}
	result := HandleVoiceXML(byteResult)
	return true, result
}

// 读取返回数据
func readReturn(conn *websocket.Conn) (bool, response) {
	response := response{}
	_, p, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("err1")
		return false, response
	}
	err = json.Unmarshal(p, &response)
	if err != nil {
		fmt.Println("err2")
		return false, response
	}
	return true, response
}
