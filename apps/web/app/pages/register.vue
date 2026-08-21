<script setup lang="ts">
import { ArrowLeft, Check, Mail, Phone, Sparkles, UserRound } from 'lucide-vue-next'

const { register, isAuthenticated } = useAuth()
const form = reactive({
  name: '',
  phone: '+7 ',
  email: '',
})
const isSubmitting = ref(false)
const error = ref('')

if (isAuthenticated.value) {
  await navigateTo('/')
}

const submit = async () => {
  error.value = ''
  if (form.name.trim().length < 2) {
    error.value = 'Расскажите, как к вам обращаться'
    return
  }
  if (form.phone.replace(/\D/g, '').length < 11) {
    error.value = 'Введите российский номер из 11 цифр'
    return
  }
  isSubmitting.value = true
  await new Promise(resolve => setTimeout(resolve, 500))
  register(form)
  await navigateTo('/')
}

useSeoMeta({ title: 'Регистрация — ЕщёЕсть' })
</script>

<template>
  <section class="register-page">
    <div class="register-page__main">
      <div class="register-card">
        <NuxtLink class="register-card__back" to="/"><ArrowLeft :size="18" /> На главную</NuxtLink>
        <div class="register-card__heading">
          <AppLogo />
          <p class="eyebrow">Новый аккаунт</p>
          <!-- <h1>Начнём спасать вкусное?</h1> -->
          <p>Регистрация занимает меньше минуты. Для пет-проекта данные сохраняются локально.</p>
        </div>

        <form class="form-grid" @submit.prevent="submit">
          <div class="field">
            <label for="register-name">Ваше логин</label>
            <div class="input-wrap">
              <UserRound :size="18" />
              <input id="register-name" v-model="form.name" class="input input--icon" type="text" autocomplete="name" placeholder="Логин">
            </div>
          </div>
          <div class="field">
            <label for="register-phone">Номер телефона</label>
            <div class="input-wrap">
              <Phone :size="18" />
              <input id="register-phone" v-model="form.phone" class="input input--icon" type="tel" inputmode="tel" autocomplete="tel" placeholder="+7 999 123-45-67">
            </div>
          </div>
          <div class="field">
            <label for="register-email">Email <span>(необязательно)</span></label>
            <div class="input-wrap">
              <Mail :size="18" />
              <input id="register-email" v-model="form.email" class="input input--icon" type="email" autocomplete="email" placeholder="you@example.ru">
            </div>
          </div>
          <p v-if="error" class="register-card__error">{{ error }}</p>
          <button class="button button--primary button--block" type="submit" :disabled="isSubmitting">
            <Sparkles :size="18" /> {{ isSubmitting ? 'Создаём…' : 'Создать аккаунт' }}
          </button>
        </form>

        <p class="register-card__switch">Уже есть аккаунт? <NuxtLink to="/login">Войти</NuxtLink></p>
      </div>
    </div>

    <aside class="register-page__aside">
      <div class="register-page__image"><img src="/images/grocery-rescue.png" alt="Сумка со свежими продуктами"></div>
      <div class="register-page__benefits">
        <p class="eyebrow">После регистрации</p>
        <h2>Всё нужное будет под рукой</h2>
        <ul>
          <li><span><Check :size="17" /></span><div><strong>Избранные места</strong><small>Сохраняйте то, что хочется попробовать.</small></div></li>
          <li><span><Check :size="17" /></span><div><strong>Удобное получение</strong><small>Код, время и адрес всегда в активном заказе.</small></div></li>
          <li><span><Check :size="17" /></span><div><strong>Личный результат</strong><small>Следите за экономией и спасёнными пакетами.</small></div></li>
        </ul>
      </div>
    </aside>
  </section>
</template>

<style scoped>
.register-page { display: grid; min-height: calc(100vh - 72px); grid-template-columns: 1.05fr .95fr; background: var(--white); }
.register-page__main { display: grid; place-items: center; padding: 48px 32px 70px; }
.register-card { width: min(100%, 470px); }
.register-card__back { display: inline-flex; align-items: center; gap: 8px; margin-bottom: 38px; color: var(--ink-700); font-size: .82rem; font-weight: 750; }
.register-card__heading > :deep(.brand) { margin-bottom: 38px; }
.register-card__heading h1 { margin-bottom: 13px; font-size: clamp(2rem, 4vw, 3.1rem); }
.register-card__heading > p:last-child { margin-bottom: 28px; color: var(--ink-700); line-height: 1.6; }
.field label span { color: var(--ink-500); font-weight: 600; }
.register-card__error { margin: -5px 0 0; color: var(--danger); font-size: .8rem; font-weight: 720; }
.register-card__switch { margin: 24px 0 0; color: var(--ink-700); font-size: .86rem; text-align: center; }
.register-card__switch a { color: var(--forest-900); font-weight: 900; }
.register-page__aside { position: relative; display: grid; align-content: end; min-height: 720px; overflow: hidden; padding: 40px; color: var(--white); background: var(--forest-900); }
.register-page__aside::after { position: absolute; inset: 0; content: ""; background: linear-gradient(180deg, rgba(16,43,36,.02), rgba(16,43,36,.98) 74%); }
.register-page__image { position: absolute; inset: 0; }
.register-page__image img { width: 100%; height: 100%; object-fit: cover; opacity: .7; }
.register-page__benefits { position: relative; z-index: 2; padding: 10px 14px; }
.register-page__benefits .eyebrow { color: var(--mint-300); }
.register-page__benefits h2 { max-width: 470px; margin-bottom: 25px; }
.register-page__benefits ul { display: grid; gap: 18px; margin: 0; padding: 0; list-style: none; }
.register-page__benefits li { display: flex; align-items: flex-start; gap: 12px; }
.register-page__benefits li > span { display: grid; width: 30px; height: 30px; flex: 0 0 auto; place-items: center; border-radius: 50%; color: var(--forest-950); background: var(--lime-300); }
.register-page__benefits li div { display: grid; gap: 4px; }
.register-page__benefits small { color: rgba(255,255,255,.65); line-height: 1.45; }
@media (max-width: 840px) { .register-page { grid-template-columns: 1fr; } .register-page__aside { display: none; } }
@media (max-width: 520px) { .register-page__main { align-items: start; padding: 30px 20px 60px; } .register-card__heading > :deep(.brand) { margin-bottom: 30px; } }
</style>
