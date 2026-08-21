export default defineNuxtRouteMiddleware(async () => {
  const { initialized, isAuthenticated, refreshSession } = useAuth()
  if (!initialized.value) {
    await refreshSession()
  }
  if (!isAuthenticated.value) {
    return navigateTo('/login', { replace: true })
  }
})
