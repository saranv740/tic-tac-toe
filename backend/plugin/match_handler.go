package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
)

var _ runtime.Match = &MatchHandler{}

type MatchHandler struct{}

const (
	PlayerX = 1
	PlayerO = 2
)

type GameState int

const (
	_ GameState = iota
	StateWaiting
	StatePlaying
	StateEnded
)

const (
	OpCodeMove               = 1
	OpCodeStateUpdate        = 2
	OpCodeGameStart          = 3
	OpCodeGameEnded          = 4
	OpCodePlayerDisconnected = 5
	OpCodePlayerReconnected  = 6
)

const TurnTimeLimitSec = 30

type MatchState struct {
	State       GameState                   `json:"state"`
	Board       []int                       `json:"board"`
	Turn        int                         `json:"turn"`
	PlayerX     string                      `json:"player_x"`
	PlayerO     string                      `json:"player_o"`
	Winner      int                         `json:"winner"`
	Deadline    int64                       `json:"deadline"`
	IsTimedMode bool                        `json:"is_timed_mode"`
	Presences   map[string]runtime.Presence `json:"-"` // Keyed by session ID
	Disconnects map[string]int64            `json:"-"` // Keyed by user ID
}

type MoveMessage struct {
	Position int `json:"position"`
}

type UserStats struct {
	Wins   int `json:"wins"`
	Losses int `json:"losses"`
	Draws  int `json:"draws"`
}

type EndMessage struct {
	Winner int    `json:"winner"`
	Reason string `json:"reason"`
}

// Invoked when a match is created as a result of the match create function and sets up the initial state of a match
func (m *MatchHandler) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]any) (any, int, string) {
	isTimedMode := false
	if mode, ok := params["mode"].(string); ok && mode == "timed" {
		isTimedMode = true
	}

	state := &MatchState{
		State:       StateWaiting,
		Board:       make([]int, 9),
		Presences:   make(map[string]runtime.Presence),
		Disconnects: make(map[string]int64),
		IsTimedMode: isTimedMode,
	}
	tickRate := 5
	label := "tic-tac-toe"

	return state, tickRate, label
}

// Executed when a user attempts to join the match using the client's match join operation
func (m *MatchHandler) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state any, presence runtime.Presence, metadata map[string]string) (any, bool, string) {
	s := state.(*MatchState)

	// Only allow 2 players
	if len(s.Presences) >= 2 {
		return s, false, "match full"
	}

	// Prevent the same player joining twice from different sessions simultaneously
	if _, ok := s.Presences[presence.GetUserId()]; ok {
		return s, false, "already joined"
	}

	return s, true, ""
}

// Executed when one or more users have successfully completed the match join process after their MatchJoinAttempt() returns true
func (m *MatchHandler) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state any, presences []runtime.Presence) any {
	s := state.(*MatchState)

	for _, p := range presences {
		s.Presences[p.GetSessionId()] = p

		isRejoin := false
		if s.PlayerX == p.GetUserId() || s.PlayerO == p.GetUserId() {
			isRejoin = true
			delete(s.Disconnects, p.GetUserId())
		} else if s.PlayerX == "" {
			s.PlayerX = p.GetUserId()
		} else if s.PlayerO == "" && s.PlayerX != p.GetUserId() {
			s.PlayerO = p.GetUserId()
		}

		if isRejoin && s.State == StatePlaying {
			data, err := json.Marshal(map[string]string{"user_id": p.GetUserId()})
			if err != nil {
				logger.Error("Failed to marshal player reconnected message: %v", err)
			} else {
				dispatcher.BroadcastMessage(OpCodePlayerReconnected, data, nil, nil, true)
			}
			s.broadcastState(dispatcher, OpCodeStateUpdate)
		}
	}

	// Check if match is ready to start
	if s.State == StateWaiting && s.PlayerX != "" && s.PlayerO != "" && len(s.Presences) == 2 {
		s.State = StatePlaying
		s.Turn = PlayerX
		if s.IsTimedMode {
			s.Deadline = time.Now().Add(TurnTimeLimitSec * time.Second).Unix()
		}

		// Broadcast game start
		s.broadcastState(dispatcher, OpCodeGameStart)
	}

	return s
}

