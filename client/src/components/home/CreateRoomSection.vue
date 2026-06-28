<script setup lang="ts">
import { computed, ref, watch } from 'vue';

import BaseButton from '../common/BaseButton.vue';

const props = withDefaults(
  defineProps<{
    disabled?: boolean;
    initialRoomName?: string;
  }>(),
  {
    disabled: false,
    initialRoomName: '',
  },
);

const emit = defineEmits<{
  create: [roomName: string];
}>();

const inputtedRoomName = ref(props.initialRoomName);
const trimmedRoomName = computed(() => inputtedRoomName.value.trim());
const canSubmit = computed(() => !props.disabled && trimmedRoomName.value.length > 0);

watch(
  () => props.initialRoomName,
  (initialRoomName) => {
    if (inputtedRoomName.value.length === 0) {
      inputtedRoomName.value = initialRoomName;
    }
  },
);

const handleSubmit = () => {
  if (!canSubmit.value) {
    return;
  }

  emit('create', trimmedRoomName.value);
};
</script>

<template>
  <section class="space-y-5 p-8">
    <h2 class="text-3xl">新たに作成</h2>
    <form class="flex gap-5 pr-8" @submit.prevent="handleSubmit">
      <input
        v-model="inputtedRoomName"
        class="outline-none text-2xl p-4 border-4 border-primary flex-1"
        type="text"
        placeholder="あなたの部屋の名前"
      />
      <BaseButton variant="primary" :disabled="!canSubmit" @btn-click="handleSubmit"
        >作成</BaseButton
      >
    </form>
  </section>
</template>
