package yaml

import (
	"log"
	"testing"
	J "utils/json"

	"github.com/ghodss/yaml"
	gy "gopkg.in/yaml.v2"
)

func TestFromFile(t *testing.T) {
	data := FromFile("test.yaml")
	log.Printf("%+v", data)
}

func Test2Json(t *testing.T) {
	{
		data := FromFile("test.yaml")
		str := J.ToJson(data)
		log.Printf("json obj: %+v", str)
	}

	{
		data := FromFile("list.yaml")
		str := J.ToJson(data)
		log.Printf("json arr: %+v", str)
	}
}

func TestStruct(t *testing.T) {
	type MyType struct {
		Name string `json:"name" mapstructure:"name"`
		Type string `json:"type" mapstructure:"type"`
		Desc string `json:"description" mapstructure:"description"`
	}
	var m MyType

	var data = `
name: name
type: string haha xx
description: this's your name
`
	err := yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Printf("err: %+v", err)
	}
	log.Printf("%+v", m)
	j := J.ToJson(m)
	log.Printf("%+v", j)
}

func Test2Slice(t *testing.T) {
	type MyType struct {
		Name string `yaml:"name" mapstructure:"name"`
		Type string `yaml:"type" mapstructure:"type"`
		Desc string `yaml:"description" mapstructure:"description"`
	}

	var m []MyType
	var data = `
- name: name
  type: string haha xx
  description: this's your name
- name: age
  type: int
  description: how old are you
- name: company
  type: string
  description: where you work
- name: email
  type: email
  description: Arbitrary key/value metadata
`
	err := yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Printf("err: %+v", err)
	}
	log.Printf("%+v", m)
	j := J.ToJson(m)
	log.Printf("%+v", j)

}

func Test_go_yaml(t *testing.T) {

	yfile := `
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: db2
  namespace: linksame-simplest
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 5Gi
  selector:
    matchLabels:
      name: "db"
`
	var d map[string]interface{}

	gy.Unmarshal([]byte(yfile), &d)

	s := d["spec"]
	if s, ok := s.(map[interface{}]interface{}); ok {
		sel := s["selector"]
		log.Printf("h1: %+v", sel)
		if sel, ok := sel.(map[interface{}]interface{}); ok {
			m := sel["matchLabels"]
			log.Printf("h2: %+v", m)
			if m, ok := m.(map[interface{}]interface{}); ok {
				m["xxx"] = 1234214
			}
		}
	}

	log.Printf("h3: %+v", d)

}
