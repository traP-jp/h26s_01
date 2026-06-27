import { createRouter, createWebHistory } from 'vue-router';

import PlayingRoomView from '@/components/room/playing/PlayingRoomView.vue';
import HomeView from '@/views/HomeView.vue';
import MockCheckView from '@/views/MockCheckView.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/mock-check',
      name: 'mock-check',
      component: MockCheckView,
    },
    {
      path: '/mock-playing-room',
      component: PlayingRoomView,
    },
  ],
});

export default router;
