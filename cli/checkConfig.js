import inquirer from "inquirer";
import * as dotenv from "dotenv";
import fs from "fs";
import path from "path";
import os from "os";

const __dirname = path.resolve();

export async function checkConfig() {
  if (
    !process.env.HRMS_HOST || // todo no account save mode
    !process.env.HRMS_USER ||
    !process.env.HRMS_PWD
  ) {
    console.log("❌ Config not completed! Please provide:");
    const configQuestions = [
      {
        type: "input",
        name: "HRMS_HOST",
        message: "HRMS Host?",
        default() {
          return process.env.HRMS_HOST;
        },
      },
      {
        type: "input",
        name: "HRMS_USER",
        message: "HRMS User?",
        default() {
          return process.env.HRMS_USER;
        },
      },
      {
        type: "password",
        name: "HRMS_PWD",
        message: "HRMS Password?",
      },
    ];

    const answers = await inquirer.prompt(configQuestions);
    console.log(answers);
    fs.writeFileSync(
      path.resolve(__dirname, ".env"),
      Object.entries(answers)
        .map(([key, value]) => `${key}=${value}`)
        .join(os.EOL)
    );
    dotenv.config();

    // console.log(JSON.stringify(answer, null, "  "));
  }
}
