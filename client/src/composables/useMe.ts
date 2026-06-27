import { apiClient } from '@/api/client';
import { getApiErrorMessage } from '@/api/error';
import { useUserStore } from '@/stores/user';

export const useMe = () => {
  const userStore = useUserStore();

  const fetchMe = async () => {
    userStore.setLoading(true);
    userStore.setError(null);

    try {
      const { data, error, response } = await apiClient.GET('/api/me');

      if (error) {
        userStore.setCurrentUser(null);
        userStore.setError(getApiErrorMessage(error, 'ユーザー情報の取得に失敗しました', response));
        return null;
      }

      userStore.setCurrentUser(data);
      return data;
    } catch (error) {
      userStore.setCurrentUser(null);
      userStore.setError(getApiErrorMessage(error, 'ユーザー情報の取得に失敗しました'));
      return null;
    } finally {
      userStore.setLoading(false);
    }
  };

  return {
    fetchMe,
    userStore,
  };
};
