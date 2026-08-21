<script setup lang="ts">
import { ChevronDown, Heart, MapPin, UserRound } from 'lucide-vue-next'

const route = useRoute()
const { user, isAuthenticated } = useAuth()

const links = [
  { label: 'Главная', to: '/' },
  { label: 'Рядом', to: '/discover' },
  { label: 'Каталог', to: '/browse' },
  { label: 'Получение', to: '/delivery' },
]

const isActive = (to: string) => to === '/' ? route.path === '/' : route.path.startsWith(to)
</script>

<template>
  <header class="header">
    <div class="container header__inner">
      <AppLogo />

      <nav v-if="isAuthenticated" class="header__nav" aria-label="Основная навигация">
        <NuxtLink
          v-for="link in links"
          :key="link.to"
          :to="link.to"
          class="header__link"
          :class="{ 'header__link--active': isActive(link.to) }"
        >
          {{ link.label }}
        </NuxtLink>
      </nav>

      <div class="header__actions">
        <NuxtLink v-if="isAuthenticated" class="location-pill" to="/discover">
          <MapPin :size="16" />
          <span>{{ user?.city || 'Москва' }}</span>
          <ChevronDown :size="14" />
        </NuxtLink>

        <template v-if="isAuthenticated">
          <NuxtLink class="icon-button header__heart" to="/profile#favorites" aria-label="Избранное">
            <Heart :size="19" />
          </NuxtLink>
          <NuxtLink class="header__profile" to="/profile" aria-label="Профиль">
            <span class="avatar">{{ user?.name?.charAt(0) || 'Я' }}</span>
            <span class="header__profile-name">{{ user?.name?.split(' ')[0] }}</span>
          </NuxtLink>
        </template>

        <template v-else>
          <NuxtLink class="button button--ghost header__login" to="/login">
            <UserRound :size="18" />
            Войти
          </NuxtLink>
          <NuxtLink class="button button--primary header__register" to="/register">
            Создать аккаунт
          </NuxtLink>
        </template>
      </div>
    </div>
  </header>
</template>

<style scoped>
.header {
  position: sticky;
  z-index: 50;
  top: 0;
  border-bottom: 1px solid rgba(226, 223, 211, 0.78);
  background: rgba(251, 250, 246, 0.88);
  backdrop-filter: blur(18px);
}

.header__inner {
  display: flex;
  height: 72px;
  align-items: center;
  gap: 36px;
}

.header__nav {
  display: flex;
  height: 100%;
  align-items: center;
  gap: 28px;
}

.header__link {
  position: relative;
  display: grid;
  height: 100%;
  place-items: center;
  color: var(--ink-700);
  font-size: 0.9rem;
  font-weight: 750;
}

.header__link::after {
  position: absolute;
  right: 0;
  bottom: -1px;
  left: 0;
  height: 3px;
  border-radius: 3px 3px 0 0;
  content: "";
  background: var(--coral-500);
  opacity: 0;
  transform: scaleX(0.35);
  transition: opacity 160ms ease, transform 160ms ease;
}

.header__link:hover,
.header__link--active {
  color: var(--forest-900);
}

.header__link--active::after {
  opacity: 1;
  transform: scaleX(1);
}

.header__actions {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-left: auto;
}

.location-pill {
  display: inline-flex;
  min-height: 40px;
  align-items: center;
  gap: 6px;
  padding: 0 13px;
  border: 1px solid var(--sand-200);
  border-radius: 999px;
  color: var(--forest-900);
  background: var(--white);
  font-size: 0.85rem;
  font-weight: 750;
}

.header__profile {
  display: flex;
  align-items: center;
  gap: 9px;
  margin-left: 3px;
  color: var(--forest-900);
  font-size: 0.9rem;
  font-weight: 800;
}

.header__profile .avatar {
  width: 40px;
  height: 40px;
}

.header__login,
.header__register {
  min-height: 44px;
}

@media (max-width: 1020px) {
  .header__nav,
  .location-pill,
  .header__profile-name {
    display: none;
  }
}

@media (max-width: 720px) {
  .header__inner {
    height: 64px;
  }

  .header__heart,
  .header__profile {
    display: none;
  }

  .header__login {
    min-height: 40px;
    padding-inline: 12px;
  }

  .header__register {
    display: none;
  }
}
</style>
