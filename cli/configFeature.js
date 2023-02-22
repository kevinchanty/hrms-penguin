import inquirer from "inquirer";
import fs from "fs";
import path from "path";

const configPath = path.resolve(
  path.dirname(process.execPath),
  "hrms-config.json"
);

export async function getConfig() {
  try {
    let config = JSON.parse(fs.readFileSync(configPath));

    if (!config.hrmsHost || !config.hrmsUser || !config.empNo) {
      throw new Error();
    }

    return config;
  } catch (_) {
    console.log("❌ Config not completed! Please provide:");
    let answers = await promptConfig();
    answers.amendTemplates = [
      {
        name: "WFH",
        inHour: "09",
        inMin: "00",
        outHour: "18",
        outMin: "00",
        remarks: "WFH",
      },
    ];
    writeConfig({
      ...answers,
      amendTemplates: [
        {
          name: "WFH",
          inHour: "09",
          inMin: "00",
          outHour: "18",
          outMin: "00",
          remarks: "WFH",
        },
      ],
    });
    return answers;
  }
}

export async function promptConfig(config = {}) {
  const configQuestions = [
    {
      type: "input",
      name: "hrmsHost",
      message: "HRMS Host?",
      default: config.hrmsHost ?? "https://hrms.some-company.com.hk",
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
      name: "empNo",
      message: "Staff Number?",
      default: config.empNo,
    },
  ];

  const answers = await inquirer.prompt(configQuestions);

  return answers;
}

export function writeConfig(config = {}) {
  console.log("Written config at:", configPath);
  fs.writeFileSync(path.resolve(configPath), JSON.stringify(config, null, 2));
}
