import chalk from "chalk";
import inquirer from "inquirer";
import { createSpinner } from "nanospinner";
import os from "os";
import { DATE_REGEX } from "../const.js";

export async function amendAttendanceRecordFeature(hrmsCore) {
  // await hrmsCore.getAction();

  // Fetch action items from HRMS
  if (hrmsCore.actionItems.length === 0) {
    const { shouldFetch } = await inquirer.prompt({
      type: "confirm",
      name: "shouldFetch",
      message: "Fetch action items from HRMS?",
    });

    if (shouldFetch) {
      const fetchSpinner = createSpinner(
        "Fetching action items from HRMS..."
      ).start();

      try {
        await hrmsCore.fetchActions();
        fetchSpinner.success({ text: "Fetch action items from HRMS :" });
        console.log(chalk.green(hrmsCore.actionItems.flat().join(os.EOL)));
      } catch (err) {
        fetchSpinner.error(err.message);
        process.exit(1);
      }
    }
  }

  let shouldManuallyInput = true;
  let amendDates = [];

  // Select date from actions items
  if (hrmsCore.actionItemDates.length > 0) {
    const { selectedDates } = await inquirer.prompt({
      type: "checkbox",
      loop: false,
      pageSize: 20,
      name: "selectedDates",
      message: "Please selected date :",
      choices: [...hrmsCore.actionItemDates, "manually input"],
    });

    if (selectedDates.at(-1) === "manually input") {
      selectedDates.splice(amendDates.length - 1, 1);
    } else {
      shouldManuallyInput = false;
    }

    amendDates.push(...selectedDates);
  }

  // Ask for manually input date
  if (shouldManuallyInput) {
    const { inputDates } = await inquirer.prompt({
      type: "input",
      name: "inputDates",
      message:
        "Please input date , separated by comma (YYYY-MM-DD,YYYY-MM-DD):",
      validate: (answer) => {
        const answerArr = answer.split(",");

        if (answerArr.every((str) => DATE_REGEX.test(str))) {
          return true;
        } else {
          return "Format incorrect";
        }
      },
    });
    amendDates.push(...inputDates.split(","));
  }

  // POST to HRMS
  console.log("amendDates logs:", amendDates);
  const postSpinner = createSpinner("Requesting HRMS...").start();
  let amendResult = [];

  try {
    const promiseArr = Promise.allSettled(
      amendDates.map((date) =>
        hrmsCore.amendAttendanceRecord(date, "09", "00", "18", "00", "WFH")
      )
    );
    amendResult = await promiseArr;
  } catch (err) {
    postSpinner.error({ text: err });
    process.exit(1);
  }

  const outCome = amendResult.reduce(
    (prev, result, index) => {
      if (result.status === "fulfilled") {
        prev.isAllFail = false;
        prev.output.push(`✔  ${amendDates[index]}`);
      } else {
        prev.isAllSuccess = false;
        prev.output.push(`❌  ${amendDates[index]}`);
      }
      return prev;
    },
    {
      isAllSuccess: true,
      isAllFail: true,
      output: [],
    }
  );
  if (outCome.isAllSuccess) {
    postSpinner.success({ text: "Requested HRMS - All success:" });
  } else if (outCome.isAllFail) {
    postSpinner.error({ text: "Requested HRMS - All failed!:" });
  } else {
    postSpinner.warn({ text: "Requested HRMS - Some failed!:" });
  }

  console.log(outCome.output.join(os.EOL));
}

export default amendAttendanceRecordFeature;
