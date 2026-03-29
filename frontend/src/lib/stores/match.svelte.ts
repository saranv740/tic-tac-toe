/**
 * Live match state store
 *
 * Updated exclusively by the WebSocket message handlers in /match/+page.svelte.
 * All values here mirror the backend's MatchState JSON structure.
 */

import type { MatchStatePayload, EndPayload } from '$lib/types';
import type { PossibleInputs, Turn } from '$lib/types';

export interface PlayerInfo {
	userId: string;
	username: string;
	wins: number;
}

export type GamePhase = 'waiting' | 'playing' | 'ended';

/** Maps the server's numeric board (0/1/2) to the UI's string board (''/X/O). */
function mapBoard(serverBoard: number[]): PossibleInputs[] {
	return serverBoard.map((cell) => {
		if (cell === 1) return 'X';
		if (cell === 2) return 'O';
		return '';
	});
}

/** Maps the server's numeric turn (1/2) to the UI Turn type. */
function mapTurn(serverTurn: number): Turn {
	return serverTurn === 1 ? 'X' : 'O';
}

export class MatchStore {
	board = $state<PossibleInputs[]>(Array(9).fill(''));
	turn = $state<Turn>('X');
	deadline = $state<number>(0); // Unix timestamp; 0 = no timer
	playerX = $state<PlayerInfo | null>(null);
	playerO = $state<PlayerInfo | null>(null);
	gamePhase = $state<GamePhase>('waiting');
	result = $state<{ winner: Turn | 'draw' | null; reason: string }>({
		winner: null,
		reason: ''
	});
	opponentDisconnected = $state(false);

	/** Applied when OpCodeGameStart (3) or OpCodeStateUpdate (2) are received. */
	applyState(payload: MatchStatePayload) {
		this.board = mapBoard(payload.board);
		this.turn = mapTurn(payload.turn);
		this.deadline = payload.deadline ?? 0;

		if (payload.state === 2 /* StatePlaying */) {
			this.gamePhase = 'playing';
		}
	}

	/** Called when OpCodeGameStart (3) is received, after player info is fetched. */
	setPlayers(x: PlayerInfo, o: PlayerInfo) {
		this.playerX = x;
		this.playerO = o;
	}

	/** Called when OpCodeGameEnded (4) is received. */
	applyEnd(payload: EndPayload) {
		this.gamePhase = 'ended';
		this.result = {
			winner: payload.winner === 0 ? 'draw' : payload.winner === 1 ? 'X' : 'O',
			reason: payload.reason
		};
	}

	setOpponentDisconnected(v: boolean) {
		this.opponentDisconnected = v;
	}
}
