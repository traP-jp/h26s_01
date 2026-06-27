<script lang="ts" setup>
import BaseButton from '../common/BaseButton.vue';
import UserIcon from '../common/UserIcon.vue';
</script>

<template>
  <li class="flex items-center justify-between text-2xl px-8 border-l-4 gap-5">
    <p>一つ目の部屋</p>
    <div class="grid grid-cols-3 items-center gap-3 w-xl">
      <div class="flex justify-center gap-1">
        <UserIcon tra-qid="cuptea" size="medium" />
        <UserIcon tra-qid="yasako" size="medium" />
        <UserIcon tra-qid="yadorigi" size="medium" />
      </div>
      <p class="text-center">三人 / 九人</p>
      <BaseButton variant="primary">参加</BaseButton>
    </div>
  </li>
</template>
