<script setup lang="ts">
import { Bike, Check, ChevronRight, Clock3, Copy, HelpCircle, MapPin, Navigation, PackageCheck, Phone, QrCode, Store, Truck } from 'lucide-vue-next'
import { demoOrder, findOffer, offers } from '~/data/marketplace'

definePageMeta({ middleware: 'auth' })

const mode = ref<'pickup' | 'delivery'>('pickup')
const copied = ref(false)
const activeOffer = findOffer(demoOrder.offerId)!

const copyCode = async () => {
  if (import.meta.client && navigator.clipboard) {
    await navigator.clipboard.writeText(demoOrder.code)
  }
  copied.value = true
  setTimeout(() => copied.value = false, 1400)
}

useSeoMeta({ title: 'Получение заказа — ЕщёЕсть' })
</script>

<template>
  <section class="delivery-page page-section page-section--compact">
    <div class="container">
      <div class="delivery-heading">
        <div>
          <p class="eyebrow">Мои заказы</p>
          <h1>Получение</h1>
          <p>Все детали активного заказа в одном месте.</p>
        </div>
        <div class="delivery-mode" role="group" aria-label="Способ получения">
          <button :class="{ active: mode === 'pickup' }" type="button" @click="mode = 'pickup'"><Store :size="18" /> Самовывоз</button>
          <button :class="{ active: mode === 'delivery' }" type="button" @click="mode = 'delivery'"><Bike :size="18" /> Доставка <small>демо</small></button>
        </div>
      </div>

      <div class="delivery-layout">
        <div class="delivery-main">
          <article class="order-status-card surface-dark">
            <div class="order-status-card__header">
              <span class="order-status-card__state"><i class="status-dot" /> {{ mode === 'pickup' ? 'Готовят к выдаче' : 'Собирают для курьера' }}</span>
              <span>Заказ #{{ demoOrder.id }}</span>
            </div>

            <div class="order-status-card__body">
              <img :src="activeOffer.image" :alt="activeOffer.title">
              <div>
                <small>{{ activeOffer.merchant }}</small>
                <h2>{{ activeOffer.title }}</h2>
                <p v-if="mode === 'pickup'"><Clock3 :size="17" /> Заберите сегодня, {{ demoOrder.eta }}</p>
                <p v-else><Truck :size="17" /> Доставим сегодня, 20:30–21:10</p>
              </div>
              <strong>{{ activeOffer.price }} ₽</strong>
            </div>

            <div class="order-progress">
              <div class="order-progress__item order-progress__item--done"><span><Check :size="16" /></span><small>Подтверждён</small></div>
              <i class="done" />
              <div class="order-progress__item order-progress__item--active"><span><PackageCheck :size="16" /></span><small>{{ mode === 'pickup' ? 'Готовят' : 'Собирают' }}</small></div>
              <i />
              <div class="order-progress__item"><span><Store v-if="mode === 'pickup'" :size="16" /><Bike v-else :size="16" /></span><small>{{ mode === 'pickup' ? 'Получен' : 'В пути' }}</small></div>
            </div>
          </article>

          <article v-if="mode === 'pickup'" class="pickup-code card">
            <div class="pickup-code__heading">
              <span><QrCode :size="24" /></span>
              <div><h3>Код получения</h3><p>Покажите сотруднику только когда получите пакет.</p></div>
            </div>
            <button class="pickup-code__value" type="button" @click="copyCode">
              <strong>{{ demoOrder.code }}</strong>
              <span><Copy :size="16" /> {{ copied ? 'Скопировано' : 'Скопировать' }}</span>
            </button>
            <p class="pickup-code__note">Не нажимайте «Получено» заранее — в настоящей версии подтверждение выполнит сотрудник места.</p>
          </article>

          <article v-else class="courier-card card">
            <div class="courier-card__person">
              <div class="avatar">М</div>
              <div><small>Ваш курьер</small><strong>Михаил · 4.9</strong><span>Будет у места через 8 минут</span></div>
              <button class="icon-button" type="button" aria-label="Позвонить курьеру"><Phone :size="18" /></button>
            </div>
            <div class="courier-map">
              <span class="courier-map__road one" /><span class="courier-map__road two" />
              <div class="courier-map__store"><Store :size="16" /></div>
              <div class="courier-map__bike"><Bike :size="17" /></div>
              <div class="courier-map__home"><MapPin :size="17" /></div>
              <span class="courier-map__eta">20–25 мин</span>
            </div>
            <p class="courier-card__demo">Доставка — только UX-прототип. Логистика и реальный курьерский backend пока не входят в пет-проект.</p>
          </article>

          <article class="pickup-details card">
            <div class="pickup-details__map">
              <span class="pickup-details__pin"><MapPin :size="21" /></span>
              <span class="pickup-details__you"><Navigation :size="15" fill="currentColor" /></span>
            </div>
            <div class="pickup-details__copy">
              <p class="eyebrow">{{ mode === 'pickup' ? 'Куда идти' : 'Адрес доставки' }}</p>
              <h3>{{ mode === 'pickup' ? activeOffer.merchant : 'Дом' }}</h3>
              <p>{{ mode === 'pickup' ? activeOffer.address : 'Покровский б-р, 5 · кв. 24' }}</p>
              <p v-if="mode === 'pickup'" class="pickup-details__hint">Вход со двора, скажите сотруднику, что вы за пакетом «ЕщёЕсть».</p>
              <button class="button button--secondary" type="button"><Navigation :size="17" /> {{ mode === 'pickup' ? 'Построить маршрут' : 'Изменить адрес' }}</button>
            </div>
          </article>
        </div>

        <aside class="delivery-aside">
          <article class="order-summary card">
            <h3>Детали заказа</h3>
            <div><span>1 × {{ activeOffer.title }}</span><strong>{{ activeOffer.price }} ₽</strong></div>
            <div><span>{{ mode === 'pickup' ? 'Самовывоз' : 'Доставка (демо)' }}</span><strong>{{ mode === 'pickup' ? '0 ₽' : '149 ₽' }}</strong></div>
            <div class="order-summary__total"><span>Итого</span><strong>{{ mode === 'pickup' ? activeOffer.price : activeOffer.price + 149 }} ₽</strong></div>
            <p>Оплачено демо-платежом · •• 2048</p>
          </article>

          <article class="need-help card">
            <span><HelpCircle :size="21" /></span>
            <div><h3>Нужна помощь?</h3><p>Если место закрыто или заказ не готов, напишите нам.</p><button type="button">Открыть поддержку <ChevronRight :size="16" /></button></div>
          </article>

          <article class="next-order">
            <p>Хотите спасти ещё?</p>
            <OfferCard :offer="offers[2]!" compact />
          </article>
        </aside>
      </div>
    </div>
  </section>
