/**
 * Global session store
 */

import type { Session } from '@heroiclabs/nakama-js';
import type { Turn } from '$lib/types';

export class SessionStore {
	session = $state<Session | null>(null);
	username = $state<string>('');
	matchId = $state<string | null>(null);
	role = $state<Turn | null>(null);

	clearMatch() {
		this.matchId = null;
		this.role = null;
	}
}

export const sessionStore = new SessionStore();
