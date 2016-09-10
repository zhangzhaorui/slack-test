package main

import (
	"fmt"
	//"github.com/asiainfoLDP/datafoundry_proxy/env"
	//log "github.com/asiainfoLDP/datafoundry_slack/log"
	"encoding/json"
	"github.com/asiainfoLDP/datafoundry_slack/ds"
	oshandler "github.com/asiainfoLDP/datafoundry_slack/handler"
	"github.com/julienschmidt/httprouter"
	"github.com/openshift/origin/pkg/cmd/util/tokencmd"
	kclient "k8s.io/kubernetes/pkg/client/unversioned"
	"net/http"
	"strings"
	"time"
	//"github.com/astaxie/beego/logs/es"
	//"github.com/astaxie/beego/logs/es"
	"bytes"
	"io/ioutil"
)

const (
	//	DATAFOUNDRY_HOST_ADDR  string = "DATAFOUNDRY_HOST_ADDR"
	//	DATAFOUNDRY_ADMIN_USER string = "DATAFOUNDRY_ADMIN_USER"
	//	DATAFOUNDRY_ADMIN_PASS string = "DATAFOUNDRY_ADMIN_PASS"
	//	DATAFOUNDRY_API_ADDR   string = "DATAFOUNDRY_API_ADDR"
	//
	SERVICE_PORT = "0.0.0.0:8080"
)

var (
//theOC *OpenshiftClient

//DataFoundryEnv = &EnvOnce{
//	envs: map[string]string{
//		DATAFOUNDRY_HOST_ADDR:  "dev.dataos.io:8443",
//		DATAFOUNDRY_ADMIN_USER: "",
//		DATAFOUNDRY_ADMIN_PASS: "",
//		DATAFOUNDRY_API_ADDR:   "",
//	},
//}
)

type OpenshiftClient struct {
	host string
	//authUrl string
	oapiUrl string
	kapiUrl string

	namespace   string
	username    string
	password    string
	bearerToken string
}

//type EnvOnce struct {
//	envs map[string]string
//	once sync.Once
//}

//func Init(ocEnv env.Env) {
//	theOC = newOpenshiftClient(
//		ocEnv.Get("DATAFOUNDRY_HOST_ADDR"),
//		ocEnv.Get("DATAFOUNDRY_ADMIN_USER"),
//		ocEnv.Get("DATAFOUNDRY_ADMIN_PASS"),
//		ocEnv.Get("NAMESPACE"),
//	)
//
//	fmt.Println("username:", theOC.username)
//	fmt.Println("password:", theOC.password)
//}

func newOpenshiftClient(host, username, password, namespace string) *OpenshiftClient {
	host = httpsAddrMaker(host)
	oc := &OpenshiftClient{
		host: host,
		//authUrl: host + "/oauth/authorize?response_type=token&client_id=openshift-challenging-client",
		oapiUrl: host + "/oapi/v1",
		kapiUrl: host + "/api/v1",

		namespace: namespace,
		username:  username,
		password:  password,
	}

	go oc.updateBearerToken()

	return oc
}

func (oc *OpenshiftClient) updateBearerToken() {
	for {
		clientConfig := &kclient.Config{}
		clientConfig.Host = oc.host
		clientConfig.Insecure = true
		//clientConfig.Version =

		println("Request Token from: ", clientConfig.Host)

		token, err := tokencmd.RequestToken(clientConfig, nil, oc.username, oc.password)
		if err != nil {
			println("RequestToken error: ", err.Error())

			time.Sleep(15 * time.Second)
		} else {
			//clientConfig.BearerToken = token
			oc.bearerToken = "Bearer " + token

			println("RequestToken token: ", token)

			time.Sleep(3 * time.Hour)
		}

		fmt.Println("token:", token)
	}
}

func httpsAddrMaker(addr string) string {
	if strings.HasSuffix(addr, "/") {
		addr = strings.TrimRight(addr, "/")
	}

	if !strings.HasPrefix(addr, "https://") {
		return fmt.Sprintf("https://%s", addr)
	}

	return addr
}

//var load bool = true

//type eventStatus struct {
//	load bool
//	events []ds.Event
//}

//var es *eventStatus

