<script setup lang="ts">
import { ArrowRight, Check, Clock3, Leaf, MapPin, PackageCheck, Search, ShieldCheck, Sparkles } from 'lucide-vue-next'
import { categories, demoOrder, findOffer, offers } from '~/data/marketplace'

const { user, isAuthenticated } = useAuth()
const selectedCategory = ref('all')
const activeOrderOffer = findOffer(demoOrder.offerId)

const recommended = computed(() => {
  const matching = selectedCategory.value === 'all'
    ? offers
    : offers.filter((offer) => offer.category === selectedCategory.value)
  return matching.slice(0, 3)
})

useSeoMeta({
  title: computed(() => isAuthenticated.value
    ? 'Главная — ЕщёЕсть'
    : 'ЕщёЕсть — хорошая еда рядом со скидкой'),
  description: 'Находите свежие сюрприз-пакеты из кафе, пекарен и магазинов рядом и забирайте их со скидкой.',
})
</script>

<template>
  <div>
    <template v-if="!isAuthenticated">
      <section class="guest-hero">
        <div class="container guest-hero__grid">
          <div class="guest-hero__copy">
            <p class="eyebrow">Вкусно. Выгодно. Бережно.</p>
            <h1>Хорошая еда.<br><span>Ещё есть.</span></h1>
            <p class="guest-hero__lead">
              Забирайте сюрприз-пакеты из любимых кафе, пекарен и магазинов со скидкой до 70%.
            </p>
            <div class="guest-hero__actions">
              <NuxtLink class="button button--accent" to="/discover">
                Найти еду рядом <ArrowRight :size="19" />
              </NuxtLink>
              <NuxtLink class="button button--secondary" to="/register">Создать аккаунт</NuxtLink>
            </div>
            <div class="guest-hero__trust">
              <span><Check :size="15" /> Без подписки</span>
              <span><Check :size="15" /> Забираете сами</span>
              <span><Check :size="15" /> Свежие предложения каждый день</span>
            </div>
          </div>

          <div class="guest-hero__visual">
            <div class="guest-hero__image-wrap">
              <img src="/images/bakery-rescue.png" alt="Свежая выпечка в бумажном пакете">
              <div class="floating-offer">
                <div class="floating-offer__icon">🥐</div>
                <div>
                  <strong>Пакет выпечки</strong>
                  <span>Забрать с 20:00</span>
                </div>
                <b>299 ₽</b>
              </div>
              <div class="floating-distance"><MapPin :size="15" /> 400 м от вас</div>
            </div>
            <div class="guest-hero__shape" aria-hidden="true" />
          </div>
        </div>
      </section>

      <section class="social-proof">
        <div class="container social-proof__grid">
          <div><strong>40–70%</strong><span>обычная скидка</span></div>
          <div><strong>15 минут</strong><span>на удобный самовывоз</span></div>
          <div><strong>1 пакет</strong><span>маленький полезный выбор</span></div>
          <p>Демо-данные пет-проекта — чтобы проверить продуктовую идею и интерфейс.</p>
        </div>
      </section>

      <section class="page-section how-it-works">
        <div class="container">
          <div class="section-heading">
            <div>
              <p class="eyebrow">Три простых шага</p>
              <h2>Забрать проще, чем выбрать ужин</h2>
            </div>
            <p class="how-it-works__intro">Места собирают хорошие остатки дня в сюрприз-пакеты — вы выбираете время и забираете.</p>
          </div>

          <div class="steps-grid">
            <article>
              <span>01</span>
              <div class="steps-grid__icon"><Search :size="24" /></div>
              <h3>Найдите рядом</h3>
              <p>Смотрите предложения на карте или выбирайте по категории и времени.</p>
            </article>
            <article>
              <span>02</span>
              <div class="steps-grid__icon"><PackageCheck :size="24" /></div>
              <h3>Забронируйте</h3>
              <p>Состав будет сюрпризом, зато цена, окно получения и место известны заранее.</p>
            </article>
            <article>
              <span>03</span>
              <div class="steps-grid__icon"><Sparkles :size="24" /></div>
              <h3>Заберите и наслаждайтесь</h3>
              <p>Покажите короткий код сотруднику в указанное время — и пакет ваш.</p>
            </article>
          </div>
        </div>
      </section>

      <section class="guest-offers page-section">
        <div class="container">
          <div class="section-heading">
            <div>
              <p class="eyebrow">Сегодня в Москве</p>
              <h2>Посмотрите, что может быть рядом</h2>
            </div>
            <NuxtLink class="text-link" to="/browse">Весь каталог <ArrowRight :size="17" /></NuxtLink>
          </div>
          <div class="offer-grid">
            <OfferCard v-for="offer in offers.slice(0, 3)" :key="offer.id" :offer="offer" />
          </div>
        </div>
      </section>

      <section class="guest-values page-section">
        <div class="container guest-values__grid">
          <div class="guest-values__visual">
            <img src="/images/grocery-rescue.png" alt="Свежие овощи, фрукты и хлеб в сумке">
            <div><Leaf :size="21" /> Хорошее получает ещё один шанс</div>
          </div>
          <div class="guest-values__copy">
            <p class="eyebrow">Зачем это нужно</p>
            <h2>Выгода для вас. Новый доход для мест рядом.</h2>
            <p>Небольшая привычка помогает кафе и магазинам продавать приготовленное, а вам — открывать новые места без лишних трат.</p>
            <ul>
              <li><ShieldCheck :size="19" /><span><strong>Понятные условия</strong>Цена, время и адрес известны до бронирования.</span></li>
              <li><Clock3 :size="19" /><span><strong>Только актуальные окна</strong>В каталоге показываются предложения, которые ещё можно забрать.</span></li>
              <li><Leaf :size="19" /><span><strong>Измеримый результат</strong>В профиле видно спасённые пакеты и вашу экономию.</span></li>
            </ul>
            <NuxtLink class="button button--primary" to="/register">Начать бесплатно</NuxtLink>
          </div>
        </div>
      </section>

      <footer class="guest-footer">
        <div class="container guest-footer__inner">
          <AppLogo inverse />
          <p>Пет-проект food rescue marketplace · 2026</p>
          <div><NuxtLink to="/browse">Каталог</NuxtLink><NuxtLink to="/login">Войти</NuxtLink></div>
        </div>
      </footer>
    </template>

    <template v-else>
      <section class="member-hero">
        <div class="container member-hero__grid">
          <div>
            <p class="member-hero__date">Сегодня в Москве</p>
            <h1>Привет, {{ user?.name?.split(' ')[0] }}!</h1>
            <p>Что хорошего спасём сегодня?</p>
          </div>
          <NuxtLink class="location-card" to="/discover">
            <span class="location-card__icon"><MapPin :size="20" /></span>
            <span><small>Ищем рядом с</small><strong>Чистыми прудами</strong></span>
            <ArrowRight :size="18" />
          </NuxtLink>
        </div>
      </section>

      <section class="container member-content">
        <div class="member-search">
          <Search :size="20" />
          <input aria-label="Поиск" placeholder="Найти кафе, магазин или блюдо..." @focus="navigateTo('/browse')">
          <NuxtLink class="button button--primary" to="/browse">Найти</NuxtLink>
        </div>

        <CategoryScroller v-model="selectedCategory" />

        <div v-if="activeOrderOffer" class="active-order">
          <div class="active-order__visual">
            <img :src="activeOrderOffer.image" :alt="activeOrderOffer.title">
          </div>
          <div class="active-order__copy">
            <span class="active-order__status"><i class="status-dot" /> Готовят ваш пакет</span>
            <h3>{{ activeOrderOffer.merchant }}</h3>
            <p><Clock3 :size="16" /> Заберите {{ demoOrder.eta }}</p>
          </div>
          <div class="active-order__code"><small>Код получения</small><strong>{{ demoOrder.code }}</strong></div>
          <NuxtLink class="button button--accent" to="/delivery">Открыть заказ</NuxtLink>
        </div>

        <div class="section-heading member-offers-heading">
          <div>
            <p class="eyebrow">Подобрано для вас</p>
            <h2>{{ selectedCategory === 'all' ? 'Успейте забрать сегодня' : categories.find(c => c.id === selectedCategory)?.label }}</h2>
          </div>
          <NuxtLink class="text-link" :to="`/browse?category=${selectedCategory}`">Показать всё <ArrowRight :size="17" /></NuxtLink>
        </div>

        <div class="offer-grid">
          <OfferCard v-for="offer in recommended" :key="offer.id" :offer="offer" />
        </div>

        <ImpactCard class="member-impact" />
      </section>
    </template>
  </div>
