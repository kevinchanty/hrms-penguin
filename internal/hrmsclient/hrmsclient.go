package hrmsclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"slices"
	"strconv"
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

type HrmsConfig struct {
	Host     string `json:"hrmsHost"`
	UserName string `json:"hrmsUser"`
	Pwd      string `json:"-"`
}

type NewClientOption struct {
	HrmsConfig
	Logger *log.Logger
}

func New(option NewClientOption) *HrmsClient {
	// http client
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{
		Jar:     jar,
		Timeout: 60 * time.Second,
	}

	logger := option.Logger.WithPrefix("[HrmsClient]")

	return &HrmsClient{host: option.Host, userName: option.UserName, pwd: option.Pwd, httpClient: client, logger: logger}
}

func (c *HrmsClient) Login() error {
	c.logger.Debug("Login start")

	formData := url.Values{}
	formData.Set("action", "login")
	formData.Set("fldEmpLoginID", c.userName)
	formData.Set("fldEmpPwd", c.pwd)
	formData.Set("code", "undefined")

	c.logger.Debug("Posting Login...", "formData", formData)
	res, err := c.httpClient.PostForm(fmt.Sprintf("%s/api/admin/login", c.host), formData)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	c.logger.Debug("Login Response:", "headers", res.Header)

	if res.StatusCode == 200 && res.Header.Get("Set-Cookie") != "" {
		c.logger.Debug("Login Success")
		c.logger.Debugf("%v", res.Header.Get("Set-Cookie"))
		c.logger.Debugf("%s", c.httpClient.Jar.Cookies(&url.URL{Host: c.host})) // todo: not working
	} else {
		return errors.New("login failed")
	}
	return nil
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
	DateStr            string `json:"fldDate"`
	OriginalInTimeStr  string `json:"fldOriIn1"`
	OriginalOutTimeStr string `json:"fldOriOut1"`
	Date               time.Time
	OriginalInTime     time.Time
	OriginalOutTime    time.Time
	IsLate             bool
}

func (c *HrmsClient) FetchAttendance(year string, month string) ([]AttendanceData, error) {
	endpointUrl := fmt.Sprintf("%s/api/Attendance/getpages", c.host)

	params := url.Values{}
	params.Add("page", "1")
	params.Add("limit", "31")
	params.Add("types", "1")
	params.Add("fldMonth", fmt.Sprintf("%s-%s", year, month))

	fullEndPointUrl := fmt.Sprintf("%s?%s", endpointUrl, params.Encode())

	c.logger.Debug("[FetchAttendance] fetching endpoint...", "url", fullEndPointUrl)
	resp, err := c.httpClient.Get(fullEndPointUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respStr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	c.logger.Debug("[FetchAttendance]", "response body", string(respStr))

	var attendanceRes AttendanceRes
	err = json.Unmarshal(respStr, &attendanceRes)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("[FetchAttendance] result", "data", attendanceRes)

	// Parse date, in/out time, add IsLate
	for i, attendanceData := range attendanceRes.Data {
		var (
			inTime  time.Time
			outTime time.Time
			date    time.Time
		)

		date, err = time.Parse("2006-01-02", attendanceData.DateStr)
		if err != nil {
			return nil, err
		}

		if attendanceData.OriginalInTimeStr != "" {
			inTime, err = time.Parse("2006-01-02 15:04", fmt.Sprintf("%v %v", attendanceData.DateStr, attendanceData.OriginalInTimeStr))
			if err != nil {
				return nil, err
			}
		}

		if attendanceData.OriginalOutTimeStr != "" {
			outTime, err = time.Parse("2006-01-02 15:04", fmt.Sprintf("%v %v", attendanceData.DateStr, attendanceData.OriginalOutTimeStr))
			if err != nil {
				return nil, err
			}
		}

		attendanceRes.Data[i].OriginalInTime = inTime
		attendanceRes.Data[i].OriginalOutTime = outTime
		attendanceRes.Data[i].Date = date

		if !inTime.IsZero() {
			lateTime, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%v %v", attendanceData.
				DateStr, "09:30"))
			if err != nil {
				return nil, err
			}
			if inTime.After(lateTime) {
				attendanceRes.Data[i].IsLate = true
			}
		}
	}
	return attendanceRes.Data, nil
}

