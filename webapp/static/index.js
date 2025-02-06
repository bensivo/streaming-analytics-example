// Create WebSocket connection
const socket = new WebSocket('ws://' + window.location.host + '/ws');

socket.addEventListener('open', (event) => {
    console.log('Connected to WebSocket server');

    document.getElementById("status").textContent = "Status: Connected";
});

socket.addEventListener('message', (event) => {
    console.log('Received:', event.data);
});

socket.addEventListener('close', (event) => {
    console.log('Disconnected from WebSocket server');
    document.getElementById("status").textContent = "Status: Disconnected";
});

function sendRating(rating) {
    socket.send(JSON.stringify({
        version: 1,
        timestamp: new Date().toISOString(),
        rating: rating
    }));
}

// Button click handlers
document.getElementById('btn-very-angry').addEventListener('click', () => {
    sendRating(0);
});

document.getElementById('btn-angry').addEventListener('click', () => {
    sendRating(1);
});

document.getElementById('btn-neutral').addEventListener('click', () => {
    sendRating(2);
});

document.getElementById('btn-happy').addEventListener('click', () => {
    sendRating(3);
});

document.getElementById('btn-very-happy').addEventListener('click', () => {
    sendRating(4);
});