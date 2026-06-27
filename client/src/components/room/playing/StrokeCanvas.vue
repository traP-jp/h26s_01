<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from 'vue';

import type { DrawStrokeEvent, Stroke } from '@/types/api';

type Point = {
  x: number;
  y: number;
};

const STROKE_ANIMATION_MS = 320;

const props = withDefaults(
  defineProps<{
    strokes: Stroke[];
    canDraw: boolean;
    disabled?: boolean;
  }>(),
  {
    disabled: false,
  },
);

const draftStroke = defineModel<DrawStrokeEvent | null>('draftStroke', {
  required: true,
});
const isDragging = ref(false);
const activePointerId = ref<number | null>(null);
const animatedStrokeKeys = ref<Set<string>>(new Set());
const previousStrokeCount = ref(props.strokes.length);
const animationTimeoutIds = new Set<number>();

const canInteract = computed(() => props.canDraw && !props.disabled);

const clampUnit = (value: number) => Math.min(Math.max(value, 0), 1);

const getSvgPoint = (event: PointerEvent, svg: SVGSVGElement): Point => {
  const rect = svg.getBoundingClientRect();

  if (rect.width === 0 || rect.height === 0) {
    return {
      x: 0,
      y: 0,
    };
  }

  return {
    x: clampUnit((event.clientX - rect.left) / rect.width),
    y: clampUnit((event.clientY - rect.top) / rect.height),
  };
};

const getStrokePath = (stroke: DrawStrokeEvent | Stroke) =>
  `M ${stroke.x1} ${stroke.y1} L ${stroke.x2} ${stroke.y2}`;

const getStrokeKey = (stroke: Stroke, index: number) =>
  `${index}:${stroke.drawerId}:${stroke.x1}:${stroke.y1}:${stroke.x2}:${stroke.y2}`;

const releasePointer = (event: PointerEvent) => {
  const svg = event.currentTarget as SVGSVGElement;

  if (svg.hasPointerCapture(event.pointerId)) {
    svg.releasePointerCapture(event.pointerId);
  }
};

const handlePointerDown = (event: PointerEvent) => {
  if (!canInteract.value || (event.pointerType === 'mouse' && event.button !== 0)) {
    return;
  }

  const svg = event.currentTarget as SVGSVGElement;
  const point = getSvgPoint(event, svg);

  draftStroke.value = {
    x1: point.x,
    x2: point.x,
    y1: point.y,
    y2: point.y,
  };
  isDragging.value = true;
  activePointerId.value = event.pointerId;
  svg.setPointerCapture(event.pointerId);
};

const updateDraftEnd = (event: PointerEvent) => {
  if (
    !isDragging.value ||
    activePointerId.value !== event.pointerId ||
    draftStroke.value === null
  ) {
    return;
  }

  const svg = event.currentTarget as SVGSVGElement;
  const point = getSvgPoint(event, svg);

  draftStroke.value = {
    ...draftStroke.value,
    x2: point.x,
    y2: point.y,
  };
};

const handlePointerMove = (event: PointerEvent) => {
  updateDraftEnd(event);
};

const handlePointerUp = (event: PointerEvent) => {
  updateDraftEnd(event);
  releasePointer(event);
  isDragging.value = false;
  activePointerId.value = null;
};

const handlePointerCancel = (event: PointerEvent) => {
  releasePointer(event);
  draftStroke.value = null;
  isDragging.value = false;
  activePointerId.value = null;
};

watch(
  () => props.strokes.length,
  (strokeCount) => {
    if (strokeCount < previousStrokeCount.value) {
      animatedStrokeKeys.value = new Set();
      previousStrokeCount.value = strokeCount;
      return;
    }

    if (strokeCount === previousStrokeCount.value) {
      return;
    }

    const nextAnimatedStrokeKeys = new Set(animatedStrokeKeys.value);

    props.strokes.slice(previousStrokeCount.value).forEach((stroke, offset) => {
      const index = previousStrokeCount.value + offset;
      const key = getStrokeKey(stroke, index);
      nextAnimatedStrokeKeys.add(key);

      const timeoutId = window.setTimeout(() => {
        const updatedAnimatedStrokeKeys = new Set(animatedStrokeKeys.value);
        updatedAnimatedStrokeKeys.delete(key);
        animatedStrokeKeys.value = updatedAnimatedStrokeKeys;
        animationTimeoutIds.delete(timeoutId);
      }, STROKE_ANIMATION_MS);
      animationTimeoutIds.add(timeoutId);
    });

    animatedStrokeKeys.value = nextAnimatedStrokeKeys;
    previousStrokeCount.value = strokeCount;
  },
);

onBeforeUnmount(() => {
  animationTimeoutIds.forEach((timeoutId) => {
    window.clearTimeout(timeoutId);
  });
  animationTimeoutIds.clear();
});
</script>

<template>
  <svg
    class="aspect-square w-full touch-none border border-primary bg-background text-primary"
    viewBox="0 0 1 1"
    @pointercancel="handlePointerCancel"
    @pointerdown="handlePointerDown"
    @pointermove="handlePointerMove"
    @pointerup="handlePointerUp"
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
      v-for="(stroke, index) in strokes"
      :key="getStrokeKey(stroke, index)"
      class="pointer-events-none stroke-canvas__stroke"
      :class="{
        'stroke-canvas__stroke--animated': animatedStrokeKeys.has(getStrokeKey(stroke, index)),
      }"
      :d="getStrokePath(stroke)"
      pathLength="1"
    />

    <path
      v-if="draftStroke !== null"
      class="pointer-events-none opacity-60"
      :d="getStrokePath(draftStroke)"
      fill="none"
      stroke="currentColor"
      stroke-dasharray="0.03 0.02"
      stroke-linecap="round"
      stroke-width="0.012"
    />
  </svg>
</template>

<style scoped>
.stroke-canvas__stroke {
  fill: none;
  stroke: currentColor;
  stroke-dasharray: 1;
  stroke-dashoffset: 0;
  stroke-linecap: round;
  stroke-width: 0.014;
}

.stroke-canvas__stroke--animated {
  animation: stroke-canvas-draw 320ms ease-out;
}

@keyframes stroke-canvas-draw {
  from {
    stroke-dashoffset: 1;
  }

  to {
    stroke-dashoffset: 0;
  }
}
</style>
