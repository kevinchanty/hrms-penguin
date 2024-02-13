package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"

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
	client := &http.Client{Jar: jar, Timeout: 60 * time.Second}

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

func (c *HrmsClient) GetAction() ([]table.Row, error) {
	c.logger.Debug("GetAction Start")

	formData := url.Values{}
	formData.Set("action", "maincontent")
	res, err := c.httpClient.PostForm(fmt.Sprintf("%s/api/Home/GetAction", c.host), formData)
	if err != nil {
		return nil, err
	}
	c.logger.Debugf("StatusCode: %d", res.StatusCode)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	return ParseMainActionForTable(string(body)), nil
}

type Action struct {
	missAttendance []string
	earlyLeave     []string
	lateness       []string
}

// Trying to use HTML parser as it may be useful for other endpoints
func ParseMainAction(actionStr string) *Action {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(actionStr))
	if err != nil {
		log.Fatal(err)
	}

	action := Action{
		missAttendance: make([]string, 0, 31),
		earlyLeave:     make([]string, 0, 31),
		lateness:       make([]string, 0, 31),
	}

	// Find the review items
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		var currentArray *[]string
		s.Contents().Each(func(i int, s *goquery.Selection) {
			if s.Is("br") {
				return
			}

			switch text := strings.TrimSpace(s.Text()); text {
			case "Missing Attendance record 欠缺出入勤紀錄:":
				currentArray = &action.missAttendance
			case "Early leave:":
				currentArray = &action.earlyLeave
			case "Lateness 遲到:":
				currentArray = &action.lateness
			default:
				*currentArray = append(*currentArray, text)
			}
		})
	})

	return &action
}

// todo: combine same date with 2 types, map and transform
// {"2023-12-31". "Missing Attendance record 欠缺出入勤紀錄, Lateness 遲到"}
func ParseMainActionForTable(actionStr string) []table.Row {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(actionStr))
	if err != nil {
		log.Fatal(err)
	}

	rows := []table.Row{}

	// Find the review items
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		var currentType string

		s.Contents().Each(func(i int, s *goquery.Selection) {
			if s.Is("br") {
				return
			}

			switch text := strings.TrimSpace(s.Text()); text {
			case "Missing Attendance record 欠缺出入勤紀錄:":
				currentType = "Missing Attendance record 欠缺出入勤紀錄"
			case "Early leave:":
				currentType = "Early leave"
			case "Lateness 遲到:":
				currentType = "Lateness 遲到"
			default:
				rows = append(rows, table.Row{text, currentType})
			}
		})
	})

	println(rows)
	return rows
}