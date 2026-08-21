<script setup lang="ts">
import { categories } from '~/data/marketplace'

defineProps<{
  modelValue: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

</script>

<template>
  <div
    class="scroll-row"
    role="list"
    aria-label="Категории"
  >
    <button
      v-for="category in categories"
      :key="category.id"
      type="button"
      class="category"
      :class="{ 'category--active': modelValue === category.id }"
      @click="emit('update:modelValue', category.id)"
    >
      <span class="category__emoji" aria-hidden="true">{{ category.emoji }}</span>
      <span>{{ category.label }}</span>
    </button>
  </div>
</template>

<style scoped>
.category {
  display: flex;
  flex: 0 0 auto;
  min-height: 48px;
  align-items: center;
  gap: 8px;
  padding: 0 17px 0 10px;
  border: 1px solid var(--sand-200);
  border-radius: 999px;
  color: var(--ink-700);
  background: var(--white);
  cursor: pointer;
  font-weight: 800;
  white-space: nowrap;
  transition: color 160ms ease, background 160ms ease, transform 160ms ease;
}

.category:hover {
  transform: translateY(-2px);
}

.category--active {
  border-color: var(--forest-900);
  color: var(--white);
  background: var(--forest-900);
}

.category__emoji {
  display: grid;
  width: 32px;
  height: 32px;
  place-items: center;
  border-radius: 50%;
  background: var(--cream-100);
  font-size: 1.04rem;
}

</style>
