# Tic-Tac-Toe

A production-ready, multiplayer Tic-Tac-Toe game with a server-authoritative architecture built using Nakama.

## Links:
- Frontend: [Vercel](https://ttt.codebysaran.in)
- Backend: [Digital Ocean](https://be.codebysaran.in)

## Setup and Installation

### Backend

The backend uses Nakama and PostgreSQL, containerized with Docker.

1. Navigate to the backend directory:
   cd backend
2. Start the services using Docker Compose:
   docker-compose up --build -d
3. The Nakama server will be available at `http://localhost:7350`.

### Frontend

The frontend is a SvelteKit application using pnpm.

1. Navigate to the frontend directory:
   cd frontend
2. Install the dependencies:
   pnpm install
3. Create a `.env` file based on the example file `.example.env`
4. Start the development server:
   pnpm run dev
5. Access the game at the localhost URL `http://localhost:5173`.

## Architecture and Design Decisions

The project consists of a server-authoritative backend and a modern SvelteKit frontend.

### Backend Architecture
- **Infrastructure**: Nakama provides the core multiplayer infrastructure, including socket connections, real-time message broadcasting, and player presence management.
- **Server-Authoritative Logic**: All game state and move validations are handled on the server. Clients submit moves, and the server validates them against the current state before broadcasting updates, preventing client-side cheating.
- **Matchmaking**: Built-in Nakama mechanisms manage game rooms, pairing players, and handling disconnections.

### Frontend Architecture
- **Framework**: SvelteKit combined with Svelte 5 for reactive, class-based state management (Session and Match stores).
- **Component Styling**: Styled using Tailwind CSS, focusing on depth through surface tiers and "frosted glass" effects for floating elements.
- **Connection drop handling**: The frontend persists the session and match state in the browser's local storage. This allows the user to close the browser and reopen it later to resume the game.

## Deployment Process

`Frontend -> HTTPS/WSS -> Nginx -> Nakama -> Postgres`

### 1. Backend deployment
- Hosted on a cloud VM (DigitalOcean droplet)
- Uses Docker Compose (Nakama + Postgres)
- Nakama runs on localhost (`127.0.0.1:7350`), not publicly exposed
- Deployment command: `docker compose up -d`

### 2. Reverse proxy and HTTPS
- Nginx used as reverse proxy
- Domain: `be.codebysaran.in`
- Proxies to `localhost:7350`
- SSL via Let's Encrypt (Certbot), supporting HTTPS and WSS

### 3. Frontend deployment
- Hosted on Vercel
- Uses environment variables:
  - `PUBLIC_NAKAMA_HOST`
  - `PUBLIC_NAKAMA_PORT`
  - `PUBLIC_NAKAMA_SERVER_KEY`
  - `PUBLIC_NAKAMA_USE_SSL`

### 4. Security/networking
- Only ports 80 and 443 exposed
- Port 7350 blocked externally
- All traffic routed through Nginx

## API/Server Configuration Details

- **Host**: `localhost`
- **Port**: `7350` (gRPC/HTTP)
- **Database**: PostgreSQL (configured via `docker-compose.yml`)
- **Nakama Server Key**: `defaultKey` (development)

## How to Test Multiplayer Functionality

1. Ensure both the Nakama backend and the SvelteKit frontend are running.
2. Open your web browser and navigate to the frontend local URL.
3. Open a second, independent session. This can be done by using an Incognito/Private window or a completely different web browser.
4. In the first window, initiate matchmaking or create a game room.
5. In the second window, join the matchmaking pool or the specific room created.
6. Observe the real-time updates as you play turns sequentially between the two browser windows.
