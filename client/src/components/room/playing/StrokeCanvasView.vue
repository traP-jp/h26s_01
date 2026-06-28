<script setup lang="ts">
import type { Stroke } from '@/types/api';

const props = defineProps<{
  strokes: Stroke[];
}>();

const getStrokePath = (stroke: Stroke) => `M ${stroke.x1} ${stroke.y1} L ${stroke.x2} ${stroke.y2}`;

const getStrokeKey = (stroke: Stroke, index: number) =>
  `${index}:${stroke.drawerId}:${stroke.x1}:${stroke.y1}:${stroke.x2}:${stroke.y2}`;
</script>

<template>
  <svg
    class="aspect-square w-full border border-primary bg-background text-primary"
    viewBox="0 0 1 1"
  >
    <line
      class="pointer-events-none opacity-20"
      stroke="currentColor"
      stroke-width="0.003"
      x1="0.5"
      x2="0.5"
      y1="0"
      y2="1"
    />
    <line
      class="pointer-events-none opacity-20"
      stroke="currentColor"
      stroke-width="0.003"
      x1="0"
      x2="1"
      y1="0.5"
      y2="0.5"
    />

    <path
      v-for="(stroke, index) in props.strokes"
      :key="getStrokeKey(stroke, index)"
      class="pointer-events-none"
      :d="getStrokePath(stroke)"
      fill="none"
      stroke="currentColor"
      stroke-linecap="round"
      stroke-width="0.014"
    />
  </svg>
</template>