</template>

<style scoped>
.guest-hero {
  position: relative;
  overflow: hidden;
  padding: 68px 0 80px;
  background: linear-gradient(180deg, var(--cream-50), #f1efe6);
}

.guest-hero::before {
  position: absolute;
  top: -180px;
  right: -130px;
  width: 520px;
  height: 520px;
  border-radius: 50%;
  content: "";
  background: var(--mint-200);
  filter: blur(1px);
  opacity: .56;
}

.guest-hero__grid {
  position: relative;
  display: grid;
  grid-template-columns: 1.02fr .98fr;
  align-items: center;
  gap: clamp(40px, 7vw, 100px);
}

.guest-hero__copy {
  position: relative;
  z-index: 2;
}

.guest-hero h1 span {
  color: var(--coral-500);
}

.guest-hero__lead {
  max-width: 580px;
  margin-bottom: 28px;
  color: var(--ink-700);
  font-size: clamp(1.04rem, 2vw, 1.25rem);
  line-height: 1.65;
}

.guest-hero__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.guest-hero__trust {
  display: flex;
  flex-wrap: wrap;
  gap: 12px 20px;
  margin-top: 26px;
  color: var(--ink-700);
  font-size: .78rem;
  font-weight: 720;
}

.guest-hero__trust span {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.guest-hero__trust svg {
  color: var(--success);
}

.guest-hero__visual {
  position: relative;
  min-height: 520px;
}

.guest-hero__image-wrap {
  position: absolute;
  z-index: 2;
  inset: 0 0 0 8%;
  overflow: hidden;
  border: 9px solid rgba(255,255,255,.82);
  border-radius: 48% 48% 28px 28px;
  box-shadow: var(--shadow-lg);
}

.guest-hero__image-wrap > img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.guest-hero__shape {
  position: absolute;
  right: -34px;
  bottom: -35px;
  width: 190px;
  height: 190px;
  border: 25px solid var(--lime-300);
  border-radius: 50%;
}

.floating-offer {
  position: absolute;
  right: 18px;
  bottom: 20px;
  left: 18px;
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  gap: 12px;
  padding: 14px;
  border: 1px solid rgba(255,255,255,.58);
  border-radius: 18px;
  background: rgba(255,255,255,.91);
  box-shadow: var(--shadow-md);
  backdrop-filter: blur(14px);
}

.floating-offer__icon {
  display: grid;
  width: 45px;
  height: 45px;
  place-items: center;
  border-radius: 14px;
  background: var(--cream-100);
  font-size: 1.35rem;
}

.floating-offer strong,
.floating-offer span { display: block; }
.floating-offer span { margin-top: 4px; color: var(--ink-700); font-size: .72rem; }
.floating-offer b { color: var(--forest-900); font-size: 1.1rem; }

.floating-distance {
  position: absolute;
  top: 28px;
  left: 50%;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 9px 12px;
  border-radius: 999px;
  color: var(--forest-900);
  background: rgba(255,255,255,.9);
  font-size: .76rem;
  font-weight: 850;
  backdrop-filter: blur(12px);
  transform: translateX(-50%);
}

.social-proof {
  color: var(--white);
  background: var(--forest-900);
}

.social-proof__grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr) 1.7fr;
  align-items: center;
  gap: 28px;
  min-height: 118px;
}

