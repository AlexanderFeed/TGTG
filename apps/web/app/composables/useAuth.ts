import type { MarketplaceUser } from '~/types/marketplace'

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

const mutationHeaders = {
  'X-Requested-With': 'eshche-est-web',
}

export const useAuth = () => {
  const user = useState<MarketplaceUser | null>('auth:user', () => null)
  const initialized = useState<boolean>('auth:initialized', () => false)
  const requestFetch = useRequestFetch()
  const isAuthenticated = computed(() => Boolean(user.value))

  const errorMessage = (error: unknown, fallback = 'Не удалось выполнить запрос. Попробуйте ещё раз.') => {
    const message = (error as ApiErrorShape)?.data?.error?.message
    return message || fallback
  }

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

  const requestLoginCode = (email: string) => $fetch<ChallengeResponse>('/api/v1/auth/login/request', {
    method: 'POST',
    headers: mutationHeaders,
    body: { email },
  })

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

  const updateProfile = async (patch: Pick<MarketplaceUser, 'name' | 'city'>) => {
    const response = await $fetch<UserResponse>('/api/v1/users/me', {
      method: 'PATCH',
      headers: mutationHeaders,
      body: patch,
    })
    user.value = response.user
    return response.user
  }

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
