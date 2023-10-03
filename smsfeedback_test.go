package smsfeedback

import(
	"fmt"
	"io/ioutil"
	"bytes"
	"encoding/json"
	"testing"
)

const TEST_DATA_FILE = ".test_data.json"

type TestData struct {
	Login string `json:"login"`
	Pwd string `json:"pwd"`
	Tel string `json:"tel"`
	SMS struct {
		Tel string `json:"tel"`
		Text string `json:"text"`
		Sign string `json:"sign"`
	} `json:"sms"`
	SMS_ids []string  `json:"sms_ids"`
}

func (d *TestData) loadData() error {
	file, err := ioutil.ReadFile(TEST_DATA_FILE)
	if err == nil {
		file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
		err = json.Unmarshal([]byte(file), d)		
	}
	return err	
}


func TestSendSMS(t *testing.T) {
	d := TestData{}
	d.loadData()
	msg_id, err := SendSMS(d.Login, d.Pwd, d.SMS.Tel, d.SMS.Text, d.SMS.Sign, "")
	if err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Println("msg_id=",msg_id)
}


func TestCheckSMSList(t *testing.T) {
	d := TestData{}
	d.loadData()
	result_list, err := GetDelivered(d.Login, d.Pwd, d.SMS_ids)
	if err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Printf("result_list=%+v\n", result_list)
}


func TestCheckBalance(t *testing.T) {
	d := TestData{}
	d.loadData()	
	res, err := GetBalance(d.Login, d.Pwd)
	if err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Printf("res=%f\n", res)
}


func TestCheckSMS(t *testing.T) {
	d := TestData{}
	d.loadData()	
	res, err := GetSMSDelivered(d.Login, d.Pwd, d.SMS_ids[0])
	if err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Printf("res=%+v\n", res)
}
