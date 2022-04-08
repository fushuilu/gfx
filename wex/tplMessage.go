package wex

import (
	"github.com/fushuilu/golibrary/appx/errorx"
)

// https://work.weixin.qq.com/api/doc/90000/90135/90236
type WorkMsgData struct {
	Touser               string                        `json:"touser"`                            // 成员 ID 列表 UserID1|UserID2|UserID3
	Toparty              string                        `json:"toparty,omitempty"`                 // 部门 ID 列表 PartyID1|PartyID2
	Totag                string                        `json:"totag,omitempty"`                   // 标签 ID 列表 TagID1 | TagID2
	Msgtype              string                        `json:"msgtype"`                           // 消息类型，固定为 text|image
	Agentid              string                        `json:"agentid,omitempty"`                 // 企业应用 ID
	Text                 *WorkMsgOfText                `json:"text,omitempty"`                    // 文本消息，支持 \n 和 a 标签
	Image                *WorkMsgOfImage               `json:"image,omitempty"`                   // 图片消息
	Voice                *WorkMsgVoice                 `json:"voice,omitempty"`                   // 语音消息
	Video                *WorkMsgVideo                 `json:"video,omitempty"`                   // 视频消息
	File                 *WorkMsgOfFile                `json:"file,omitempty"`                    // 文件消息
	Textcard             *WorkMsgOfTextcard            `json:"textcard,omitempty"`                // 文本卡片消息
	News                 *WorkMsgOfNews                `json:"news,omitempty"`                    // 图文消息
	Mpnews               *WorkMsgOfMpnews              `json:"mpnews,omitempty"`                  // 图文消息
	Markdown             *WorkMsgOfMarkdown            `json:"markdown,omitempty"`                // Markdown
	MiniprogramNotice    *WorkMsgOfMiniprogramNotice   `json:"miniprogram_notice ,omitempty"`     // 小程序通知消息
	InteractiveTaskcard  *WorkMsgOfInteractiveTaskcard `json:"interactive_taskcard ,omitempty"`   // 模板卡片
	Safe                 int                           `json:"safe ,omitempty"`                   // 是否保密信息
	EnableIdTrans        int                           `json:"enable_id_trans ,omitempty"`        // 是否开启 id 转译
	EnableDuplicateCheck int                           `json:"enable_duplicate_check ,omitempty"` // 是否开启重复消息检查
}

type WorkMsgOfImage struct {
	MediaId string `json:"media_id"` // 图片媒体文件id，可以调用上传临时素材接口获取
}
type WorkMsgOfText struct {
	Content string `json:"content"` // 消息内容
}
type WorkMsgOfInteractiveTaskcard struct {
	Title       string `json:"title"`       // 标题 128
	Description string `json:"description"` // 描述 512
	Url         string `json:"url"`         // 链接 2048
	TaskId      string `json:"task_id"`     // 任务 ID 128
	Btn         struct { // 按钮列表 1~2
		Key    string `json:"key"`     // 回调事件名称 128
		Name   string `json:"name"`    // 按钮名称
		Color  string `json:"color"`   // 按钮字体颜色 red/blue(默认)
		IsBold bool   `json:"is_bold"` // 加粗 (false)
	} `json:"btn"`
}
type WorkMsgOfMiniprogramNotice struct { // 小程序通知消息
	Appid             string `json:"appid"`               // 小程序 appid
	Page              string `json:"page"`                // 小程序页面
	Title             string `json:"title"`               // 标题 4~12
	Description       string `json:"description"`         // 描述 4~12
	EmphasisFirstItem bool   `json:"emphasis_first_item"` // 是否放大第一个 content_item
	ContentItem       []struct { // 消息内容，最多允许 10 个 item
		Key   string `json:"key"`   // ~10
		Value string `json:"value"` // ~30
	} `json:"content_item"`
}
type WorkMsgOfMarkdown struct {
	Content string `json:"content"` // markdown内容，最长不超过2048个字节，必须是utf8编码
}
type WorkMsgOfMpnews struct { // 图文消息
	Articles []struct { // 1~8 条
		Title            string `json:"title"`              // 标题 128
		ThumbMediaId     string `json:"thumb_media_id"`     // 图文消息缩略图的 media_id, 可以通过素材管理接口获得。此处 thumb_media_id 即上传接口返回的media_id
		Author           string `json:"author"`             // 作者 64
		ContentSourceUrl string `json:"content_source_url"` // 点击"阅读原文"之后的页面链接
		Content          string `json:"content"`            // 图文消息的内容，support html, 666K
		Digest           string `json:"digest,omitempty"`   // 图文消息描述 512
	} `json:"articles"`
}
type WorkMsgOfNews struct { // 图文消息
	Articles []struct { // 1~8 条图文
		Title       string `json:"title"`       // 标题 128
		Description string `json:"description"` // 描述 512
		Url         string `json:"url"`         // 链接  2048
		Picurl      string `json:"picurl"`      // 图片链接 jpg/png, 1068*455 或 150*150
	} `json:"articles"`
}
type WorkMsgOfTextcard struct { // 卡片
	Title       string `json:"title"`       // 标题 128
	Description string `json:"description"` // 描述 512
	Url         string `json:"url"`         // 跳转链接 2048
	Btntxt      string `json:"btntxt"`      // 详情 4
}
type WorkMsgOfFile struct { // 文件
	MediaId string `json:"media_id"` // 文件id，可以调用上传临时素材接口获取
}
type WorkMsgVideo struct { // 视频
	MediaId     string `json:"media_id"`    // 视频媒体文件id，可以调用上传临时素材接口获取
	Title       string `json:"title"`       // 视频消息的标题（128）
	Description string `json:"description"` // 视频消息的描述（512）
}
type WorkMsgVoice struct {
	MediaId string `json:"media_id"` // 语音文件id，可以调用上传临时素材接口获取
}

