import { readonly, ref } from 'vue';

import {
  connectSocket as connectAppSocket,
  disconnectSocket as disconnectAppSocket,
  getSocket,
} from '@/api/socket';

const isConnected = ref(false);
const connectionError = ref<string | null>(null);
let hasRegisteredListeners = false;

const registerConnectionListeners = () => {
  if (hasRegisteredListeners) {
    return;
  }

  const socket = getSocket();

  socket.on('connect', () => {
    isConnected.value = true;
    connectionError.value = null;
  });

  socket.on('disconnect', () => {
    isConnected.value = false;
  });

  socket.on('connect_error', (error) => {
    isConnected.value = false;
    connectionError.value = error.message;
  });

  hasRegisteredListeners = true;
};

export const useSocketConnection = () => {
  registerConnectionListeners();

  const connect = () => {
    connectionError.value = null;
    const socket = connectAppSocket();
    isConnected.value = socket.connected;
  };

  const disconnect = () => {
    disconnectAppSocket();
    isConnected.value = false;
  };

  return {
    connect,
    connectionError: readonly(connectionError),
    disconnect,
    isConnected: readonly(isConnected),
  };
};
