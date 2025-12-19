
const API_URL = 'http://localhost:3000/api'; // endpoint to get totals

const COLORS = [
    '#3B82F6',
    '#EF4444',
    '#10B981',
    '#F59E0B',
    '#8B5CF6', 
    '#EC4899' 
];

async function updateDashboard(evictionId) {
    try {
        const response = await fetch(`${API_URL}/stats/${evictionId}`);
        if (!response.ok) throw new Error('Failed to fetch stats');

        // {"total_votes":11,"nominee_votes":[{"nominee_id":"bob","votes":7},{"nominee_id":"cob","votes":1},{"nominee_id":"ana","votes":3}],"evicted":{"nominee_id":"bob","votes":7}}
        const data = await response.json();
    
        renderDashboard(data);
    } catch (error) {
        console.error('Error:', error);
    }
}

function renderDashboard(evictionData) {
    loadData();
    // 1. Calculate Totals
    // Convert object {Name: Count} to Array [{name: "Ana", count: 10}, ...]
    // console.log(evictionData)
    const entries = evictionData.nominee_votes;
    const totalVotes = evictionData.total_votes;

    // 2. Update Headline
    document.getElementById('count-total').innerText = totalVotes;
    
    // 3. Prepare DOM Elements
    const barContainer = document.getElementById('stacked-bar');
    const legendContainer = document.getElementById('chart-legend');
    
    barContainer.innerHTML = '';
    legendContainer.innerHTML = '';

    if (totalVotes === 0) {
        barContainer.innerHTML = '<div style="width:100%; text-align:center; line-height:40px; color:#888;">No votes yet</div>';
        return;
    }

    // 4. Sort entries (optional: highest votes first looks better)
    entries.sort((a, b) => b.votes - a.votes);

    // 5. Render Bar Segments and Legend
    entries.forEach((item, index) => {
        const percentage = ((item.votes / totalVotes) * 100).toFixed(1);
        const color = COLORS[index % COLORS.length]; // Cycle colors if > 6 nominees

        // --- Render Bar Segment ---
        const segment = document.createElement('div');
        segment.className = 'bar-segment';
        segment.style.width = `${percentage}%`;
        segment.style.backgroundColor = color;
        segment.title = `${item.nominee_id}: ${item.votes} votes (${percentage}%)`; // Tooltip
        barContainer.appendChild(segment);
        console.log(barContainer);

        // --- Render Legend Item ---
        const legendItem = document.createElement('div');
        legendItem.className = 'legend-item';
        legendItem.innerHTML = `
            <span class="color-dot" style="background-color: ${color}"></span>
            <strong>${item.nominee_id}</strong>: ${item.votes} <span class="percentage">(${percentage}%)</span>
        `;
        legendContainer.appendChild(legendItem);
    });
}

async function loadData() {
    try {
        const response = await fetch(`${API_URL}/analytics/minutely`);
        const data = await response.json();
        
        processAndRender(data);
    } catch (error) {
        console.error('Error:', error);
    }
}

function processAndRender(rawData) {
    // 3. Process Data: Group by Nominee
    // Chart.js needs datasets array: [{label: 'Nominee A', data: [...]}, ...]
    const nominees = {}; 
    
    rawData.forEach(row => {
        if (!nominees[row.nominee_id]) {
            nominees[row.nominee_id] = [];
        }
        nominees[row.nominee_id].push({
            x: row.timedate, // Time on X axis
            y: row.votes // Votes on Y axis
        });
    });

    // Convert to Chart.js Datasets
    const datasets = Object.keys(nominees).map((nomId, index) => {
        return {
            label: nomId,
            data: nominees[nomId],
            borderColor: COLORS[index % COLORS.length],
            borderWidth: 3,
            backgroundColor: COLORS[index % COLORS.length],
            fill: false,
            tension: 0.0
        };
    });

    // 4. Render Chart
    const ctx = document.getElementById('voteChart').getContext('2d');
    new Chart(ctx, {
        type: 'line',
        data: { datasets: datasets },
        options: {
            responsive: true,
            scales: {
                x: {
                    type: 'time', // Requires the date-fns adapter loaded above
                    time: { unit: 'minute' },
                    title: { display: true, text: 'Time', font: { size: 20 } },
                    ticks: {
                        font: {
                            size: 18
                        }
                    }
                },
                y: {
                    beginAtZero: true,
                    title: { display: true, text: 'Votes', font: { size: 20 } },
                    ticks: { font: { size: 18 } }
                }
            },
            plugins: {
                tooltip: {
                    titleFont: { size: 16 },
                    bodyFont: { size: 14 },
                    titleMarginBottom: 10,
                    bodyAlign: 'right',
                    usePointStyle: true
                },
                legend: { 
                    position: 'bottom',
                    title: {
                        padding: { bottom: 100 }
                    },
                    labels: { font: { size: 20 }, usePointStyle: true, boxWidth: 10, boxHeight: 10 } 
                }
            }
        }
    });
}

const selectEvictions = document.getElementById('evictions');
const  evictionId = selectEvictions.options[0].value || '';

loadData();
setInterval(updateDashboard, 500, evictionId);
