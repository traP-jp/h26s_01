import { defineStore } from 'pinia';

import type { User } from '@/types/api';

type UserState = {
  currentUser: User | null;
  isLoading: boolean;
  error: string | null;
};

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    currentUser: null,
    isLoading: false,
    error: null,
  }),
  getters: {
    userId: (state) => state.currentUser?.id ?? null,
    isLoggedIn: (state) => state.currentUser !== null,
  },
  actions: {
    setLoading(isLoading: boolean) {
      this.isLoading = isLoading;
    },
    setError(error: string | null) {
      this.error = error;
    },
    setCurrentUser(user: User | null) {
      this.currentUser = user;
    },
  },
});
