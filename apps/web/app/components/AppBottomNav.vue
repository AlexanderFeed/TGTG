<script setup lang="ts">
import { Compass, Home, PackageCheck, Search, UserRound } from 'lucide-vue-next'

const route = useRoute()

const links = [
  { label: 'Главная', to: '/browse', icon: Home },
  { label: 'Рядом', to: '/discover', icon: Compass },
  { label: 'Каталог', to: '/', icon: Search },
  { label: 'Заказ', to: '/delivery', icon: PackageCheck },
  { label: 'Профиль', to: '/profile', icon: UserRound },
]

const isActive = (to: string) => to === '/' ? route.path === '/' : route.path.startsWith(to)
</script>

<template>
  <nav class="bottom-nav" aria-label="Мобильная навигация">
    <NuxtLink
      v-for="link in links"
      :key="link.to"
      :to="link.to"
      class="bottom-nav__link"
      :class="{ 'bottom-nav__link--active': isActive(link.to) }"
    >
      <component :is="link.icon" :size="20" :stroke-width="isActive(link.to) ? 2.6 : 2" />
      <span>{{ link.label }}</span>
    </NuxtLink>
  </nav>
</template>

<style scoped>
.bottom-nav {
  position: fixed;
  z-index: 60;
  right: 10px;
  bottom: max(10px, env(safe-area-inset-bottom));
  left: 10px;
  display: none;
  min-height: 66px;
  grid-template-columns: repeat(5, 1fr);
  padding: 6px;
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 22px;
  background: rgba(16, 43, 36, 0.96);
  box-shadow: var(--shadow-lg);
  backdrop-filter: blur(16px);
}

.bottom-nav__link {
  display: grid;
  align-content: center;
  justify-items: center;
  gap: 3px;
  border-radius: 16px;
  color: rgba(255, 255, 255, 0.65);
  font-size: 0.64rem;
  font-weight: 750;
}

.bottom-nav__link--active {
  color: var(--forest-950);
  background: var(--lime-300);
}

@media (max-width: 720px) {
  .bottom-nav {
    display: grid;
  }
}
</style>
