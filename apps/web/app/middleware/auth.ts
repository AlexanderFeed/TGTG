// A page opts into this route guard with:
//   definePageMeta({ middleware: 'auth' })
// profile.vue does this, so unauthenticated visitors cannot render the page.
export default defineNuxtRouteMiddleware(async () => {
  const { initialized, isAuthenticated, refreshSession } = useAuth()

  // On a first/direct visit we must ask Go about the cookie. On later client-side
  // navigations initialized avoids repeating the same request unnecessarily.
  if (!initialized.value) {
    await refreshSession()
  }
  if (!isAuthenticated.value) {
    // replace avoids leaving the protected page in browser history.
    return navigateTo('/login', { replace: true })
  }
})
