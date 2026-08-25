<script setup lang="ts">
import { Bell, ChevronRight, CircleHelp, Edit3, Heart, LogOut, Mail, MapPin, PackageCheck, Save, Settings, ShieldCheck, Smartphone, UserRound } from 'lucide-vue-next'

// Nuxt runs app/middleware/auth.ts before this protected page is rendered.
definePageMeta({ middleware: 'auth' })

const { user, updateProfile, errorMessage, logout } = useAuth()
const { favorites } = useMarketplace()
const editMode = ref(false)
const saved = ref(false)
const saving = ref(false)
const saveError = ref('')
const form = reactive({
  name: user.value?.name || '',
  city: user.value?.city || 'Москва',
})
const notifications = reactive({
  nearby: true,
  pickup: true,
  news: false,
})

// Keep the form synchronized when refreshSession replaces the shared user, but
// do not overwrite fields while the visitor is actively editing them.
watch(user, (current) => {
  if (!current || editMode.value) return
  form.name = current.name
  form.city = current.city
}, { immediate: true })

const memberSince = computed(() => user.value?.createdAt
  ? new Intl.DateTimeFormat('ru-RU', { month: 'long', year: 'numeric' }).format(new Date(user.value.createdAt))
  : '')

// The template's @submit.prevent calls save. updateProfile sends PATCH through
// useAuth and updates the shared user so the header/profile change together.
const save = async () => {
  saveError.value = ''
  saving.value = true
  try {
    await updateProfile(form)
    editMode.value = false
    saved.value = true
    setTimeout(() => saved.value = false, 1800)
  } catch (requestError) {
    saveError.value = errorMessage(requestError)
  } finally {
    saving.value = false
  }
}

useSeoMeta({ title: 'Профиль — ЕщёЕсть' })
</script>

<template>
  <section class="profile-page page-section page-section--compact">
    <div class="container">
      <div class="profile-hero surface-dark">
        <div class="profile-hero__avatar">{{ user?.name?.charAt(0) }}</div>
        <div class="profile-hero__copy">
          <p>Ваш профиль</p>
          <h1>{{ user?.name }}</h1>
          <span><MapPin :size="15" /> {{ user?.city }} · с нами с {{ memberSince }}</span>
        </div>
        <button class="button button--accent" type="button" @click="editMode = !editMode"><Edit3 :size="17" /> {{ editMode ? 'Отменить' : 'Редактировать' }}</button>
        <div class="profile-hero__stats">
          <div><strong>{{ user?.impact.rescued }}</strong><span>пакетов спасено</span></div>
          <div><strong>{{ user?.impact.savedRubles.toLocaleString('ru-RU') }} ₽</strong><span>сэкономлено</span></div>
          <div><strong>{{ user?.impact.co2Kg }} кг</strong><span>CO₂ сбережено</span></div>
        </div>
      </div>

      <p v-if="saved" class="profile-saved"><Save :size="16" /> Изменения сохранены</p>

      <div class="profile-layout">
        <div class="profile-main">
          <article class="profile-card card">
            <div class="profile-card__heading"><div><h2>Личные данные</h2><p>Профиль хранится в PostgreSQL, а email подтверждён.</p></div><UserRound :size="23" /></div>
            <form class="profile-form" @submit.prevent="save">
              <div class="field"><label for="profile-name">Имя</label><input id="profile-name" v-model="form.name" class="input" :disabled="!editMode"></div>
              <div class="field"><label for="profile-email">Подтверждённый email</label><input id="profile-email" class="input" type="email" :value="user?.email" disabled></div>
              <div class="field"><label for="profile-city">Город</label><select id="profile-city" v-model="form.city" class="select" :disabled="!editMode"><option>Москва</option><option>Санкт-Петербург</option><option>Казань</option></select></div>
              <p v-if="saveError" class="profile-form__error">{{ saveError }}</p>
              <button v-if="editMode" class="button button--primary" type="submit" :disabled="saving"><Save :size="17" /> {{ saving ? 'Сохраняем…' : 'Сохранить' }}</button>
            </form>
          </article>

          <article id="favorites" class="profile-card card">
            <div class="profile-card__heading"><div><h2>Избранное</h2><p>{{ favorites.length ? 'Места и пакеты, которые вы сохранили.' : 'Пока здесь ничего нет.' }}</p></div><Heart :size="23" /></div>
            <div v-if="favorites.length" class="profile-favorites">
              <OfferCard v-for="offer in favorites" :key="offer.id" :offer="offer" compact />
            </div>
            <div v-else class="profile-empty">
              <span>♡</span><p>Нажмите сердечко на карточке предложения, чтобы вернуться к нему позже.</p><NuxtLink class="button button--secondary" to="/">Открыть каталог</NuxtLink>
            </div>
          </article>

          <article class="profile-card card">
            <div class="profile-card__heading"><div><h2>Уведомления</h2><p>Решите, о чём напоминать в демо-приложении.</p></div><Bell :size="23" /></div>
            <div class="settings-list">
              <label><span><strong>Новые пакеты рядом</strong><small>Когда любимое место публикует предложение</small></span><input v-model="notifications.nearby" type="checkbox"><i /></label>
              <label><span><strong>Напоминания о получении</strong><small>За 30 минут до начала окна</small></span><input v-model="notifications.pickup" type="checkbox"><i /></label>
              <label><span><strong>Новости продукта</strong><small>Редкие обновления пет-проекта</small></span><input v-model="notifications.news" type="checkbox"><i /></label>
            </div>
          </article>
        </div>

        <aside class="profile-aside">
          <article class="profile-menu card">
            <NuxtLink to="/delivery"><span><PackageCheck :size="19" /> Мои заказы</span><ChevronRight :size="18" /></NuxtLink>
            <button type="button"><span><MapPin :size="19" /> Адреса</span><ChevronRight :size="18" /></button>
            <button type="button"><span><Smartphone :size="19" /> Устройства</span><ChevronRight :size="18" /></button>
            <button type="button"><span><ShieldCheck :size="19" /> Приватность</span><ChevronRight :size="18" /></button>
            <button type="button"><span><Settings :size="19" /> Настройки</span><ChevronRight :size="18" /></button>
          </article>

          <article class="profile-contact card">
            <span><CircleHelp :size="22" /></span><div><h3>Есть вопрос?</h3><p>Пока это пет-проект, но интерфейс поддержки уже можно проверить.</p><button type="button"><Mail :size="16" /> Написать нам</button></div>
          </article>

          <button class="profile-logout" type="button" @click="logout"><LogOut :size="18" /> Выйти из аккаунта</button>
        </aside>
      </div>
    </div>
  </section>
