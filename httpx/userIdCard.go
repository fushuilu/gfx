package httpx

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/fushuilu/golibrary/appx/errorx"
	"github.com/fushuilu/golibrary/libx"
	"github.com/fushuilu/golibrary/libx/cachex"

	"github.com/gogf/gf/v2/util/gconv"
)

// 接收到请求的接口，如何判断并得到当前请求用户的身份 (jwt token)
type UserIdCard interface {
	GetIdCard(token string) (IdCard, error) // 必须返回一个身份
	CreateJwtToken(iat int64, card IdCard) (string, error)
	RemoveJwtToken(card IdCard) error
}

func NewJwtUserIdCard(secret string, redis *cachex.XRedis) UserIdCard {
	return &jwtUserIdCard{
		jwt:   cmn.NewJwt(secret),
		redis: redis,
	}
}

type jwtUserIdCard struct {
	jwt   cmn.Jwt
	redis *cachex.XRedis
}

func (j *jwtUserIdCard) key(userid int64, name string) string {
	return fmt.Sprintf("jwt:%d:%s", userid, name)
}

func (j *jwtUserIdCard) GetIdCard(token string) (IdCard, error) {
	if token == "" {
		return IdCard{}, nil
	}
	if claims, err := j.jwt.Decode(token); err != nil {
		return IdCard{}, errorx.New("登录凭证过期或无效，请重新登录", err)
	} else {
		if j.redis != nil {
			key := j.key(gconv.Int64(claims["userid"]), gconv.String(claims["name"]))
			if content, err := j.redis.GetString(key); err != nil {
				return IdCard{}, errorx.New("读取 id 缓存错误")
			} else if content != "1" {
				return IdCard{}, errorx.New("id card 已经过期或不存在")
			}
		}
		return IdCard{UserId: gconv.Int64(claims["userid"]),
			Uid:  gconv.String(claims["uid"]),
			Name: gconv.String(claims["name"]),
		}, nil
	}
}

func (j *jwtUserIdCard) CreateJwtToken(iat int64, card IdCard) (string, error) {
	now := jwt.TimeFunc().Unix()

	claims := make(jwt.MapClaims)
	claims["exp"] = now + iat
	claims["iat"] = iat
	claims["userid"] = card.UserId
	claims["uid"] = card.Uid
	claims["name"] = card.Name

	if token, err := j.jwt.Encode(claims); err != nil {
		return "", err
	} else {
		if j.redis != nil {
			if err = j.redis.SetString(j.key(card.UserId, gconv.String(claims["name"])), "1", iat); err != nil {
				return "", errorx.New("设置 jwt 缓存错误", err, iat)
			}
		}
		return token, nil
	}
}

func (j *jwtUserIdCard) RemoveJwtToken(card IdCard) error {
	if j.redis != nil {
		return j.redis.Delete(j.key(card.UserId, card.Name))
	}
	return nil
}
