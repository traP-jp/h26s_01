<script setup lang="ts">
import { MAX_LIVES } from '@/constants/game';
import { toKanjiNumber } from '@/utils/to-kanji-number';

const props = defineProps<{
  remainingLives: number;
}>();
</script>

<template>
  <div class="flex" role="img" :aria-label="`残機${toKanjiNumber(props.remainingLives)}`">
    <svg
      v-for="n in MAX_LIVES"
      :key="n"
      aria-hidden="true"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
      stroke-width="1"
      stroke="currentColor"
      class="size-18"
      :class="n <= props.remainingLives ? 'fill-primary' : 'fill-transparent opacity-40'"
    >
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12Z"
      />
    </svg>
  </div>
</template>
