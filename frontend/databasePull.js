function fetchAndPopulateLists() {
    // Replace with the actual URL of your backend endpoint
    fetch('http://52.86.177.235:8080/settings')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            // Assuming 'data' is an array and you're interested in the first item
            var WhiteList =  document.getElementById('whitelist')
            var BlackList = document.getElementById('blacklist')
            data.map(settings=> {
                const whitelist = settings.whitelist.join('\n');
                const blacklist = settings.blacklist.join('\n');
    
                // Populate the textarea elements
                WhiteList.value += whitelist + "\n";
                BlackList.value += blacklist + "\n";
            }
            );
            // Assuming the backend returns the lists directly as arrays of strings

        })
        .catch(error => console.error('There has been a problem with your fetch operation:', error));
    }

function setupTextareaListeners() {
    const whitelistTextarea = document.getElementById('whitelist');
    const blacklistTextarea = document.getElementById('blacklist');

    const saveSettings = () => {
        const whitelist = whitelistTextarea.value.split('\n');
        const blacklist = blacklistTextarea.value.split('\n');

        fetch('http://52.86.177.235:8080/settings', {
            method: 'PUT', // Assuming you're using PUT for the update endpoint
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ whitelist, blacklist })
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            console.log('Settings updated successfully');
        })
        .catch(error => console.error('Error updating settings:', error));
    };

    // Save settings when the user clicks away from the textarea
    whitelistTextarea.addEventListener('blur', saveSettings);
    blacklistTextarea.addEventListener('blur', saveSettings);

    // Save settings when the user presses Enter in the textarea
    const handleEnterPress = (event) => {
        if (event.key === 'Enter') {
            saveSettings();
            // Prevent default to avoid creating a new line after saving
            event.preventDefault(); 
        }
    };
    whitelistTextarea.addEventListener('keydown', handleEnterPress);
    blacklistTextarea.addEventListener('keydown', handleEnterPress);
}
    
    // Call the function to setup listeners
    setupTextareaListeners();
    

    // Call the function to fetch and populate when the page loads
    document.addEventListener('DOMContentLoaded', fetchAndPopulateLists);