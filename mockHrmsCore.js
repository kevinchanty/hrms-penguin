import fetch, { FormData, Headers } from "node-fetch";
import chalk from "chalk";

function sleep(ms) {
  return new Promise((resolve) => {
    setTimeout(resolve, ms);
  });
}

const dateRegex = /^\d{4}\-(0[1-9]|1[012])\-(0[1-9]|[12][0-9]|3[01])$/;

export class MockHrmsCore {
  cookie = "";
  headers = undefined;
  actionItems = [];
  sleepTime = 1;

  constructor() {
    console.log(chalk.yellow("❗USING MOCK CORE"));
  }

  async login(fail) {
    await sleep(this.sleepTime);

    this.cookie = "set!";
    this.headers = new Headers({ Cookie: "set!" });
  }

  async fetchActions(fail) {
    const formData = new FormData();
    formData.set("action", "maincontent");

    await sleep(this.sleepTime);
    if (fail) {
      throw new Error("HRMS-CORE: Fetch Action Items fail!");
    }

    const resStr =
      "<p>Missing Attendance record 欠缺出入勤紀錄:<br /> 2022-12-19<br />2022-12-22<br />2022-12-29</p><p>Early leave:<br /> 2022-12-23</p>";
    this.handleActionRes(resStr);
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
    const response = fetch("", {
      method: "POST",
      body: formData,
    });
  }
}
