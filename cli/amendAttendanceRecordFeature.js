import chalk from "chalk";
import inquirer from "inquirer";
import { createSpinner } from "nanospinner";
import os from "os";

export async function amendAttendanceRecordFeature(hrmsCore) {
  // await hrmsCore.getAction();

  if (hrmsCore.actionItems.length === 0) {
    const { shouldFetch } = await inquirer.prompt({
      type: "confirm",
      name: "shouldFetch",
      message: "Fetch action items from HRMS?",
    });

    if (shouldFetch) {
      const fetchSpinner = createSpinner(
        "Fetching action items HRMS..."
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

    const { selectedDates } = await inquirer.prompt({
      type: "checkbox",
      name: "selectedDates",
      message: "Please selected date ",
    });
  }
}

export default amendAttendanceRecordFeature;
