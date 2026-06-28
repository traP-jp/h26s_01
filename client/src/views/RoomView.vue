<script setup lang="ts">
import LastResult, { type LastResultProps } from '@/components/room/last-result/LastResult.vue';
import ResultMojigoto from '@/components/room/ResultMojigoto.vue';
import WaitingRoom from '@/components/room/waiting/WaitingRoom.vue';

const retest: LastResultProps = {
  successGame: 7,
  loseGame: 2,
};
</script>

<template>
  <WaitingRoom />
  <ResultMojigoto />
  <LastResult :lastData="retest" />
</template>
