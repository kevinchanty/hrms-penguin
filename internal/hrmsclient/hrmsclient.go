package hrmsclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
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
	logger     *log.Logger
	httpClient *http.Client
}

type ClientOption struct {
	Host     string
	UserName string
	Pwd      string
	Logger   *log.Logger
}

func New(option ClientOption) *HrmsClient {
	// http client
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{
		Jar:     jar,
		Timeout: 60 * time.Second,
	}

	logger := option.Logger.WithPrefix("HrmsClient")

	return &HrmsClient{host: option.Host, userName: option.UserName, pwd: option.Pwd, httpClient: client, logger: logger}
}

func (c *HrmsClient) Login() {
	c.logger.Debug("Login start")

	formData := url.Values{}
	formData.Set("action", "login")
	formData.Set("fldEmpLoginID", c.userName)
	formData.Set("fldEmpPwd", c.pwd)
	formData.Set("code", "undefined")

	c.logger.Debug("Posting Login...", "formData", formData)
	res, err := c.httpClient.PostForm(fmt.Sprintf("%s/api/admin/login", c.host), formData)
	if err != nil {
		c.logger.Fatal(err)
	}
	defer res.Body.Close()

	c.logger.Debug(res.Header)

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
	c.logger.Debugf("[HrmsClient] GetAction Post %s StatusCode: %d", fmt.Sprintf("%s/api/Home/GetAction", c.host), res.StatusCode)

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

type AttendanceRes struct {
	Msg  string           `json:"msg"`
	Data []AttendanceData `json:"data"`
}

type AttendanceData struct {
	Date            string `json:"fldDate"`
	OriginalInTime  string `json:"fldOriIn1"`
	OriginalOutTime string `json:"fldOriOut1"`
}

func (c *HrmsClient) GetAttendance(year string, month string) ([]AttendanceData, error) {
	endpointUrl := fmt.Sprintf("%s/api/Attendance/getpages", c.host)

	params := url.Values{}
	params.Add("page", "1")
	params.Add("limit", "31")
	params.Add("types", "1")
	params.Add("fldMonth", fmt.Sprintf("%s-%s", year, month))

	fullEndPointUrl := fmt.Sprintf("%s?%s", endpointUrl, params.Encode())

	c.logger.Debug("Fetching Attendance...", "url", fullEndPointUrl)
	resp, err := c.httpClient.Get(fullEndPointUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respStr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	c.logger.Debugf("Fetch Attendance response body: %v", string(respStr))

	var attendanceRes AttendanceRes
	err = json.Unmarshal(respStr, &attendanceRes)
	if err != nil {
		return nil, err
	}

	fmt.Printf("result: %v", attendanceRes)

	return attendanceRes.Data, nil
}
