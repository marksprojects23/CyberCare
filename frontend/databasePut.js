function setupTextareaListeners() {
  const whitelistTextarea = document.getElementById("whitelist");
  const blacklistTextarea = document.getElementById("blacklist");

  const saveSettings = () => {
    const whitelist = whitelistTextarea.value.split("\n");
    const blacklist = blacklistTextarea.value.split("\n");

    fetch("http://52.86.177.235:8080/settings", {
      method: "PUT", // Assuming you're using PUT for the update endpoint
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ whitelist, blacklist }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("Network response was not ok");
        }
        console.log("Settings updated successfully");
      })
      .catch((error) => console.error("Error updating settings:", error));
  };

  // Save settings when the user clicks away from the textarea
  whitelistTextarea.addEventListener("blur", saveSettings);
  blacklistTextarea.addEventListener("blur", saveSettings);

  // Save settings when the user presses Enter in the textarea
  const handleEnterPress = (event) => {
    if (event.key === "Enter") {
      saveSettings();
    }
  };
  whitelistTextarea.addEventListener("keydown", handleEnterPress);
  blacklistTextarea.addEventListener("keydown", handleEnterPress);
}

// Call the function to setup listeners
setupTextareaListeners();
