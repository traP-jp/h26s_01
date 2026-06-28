<template>
  <p class="bg-primary text-background py-2 px-6 text-lg">
    <slot name="title" />
  </p>
  <div class="min-w-0 p-6">
    <slot name="default" />
  </div>
</template>
