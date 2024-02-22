// main.js
// 1 main and 4 functions
import { dashboard } from "./dashboard.js";
import { learn } from "./learn.js";
import { help } from "./help.js";
import { settings } from "./settings.js";

function main() {
  dashboard();
  learn();
  help();
  settings();
}

main();
