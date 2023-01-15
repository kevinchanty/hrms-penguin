import { Command } from "commander";
import { checkConfig } from "./cli/checkConfig.js";
import inquirer from "inquirer";
import { createSpinner } from "nanospinner";
import * as dotenv from "dotenv";
import amendAttendanceRecordFeature from "./cli/amendAttendanceRecordFeature.js";
import { HrmsCore } from "./hrmsCore.js";
import { MockHrmsCore } from "./mockHrmsCore.js";
dotenv.config();
console.clear();

const program = new Command();

program
  .name(process.env.npm_package_name)
  .description("CLI helper to make lives easier in a company")
  .version(process.env.npm_package_version)
  .option("-m, --mockCore", "use mock HRMS core, for developing")
  .action(mainAction);
await program.parseAsync();

async function mainAction(options) {
  console.log("🐧 HRMS Penguin");

  await checkConfig();

  const hrmsCore = options.mockCore
    ? new MockHrmsCore()
    : new HrmsCore(
        process.env.HRMS_HOST,
        process.env.HRMS_USER,
        process.env.HRMS_PWD
      );
  const loginSpinner = createSpinner("Logging in to HRMS...").start();
  try {
    await hrmsCore.login();
    loginSpinner.success({ text: "Logged in!" });
  } catch (err) {
    loginSpinner.error({ text: err.message });
    process.exit(1);
  }

  const { feature } = await inquirer.prompt({
    type: "list",
    name: "feature",
    message: "Please select function:",
    choices: [
      { name: "Amend Attendance Record", value: "amendAttendanceRecord" },
      { name: "Change Config", value: "changeConfig" },
      { name: "Exit", value: "exit" },
    ],
  });

  switch (feature) {
    case "amendAttendanceRecord":
      await amendAttendanceRecordFeature(hrmsCore);
      break;
    case "changeConfig":
      console.log("todo!");
      process.exit(0);
    case "exit":
      console.log("🐧 byebye");
      process.exit(0);
    default:
      break;
  }
}
