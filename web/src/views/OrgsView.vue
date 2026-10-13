<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { orgsApi } from '@/api'
import type { Organization } from '@/api/types'
import ModalForm from '@/components/ModalForm.vue'
import PageState from '@/components/ui/PageState.vue'
import { useBreadcrumbs } from '@/composables/useBreadcrumbs'
import { i18n, localizeError } from '@/i18n'

const router = useRouter()
const { setBreadcrumbs, clearBreadcrumbs } = useBreadcrumbs()

const orgs = ref<Organization[]>([])
const loading = ref(true)
const error = ref('')
const showForm = ref(false)
const slug = ref('')
const name = ref('')
const description = ref('')
const creating = ref(false)
const formError = ref('')

async function load() {
  loading.value = true
  error.value = ''
  try {
    setBreadcrumbs([{
      label: 'orgs',
      labelKey: 'nav.orgs',
    }])
    const res = await orgsApi.listMine()
    orgs.value = res.organizations ?? []
  } catch (err) {
    error.value = localizeError(err, 'errors.loadOrgs')
  } finally {
    loading.value = false
  }
}

async function create() {
  creating.value = true
  formError.value = ''
  try {
    const created = await orgsApi.create({
      slug: slug.value.trim(),
      name: name.value.trim(),
      description: description.value.trim(),
    })
    slug.value = ''
    name.value = ''
    description.value = ''
    showForm.value = false
    await router.push({ 
      name: 'org', 
      params: { 
        slug: created.slug 
      } 
    })
  } catch (err) {
    formError.value = localizeError(err, 'errors.createOrg')
  } finally {
    creating.value = false
  }
}

function roleLabel(role?: string) {
  if (!role) {
    return ''
  }

  const key = `orgs.role_${role}` as 'orgs.role_owner' | 'orgs.role_admin' | 'orgs.role_member'
  return i18n.global.te(key) ? i18n.global.t(key) : role
}

onMounted(load)
onUnmounted(clearBreadcrumbs)
</script>

<template>
  <div class="mx-auto max-w-3xl px-4 py-6">
    <div class="mb-6 flex flex-wrap items-end justify-between gap-3">
      <div>
        <h1 class="text-2xl font-semibold">{{ $t('orgs.title') }}</h1>
      </div>
      <button type="button" class="btn-primary" @click="showForm = true">
        {{ $t('orgs.new') }}
      </button>
    </div>

    <PageState :loading="loading" :error="error">
      <section class="overflow-hidden rounded-md border border-line bg-white">
        <h2 class="border-b border-line bg-paper px-4 py-2 text-sm font-semibold">
          {{ $t('orgs.yours') }}
        </h2>
        <ul v-if="orgs.length" class="divide-y divide-line">
          <li v-for="org in orgs" :key="org.id">
            <RouterLink
              :to="{ name: 'org', params: { slug: org.slug } }"
              class="flex items-center gap-3 px-4 py-3 hover:bg-paper"
            >
              <span
                class="flex h-10 w-10 shrink-0 items-center justify-center rounded-md bg-accent-soft text-sm font-semibold text-accent"
              >
                {{ org.name.slice(0, 1).toUpperCase() }}
              </span>
              <div class="min-w-0 flex-1">
                <p class="truncate font-semibold text-accent">{{ org.name }}</p>
                <p class="truncate font-mono text-xs text-ink-muted">/{{ org.slug }}</p>
                <p v-if="org.description" class="mt-0.5 truncate text-xs text-ink-muted">
                  {{ org.description }}
                </p>
              </div>
              <span
                v-if="org.role"
                class="shrink-0 rounded-full border border-line bg-paper px-2 py-0.5 text-xs text-ink-muted"
              >
                {{ roleLabel(org.role) }}
              </span>
            </RouterLink>
          </li>
        </ul>
        <div v-else class="px-4 py-12 text-center">
          <p class="text-sm text-ink-muted">{{ $t('orgs.empty') }}</p>
          <button type="button" class="btn-primary mt-4" @click="showForm = true">
            {{ $t('orgs.new') }}
          </button>
        </div>
      </section>
    </PageState>

    <ModalForm
      v-if="showForm"
      :title="$t('orgs.new')"
      :submit-label="creating ? $t('orgs.creating') : $t('orgs.create')"
      @close="showForm = false"
      @submit="create"
    >
      <label class="block text-sm">
        <span class="mb-1 block font-semibold">{{ $t('orgs.slug') }}</span>
        <input
          v-model="slug"
          class="input font-mono"
          required
          minlength="3"
          maxlength="39"
          pattern="[a-z0-9_-]+"
          placeholder="acme"
          autofocus
        >
        <span class="mt-1 block text-xs text-ink-muted">{{ $t('orgs.slugHint') }}</span>
      </label>
      <label class="block text-sm">
        <span class="mb-1 block font-semibold">{{ $t('orgs.name') }}</span>
        <input
          v-model="name"
          class="input"
          required
          maxlength="100"
          placeholder="Acme Inc"
        >
      </label>
      <label class="block text-sm">
        <span class="mb-1 block font-semibold">{{ $t('orgs.description') }}</span>
        <input
          v-model="description"
          class="input"
          :placeholder="$t('modal.optional')"
        >
      </label>
      <p v-if="formError" class="text-sm text-warn">{{ formError }}</p>
    </ModalForm>
  </div>
</template>
