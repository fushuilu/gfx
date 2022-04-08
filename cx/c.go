package cx

import (
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	isDebug = false
)

func IsDebug() bool {
	return isDebug
}

func SetDebug(open bool) {
	isDebug = open
}

// 默认配置
var ctx = gctx.New()

func Get(pattern string, def ...interface{}) (*gvar.Var, error) {
	return g.Cfg().Get(ctx, pattern, def)
}

func Contains(pattern string) (bool, error) {
	v, err := Get(pattern)
	if err != nil {
		return false, err
	}
	return v.Bool(), nil
}

func ContainsWith(key string) bool {
	data, _ := Contains(key)
	return data
}

func GetArray(pattern string, def ...interface{}) ([]interface{}, error) {
	if v, err := Get(pattern, def); err != nil {
		return nil, err
	} else {
		return v.Interfaces(), nil
	}
}

func GetString(pattern string, def ...interface{}) (string, error) {
	if v, err := Get(pattern, def); err != nil {
		return "", err
	} else {
		return v.String(), nil
	}
}
func GetStringWith(key string, def ...interface{}) string {
	data, _ := GetString(key, def)
	return data
}

func GetStrings(pattern string, def ...interface{}) ([]string, error) {
	if v, err := Get(pattern, def); err != nil {
		return nil, err
	} else {
		return v.Strings(), nil
	}
}

func GetBool(pattern string) (bool, error) {
	if v, err := Get(pattern); err != nil {
		return false, err
	} else {
		return v.Bool(), nil
	}
}

func GetBoolWith(key string) bool {
	data, _ := GetBool(key)
	return data
}
func GetInt(pattern string, def ...interface{}) (int, error) {
	if v, err := Get(pattern, def); err != nil {
		return 0, err
	} else {
		return v.Int(), nil
	}
}

func GetIntWith(key string, def ...interface{}) int {
	data, _ := GetInt(key)
	return data
}

func GetInt64(pattern string, def ...interface{}) (int64, error) {
	if v, err := Get(pattern, def); err != nil {
		return 0, err
	} else {
		return v.Int64(), nil
	}
}

func GetStruct(pattern string, pointer interface{}) error {
	if v, err := Get(pattern); err != nil {
		return err
	} else {
		return v.Struct(pointer)
	}
}

func GetStructs(pattern string, pointer interface{}) error {
	if v, err := Get(pattern); err != nil {
		return err
	} else {
		return v.Structs(pointer)
	}
}
