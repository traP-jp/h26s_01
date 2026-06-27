<script setup lang="ts">
import BaseButton from '../common/BaseButton.vue';
</script>

<template>
  <section class="space-y-5 p-8">
    <h2 class="text-3xl">新たに作成</h2>
    <form class="flex gap-5 pr-8">
      <input
        class="outline-none text-2xl p-4 border-4 border-primary flex-1"
        type="text"
        placeholder="あなたの部屋の名前"
      />
      <BaseButton variant="primary">作成</BaseButton>
    </form>
  </section>
</template>
