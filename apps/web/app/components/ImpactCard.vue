<script setup lang="ts">
import { BadgeRussianRuble, Leaf, PackageCheck } from 'lucide-vue-next'

const { user } = useAuth()
</script>

<template>
  <section class="impact-card surface-dark">
    <div>
      <p class="impact-card__kicker">Ваш вклад</p>
      <h2>Хорошая еда не пропала зря</h2>
      <p>Небольшие покупки складываются в заметный результат.</p>
    </div>
    <div class="impact-card__stats">
      <div>
        <PackageCheck :size="21" />
        <strong>{{ user?.impact.rescued ?? 0 }}</strong>
        <span>пакетов спасено</span>
      </div>
      <div>
        <BadgeRussianRuble :size="21" />
        <strong>{{ (user?.impact.savedRubles ?? 0).toLocaleString('ru-RU') }} ₽</strong>
        <span>сэкономлено</span>
      </div>
      <div>
        <Leaf :size="21" />
        <strong>{{ user?.impact.co2Kg ?? 0 }} кг</strong>
        <span>CO₂ не потрачено зря</span>
      </div>
    </div>
  </section>
</template>

<style scoped>
.impact-card {
  display: grid;
  grid-template-columns: 0.9fr 1.4fr;
  gap: 44px;
  padding: 38px;
  border-radius: 30px;
  background-image: radial-gradient(circle at 90% 0%, rgba(221, 244, 100, .18), transparent 34%);
}

.impact-card__kicker {
  margin-bottom: 10px;
  color: var(--lime-300);
  font-size: .76rem;
  font-weight: 900;
  letter-spacing: .13em;
  text-transform: uppercase;
}

.impact-card h2 {
  max-width: 430px;
  margin-bottom: 12px;
  font-size: clamp(1.65rem, 3vw, 2.4rem);
}

.impact-card p:last-child {
  margin-bottom: 0;
  color: rgba(255,255,255,.68);
  line-height: 1.6;
}

.impact-card__stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  align-items: stretch;
  gap: 12px;
}

.impact-card__stats > div {
  display: grid;
  align-content: center;
  gap: 8px;
  min-height: 150px;
  padding: 20px;
  border: 1px solid rgba(255,255,255,.1);
  border-radius: 20px;
  background: rgba(255,255,255,.07);
}

.impact-card__stats svg {
  color: var(--lime-300);
}

.impact-card__stats strong {
  font-size: 1.45rem;
}

.impact-card__stats span {
  color: rgba(255,255,255,.65);
  font-size: .77rem;
  line-height: 1.35;
}

@media (max-width: 860px) {
  .impact-card {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 620px) {
  .impact-card {
    gap: 25px;
    padding: 25px 20px;
    border-radius: 24px;
  }

  .impact-card__stats {
    grid-template-columns: 1fr;
  }

  .impact-card__stats > div {
    min-height: auto;
    grid-template-columns: auto 1fr;
    align-items: center;
  }

  .impact-card__stats span {
    grid-column: 2;
  }
}
</style>
