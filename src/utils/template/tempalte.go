package template

import (
	"bytes"
	"html/template"
)

// 解析到字符串
func File2Str(file string, vals interface{}) ( ret string, err error) {

	tpl, err := template.ParseFiles(file)
	if err != nil {
		return
	}

	bp := new(bytes.Buffer)

	err = tpl.Execute(bp, vals)
	if err != nil {
		return
	}

	ret =  bp.String()
	return
}


