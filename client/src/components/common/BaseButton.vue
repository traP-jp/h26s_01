<script lang="ts" setup>
const props = withDefaults(
  defineProps<{
    variant: 'primary' | 'secondary';
    disabled?: boolean;
  }>(),
  { disabled: false },
);
</script>

<template>
  <button
    type="button"
    :disabled="props.disabled"
    class="text-2xl px-12 py-2 border-4 border-primary transition"
    :class="{
      'opacity-40': props.disabled,
      'cursor-pointer': !props.disabled,
      'text-background bg-primary': props.variant == 'primary',
      'hover:text-primary hover:bg-background': props.variant == 'primary' && !props.disabled,
      'text-primary bg-background': props.variant == 'secondary',
      'hover:text-background hover:bg-primary': props.variant == 'secondary' && !props.disabled,
    }"
  >
    <slot />
  </button>
</template>
