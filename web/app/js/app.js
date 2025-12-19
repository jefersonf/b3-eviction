const canvas = document.getElementById('captchaCanvas');
const ctx = canvas.getContext('2d');
const captchaInput = document.getElementById('captchaInput');
const refreshBtn = document.getElementById('refreshCaptcha');
const statusDiv = document.getElementById('status-msg');

const API_URL = 'http://localhost:3000/api';

let currentCaptcha = '';
let currentNominee = null;
let currentEvictionId = 'bbb-load-test-id'; // TEST


async function castVote(nomineeName) {
    currentNominee = nomineeName;
    console.log("Selected:", nomineeName);

    statusDiv.textContent = '';
    statusDiv.className = '';

    const allCards = document.querySelectorAll('.card');
    allCards.forEach(card => {
        card.classList.remove('selected');
        card.classList.add('shadow');
    });

    const activeCard = document.getElementById(`card-${nomineeName.toLowerCase()}`);
    if (activeCard) {
        activeCard.classList.add('selected');
        const voteConfirmationBox = document.getElementById("vote-confirmation-box");
        if (voteConfirmationBox) { 
            voteConfirmationBox.classList.remove('hidden')
        }
    }
}

// Random character generator
function getRandomChar() {
    const chars = 'ABCDEFGHJKLMNPQRSTUVWXYZabdefghjmnpqrstuvwxyz23456789';
    return chars.charAt(Math.floor(Math.random() * chars.length));
}

// Draw captcha
function drawCaptcha() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    
    // Add Noise: Random background lines
    for (let i = 0; i < 11; i++) {
        ctx.beginPath();
        ctx.moveTo(Math.random() * canvas.width, Math.random() * canvas.height);
        ctx.lineTo(Math.random() * canvas.width, Math.random() * canvas.height);
        ctx.strokeStyle = '#888';
        ctx.lineWidth = Number(1 + Math.random()*2).toFixed(0);
        ctx.stroke();
    }

    // Generate and Draw Text
    currentCaptcha = '';
    const fontSize = 30;
    ctx.font = `${fontSize}px Arial`;
    ctx.textBaseline = 'middle';

    for (let i = 0; i < 6; i++) {
        const char = getRandomChar();
        currentCaptcha += char;
        // Save state to rotate specific character
        ctx.save();
        // Random position
        const x = 30 * i + 20;
        const y = canvas.height / 2;
        // Move to position and rotate randomly
        ctx.translate(x, y);
        const angle = (Math.random() - 0.7) * 0.4;
        ctx.rotate(angle);
        ctx.fillStyle = '#444';
        ctx.fillText(char, 0, 0);
        ctx.restore();
    }
}

// Vote validation Logic
document.getElementById('voteForm').addEventListener('submit', function(e) {
    e.preventDefault();

    const userEntry = captchaInput.value;

    if (userEntry.toLowerCase() === currentCaptcha.toLowerCase()) {
        submitVote(); 
    } else {
        drawCaptcha(); // Regenerate on failure
        statusDiv.textContent = "Código incorreto. Por favor, tente novamente.";
        statusDiv.className = 'error';
        captchaInput.value = '';
    }
});

async function submitVote() {

    console.log("Voting for:", currentNominee);
    const payload = { evictionId: currentEvictionId, nominee_id: candidateName };

    try {
        const response = await fetch(`${API_URL}/vote`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(payload)
        });

        if (response.ok) {
            statusDiv.textContent = `Vote for ${candidateName} confirmed!`;
            statusDiv.className = 'success';
        } else {
            throw new Error('Server error');
        }
    } catch (error) {
        console.error(error);
        statusDiv.textContent = "Erro ao enviar voto. Tente novamente.";
        statusDiv.className = 'error';
    } finally {
        statusDiv.textContent = 'Voto realizado!';
        statusDiv.className = 'success';
    }

    // Reset after success
    drawCaptcha();
    captchaInput.value = '';
}

// Initialize everything...
refreshBtn.addEventListener('click', drawCaptcha);
canvas.addEventListener('click', drawCaptcha); // Allow clicking image to refresh
drawCaptcha();