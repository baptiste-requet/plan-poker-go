<script lang="ts">
	import { createEventDispatcher, onDestroy, onMount } from 'svelte';
	import { webSocketConnection } from './webSocketStore';
	import { MessageType } from './message';
	import type { JoinRoomMessage, Message } from './message';
	import type { CurrentUser } from './room';

	export let currentUser: CurrentUser;
	export let roomCode: string;
	export let webSocketUrl: string;

	const dispatch = createEventDispatcher();

	let socket: WebSocket;

	onMount(() => {
		socket = new WebSocket(webSocketUrl);
		console.debug(`Connecting websocket to: ${webSocketUrl}..`)
		webSocketConnection.set(socket);

		socket.onopen = () => {
			console.log('WebSocket connected!');

			const joinRoomMessage: JoinRoomMessage = {
				type: MessageType.JOIN_ROOM,
				payload: {
					roomCode: roomCode,
					nickname: currentUser.nickname
				}
			};

			socket.send(JSON.stringify(joinRoomMessage));
		};

		socket.onmessage = (event) => {
			const message: Message = JSON.parse(event.data);
			dispatch('message', message);
		};

		socket.onerror = (error) => {
			console.error('WebSocket error:', error);
			webSocketConnection.set(undefined);
		};

		socket.onclose = (event) => {
			console.log('WebSocket connection closed:', event.code, event.reason);
			webSocketConnection.set(undefined);
		};
	});

	onDestroy(() => {
		if (socket !== null && socket !== undefined) {
			socket.close();
		}
	});
</script>
