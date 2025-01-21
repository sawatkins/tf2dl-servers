const POLL_INTERVAL = 35000; // 35 seconds

// Function to fetch server IPs
async function fetchServerIPs() {
    try {
        const response = await fetch("/api/server-ips");
        const ips = await response.json();
        // console.log("ips", ips);
        return ips;
    } catch (error) {
        console.error('Error fetching server IPs:', error);
        return [];
    }
}

// Function to fetch server info
async function fetchServerInfo(ip) {
    try {
        const response = await fetch(`/api/server-info?ip=${ip}`);
        const serverInfo = await response.json();
        // console.log("serverInfo", serverInfo);
        return serverInfo;
    } catch (error) {
        console.error(`Error fetching server info for ${ip}:`, error);
        return null;
    }
}

// Function to update table row
function updateTableRow(serverInfo) {
    const tableBody = document.querySelector('#server-table table tbody');
    let row = tableBody.querySelector(`tr[data-ip="${serverInfo.public_ip}"]`);
    
    if (!row) {
        row = document.createElement('tr');
        row.setAttribute('data-ip', serverInfo.public_ip);
        tableBody.appendChild(row);
    }

    // console.log("serverInfo", serverInfo);

    // temportary fix
    if (serverInfo.public_ip === "54.193.198.90") {
        region = "us-west";
    }
    else {
        region = "eu-central";
    }

    row.innerHTML = `
        <td class="region ${region}">
            <img src="/img/${region}.svg" alt="${region} flag" class="flag-icon">
            <span>${region}</span>
        </td>
        <td class="online-status">Online</td>
        <td>${serverInfo.map}</td>
        <td>${serverInfo.players}/${serverInfo.max_players}</td>
        <td>
            <a href="steam://connect/${serverInfo.public_ip}:27015">Connect</a>
            <span class="copy-icon" title="Copy ${serverInfo.public_ip}" onclick="copyToClipboard(this, '${serverInfo.public_ip}')">
                <img src="/img/copy.svg" alt="Copy ${serverInfo.public_ip}" style="width: 16px; height: 16px; vertical-align: middle;">
            </span>
        </td>
    `;
}

// Function to copy IP to clipboard
function copyToClipboard(element, ip) {
    navigator.clipboard.writeText(ip).catch(err => {
        console.error('Failed to copy IP: ', err);
    });
}

// Function to display refreshing state
function setRefreshingState() {
    const refreshButton = document.getElementById('refresh');
    refreshButton.innerHTML = 'Refreshing...';
    refreshButton.style.cursor = 'default';
    refreshButton.style.textDecoration = 'none';
}

// Function to reset to ready state
function setReadyState() {
    const refreshButton = document.getElementById('refresh');
    refreshButton.innerHTML = 'Manual refresh';
    refreshButton.style.cursor = 'pointer';
    refreshButton.style.textDecoration = 'underline';
}

// Modified pollServers function
async function pollServers() {
    console.log("polling servers...");
    setRefreshingState();

    const ips = await Promise.race([
        fetchServerIPs(),
        new Promise(resolve => setTimeout(() => resolve([]), 5000)) // 5 seconds timeout
    ]);

    const headerRow = document.querySelector('#header-row');
    const defaultRow = document.querySelector('#default-row');
    const table = document.querySelector('#server-table table');

    if (ips.length > 0) {
        defaultRow.style.display = 'none';

        ips.forEach(ip => {
            fetchServerInfo(ip).then(serverInfo => {
                updateTableRow(serverInfo);
            });
        });
    } else {
        defaultRow.style.display = 'table-row';
        // Remove all rows except the header and default row
        Array.from(table.rows).forEach(row => {
            if (row !== headerRow && row !== defaultRow) {
                row.remove();
            }
        });
    }

    setReadyState();
}

// Start polling
setInterval(pollServers, POLL_INTERVAL);
pollServers();

// Add event listener for manual refresh
document.addEventListener('DOMContentLoaded', function() {
    const refreshButton = document.getElementById('refresh');
    refreshButton.addEventListener('click', pollServers);
});

