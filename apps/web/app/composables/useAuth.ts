import type { MarketplaceUser } from '~/types/marketplace'

// These TypeScript interfaces describe the JSON contract returned by Go. They
// do not perform runtime validation; they help the editor/compiler catch wrong
// property names while we write frontend code.
interface ChallengeResponse {
  challengeId: string
  expiresInSeconds: number
  delivery: 'email'
  devCode?: string
}

interface UserResponse {
  user: MarketplaceUser
}

interface ApiErrorShape {
  data?: {
    error?: {
      message?: string
    }
  }
}

// Go middleware requires this marker on state-changing requests. All auth and
// profile mutations use this shared object so one request cannot forget it.
const mutationHeaders = {
  'X-Requested-With': 'eshche-est-web',
}

// useAuth is a Nuxt composable: pages/components call it to share authenticated
// user state and the functions that talk to the backend. useState survives page
// navigation and is also transferred from server rendering to the browser.
export const useAuth = () => {
  const user = useState<MarketplaceUser | null>('auth:user', () => null)
  const initialized = useState<boolean>('auth:initialized', () => false)

  // During server-side rendering useRequestFetch forwards the visitor's request
  // headers/cookies to our /api route. Plain $fetch is sufficient for button/form
  // actions that happen in the browser after the page is interactive.
  const requestFetch = useRequestFetch()
  const isAuthenticated = computed(() => Boolean(user.value))

  const errorMessage = (error: unknown, fallback = 'Не удалось выполнить запрос. Попробуйте ещё раз.') => {
    const message = (error as ApiErrorShape)?.data?.error?.message
    return message || fallback
  }

  // Ask the protected Go /auth/me endpoint whether this browser has a valid
  // HttpOnly session. JavaScript cannot inspect that cookie directly by design.
  const refreshSession = async (force = false) => {
    if (initialized.value && !force) return user.value
    try {
      const response = await requestFetch<UserResponse>('/api/v1/auth/me')
      user.value = response.user
    } catch {
      user.value = null
    } finally {
      initialized.value = true
    }
    return user.value
  }

  // Every URL starts with /api, so Nuxt's server/api/[...path].ts receives it
  // first and forwards /v1/... to Go.
  const requestLoginCode = (email: string) => $fetch<ChallengeResponse>('/api/v1/auth/login/request', {
    method: 'POST',
    headers: mutationHeaders,
    body: { email },
  })

  // Successful verification returns the user and a Set-Cookie header. The
  // browser stores the cookie automatically; this code stores only public user
  // data so the UI can react immediately.
  const verifyLoginCode = async (challengeId: string, email: string, code: string) => {
    const response = await $fetch<UserResponse>('/api/v1/auth/login/verify', {
      method: 'POST',
      headers: mutationHeaders,
      body: { challengeId, email, code },
    })
    user.value = response.user
    initialized.value = true
    return response.user
  }

  const requestRegistrationCode = (name: string, email: string) => $fetch<ChallengeResponse>('/api/v1/auth/register/request', {
    method: 'POST',
    headers: mutationHeaders,
    body: { name, email },
  })

  const verifyRegistrationCode = async (challengeId: string, email: string, code: string) => {
    const response = await $fetch<UserResponse>('/api/v1/auth/register/verify', {
      method: 'POST',
      headers: mutationHeaders,
      body: { challengeId, email, code },
    })
    user.value = response.user
    initialized.value = true
    return response.user
  }

  // Pick prevents this call from accidentally submitting read-only user fields
  // such as id, role, email, or impact.
  const updateProfile = async (patch: Pick<MarketplaceUser, 'name' | 'city'>) => {
    const response = await $fetch<UserResponse>('/api/v1/users/me', {
      method: 'PATCH',
      headers: mutationHeaders,
      body: patch,
    })
    user.value = response.user
    return response.user
  }

  // finally clears local UI state even when the network request fails. A failed
  // server request may leave the cookie valid, so refreshSession can still find
  // it on a later full reload; for this pet project the immediate UI exit wins.
  const logout = async () => {
    try {
      await $fetch('/api/v1/auth/logout', {
        method: 'POST',
        headers: mutationHeaders,
      })
    } finally {
      user.value = null
      initialized.value = true
      await navigateTo('/')
    }
  }

  // Returning refs (state) and actions gives every caller the same small auth API.
  return {
    user,
    initialized,
    isAuthenticated,
    errorMessage,
    refreshSession,
    requestLoginCode,
    verifyLoginCode,
    requestRegistrationCode,
    verifyRegistrationCode,
    updateProfile,
    logout,
  }
}
