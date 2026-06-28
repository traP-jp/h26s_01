import { getSocket } from '@/api/socket';
import { useRoomStore } from '@/stores/room';
import type { RoomUpdatedEvent } from '@/types/api';

let hasRegisteredListeners = false;
let activeRegistrations = 0;
const roomUpdatedHandlers = new Set<(event: RoomUpdatedEvent) => void>();

const handleRoomUpdated = (event: RoomUpdatedEvent) => {
  const roomStore = useRoomStore();
  const currentRoomId = roomStore.roomId;

  if (currentRoomId !== null && event.room.id !== currentRoomId) {
    roomUpdatedHandlers.forEach((handler) => {
      handler(event);
    });
    return;
  }

  roomStore.setRoom(event.room);
  roomUpdatedHandlers.forEach((handler) => {
    handler(event);
  });
};

export const useRoomSocket = () => {
  const roomStore = useRoomStore();
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

    socket.on('room:updated', handleRoomUpdated);
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
