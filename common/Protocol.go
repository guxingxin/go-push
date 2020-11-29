package common

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

// 推送类型
const (
	PUSH_TYPE_ROOM = 1 // 推送房间
	PUSH_TYPE_ALL  = 2 // 推送在线
)

// websocket的Message对象
type WSMessage struct {
	MsgType int
	MsgData []byte
}

// 业务消息的固定格式(type+data)
type BizMessage struct {
	Type string          `json:"type"` // type消息类型: PING, PONG, JOIN, LEAVE, PUSH
	Data json.RawMessage `json:"data"` // data数据字段
}

// Data数据类型

// PUSH
type BizPushData struct {
	Items []*json.RawMessage `json:"items"`
}

// PING
type BizPingData struct{}

// PONG
type BizPongData struct{}

// JOIN
type BizJoinData struct {
	Token string `json:"token"`
	Room  string `json:"room"`
}

type AuthCheckArgs struct {
	//token
	AccessToken string `json:"access_token" form:"access_token"`
	//room id
	RoomId string `json:"room_id" form:"room_id"`
}

//接口业务响应码和说明信息
type Status struct {
	//响应码：0表示成功，否则失败
	State int `json:"state"`
	//说明信息：成功为OK，否则为错误描述
	Msg string `json:"msg"`
}

//接口业务响应
type AuthCheckReply struct {
	//响应码和说明信息
	Status Status `json:"status"`
	//响应数据
	Data AuthCheckData `json:"data"`
}

type AuthCheckData struct {
	UserId string `json:"user_id" form:"user_id"`
}

type RoomBehaviorArgs struct {
	//room id
	RoomId string `json:"room_id" form:"room_id"`
	//行为 1进入 2离开
	Behavior int `json:"behavior" form:"behavior"`
}

type RoomBehaviorReply struct {
	//响应码和说明信息
	Status Status `json:"status"`
}

// LEAVE
type BizLeaveData struct {
	Room string `json:"room"`
}

func BuildWSMessage(msgType int, msgData []byte) (wsMessage *WSMessage) {
	return &WSMessage{
		MsgType: msgType,
		MsgData: msgData,
	}
}

func EncodeWSMessage(bizMessage *BizMessage) (wsMessage *WSMessage, err error) {
	var (
		buf []byte
	)
	if buf, err = json.Marshal(*bizMessage); err != nil {
		return
	}
	wsMessage = &WSMessage{websocket.TextMessage, buf}
	return
}

// 解析{"type": "PING", "data": {...}}的包
func DecodeBizMessage(buf []byte) (bizMessage *BizMessage, err error) {
	var (
		bizMsgObj BizMessage
	)

	if err = json.Unmarshal(buf, &bizMsgObj); err != nil {
		return
	}

	bizMessage = &bizMsgObj
	return
}
