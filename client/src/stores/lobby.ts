import { defineStore } from 'pinia';

import type { Room, RoomListUpdatedEvent } from '@/types/api';

type LobbyState = {
  rooms: Room[];
  isLoading: boolean;
  error: string | null;
};

export const useLobbyStore = defineStore('lobby', {
  state: (): LobbyState => ({
    rooms: [],
    isLoading: false,
    error: null,
  }),
  getters: {
    waitingRooms: (state) => state.rooms.filter((room) => room.status === 'waiting'),
    playingRooms: (state) => state.rooms.filter((room) => room.status === 'playing'),
    getRoomById: (state) => (roomId: string) =>
      state.rooms.find((room) => room.id === roomId) ?? null,
  },
  actions: {
    setLoading(isLoading: boolean) {
      this.isLoading = isLoading;
    },
    setError(error: string | null) {
      this.error = error;
    },
    setRooms(rooms: Room[]) {
      this.rooms = rooms;
    },
    upsertRoom(room: Room) {
      const index = this.rooms.findIndex((item) => item.id === room.id);

      if (index === -1) {
        this.rooms.push(room);
        return;
      }

      this.rooms.splice(index, 1, room);
    },
    removeRoom(roomId: string) {
      this.rooms = this.rooms.filter((room) => room.id !== roomId);
    },
    applyRoomListUpdatedEvent(event: RoomListUpdatedEvent) {
      switch (event.eventType) {
        case 'room_created':
        case 'room_updated':
          this.upsertRoom(event.room);
          break;
        case 'room_deleted':
          this.removeRoom(event.roomId);
          break;
      }
    },
  },
});