.social-proof__grid > div { display: grid; gap: 4px; }
.social-proof__grid strong { color: var(--lime-300); font-size: 1.35rem; }
.social-proof__grid span { color: rgba(255,255,255,.64); font-size: .76rem; }
.social-proof__grid p { margin: 0; color: rgba(255,255,255,.58); font-size: .76rem; line-height: 1.55; }

.how-it-works__intro {
  max-width: 440px;
  margin-bottom: 4px;
  color: var(--ink-700);
  line-height: 1.65;
}

.steps-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 18px;
}

.steps-grid article {
  position: relative;
  min-height: 250px;
  padding: 28px;
  overflow: hidden;
  border: 1px solid var(--sand-200);
  border-radius: 24px;
  background: var(--white);
}

.steps-grid article > span {
  position: absolute;
  top: 6px;
  right: 16px;
  color: var(--cream-100);
  font-size: 4.6rem;
  font-weight: 950;
  letter-spacing: -.08em;
}

.steps-grid__icon {
  position: relative;
  display: grid;
  width: 50px;
  height: 50px;
  margin-bottom: 35px;
  place-items: center;
  border-radius: 16px;
  color: var(--forest-950);
  background: var(--lime-300);
}

.steps-grid p { margin: 0; color: var(--ink-700); font-size: .9rem; line-height: 1.65; }

