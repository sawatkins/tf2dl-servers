// Configuration
const API_ENDPOINT = 'https://api.example.com/servers'; // Replace with your actual API endpoint
const POLL_INTERVAL = 30000; // Poll every 30 seconds

// Function to fetch server IPs
async function fetchServerIPs() {
    try {
        const response = await fetch(API_ENDPOINT);
        return await response.json();
    } catch (error) {
        console.error('Error fetching server IPs:', error);
        return [];
    }
}

// Function to fetch server info
async function fetchServerInfo(ip) {
    try {
        const response = await fetch(`http://${ip}:8000/server-info`);
        return await response.json();
    } catch (error) {
        console.error(`Error fetching server info for ${ip}:`, error);
        return null;
    }
}

// Function to update table row
function updateTableRow(serverInfo) {
    const tableBody = document.querySelector('#server-table table tbody');
    let row = tableBody.querySelector(`tr[data-ip="${serverInfo.ip}"]`);
    
    if (!row) {
        row = document.createElement('tr');
        row.setAttribute('data-ip', serverInfo.ip);
        tableBody.appendChild(row);
    }

    row.innerHTML = `
        <td class="region ${serverInfo.region.toLowerCase()}">
            <img src="${serverInfo.region.toLowerCase()}.svg" alt="${serverInfo.region} flag" class="flag-icon">
            <span>${serverInfo.region}</span>
        </td>
        <td>${serverInfo.status}</td>
        <td>${serverInfo.map}</td>
        <td>${serverInfo.players}/${serverInfo.maxPlayers}</td>
        <td><a href="steam://connect/${serverInfo.ip}">Connect</a></td>
    `;
}

// Main polling function
async function pollServers() {
    const ips = await fetchServerIPs();
    for (const ip of ips) {
        const serverInfo = await fetchServerInfo(ip);
        if (serverInfo) {
            updateTableRow(serverInfo);
        }
    }
}

// Start polling
setInterval(pollServers, POLL_INTERVAL);
pollServers(); // Initial poll

