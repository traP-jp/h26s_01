<script setup lang="ts">
import BaseButton from '@/components/common/BaseButton.vue';
import UserIcon from '@/components/common/UserIcon.vue';

import CharLikeCanvas from './CharLikeCanvas.vue';
import LeftSideSection from './LeftSideSection.vue';
import LifeGage from './LifeGage.vue';
import UserCard from './UserCard.vue';
</script>

<template>
  <div class="flex flex-col w-xs p-8 justify-between">
    <LifeGage />
    <div class="flex flex-col gap-3">
      <BaseButton variant="secondary">書き直し</BaseButton>
      <BaseButton variant="primary">確定</BaseButton>
    </div>
  </div>
</template>