// Executed when one or more users have left the match for any reason including connection loss.
func (m *MatchHandler) MatchLeave(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state any, presences []runtime.Presence) any {
	s := state.(*MatchState)

	for _, p := range presences {
		delete(s.Presences, p.GetSessionId())

		// Assign a 30s death timer instead of instant forfeit
		if s.State == StatePlaying {
			s.Disconnects[p.GetUserId()] = time.Now().Add(TurnTimeLimitSec * time.Second).Unix()
			data, err := json.Marshal(map[string]string{"user_id": p.GetUserId()})
			if err != nil {
				logger.Error("Failed to marshal player disconnected message: %v", err)
			} else {
				dispatcher.BroadcastMessage(OpCodePlayerDisconnected, data, nil, nil, true)
			}
		}
	}
	if len(s.Presences) == 0 {
		return nil
	}

	return s
}

// Executed on an interval based on the tick rate returned by MatchInit().
func (m *MatchHandler) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state any, messages []runtime.MatchData) any {
	s := state.(*MatchState)

	if s.State != StatePlaying {
		// Stop tick loop if match ended or waiting
		return s
	}

	// 1. Process Disconnections
	now := time.Now().Unix()
	for userID, deadline := range s.Disconnects {
		if now >= deadline {
			s.State = StateEnded
			if userID == s.PlayerX {
				s.Winner = PlayerO
				s.broadcastGameEnded(dispatcher, PlayerO, "Player X disconnected")
				updatePlayerStats(ctx, logger, nk, s.PlayerO, "win")
				updatePlayerStats(ctx, logger, nk, s.PlayerX, "loss")
			} else {
				s.Winner = PlayerX
				s.broadcastGameEnded(dispatcher, PlayerX, "Player O disconnected")
				updatePlayerStats(ctx, logger, nk, s.PlayerX, "win")
				updatePlayerStats(ctx, logger, nk, s.PlayerO, "loss")
			}
			return s
		}
	}

	// 2. Process Turn Timeouts
	if s.IsTimedMode && s.Deadline > 0 {
		now := time.Now().Unix()
		if now >= s.Deadline {
			s.State = StateEnded
			// Forfeit for the player whose turn it was
			if s.Turn == PlayerX {
				s.Winner = PlayerO
				s.broadcastGameEnded(dispatcher, PlayerO, "Player X timeout")

				updatePlayerStats(ctx, logger, nk, s.PlayerO, "win")
				updatePlayerStats(ctx, logger, nk, s.PlayerX, "loss")
			} else {
				s.Winner = PlayerX
				s.broadcastGameEnded(dispatcher, PlayerX, "Player O timeout")

				updatePlayerStats(ctx, logger, nk, s.PlayerX, "win")
				updatePlayerStats(ctx, logger, nk, s.PlayerO, "loss")
			}
			return s
		}
	}

	// 2. Process Messages
	for _, msg := range messages {
		if msg.GetOpCode() == OpCodeMove {
			// Validate whose turn it is
			isPlayerX := msg.GetUserId() == s.PlayerX
			isPlayerO := msg.GetUserId() == s.PlayerO

			if (s.Turn == PlayerX && !isPlayerX) || (s.Turn == PlayerO && !isPlayerO) {
				continue // Not their turn
			}

			var moveData MoveMessage
			if err := json.Unmarshal(msg.GetData(), &moveData); err != nil {
				continue // Invalid format
			}

			// Validate move (0-8) and check if cell is empty
			pos := moveData.Position
			if pos < 0 || pos > 8 || s.Board[pos] != 0 {
				continue // Invalid move
			}

			// Apply move
			s.Board[pos] = s.Turn

			// Check win / draw
			winner := checkWin(s.Board)
			if winner != 0 {
				s.State = StateEnded
				s.Winner = winner
				s.broadcastState(dispatcher, OpCodeStateUpdate) // Send board state before end

				reason := "Player X wins!"
				if winner == PlayerO {
					reason = "Player O wins!"
				}

				s.broadcastGameEnded(dispatcher, winner, reason)
				if winner == PlayerX {
					updatePlayerStats(ctx, logger, nk, s.PlayerX, "win")
					updatePlayerStats(ctx, logger, nk, s.PlayerO, "loss")
				} else {
					updatePlayerStats(ctx, logger, nk, s.PlayerO, "win")
					updatePlayerStats(ctx, logger, nk, s.PlayerX, "loss")
				}
			} else if isBoardFull(s.Board) {
				s.State = StateEnded
				s.Winner = 0 // Draw

				s.broadcastState(dispatcher, OpCodeStateUpdate)
				s.broadcastGameEnded(dispatcher, 0, "Draw")

				updatePlayerStats(ctx, logger, nk, s.PlayerX, "draw")
				updatePlayerStats(ctx, logger, nk, s.PlayerO, "draw")
			} else {
				// Switch turns and reset timer
				if s.Turn == PlayerX {
					s.Turn = PlayerO
				} else {
					s.Turn = PlayerX
				}

				if s.IsTimedMode {
					s.Deadline = time.Now().Add(TurnTimeLimitSec * time.Second).Unix()
				}

				// Broadcast state
				s.broadcastState(dispatcher, OpCodeStateUpdate)
			}
		}
	}

	return s
}

