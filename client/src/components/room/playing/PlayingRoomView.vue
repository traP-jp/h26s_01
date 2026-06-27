<script setup lang="ts">
import BottomGrid from './BottomGrid.vue';
import CenterGrid from './CenterGrid.vue';
import LeftSideGrid from './LeftSideGrid.vue';
import RightSideGrid from './RightSideGrid.vue';
</script>

<template>
  <div class="h-dvh flex flex-col">
    <div class="flex-1 flex justify-between">
      <LeftSideGrid />
      <CenterGrid />
      <RightSideGrid />
    </div>
    <BottomGrid />
  </div>
</template>