func getEvent(watchType string) {
	//es.load = true

	//events := make([]ds.Event, 0)

	oc := oshandler.NewOpenshiftClient("Bearer h4IqusemStY7JIPgVL0fmUBxYf2Vo2S8n1s3hJliuSY")
	fmt.Println("new client.....")
	uri := "/namespaces/team12/" + watchType
	watchStatus, _, err := oc.KWatch(uri)
	if err != nil {
		fmt.Println("KWatch err:", err)
		return
	}

	go func() {
		fmt.Println("Watch......")
		event := ds.Event{}
		//object := ds.ResourcequotasObject{}
		//event.Object = &object

		for {
			status, _ := <-watchStatus
			if status.Err == nil {
				//load = true
				switch watchType {
				case "events":
					object := ds.EventObject{}
					event.Object = &object
					fmt.Println("watch type:", watchType)
					err := json.Unmarshal(status.Info, &event)
					if err != nil {
						fmt.Println("Unmarshal err:", err)
						return
					}

					if event.Type == "ADDED" && event.Object.(*ds.EventObject).Reason == "Started" && judgeEvent(event) {
						sendToSlack(event)
					}

				case "resourcequotas":
					object := ds.ResourcequotasObject{}
					event.Object = &object
					fmt.Println("watch type:", watchType)
					json.Unmarshal(status.Info, &event)
					if err != nil {
						fmt.Println("Unmarshal err:", err)
						return
					}
					//fmt.Println(event.Object)
					sendToSlack(event)
				}

				//events = append(events, event)
				//if es.load {
				//	events = append(events, event)
				//}

				//if event.Type == "ADDED" && event.Object.Reason == "Started" && judgeEvent(event) {
				//	//sendToSlack(event)
				//}
			}

			//if es.load == false && len(es.events) != 0 {
			//	fmt.Println(len(es.events), es.events[len(es.events)-1].Object.Reason)
			//	//fmt.Println(es.events)
			//	es.events = make([]ds.Event, 0)
			//	fmt.Println(len(es.events))
			//}
		}
	}()
}

//func read()  {
//	for {
//		if es.load == false && len(es.events) != 0 {
//			fmt.Println(len(es.events), es.events[len(es.events)-1].Object.Reason)
//			//fmt.Println(es.events)
//			es.events = make([]ds.Event, 0)
//			fmt.Println(len(es.events))
//		}
//	}
//
//}

func judgeEvent(event ds.Event) bool {
	pod := event.Object.(*ds.EventObject).InvolvedObject.Name
	strs := strings.Split(pod, "-")
	if strs[len(strs)-1] != "deploy" {
		return true
	} else {
		return false
	}
}

func sendToSlack(event ds.Event) {
	fmt.Println("begin sendToSlack......")
	defer fmt.Println("end sendToSlack......")

	msg := ds.SlackMsg{Channel: "#datafoundry", Username: "webhookbot", Icon_emoji: ":ghost:"}

	//fmt.Println(event.Object.(*ds.ResourcequotasObject).Metadata)

	//value, ok := event.Object.(*ds.ResourcequotasObject)
	//fmt.Println(ok)
	//if ok {
	//	fmt.Println("this is ResourcequotasObject")
	//		text := fmt.Sprintf("project: %s\ntotal(cpu: %s, memory: %s)\nused(cpu: %s, memory: %s)\nrest: ",
	//			value.Metadata.Namespace,
	//			value.Spec.Hard.RequestsCpu, value.Spec.Hard.RequestsMemory,
	//			value.Status.Used.LimitsCpu, value.Status.Used.LimitsMemory)
	//		msg.Text = text
	//}

	switch t := event.Object.(type) {
	case *ds.ResourcequotasObject:
		fmt.Println("this is ResourcequotasObject")
		text := fmt.Sprintf("project: %s\ntotal(cpu: %s, memory: %s)\nused(cpu: %s, memory: %s) ",
			t.Metadata.Namespace,
			t.Spec.Hard.RequestsCpu, t.Spec.Hard.RequestsMemory,
			t.Status.Used.LimitsCpu, t.Status.Used.LimitsMemory)
		msg.Text = text
	case *ds.EventObject:
		fmt.Println("this is EventObject")
		text := fmt.Sprintf("projiect: %s\npod: %s\nhost: %s\nbuild succeed!",
			t.InvolvedObject.Namespace,
			t.InvolvedObject.Name,
			t.Source.Host)
		msg.Text = text
	}

	//text := fmt.Sprintf("projiect: %s\npod: %s\nhost: %s\nbuild succeed!",
	//	event.Object.(ds.EventObject).InvolvedObject.Namespace,
	//	event.Object.(ds.EventObject).InvolvedObject.Name,
	//	event.Object.(ds.EventObject).Source.Host)
	//msg.Text = text

	body, err := json.Marshal(&msg)

	client := &http.Client{}
	url := "https://hooks.slack.com/services/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"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("http.NewRequest err:", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("client.Do err:", err)
		return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll err:", err)
		return
	}
	fmt.Println(string(respBody))
}

//func computerRest(event ds.Event)  {
//	totalCup := event.Object.(ds.ResourcequotasObject).Spec.Hard.RequestsCpu
//	for i := 0; i <= len(totalCup);i++  {
//		if totalCup[i] > 39 {
//			totalCup = totalCup[:i]
//		}
//	}
//}

func main() {

	oshandler.Init(DataFoundryEnv)

	//time.Sleep(time.Second * 20)
	//fmt.Println("week......")

	//fmt.Println(theOC)
	//fmt.Println(theOC.bearerToken)

	go getEvent("resourcequotas")

	go getEvent("events")

	router := httprouter.New()
	//router.GET("/alarm", rootHandler)
	//router.POST("/alarm", sendMessageHandler)

	//log.Info("listening on", SERVICE_PORT)
	err := http.ListenAndServe(SERVICE_PORT, router)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
		return
	}

	//log.Println("watchStatus:", watchStatus, "  info:", info)
}

//func init() {
//	es = new(eventStatus)
//	es.events = make([]ds.Event, 0)
//	es.load = true
//
//	fmt.Println(len(es.events))
//}
