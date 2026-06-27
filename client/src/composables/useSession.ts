import { useMe } from './useMe';
import { useSocketConnection } from './useSocketConnection';

export const useSession = () => {
  const { fetchMe, userStore } = useMe();
  const socketConnection = useSocketConnection();

  const initializeSession = async () => {
    const user = await fetchMe();

    if (user === null) {
      return null;
    }

    socketConnection.connect();
    return user;
  };

  return {
    initializeSession,
    socketConnection,
    userStore,
  };
};
