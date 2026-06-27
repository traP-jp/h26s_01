import { getSocket } from '@/api/socket';
import { useLobbyStore } from '@/stores/lobby';
import { useRoomStore } from '@/stores/room';
import type { RoomListUpdatedEvent } from '@/types/api';

let hasRegisteredListeners = false;
let activeRegistrations = 0;
const roomListUpdatedHandlers = new Set<(event: RoomListUpdatedEvent) => void>();

const handleRoomListUpdated = (event: RoomListUpdatedEvent) => {
  const lobbyStore = useLobbyStore();
  const roomStore = useRoomStore();

  lobbyStore.applyRoomListUpdatedEvent(event);

  if (event.eventType === 'room_deleted' && event.roomId === roomStore.roomId) {
    roomStore.resetRoom();
  }

  roomListUpdatedHandlers.forEach((handler) => {
    handler(event);
  });
};

export const useLobbySocketEvents = () => {
  const socket = getSocket();
  let isRegistered = false;

  const register = () => {
    if (isRegistered) {
      return;
    }

    activeRegistrations += 1;
    isRegistered = true;

    if (hasRegisteredListeners) {
      return;
    }

    socket.on('room_list:updated', handleRoomListUpdated);
    hasRegisteredListeners = true;
  };

  const cleanup = () => {
    if (!isRegistered) {
      return;
    }

    activeRegistrations = Math.max(activeRegistrations - 1, 0);
    isRegistered = false;

    if (activeRegistrations > 0) {
      return;
    }

    socket.off('room_list:updated', handleRoomListUpdated);
    hasRegisteredListeners = false;
  };

  const onRoomListUpdated = (handler: (event: RoomListUpdatedEvent) => void) => {
    roomListUpdatedHandlers.add(handler);
  };

  const offRoomListUpdated = (handler: (event: RoomListUpdatedEvent) => void) => {
    roomListUpdatedHandlers.delete(handler);
  };

  return {
    cleanup,
    offRoomListUpdated,
    onRoomListUpdated,
    register,
  };
};
