<script setup lang="ts">
import { ArrowLeft, Check, ChevronRight, Clock3, Heart, Info, MapPin, Minus, PackageCheck, Plus, ShieldCheck, Star, Store, Truck } from 'lucide-vue-next'
import { findOffer, formatPrice, offers } from '~/data/marketplace'

const route = useRoute()
const { isAuthenticated } = useAuth()
const { isFavorite, toggleFavorite } = useMarketplace()
const quantity = ref(1)
const mode = ref<'pickup' | 'delivery'>('pickup')
const showConfirmation = ref(false)

const offer = computed(() => findOffer(String(route.params.id)))

if (!offer.value) {
  throw createError({ statusCode: 404, statusMessage: 'Предложение не найдено' })
}

const total = computed(() => (offer.value?.price || 0) * quantity.value + (mode.value === 'delivery' ? 149 : 0))
const related = computed(() => offers.filter(item => item.id !== offer.value?.id).slice(0, 3))

const reserve = async () => {
  if (!isAuthenticated.value) {
    await navigateTo(`/login?next=${encodeURIComponent(route.fullPath)}`)
    return
  }
  showConfirmation.value = true
}

useSeoMeta({
  title: computed(() => `${offer.value?.title} — ЕщёЕсть`),
  description: computed(() => offer.value?.description),
})
</script>

<template>
  <section v-if="offer" class="offer-page page-section page-section--compact">
    <div class="container">
      <NuxtLink class="offer-page__back" to="/"><ArrowLeft :size="18" /> Вернуться к предложениям</NuxtLink>

      <div class="offer-detail">
        <div class="offer-gallery">
          <img :src="offer.image" :alt="offer.title">
          <div class="offer-gallery__badges">
            <span v-for="tag in offer.tags" :key="tag" class="chip chip--accent">{{ tag }}</span>
          </div>
          <button class="icon-button offer-gallery__favorite" :class="{ active: isFavorite(offer.id) }" type="button" @click="toggleFavorite(offer.id)">
            <Heart :size="21" :fill="isFavorite(offer.id) ? 'currentColor' : 'none'" />
          </button>
        </div>

        <div class="offer-info">
          <div class="offer-info__merchant">
            <span class="offer-info__merchant-logo">{{ offer.merchant.charAt(0) }}</span>
            <span><small>Предложение от</small><strong>{{ offer.merchant }}</strong></span>
            <span class="offer-info__rating"><Star :size="15" fill="currentColor" /> {{ offer.rating }} <small>({{ offer.reviewCount }})</small></span>
          </div>
          <h1>{{ offer.title }}</h1>
          <p class="offer-info__description">{{ offer.description }}</p>

          <div class="offer-info__facts">
            <div><span><Clock3 :size="20" /></span><div><small>Получение</small><strong>{{ offer.pickupWindow }}</strong></div></div>
            <div><span><MapPin :size="20" /></span><div><small>Место</small><strong>{{ offer.address }} · {{ offer.distanceKm.toLocaleString('ru-RU') }} км</strong></div></div>
            <div><span><PackageCheck :size="20" /></span><div><small>Доступно</small><strong>{{ offer.available }} {{ offer.available === 1 ? 'пакет' : 'пакета' }}</strong></div></div>
          </div>

          <div class="offer-info__section">
            <h3>Что может быть внутри</h3>
            <p>{{ offer.contents }}</p>
            <span><Info :size="15" /> Точный состав — сюрприз и зависит от остатков дня.</span>
          </div>

          <div class="receive-mode">
            <button :class="{ active: mode === 'pickup' }" type="button" @click="mode = 'pickup'"><Store :size="18" /><span><strong>Самовывоз</strong><small>Бесплатно</small></span><Check v-if="mode === 'pickup'" :size="17" /></button>
            <button v-if="offer.delivery" :class="{ active: mode === 'delivery' }" type="button" @click="mode = 'delivery'"><Truck :size="18" /><span><strong>Доставка</strong><small>Демо · 149 ₽</small></span><Check v-if="mode === 'delivery'" :size="17" /></button>
          </div>

          <div class="offer-checkout">
            <div class="offer-checkout__price"><span><strong>{{ formatPrice(offer.price) }}</strong><s>{{ formatPrice(offer.originalPrice) }}</s></span><small>за один пакет</small></div>
            <div class="quantity-picker"><button type="button" :disabled="quantity <= 1" @click="quantity--"><Minus :size="17" /></button><strong>{{ quantity }}</strong><button type="button" :disabled="quantity >= offer.available" @click="quantity++"><Plus :size="17" /></button></div>
            <button class="button button--coral offer-checkout__button" type="button" @click="reserve">Забронировать · {{ formatPrice(total) }}</button>
          </div>
          <p class="offer-info__safe"><ShieldCheck :size="17" /> В пет-проекте бронирование и оплата демонстрационные.</p>
        </div>
      </div>

      <section class="offer-location card">
        <div class="offer-location__map"><span><MapPin :size="22" /></span></div>
        <div><p class="eyebrow">Где забрать</p><h2>{{ offer.merchant }}</h2><p>{{ offer.address }}, Москва</p><small>Откройте страницу получения после бронирования — там будет маршрут и короткий код.</small></div>
        <button class="button button--secondary" type="button">Показать маршрут <ChevronRight :size="18" /></button>
      </section>

      <section class="related-offers">
        <div class="section-heading"><div><p class="eyebrow">Ещё рядом</p><h2>Возможно, вам понравится</h2></div><NuxtLink class="text-link" to="/">Весь каталог <ChevronRight :size="17" /></NuxtLink></div>
        <div class="offer-grid"><OfferCard v-for="item in related" :key="item.id" :offer="item" /></div>
      </section>
    </div>

    <Transition name="modal">
      <div v-if="showConfirmation" class="confirmation" role="dialog" aria-modal="true" aria-labelledby="confirmation-title" @click.self="showConfirmation = false">
        <div class="confirmation__card">
          <span class="confirmation__icon"><Check :size="28" /></span>
          <p class="eyebrow">Демо-бронирование</p>
          <h2 id="confirmation-title">Пакет за вами!</h2>
          <p>Мы добавили демонстрационный заказ. На странице получения можно посмотреть код, маршрут и прототип доставки.</p>
          <div class="confirmation__summary"><img :src="offer.image" :alt="offer.title"><div><strong>{{ offer.merchant }}</strong><span>{{ mode === 'pickup' ? offer.pickupWindow : 'Доставка 20:30–21:10' }}</span></div><b>{{ formatPrice(total) }}</b></div>
          <NuxtLink class="button button--primary button--block" to="/delivery">Перейти к заказу</NuxtLink>
          <button class="button button--ghost button--block" type="button" @click="showConfirmation = false">Продолжить смотреть</button>
        </div>
      </div>
    </Transition>
  </section>
