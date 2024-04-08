let activeTabId = null;
let startTime = new Map(); // Using a Map to keep track of start times for each tab
let domainTime = {};
var whitelist = [];
var blacklist = [];

// Dashboard configuration belows
chrome.tabs.onActivated.addListener((activeInfo) => {
  if (activeTabId !== null) {
    updateDomainTime(activeTabId);
  }

  activeTabId = activeInfo.tabId;
  startTime.set(activeTabId, performance.now()); // Set start time for new active tab

  // Initialize domain time for the new tab if necessary
  chrome.tabs.get(activeTabId, (tab) => {
    if (!tab.url) return; // Early return if tab.url is undefined
    let url = new URL(tab.url);
    let domain = url.hostname;
    if (!domainTime[domain]) {
      domainTime[domain] = 0;
    }
  });
});

chrome.tabs.onUpdated.addListener((tabId, changeInfo, tab) => {
  if (tabId === activeTabId && changeInfo.status === "complete") {
    startTime.set(activeTabId, performance.now());
  }
});

chrome.windows.onFocusChanged.addListener(() => {
  if (activeTabId !== null) {
    updateDomainTime(activeTabId);
    startTime.set(activeTabId, performance.now());
  }
});

function updateDomainTime(tabId) {
  chrome.tabs.get(tabId, (tab) => {
    if (chrome.runtime.lastError || !tab.url) {
      // Log error or simply return if there's an issue getting the tab or if tab.url is undefined
      console.log(chrome.runtime.lastError?.message);
      return;
    }
    let url = new URL(tab.url);
    let domain = url.hostname;
    let currentTime = performance.now();
    let start = startTime.get(tabId) || currentTime; // Use current time if start time is not set (shouldn't happen)
    let timeSpent = currentTime - start; // Calculate time spent in seconds
    domainTime[domain] += timeSpent;

    // Reset start time for the current tab
    startTime.set(tabId, currentTime);
  });
}

chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.query === "domainTime") {
    sendResponse(domainTime);
  }
});

// FILTERING STUFF BELOW THIS LINE

function fetchLists() {
  fetch("http://52.86.177.235:8080/settings")
    .then((response) => {
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      return response.json();
    })
    .then((data) => {
      // Assuming the first element of the array contains the lists
      const settings = data[0];

      // Clean the lists by filtering out empty strings
      whitelist = settings.whitelist.filter((item) => item);
      blacklist = settings.blacklist.filter((item) => item);

      // Log the lists to the console to confirm they are correct
      console.log("Whitelist:", whitelist);
      console.log("Blacklist:", blacklist);
    })
    .catch((error) => console.error("Failed to fetch settings:", error));
}

// Listen for web requests and block them based on the whitelist and blacklist
chrome.webRequest.onBeforeRequest.addListener(
  function (details) {
    const url = new URL(details.url);
    const domain = url.hostname.replace("www.", ""); // Normalize domain name

    // If there's a whitelist, block everything that's not on it
    if (whitelist.length > 0 && !whitelist.includes(domain)) {
      return { cancel: true }; // Not on whitelist, so block it
    }

    // If there's a blacklist, block everything that's on it
    if (blacklist.includes(domain)) {
      return { cancel: true }; // On blacklist, so block it
    }

    // Otherwise, allow the request
    return { cancel: false };
  },
  { urls: ["<all_urls>"] },
  ["blocking"]
);

// Initially fetch lists and then set an interval to regularly update them
fetchLists();
setInterval(fetchLists, 300000); // Update every 5 minutes

chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.action === "fetchLists") {
    fetchLists();
    sendResponse({ status: "Lists updated" });
  }
});
