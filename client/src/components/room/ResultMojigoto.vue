<script setup lang="ts">
import BaseButton from '../common/BaseButton.vue';
import Result1Moji from './Result1Moji.vue';
</script>

<template>
  <div class="min-h-screen bg-background">
    <div class="w-full h-1 bg-background mt-16"></div>
    <!--上に余白できない問題の一時的解決-->
    <div class="flex flex-row mt-16 ml-24 gap-208">
      <p class="text-primary text-6xl">途中結果</p>
      <BaseButton variant="primary">次の問題へ</BaseButton>
    </div>
    <div class="w-full h-1 bg-primary mt-16"></div>
    <div class="mt-12">
      <Result1Moji />
    </div>
  </div>
</template>
