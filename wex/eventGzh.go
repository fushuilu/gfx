package wex

import (
	"github.com/silenceper/wechat/v2/officialaccount"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/silenceper/wechat/v2/officialaccount/user"
)

// 获取公众号用户信息
type GzhGetUserInfo func(gzh *officialaccount.OfficialAccount, openId string) (userInfo *user.Info, err error)

func DefaultGzhGetUserInfo(gzh *officialaccount.OfficialAccount, openId string) (userInfo *user.Info, err error) {
	return gzh.GetUser().GetUserInfo(openId)
}

func ResponseText(text string) *message.Reply {
	return &message.Reply{
		MsgType: message.MsgTypeText,
		MsgData: message.NewText(text),
	}
}



// 微信消息推送
const (
	MsgTypeText       = "text"
	MsgTypeImage      = "image"
	MsgTypeVoice      = "voice"
	MsgTypeVideo      = "video"
	MsgTypeShortVideo = "shortvideo"
	MsgTypeLocation   = "location"
	MsgTypeLink       = "link"
	MsgTypeEvent      = "event"

	EventSubscribe   = "subscribe" // 关注公众号，或描述带参数二维码事件
	EventUnsubscribe = "unsubscribe"
	EventClick       = "CLICK" // 点击菜单
	EventView        = "VIEW"  // 点击菜单跳转
	EventSCAN        = "SCAN"
	EventLocation    = "LOCATION"              // latitude, longitude, precision 存放在上面的 LocationX
	EventTempMsg     = "TEMPLATESENDJOBFINISH" // 模板消息 msgId, status=success

	RespTypeText = "text"
	RespTypeNews = "news" // 图片消息
	// 以下格式通常不用于回复，因为使用的是 mediaId
	RespTypeImage = "image"
	RespTypeVoice = "voice"
	RespTypeVideo = "video"
	RespTypeMusic = "music"

	SUCCESS = "success"
)

type GzhMessageNews struct {
	Title       string
	Description string
	PicUrl      string // 图片链接，支持JPG、PNG格式，较好的效果为大图360*200，小图200*200
	Url         string
}

// 微信订阅事件
type GzhMessage struct {
	UserId  int64  // 用户 id
	Appid   string // 微信应用 id
	Openid  string
	Unionid string
	// 消息类型
	MsgType string
	// 普通消息
	MsgID        int64
	Content      string  // text
	MediaID      string  // image | voice | video | shortvideo
	PicURL       string  // image
	Format       string  // voice 语音格式
	Recognition  string  // 语音识别结果 UTF8 编码
	ThumbMediaId string  // video | shortvideo 视频缩略图
	LocationX    float64 // location 纬度
	LocationY    float64 // location 经度
	Scale        float64 // 地图缩放大小
	Precision    float64 // 地理位置经度
	Label        string  // 地理位置信息
	Title        string  // link 链接标题
	Description  string  // 描述
	URL          string  // 链接
	// 事件相关
	Event       string
	EventKey    string // 事件 key 值
	EventTicket string // 二维码的 ticket ，用于换取二维码图片
	Status      string
}

type EventMessageResult struct {
	Stop    bool        // 停止处理
	MsgType string      // 响应的消息类型
	MsgData interface{} // 响应的内容
}

type GzhMessageEvent func(pd GzhMessage) (EventMessageResult, error)

// 公众号接收消息回调
var messageEvents = make([]GzhMessageEvent, 0)

// 添加微信回调函数
func BindGzhMessageEvent(event GzhMessageEvent) {
	messageEvents = append(messageEvents, event)
}

// msg.MsgType == message.MsgTypeEvent && msg.Event == message.EventSubscribe {
func TriggerGzhMessageEvent(pd GzhMessage) (EventMessageResult, error) {
	for i := range messageEvents {
		if rst, err := messageEvents[i](pd); err != nil {
			return rst, err
		} else if rst.Stop {
			return rst, nil
		}
	}
	return EventMessageResult{ // 空字符串，不处理
		MsgType: RespTypeText, MsgData: "感谢关注",
	}, nil
}
