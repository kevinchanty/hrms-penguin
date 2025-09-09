package hrmsclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"time"
)

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

type LeaveApplicationRes struct {
	Message string                 `json:"msg"`
	Data    []LeaveApplicationData `json:"data"`
}

type LeaveApplicationData struct {
	StartTime time.Time
	EndTime   time.Time
}

func (rd *LeaveApplicationData) UnmarshalJSON(data []byte) error {
	var raw struct {
		StartDate string `json:"fldAttDate"`
		EndDate   string `json:"fldAttDateTo"`
		StartTime string `json:"fldAttTime"`
		EndTime   string `json:"fldAttToTime"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	startTime, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%v %v", raw.StartDate, raw.StartTime))
	if err != nil {
		return fmt.Errorf("failed to parse StartDate: %w", err)
	}
	rd.StartTime = startTime

	endTime, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%v %v", raw.EndDate, raw.EndTime))
	if err != nil {
		return fmt.Errorf("failed to parse EndDate: %w", err)
	}
	rd.EndTime = endTime

	return nil
}

func (c *HrmsClient) FetchLeaveApplicationRecord(dateFrom string, dateTo string, page int, limit int) ([]LeaveApplicationData, error) {
	endpointUrl := fmt.Sprintf("%s/api/Leave/getpages", c.host)

	params := url.Values{}
	params.Add("fldDateFrom", dateFrom)
	params.Add("fldDateTo", dateTo)
	params.Add("page", strconv.Itoa(page))
	params.Add("limit", strconv.Itoa(limit))
	params.Add("types", "1")

	fullEndPointUrl := fmt.Sprintf("%s?%s", endpointUrl, params.Encode())

	c.logger.Debug("[FetchLeaveApplicationRecord] fetching endpoint...", "url", fullEndPointUrl)
	resp, err := c.httpClient.Get(fullEndPointUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respStr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	c.logger.Debug("[FetchLeaveApplicationRecord]", "response body", string(respStr))

	var leaveApplicationRes LeaveApplicationRes
	err = json.Unmarshal(respStr, &leaveApplicationRes)

	c.logger.Debug("[FetchLeaveApplicationRecord] result", "data", leaveApplicationRes)

	return leaveApplicationRes.Data, nil
}