</template>

<style scoped>
.delivery-heading { display: flex; align-items: end; justify-content: space-between; gap: 25px; margin-bottom: 30px; }
.delivery-heading h1 { margin-bottom: 9px; font-size: clamp(2.3rem, 5vw, 4.2rem); }
.delivery-heading > div:first-child > p:last-child { margin: 0; color: var(--ink-700); }
.delivery-mode { display: flex; padding: 5px; border-radius: 16px; background: var(--cream-100); }
.delivery-mode button { position: relative; display: flex; min-height: 44px; align-items: center; gap: 7px; padding: 0 15px; border-radius: 12px; color: var(--ink-700); background: transparent; cursor: pointer; font-size: .82rem; font-weight: 800; }
.delivery-mode button.active { color: var(--forest-900); background: var(--white); box-shadow: var(--shadow-sm); }
.delivery-mode small { position: absolute; top: -13px; right: 6px; padding: 2px 5px; border-radius: 999px; color: var(--forest-950); background: var(--lime-300); font-size: .52rem; }
.delivery-layout { display: grid; grid-template-columns: minmax(0, 1.55fr) minmax(280px, .72fr); align-items: start; gap: 24px; }
.delivery-main, .delivery-aside { display: grid; gap: 20px; }
.order-status-card { padding: 24px; border-radius: 28px; background-image: radial-gradient(circle at 100% 0%, rgba(221,244,100,.18), transparent 35%); }
.order-status-card__header { display: flex; align-items: center; justify-content: space-between; gap: 16px; padding-bottom: 18px; border-bottom: 1px solid rgba(255,255,255,.1); color: rgba(255,255,255,.55); font-size: .72rem; }
.order-status-card__state { display: flex; align-items: center; gap: 9px; color: var(--mint-300); font-weight: 850; letter-spacing: .06em; text-transform: uppercase; }
.order-status-card__body { display: grid; grid-template-columns: 112px 1fr auto; align-items: center; gap: 18px; padding: 22px 0 30px; }
.order-status-card__body img { width: 112px; height: 98px; object-fit: cover; border-radius: 18px; }
.order-status-card__body small { color: var(--mint-300); font-size: .68rem; font-weight: 850; text-transform: uppercase; }
.order-status-card__body h2 { margin: 6px 0 10px; font-size: clamp(1.25rem, 2.4vw, 1.9rem); }
.order-status-card__body p { display: flex; align-items: center; gap: 7px; margin: 0; color: rgba(255,255,255,.68); font-size: .82rem; }
.order-status-card__body > strong { color: var(--lime-300); font-size: 1.25rem; }
.order-progress { display: grid; grid-template-columns: auto 1fr auto 1fr auto; align-items: start; }
.order-progress > i { height: 3px; margin-top: 18px; background: rgba(255,255,255,.18); }
.order-progress > i.done { background: var(--lime-300); }
.order-progress__item { display: grid; width: 72px; justify-items: center; gap: 7px; color: rgba(255,255,255,.48); text-align: center; }
.order-progress__item span { display: grid; width: 36px; height: 36px; place-items: center; border: 2px solid rgba(255,255,255,.25); border-radius: 50%; }
.order-progress__item small { font-size: .64rem; }
.order-progress__item--done, .order-progress__item--active { color: var(--white); }
.order-progress__item--done span { color: var(--forest-950); border-color: var(--lime-300); background: var(--lime-300); }
.order-progress__item--active span { border-color: var(--lime-300); box-shadow: 0 0 0 5px rgba(221,244,100,.1); }
.pickup-code { display: grid; grid-template-columns: 1fr auto; gap: 18px 30px; padding: 25px; }
.pickup-code__heading { display: flex; align-items: center; gap: 13px; }
.pickup-code__heading > span { display: grid; width: 50px; height: 50px; flex: 0 0 auto; place-items: center; border-radius: 16px; color: var(--forest-950); background: var(--lime-300); }
.pickup-code h3 { margin-bottom: 5px; }
.pickup-code p { margin: 0; color: var(--ink-700); font-size: .8rem; line-height: 1.5; }
.pickup-code__value { display: flex; min-width: 220px; align-items: center; justify-content: space-between; gap: 15px; padding: 12px 15px; border-radius: 15px; color: var(--forest-900); background: var(--cream-100); cursor: pointer; }
.pickup-code__value strong { font-size: 1.7rem; letter-spacing: .16em; }
.pickup-code__value span { display: flex; align-items: center; gap: 5px; color: var(--ink-700); font-size: .68rem; font-weight: 750; }
.pickup-code__note { grid-column: 1 / -1; padding-top: 15px; border-top: 1px solid var(--cream-100); }
.courier-card { padding: 22px; }
.courier-card__person { display: grid; grid-template-columns: auto 1fr auto; align-items: center; gap: 12px; margin-bottom: 18px; }
.courier-card__person > div:nth-child(2) { display: grid; gap: 3px; }
.courier-card__person small { color: var(--ink-500); font-size: .65rem; }
.courier-card__person span { color: var(--ink-700); font-size: .72rem; }
.courier-map { position: relative; height: 210px; overflow: hidden; border-radius: 20px; background: #e8eee6; }
.courier-map__road { position: absolute; top: 50%; left: -10%; width: 120%; height: 6px; background: var(--white); transform: rotate(-8deg); }
.courier-map__road.two { top: -20%; left: 48%; width: 100%; transform: rotate(84deg); }
.courier-map__store, .courier-map__bike, .courier-map__home { position: absolute; z-index: 2; display: grid; width: 38px; height: 38px; place-items: center; border: 3px solid var(--white); border-radius: 50%; box-shadow: var(--shadow-sm); }
.courier-map__store { top: 27%; left: 17%; color: var(--white); background: var(--forest-900); }
.courier-map__bike { top: 46%; left: 49%; color: var(--forest-950); background: var(--lime-300); }
.courier-map__home { right: 17%; bottom: 19%; color: var(--white); background: var(--coral-500); }
.courier-map__eta { position: absolute; z-index: 2; top: 21px; left: 50%; padding: 7px 10px; border-radius: 999px; color: var(--forest-900); background: var(--white); box-shadow: var(--shadow-sm); font-size: .7rem; font-weight: 850; transform: translateX(-50%); }
.courier-card__demo { margin: 14px 0 0; color: var(--ink-500); font-size: .69rem; line-height: 1.5; }
.pickup-details { display: grid; grid-template-columns: .85fr 1.15fr; overflow: hidden; }
.pickup-details__map { position: relative; min-height: 270px; background-color: #e8eee6; background-image: linear-gradient(35deg, transparent 47%, #fff 48%, #fff 52%, transparent 53%), linear-gradient(112deg, transparent 47%, rgba(255,255,255,.8) 48%, rgba(255,255,255,.8) 52%, transparent 53%); background-size: 110px 110px, 150px 150px; }
.pickup-details__pin, .pickup-details__you { position: absolute; display: grid; place-items: center; border: 4px solid var(--white); border-radius: 50%; box-shadow: var(--shadow-sm); }
.pickup-details__pin { top: 28%; left: 30%; width: 48px; height: 48px; color: var(--white); background: var(--coral-500); }
.pickup-details__you { right: 20%; bottom: 22%; width: 38px; height: 38px; color: var(--white); background: #3389d8; }
.pickup-details__copy { align-self: center; padding: 28px; }
.pickup-details__copy h3 { margin-bottom: 7px; font-size: 1.3rem; }
.pickup-details__copy > p:not(.eyebrow) { margin-bottom: 8px; color: var(--ink-700); font-size: .82rem; line-height: 1.5; }
.pickup-details__hint { padding: 10px 12px; border-radius: 12px; background: var(--cream-100); }
.pickup-details__copy .button { min-height: 44px; margin-top: 10px; }
.order-summary { padding: 22px; }
.order-summary h3 { margin-bottom: 22px; }
.order-summary > div { display: flex; justify-content: space-between; gap: 18px; margin-bottom: 14px; color: var(--ink-700); font-size: .79rem; }
.order-summary > div strong { flex: 0 0 auto; color: var(--ink-950); }
.order-summary__total { margin-top: 19px; padding-top: 18px; border-top: 1px solid var(--sand-200); font-size: .95rem !important; font-weight: 850; }
.order-summary p { margin: 0; padding-top: 14px; border-top: 1px solid var(--cream-100); color: var(--ink-500); font-size: .67rem; }
.need-help { display: flex; align-items: flex-start; gap: 12px; padding: 20px; }
.need-help > span { display: grid; width: 40px; height: 40px; flex: 0 0 auto; place-items: center; border-radius: 13px; color: var(--forest-900); background: var(--mint-200); }
.need-help h3 { margin-bottom: 6px; }
.need-help p { color: var(--ink-700); font-size: .75rem; line-height: 1.5; }
.need-help button { display: flex; align-items: center; gap: 4px; padding: 0; color: var(--forest-700); background: transparent; cursor: pointer; font-size: .73rem; font-weight: 850; }
.next-order > p { margin: 5px 0 12px; color: var(--ink-700); font-size: .75rem; font-weight: 750; }
@media (max-width: 940px) { .delivery-layout { grid-template-columns: 1fr; } .delivery-aside { grid-template-columns: repeat(2, 1fr); } .next-order { display: none; } }
@media (max-width: 680px) { .delivery-heading { display: block; } .delivery-mode { margin-top: 22px; } .delivery-mode button { flex: 1; justify-content: center; } .order-status-card { padding: 20px 16px; } .order-status-card__body { grid-template-columns: 82px 1fr; } .order-status-card__body img { width: 82px; height: 78px; } .order-status-card__body > strong { grid-column: 2; } .pickup-code { grid-template-columns: 1fr; } .pickup-code__value { min-width: 0; } .pickup-code__note { grid-column: auto; } .pickup-details { grid-template-columns: 1fr; } .pickup-details__map { min-height: 200px; } .delivery-aside { grid-template-columns: 1fr; } }
</style>
