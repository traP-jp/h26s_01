<script setup lang="ts">
import BaseButton from '@/components/common/BaseButton.vue';
import TheHeader from '@/components/common/TheHeader.vue';
import UserIcon from '@/components/common/UserIcon.vue';
</script>

<template>
  <div class="min-h-screen bg-background">
    <TheHeader />
    <div>
      <div class="flex flex-row gap-24 items-baseline mt-20 ml-28">
        <p class="text-6xl text-primary">部屋名</p>
        <p class="text-4xl text-primary">六人 / 八人</p>
      </div>
      <div class="grid grid-cols-4 gap-x-28 gap-y-12 mt-16 ml-32 w-fit">
        <div v-for="n in 7" :key="n" class="flex flex-col items-center w-fit gap-4">
          <UserIcon tra-qid="" size="large" />
          <p class="text-3xl text-center text-primary">traQid</p>
        </div>
      </div>
    </div>
    <div class="fixed bottom-16 right-24">
      <BaseButton variant="primary"><p class="text-background">準備完了</p></BaseButton>
    </div>
  </div>
</template>
