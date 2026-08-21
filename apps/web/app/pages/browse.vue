<script setup lang="ts">
import { ArrowDownUp, Search, SlidersHorizontal, Truck, X } from 'lucide-vue-next'
import { categories, offers } from '~/data/marketplace'

const route = useRoute()
const query = ref('')
const selectedCategory = ref(typeof route.query.category === 'string' ? route.query.category : 'all')
const sort = ref<'distance' | 'price' | 'rating'>('distance')
const deliveryOnly = ref(false)
const showFilters = ref(false)

const filtered = computed(() => {
  const search = query.value.trim().toLocaleLowerCase('ru-RU')
  const result = offers.filter((offer) => {
    const category = selectedCategory.value === 'all' || offer.category === selectedCategory.value
    const text = !search || [offer.title, offer.merchant, offer.district, ...offer.tags].some(value => value.toLocaleLowerCase('ru-RU').includes(search))
    const delivery = !deliveryOnly.value || offer.delivery
    return category && text && delivery
  })

  return [...result].sort((a, b) => {
    if (sort.value === 'price') return a.price - b.price
    if (sort.value === 'rating') return b.rating - a.rating
    return a.distanceKm - b.distanceKm
  })
})

const clearFilters = () => {
  query.value = ''
  selectedCategory.value = 'all'
  deliveryOnly.value = false
  sort.value = 'distance'
}

useSeoMeta({ title: 'Каталог предложений — ЕщёЕсть' })
</script>

<template>
  <section class="browse-page page-section page-section--compact">
    <div class="container">
      <div class="browse-heading">
        <div>
          <p class="eyebrow">Каталог</p>
          <h1>Выберите что-нибудь хорошее</h1>
          <p>Актуальные пакеты из мест рядом с вами.</p>
        </div>
        <div class="browse-heading__art" aria-hidden="true"><span>🥐</span><span>🥬</span><span>🍣</span></div>
      </div>

      <div class="browse-toolbar card">
        <div class="input-wrap browse-search">
          <Search :size="19" />
          <input v-model="query" class="input input--icon" type="search" placeholder="Кафе, блюдо или район...">
        </div>
        <label class="browse-sort">
          <ArrowDownUp :size="17" />
          <select v-model="sort" aria-label="Сортировка">
            <option value="distance">Сначала ближе</option>
            <option value="price">Сначала дешевле</option>
            <option value="rating">По рейтингу</option>
          </select>
        </label>
        <button class="button button--secondary browse-filter-button" :class="{ active: showFilters }" type="button" @click="showFilters = !showFilters">
          <SlidersHorizontal :size="18" /> Фильтры
        </button>
      </div>

      <div v-if="showFilters" class="browse-filter-panel card">
        <label class="toggle-row">
          <span><Truck :size="19" /><span><strong>Есть доставка</strong><small>Показывать места с демо-доставкой</small></span></span>
          <input v-model="deliveryOnly" type="checkbox">
          <i />
        </label>
        <button class="text-link" type="button" @click="clearFilters"><X :size="16" /> Сбросить всё</button>
      </div>

      <CategoryScroller v-model="selectedCategory" />

      <div class="browse-summary">
        <p><strong>{{ filtered.length }}</strong> {{ filtered.length === 1 ? 'предложение' : 'предложений' }}</p>
        <div v-if="deliveryOnly" class="chip chip--accent"><Truck :size="15" /> С доставкой</div>
      </div>

      <div v-if="filtered.length" class="offer-grid browse-grid">
        <OfferCard v-for="offer in filtered" :key="offer.id" :offer="offer" />
      </div>

      <div v-else class="browse-empty card">
        <span>🔎</span>
        <h2>Ничего не нашли</h2>
        <p>Попробуйте другую категорию или очистите фильтры.</p>
        <button class="button button--primary" type="button" @click="clearFilters">Сбросить фильтры</button>
      </div>
    </div>
  </section>
</template>

