<script setup lang="ts">
import { storeToRefs } from 'pinia';

import { useUserStore } from '@/stores/user.ts';

import UserIcon from './UserIcon.vue';

const { userId } = storeToRefs(useUserStore());
</script>

<template>
  <header
    class="sticky top-0 bg-primary px-8 py-3 flex justify-between text-background items-center"
  >
    <img src="/public/favicon.ico" alt="以線伝心のロゴ" class="size-24" />
    <div class="flex gap-3 items-center" v-if="userId">
      <UserIcon :tra-qid="userId" size="small" />
      <p>{{ userId }}</p>
    </div>
  </header>
</template>