func (c *HrmsClient) GetTodayAttendance() (AttendanceData, error) {
	currentTime := time.Now()
	day := currentTime.Day()
	var accountingMonth time.Month

	if day <= 15 {
		accountingMonth = currentTime.Month()
	} else {
		accountingMonth = currentTime.Month() + 1
	}

	attendanceDataList, err := c.FetchAttendance(strconv.Itoa(currentTime.Year()), strconv.Itoa(int(accountingMonth)))
	if err != nil {
		return AttendanceData{}, err
	}

	for _, attendanceData := range attendanceDataList {
		recordDate, err := time.Parse(time.DateOnly, attendanceData.DateStr)
		if err != nil {
			return AttendanceData{}, err
		}

		if recordDate.Day() == currentTime.Day() && recordDate.Month() == currentTime.Month() {
			return attendanceData, nil
		}
	}

	return AttendanceData{}, errors.New("[GetTodayAttendance] No record found.")

}

func (c *HrmsClient) GetRecentAttendance() ([]AttendanceData, error) {
	currentTime := time.Now()
	day := currentTime.Day()
	var accountingMonths []time.Month

	accountingMonths = append(accountingMonths, currentTime.Month())

	if day > 15 {
		accountingMonths = append(accountingMonths, currentTime.Month()+1)
	}

	c.logger.Debug("[GetRecentAttendance]", "accountingMonths", accountingMonths)

	type fetchResult struct {
		record []AttendanceData
		err    error
	}

	resultCh := make(chan fetchResult)
	attendanceDataList := make([]AttendanceData, 0)

	for _, month := range accountingMonths {
		go func() {
			record, err := c.FetchAttendance(strconv.Itoa(currentTime.Year()), strconv.Itoa(int(month)))
			resultCh <- fetchResult{record, err}
		}()
	}

	for range accountingMonths {
		result := <-resultCh
		if result.err != nil {
			return nil, result.err
		}
		attendanceDataList = append(attendanceDataList, result.record...)
	}

	filteredList := make([]AttendanceData, 0, len(attendanceDataList))
	for _, data := range attendanceDataList {
		recordDate, err := time.Parse(time.DateOnly, data.DateStr)
		if err != nil {
			return nil, err
		}
		if recordDate.After(time.Now()) {
			continue
		}
		if recordDate.Weekday() == 0 || recordDate.Weekday() == 6 {
			continue
		}

		filteredList = append(filteredList, data)
	}
	slices.SortFunc(filteredList, func(a AttendanceData, b AttendanceData) int {
		if a.Date.Before(b.Date) {
			return -1
		} else {
			return 1
		}
	})

	return filteredList, nil
}

type CreateLeaveApplicationRes struct {
	Status string
}

func (c *HrmsClient) CreateLeaveApplication(startTime time.Time, endTime time.Time) error {
	endpointUrl := fmt.Sprintf("%s/api/Leave/CreateLeave", c.host)

	formData := url.Values{}
	formData.Set("fldAttInOut", "VL")
	formData.Set("fldFromDate", startTime.Format(time.DateOnly))
	formData.Set("fldDateFromHour", startTime.Format("15"))
	formData.Set("fldDateFromMin", startTime.Format("04"))
	formData.Set("fldToDate", endTime.Format(time.DateOnly))
	formData.Set("fldDateToHour", endTime.Format("15"))
	formData.Set("fldDateToMin", endTime.Format("04"))
	formData.Set("ReferenceDate", "")

	c.logger.Debug("[CreateLeaveApplication] Posting CreateLeave...", "formData", formData)

	res, err := c.httpClient.PostForm(endpointUrl, formData)
	if err != nil || res.StatusCode != 200 {
		return err
	}

	resStr, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var createResponse CreateLeaveApplicationRes
	err = json.Unmarshal(resStr, &createResponse)
	if err != nil {
		return err
	}

	if createResponse.Status == "fail" {
		return fmt.Errorf("create leave fails, res: %+v", createResponse)
	}
	defer res.Body.Close()

	return nil
}