</template>

<style scoped>
.offer-page__back { display: inline-flex; align-items: center; gap: 8px; margin-bottom: 25px; color: var(--ink-700); font-size: .82rem; font-weight: 750; }
.offer-detail { display: grid; grid-template-columns: minmax(0, 1.02fr) minmax(400px, .98fr); gap: clamp(35px, 6vw, 75px); align-items: start; }
.offer-gallery { position: sticky; top: 98px; height: min(680px, calc(100vh - 125px)); min-height: 540px; overflow: hidden; border-radius: 30px; background: var(--cream-100); }
.offer-gallery > img { width: 100%; height: 100%; object-fit: cover; }
.offer-gallery__badges { position: absolute; top: 18px; left: 18px; display: flex; flex-wrap: wrap; gap: 8px; }
.offer-gallery__favorite { position: absolute; top: 17px; right: 17px; }
.offer-gallery__favorite.active { color: var(--coral-500); }
.offer-info { padding-top: 5px; }
.offer-info__merchant { display: flex; align-items: center; gap: 11px; margin-bottom: 24px; }
.offer-info__merchant-logo { display: grid; width: 43px; height: 43px; place-items: center; border-radius: 14px; color: var(--forest-950); background: var(--lime-300); font-weight: 950; }
.offer-info__merchant > span:nth-child(2) { display: grid; gap: 2px; }
.offer-info__merchant small { color: var(--ink-500); font-size: .66rem; }
.offer-info__merchant strong { font-size: .86rem; }
.offer-info__rating { display: inline-flex; align-items: center; gap: 4px; margin-left: auto; color: #aa7104; font-size: .82rem; font-weight: 850; }
.offer-info__rating small { color: var(--ink-500); }
.offer-info h1 { margin-bottom: 17px; font-size: clamp(2.3rem, 5vw, 4.3rem); }
.offer-info__description { color: var(--ink-700); font-size: .98rem; line-height: 1.7; }
.offer-info__facts { display: grid; gap: 14px; margin: 28px 0; padding: 20px 0; border-top: 1px solid var(--sand-200); border-bottom: 1px solid var(--sand-200); }
.offer-info__facts > div { display: flex; align-items: center; gap: 12px; }
.offer-info__facts > div > span { display: grid; width: 42px; height: 42px; flex: 0 0 auto; place-items: center; border-radius: 13px; color: var(--forest-900); background: var(--mint-200); }
.offer-info__facts > div > div { display: grid; gap: 3px; }
.offer-info__facts small { color: var(--ink-500); font-size: .66rem; }
.offer-info__facts strong { font-size: .82rem; }
.offer-info__section { margin: 27px 0; }
.offer-info__section h3 { margin-bottom: 9px; }
.offer-info__section p { color: var(--ink-700); font-size: .86rem; line-height: 1.65; }
.offer-info__section span { display: flex; align-items: center; gap: 6px; color: var(--ink-500); font-size: .69rem; }
.receive-mode { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; margin: 25px 0; }
.receive-mode button { display: grid; grid-template-columns: auto 1fr auto; align-items: center; gap: 10px; padding: 14px; border: 1px solid var(--sand-200); border-radius: 16px; color: var(--ink-700); text-align: left; background: var(--white); cursor: pointer; }
.receive-mode button.active { border-color: var(--forest-700); color: var(--forest-900); background: #eef7f0; box-shadow: 0 0 0 3px rgba(184,224,194,.3); }
.receive-mode button > span { display: grid; gap: 3px; }
.receive-mode small { color: var(--ink-500); font-size: .65rem; }
.offer-checkout { display: grid; grid-template-columns: 1fr auto; gap: 16px; align-items: center; padding: 20px; border-radius: 22px; background: var(--cream-100); }
.offer-checkout__price { display: grid; gap: 4px; }
.offer-checkout__price > span { display: flex; align-items: baseline; gap: 9px; }
.offer-checkout__price strong { color: var(--forest-900); font-size: 1.55rem; }
.offer-checkout__price s, .offer-checkout__price small { color: var(--ink-500); font-size: .7rem; }
.quantity-picker { display: flex; align-items: center; gap: 12px; padding: 5px; border-radius: 999px; background: var(--white); }
.quantity-picker button { display: grid; width: 34px; height: 34px; place-items: center; border-radius: 50%; color: var(--forest-900); background: var(--cream-100); cursor: pointer; }
.quantity-picker button:disabled { opacity: .35; cursor: not-allowed; }
.offer-checkout__button { grid-column: 1 / -1; width: 100%; }
.offer-info__safe { display: flex; align-items: center; justify-content: center; gap: 7px; margin: 13px 0 0; color: var(--ink-500); font-size: .68rem; }
.offer-location { display: grid; grid-template-columns: 260px 1fr auto; align-items: center; gap: 30px; margin: 80px 0; padding: 18px; }
.offer-location__map { position: relative; height: 185px; overflow: hidden; border-radius: 18px; background-color: #e8eee6; background-image: linear-gradient(34deg, transparent 47%, #fff 48%, #fff 52%, transparent 53%), linear-gradient(112deg, transparent 47%, rgba(255,255,255,.8) 48%, rgba(255,255,255,.8) 52%, transparent 53%); background-size: 95px 95px, 120px 120px; }
.offer-location__map span { position: absolute; top: 50%; left: 50%; display: grid; width: 50px; height: 50px; place-items: center; border: 4px solid var(--white); border-radius: 50%; color: var(--white); background: var(--coral-500); box-shadow: var(--shadow-md); transform: translate(-50%, -50%); }
.offer-location h2 { margin-bottom: 7px; font-size: 1.5rem; }
.offer-location > div:nth-child(2) > p:not(.eyebrow) { margin-bottom: 7px; color: var(--ink-700); }
.offer-location small { color: var(--ink-500); line-height: 1.45; }
.related-offers { padding-bottom: 30px; }
.confirmation { position: fixed; z-index: 200; inset: 0; display: grid; place-items: center; padding: 20px; background: rgba(16,43,36,.63); backdrop-filter: blur(8px); }
.confirmation__card { width: min(100%, 480px); padding: 32px; border-radius: 28px; background: var(--white); box-shadow: var(--shadow-lg); text-align: center; }
.confirmation__icon { display: grid; width: 62px; height: 62px; margin: 0 auto 20px; place-items: center; border-radius: 20px; color: var(--forest-950); background: var(--lime-300); }
.confirmation__card .eyebrow { justify-content: center; }
.confirmation__card h2 { margin-bottom: 10px; }
.confirmation__card > p:not(.eyebrow) { color: var(--ink-700); font-size: .85rem; line-height: 1.6; }
.confirmation__summary { display: grid; grid-template-columns: 60px 1fr auto; align-items: center; gap: 11px; margin: 22px 0; padding: 10px; border-radius: 16px; text-align: left; background: var(--cream-100); }
.confirmation__summary img { width: 60px; height: 56px; object-fit: cover; border-radius: 11px; }
.confirmation__summary div { display: grid; gap: 4px; }
.confirmation__summary span { color: var(--ink-500); font-size: .67rem; }
.confirmation__summary b { color: var(--forest-900); font-size: .85rem; }
.modal-enter-active, .modal-leave-active { transition: opacity .2s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
@media (max-width: 960px) { .offer-detail { grid-template-columns: 1fr 1fr; gap: 30px; } .offer-location { grid-template-columns: 210px 1fr; } .offer-location > .button { grid-column: 2; justify-self: start; } }
@media (max-width: 760px) { .offer-detail { grid-template-columns: 1fr; } .offer-gallery { position: relative; top: auto; height: 430px; min-height: 0; } .offer-info h1 { font-size: 2.5rem; } .offer-location { grid-template-columns: 1fr; margin: 50px 0; } .offer-location__map { height: 210px; } .offer-location > .button { grid-column: auto; justify-self: stretch; } }
@media (max-width: 430px) { .offer-gallery { height: 340px; border-radius: 22px; } .receive-mode { grid-template-columns: 1fr; } .offer-checkout { grid-template-columns: 1fr; } .quantity-picker { justify-self: start; } .confirmation__card { padding: 25px 18px; } }
</style>
