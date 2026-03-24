## Tic-Tac-Toe

Build a production-ready, multiplayer Tic-Tac-Toe game with server-authoritative architecture using Nakama as the backend infrastructure.

## Technical requirements:

### Frontend:

- Choose your preferred tech stack (React, React Native, Flutter, Unity)
  - choosen svelte/sveltekit
- Implement responsive UI optimized for mobile devices
- Display real-time game state updates
- Show player information and match status

### Backend:

#### Core requirements

1. Server-Authoritative Game Logic

- Implement all game state management on the server side
- Validate all player moves server-side before applying them
- Prevent client-side manipulation or cheating
- Broadcast validated game state updates to connected clients

2. Matchmaking System

- Enable players to create new game rooms
- Implement automatic matchmaking to pair players
- Support game room discovery and joining
- Handle player connections and disconnections gracefully

3. Deployment

- Deploy Nakama server to a cloud provider (AWS, GCP, Azure, DigitalOcean,
  etc.)
- Deploy the frontend application with public accessibility or share mobile
  application with us
- Provide deployment documentation

#### Optional requirements

1. Concurrent Game Support

- Handle multiple simultaneous game sessions
- Implement proper game room isolation
- Ensure scalable architecture for multiple concurrent players

2. Leaderboard System

- Track player wins, losses, and win streaks
- Implement global ranking system
- Display top players with statistics
- Persist player performance data

3. Timer-Based Game Mode

- Add time limits for each player's turn (e.g., 30 seconds per move)
- Implement automatic forfeit on timeout
- Extend matchmaking to support mode selection (classic vs. timed)
- Display countdown timers in the UI
