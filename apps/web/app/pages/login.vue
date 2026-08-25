<script setup lang="ts">
import { ArrowLeft, KeyRound, LockKeyhole, Mail, ShieldCheck } from 'lucide-vue-next'

const { requestLoginCode, verifyLoginCode, errorMessage, isAuthenticated } = useAuth()

// This page is a two-step state machine. `email` shows the first form; after Go
// creates a challenge, `code` shows the verification form for that challenge.
const step = ref<'email' | 'code'>('email')
const email = ref('')
const code = ref('')
const challengeId = ref('')
const devCode = ref('')
const isSubmitting = ref(false)
const error = ref('')

if (isAuthenticated.value) {
  await navigateTo('/')
}

// Vue calls this because the template uses @submit.prevent="requestCode".
// `.prevent` stops the browser's normal full-page form submission.
const requestCode = async () => {
  error.value = ''
  if (!email.value.trim() || !email.value.includes('@')) {
    error.value = 'Введите корректный email'
    return
  }
  isSubmitting.value = true
  try {
    // useAuth -> Nuxt /api proxy -> Go requestLogin handler -> PostgreSQL/mailer.
    const response = await requestLoginCode(email.value)
    challengeId.value = response.challengeId
    devCode.value = response.devCode || ''
    if (response.devCode) code.value = response.devCode
    step.value = 'code'
  } catch (requestError) {
    error.value = errorMessage(requestError)
  } finally {
    isSubmitting.value = false
  }
}

// The second request combines the challenge ID, same email, and code. On
// success Go sets the HttpOnly cookie and useAuth stores the returned user.
const verifyCode = async () => {
  error.value = ''
  if (!/^\d{6}$/.test(code.value)) {
    error.value = 'Введите код из 6 цифр'
    return
  }
  isSubmitting.value = true
  try {
    await verifyLoginCode(challengeId.value, email.value, code.value)
    await navigateTo('/')
  } catch (verifyError) {
    error.value = errorMessage(verifyError)
  } finally {
    isSubmitting.value = false
  }
}

useSeoMeta({ title: 'Войти — ЕщёЕсть' })
</script>

<template>
  <section class="auth-page">
    <div class="auth-page__aside">
      <NuxtLink class="auth-page__back" to="/"><ArrowLeft :size="18" /> На главную</NuxtLink>
      <div class="auth-page__aside-content">
        <AppLogo inverse />
        <blockquote>«Каждый спасённый пакет — это вкусный план на вечер и чуть меньше лишних отходов.»</blockquote>
        <div class="auth-page__aside-points">
          <span><ShieldCheck :size="19" /> Демо без настоящих платежей</span>
          <span><LockKeyhole :size="19" /> Защищённая серверная сессия</span>
        </div>
      </div>
      <img src="/images/cafe-rescue.png" alt="Десерты и кофе в уютном кафе">
    </div>

    <div class="auth-page__main">
      <div class="auth-card">
        <div class="auth-card__mobile-logo"><AppLogo /></div>
        <div class="auth-card__icon"><Mail v-if="step === 'email'" :size="25" /><KeyRound v-else :size="25" /></div>
        <p class="eyebrow">С возвращением</p>
        <h1>{{ step === 'email' ? 'Войти по email' : 'Введите код' }}</h1>
        <p class="auth-card__lead">
          {{ step === 'email' ? 'Мы отправим шестизначный код подтверждения.' : `Код отправлен на ${email}.` }}
          <strong v-if="devCode"> Локальный код: {{ devCode }}</strong>
        </p>

        <form v-if="step === 'email'" class="form-grid" @submit.prevent="requestCode">
          <div class="field">
            <label for="login-email">Email</label>
            <div class="input-wrap">
              <Mail :size="18" />
              <input id="login-email" v-model="email" class="input input--icon" type="email" autocomplete="email" placeholder="you@example.ru">
            </div>
          </div>
          <p v-if="error" class="auth-card__error">{{ error }}</p>
          <button class="button button--primary button--block" type="submit" :disabled="isSubmitting">
            {{ isSubmitting ? 'Отправляем…' : 'Получить код' }}
          </button>
        </form>

        <form v-else class="form-grid" @submit.prevent="verifyCode">
          <div class="field">
            <label for="login-code">Код из письма</label>
            <input id="login-code" v-model="code" class="input auth-card__code" type="text" inputmode="numeric" maxlength="6" autocomplete="one-time-code" placeholder="••••••" autofocus>
          </div>
          <p v-if="error" class="auth-card__error">{{ error }}</p>
          <button class="button button--primary button--block" type="submit" :disabled="isSubmitting">
            {{ isSubmitting ? 'Проверяем…' : 'Войти' }}
          </button>
          <button class="button button--ghost button--block" type="button" @click="step = 'email'; code = ''; devCode = ''; error = ''">Изменить email</button>
        </form>

        <p class="auth-card__switch">Впервые здесь? <NuxtLink to="/register">Создать аккаунт</NuxtLink></p>
        <p class="auth-card__demo">Пет-проект: аккаунт и сессия сохраняются в PostgreSQL. В локальном режиме код показывается на экране; на VPS он отправляется через SMTP.</p>
      </div>
    </div>
  </section>