</template>

<style scoped>
.profile-hero { position: relative; display: grid; grid-template-columns: auto 1fr auto; align-items: center; gap: 20px; padding: 34px; overflow: hidden; border-radius: 30px; background-image: radial-gradient(circle at 88% 5%, rgba(221,244,100,.25), transparent 30%); }
.profile-hero::after { position: absolute; right: -65px; bottom: -85px; width: 230px; height: 230px; border: 26px solid rgba(255,255,255,.07); border-radius: 50%; content: ""; }
.profile-hero__avatar { position: relative; z-index: 1; display: grid; width: 86px; height: 86px; place-items: center; border: 5px solid rgba(255,255,255,.14); border-radius: 27px; color: var(--forest-950); background: var(--lime-300); font-size: 2rem; font-weight: 950; }
.profile-hero__copy { position: relative; z-index: 1; }
.profile-hero__copy p { margin-bottom: 5px; color: var(--mint-300); font-size: .75rem; font-weight: 850; text-transform: uppercase; }
.profile-hero h1 { margin-bottom: 9px; font-size: clamp(2.1rem, 4vw, 3.5rem); }
.profile-hero__copy span { display: flex; align-items: center; gap: 6px; color: rgba(255,255,255,.62); font-size: .76rem; }
.profile-hero > .button { position: relative; z-index: 1; }
.profile-hero__stats { position: relative; z-index: 1; display: grid; grid-column: 1 / -1; grid-template-columns: repeat(3, 1fr); gap: 10px; padding-top: 25px; border-top: 1px solid rgba(255,255,255,.1); }
.profile-hero__stats div { display: grid; gap: 4px; }
.profile-hero__stats strong { color: var(--lime-300); font-size: 1.35rem; }
.profile-hero__stats span { color: rgba(255,255,255,.56); font-size: .7rem; }
.profile-saved { position: fixed; z-index: 100; right: 25px; bottom: 25px; display: flex; align-items: center; gap: 7px; padding: 12px 15px; border-radius: 999px; color: var(--white); background: var(--success); box-shadow: var(--shadow-md); font-size: .78rem; font-weight: 800; }
.profile-layout { display: grid; grid-template-columns: minmax(0, 1.6fr) minmax(260px, .65fr); align-items: start; gap: 24px; margin-top: 24px; }
.profile-main, .profile-aside { display: grid; gap: 20px; }
.profile-card { padding: 25px; }
.profile-card__heading { display: flex; align-items: start; justify-content: space-between; gap: 20px; margin-bottom: 24px; padding-bottom: 19px; border-bottom: 1px solid var(--cream-100); }
.profile-card__heading h2 { margin-bottom: 5px; font-size: 1.35rem; }
.profile-card__heading p { margin: 0; color: var(--ink-500); font-size: .75rem; }
.profile-card__heading > svg { color: var(--coral-500); }
.profile-form { display: grid; grid-template-columns: 1fr 1fr; gap: 18px; }
.profile-form .button { justify-self: start; }
.profile-form__error { grid-column: 1 / -1; margin: 0; color: var(--danger); font-size: .8rem; font-weight: 720; }
.profile-form .input:disabled, .profile-form .select:disabled { color: var(--ink-700); background: var(--cream-50); opacity: 1; }
.profile-favorites { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.profile-empty { display: grid; justify-items: center; padding: 30px 15px 12px; text-align: center; }
.profile-empty > span { color: var(--mint-300); font-size: 4rem; line-height: 1; }
.profile-empty p { max-width: 430px; margin: 10px 0 20px; color: var(--ink-700); font-size: .82rem; line-height: 1.55; }
.settings-list { display: grid; }
.settings-list label { display: flex; align-items: center; justify-content: space-between; gap: 20px; padding: 15px 0; cursor: pointer; }
.settings-list label + label { border-top: 1px solid var(--cream-100); }
.settings-list label > span { display: grid; gap: 4px; }
.settings-list small { color: var(--ink-500); font-size: .72rem; }
.settings-list input { position: absolute; opacity: 0; }
.settings-list i { position: relative; width: 46px; height: 26px; flex: 0 0 auto; border-radius: 999px; background: var(--sand-200); transition: background 160ms ease; }
.settings-list i::after { position: absolute; top: 4px; left: 4px; width: 18px; height: 18px; border-radius: 50%; content: ""; background: var(--white); box-shadow: 0 2px 5px rgba(0,0,0,.14); transition: transform 160ms ease; }
.settings-list input:checked + i { background: var(--forest-700); }
.settings-list input:checked + i::after { transform: translateX(20px); }
.profile-menu { overflow: hidden; }
.profile-menu a, .profile-menu button { display: flex; width: 100%; min-height: 54px; align-items: center; justify-content: space-between; gap: 15px; padding: 0 17px; color: var(--ink-700); background: var(--white); cursor: pointer; font-size: .79rem; font-weight: 750; }
.profile-menu a + *, .profile-menu button + * { border-top: 1px solid var(--cream-100); }
.profile-menu a:hover, .profile-menu button:hover { color: var(--forest-900); background: var(--cream-50); }
.profile-menu span { display: flex; align-items: center; gap: 10px; }
.profile-menu span svg { color: var(--forest-700); }
.profile-contact { display: flex; align-items: flex-start; gap: 12px; padding: 20px; }
.profile-contact > span { display: grid; width: 42px; height: 42px; flex: 0 0 auto; place-items: center; border-radius: 14px; color: var(--forest-950); background: var(--lime-300); }
.profile-contact h3 { margin-bottom: 7px; }
.profile-contact p { color: var(--ink-700); font-size: .74rem; line-height: 1.5; }
.profile-contact button { display: flex; align-items: center; gap: 6px; padding: 0; color: var(--forest-700); background: transparent; cursor: pointer; font-size: .75rem; font-weight: 850; }
.profile-logout { display: flex; min-height: 48px; align-items: center; justify-content: center; gap: 8px; border: 1px solid var(--coral-100); border-radius: 15px; color: var(--danger); background: #fff8f6; cursor: pointer; font-weight: 800; }
@media (max-width: 900px) { .profile-layout { grid-template-columns: 1fr; } .profile-aside { grid-template-columns: 1fr 1fr; } .profile-logout { grid-column: 1 / -1; } }
@media (max-width: 650px) { .profile-hero { grid-template-columns: auto 1fr; padding: 24px 20px; } .profile-hero__avatar { width: 64px; height: 64px; border-radius: 21px; } .profile-hero > .button { grid-column: 1 / -1; } .profile-hero__stats { gap: 5px; } .profile-hero__stats strong { font-size: 1rem; } .profile-hero__stats span { font-size: .58rem; } .profile-form { grid-template-columns: 1fr; } .profile-favorites { grid-template-columns: 1fr; } .profile-aside { grid-template-columns: 1fr; } .profile-logout { grid-column: auto; } }
</style>
