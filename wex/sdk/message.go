package sdk

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/fushuilu/gfx/wex/sdk/wxbizmsgcrypt"

	"github.com/fushuilu/gfx/httpx"
	"github.com/fushuilu/gfx/logx"
	"github.com/fushuilu/golibrary/appx/errorx"
	"github.com/fushuilu/golibrary/lerror"
)

// https://open.work.weixin.qq.com/api/doc/90000/90135/90240
type WorkMsg struct {
	ToUserName   string `json:"ToUserName"`   // 企业微信 CorpID
	FromUserName string `json:"FromUserName"` // 成员 UserID
	AgentID      int    `json:"AgentID"`      // 企业应用 ID
	CreateTime   int    `json:"CreateTime"`   // 消息创建时间(整型)
	MsgType      string `json:"MsgType"`
	MsgId        string `json:"MsgId"`
	// MsgType = text 文本
	Content string `json:"Content,omitempty"`
	// MsgType = image 图片
	PicUrl  string `json:"PicUrl,omitempty"`  // image/link 封面缩放图的 rul
	MediaId int64  `json:"MediaId,omitempty"` // image/voice/video
	// MsgType = voice 语音
	Format string `json:"Format,omitempty"` // 语音格式
	// MsgType = video
	ThumbMediaId int64 `json:"ThumbMediaId,omitempty"` // 视频消息缩略图媒体
	// MsgType = location
	LocationX float64 `json:"Location_X,omitempty"` // 地理位置纬度
	LocationY float64 `json:"Location_Y,omitempty"` // 地理位置经度
	Scale     int     `json:"Scale,omitempty"`      // 地图缩放大小
	Label     string  `json:"Label,omitempty"`      // 地理位置信息
	// MsgType = link
	Title       string `json:"Title"`       // 标题
	Description string `json:"Description"` // 描述
	Url         string `json:"Url"`         // 链接跳转的 url
	// MsgType = event
	Event string `json:"Event"`
	// event = subscribe(关注)/unsubscribe(取消关注)
	// event = enter_agent 进入应用
	// event = LOCATION 上报地理位置
	Latitude  float64 `json:"Latitude"`  // 纬度
	Longitude float64 `json:"Longitude"` // 经度
	Precision float64 `json:"Precision"` // 精度
	// event = change_contact 新增部门
	ChangeType string `json:"ChangeType"` // create_party|update_party|delete_party
	Id         int64  `json:"Id"`         // 部门 ID
	Name       string `json:"Name"`       // 部门名称
	ParentId   int64  `json:"ParentId"`   // 父部门 ID
	Order      int    `json:"Order"`      // 部门排序
	// 点击菜单拉取消息的事件推送 event=click
	// 点击菜单跳转链接事件推送 event=view
	EventKey string `json:"EventKey"` // 自定义菜单接口中的 KEY 值；跳转 URL
	// ...
	// 追加的数据
	Uid    string `json:"-"`
	UserId int64  `json:"-"`
}

// https://work.weixin.qq.com/api/doc/90000/90135/90930
// 支持 Http Get请求验证URL有效性
type WxCpt struct {
	biz *wxbizmsgcrypt.WXBizMsgCrypt
}

// https://work.weixin.qq.com/api/doc/10514
// https://work.weixin.qq.com/api/doc/90000/90139/90968#%E9%99%84%E6%B3%A8
// receiverId => 1: 企业应用的回调，表示corpid; 2: 第三方事件的回调，表示suiteid
func NewWxCpt(token, encodingAeskey, receiverId string) WxCpt {
	return WxCpt{
		biz: wxbizmsgcrypt.NewWXBizMsgCrypt(token,
			encodingAeskey,
			receiverId, wxbizmsgcrypt.XmlType),
	}
}

type ReqEchostr struct {
	MsgSignature string `json:"msg_signature"`     // "5c45ff5e21c57e6ad56bac8758b79b1d9ac89fd3"
	Timestamp    string `json:"timestamp"`         // "1409659589"
	Nonce        string `json:"nonce"`             // "263014780"
	Echostr      string `json:"echostr,omitempty"` // ""
}

func (pd *ReqEchostr) Invalid() error {
	if pd.MsgSignature == "" {
		return errorx.New("msg signature is empty")
	}
	if pd.Timestamp == "" {
		return errorx.New("timestamp is empty")
	}
	if pd.Nonce == "" {
		return errorx.New("nonce is empty")
	}
	return nil
}

// 验证回调URL
func (w *WxCpt) Verify(pd ReqEchostr) (string, error) {
	echoStr, cryptErr := w.biz.VerifyURL(pd.MsgSignature, pd.Timestamp, pd.Nonce, pd.Echostr)
	if cryptErr != nil {
		logx.Debugf("|<=== 验证 URL 失败: %d : %s", cryptErr.ErrCode, cryptErr.ErrMsg)
		return "", errorx.New("验证 URL 失败")
	}
	return string(echoStr), nil
}

// 对用户回复的消息解密
func (w *WxCpt) Msg(pd ReqEchostr, data []byte) (msg WorkMsg, err error) {
	decryptMsg, cryptError := w.biz.DecryptMsg(pd.MsgSignature, pd.Timestamp, pd.Nonce, data)
	if cryptError != nil {
		logx.Debugf("|<=== 对用户回复的消息解密: %d : %s", cryptError.ErrCode, cryptError.ErrMsg)
		return msg, errorx.New("验证用户回复消息失败")
	}
	//glog.Debug("|<==== 解密数据:", string(msg))
	if err := xml.Unmarshal(decryptMsg, &msg); err != nil {
		return msg, lerror.Wrap(err, "解析消息错误")
	}
	return
}

// https://work.weixin.qq.com/api/doc/90000/90135/90236
const (
	//messageSendURL = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s&debug=1"
	messageSendURL = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s"
)

// 企业回复用户消息的加密
func (w *WxCpt) Send(token string, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return lerror.Wrap(err, "编码企业微信信息错误")
	}

	//glog.Debug("postData:", string(bytes))
	//encryptMsg, cryptError := w.biz.EncryptMsg(string(bytes),
	//	gconv.String(time.Now().Unix()),
	//	grand.Letters(8))
	//if cryptError != nil {
	//	glog.Debugf("|<=== 发送消息失败: %d : %s", cryptError.ErrCode, cryptError.ErrMsg)
	//	return errors.New("加密发送消息失败")
	//}
	// 发送请求
	resp := RespWorkCommon{}
	if err = httpx.NewClientHttp().Post(fmt.Sprintf(messageSendURL, token), string(bytes), &resp); err != nil {
		return lerror.Wrap(err, "发送企业微信消息失败")
	}
	if resp.Errcode != 0 {
		return errorx.New(resp.Errmsg)
	}
	return nil
}
