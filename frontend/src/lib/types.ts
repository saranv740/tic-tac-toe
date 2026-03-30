export type Turn = 'X' | 'O';

export type PossibleInputs = 'X' | 'O' | '';
export type BoardInput = PossibleInputs[];

export interface UserStats {
	name: string;
	wins: number;
	losses: number;
	draws: number;
	current_streak: number;
}

export type MatchType = 'classic' | 'timed';
export interface MatchLabel {
	open: 0 | 1;
	mode: MatchType;
	name: string;
}

export interface MatchStatePayload {
	state: number; // 1=waiting, 2=playing, 3=ended
	board: number[]; // 9 elements: 0=empty, 1=playerX, 2=playerO
	turn: number; // 1=playerX, 2=playerO
	player_x: string; // userId
	player_o: string; // userId
	winner: number; // 0=none/draw, 1=X, 2=O
	deadline: number; // Unix timestamp (0 if no timer)
	is_timed_mode: boolean;
}

/** Payload of OpCodeGameEnded (4). */
export interface EndPayload {
	winner: number; // 0=draw, 1=X, 2=O
	reason: string;
}
