package services

import (
	"bytes"
	crypto_rand "crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/llotosl/manycalls2.0/pkg/requests"
)

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

func Captcha(key string, googlekey string, url string, invisible string, method string) string {
	request := requests.NewRequest(http.Client{})

	if method == "recaptcha2" {

		a, _, _ := request.Get("http://api.captcha.guru/in.php?key="+key+"&method=userrecaptcha&googlekey="+googlekey+"&invisible="+invisible+"&pageurl="+url, map[string]string{})
		re := regexp.MustCompile(`\d+`)
		recaptcha_id := re.FindAllString(string(a), -1)
		fmt.Println(recaptcha_id[0])
		a, _, _ = request.Get("http://api.captcha.guru/res.php?key="+key+"&action=get&id="+recaptcha_id[0], map[string]string{})

		for string(a) == "CAPCHA_NOT_READY" {
			time.Sleep(5 * time.Second)
			a, _, _ = request.Get("http://api.captcha.guru/res.php?key="+key+"&action=get&id="+recaptcha_id[0], map[string]string{})
		}
		fmt.Println(string(a))
		if string(a) == "ERROR_CAPTCHA_UNSOLVABLE" {
			return "ERROR_CAPTCHA_UNSOLVABLE"
		}
		re = regexp.MustCompile(`OK\|(\S+)$`)
		captcha_token := re.FindStringSubmatch(string(a))
		fmt.Println(captcha_token)
		return captcha_token[1]

	} else if method == "hcaptcha" {

		a, _, _ := request.Get("http://api.captcha.guru/in.php?key="+key+"&method=hcaptcha&sitekey="+googlekey+"&pageurl="+url, map[string]string{})
		re := regexp.MustCompile(`\d+`)
		recaptcha_id := re.FindAllString(string(a), -1)
		fmt.Println(recaptcha_id[0])
		a, _, _ = request.Get("http://api.captcha.guru/res.php?key="+key+"&action=get&id="+recaptcha_id[0], map[string]string{})

		for string(a) == "CAPCHA_NOT_READY" {
			time.Sleep(5 * time.Second)
			a, _, _ = request.Get("http://api.captcha.guru/res.php?key="+key+"&action=get&id="+recaptcha_id[0], map[string]string{})

		}

		re = regexp.MustCompile(`OK\|(\S+)$`)
		captcha_token := re.FindStringSubmatch(string(a))
		fmt.Println(captcha_token)
		return captcha_token[1]
	}
	return "error"
}

// MailRu service. Take proxies and CaptchaGuru token.
type MailRu struct {
	token string
}

// NewMailRu create MailRu struct.
func NewMailRu(token string) *MailRu {
	return &MailRu{token: token}
}

// Call make call to phone.
func (m *MailRu) Call(phone string, proxy string) error {
	var b bytes.Buffer
	mailLogin := randStr(10, "ABCDEFGHIJKLMNOPQRSTYVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

	_, clientCookie, err := requests.MakeClient("")
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

	data, contentType, err := requests.MakeBoundary("fdsQEFFJjffjgHkf", dataHead)
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

	body, _, err := request.Post("https://account.mail.ru/api/v1/user/signup", headers, []byte(data))
	if err != nil {
		return err
	}

	fmt.Println(body)

	return nil
}
