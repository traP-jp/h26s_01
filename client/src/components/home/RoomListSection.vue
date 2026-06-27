<script lang="ts" setup>
import RoomListItem from './RoomListItem.vue';
import RoomListItem_already from './RoomListItem_already.vue';
</script>

<template>
  <section class="space-y-5 p-8">
    <h2 class="text-3xl">部屋の一覧</h2>
    <ul class="flex flex-col gap-2">
      <RoomListItem />
      <RoomListItem_already />
    </ul>
  </section>
</template>
