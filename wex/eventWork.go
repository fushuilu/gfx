package wex

import "github.com/fushuilu/gfx/wex/sdk"

type WorkActionEvent func(msg sdk.WorkMsg) (EventMessageResult, error)

var workEvents = make([]WorkActionEvent, 0)

func BindWorkActionEvent(listener WorkActionEvent) {
	workEvents = append(workEvents, listener)
}

// 处理接收到的企业微信消息
func TriggerWorkActionEvent(msg sdk.WorkMsg) error {
	for i := range workEvents {
		if rst, err := workEvents[i](msg); err != nil {
			return err
		} else if rst.Stop {
			return nil
		}
	}
	return nil
}
