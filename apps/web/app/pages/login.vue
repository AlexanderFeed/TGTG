<script setup lang="ts">
import { ArrowLeft, KeyRound, LockKeyhole, Phone, ShieldCheck } from 'lucide-vue-next'

const { login, isAuthenticated } = useAuth()
const step = ref<'phone' | 'code'>('phone')
const phone = ref('+7 ')
const code = ref('')
const isSubmitting = ref(false)
const error = ref('')

if (isAuthenticated.value) {
  await navigateTo('/')
}

const normalizedPhone = computed(() => phone.value.replace(/\D/g, ''))

const requestCode = async () => {
  error.value = ''
  if (normalizedPhone.value.length < 11) {
    error.value = 'Введите российский номер из 11 цифр'
    return
  }
  isSubmitting.value = true
  await new Promise(resolve => setTimeout(resolve, 450))
  isSubmitting.value = false
  step.value = 'code'
}

const verifyCode = async () => {
  error.value = ''
  if (code.value.replace(/\D/g, '').length !== 4) {
    error.value = 'Введите любые 4 цифры для демо-входа'
    return
  }
  isSubmitting.value = true
  await new Promise(resolve => setTimeout(resolve, 450))
  login(phone.value)
  await navigateTo('/')
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
          <span><LockKeyhole :size="19" /> Данные остаются в вашем браузере</span>
        </div>
      </div>
      <img src="/images/cafe-rescue.png" alt="Десерты и кофе в уютном кафе">
    </div>

    <div class="auth-page__main">
      <div class="auth-card">
        <div class="auth-card__mobile-logo"><AppLogo /></div>
        <div class="auth-card__icon"><Phone v-if="step === 'phone'" :size="25" /><KeyRound v-else :size="25" /></div>
        <p class="eyebrow">С возвращением</p>
        <h1>{{ step === 'phone' ? 'Войти по номеру' : 'Введите код' }}</h1>
        <p class="auth-card__lead">
          {{ step === 'phone' ? 'Мы отправим короткий код подтверждения.' : `Демо-код отправлен на ${phone}. Подойдут любые 4 цифры.` }}
        </p>

        <form v-if="step === 'phone'" class="form-grid" @submit.prevent="requestCode">
          <div class="field">
            <label for="login-phone">Номер телефона</label>
            <div class="input-wrap">
              <Phone :size="18" />
              <input id="login-phone" v-model="phone" class="input input--icon" type="tel" inputmode="tel" autocomplete="tel" placeholder="+7 999 123-45-67">
            </div>
          </div>
          <p v-if="error" class="auth-card__error">{{ error }}</p>
          <button class="button button--primary button--block" type="submit" :disabled="isSubmitting">
            {{ isSubmitting ? 'Отправляем…' : 'Получить код' }}
          </button>
        </form>

        <form v-else class="form-grid" @submit.prevent="verifyCode">
          <div class="field">
            <label for="login-code">Код из SMS</label>
            <input id="login-code" v-model="code" class="input auth-card__code" type="text" inputmode="numeric" maxlength="4" autocomplete="one-time-code" placeholder="••••" autofocus>
          </div>
          <p v-if="error" class="auth-card__error">{{ error }}</p>
          <button class="button button--primary button--block" type="submit" :disabled="isSubmitting">
            {{ isSubmitting ? 'Проверяем…' : 'Войти' }}
          </button>
          <button class="button button--ghost button--block" type="button" @click="step = 'phone'; code = ''; error = ''">Изменить номер</button>
        </form>

        <p class="auth-card__switch">Впервые здесь? <NuxtLink to="/register">Создать аккаунт</NuxtLink></p>
        <p class="auth-card__demo">Пет-проект: OTP и авторизация сейчас работают локально и предназначены только для демонстрации интерфейса.</p>
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
