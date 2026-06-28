<script setup lang="ts">
import BottomGrid from '@/components/room/playing/BottomGrid.vue';
import CenterGrid from '@/components/room/playing/CenterGrid.vue';
import LeftSideGrid from '@/components/room/playing/LeftSideGrid.vue';
import RightSideGrid from '@/components/room/playing/RightSideGrid.vue';
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
