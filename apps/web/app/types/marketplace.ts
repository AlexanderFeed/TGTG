export interface MarketplaceUser {
  id: string
  name: string
  email: string
  city: string
  role: 'customer' | 'merchant' | 'admin'
  verifiedAt: string
  createdAt: string
  impact: {
    rescued: number
    savedRubles: number
    co2Kg: number
  }
}

export interface Offer {
  id: string
  title: string
  merchant: string
  category: string
  image: string
  price: number
  originalPrice: number
  pickupWindow: string
  pickupStart: string
  distanceKm: number
  rating: number
  reviewCount: number
  available: number
  address: string
  district: string
  description: string
  contents: string
  tags: string[]
  marker: { x: number; y: number }
  delivery: boolean
}

export interface Category {
  id: string
  label: string
  emoji: string
}

export interface DemoOrder {
  id: string
  offerId: string
  status: 'confirmed' | 'preparing' | 'on_the_way' | 'ready'
  code: string
  eta: string
  mode: 'pickup' | 'delivery'
}
