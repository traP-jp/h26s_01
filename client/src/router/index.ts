import { createRouter, createWebHistory } from 'vue-router';

import HomeView from '@/views/HomeView.vue';
import MockCheckView from '@/views/MockCheckView.vue';
import RoomView from '@/views/RoomView.vue';
import SocketSampleView from '@/views/SocketSampleView.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/room/:roomId',
      name: 'room',
      component: RoomView,
    },
    {
      path: '/mock-check',
      name: 'mock-check',
      component: MockCheckView,
    },
    {
      path: '/socket-sample',
      name: 'socket-sample',
      component: SocketSampleView,
    },
  ],
});

export default router;
