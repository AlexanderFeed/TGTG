import { offers } from '~/data/marketplace'

// Unlike useAuth, this composable does not call Go yet. It reads static mock
// offers from app/data/marketplace.ts and stores favorite IDs in a normal
// browser-readable cookie. Connecting the catalogue to GET /api/v1/offers is a
// future integration step recorded in PLAN.md.
export const useMarketplace = () => {
  const favoriteIds = useCookie<string[]>('eshche-est-favorites', {
    default: () => [],
    maxAge: 60 * 60 * 24 * 90,
    sameSite: 'lax',
  })

  const isFavorite = (id: string) => favoriteIds.value.includes(id)

  const toggleFavorite = (id: string) => {
    favoriteIds.value = isFavorite(id)
      ? favoriteIds.value.filter((favoriteId) => favoriteId !== id)
      : [...favoriteIds.value, id]
  }

  const favorites = computed(() => offers.filter((offer) => favoriteIds.value.includes(offer.id)))

  return {
    favoriteIds,
    favorites,
    isFavorite,
    toggleFavorite,
  }
}
