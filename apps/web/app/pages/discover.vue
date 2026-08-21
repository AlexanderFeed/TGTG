<script setup lang="ts">
import { Filter, List, LocateFixed, Map, MapPin, Search, SlidersHorizontal } from 'lucide-vue-next'
import { offers } from '~/data/marketplace'

const selectedCategory = ref('all')
const activeId = ref(offers[0]?.id || '')
const query = ref('')
const mobileView = ref<'map' | 'list'>('map')
const radius = ref('3')

const filteredOffers = computed(() => {
  const search = query.value.trim().toLocaleLowerCase('ru-RU')
  return offers.filter((offer) => {
    const matchesCategory = selectedCategory.value === 'all' || offer.category === selectedCategory.value
    const matchesSearch = !search || [offer.title, offer.merchant, offer.district].some(value => value.toLocaleLowerCase('ru-RU').includes(search))
    const matchesRadius = offer.distanceKm <= Number(radius.value)
    return matchesCategory && matchesSearch && matchesRadius
  })
})

const activeOffer = computed(() => filteredOffers.value.find(offer => offer.id === activeId.value) || filteredOffers.value[0])

watch(filteredOffers, (next) => {
  if (!next.some(offer => offer.id === activeId.value)) {
    activeId.value = next[0]?.id || ''
  }
})

useSeoMeta({
  title: 'Еда рядом — ЕщёЕсть',
  description: 'Посмотрите актуальные сюрприз-пакеты на карте рядом с вами.',
})
</script>

<template>
  <section class="discover-page">
    <div class="container discover-page__heading">
      <div>
        <p class="eyebrow">Карта предложений</p>
        <h1>Что есть рядом</h1>
        <p>Чистые пруды · Москва <button type="button">Изменить</button></p>
      </div>
      <div class="discover-page__view-switch" role="group" aria-label="Вид результатов">
        <button :class="{ active: mobileView === 'map' }" type="button" @click="mobileView = 'map'"><Map :size="17" /> Карта</button>
        <button :class="{ active: mobileView === 'list' }" type="button" @click="mobileView = 'list'"><List :size="17" /> Список</button>
      </div>
    </div>

    <div class="container discover-controls">
      <div class="discover-search input-wrap">
        <Search :size="18" />
        <input v-model="query" class="input input--icon" type="search" placeholder="Место, блюдо или район">
      </div>
      <label class="discover-radius">
        <LocateFixed :size="17" />
        <span>Радиус</span>
        <select v-model="radius" aria-label="Радиус поиска">
          <option value="1">до 1 км</option>
          <option value="3">до 3 км</option>
          <option value="5">до 5 км</option>
        </select>
      </label>
      <button class="button button--secondary discover-filter" type="button"><SlidersHorizontal :size="18" /> Фильтры</button>
    </div>

    <div class="container discover-categories"><CategoryScroller v-model="selectedCategory" /></div>

    <div class="container discover-layout">
      <div class="discover-results" :class="{ 'discover-results--mobile-hidden': mobileView !== 'list' }">
        <div class="discover-results__summary">
          <span>{{ filteredOffers.length }} предложений</span>
          <button type="button"><Filter :size="15" /> Сначала ближе</button>
        </div>

        <div v-if="filteredOffers.length" class="discover-results__list">
          <button
            v-for="offer in filteredOffers"
            :key="offer.id"
            class="map-result"
            :class="{ 'map-result--active': activeOffer?.id === offer.id }"
            type="button"
            @click="activeId = offer.id; mobileView = 'map'"
          >
            <img :src="offer.image" :alt="offer.title">
            <span class="map-result__body">
              <small>{{ offer.merchant }}</small>
              <strong>{{ offer.title }}</strong>
              <span><MapPin :size="13" /> {{ offer.distanceKm.toLocaleString('ru-RU') }} км · {{ offer.pickupStart }}</span>
            </span>
            <b>{{ offer.price }} ₽</b>
          </button>
        </div>
        <div v-else class="discover-empty">
          <span>🧺</span>
          <h3>В этом радиусе пока пусто</h3>
          <p>Увеличьте расстояние или попробуйте другую категорию.</p>
        </div>
      </div>

      <div class="discover-map-wrap" :class="{ 'discover-map-wrap--mobile-hidden': mobileView !== 'map' }">
        <MarketplaceMap :offers="filteredOffers" :active-id="activeOffer?.id" @select="activeId = $event" />

        <article v-if="activeOffer" class="map-preview">
          <img :src="activeOffer.image" :alt="activeOffer.title">
          <div>
            <small>{{ activeOffer.merchant }}</small>
            <strong>{{ activeOffer.title }}</strong>
            <span>{{ activeOffer.pickupWindow }}</span>
          </div>
          <div class="map-preview__price"><b>{{ activeOffer.price }} ₽</b><s>{{ activeOffer.originalPrice }} ₽</s></div>
          <NuxtLink class="button button--primary" :to="`/offers/${activeOffer.id}`">Открыть</NuxtLink>
        </article>
      </div>
    </div>
  </section>
</template>

