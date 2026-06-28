<template>
  <p class="bg-primary text-background py-2 px-7 text-lg">
    <slot name="title" />
  </p>
  <div class="min-w-0 px-7 py-4">
    <slot name="default" />
  </div>
</template>
