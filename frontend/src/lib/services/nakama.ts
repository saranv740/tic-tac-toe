/**
 * The Client instance is created ONCE at module load time.
 */

import { PUBLIC_NAKAMA_HOST, PUBLIC_NAKAMA_PORT, PUBLIC_NAKAMA_SERVER_KEY, PUBLIC_NAKAMA_USE_SSL } from '$env/static/public';
import { Client, type Session, type Socket } from '@heroiclabs/nakama-js';

const HOST = PUBLIC_NAKAMA_HOST ?? 'localhost';
const PORT = PUBLIC_NAKAMA_PORT ?? '7350';
const USE_SSL = PUBLIC_NAKAMA_USE_SSL === 'true';
const SERVER_KEY = PUBLIC_NAKAMA_SERVER_KEY ?? 'defaultkey';
const DEVICE_ID_KEY = 'ttt_device_id';

const client = new Client(SERVER_KEY, HOST, PORT, USE_SSL);

let _session: Session | null = null;
let _socket: Socket | null = null;

/** Returns the cached Nakama session, or creates one via Device auth. */
async function authenticate(): Promise<Session> {
	// Return cached session if still valid
	if (_session && !_session.isexpired(Date.now() / 1000)) {
		return _session;
	}

	let deviceId = localStorage.getItem(DEVICE_ID_KEY);
	if (!deviceId) {
		// Generate a stable random ID for this browser
		deviceId = crypto.randomUUID();
		localStorage.setItem(DEVICE_ID_KEY, deviceId);
	}

	_session = await client.authenticateDevice(deviceId, true);
	return _session;
}

/** Returns the currently cached session (throws if not yet authenticated). */
function getSession(): Session {
	if (!_session) throw new Error('Not authenticated. Call authenticate() first.');
	return _session;
}

/**
 * Returns the (lazily created) WebSocket connection.
 * The socket is created once and reused on all subsequent calls.
 */
async function getSocket(): Promise<Socket> {
	if (_socket) return _socket;

	const session = getSession();
	_socket = client.createSocket(USE_SSL);
	await _socket.connect(session, true);
	return _socket;
}

/** Clears the cached socket so the next getSocket() call creates a fresh connection. */
function resetSocket() {
	_socket = null;
}

export { client, getSession, getSocket, resetSocket, authenticate };
