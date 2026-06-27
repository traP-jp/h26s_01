import { io } from 'socket.io-client';

import type { AppSocket } from '@/types/socket';

let socket: AppSocket | null = null;

export const getSocket = (): AppSocket => {
  if (socket === null) {
    socket = io({
      autoConnect: false,
      path: '/socket.io',
    });
  }

  return socket;
};

export const connectSocket = (): AppSocket => {
  const currentSocket = getSocket();

  if (!currentSocket.connected) {
    currentSocket.connect();
  }

  return currentSocket;
};

export const disconnectSocket = () => {
  socket?.disconnect();
};
