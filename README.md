# RingRiot

RingRiot is a real-time multiplayer game. The backend is written in Go and hosted on Google Cloud Platform (GCP), while the client is built using the Godot Engine with GDScript.
Live demo: [RingRiot on GCP](https://harshgharat.itch.io/ring-riot)

---

## Table of Contents

- [Screenshots](#screenshots)
- [Features](#features)
- [Architecture](#architecture)
- [Requirements](#requirements)
- [Installation](#installation)
  - [Backend (Server)](#backend-server)
  - [Client (Godot)](#client-godot)
- [Running the Game](#running-the-game)
- [Environment Variables](#environment-variables)
- [Deployment](#deployment)
- [Contributing](#contributing)
- [License](#license)

---

## Screenshots
loginpage: 
<h1 align="center">
  <a href="https://harshgharat.itch.io/ring-riot"><img width="250" src="https://github.com/harsh16coder/RingRiot/blob/master/LoginPage.png" /></a>
</h1>
gameplay: 
<h1 align="center">
  <a href="https://harshgharat.itch.io/ring-riot"><img width="250" src="https://github.com/harsh16coder/RingRiot/blob/master/Playground.png" /></a>
</h1>

## Features

- Real-time multiplayer gameplay
- WebSocket-based communication
- Player authentication and hiscore tracking
- Profanity filter for chat
- Dockerized backend for easy deployment
- Hosted on GCP for scalability

---

## Architecture

- **Backend:** Go, Docker, GCP (Cloud Run/VM/Container Registry)
- **Client:** Godot Engine (GDScript)
- **Communication:** WebSockets, Protocol Buffers

---

## Requirements

### Backend

- [Go 1.20+](https://golang.org/dl/)
- [Docker](https://www.docker.com/get-started)
- [Google Cloud SDK](https://cloud.google.com/sdk) (for deployment)
- GCP Project with billing enabled

### Client

- [Godot Engine 4.x](https://godotengine.org/download)
- (Optional) [Git](https://git-scm.com/) for source control

---

## Installation

### Backend (Server)

1. **Clone the repository:**
    ```sh
    git clone https://github.com/harsh16coder/ringriot.git
    cd ringriot/server
    ```

2. **Set up environment variables:**
    - Copy `.env.example` to `.env` and edit as needed (e.g., `PORT`, `DATA_PATH`).

3. **Build and run with Docker:**
    ```sh
    docker compose up --build
    ```
    This will:
    - Build the Go backend
    - Mount `banned_words.txt` and data directory
    - Expose the server on the port specified in `.env`

4. **(Optional) Run locally without Docker:**
    ```sh
    go mod download
    go run ./cmd/main.go --config .env
    ```

### Client (Godot)

1. **Clone the repository (if not already):**
    ```sh
    git clone https://github.com/harsh16coder/ringriot.git
    cd ringriot/client
    ```

2. **Open the project in Godot:**
    - Launch Godot Engine.
    - Click "Import", select `client/project.godot`.

3. **Run the game:**
    - Click "Play" (F5) in Godot.
    - The client will attempt to connect to the backend WebSocket server (update the URL in `states/entered/entered.gd` if needed).

---

## Running the Game

- **Start the backend** (see above).
- **Run the client** from Godot.
- Multiple clients can connect to the same backend for multiplayer gameplay.

---

## Environment Variables

Backend `.env` example:
```
PORT=8080
DATA_PATH=/gameserver/data
```

- `PORT`: Port to run the backend server on.
- `DATA_PATH`: Path for persistent data (mounted as a Docker volume).

---

## Deployment

### Backend (GCP Cloud Run Example)

1. **Build and push Docker image:**
    ```sh
    gcloud builds submit --tag gcr.io/YOUR_PROJECT_ID/ringriot-backend
    ```

2. **Deploy to Cloud Run:**
    ```sh
    gcloud run deploy ringriot-backend \
      --image gcr.io/YOUR_PROJECT_ID/ringriot-backend \
      --platform managed \
      --region YOUR_REGION \
      --allow-unauthenticated \
      --set-env-vars PORT=8080,DATA_PATH=/gameserver/data
    ```

3. **Mount persistent storage and banned words file as needed.**

### Client

- Distribute the Godot project or export to HTML5/Windows/Linux/Mac as needed.
- Update the WebSocket URL in the client to point to your GCP backend.

---

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/foo`)
3. Commit your changes (`git commit -am 'Add foo'`)
4. Push to the branch (`git push origin feature/foo`)
5. Create a Pull Request

---

## License

This project is licensed under the MIT License. See the [LICENSE](../LICENSE) file for details.

---

## Troubleshooting

- **WebSocket connection fails:** Check backend logs, firewall rules, and ensure the client is using the correct WebSocket URL.
- **Banned words not filtered:** Ensure `banned_words.txt` is mounted and readable by the backend container.
- **Database errors:** Check `DATA_PATH` permissions and existence.

---

## Credits

- Backend: Go, Docker, GCP
- Client: Godot Engine, GDScript
- Protocol Buffers: [godobuf](https://github.com/kittenseater/godobuf)