import { useRouter } from 'vue-router';

import { useGameStore } from '@/stores/game';
import { useRoomStore } from '@/stores/room';

import { useSession } from './useSession';

export const useReturnHomeByDisconnect = () => {
  const router = useRouter();
  const roomStore = useRoomStore();
  const gameStore = useGameStore();
  const { initializeSession, socketConnection } = useSession();

  const returnHomeByDisconnect = async () => {
    socketConnection.disconnect();
    roomStore.resetRoom();
    gameStore.resetGame();
    await router.replace('/');
    await initializeSession();
  };

  return {
    returnHomeByDisconnect,
  };
};
