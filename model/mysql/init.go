package mysql

import "strings"

type (
	// Field 为字段查询结构体
	Field struct {
		// Operate MySQL中查询条件，如like,=,in
		Operate string
		// Value 查询条件对应的值
		Value interface{}
	}

	// Fields 为Field类型map，用于定义Where方法参数
	Fields map[string]Field
)

var (
	// PageSize 指定了分页大小
	PageSize = 20
)

// Where 返回用于进行where查询的sql和bindParam
func Where(fields Fields) (sql string, bindParam []interface{}) {
	return where(fields, "1=1 ", append(bindParam))
}

func where(fields Fields, sql string, bindParam []interface{}) (retSQL string, retBindParam []interface{}) {
	retSQL = sql
	retBindParam = bindParam

	// TODO 优化
	for name, field := range fields {
		switch strings.ToLower(field.Operate) {
		case "like":
			val := field.Value.(string)
			if field.Value != "" {
				retSQL += " AND `" + name + "` like ? "
				nameParam := "%" + val + "%"
				retBindParam = append(retBindParam, &nameParam)
			}
		case "in", "not in":
			retSQL += " AND `" + name + "` " + field.Operate + " (?) "
			nameParam := field.Value
			retBindParam = append(retBindParam, nameParam)
		default:
			retSQL += " AND `" + name + "` " + field.Operate + " ? "
			nameParam := field.Value
			retBindParam = append(retBindParam, &nameParam)
		}
	}

	return
}
