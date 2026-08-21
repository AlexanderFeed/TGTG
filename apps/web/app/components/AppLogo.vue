<script setup lang="ts">
defineProps<{
  inverse?: boolean
  compact?: boolean
}>()
</script>

<template>
  <NuxtLink class="brand" :class="{ 'brand--inverse': inverse }" to="/" aria-label="ЕщёЕсть — на главную">
    <span class="brand__mark" aria-hidden="true">
      <span class="brand__leaf" />
      <span class="brand__stem" />
    </span>
    <span v-if="!compact" class="brand__name">ЕщёЕсть</span>
  </NuxtLink>
</template>

<style scoped>
.brand {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  color: var(--forest-900);
  font-size: 1.16rem;
  font-weight: 950;
  letter-spacing: -0.05em;
}

.brand__mark {
  position: relative;
  display: block;
  width: 34px;
  height: 34px;
  overflow: hidden;
  border-radius: 11px;
  background: var(--forest-900);
}

.brand__leaf {
  position: absolute;
  top: 7px;
  left: 9px;
  width: 17px;
  height: 15px;
  border-radius: 12px 2px 12px 2px;
  background: var(--lime-300);
  transform: rotate(-12deg);
}

.brand__stem {
  position: absolute;
  top: 15px;
  left: 14px;
  width: 13px;
  height: 2px;
  border-radius: 2px;
  background: var(--white);
  transform: rotate(-42deg);
  transform-origin: left center;
}

.brand--inverse {
  color: var(--white);
}

.brand--inverse .brand__mark {
  background: var(--lime-300);
}

.brand--inverse .brand__leaf {
  background: var(--forest-900);
}

.brand--inverse .brand__stem {
  background: var(--forest-900);
}
</style>
