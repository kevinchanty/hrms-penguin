import inquirer from "inquirer";
import fs from "fs";
import path from "path";

const configPath = path.resolve(process.cwd(), "hrms-config.json");

export async function getConfig() {
  try {
    let config = JSON.parse(
      fs.readFileSync(`${process.cwd()}/hrms-config.json`)
    );

    if (!config.hrmsHost || !config.hrmsUser || !config.staffNo) {
      console.log("❌ Config not completed! Please provide:");
      const answers = await promptConfig();
      writeConfig(answers);
      return answers;
    }

    return config;
  } catch (_) {
    console.log("❌ Config not completed! Please provide:");
    const answers = await promptConfig();
    writeConfig(answers);
    return answers;
  }
}

export async function promptConfig(config = {}) {
  const configQuestions = [
    {
      type: "input",
      name: "hrmsHost",
      message: "HRMS Host?",
      default: config.hrmsHost ?? "hrms.some-company.com.hk",
    },
    {
      type: "input",
      name: "hrmsUser",
      message: "HRMS User?",
      default: config.hrmsUser,
    },
    // {
    //   type: "password",
    //   name: "hrmsPwd",
    //   message: "HRMS Password?",
    // },
    {
      type: "input",
      name: "staffNo",
      message: "Staff Number?",
      default: config.staffNo,
    },
  ];

  const answers = await inquirer.prompt(configQuestions);

  return answers;
}

export function writeConfig(config = {}) {
  fs.writeFileSync(path.resolve(configPath), JSON.stringify(config, null, 2));
}
