<script setup lang="ts">
import CharLikeCanvas from './CharLikeCanvas.vue';

const props = defineProps<{
  kanji: string | null;
  isGuesser: boolean;
}>();
</script>

<template>
  <div class="flex flex-col gap-6 p-8">
    <div class="flex items-center gap-3 text-2xl">
      <template v-if="!props.isGuesser && props.kanji !== null">
        <p>お題は</p>
        <CharLikeCanvas>{{ props.kanji }}</CharLikeCanvas>
        <p>です</p>
      </template>
      <p v-else>お題は読み手には表示されません</p>
    </div>
    <slot />
  </div>
</template>