// Called when the match handler receives a runtime signal
func (m *MatchHandler) MatchSignal(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state any, data string) (any, string) {
	return state, "received " + data
}

// Called when the server begins a graceful shutdown process
func (m *MatchHandler) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state any, graceSeconds int) any {
	return state
}

// Helpers
func (s *MatchState) broadcastState(dispatcher runtime.MatchDispatcher, opcode int64) {
	data, err := json.Marshal(s)
	if err == nil {
		dispatcher.BroadcastMessage(opcode, data, nil, nil, true)
	}
}

func (s *MatchState) broadcastGameEnded(dispatcher runtime.MatchDispatcher, winner int, reason string) {
	msg := EndMessage{
		Winner: winner,
		Reason: reason,
	}
	data, err := json.Marshal(msg)
	if err == nil {
		dispatcher.BroadcastMessage(OpCodeGameEnded, data, nil, nil, true)
	}
}

func checkWin(board []int) int {
	winningLines := [][]int{
		// Rows
		{0, 1, 2},
		{3, 4, 5},
		{6, 7, 8},
		// Cols
		{0, 3, 6},
		{1, 4, 7},
		{2, 5, 8},
		// Diagonals
		{0, 4, 8},
		{2, 4, 6},
	}

	for _, line := range winningLines {
		if board[line[0]] != 0 && board[line[0]] == board[line[1]] && board[line[1]] == board[line[2]] {
			return board[line[0]]
		}
	}

	return 0
}

func isBoardFull(board []int) bool {
	for _, cell := range board {
		if cell == 0 {
			return false
		}
	}
	return true
}

func updatePlayerStats(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, userID string, result string) {
	if userID == "" {
		return
	}

	// 1. Read existing stats
	stats := UserStats{Wins: 0, Losses: 0, Draws: 0}
	reads := []*runtime.StorageRead{
		{
			Collection: "stats",
			Key:        "tic-tac-toe",
			UserID:     userID,
		},
	}

	objects, err := nk.StorageRead(ctx, reads)
	if err != nil {
		logger.Error("Error reading stats: %v", err)
		return
	}

	if len(objects) > 0 {
		if err := json.Unmarshal([]byte(objects[0].Value), &stats); err != nil {
			logger.Error("Error unmarshalling stats: %v", err)
		}
	}

	// 2. Increment stats
	switch result {
	case "win":
		stats.Wins++
	case "loss":
		stats.Losses++
	case "draw":
		stats.Draws++
	}

	// 3. Save directly to Storage
	valueBytes, err := json.Marshal(stats)
	if err != nil {
		logger.Error("Error marshalling stats: %v", err)
		return
	}

	writes := []*runtime.StorageWrite{
		{
			Collection:      "stats",
			Key:             "tic-tac-toe",
			UserID:          userID,
			Value:           string(valueBytes),
			PermissionRead:  2, // Public read
			PermissionWrite: 0, // Only server logic write
		},
	}
	if _, err := nk.StorageWrite(ctx, writes); err != nil {
		logger.Error("Error writing stats: %v", err)
	}

	// 4. Update Leaderboard score specifically for "wins"
	if result == "win" {
		score := int64(1)
		subscore := int64(0)
		_, err := nk.LeaderboardRecordWrite(ctx, "tic_tac_toe_global", userID, "", score, subscore, nil, nil)
		if err != nil {
			logger.Error("Error writing to leaderboard: %v", err)
		}
	}
}
