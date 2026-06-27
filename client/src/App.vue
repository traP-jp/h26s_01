<script setup lang="ts"></script>

<template>
  <main>
    <nav>
      <RouterLink to="/">Home</RouterLink>
      <RouterLink to="/mock-check">Mock Check</RouterLink>
      <RouterLink to="/socket-sample">Socket Sample</RouterLink>
    </nav>

    <RouterView />
  </main>
</template>

<style scoped></style>
