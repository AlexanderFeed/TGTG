<script setup lang="ts">
const { refreshSession } = useAuth()

await callOnce('auth-session', () => refreshSession())

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
