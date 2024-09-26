const waitingScreen = document.getElementById('waiting-screen');
const chatScreen = document.getElementById('chat-screen');
const chatMessages = document.getElementById('chat-messages');
const chatForm = document.getElementById('chat-form');
const chatInput = document.getElementById('chat-input');

let socket;

function connectWebSocket() {
    socket = new WebSocket('ws://' + window.location.host + '/ws');

    socket.onopen = function(e) {
        console.log('WebSocket connection established');
        waitingScreen.textContent = 'Waiting for a partner...';
    };

    socket.onmessage = function(event) {
        if (event.data === "Waiting for a partner...") {
            waitingScreen.textContent = event.data;
        } else if (event.data === "You've been paired with a stranger. Start chatting!") {
            waitingScreen.style.display = 'none';
            chatScreen.style.display = 'flex';
            addMessage('System', event.data);
        } else if (event.data === "Your partner has disconnected.") {
            addMessage('System', event.data);
            setTimeout(() => {
                chatScreen.style.display = 'none';
                waitingScreen.style.display = 'flex';
                waitingScreen.textContent = 'Waiting for a new partner...';
                chatMessages.innerHTML = '';
            }, 3000);
        } else {
            addMessage('Stranger', event.data);
        }
    };

    socket.onclose = function(event) {
        console.log('WebSocket connection closed');
        addMessage('System', 'Connection lost. Trying to reconnect...');
        setTimeout(connectWebSocket, 3000);
    };

    socket.onerror = function(error) {
        console.log('WebSocket error: ' + error.message);
        addMessage('System', 'An error occurred. Please try refreshing the page.');
    };
}

function addMessage(sender, message) {
    const messageElement = document.createElement('div');
    messageElement.innerHTML = `<strong>${sender}:</strong> ${message}`;
    chatMessages.appendChild(messageElement);
    chatMessages.scrollTop = chatMessages.scrollHeight;
}

connectWebSocket();

chatForm.addEventListener('submit', function(e) {
    e.preventDefault();
    if (chatInput.value) {
        socket.send(chatInput.value);
        addMessage('You', chatInput.value);
        chatInput.value = '';
    }
});
