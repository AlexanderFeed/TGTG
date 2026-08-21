<script setup lang="ts">
import { LocateFixed, Minus, Navigation, Plus } from 'lucide-vue-next'
import type { Offer } from '~/types/marketplace'

defineProps<{
  offers: Offer[]
  activeId?: string
}>()

const emit = defineEmits<{
  select: [id: string]
}>()
</script>

<template>
  <div class="market-map" aria-label="Демонстрационная карта предложений">
    <div class="market-map__river" />
    <div class="market-map__road market-map__road--one" />
    <div class="market-map__road market-map__road--two" />
    <div class="market-map__road market-map__road--three" />
    <div class="market-map__road market-map__road--four" />

    <button
      v-for="offer in offers"
      :key="offer.id"
      class="map-pin"
      :class="{ 'map-pin--active': activeId === offer.id }"
      :style="{ left: `${offer.marker.x}%`, top: `${offer.marker.y}%` }"
      type="button"
      :aria-label="`${offer.merchant}, ${offer.price} рублей`"
      @click="emit('select', offer.id)"
    >
      <span>{{ offer.price }} ₽</span>
      <i />
    </button>

    <div class="market-map__label market-map__label--one">Чистые пруды</div>
    <div class="market-map__label market-map__label--two">Покровка</div>
    <div class="market-map__label market-map__label--three">Басманный</div>

    <div class="market-map__controls">
      <button type="button" aria-label="Приблизить"><Plus :size="18" /></button>
      <button type="button" aria-label="Отдалить"><Minus :size="18" /></button>
    </div>
    <button class="market-map__locate" type="button" aria-label="Моё местоположение">
      <LocateFixed :size="20" />
    </button>
    <div class="market-map__me"><Navigation :size="14" fill="currentColor" /></div>
  </div>
</template>

<style scoped>
.market-map {
  position: relative;
  min-height: 570px;
  overflow: hidden;
  border-radius: 28px;
  background-color: #e8eee6;
  background-image:
    linear-gradient(33deg, transparent 48%, rgba(255,255,255,.75) 49%, rgba(255,255,255,.75) 52%, transparent 53%),
    linear-gradient(118deg, transparent 48%, rgba(255,255,255,.66) 49%, rgba(255,255,255,.66) 52%, transparent 53%),
    radial-gradient(circle at 25% 20%, rgba(176, 204, 172, .45), transparent 18%),
    radial-gradient(circle at 70% 76%, rgba(176, 204, 172, .5), transparent 24%);
  background-size: 150px 150px, 190px 190px, auto, auto;
  box-shadow: inset 0 0 0 1px rgba(23, 62, 50, 0.06);
}

.market-map__river {
  position: absolute;
  top: -20%;
  left: 8%;
  width: 14%;
  height: 150%;
  border: 12px solid rgba(206, 229, 236, 0.9);
  border-top: 0;
  border-bottom: 0;
  border-radius: 45%;
  transform: rotate(26deg);
}

.market-map__road {
  position: absolute;
  height: 5px;
  border: 1px solid rgba(210, 209, 197, 0.8);
  background: rgba(255,255,255,.9);
}

.market-map__road--one { top: 24%; left: -10%; width: 120%; transform: rotate(-8deg); }
.market-map__road--two { top: 61%; left: -10%; width: 120%; transform: rotate(11deg); }
.market-map__road--three { top: 0; left: 59%; width: 78%; transform: rotate(83deg); }
.market-map__road--four { top: 0; left: 31%; width: 82%; transform: rotate(94deg); }

.map-pin {
  position: absolute;
  z-index: 4;
  display: grid;
  padding: 8px 11px;
  border: 3px solid var(--white);
  border-radius: 999px;
  color: var(--white);
  background: var(--forest-900);
  box-shadow: 0 8px 22px rgba(16, 43, 36, .27);
  cursor: pointer;
  font-size: .72rem;
  font-weight: 900;
  transform: translate(-50%, -50%);
  transition: transform 160ms ease, background 160ms ease;
}

.map-pin i {
  position: absolute;
  right: 50%;
  bottom: -7px;
  width: 9px;
  height: 9px;
  border-right: 3px solid var(--white);
  border-bottom: 3px solid var(--white);
  background: inherit;
  transform: translateX(50%) rotate(45deg);
}

.map-pin:hover,
.map-pin--active {
  z-index: 5;
  background: var(--coral-500);
  transform: translate(-50%, -50%) scale(1.12);
}

.market-map__label {
  position: absolute;
  color: rgba(43, 72, 62, .54);
  font-size: .75rem;
  font-weight: 750;
  letter-spacing: .02em;
}

.market-map__label--one { top: 15%; left: 44%; }
.market-map__label--two { top: 48%; left: 12%; }
.market-map__label--three { right: 12%; bottom: 15%; }

.market-map__controls {
  position: absolute;
  z-index: 6;
  top: 18px;
  right: 18px;
  display: grid;
  overflow: hidden;
  border-radius: 12px;
  background: var(--white);
  box-shadow: var(--shadow-sm);
}

.market-map__controls button,
.market-map__locate {
  display: grid;
  width: 40px;
  height: 40px;
  place-items: center;
  color: var(--forest-900);
  background: var(--white);
  cursor: pointer;
}

.market-map__controls button + button {
  border-top: 1px solid var(--sand-200);
}

.market-map__locate {
  position: absolute;
  z-index: 6;
  right: 18px;
  bottom: 18px;
  border-radius: 50%;
  box-shadow: var(--shadow-sm);
}

.market-map__me {
  position: absolute;
  z-index: 4;
  top: 49%;
  left: 48%;
  display: grid;
  width: 32px;
  height: 32px;
  place-items: center;
  border: 4px solid var(--white);
  border-radius: 50%;
  color: var(--white);
  background: #3389d8;
  box-shadow: 0 0 0 12px rgba(51, 137, 216, .15);
}

@media (max-width: 720px) {
  .market-map {
    min-height: 390px;
    border-radius: 22px;
  }
}
</style>