.guest-offers { background: var(--cream-100); }

.guest-values__grid {
  display: grid;
  grid-template-columns: .98fr 1.02fr;
  align-items: center;
  gap: clamp(40px, 8vw, 110px);
}

.guest-values__visual {
  position: relative;
  height: 500px;
}

.guest-values__visual img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 30px 30px 50% 30px;
}

.guest-values__visual div {
  position: absolute;
  right: -22px;
  bottom: 28px;
  display: inline-flex;
  align-items: center;
  gap: 9px;
  padding: 14px 18px;
  border-radius: 999px;
  color: var(--forest-950);
  background: var(--lime-300);
  box-shadow: var(--shadow-md);
  font-size: .82rem;
  font-weight: 900;
}

.guest-values__copy > p:not(.eyebrow) { color: var(--ink-700); line-height: 1.7; }
.guest-values__copy ul { display: grid; gap: 18px; margin: 28px 0; padding: 0; list-style: none; }
.guest-values__copy li { display: flex; align-items: flex-start; gap: 12px; }
.guest-values__copy li > svg { flex: 0 0 auto; margin-top: 2px; color: var(--coral-500); }
.guest-values__copy li span { display: grid; gap: 3px; color: var(--ink-700); font-size: .88rem; line-height: 1.5; }
.guest-values__copy li strong { color: var(--ink-950); }

.guest-footer { padding: 38px 0; color: var(--white); background: var(--forest-950); }
.guest-footer__inner { display: flex; align-items: center; justify-content: space-between; gap: 20px; }
.guest-footer p { margin: 0; color: rgba(255,255,255,.52); font-size: .77rem; }
.guest-footer__inner > div { display: flex; gap: 22px; color: rgba(255,255,255,.74); font-size: .82rem; font-weight: 750; }

.member-hero {
  padding: 40px 0 84px;
  color: var(--white);
  background: var(--forest-900);
  background-image: radial-gradient(circle at 80% -40%, rgba(221,244,100,.28), transparent 38%);
}

.member-hero__grid { display: flex; align-items: center; justify-content: space-between; gap: 30px; }
.member-hero__date { margin-bottom: 8px; color: var(--mint-300); font-size: .8rem; font-weight: 800; }
.member-hero h1 { margin-bottom: 9px; font-size: clamp(2.2rem, 5vw, 4.2rem); }
.member-hero h1 + p { margin: 0; color: rgba(255,255,255,.68); font-size: 1.05rem; }

.location-card {
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  gap: 12px;
  min-width: 310px;
  padding: 14px;
  border: 1px solid rgba(255,255,255,.12);
  border-radius: 18px;
  background: rgba(255,255,255,.08);
}
.location-card__icon { display: grid; width: 44px; height: 44px; place-items: center; border-radius: 14px; color: var(--forest-950); background: var(--lime-300); }
.location-card small, .location-card strong { display: block; }
.location-card small { margin-bottom: 3px; color: rgba(255,255,255,.55); font-size: .7rem; }
.location-card strong { font-size: .9rem; }

