<script setup lang="ts">
import CharLikeCanvas from './CharLikeCanvas.vue';
</script>

<template>
  <div class="flex flex-col p-8">
    <div class="flex items-center gap-3 text-2xl">
      <p>お題は</p>
      <CharLikeCanvas>読</CharLikeCanvas>
      <p>です</p>
    </div>
  </div>
</template>
