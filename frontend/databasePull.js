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

    // Call the function to fetch and populate when the page loads
    document.addEventListener('DOMContentLoaded', fetchAndPopulateLists);