import { getSocket } from '@/api/socket';
import { useRoomStore } from '@/stores/room';
import type { RoomUpdatedEvent } from '@/types/api';

let hasRegisteredListeners = false;

export const useRoomSocket = () => {
  const roomStore = useRoomStore();
  const socket = getSocket();

  const handleRoomUpdated = (event: RoomUpdatedEvent) => {
    roomStore.setRoom(event.room);
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

  const joinRoom = (roomId: string) => {
    roomStore.setJoining(true);
    roomStore.setError(null);
    socket.emit('room:join', {
      roomId,
    });
    roomStore.setJoining(false);
  };

  const sendReady = () => {
    if (!roomStore.canSendReady) {
      return;
    }

    roomStore.setSendingReady(true);
    roomStore.setError(null);
    socket.emit('game:ready');
    roomStore.markReadySent();
    roomStore.setSendingReady(false);
  };

  return {
    cleanup,
    joinRoom,
    register,
    sendReady,
  };
};
