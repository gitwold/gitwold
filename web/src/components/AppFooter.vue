<script setup lang="ts">
import LocaleSwitch from '@/components/LocaleSwitch.vue'
</script>

<template>
  <footer class="mt-auto border-t border-line bg-white">
    <div class="mx-auto flex max-w-7xl flex-col gap-4 px-4 py-6 sm:flex-row sm:items-center sm:justify-between">
      <div class="flex flex-wrap items-center gap-x-4 gap-y-2 text-sm text-ink-muted">
        <span class="font-semibold text-ink">Repo</span>
      </div>
      <div class="flex flex-wrap items-center gap-3">
        <span class="text-xs text-ink-muted">{{ $t('nav.language') }}</span>
        <LocaleSwitch />
      </div>
    </div>
  </footer>
</template>
