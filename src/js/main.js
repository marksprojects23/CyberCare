// main.js
// 1 main and 4 functions
import { dashboard as dashboardFunc } from "./dashboard.js";
import { learn as learnFunc } from "./learn.js";
import { help as helpFunc } from "./help.js";
import { settings as settingsFunc } from "./settings.js";

function main() {
  dashboardFunc;
  learnFunc;
  helpFunc;
  settingsFunc;

  // Create four button elements
  var dashboardButton = document.createElement("button");
  var learnButton = document.createElement("button");
  var helpButton = document.createElement("button");
  var settingsButton = document.createElement("button");

  // Set the button's text content
  dashboardButton.textContent = "Dashboard";
  learnButton.textContent = "Learn";
  helpButton.textContent = "Help";
  settingsButton.textContent = "Settings";

  // Add a click event listener to the button
  dashboardButton.addEventListener("click", dashboardFunc);
  learnButton.addEventListener("click", learnFunc);
  helpButton.addEventListener("click", helpFunc);
  settingsButton.addEventListener("click", settingsFunc);

  // Append the button to the main content div
  //var mainContent = document.getElementById("main-content");
  document.getElementById("main-content").appendChild(dashboardButton);
  document.getElementById("main-content").appendChild(learnButton);
  document.getElementById("main-content").appendChild(helpButton);
  document.getElementById("main-content").appendChild(settingsButton);
}

main();