.member-content { position: relative; margin-top: -40px; padding-bottom: 90px; }
.member-search {
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
  padding: 9px 9px 9px 19px;
  border-radius: 20px;
  background: var(--white);
  box-shadow: var(--shadow-md);
}
.member-search > svg { color: var(--ink-500); }
.member-search input { height: 46px; border: 0; outline: 0; color: var(--ink-950); background: transparent; }
.member-search .button { min-height: 46px; }

.active-order {
  display: grid;
  grid-template-columns: 92px 1fr auto auto;
  align-items: center;
  gap: 18px;
  margin: 32px 0 52px;
  padding: 14px 16px 14px 14px;
  border: 1px solid var(--sand-200);
  border-radius: 22px;
  background: var(--white);
  box-shadow: var(--shadow-sm);
}
.active-order__visual { width: 92px; height: 82px; overflow: hidden; border-radius: 15px; }
.active-order__visual img { width: 100%; height: 100%; object-fit: cover; }
.active-order__status { display: flex; align-items: center; gap: 9px; margin-bottom: 8px; color: var(--success); font-size: .72rem; font-weight: 850; text-transform: uppercase; }
.active-order h3 { margin-bottom: 7px; }
.active-order__copy p { display: flex; align-items: center; gap: 6px; margin: 0; color: var(--ink-700); font-size: .8rem; }
.active-order__code { display: grid; min-width: 105px; justify-items: center; gap: 4px; padding: 10px 16px; border-right: 1px solid var(--sand-200); border-left: 1px solid var(--sand-200); }
.active-order__code small { color: var(--ink-500); font-size: .67rem; }
.active-order__code strong { color: var(--forest-900); font-size: 1.5rem; letter-spacing: .13em; }

.member-offers-heading { margin-top: 44px; }
.member-impact { margin-top: 62px; }

@media (max-width: 960px) {
  .guest-hero__grid, .guest-values__grid { grid-template-columns: 1fr 1fr; gap: 38px; }
  .guest-hero__visual { min-height: 430px; }
  .social-proof__grid { grid-template-columns: repeat(3, 1fr); }
  .social-proof__grid p { display: none; }
  .active-order { grid-template-columns: 82px 1fr auto; }
  .active-order__code { display: none; }
}

@media (max-width: 720px) {
  .guest-hero { padding: 42px 0 58px; }
  .guest-hero__grid, .guest-values__grid { grid-template-columns: 1fr; }
  .guest-hero__visual { min-height: 400px; }
  .guest-hero__image-wrap { inset: 0; border-radius: 28px 28px 42% 28px; }
  .floating-distance { top: 22px; left: 22px; transform: none; }
  .guest-hero__shape { right: -20px; }
  .social-proof__grid { grid-template-columns: repeat(3, 1fr); min-height: 96px; gap: 10px; }
  .social-proof__grid strong { font-size: 1rem; }
  .social-proof__grid span { font-size: .6rem; }
  .steps-grid { grid-template-columns: 1fr; }
  .steps-grid article { min-height: 215px; }
  .guest-values__visual { height: 360px; }
  .guest-values__visual div { right: 8px; bottom: 12px; }
  .guest-footer__inner { flex-wrap: wrap; }
  .guest-footer p { order: 3; width: 100%; }
  .member-hero { padding: 28px 0 76px; }
  .member-hero__grid { display: block; }
  .location-card { min-width: 0; margin-top: 22px; }
  .member-search { grid-template-columns: auto 1fr; padding-right: 16px; }
  .member-search .button { display: none; }
  .active-order { grid-template-columns: 72px 1fr; }
  .active-order__visual { width: 72px; height: 72px; }
  .active-order > .button { grid-column: 1 / -1; }
}

@media (max-width: 420px) {
  .guest-hero__actions .button { width: 100%; }
  .floating-offer { grid-template-columns: auto 1fr; }
  .floating-offer b { display: none; }
}
</style>
