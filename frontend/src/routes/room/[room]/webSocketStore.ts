import { writable, type Writable } from "svelte/store";

export const webSocketConnection: Writable<WebSocket | undefined> = writable()
