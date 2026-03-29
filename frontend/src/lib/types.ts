export type Turn = 'X' | 'O';

export type PossibleInputs = 'X' | 'O' | '';
export type BoardInput = PossibleInputs[];

export interface UserStats {
	name: string;
	wins: number;
	losses: number;
	draws: number;
}

export type MatchType = 'classic' | 'timed';