</template>

<style scoped>
.auth-page { display: grid; min-height: calc(100vh - 72px); grid-template-columns: .9fr 1.1fr; background: var(--white); }
.auth-page__aside { position: relative; min-height: 720px; overflow: hidden; padding: 34px; color: var(--white); background: var(--forest-950); }
.auth-page__aside::after { position: absolute; inset: 0; content: ""; background: linear-gradient(180deg, rgba(16,43,36,.02) 20%, rgba(16,43,36,.96) 88%); }
.auth-page__aside > img { position: absolute; inset: 0; width: 100%; height: 100%; object-fit: cover; opacity: .62; }
.auth-page__back { position: relative; z-index: 2; display: inline-flex; align-items: center; gap: 8px; color: rgba(255,255,255,.82); font-size: .83rem; font-weight: 750; }
.auth-page__aside-content { position: absolute; z-index: 2; right: 44px; bottom: 50px; left: 44px; }
.auth-page__aside-content blockquote { max-width: 550px; margin: 32px 0; font-size: clamp(1.5rem, 3vw, 2.4rem); font-weight: 750; letter-spacing: -.04em; line-height: 1.25; }
.auth-page__aside-points { display: flex; flex-wrap: wrap; gap: 14px 22px; color: rgba(255,255,255,.69); font-size: .75rem; }
.auth-page__aside-points span { display: flex; align-items: center; gap: 7px; }
.auth-page__aside-points svg { color: var(--lime-300); }
.auth-page__main { display: grid; place-items: center; padding: 54px 32px; }
.auth-card { width: min(100%, 440px); }
.auth-card__mobile-logo { display: none; margin-bottom: 48px; }
.auth-card__icon { display: grid; width: 54px; height: 54px; margin-bottom: 24px; place-items: center; border-radius: 17px; color: var(--forest-950); background: var(--lime-300); }
.auth-card h1 { margin-bottom: 12px; font-size: clamp(2rem, 4vw, 3.1rem); }
.auth-card__lead { margin-bottom: 30px; color: var(--ink-700); line-height: 1.65; }
.auth-card__error { margin: -5px 0 0; color: var(--danger); font-size: .8rem; font-weight: 720; }
.auth-card__code { text-align: center; font-size: 1.65rem; font-weight: 900; letter-spacing: .55em; }
.auth-card__switch { margin: 26px 0 0; text-align: center; color: var(--ink-700); font-size: .88rem; }
.auth-card__switch a { color: var(--forest-900); font-weight: 900; }
.auth-card__demo { margin: 32px 0 0; padding-top: 20px; border-top: 1px solid var(--cream-100); color: var(--ink-500); font-size: .7rem; line-height: 1.55; text-align: center; }
@media (max-width: 820px) { .auth-page { grid-template-columns: 1fr; min-height: calc(100vh - 64px); } .auth-page__aside { display: none; } .auth-card__mobile-logo { display: block; } }
@media (max-width: 520px) { .auth-page__main { align-items: start; padding: 35px 20px 60px; } }
</style>