<style scoped>
.browse-heading { display: flex; min-height: 210px; align-items: center; justify-content: space-between; gap: 35px; margin-bottom: 28px; padding: 35px 42px; overflow: hidden; border-radius: 30px; color: var(--white); background: var(--forest-900); background-image: radial-gradient(circle at 80% 20%, rgba(221,244,100,.25), transparent 35%); }
.browse-heading h1 { max-width: 750px; margin-bottom: 12px; font-size: clamp(2.25rem, 5vw, 4.3rem); }
.browse-heading p:last-child { margin: 0; color: rgba(255,255,255,.65); }
.browse-heading .eyebrow { color: var(--mint-300); }
.browse-heading__art { position: relative; display: grid; width: 190px; height: 140px; flex: 0 0 auto; place-items: center; border: 1px solid rgba(255,255,255,.12); border-radius: 50%; background: rgba(255,255,255,.07); }
.browse-heading__art span { position: absolute; display: grid; width: 62px; height: 62px; place-items: center; border-radius: 20px; background: var(--white); box-shadow: var(--shadow-md); font-size: 1.7rem; transform: rotate(-8deg); }
.browse-heading__art span:nth-child(1) { top: -5px; left: 15px; }
.browse-heading__art span:nth-child(2) { right: 8px; bottom: -4px; transform: rotate(9deg); }
.browse-heading__art span:nth-child(3) { right: 6px; top: -16px; width: 52px; height: 52px; font-size: 1.35rem; transform: rotate(7deg); }
.browse-toolbar { display: grid; grid-template-columns: 1fr auto auto; gap: 12px; margin-bottom: 16px; padding: 10px; }
.browse-search .input { border-color: transparent; background: var(--cream-50); }
.browse-sort { display: flex; min-height: 52px; align-items: center; gap: 7px; padding: 0 14px; border: 1px solid var(--sand-200); border-radius: 15px; color: var(--ink-700); }
.browse-sort select { border: 0; color: var(--forest-900); background: transparent; outline: 0; font-weight: 800; }
.browse-filter-button { border-radius: 15px; }
.browse-filter-button.active { color: var(--white); background: var(--forest-900); box-shadow: none; }
.browse-filter-panel { display: flex; align-items: center; justify-content: space-between; gap: 20px; margin-bottom: 16px; padding: 16px 18px; }
.toggle-row { display: flex; align-items: center; gap: 14px; cursor: pointer; }
.toggle-row > span { display: flex; align-items: center; gap: 11px; }
.toggle-row > span > svg { color: var(--coral-500); }
.toggle-row > span span { display: grid; gap: 2px; }
.toggle-row small { color: var(--ink-500); font-size: .7rem; }
.toggle-row input { position: absolute; opacity: 0; }
.toggle-row i { position: relative; width: 46px; height: 26px; border-radius: 999px; background: var(--sand-200); transition: background 160ms ease; }
.toggle-row i::after { position: absolute; top: 4px; left: 4px; width: 18px; height: 18px; border-radius: 50%; content: ""; background: var(--white); box-shadow: 0 2px 5px rgba(0,0,0,.12); transition: transform 160ms ease; }
.toggle-row input:checked + i { background: var(--forest-700); }
.toggle-row input:checked + i::after { transform: translateX(20px); }
.browse-filter-panel .text-link { border: 0; background: transparent; cursor: pointer; }
.browse-summary { display: flex; min-height: 65px; align-items: center; justify-content: space-between; gap: 14px; }
.browse-summary p { margin: 0; color: var(--ink-700); font-size: .87rem; }
.browse-summary strong { color: var(--ink-950); }
.browse-grid { margin-bottom: 35px; }
.browse-empty { display: grid; justify-items: center; padding: 70px 24px; text-align: center; }
.browse-empty > span { font-size: 2.5rem; }
.browse-empty h2 { margin: 14px 0 8px; }
.browse-empty p { margin-bottom: 22px; color: var(--ink-700); }
@media (max-width: 760px) { .browse-heading { min-height: 180px; padding: 28px 24px; } .browse-heading__art { display: none; } .browse-toolbar { grid-template-columns: 1fr auto; } .browse-sort { grid-row: 2; grid-column: 1 / -1; } .browse-filter-button { width: 52px; padding: 0; font-size: 0; } .browse-filter-panel { align-items: flex-start; } .toggle-row { align-items: flex-start; } .toggle-row > span span small { display: none; } }
</style>
