<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { usersApi } from '@/api'
import { useAuth } from '@/composables/useAuth'
import { useBreadcrumbs } from '@/composables/useBreadcrumbs'
import Icon from '@/components/ui/Icon.vue'

const auth = useAuth()
const { crumbs } = useBreadcrumbs()
const router = useRouter()
const searchQ = ref('')
const menuOpen = ref(false)
const menuRef = ref<HTMLElement | null>(null)

const username = computed(() => auth.user.value?.username || '')
const initial = computed(() => username.value.slice(0, 1).toUpperCase() || '?')
const avatarSrc = computed(() => {
  const u = auth.user.value
  if (!u?.has_avatar || !u.username) {
    return ''
  }

  return usersApi.avatarURL(u.username)
})

function toggleMenu() {
  menuOpen.value = !menuOpen.value
}

function closeMenu() {
  menuOpen.value = false
}

async function onLogout() {
  closeMenu()
  await auth.logout()
  await router.push({ 
    name: 'login' 
  })
}

async function onSearch() {
  const q = searchQ.value.trim()
  await router.push({
    name: 'search',
    query: q ? { q } : {},
  })
}

function onDocClick(ev: MouseEvent) {
  if (!menuOpen.value || !menuRef.value) {
    return
  }

  if (!menuRef.value.contains(ev.target as Node)) {
    closeMenu()
  }
}

function onKeydown(ev: KeyboardEvent) {
  if (ev.key === 'Escape') {
    closeMenu()
  }
}

onMounted(() => {
  document.addEventListener('click', onDocClick)
  document.addEventListener('keydown', onKeydown)
})

onUnmounted(() => {
  document.removeEventListener('click', onDocClick)
  document.removeEventListener('keydown', onKeydown)
})
</script>

<template>
  <header class="bg-header text-white">
    <div class="mx-auto flex h-14 max-w-7xl items-center gap-4 px-4">
      <RouterLink to="/" class="flex cursor-pointer items-center gap-2 font-semibold text-white">
        <Icon name="mark" />
        Gitwold
      </RouterLink>

      <nav
        v-if="crumbs.length"
        class="hidden min-w-0 flex-1 items-center gap-1.5 overflow-x-auto text-sm text-white/70 sm:flex"
        :aria-label="$t('nav.breadcrumb')"
      >
        <template
          v-for="(crumb, i) in crumbs"
          :key="`${crumb.labelKey || crumb.label}-${i}`"
        >
          <span v-if="i > 0">/</span>
          <RouterLink
            v-if="crumb.to"
            :to="crumb.to"
            class="cursor-pointer truncate text-white hover:underline"
          >
            {{ crumb.labelKey ? $t(crumb.labelKey) : crumb.label }}
          </RouterLink>
          <span v-else class="truncate text-white">{{ crumb.labelKey ? $t(crumb.labelKey) : crumb.label }}</span>
        </template>
      </nav>
      <div v-else class="flex-1" />

      <form
        class="hidden min-w-0 max-w-xs flex-1 sm:block"
        @submit.prevent="onSearch"
      >
        <input
          v-model="searchQ"
          type="search"
          class="w-full rounded-md border-0 bg-white/10 px-3 py-1.5 text-sm text-white placeholder:text-white/50 focus:bg-white/15 focus:outline-none focus:ring-1 focus:ring-white/40"
          :placeholder="$t('search.placeholder')"
          :aria-label="$t('search.title')"
        >
      </form>
      <RouterLink
        :to="{ name: 'orgs' }"
        class="hidden cursor-pointer text-sm text-white/80 hover:text-white sm:inline"
      >
        {{ $t('nav.orgs') }}
      </RouterLink>

      <div ref="menuRef" class="relative">
        <button
          type="button"
          class="flex cursor-pointer items-center gap-2 rounded-full p-0.5 text-sm text-white/90 hover:bg-white/10"
          :aria-expanded="menuOpen"
          aria-haspopup="menu"
          :aria-label="$t('nav.profile')"
          @click.stop="toggleMenu"
        >
          <span
            class="flex h-7 w-7 items-center justify-center overflow-hidden rounded-full bg-white/15 text-xs font-semibold"
          >
            <img
              v-if="avatarSrc"
              :src="avatarSrc"
              alt=""
              class="h-full w-full object-cover"
            >
            <span v-else>{{ initial }}</span>
          </span>
        </button>

        <div
          v-if="menuOpen"
          class="absolute right-0 z-50 mt-2 w-52 overflow-hidden rounded-md border border-line bg-white py-1 text-ink shadow-lg"
          role="menu"
        >
          <p class="border-b border-line px-3 py-2 text-xs text-ink-muted">
            {{ username }}
          </p>
          <RouterLink
            :to="{ name: 'profile' }"
            class="block cursor-pointer px-3 py-2 text-sm hover:bg-paper"
            role="menuitem"
            @click="closeMenu"
          >
            {{ $t('nav.profile') }}
          </RouterLink>
          <RouterLink
            :to="{ name: 'ssh-keys' }"
            class="block cursor-pointer px-3 py-2 text-sm hover:bg-paper"
            role="menuitem"
            @click="closeMenu"
          >
            {{ $t('nav.sshKeys') }}
          </RouterLink>
          <RouterLink
            :to="{ name: 'tokens' }"
            class="block cursor-pointer px-3 py-2 text-sm hover:bg-paper"
            role="menuitem"
            @click="closeMenu"
          >
            {{ $t('nav.accessTokens') }}
          </RouterLink>
          <div class="my-1 border-t border-line" />
          <button
            type="button"
            class="block w-full cursor-pointer px-3 py-2 text-left text-sm text-warn hover:bg-paper"
            role="menuitem"
            @click="onLogout"
          >
            {{ $t('nav.logout') }}
          </button>
        </div>
      </div>
    </div>
  </header>
</template>
