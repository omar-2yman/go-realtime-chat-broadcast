# Real-Time Concurrent Chat System in Go

This project is an upgraded version of a basic RPC chat system.  
Instead of returning chat history on request, the system now uses **real-time broadcasting** with **Go concurrency**, **goroutines**, **channels**, and **Mutex synchronization**.

---

## 🚀 Features

### 🧵 Full Concurrency
- Each client handled in its own goroutine.
- Server broadcasts in parallel using multiple goroutines.

### 📡 Instant Message Broadcast
- Messages are instantly sent to **all other users**.
- Sender does **not** receive their own message (no self-echo).

### 👥 Join Notifications
- When a client joins, all clients receive:  
  **`User [ID] joined`**

### 🔒 Safe Shared State
- Server stores all clients in a shared map.
- `sync.Mutex` ensures safe concurrent access.

### 🔄 Real-Time RPC Callbacks
- Each client opens a small RPC listener.
- Server calls `ClientRPC.Receive` to deliver messages instantly.

---

## 🧱 Project Structure
REALTIME-RPC-CHAT/
│
├── client.go # Real-time RPC client
├── server.go # RPC server with join + broadcast logic
└── README.md # Project documentation

---

## ⚙️ How It Works

### 🔧 Server Logic
- Runs on port `9000`
- Assigns unique user IDs
- Stores active clients in a Mutex-protected map
- Broadcasts:
  - Join notifications  
  - Chat messages  
- Each broadcast is done in a separate goroutine

### 🔧 Client Logic
- Connects to the RPC server
- Gets an auto-assigned user ID
- Starts a callback RPC listener
- Sends messages → server broadcasts to others  
- Receives all incoming messages instantly

---

## ▶️ How to Run

### 1️⃣ Start the Server

Open a terminal inside the project folder:

```bash
go run server.go
```
Server will start on port `:9000` and print incoming messages.

### 2️⃣ Run the Client

open another terminal (many as you want):

```bash
go run client.go
```

### 🖼 Example Output
**Server Terminal:**
```bash
Server running on port 9000...
User 1 joined
User 2 joined
[1]: hello
```

**Client Terminal:**
```bash
You joined as User 2
[User 1] hello
Message:
```
