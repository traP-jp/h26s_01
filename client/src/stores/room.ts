import { defineStore } from 'pinia';

import type { Room, RoomMember } from '@/types/api';

import { useUserStore } from './user';

type RoomState = {
  currentRoom: Room | null;
  hasSentReady: boolean;
  isJoining: boolean;
  isSendingReady: boolean;
  error: string | null;
};

export const useRoomStore = defineStore('room', {
  state: (): RoomState => ({
    currentRoom: null,
    hasSentReady: false,
    isJoining: false,
    isSendingReady: false,
    error: null,
  }),
  getters: {
    roomId: (state) => state.currentRoom?.id ?? null,
    members: (state) => state.currentRoom?.members ?? [],
    memberCount: (state) => state.currentRoom?.members.length ?? 0,
    isPlaying: (state) => state.currentRoom?.status === 'playing',
    myRoomMember(state): RoomMember | null {
      const userStore = useUserStore();
      const userId = userStore.currentUser?.id;

      if (!userId) {
        return null;
      }

      return state.currentRoom?.members.find((member) => member.id === userId) ?? null;
    },
    isReady(): boolean {
      return this.hasSentReady || (this.myRoomMember?.isReady ?? false);
    },
    allReady: (state) =>
      state.currentRoom !== null &&
      state.currentRoom.members.length > 0 &&
      state.currentRoom.members.every((member) => member.isReady),
    canSendReady(): boolean {
      return (
        this.currentRoom !== null &&
        this.currentRoom.status === 'waiting' &&
        !this.isReady &&
        !this.isSendingReady
      );
    },
  },
  actions: {
    setJoining(isJoining: boolean) {
      this.isJoining = isJoining;
    },
    setSendingReady(isSendingReady: boolean) {
      this.isSendingReady = isSendingReady;
    },
    setError(error: string | null) {
      this.error = error;
    },
    setRoom(room: Room | null) {
      this.currentRoom = room;

      const userId = useUserStore().currentUser?.id;
      const myMember = room?.members.find((member) => member.id === userId);

      if (myMember?.isReady) {
        this.hasSentReady = false;
      }
    },
    markReadySent() {
      this.hasSentReady = true;
    },
    resetRoom() {
      this.currentRoom = null;
      this.hasSentReady = false;
      this.isJoining = false;
      this.isSendingReady = false;
      this.error = null;
    },
  },
});
