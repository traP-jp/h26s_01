import { getSocket } from '@/api/socket';
import { useLobbyStore } from '@/stores/lobby';
import { useRoomStore } from '@/stores/room';
import type { RoomListUpdatedEvent } from '@/types/api';

let hasRegisteredListeners = false;
const roomListUpdatedHandlers = new Set<(event: RoomListUpdatedEvent) => void>();

export const useLobbySocketEvents = () => {
  const lobbyStore = useLobbyStore();
  const roomStore = useRoomStore();
  const socket = getSocket();

  const handleRoomListUpdated = (event: RoomListUpdatedEvent) => {
    lobbyStore.applyRoomListUpdatedEvent(event);

    if (event.eventType === 'room_deleted' && event.roomId === roomStore.roomId) {
      roomStore.resetRoom();
    }

    roomListUpdatedHandlers.forEach((handler) => {
      handler(event);
    });
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
