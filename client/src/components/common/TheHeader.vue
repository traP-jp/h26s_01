<script setup lang="ts">
import { storeToRefs } from 'pinia';

import { useUserStore } from '@/stores/user.ts';

import UserIcon from './UserIcon.vue';

const { userId } = storeToRefs(useUserStore());
</script>

<template>
  <header class="sticky top-0 bg-primary p-8 flex justify-between text-background items-center">
    <h1 class="text-2xl">ゲームタイトル</h1>
    <div class="flex gap-3 items-center" v-if="userId">
      <UserIcon :tra-qid="userId" size="small" />
      <p>{{ userId }}</p>
    </div>
  </header>
</template>
