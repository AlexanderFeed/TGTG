<script setup lang="ts">
import { ArrowLeft, Check, KeyRound, Mail, Sparkles, UserRound } from 'lucide-vue-next'

const { requestRegistrationCode, verifyRegistrationCode, errorMessage, isAuthenticated } = useAuth()

// Registration also has two steps. The first request only creates an email
// challenge; the `users` row is created by Go after successful verification.
const step = ref<'details' | 'code'>('details')
const form = reactive({
  name: '',
  email: '',
})
const code = ref('')
const challengeId = ref('')
const devCode = ref('')
const isSubmitting = ref(false)
const error = ref('')

if (isAuthenticated.value) {
  await navigateTo('/')
}

// This is the frontend entry point for POST /api/v1/auth/register/request.
// Lightweight checks give instant feedback; Go repeats authoritative validation.
const requestCode = async () => {
  error.value = ''
  if (form.name.trim().length < 2) {
    error.value = 'Расскажите, как к вам обращаться'
    return
  }
  if (!form.email.trim() || !form.email.includes('@')) {
    error.value = 'Введите корректный email'
    return
  }
  isSubmitting.value = true
  try {
    const response = await requestRegistrationCode(form.name, form.email)
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

// This is the frontend entry point for POST /api/v1/auth/register/verify.
// The composable keeps HTTP details out of this presentational page.
const verifyCode = async () => {
  error.value = ''
  if (!/^\d{6}$/.test(code.value)) {
    error.value = 'Введите код из 6 цифр'
    return
  }
  isSubmitting.value = true
  try {
    await verifyRegistrationCode(challengeId.value, form.email, code.value)
    await navigateTo('/')
  } catch (verifyError) {
    error.value = errorMessage(verifyError)
  } finally {
    isSubmitting.value = false
  }
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
          <p>{{ step === 'details' ? 'Регистрация занимает меньше минуты.' : `Код подтверждения отправлен на ${form.email}.` }} <strong v-if="devCode">Локальный код: {{ devCode }}</strong></p>
        </div>

        <form v-if="step === 'details'" class="form-grid" @submit.prevent="requestCode">
          <div class="field">
            <label for="register-name">Ваше имя</label>
            <div class="input-wrap">
              <UserRound :size="18" />
              <input id="register-name" v-model="form.name" class="input input--icon" type="text" autocomplete="name" placeholder="Алексей">
            </div>
          </div>
          <div class="field">
            <label for="register-email">Email</label>
            <div class="input-wrap">
              <Mail :size="18" />
              <input id="register-email" v-model="form.email" class="input input--icon" type="email" autocomplete="email" placeholder="you@example.ru">
            </div>
          </div>
          <p v-if="error" class="register-card__error">{{ error }}</p>
          <button class="button button--primary button--block" type="submit" :disabled="isSubmitting">
            <Sparkles :size="18" /> {{ isSubmitting ? 'Отправляем…' : 'Получить код' }}
          </button>
        </form>

        <form v-else class="form-grid" @submit.prevent="verifyCode">
          <div class="field">
            <label for="register-code">Код из письма</label>
            <div class="input-wrap">
              <KeyRound :size="18" />
              <input id="register-code" v-model="code" class="input input--icon auth-card__code" type="text" inputmode="numeric" maxlength="6" autocomplete="one-time-code" placeholder="••••••" autofocus>
            </div>
          </div>
          <p v-if="error" class="register-card__error">{{ error }}</p>
          <button class="button button--primary button--block" type="submit" :disabled="isSubmitting">
            <Sparkles :size="18" /> {{ isSubmitting ? 'Проверяем…' : 'Создать аккаунт' }}
          </button>
          <button class="button button--ghost button--block" type="button" @click="step = 'details'; code = ''; devCode = ''; error = ''">Изменить данные</button>
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
#register-code { text-align: center; font-size: 1.35rem; font-weight: 900; letter-spacing: .35em; }
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
