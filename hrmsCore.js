import fetch, { FormData, Headers } from "node-fetch";

function sleep(ms) {
  return new Promise((resolve) => {
    setTimeout(resolve, ms);
  });
}

const dateRegex = /^\d{4}\-(0[1-9]|1[012])\-(0[1-9]|[12][0-9]|3[01])$/;

export class HrmsCore {
  cookie = "";
  headers = undefined;
  actionItems = [];
  actionItemDates = [];

  async login() {
    const formData = new FormData();
    formData.set("action", "login");
    formData.set("fldEmpLoginID", process.env.HRMS_USER);
    formData.set("fldEmpPwd", process.env.HRMS_PWD);
    formData.set("code", undefined);

    const response = await fetch(`${HRMS_HOST}/api/admin/login`, {
      method: "POST",
      body: formData,
    });

    if (response.status === 200) {
      this.cookie = response.headers.raw()["set-cookie"][0];
      this.headers = new Headers({ Cookie: resCookie });
    } else {
      throw new Error("HRMS-CORE: Login fail!");
    }
  }

  async fetchActions() {
    const formData = new FormData();
    formData.set("action", "maincontent");

    const response = await fetch(`${HRMS_HOST}/api/Home/GetAction`, {
      method: "POST",
      body: formData,
      headers: this.headers,
    });

    if (response.status === 200) {
      const resStr = await response.text();
      this.handleActionRes(resStr);
    } else {
      throw new Error("HRMS-CORE: Fetch Action Items fail!");
    }
  }

  handleActionRes(resStr) {
    this.actionItems = resStr
      .replaceAll("<p>", "")
      .split("</p>")
      .map((sectionStr) => {
        return sectionStr.split("<br />").map((str) => {
          return str.trim();
        });
      });
    this.actionItems.pop();
    this.actionItemDates = this.actionItems
      .flat()
      .filter((str) => dateRegex.test(str));
  }

  async getAttendanceRecord() {}

  async getAttendanceAmendRecord() {}

  async amendAttendanceRecord() {
    const data = {};
    const formData = new FormData();
    const response = fetch(
      `${HRMS_HOST}/api/Home/GetAction/api/Attendance/CreateMissAttendance`,
      {
        method: "POST",
        body: formData,
      }
    );
  }
}

// TO do
// Batch Submit
// CLI Interface

/* 
fldAttID: 
0
fldSubmitSuc: 
fldEmpNo: 
SF220083
AttDate: 
2022-12-19
fldLocation1: 
fldStartWorkHour: 
09
fldStartWorkMin: 
00
OutDate1: 
2022-12-19
fldLunchOutHour: 
18
fldLunchOutMin: 
00
fldLocation2: 
fldLunchInHour: 
fldLunchInMin: 
fldDinnerOutHour: 
fldDinnerOutMin: 
fldLocation3: 
fldDinnerInHour: 
fldDinnerInMin: 
fldFinishWorkTimeHour: 
fldFinishWorkTimeMin: 
fldAttRemark: 
WFH 
*/
