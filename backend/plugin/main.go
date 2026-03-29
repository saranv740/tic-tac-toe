package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
)

func rpcCreateMatch(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !ok || userID == "" {
		return "", runtime.NewError("must be authenticated", 16) // Unauthenticated
	}

	// Parse payload safely
	var request map[string]any
	if payload != "" {
		if err := json.Unmarshal([]byte(payload), &request); err != nil {
			return "", runtime.NewError("invalid payload format", 3) // Invalid Argument
		}
	}

	mode := "classic"
	if m, ok := request["mode"].(string); ok && m == "timed" {
		mode = m
	}

	// Get username to name the room appropriately
	account, err := nk.AccountGetId(ctx, userID)
	if err != nil {
		logger.Error("Failed to get account: %v", err)
		return "", runtime.NewError("failed to get user account", 13) // Internal
	}

	roomName := account.User.Username + "'s room"

	matchID, err := nk.MatchCreate(ctx, "tic-tac-toe", map[string]any{
		"mode": mode,
		"name": roomName,
	})

	if err != nil {
		logger.Error("Failed to create match: %v", err)
		return "", runtime.NewError("failed to create match", 13) // Internal
	}

	response, _ := json.Marshal(map[string]string{
		"match_id": matchID,
	})

	return string(response), nil
}

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	initStart := time.Now()

	// 1. Register the Match Handler
	if err := initializer.RegisterMatch("tic-tac-toe", func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
		return &MatchHandler{}, nil
	}); err != nil {
		logger.Error("Failed to register match: %v", err)
		return err
	}

	// 1.5 Register the Match Creation RPC
	if err := initializer.RegisterRpc("create_match", rpcCreateMatch); err != nil {
		logger.Error("Failed to register rpc create_match: %v", err)
		return err
	}

	// 2. Register Matchmaker Matched Hook
	// This automatically creates a server-authoritative match when users are paired
	err := initializer.RegisterMatchmakerMatched(func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, entries []runtime.MatchmakerEntry) (string, error) {
		mode := "classic"
		if len(entries) > 0 {
			props := entries[0].GetProperties()
			if props != nil {
				if m, ok := props["mode"].(string); ok {
					mode = m
				}
			}
		}

		matchId, err := nk.MatchCreate(ctx, "tic-tac-toe", map[string]any{
			"mode": mode,
		})

		if err != nil {
			logger.Error("Error creating match: %v", err)
			return "", err
		}

		return matchId, nil
	})

	if err != nil {
		logger.Error("Failed to register MatchmakerMatched hook: %v", err)
		return err
	}

	// 3. Set up the Leaderboard
	leaderboardId := "tic_tac_toe_global"
	authoritative := true
	sortOrder := "desc"
	operator := "incr"
	resetSchedule := ""
	metadata := make(map[string]any)

	if err := nk.LeaderboardCreate(ctx, leaderboardId, authoritative, sortOrder, operator, resetSchedule, metadata, false); err != nil {
		logger.Error("Failed to create leaderboard: %v", err)
		return err
	}

	logger.Info("Plugin loaded in '%d' msec.", time.Now().Sub(initStart).Milliseconds())
	return nil
}
