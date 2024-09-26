# Random Chat Application

This is a simple random text chat application inspired by Omegle, where users can connect with random strangers via WebSocket for a text-only conversation. The server assigns pairs of users into chat sessions automatically once two users are available.

## Features
- Random pairing of users for text chat.
- WebSocket-based communication between users.
- Real-time messaging.
- Lightweight, with minimal dependencies.
- User notifications when paired and when disconnected.

## Project Structure
The project is organized into several packages and modules to maintain clean and modular code.

```
.
├── cmd
│   └── server
│       └── main.go       # Main entry point for the server
├── internal
│   ├── handlers
│   │   └── chat.go       # Chat handlers, including WebSocket and pairing logic
│   └── websocket
│       └── client.go     # WebSocket client and connection logic
├── pkg
│   └── queue
│       └── queue.go      # A simple thread-safe queue for managing connected users
├── static
│   └── js
│       └── chat.js       # JS for realtime chat rendering
│   └── css
│       └── style.css     # Basic CSS
└── templates
    └── index.html        # Home page for the chat application
```

## How It Works

1. **User Queue Management**: 
   - When a user connects to the WebSocket, they are pushed into a queue. The server checks the queue continuously, and once there are two users in the queue, they are paired into a chat.
   - The `Queue` implementation is thread-safe, ensuring that multiple users can be queued concurrently.

2. **WebSocket Communication**:
   - Users interact through WebSockets. Upon connection, each user is notified that they are waiting for a partner. Once paired, they can start chatting in real-time.

3. **Client Lifecycle**:
   - Each user has a `Client` object representing their WebSocket connection.
   - If a client disconnects, the partner is notified, and the session is terminated.

4. **Chat Handling**:
   - Messages are sent between paired users in real-time using the `readPump` and `writePump` functions for each WebSocket connection. These ensure that messages are read from one user and passed to the other.

## Dependencies

- **Gorilla WebSocket**: Used for WebSocket-based real-time communication.
  ```
  go get github.com/gorilla/websocket
  ```

## Installation and Usage

### Prerequisites

- Go 1.18+ installed on your machine.
- Basic knowledge of how to run Go applications.

### Steps to Run the Application

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/your-username/random-chat.git
   cd random-chat
   ```

2. **Install Dependencies**:
   Install the necessary Go packages:
   ```bash
   go mod tidy
   ```

3. **Run the Server**:
   To start the server, run the following command from the root directory:
   ```bash
   go run ./cmd/server/main.go
   ```

4. **Access the Chat Application**:
   Open your web browser and navigate to:
   ```
   http://localhost:8080
   ```

   You can now use the chat interface. The server will pair users once two users are connected.

### Folder Explanation

- **cmd/server/main.go**:
  - This is the entry point of the application. It initializes the user queue and WebSocket hub, sets up HTTP handlers, and starts the HTTP server on port `8080`.

- **internal/handlers/chat.go**:
  - Handles WebSocket connections, reads and writes messages, pairs users, and manages the lifecycle of the WebSocket sessions.

- **pkg/queue/queue.go**:
  - Implements a thread-safe queue to manage users waiting for a chat partner.

- **internal/websocket/client.go**:
  - Contains logic for handling the WebSocket connection and managing client messaging.

### Static Files and Templates

- **static/**: This folder contains any static assets (like CSS or JS files).
- **templates/index.html**: The home page template for the chat application.

### API Endpoints

- `/`: Serves the home page for the application.
- `/ws`: Establishes a WebSocket connection for real-time messaging.

### Known Limitations

- No user authentication, all users are anonymous.
- Text-only chat (no multimedia support).
- Limited error handling and reconnection logic for WebSocket disconnections.
- Basic UI and user experience.

### Future Improvements

- Add support for multimedia messages (images, videos, etc.).
- Improve user experience with a better frontend interface.
- Implement user authentication or identification.
- Add reconnection handling for WebSocket connections.

### License

This project is licensed under the MIT License.

---

Feel free to contribute to this project and improve it!