<style scoped>
.discover-page { padding: 38px 0 70px; }
.discover-page__heading { display: flex; align-items: end; justify-content: space-between; gap: 24px; margin-bottom: 26px; }
.discover-page h1 { margin-bottom: 10px; font-size: clamp(2.25rem, 5vw, 4.2rem); }
.discover-page__heading p:last-child { margin: 0; color: var(--ink-700); }
.discover-page__heading p button { margin-left: 7px; padding: 0; color: var(--forest-700); background: transparent; cursor: pointer; font-size: .82rem; font-weight: 850; text-decoration: underline; }
.discover-page__view-switch { display: none; padding: 4px; border-radius: 14px; background: var(--cream-100); }
.discover-page__view-switch button { display: flex; min-height: 38px; align-items: center; gap: 6px; padding: 0 12px; border-radius: 10px; color: var(--ink-700); background: transparent; }
.discover-page__view-switch button.active { color: var(--forest-900); background: var(--white); box-shadow: var(--shadow-sm); }
.discover-controls { display: grid; grid-template-columns: minmax(260px, 1fr) auto auto; gap: 12px; margin-bottom: 18px; }
.discover-radius { display: flex; min-height: 52px; align-items: center; gap: 8px; padding: 0 12px 0 15px; border: 1px solid var(--sand-200); border-radius: 15px; color: var(--ink-700); background: var(--white); font-size: .82rem; font-weight: 750; }
.discover-radius select { border: 0; color: var(--forest-900); background: transparent; outline: 0; font-weight: 850; }
.discover-filter { border-radius: 15px; }
.discover-categories { margin-bottom: 22px; }
.discover-layout { display: grid; grid-template-columns: 375px 1fr; gap: 20px; }
.discover-results { height: 570px; padding: 17px 10px 17px 17px; overflow: hidden; border: 1px solid var(--sand-200); border-radius: 25px; background: var(--white); }
.discover-results__summary { display: flex; align-items: center; justify-content: space-between; padding: 2px 7px 14px 0; color: var(--ink-700); font-size: .76rem; font-weight: 750; }
.discover-results__summary button { display: flex; align-items: center; gap: 5px; padding: 0; color: var(--forest-700); background: transparent; cursor: pointer; font-size: .72rem; font-weight: 850; }
.discover-results__list { display: grid; gap: 10px; max-height: 512px; padding-right: 7px; overflow-y: auto; }
.map-result { display: grid; width: 100%; grid-template-columns: 80px 1fr auto; align-items: center; gap: 12px; padding: 9px; border: 1px solid transparent; border-radius: 17px; text-align: left; background: var(--cream-50); cursor: pointer; transition: background 160ms ease, border-color 160ms ease; }
.map-result:hover, .map-result--active { border-color: var(--mint-300); background: #edf6ef; }
.map-result > img { width: 80px; height: 75px; object-fit: cover; border-radius: 12px; }
.map-result__body { display: grid; min-width: 0; gap: 5px; }
.map-result__body small { overflow: hidden; color: var(--forest-700); font-size: .67rem; font-weight: 850; text-overflow: ellipsis; text-transform: uppercase; white-space: nowrap; }
.map-result__body strong { overflow: hidden; font-size: .86rem; text-overflow: ellipsis; white-space: nowrap; }
.map-result__body span { display: flex; align-items: center; gap: 4px; color: var(--ink-500); font-size: .67rem; }
.map-result > b { align-self: end; color: var(--forest-900); font-size: .84rem; }
.discover-map-wrap { position: relative; min-width: 0; }
.map-preview { position: absolute; z-index: 10; right: 20px; bottom: 20px; left: 20px; display: grid; grid-template-columns: 74px 1fr auto auto; align-items: center; gap: 13px; padding: 12px; border: 1px solid rgba(255,255,255,.6); border-radius: 20px; background: rgba(255,255,255,.94); box-shadow: var(--shadow-md); backdrop-filter: blur(12px); }
.map-preview img { width: 74px; height: 65px; object-fit: cover; border-radius: 12px; }
.map-preview > div:nth-child(2) { display: grid; min-width: 0; gap: 4px; }
.map-preview small { color: var(--forest-700); font-size: .65rem; font-weight: 850; text-transform: uppercase; }
.map-preview strong { overflow: hidden; font-size: .9rem; text-overflow: ellipsis; white-space: nowrap; }
.map-preview span { color: var(--ink-500); font-size: .69rem; }
.map-preview__price { display: grid; gap: 3px; text-align: right; }
.map-preview__price b { color: var(--forest-900); }
.map-preview__price s { color: var(--ink-500); font-size: .68rem; }
.map-preview .button { min-height: 43px; padding-inline: 16px; }
.discover-empty { display: grid; justify-items: center; padding: 70px 20px; text-align: center; }
.discover-empty > span { font-size: 2rem; }
.discover-empty h3 { margin: 14px 0 7px; }
.discover-empty p { color: var(--ink-500); font-size: .82rem; line-height: 1.5; }
@media (max-width: 940px) { .discover-layout { grid-template-columns: 310px 1fr; } .map-preview { grid-template-columns: 64px 1fr auto; } .map-preview .button { display: none; } }
@media (max-width: 720px) {
  .discover-page { padding-top: 28px; }
  .discover-page__heading { align-items: center; }
  .discover-page__view-switch { display: flex; }
  .discover-controls { grid-template-columns: 1fr auto; }
  .discover-radius { display: none; }
  .discover-filter { width: 52px; padding: 0; }
  .discover-filter :deep(span) { display: none; }
  .discover-filter { font-size: 0; }
  .discover-layout { display: block; }
  .discover-results { height: auto; min-height: 390px; padding: 13px; }
  .discover-results__list { max-height: none; }
  .discover-results--mobile-hidden, .discover-map-wrap--mobile-hidden { display: none; }
  .map-preview { right: 10px; bottom: 10px; left: 10px; grid-template-columns: 58px 1fr auto; }
  .map-preview img { width: 58px; height: 55px; }
  .map-preview__price s { display: none; }
}
@media (max-width: 430px) { .discover-page__heading { align-items: end; } .discover-page__view-switch button { width: 42px; justify-content: center; padding: 0; font-size: 0; } }
</style>
