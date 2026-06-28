<script lang="ts" setup>
import UserIcon from '@/components/common/UserIcon.vue';
import type { User } from '@/types/api';

defineProps<{ user: User }>();
</script>

<template>
  <div class="flex gap-3 items-center text-xl">
    <UserIcon :tra-qid="user.id" size="small" />
    <p>{{ user.id }} さん</p>
  </div>
</template>
