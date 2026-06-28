<script setup lang="ts">
import CharLikeCanvas from './CharLikeCanvas.vue';

const props = defineProps<{
  kanji: string | null;
  isGuesser: boolean;
}>();
</script>

<template>
  <main class="flex min-h-0 min-w-0 flex-col gap-6 overflow-hidden p-8">
    <div class="flex items-center justify-center gap-3 text-2xl">
      <template v-if="!props.isGuesser && props.kanji !== null">
        <p>お題は</p>
        <CharLikeCanvas>{{ props.kanji }}</CharLikeCanvas>
        <p>です</p>
      </template>
      <p v-else>お題は読み手には表示されません</p>
    </div>
    <slot />
  </main>
</template>
