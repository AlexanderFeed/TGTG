import type { MarketplaceUser } from '~/types/marketplace'

interface RegisterInput {
  name: string
  phone: string
  email?: string
}

export const useAuth = () => {
  const user = useCookie<MarketplaceUser | null>('eshche-est-user', {
    default: () => null,
    maxAge: 60 * 60 * 24 * 30,
    sameSite: 'lax',
  })

  const isAuthenticated = computed(() => Boolean(user.value))

  const createDemoUser = (input: RegisterInput): MarketplaceUser => ({
    id: 'demo-user-1',
    name: input.name.trim() || 'Алексей',
    phone: input.phone,
    email: input.email,
    city: 'Москва',
    impact: {
      rescued: 12,
      savedRubles: 6840,
      co2Kg: 18.6,
    },
  })

  const login = (phone: string) => {
    user.value = createDemoUser({
      name: 'Алексей',
      phone,
      email: 'alexey@example.ru',
    })
  }

  const register = (input: RegisterInput) => {
    user.value = createDemoUser(input)
  }

  const updateProfile = (patch: Partial<Pick<MarketplaceUser, 'name' | 'email' | 'city'>>) => {
    if (!user.value) return
    user.value = { ...user.value, ...patch }
  }

  const logout = async () => {
    user.value = null
    await navigateTo('/')
  }

  return {
    user,
    isAuthenticated,
    login,
    register,
    updateProfile,
    logout,
  }
}
