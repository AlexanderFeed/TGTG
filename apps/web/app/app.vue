<script setup lang="ts">
useHead({
  bodyAttrs: {
    class: 'app-body',
  },
})
</script>

<template>
  <NuxtLoadingIndicator color="#ff6b4a" :height="3" />
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>
