package session

import (
	"reflect"
)

type Session struct {
	// tmux or zoxide
	Src string
	// The display name
	Name string
	// The absolute directory path
	Path string
}

type Srcs struct {
	Tmux   bool
	Zoxide bool
}

func checkAnyTrue(s interface{}) bool {
	val := reflect.ValueOf(s)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.Bool && field.Bool() {
			return true
		}
	}
	return false
}