func (pd *WorkMsgData) Invalid() error {
	if pd.Agentid == "" {
		return errorx.New("必须填写 agentid")
	}
	if pd.Msgtype == "" {
		pd.Msgtype = "text"
	}
	if pd.Msgtype == "text" {
		if pd.Text.Content == "" {
			return errorx.New("消息内容不能为空")
		}
	}
	if pd.Touser == "" && pd.Toparty == "" && pd.Totag == "" {
		return errorx.New("接收对象全部为空")
	}
	return nil
}

// 微信模板消息项
type TplDataItem struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
}

type TplMsgDataReceiver struct {
	Appid  string // 微信应用 appid
	Openid string // 用户的 openid
	Url    string // 模板链接
}

// https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html
type TplMsgData struct {
	Appid  string // 微信应用 appid
	Openid string // 用户的 openid

	TemplateID string // 模板 ID

	Url  string // 模板链接
	Data map[string]TplDataItem

	MiniProgram TplMsgDataMiniItem
}

type TplMsgDataMiniItem struct {
	AppID    string `json:"appid"`
	PagePath string `json:"pagepath"`
}

func (pd *TplMsgData) Invalid() error {
	if pd.Appid == "" {
		return errorx.New("必须指定微信 appid")
	}
	if pd.Openid == "" {
		return errorx.New("必须指定接收用户 ID 或者 Openid")
	}
	if pd.TemplateID == "" {
		return errorx.New("必须指定消息模板名称 或者 ID")
	}
	return nil
}

const (
	MiniStateDeveloper = "developer" // 开发版
	MiniStateTrial     = "trial"     // 体验版
	MiniStateFormal    = "formal"    // 正式版
)

// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/subscribe-message/subscribeMessage.send.html
type MiniTplMsgData struct {
	Appid            string                 // 小程序 appid
	Openid           string                 // 小程序用户的 openid
	TemplateID       string                 // 小程序的模板 ID
	Page             string                 // 小程序页面路径
	Data             map[string]TplDataItem // 小程序模板数据
	MiniprogramState string                 // 小程序版本 developer(开发版)|trial(体验版)|formal(正式版默认)
}

func (pd *MiniTplMsgData) Invalid() error {
	if pd.Appid == "" {
		return errorx.New("请填写小程序 appid")
	}
	if pd.Openid == "" {
		return errorx.New("请提供小程序用户 openid")
	}
	if pd.TemplateID == "" {
		return errorx.New("请提供小程序模板 ID")
	}
	if pd.Page == "" {
		return errorx.New("请提供小程序页面路径")
	}

	if len(pd.Data) < 1 {
		return errorx.New("请提供小程序模板数据")
	}

	return nil
}

/**
微信发送模板消息, access token 需要 ip 白名单
*/
type TplMessage interface {
	// https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html
	PostTplMsg(data TplMsgData) (int64, error) // 发送模板消息给指定的用户
	// https://work.weixin.qq.com/api/doc/90000/90135/90236
	PostWorkMsg(data WorkMsgData) error // 发送企业微信消息
	// https://developers.weixin.qq.com/miniprogram/dev/framework/open-ability/subscribe-message.html
	// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/subscribe-message/subscribeMessage.send.html
	PostMiniTplMsg(msg MiniTplMsgData) (bool, error) // 发送小程序订阅消息
}

