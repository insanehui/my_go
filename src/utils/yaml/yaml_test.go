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

func Test_ok(t *testing.T) {
	{
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
name: "db
		`
		log.Println(Ok(yfile))
	}

	{
		yfile := "" // 空字符串，也算是合法的yaml
		log.Println(Ok(yfile))
	}

}

func Test_decode(t *testing.T) {

	{
		// 检查多行字符串能不能解析
		yfile := `
blueprintVersion: 20161109 

# 应用信息
app:
  id: xxxx # 应用id
  vendor: LinkSame # 应用提供商
  version: 1.6 # 应用的版本号

deploy:
  id: uuid
  domain: linksame

# 公有云部署信息
cloud:
  provider: ali # 供应商
  geo: xx # 机房所在地区

SLA: 
  performance: high # 可选值：high, medium, low

k8s: |
 apiVersion: v1
 kind: Namespace
 metadata:
   name: {{.Namespace}}
 ---
 apiVersion: extensions/v1beta1
 kind: Deployment
 metadata:
   name: user
   namespace: {{.Namespace}}
   labels:
     name: user
     namespace: {{.Namespace}}
 spec:
   replicas: 1
   template:
     metadata:
       labels:
         name: user
         namespace: {{.Namespace}}
     spec:
       containers:
       - name: user
         image: weaveworksdemos/user:0.3.0
         ports:
         - containerPort: 80
 ---
 apiVersion: v1
 kind: Service
 metadata:
   name: user
   namespace: {{.Namespace}}
   labels:
     name: user
     namespace: {{.Namespace}}
 spec:
   ports:
     # the port that this service should serve on
   - port: 80
     targetPort: 80
   selector:
     name: user
     namespace: {{.Namespace}}
 ---
 apiVersion: extensions/v1beta1
 kind: Deployment
 metadata:
   name: user-db
   namespace: {{.Namespace}}
   labels:
     name: user-db
     namespace: {{.Namespace}}
 spec:
   replicas: 1
   template:
     metadata:
       labels:
         name: user-db
         namespace: {{.Namespace}}
     spec:
       containers:
       - name: user-db
         image: weaveworksdemos/user-db:0.3.0
         ports:
         - name: mongo
           containerPort: 27017
 ---
 apiVersion: v1
 kind: Service
 metadata:
   name: user-db
   namespace: {{.Namespace}}
   labels:
     name: user-db
     namespace: {{.Namespace}}
 spec:
   ports:
     # the port that this service should serve on
   - port: 27017
     targetPort: 27017
   selector:
     name: user-db
     namespace: {{.Namespace}}
 ---
 apiVersion: extensions/v1beta1
 kind: Deployment
 metadata:
   name: cart
   namespace: {{.Namespace}}
   labels:
     name: cart
     namespace: {{.Namespace}}
 spec:
   replicas: 1
   template:
     metadata:
       labels:
         name: cart
         namespace: {{.Namespace}}
     spec:
       containers:
       - name: cart
         image: weaveworksdemos/cart:0.3.0
         ports:
         - containerPort: 80
 ---
 apiVersion: v1
 kind: Service
 metadata:
   name: cart
   namespace: {{.Namespace}}
   labels:
     name: cart
     namespace: {{.Namespace}}
   annotations:
     prometheus.io/path: "/prometheus"
 spec:
   ports:
     # the port that this service should serve on
   - port: 80
     targetPort: 80
   selector:
     name: cart
     namespace: {{.Namespace}}
 ---
 apiVersion: extensions/v1beta1
 kind: Deployment
 metadata:
   name: cart-db
   namespace: {{.Namespace}}
   labels:
     name: cart-db
     namespace: {{.Namespace}}
 spec:
   replicas: 1
   template:
     metadata:
       labels:
         name: cart-db
         namespace: {{.Namespace}}
     spec:
       containers:
       - name: cart-db
         image: mongo:3.4
         ports:
         - name: mongo
           containerPort: 27017
 ---
 apiVersion: v1
 kind: Service
 metadata:
   name: cart-db
   namespace: {{.Namespace}}
   labels:
     name: cart-db
     namespace: {{.Namespace}}
 spec:
   ports:
     # the port that this service should serve on
   - port: 27017
     targetPort: 27017
   selector:
     name: cart-db
     namespace: {{.Namespace}}
 ---
 apiVersion: extensions/v1beta1
 kind: Deployment
 metadata:
   name: catalogue
   namespace: {{.Namespace}}
   labels:
     name: catalogue
     namespace: {{.Namespace}}
 spec:
   replicas: 1
   template:
     metadata:
       labels:
         name: catalogue
         namespace: {{.Namespace}}
     spec:
       containers:
       - name: catalogue
         image: weaveworksdemos/catalogue:0.2.0
         ports:
         - containerPort: 80
 ---
 apiVersion: v1
 kind: Service
 metadata:
   name: catalogue
   namespace: {{.Namespace}}
   labels:
     name: catalogue
     namespace: {{.Namespace}}
 spec:
   ports:
     # the port that this service should serve on
   - port: 80
     targetPort: 80
   selector:
     name: catalogue
     namespace: {{.Namespace}}
 ---
 apiVersion: extensions/v1beta1
 kind: Deployment
 metadata:
   name: catalogue-db
   namespace: {{.Namespace}}
   labels:
     name: catalogue-db
     namespace: {{.Namespace}}
 spec:
   replicas: 1
   template:
     metadata:
       labels:
         name: catalogue-db
         namespace: {{.Namespace}}
     spec:
       containers:
       - name: catalogue-db
         image: weaveworksdemos/catalogue-db:0.2.0
         env:
           - name: MYSQL_ROOT_PASSWORD
             value: fake_password
           - name: MYSQL_DATABASE
             value: socksdb
         ports:
         - name: mysql
           containerPort: 3306
 ---
 apiVersion: v1
 kind: Service
 metadata:
   name: catalogue-db
   namespace: {{.Namespace}}
   labels:
     name: catalogue-db
     namespace: {{.Namespace}}
 spec:
   ports:
     # the port that this service should serve on
   - port: 3306
     targetPort: 3306
   selector:
     name: catalogue-db
     namespace: {{.Namespace}}
 ---
 apiVersion: extensions/v1beta1
 kind: Deployment
 metadata:
   name: front-end
   namespace: {{.Namespace}}
 spec:
   replicas: 1
   template:
     metadata:
       labels:
         name: front-end
         namespace: {{.Namespace}}
     spec:
       containers:
       - name: front-end
         image: weaveworksdemos/front-end:0.2.0
         resources:
           requests:
             cpu: 100m
             memory: 100Mi
         ports:
         - containerPort: 8079
 ---
 apiVersion: v1
 kind: Service
 metadata:
   name: front-end
   namespace: {{.Namespace}}
   labels:
     name: front-end
     namespace: {{.Namespace}}
 spec:
   type: LoadBalancer
   ports:
   - port: 80
     targetPort: 8079
   selector:
     name: front-end
     namespace: {{.Namespace}}
 ---
 apiVersion: extensions/v1beta1
 kind: Deployment
 metadata:
   name: orders
   namespace: {{.Namespace}}
   labels:
     name: orders
     namespace: {{.Namespace}}
 spec:
   replicas: 1
   template:
     metadata:
       labels:
         name: orders
         namespace: {{.Namespace}}
     spec:
       containers:
       - name: orders
         image: weaveworksdemos/orders:0.3.0
         ports:
         - containerPort: 80
 ---
 apiVersion: v1
 kind: Service
 metadata:
   name: orders
   namespace: {{.Namespace}}
   labels:
     name: orders
     namespace: {{.Namespace}}
   annotations:
     prometheus.io/path: "/prometheus"
 spec:
   ports:
     # the port that this service should serve on
   - port: 80
     targetPort: 80
   selector:
     name: orders
     namespace: {{.Namespace}}
 ---
 apiVersion: extensions/v1beta1
 kind: Deployment
 metadata:
   name: orders-db
   namespace: {{.Namespace}}
   labels:
     name: orders-db
     namespace: {{.Namespace}}
 spec:
   replicas: 1
   template:
     metadata:
       labels:
         name: orders-db
         namespace: {{.Namespace}}
     spec:
       containers:
       - name: orders-db
         image: mongo:3.4
         ports:
         - name: mongo
           containerPort: 27017
 ---
 apiVersion: v1
 kind: Service
 metadata:
   name: orders-db
   namespace: {{.Namespace}}
   labels:
     name: orders-db
     namespace: {{.Namespace}}
 spec:
   ports:
     # the port that this service should serve on
   - port: 27017
     targetPort: 27017
   selector:
     name: orders-db
     namespace: {{.Namespace}}
 ---
 apiVersion: extensions/v1beta1
 kind: Deployment
 metadata:
   name: payment
   namespace: {{.Namespace}}
   labels:
     name: payment
     namespace: {{.Namespace}}
 spec:
   replicas: 1
   template:
     metadata:
       labels:
         name: payment
         namespace: {{.Namespace}}
     spec:
       containers:
       - name: payment
         image: weaveworksdemos/payment:0.3.0
         ports:
         - containerPort: 80
 ---
 apiVersion: v1
 kind: Service
 metadata:
   name: payment
   namespace: {{.Namespace}}
   labels:
     name: payment
     namespace: {{.Namespace}}
 spec:
   ports:
     # the port that this service should serve on
   - port: 80
     targetPort: 80
   selector:
     name: payment
     namespace: {{.Namespace}}
 ---
 apiVersion: extensions/v1beta1
 kind: Deployment
 metadata:
   name: queue-master
   namespace: {{.Namespace}}
   labels:
     name: queue-master
     namespace: {{.Namespace}}
 spec:
   replicas: 1
   template:
     metadata:
       labels:
         name: queue-master
         namespace: {{.Namespace}}
     spec:
       containers:
       - name: queue-master
         image: weaveworksdemos/queue-master:0.3.0
         ports:
         - containerPort: 80
 ---
 apiVersion: v1
 kind: Service
 metadata:
   name: queue-master
   namespace: {{.Namespace}}
   labels:
     name: queue-master
     namespace: {{.Namespace}}
   annotations:
     prometheus.io/path: "/prometheus"
 spec:
   ports:
     # the port that this service should serve on
   - port: 80
     targetPort: 80
   selector:
     name: queue-master
     namespace: {{.Namespace}}
 ---
 apiVersion: extensions/v1beta1
 kind: Deployment
 metadata:
   name: rabbitmq
   namespace: {{.Namespace}}
   labels:
     name: rabbitmq
     namespace: {{.Namespace}}
 spec:
   replicas: 1
   template:
     metadata:
       labels:
         name: rabbitmq
         namespace: {{.Namespace}}
     spec:
       containers:
       - name: rabbitmq
         image: rabbitmq:3
         ports:
         - containerPort: 5672
 ---
 apiVersion: v1
 kind: Service
 metadata:
   name: rabbitmq
   namespace: {{.Namespace}}
   labels:
     name: rabbitmq
     namespace: {{.Namespace}}
 spec:
   ports:
     # the port that this service should serve on
   - port: 5672
     targetPort: 5672
   selector:
     name: rabbitmq
     namespace: {{.Namespace}}
 ---
 apiVersion: extensions/v1beta1
 kind: Deployment
 metadata:
   name: shipping
   namespace: {{.Namespace}}
   labels:
     name: shipping
     namespace: {{.Namespace}}
 spec:
   replicas: 1
   template:
     metadata:
       labels:
         name: shipping
         namespace: {{.Namespace}}
     spec:
       containers:
       - name: shipping
         image: weaveworksdemos/shipping:0.3.0
         ports:
         - containerPort: 80
 ---
 apiVersion: v1
 kind: Service
 metadata:
   name: shipping
   namespace: {{.Namespace}}
   labels:
     name: shipping
     namespace: {{.Namespace}}
   annotations:
     prometheus.io/path: "/prometheus"
 spec:
   ports:
     # the port that this service should serve on
   - port: 80
     targetPort: 80
   selector:
     name: shipping
     namespace: {{.Namespace}}
 
`
		var a struct {
			K8s string
		}
		log.Println(Ok(yfile))
		yaml.Unmarshal([]byte(yfile), &a)
		log.Printf("%+v", a)
	}
}
