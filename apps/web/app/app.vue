<script setup lang="ts">
const { refreshSession } = useAuth()

// app.vue is the frontend application's root entry component. callOnce hydrates
// auth state once during initial rendering, so the header/home page know whether
// to show login actions or the signed-in experience.
await callOnce('auth-session', () => refreshSession())

useHead({
  bodyAttrs: {
    class: 'app-body',
  },
})
</script>

<template>
  <NuxtLoadingIndicator />
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>
