<script setup lang="ts">
import { Clock3, Heart, MapPin, Star } from 'lucide-vue-next'
import { formatPrice } from '~/data/marketplace'
import type { Offer } from '~/types/marketplace'

const props = withDefaults(defineProps<{
  offer: Offer
  compact?: boolean
}>(), {
  compact: false,
})

const { isFavorite, toggleFavorite } = useMarketplace()

const discount = computed(() => Math.round((1 - props.offer.price / props.offer.originalPrice) * 100))
</script>

<template>
  <article class="offer-card" :class="{ 'offer-card--compact': compact }">
    <div class="offer-card__visual">
      <NuxtLink :to="`/offers/${offer.id}`" :aria-label="`Открыть ${offer.title}`">
        <img :src="offer.image" :alt="offer.title" class="offer-card__image">
      </NuxtLink>
      <div class="offer-card__discount">−{{ discount }}%</div>
      <button
        class="icon-button offer-card__favorite"
        :class="{ 'offer-card__favorite--active': isFavorite(offer.id) }"
        type="button"
        :aria-label="isFavorite(offer.id) ? 'Убрать из избранного' : 'Добавить в избранное'"
        @click="toggleFavorite(offer.id)"
      >
        <Heart :size="19" :fill="isFavorite(offer.id) ? 'currentColor' : 'none'" />
      </button>
      <span v-if="offer.available <= 2" class="offer-card__scarcity">
        {{ offer.available === 1 ? 'Последний пакет' : `Осталось ${offer.available}` }}
      </span>
    </div>

    <div class="offer-card__body">
      <div class="offer-card__merchant-row">
        <span class="offer-card__merchant">{{ offer.merchant }}</span>
        <span class="offer-card__rating"><Star :size="14" fill="currentColor" /> {{ offer.rating }}</span>
      </div>

      <NuxtLink :to="`/offers/${offer.id}`">
        <h3>{{ offer.title }}</h3>
      </NuxtLink>

      <div class="offer-card__meta">
        <span><Clock3 :size="15" /> {{ offer.pickupWindow }}</span>
        <span><MapPin :size="15" /> {{ offer.distanceKm.toLocaleString('ru-RU') }} км</span>
      </div>

      <div class="offer-card__footer">
        <div class="offer-card__price">
          <strong>{{ formatPrice(offer.price) }}</strong>
          <s>{{ formatPrice(offer.originalPrice) }}</s>
        </div>
        <NuxtLink class="offer-card__arrow" :to="`/offers/${offer.id}`" aria-label="Подробнее">→</NuxtLink>
      </div>
    </div>
  </article>
</template>

<style scoped>
.offer-card {
  overflow: hidden;
  border: 1px solid rgba(225, 222, 211, 0.82);
  border-radius: 24px;
  background: var(--white);
  box-shadow: var(--shadow-sm);
  transition: transform 180ms ease, box-shadow 180ms ease;
}

.offer-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-5px);
}

.offer-card__visual {
  position: relative;
  height: 220px;
  overflow: hidden;
  background: var(--cream-100);
}

.offer-card--compact .offer-card__visual {
  height: 185px;
}

.offer-card__image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 500ms ease;
}

.offer-card:hover .offer-card__image {
  transform: scale(1.035);
}

.offer-card__discount,
.offer-card__scarcity {
  position: absolute;
  z-index: 2;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 900;
}

.offer-card__discount {
  top: 14px;
  left: 14px;
  padding: 7px 10px;
  color: var(--forest-950);
  background: var(--lime-300);
}

.offer-card__favorite {
  position: absolute;
  z-index: 2;
  top: 12px;
  right: 12px;
}

.offer-card__favorite--active {
  color: var(--coral-500);
}

.offer-card__scarcity {
  bottom: 12px;
  left: 12px;
  padding: 7px 11px;
  color: var(--white);
  background: rgba(16, 43, 36, 0.88);
  backdrop-filter: blur(8px);
}

.offer-card__body {
  padding: 18px;
}

.offer-card__merchant-row,
.offer-card__footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.offer-card__merchant {
  overflow: hidden;
  color: var(--forest-700);
  font-size: 0.77rem;
  font-weight: 850;
  letter-spacing: 0.04em;
  text-overflow: ellipsis;
  text-transform: uppercase;
  white-space: nowrap;
}

.offer-card__rating {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  color: #b47905;
  font-size: 0.8rem;
  font-weight: 850;
}

.offer-card h3 {
  margin: 9px 0 13px;
  font-size: 1.12rem;
}

.offer-card__meta {
  display: grid;
  gap: 7px;
  padding-bottom: 16px;
  color: var(--ink-700);
  font-size: 0.82rem;
}

.offer-card__meta span {
  display: flex;
  align-items: center;
  gap: 7px;
}

.offer-card__footer {
  padding-top: 15px;
  border-top: 1px solid var(--cream-100);
}

.offer-card__price {
  display: flex;
  align-items: baseline;
  gap: 8px;
}

.offer-card__price strong {
  color: var(--forest-900);
  font-size: 1.23rem;
}

.offer-card__price s {
  color: var(--ink-500);
  font-size: 0.8rem;
}

.offer-card__arrow {
  display: grid;
  width: 38px;
  height: 38px;
  place-items: center;
  border-radius: 50%;
  color: var(--forest-950);
  background: var(--mint-200);
  font-size: 1.2rem;
  font-weight: 900;
}

@media (max-width: 720px) {
  .offer-card__visual {
    height: 210px;
  }
}
</style>
