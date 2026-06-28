<script setup lang="ts">
import UserIcon from '../common/UserIcon.vue';
</script>

<template>
  <div class="w-full h-156 bg-background">
    <p class="text-primary text-6xl ml-48">一文字目</p>
    <div class="flex flex-row gap-24">
      <div class="w-128 h-128 bg-white mt-8 ml-48"></div>
      <!--キャンバス代わりの白正方形-->
      <!--ここにキャンバス-->
      <!-- 回答漢字nなしver.
        <div class="flex flex-col mt-12">
            <p class="text-primary text-8xl">推測成功</p>
            <div class="flex flex-row items-center gap-32 mt-24 ml-8">
                <p class="text-primary text-4xl">お題</p>
                <div class="w-32 h-32 bg-white flex items-center justify-center">
                    <p class="text-primary text-8xl">読</p>
                </div>
            </div>
            <div class="flex flex-row items-center gap-28 mt-16 ml-8">
                <p class="text-primary text-4xl">回答者</p>
                <UserIcon tra-qid="yadorigi" size="medium"/>
            </div> 
      -->
      <div class="flex flex-col">
        <p class="text-primary text-8xl">推測成功</p>
        <div class="flex flex-row items-center gap-16 mt-12 ml-16">
          <p class="text-primary text-4xl">お題</p>
          <div class="w-32 h-32 bg-white flex items-center justify-center">
            <p class="text-primary text-8xl">読</p>
          </div>
        </div>
        <div class="flex flex-row items-center gap-16 mt-4 ml-16">
          <p class="text-primary text-4xl">回答</p>
          <div class="w-32 h-32 bg-white flex items-center justify-center">
            <p class="text-primary text-8xl">読</p>
          </div>
        </div>
        <div class="flex flex-row items-center gap-12 mt-8 ml-16">
          <p class="text-primary text-4xl">回答者</p>
          <UserIcon tra-qid="yadorigi" size="medium" />
        </div>
      </div>
    </div>
  </div>
</template>
