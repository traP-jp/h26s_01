<script lang="ts" setup>
import UserIcon from '@/components/common/UserIcon.vue';
defineProps<{ user: User }>();
</script>

<template>
  <div class="flex gap-3 items-center text-xl">
    <UserIcon tra-qid="cuptea" size="medium" />
    <p>{{ 'cuptea' }} さん</p>
  </div>
</template>
