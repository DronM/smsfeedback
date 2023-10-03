package smsfeedback

import(
	"fmt"
	"errors"
	"strings"
	"net/http"
	"net/url"
	"io/ioutil"
	"strconv"
	"math"
	"encoding/json"
	b64 "encoding/base64"
)

const (
	HOST = "api.smsfeedback.ru/messages/v2"
	PORT = 80
	RESP_DELIM = ";"
	RESP_SENT = "delivered"
	
	ERR_INCOR_NUM = "Неверный номер для SMS сообщения %s"
	ERR_RESP_UNKONWN = "Неверный ответ сервера СМС"
	
	RESP_ACCEPTED = "accepted"	
	RESP_INVALID = "invalid mobile phone"
	RESP_INVALID_DESCR = "Неверно задан номер телефона"	
	RESP_ER = "error authorization"
	RESP_ER_DESCR = "Неверный логин или пароль. Ошибка авторизации"
	RESP_EMPTY = "text is empty"
	RESP_EMPTY_DESCR = "Отсутствует текст"
	RESP_NOT_STR = "text must be string"
	RESP_NOT_STR_DESCR = "Текст не на латинице или не в utf-8"
	REP_SENDER_INVALID = "sender address invalid"
	REP_SENDER_INVALID_DESCR = "Неверная (незарегистрированная) подпись отправителя"
	REP_WAPURL_INVALID = "wapurl invalid"
	REP_WAPURL_INVALID_DESCR = "Неправильный формат wap-push ссылки"
	REP_TIME_INVALID = "invalid schedule time format"
	REP_TIME_INVALID_DESCR = "Неверный формат даты отложенной отправки сообщения"
 	REP_STATUS_INVALID = "invalid status queue name"
 	REP_STATUS_INVALID_DESCR = "Неверное название очереди статусов сообщений"
 	REP_BAL = "not enough balance"
 	REP_BAL_DESCR = "Баланс пуст"
)

/**
 * send request to server
 */
func sendRequest(login, pwd, cmd, params string) (string, error) {
	
	url := fmt.Sprintf("http://%s/%s/", HOST, cmd)
	if params != "" {
		url += "?"+params
	}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	
	if login != "" && pwd != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", login, pwd)))) )
	}
	
	resp, err := client.Do(req)	
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		msg := struct {
			Code string `json:"code"`
			Description string `json:"description"`
			Status string `json:"status"`
		}{}
	
		if err := json.Unmarshal(body, &msg); err != nil {
			return "", err
		}
		return "", errors.New(msg.Description)
	}
	
	return string(body), nil
}

//returns SMS id
func SendSMS(login, pwd, phone, text, sender, wapurl string) (string, error) {
	if phone == "" {
		return "", errors.New(fmt.Sprintf(ERR_INCOR_NUM, phone))
	}
	//correct phone number
	if phone[0:1] == "8" || phone[0:1] == "7"{
		phone = phone[1:]
		
	}else if phone[0:2] == "+7" {
		phone = phone[2:]
	}
	phone = strings.ReplaceAll(phone, "_", "")
	phone = "7" + strings.ReplaceAll(phone, "-", "")
	
	allowed_ext := []string{"900","901","902","903","904","905","906","908","909","910","911","912","913","914","915","916","917","918","919","920","921","922","923","924","925","926","927","928","929","930","931","932","933","934","936","937","938","939","941","950","951","952","953","954","955","956","958","960","961","962","963","964","965","966","967","968","969","970","971","980","981","982","983","984","985","987","988","989","991","992","993","994","995","996","997","999"}
	
	if len(phone) != 11 {
		return "", errors.New(fmt.Sprintf(ERR_INCOR_NUM, phone))
	}
	f := false
	for _, v := range allowed_ext {
		if v == phone[1:4] {
			f = true
			break;
		}
	}
	if !f {
		return "", errors.New(fmt.Sprintf(ERR_INCOR_NUM, phone))
	}
	
	params := "phone=" + rawurlencode(phone) + "&text=" + rawurlencode(text)
	if sender != "" {
		params+= "&sender=" + rawurlencode(sender)
	}
	if wapurl != "" {
		params+= "&wapurl=" + rawurlencode(wapurl)
	}	
	resp, err := sendRequest(login, pwd, "send", params)
	if err != nil {
		return "", err
	}
//fmt.Println(resp)	
	resp_val := strings.Split(resp, ";")
	if len(resp_val) == 2 {
		switch resp_val[0] {
		case RESP_ACCEPTED:
			return resp_val[1], nil
		case RESP_INVALID:
			return "",errors.New(RESP_INVALID_DESCR)
		case RESP_ER:
			return "",errors.New(RESP_ER_DESCR)
		case RESP_EMPTY:
			return "",errors.New(RESP_EMPTY_DESCR)
		case RESP_NOT_STR:
			return "",errors.New(RESP_NOT_STR_DESCR)
		case REP_SENDER_INVALID:
			return "",errors.New(REP_SENDER_INVALID_DESCR)
		case REP_WAPURL_INVALID:
			return "",errors.New(REP_WAPURL_INVALID_DESCR)
		case REP_TIME_INVALID:
			return "",errors.New(REP_TIME_INVALID_DESCR)
	 	case REP_STATUS_INVALID:
	 		return "",errors.New(REP_STATUS_INVALID_DESCR)
	 	case REP_BAL:
	 		return "",errors.New(REP_BAL_DESCR)			
		}
	}
	
	return "",errors.New(ERR_RESP_UNKONWN)
	
}

func GetDelivered(login, pwd string, ids []string) (map[string]bool, error) {
	params := ""
	for _,id := range ids {
		if params != "" {
			params+= "&"			
		}
		params+= "id=" + rawurlencode(id)
	}
	resp, err := sendRequest(login, pwd, "status", params)
	if err != nil {
		return nil, err
	}
//fmt.Println("GetDelivered resp=",resp)	
	resp_list := strings.Split(resp, "\n")
	
	res := make(map[string]bool)
	for _, resp_v := range resp_list {
		resp_val := strings.Split(resp_v, ";")
		if len(resp_val) == 2 {
			var delivered bool
			if resp_val[1] == RESP_SENT {
				delivered = true
			}
			res[resp_val[0]] = delivered
		}
	}
	return res, err
}

func GetSMSDelivered(login, pwd string, id string) (bool, error) {
	resp, err := sendRequest(login, pwd, "status", "id="+id)
	if err != nil {
		return false, err
	}
	resp_list := strings.Split(resp, "\n")
	if len(resp_list)>=1 {
		resp_val := strings.Split(resp_list[0], ";")
		if len(resp_val) == 2 && resp_val[1] == RESP_SENT{
			return true, nil
		}
	}
	return false, nil
}

//returns 
func GetBalance(login, pwd string) (float32, error) {
	resp, err := sendRequest(login, pwd, "balance", "")
	if err != nil {
		return 0.0, err
	}
	resp_list := strings.Split(resp, "\n")
	if len(resp_list) >= 1 {
		for _, resp_v := range resp_list {
			resp_val := strings.Split(resp_v, ";")
			if len(resp_val) >= 2 && resp_val[0] == "RUB" {
				f, err := strconv.ParseFloat(resp_val[1], 32)
				if err != nil {
					return 0.0, err
				}				
				return float32(math.Floor(f*100)/100), nil
			}
		}
	}
	return 0.0, err

}

func rawurlencode(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}
