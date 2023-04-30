package msgCmd

import (
	"fmt"
	"github.com/charmbracelet/log"
	"reflect"
	"strconv"
	"strings"
)

// 用于解析聊天命令

func ShouldBind(msg string, obj any) error {
	rv := reflect.ValueOf(obj)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("obj should be a Pointer and not nil")
	}
	cmds := strings.Split(msg, " ")
	if len(cmds) < rv.Elem().NumField() {
		return fmt.Errorf("命令长度不足")
	}
	var index = 0
	for i := 0; i < rv.Elem().NumField(); i++ {
		if index >= len(cmds) {
			return fmt.Errorf("命令长度不足")
		}
		if !rv.Elem().Field(i).CanSet() {
			continue
		}
		switch rv.Elem().Field(i).Type().Kind() {
		case reflect.String:
			rv.Elem().Field(i).SetString(cmds[index])
			index++
		case reflect.Int:
			value, err := strconv.Atoi(cmds[index])
			if err != nil {
				return err
			}
			rv.Elem().Field(i).Set(reflect.ValueOf(value))
			index++
		case reflect.Int64:
			value, err := strconv.ParseInt(cmds[index], 10, 64)
			if err != nil {
				return err
			}
			rv.Elem().Field(i).Set(reflect.ValueOf(value))
			index++
		case reflect.Bool:
			value, err := strconv.ParseBool(cmds[index])
			if err != nil {
				return err
			}
			rv.Elem().Field(i).Set(reflect.ValueOf(value))
			index++
		case reflect.Float64:
			value, err := strconv.ParseFloat(cmds[index], 64)
			if err != nil {
				return err
			}
			rv.Elem().Field(i).Set(reflect.ValueOf(value))
			index++
		case reflect.Slice:
			// 数组 判断长度
			lengthStr := rv.Elem().Type().Field(i).Tag.Get("opq")

			if lengthStr != "" {
				length, err := strconv.Atoi(lengthStr)
				if err != nil {
					return err
				}
				if len(cmds)-index < length {
					return fmt.Errorf("%s 长度不够", rv.Elem().Type().Field(i).Name)
				}

				// 类型变换
				switch rv.Elem().Field(i).Type().Elem().Kind() {
				case reflect.String:
					elem := make([]string, 0)
					for j := 0; j < length; j++ {
						elem = append(elem, cmds[index])
						index++
					}
					rv.Elem().Field(i).Set(reflect.ValueOf(elem))
				case reflect.Int:
					elem := make([]int, 0)
					for j := 0; j < length; j++ {
						v, err := strconv.Atoi(cmds[index])
						if err != nil {
							return err
						}
						elem = append(elem, v)
						index++
					}
					rv.Elem().Field(i).Set(reflect.ValueOf(elem))
				case reflect.Int64:
					elem := make([]int64, 0)
					for j := 0; j < length; j++ {
						v, err := strconv.ParseInt(cmds[index], 10, 64)
						if err != nil {
							return err
						}
						elem = append(elem, v)
						index++
					}
					rv.Elem().Field(i).Set(reflect.ValueOf(elem))
				case reflect.Bool:
					elem := make([]bool, 0)
					for j := 0; j < length; j++ {
						v, err := strconv.ParseBool(cmds[index])
						if err != nil {
							return err
						}
						elem = append(elem, v)
						index++
					}
					rv.Elem().Field(i).Set(reflect.ValueOf(elem))
				case reflect.Float64:
					elem := make([]float64, 0)
					for j := 0; j < length; j++ {
						v, err := strconv.ParseFloat(cmds[index], 64)
						if err != nil {
							return err
						}
						elem = append(elem, v)
						index++
					}
					rv.Elem().Field(i).Set(reflect.ValueOf(elem))
				}
			} else {
				for ; index < len(cmds); index++ {
					// 类型变换
					switch rv.Elem().Field(i).Type().Elem().Kind() {
					case reflect.String:
						rv.Elem().Field(i).Set(reflect.ValueOf(cmds[index]))
					case reflect.Int:
						v, err := strconv.Atoi(cmds[index])
						if err != nil {
							return err
						}
						rv.Elem().Field(i).Set(reflect.ValueOf(v))
					case reflect.Int64:
						v, err := strconv.ParseInt(cmds[index], 10, 64)
						if err != nil {
							return err
						}
						rv.Elem().Field(i).Set(reflect.ValueOf(v))
					case reflect.Bool:
						v, err := strconv.ParseBool(cmds[index])
						if err != nil {
							return err
						}
						rv.Elem().Field(i).Set(reflect.ValueOf(v))
					case reflect.Float64:

						v, err := strconv.ParseFloat(cmds[index], 64)
						if err != nil {
							return err
						}

						rv.Elem().Field(i).Set(reflect.ValueOf(v))
					}
				}

			}

		}
		log.Info(rv.Elem().Field(i).Type().Kind())
	}

	return nil

}
