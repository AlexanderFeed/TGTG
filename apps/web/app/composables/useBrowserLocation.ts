export interface BrowserCoordinates {
  latitude: number
  longitude: number
  accuracy: number
}

type LocationStatus = 'idle' | 'loading' | 'ready' | 'denied' | 'unavailable' | 'error'

// Browser location is public application state, not authentication state.
// A guest can grant geolocation permission without creating an account.
export const useBrowserLocation = () => {
  const coordinates = useState<BrowserCoordinates | null>('location:coordinates', () => null)
  const status = useState<LocationStatus>('location:status', () => 'idle')
  const locationError = useState<string>('location:error', () => '')

  const locationLabel = computed(() => {
    if (status.value === 'loading') return 'Определяем место…'
    if (coordinates.value) return 'Моя геопозиция'
    return 'Определить место'
  })

  const requestLocation = async () => {
    // navigator exists only in the browser, not while Nuxt renders on the server.
    if (!import.meta.client) return false

    locationError.value = ''
    if (!window.isSecureContext) {
      status.value = 'unavailable'
      locationError.value = 'Геопозиция доступна на localhost или через HTTPS.'
      return false
    }
    if (!navigator.geolocation) {
      status.value = 'unavailable'
      locationError.value = 'Этот браузер не поддерживает геопозицию.'
      return false
    }

    status.value = 'loading'
    return await new Promise<boolean>((resolve) => {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          coordinates.value = {
            latitude: position.coords.latitude,
            longitude: position.coords.longitude,
            accuracy: position.coords.accuracy,
          }
          status.value = 'ready'
          resolve(true)
        },
        (error) => {
          coordinates.value = null
          if (error.code === error.PERMISSION_DENIED) {
            status.value = 'denied'
            locationError.value = 'Доступ к геопозиции запрещён в настройках браузера.'
          } else if (error.code === error.POSITION_UNAVAILABLE) {
            status.value = 'unavailable'
            locationError.value = 'Браузер не смог определить местоположение.'
          } else {
            status.value = 'error'
            locationError.value = 'Определение местоположения заняло слишком много времени.'
          }
          resolve(false)
        },
        {
          enableHighAccuracy: true,
          timeout: 10_000,
          maximumAge: 5 * 60_000,
        },
      )
    })
  }

  return {
    coordinates,
    status,
    locationError,
    locationLabel,
    requestLocation,
  }
}
