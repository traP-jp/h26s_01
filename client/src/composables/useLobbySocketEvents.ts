import { getSocket } from '@/api/socket';
import { useLobbyStore } from '@/stores/lobby';
import type { RoomListUpdatedEvent } from '@/types/api';

let hasRegisteredListeners = false;

export const useLobbySocketEvents = () => {
  const lobbyStore = useLobbyStore();
  const socket = getSocket();

  const handleRoomListUpdated = (event: RoomListUpdatedEvent) => {
    lobbyStore.applyRoomListUpdatedEvent(event);
  };

  const register = () => {
    if (hasRegisteredListeners) {
      return;
    }

    socket.on('room_list:updated', handleRoomListUpdated);
    hasRegisteredListeners = true;
  };

  const cleanup = () => {
    socket.off('room_list:updated', handleRoomListUpdated);
    hasRegisteredListeners = false;
  };

  return {
    cleanup,
    register,
  };
};
