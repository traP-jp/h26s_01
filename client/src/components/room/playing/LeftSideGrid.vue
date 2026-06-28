<script setup lang="ts">
import UserIcon from '@/components/common/UserIcon.vue';

import LeftSideSection from './LeftSideSection.vue';
import UserCard from './UserCard.vue';
</script>

<template>
  <div class="h-full flex flex-col w-xs border-r-4 border-primary">
    <div class="p-10 text-4xl text-center bg-primary text-background">
      <p>遊戯中</p>
    </div>
    <LeftSideSection>
      <template #title>今の状況</template>
      <template #default>
        <p class="text-2xl">九文字目の一画目<br />二分五十秒が経過</p>
      </template>
    </LeftSideSection>
    <LeftSideSection>
      <template #title>読み手</template>
      <template #default>
        <UserIcon tra-qid="%ここにGuesserのid%" size="small" />
      </template>
    </LeftSideSection>
    <LeftSideSection>
      <template #title>書き手</template>
      <template #default>
        <div class="flex flex-col gap-3">
          <p class="underline text-lg">只今作業中</p>
          <UserIcon tra-qid="%ここに次の書き手のid%" size="small" />
          <p class="text-lg">以降</p>
          <UserIcon tra-qid="%ここに以降の書き手のid%, v-forで回す" size="small" />
        </div>
      </template>
    </LeftSideSection>
  </div>
</template>
