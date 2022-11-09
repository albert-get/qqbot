package command

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var Token = "Bot 102030298.nokZ62YTReCz94WuwyMwTyop21ajbilV"

func Cluctue(mess map[string]interface{}) {

	d := mess["d"]
	messD := d.(map[string]interface{})
	content := fmt.Sprintf("%v", messD["content"])
	compileRegex := regexp.MustCompile("<@!10021129508834684744> +/算数(.*)")
	matchArr := compileRegex.FindStringSubmatch(content)
	send := map[string]interface{}{}
	defer func() {
		jsonSend, errJson := json.Marshal(send)
		if errJson != nil {
		}
		payload := strings.NewReader(string(jsonSend))
		req, e := http.NewRequest("POST", "https://sandbox.api.sgroup.qq.com/channels/13116632/messages", payload)
		req.Header.Add("Authorization", Token)
		req.Header.Add("accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		http.DefaultClient.Do(req)
		if e != nil {
			fmt.Println(e)
		}
	}()
	defer func() {
		err := recover()
		if err != nil {
			send["content"] = "表达式非法"
			return
		}

	}()

	if matchArr == nil {
		send["content"] = "表达式非法"
		return
	}
	getHou := toHouZui(matchArr[len(matchArr)-1])
	if getHou == nil {
		send["content"] = "表达式非法"
		return
	}
	suanshuMap := map[string]bool{
		"+": true,
		"-": true,
		"*": true,
		"/": true,
	}
	res := []float64{}
	for _, v := range getHou {
		if suanshuMap[v] == true {
			nPrev := res[len(res)-2]
			nNext := res[len(res)-1]
			switch v {
			case "+":
				res = append(res[:len(res)-2], nPrev+nNext)
			case "-":
				res = append(res[:len(res)-2], nPrev-nNext)
			case "*":
				res = append(res[:len(res)-2], nPrev*nNext)
			case "/":
				res = append(res[:len(res)-2], nPrev/nNext)
			}
		} else {
			s, e := strconv.ParseFloat(v, 64)
			if e != nil {
				send["content"] = "我算不出来"
				return
			}
			res = append(res, s)
		}
	}
	send["content"] = strconv.FormatFloat(res[0], 'f', 5, 64)
	return

}

func toHouZui(s string) []string {
	houQ := []string{}
	num := []rune{}
	runes := []rune(s)
	var n string = ""
	suanshuMap := map[rune]bool{
		'+': true,
		'-': true,
		'*': true,
		'/': true,
	}
	kuohaoMap := map[rune]bool{
		'(': true,
		')': true,
		'[': true,
		']': true,
		'{': true,
		'}': true,
	}
	piPeiMap := map[rune]rune{
		'(': ')',
		'[': ']',
		'{': '}',
	}
	numMap := map[rune]bool{
		'0': true,
		'1': true,
		'2': true,
		'3': true,
		'4': true,
		'5': true,
		'6': true,
		'7': true,
		'8': true,
		'9': true,
		'.': true,
	}
	for _, v := range runes {
		if v == ' ' {
			continue
		} else if suanshuMap[v] == true {
			if n != "" {
				houQ = append(houQ, n)
				n = ""
			}

			if len(num) == 0 || num[len(num)-1] == '(' || num[len(num)-1] == '[' || num[len(num)-1] == '{' {
				num = append(num, v)
			} else if (num[len(num)-1] == '*' || num[len(num)-1] == '/') && (v == '+' || v == '-') {
				i := len(num) - 1
				for ; i >= 0; i-- {
					if (num[i] == '*' || num[i] == '/') && (v == '+' || v == '-') {
						houQ = append(houQ, string(num[i]))
					} else {
						break
					}
				}
				if i >= 0 {
					num = append(num[:i+1])
					num = append(num, v)
				} else {
					num = append([]rune{}, v)
				}

			} else {
				num = append(num, v)
			}

		} else if kuohaoMap[v] == true {
			if n != "" {
				houQ = append(houQ, n)
				n = ""
			}
			if v == '(' || v == '[' || v == '{' {
				num = append(num, v)
			} else {
				i := len(num) - 1
				pi := false
				for ; i >= 0; i-- {
					if num[i] == '(' || num[i] == '[' || num[i] == '{' {
						if v == piPeiMap[num[i]] {
							pi = true
							break
						} else {
							fmt.Println("表达式非法")
							return nil
						}
					} else {
						houQ = append(houQ, string(num[i]))
					}
				}
				if pi {
					num = append(num[:i])
				} else {
					fmt.Println("表达式非法")
					return nil
				}

			}
		} else if numMap[v] == true {
			n = n + string(v)
		} else {
			fmt.Println("表达式非法")
			return nil
		}
	}
	if n != "" {
		houQ = append(houQ, n)
		n = ""
	}
	if len(num) > 0 {
		for i := len(num) - 1; i >= 0; i-- {
			if kuohaoMap[num[i]] == true {
				fmt.Println("表达式非法")
				return nil
			} else {
				houQ = append(houQ, string(num[i]))
			}
		}
	}

	return houQ
}
