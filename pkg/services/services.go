package services

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
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
		i += 1
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
		a, _, _ = Get("http://api.captcha.guru/res.php?key="+key+"&action=get&id="+recaptcha_id[0], map[string]string{})

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

// Create MailRu struct.
func NewMailRu(token string) *MailRu {
	return &MailRu{token: token}
}

func (m *MailRu) Call(phone string, proxy string) error {
	var client *http.Client
	var clientCookie *http.Client

	if proxy != "" {
		// Создание Transport с прокси.
		proxyUrl, err := url.Parse("http://" + proxy)
		if err != nil {
			panic(err)
		}
		transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
		// Создание jar-файла куки.
		jar, err := cookiejar.New(nil)
		if err != nil {
			panic(err)
		}

		client = &http.Client{
			Transport: transport,
		}
		clientCookie = &http.Client{
			Transport: transport,
			Jar:       jar,
		}
	} else {
		// Создание jar-файла куки.
		jar, err := cookiejar.New(nil)
		if err != nil {
			panic(err)
		}

		client = &http.Client{}
		clientCookie = &http.Client{
			Jar: jar,
		}
	}

	request := requests.NewRequest(clientCookie)
}
