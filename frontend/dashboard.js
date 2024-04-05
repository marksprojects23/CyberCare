document.addEventListener("DOMContentLoaded", () => {
  chrome.runtime.sendMessage({ query: "domainTime" }, (response) => {
    // Assume response is an object where keys are domain names and values are time in milliseconds
    const domainNames = Object.keys(response);
    const timeSpent = domainNames.map((domain) =>
      Math.round(response[domain] / 60000)
    ); // Convert to minutes

    // Mock data for the total time by hour (this should come from your actual data)
    const totalTimeByHour = [30, 45, 60, 15, 90, 30, 60, 90, 120]; // Example data

    // Calculate the current total time
    const totalMinutes = totalTimeByHour.reduce((acc, val) => acc + val, 0);
    const hours = Math.floor(totalMinutes / 60);
    const minutes = totalMinutes % 60;

    const ctx = document.getElementById("timeChart").getContext("2d");
    const timeChart = new Chart(ctx, {
      type: "bar",
      data: {
        labels: domainNames,
        datasets: [
          {
            label: "Minutes on Domain",
            data: timeSpent,
            backgroundColor: domainNames.map((domain) => getBrandColor(domain)),
            // You would define the getBrandColor function to match domains with brand colors
            borderColor: domainNames.map((domain) => getBrandColor(domain)),
            borderWidth: 1,
          },
          {
            label: "Total Minutes by Hour",
            data: totalTimeByHour,
            // ... additional dataset configuration for line chart part
          },
        ],
      },
      options: {
        plugins: {
          legend: {
            display: false, // set to true if you want to display legend
          },
          title: {
            display: true,
            text: `Total: ${hours}h ${minutes}m`,
          },
        },
        // Additional chart options go here
      },
    });
  });
});
function getBrandColor(domain) {
  const brandColors = {
    "facebook.com": "#3b5998",
    "youtube.com": "#FF0000",
    "twitter.com": "#1DA1F2",
    // Add more mappings as necessary
  };
  return brandColors[domain] || "#999"; // Fallback color
}
