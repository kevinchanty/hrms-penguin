package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/charmbracelet/log"
	"golang.org/x/net/publicsuffix"
)

type HrmsClient struct {
	host       string
	userName   string
	pwd        string
	logger     log.Logger
	httpClient *http.Client
}

type ClientOption struct {
	Host     string
	UserName string
	Pwd      string
}

func NewHrmsClient(option ClientOption) *HrmsClient {
	// http client
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{Jar: jar}

	// logger
	logger := log.NewWithOptions(os.Stderr, log.Options{
		Prefix: "HRMS Client:",
		Level:  log.DebugLevel,
	})

	return &HrmsClient{host: option.Host, userName: option.UserName, pwd: option.Pwd, httpClient: client, logger: *logger}
}

func (c *HrmsClient) Login() {
	c.logger.Debug("Login start")

	formData := url.Values{}
	formData.Set("action", "login")
	formData.Set("fldEmpLoginID", c.userName)
	formData.Set("fldEmpPwd", c.pwd)
	formData.Set("code", "undefined")

	res, err := c.httpClient.PostForm(fmt.Sprintf("%s/api/admin/login", c.host), formData)
	if err != nil {
		c.logger.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 && res.Header.Get("Set-Cookie") != "" {
		c.logger.Debug("Login Success")
		c.logger.Debugf("%v", res.Header.Get("Set-Cookie"))
		c.logger.Debugf("%s", c.httpClient.Jar.Cookies(&url.URL{Host: c.host}))
	}

	if err != nil {
		c.logger.Fatal(err)
	}
}

func (c *HrmsClient) GetAction() {
	c.logger.Debug("GetAction Start")

	formData := url.Values{}
	formData.Set("action", "maincontent")
	res, err := c.httpClient.PostForm(fmt.Sprintf("%s/api/Home/GetAction", c.host), formData)
	if err != nil {
		log.Fatal(err)
	}
	c.logger.Debugf("StatusCode: %d", res.StatusCode)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.logger.Fatal(err)
	}
	res.Body.Close()

	c.logger.Infof("%s", body)
}

type Action struct {
	missAttendance []string
	earlyLeave     []string
	lateness       []string
}

// Trying to use HTML parser as it may be useful for other endpoints
func (c *HrmsClient) ParseMainAction(actionStr string) Action {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(actionStr))
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title

		s.Contents().Each(func(i int, s *goquery.Selection) {
			if !s.Is("br") {
				fmt.Println(strings.TrimSpace(s.Text()))
			}
			// if s.Text() == "Miss Attendance" {
			// 	c.missAttendance = append(c.missAttendance, title)
			// } else if s.Text() == "Early Leave" {
			// 	c.earlyLeave = append(c.earlyLeave, title)
			// } else if s.Text() == "Lateness" {
			// 	c.lateness = append(c.lateness, title)
			// }
		})
	})

	return Action{}
}
