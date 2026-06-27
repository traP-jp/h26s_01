<template>
  <div
    class="relative grid size-36 place-items-center border-2 border-olive-200 bg-olive-50 before:absolute before:inset-x-0 before:top-1/2 before:h-0.5 before:-translate-y-1/2 before:bg-olive-200 after:absolute after:inset-y-0 after:left-1/2 after:w-0.5 after:-translate-x-1/2 after:bg-olive-200"
  >
    <span class="relative z-10 text-8xl leading-none text-primary">
      <slot />
    </span>
  </div>
</template>
