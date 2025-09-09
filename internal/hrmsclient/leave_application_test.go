package hrmsclient

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestParseLeaveApplicationRecord(t *testing.T) {
	jsonString := `{
      "AttachPath": "",
      "fldEmpAlias": "",
      "fldEmpShop": "",
      "AttachPathURLForAPI": "",
      "fldEmpPerferredName": "",
      "bq": "",
      "YearBq": "",
      "fldEmploymentNature": "",
      "fldAttachmentWebPath": "",
      "fldAttRemark": "",
      "fldAttInOut": "VL",
      "fldLastMDate": "2025-08-20",
      "fldAttDate": "2025-08-19",
      "fldAttDateTo": "2025-08-19",
      "fldAttachment": "",
      "fldAttTime": "09:29",
      "fldAttToTime": "13:52",
      "fldHours": 4.0,
      "LeaveDays": 0.0,
      "fldDays": "4H23M",
      "fldMinutes": 23.0,
      "fldApprovalStatus": "Approved",
      "fldCancleStatus": "Application",
      "ApprovalNum": "1",
      "LeaveName": "Vacation Leave",
      "fldEmpName1": "",
      "fldEmpChiName1": "",
      "fldEmpRefNo1": "",
      "fldInCompany": "",
      "fldInBranch": "",
      "fldStatus": "Approved",
      "fldAttInOut1": "VL",
      "fldAttDate1": null,
      "fldAttInOut2": "AM",
      "fldAttTime1": "09:29",
      "fldAttToDate": null,
      "fldAttToInOut2": "PM",
      "fldAttToTime1": "13:52",
      "fldDays1": 0.0,
      "fldHours1": 4.0,
      "fldMinutes1": 23.0,
      "fldOT": 0.0,
      "fldLate": 0.0,
      "fldEarly": 0.0,
      "fldAttRemark1": "",
      "fldRejectBy": "",
      "fldRejectDate": "",
      "fldLeaveApvCatID2": "",
      "fldLeaveApvCatID3": "",
      "fldLeaveApvDate": "2025-08-21 12:36:36",
      "fldLeaveApvEmpName2": "",
      "fldLeaveApvDate2": "",
      "fldLeaveApvEmpName3": "",
      "fldLeaveApvDate3": "",
      "fldComID": "",
      "fldLastMDate1": null,
      "fldPostDate": null,
      "fldRefDate": null,
      "fldCancleStatus1": "",
      "fldEmpWorkHrsPerD": 0.0,
      "fldLeaveDays": 0.48703703703703705,
      "fldAttachment1": "",
      "fldAttachmentWebPath1": "",
      "fldAttTimeTo": "",
      "fldApprovalDate": null,
      "fldLeaveApvStatus": "Approved",
      "fldLeaveApvStatus2": "",
      "fldLeaveApvStatus3": "",
      "fldLeaveApvCatIDB": "",
      "fldLeaveApvCatIDC": "",
      "fldLeaveApvCatIDD": "",
      "fldLeaveApvCatID2B": "",
      "fldLeaveApvCatID2C": "",
      "fldLeaveApvCatID2D": "",
      "fldLeaveApvCatID3B": "",
      "fldLeaveApvCatID3C": "",
      "fldLeaveApvCatID3D": "",
      "fldCancleBy": "",
      "fldCancleDate": null,
      "fldApvEmpName": "",
      "fldApvDate": null,
      "fldApvStatus": "",
      "fldApvEmpName2": "",
      "fldApvDate2": null,
      "fldApvStatus2": "",
      "fldApvEmpName3": "",
      "fldApvDate3": null,
      "fldApvStatus3": "",
      "fldSubmitDate": null,
      "fldReferenceDate": null,
      "fldSubmittedDate": "",
      "fldFromDate": "",
      "fldResult": "",
      "fldAMPM": ""
    }`

	var got LeaveApplicationData
	err := json.Unmarshal([]byte(jsonString), &got)
	if err != nil {
		t.Errorf("Failed to unmarshal: %s", err.Error())
	}

	wantStartTime, err := time.Parse("2006-01-02 15:04", "2025-08-19 09:29")
	if err != nil {
		t.Errorf("Failed to parse wantStartTime")
	}
	wantEndTime, err := time.Parse("2006-01-02 15:04", "2025-08-19 13:52")
	if err != nil {
		t.Errorf("Failed to parse wantEndTime")
	}

	want := LeaveApplicationData{
		StartTime: wantStartTime,
		EndTime:   wantEndTime,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got = %+v, want %+v", got, want)
	}
}
