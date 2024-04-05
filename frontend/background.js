let activeTabId = null;
let startTime = new Map(); // Using a Map to keep track of start times for each tab
let domainTime = {};

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
