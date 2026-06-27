import { apiClient } from '@/api/client';
import { getApiErrorMessage } from '@/api/error';
import { useLobbyStore } from '@/stores/lobby';
import type { CreateRoomRequest } from '@/types/api';

export const useLobby = () => {
  const lobbyStore = useLobbyStore();

  const fetchRooms = async () => {
    lobbyStore.setLoading(true);
    lobbyStore.setError(null);

    try {
      const { data, error, response } = await apiClient.GET('/api/rooms');

      if (error) {
        lobbyStore.setRooms([]);
        lobbyStore.setError(getApiErrorMessage(error, '部屋一覧の取得に失敗しました', response));
        return null;
      }

      lobbyStore.setRooms(data);
      return data;
    } catch (error) {
      lobbyStore.setRooms([]);
      lobbyStore.setError(getApiErrorMessage(error, '部屋一覧の取得に失敗しました'));
      return null;
    } finally {
      lobbyStore.setLoading(false);
    }
  };

  const createRoom = async (request: CreateRoomRequest) => {
    lobbyStore.setLoading(true);
    lobbyStore.setError(null);

    try {
      const { data, error, response } = await apiClient.POST('/api/rooms', {
        body: request,
      });

      if (error) {
        lobbyStore.setError(getApiErrorMessage(error, '部屋の作成に失敗しました', response));
        return null;
      }

      lobbyStore.upsertRoom(data);
      return data;
    } catch (error) {
      lobbyStore.setError(getApiErrorMessage(error, '部屋の作成に失敗しました'));
      return null;
    } finally {
      lobbyStore.setLoading(false);
    }
  };

  return {
    createRoom,
    fetchRooms,
    lobbyStore,
  };
};
