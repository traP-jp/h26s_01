import { getSocket } from '@/api/socket';
import { useRoomStore } from '@/stores/room';
import type { RoomUpdatedEvent } from '@/types/api';

let hasRegisteredListeners = false;
const roomUpdatedHandlers = new Set<(event: RoomUpdatedEvent) => void>();

export const useRoomSocket = () => {
  const roomStore = useRoomStore();
  const socket = getSocket();

  const handleRoomUpdated = (event: RoomUpdatedEvent) => {
    roomStore.setRoom(event.room);
    roomUpdatedHandlers.forEach((handler) => {
      handler(event);
    });
  };

  const register = () => {
    if (hasRegisteredListeners) {
      return;
    }

    socket.on('room:updated', handleRoomUpdated);
    hasRegisteredListeners = true;
  };

  const cleanup = () => {
    socket.off('room:updated', handleRoomUpdated);
    hasRegisteredListeners = false;
  };

  const onRoomUpdated = (handler: (event: RoomUpdatedEvent) => void) => {
    roomUpdatedHandlers.add(handler);
  };

  const offRoomUpdated = (handler: (event: RoomUpdatedEvent) => void) => {
    roomUpdatedHandlers.delete(handler);
  };

  const joinRoom = (roomId: string) => {
    if (roomId.length === 0 || roomStore.currentRoom !== null) {
      return false;
    }

    roomStore.setJoining(true);
    roomStore.setError(null);
    socket.emit('room:join', {
      roomId,
    });
    roomStore.setJoining(false);
    return true;
  };

  const sendReady = () => {
    if (!roomStore.canSendReady) {
      return false;
    }

    roomStore.setSendingReady(true);
    roomStore.setError(null);
    socket.emit('game:ready');
    roomStore.markReadySent();
    roomStore.setSendingReady(false);
    return true;
  };

  return {
    cleanup,
    joinRoom,
    offRoomUpdated,
    onRoomUpdated,
    register,
    sendReady,
  };
};
