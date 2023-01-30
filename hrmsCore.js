import fetch, { FormData, Headers } from "node-fetch";
import { DATE_REGEX } from "./const.js";

export class HrmsCore {
  #hrmsHost = "";
  #hrmsUser = "";
  #hrmsPwd = "";
  #empNo = "";
  #cookie = "";
  #headers = undefined;
  actionItems = [];
  actionItemDates = [];

  constructor(hrmsHost, hrmsUser, hrmsPwd) {
    this.#hrmsHost =
      hrmsHost.at(-1) === "/"
        ? hrmsHost.slice(0, hrmsHost.length - 1)
        : hrmsHost;
    this.#hrmsUser = hrmsUser;
    this.#hrmsPwd = hrmsPwd;
  }

  async login() {
    const formData = new FormData();
    formData.set("action", "login");
    formData.set("fldEmpLoginID", this.#hrmsUser);
    formData.set("fldEmpPwd", this.#hrmsPwd);
    formData.set("code", undefined);

    const response = await fetch(`${this.#hrmsHost}/api/admin/login`, {
      method: "POST",
      body: formData,
    });

    if (response.status === 200 && response.headers.raw()["set-cookie"]) {
      this.#cookie = response.headers.raw()["set-cookie"][0];
      this.#headers = new Headers({ Cookie: this.#cookie });
    } else {
      throw new Error("HRMS-CORE: Login fail!");
    }
  }

  async fetchActions() {
    const formData = new FormData();
    formData.set("action", "maincontent");

    const response = await fetch(`${this.#hrmsHost}/api/Home/GetAction`, {
      method: "POST",
      body: formData,
      headers: this.#headers,
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
      .filter((str) => DATE_REGEX.test(str));
  }

  async getAttendanceRecord() {}

  async getAttendanceAmendRecord() {}

  amendAttendanceRecord(date, inHour, inMin, outHour, outMin, remarks) {
    return new Promise(async (resolve, reject) => {
      const data = {};
      const formData = new FormData();
      formData.set("fldAttID", 0);
      formData.set("AttDate", date);
      formData.set("OutDate1", date);
      formData.set("fldEmpNo", this.#empNo);
      formData.set("fldStartWorkHour", inHour);
      formData.set("fldStartWorkMin", inMin);
      formData.set("fldLunchOutHour", outHour);
      formData.set("fldLunchOutMin", outMin);
      formData.set("fldAttRemark", remarks);

      const response = await fetch(
        `${this.#hrmsHost}/api/Attendance/CreateMissAttendance`,
        {
          method: "POST",
          body: formData,
          headers: this.#headers,
        }
      );

      if (response.status === 200) {
        resolve();
      } else {
        console.log(response);

        reject(`HRMS-CORE: Amend Record Fail ! (${date})`);
      }
    });
  }
}
