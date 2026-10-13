<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { activityApi, foldersApi, orgsApi } from '@/api'
import type { ActivityEvent, Folder, OrgMember, Organization, Repository } from '@/api/types'
import ActivityFeed from '@/components/ActivityFeed.vue'
import ContentsList from '@/components/ContentsList.vue'
import CreateFolderModal from '@/components/CreateFolderModal.vue'
import CreateRepoModal from '@/components/CreateRepoModal.vue'
import Icon from '@/components/ui/Icon.vue'
import PageState from '@/components/ui/PageState.vue'
import { useBreadcrumbs } from '@/composables/useBreadcrumbs'
import { i18n, localizeError } from '@/i18n'

const props = defineProps<{ slug: string }>()
const router = useRouter()
const { setBreadcrumbs, clearBreadcrumbs } = useBreadcrumbs()

const org = ref<Organization | null>(null)
const repos = ref<Repository[]>([])
const folders = ref<Folder[]>([])
const activity = ref<ActivityEvent[]>([])
const activityError = ref('')
const activityLoading = ref(false)
const members = ref<OrgMember[]>([])
const tab = ref<'repos' | 'activity' | 'members' | 'settings'>('repos')
const loading = ref(true)
const error = ref('')
const showRepoModal = ref(false)
const showFolderModal = ref(false)
const memberUsername = ref('')
const memberRole = ref('member')
const busy = ref(false)
const name = ref('')
const description = ref('')
const saved = ref(false)

const canAdmin = computed(() => org.value?.role === 'owner' || org.value?.role === 'admin')
const canOwner = computed(() => org.value?.role === 'owner')
const rootRepos = computed(() => repos.value.filter((r) => r.folder_id == null))
const initial = computed(() => (org.value?.name || org.value?.slug || '?').slice(0, 1).toUpperCase())

function roleLabel(role: string) {
  const key = `orgs.role_${role}` as 'orgs.role_owner' | 'orgs.role_admin' | 'orgs.role_member'
  return i18n.global.te(key) ? i18n.global.t(key) : role
}

async function loadActivity() {
  activityLoading.value = true
  activityError.value = ''
  try {
    const res = await activityApi.listOwner(props.slug)
    activity.value = res.events ?? []
  } catch (err) {
    activityError.value = localizeError(err, 'errors.loadActivity')
  } finally {
    activityLoading.value = false
  }
}

async function load() {
  loading.value = true
  error.value = ''
  saved.value = false
  try {
    org.value = await orgsApi.get(props.slug)
    name.value = org.value.name
    description.value = org.value.description
    setBreadcrumbs([
      {
        label: 'orgs',
        labelKey: 'nav.orgs',
        to: { 
          name: 'orgs' 
        },
      },
      { 
        label: org.value.slug 
      },
    ])
    const [repoRes, memberRes, folderRes] = await Promise.all([
      orgsApi.listRepos(props.slug),
      org.value.role
        ? orgsApi.listMembers(props.slug)
        : Promise.resolve({ members: [] as OrgMember[] }),
      org.value.role
        ? foldersApi.list(null, props.slug)
        : Promise.resolve({ folders: [] as Folder[] }),
    ])
    repos.value = repoRes.repos ?? []
    members.value = memberRes.members ?? []
    folders.value = folderRes.folders ?? []
  } catch (err) {
    error.value = localizeError(err, 'errors.loadOrg')
  } finally {
    loading.value = false
  }
}

async function addMember() {
  busy.value = true
  error.value = ''
  try {
    await orgsApi.addMember(props.slug, {
      username: memberUsername.value.trim(),
      role: memberRole.value,
    })
    memberUsername.value = ''
    await load()
  } catch (err) {
    error.value = localizeError(err, 'errors.orgMember')
  } finally {
    busy.value = false
  }
}

async function saveSettings() {
  busy.value = true
  error.value = ''
  saved.value = false
  try {
    org.value = await orgsApi.update(props.slug, {
      name: name.value.trim(),
      description: description.value.trim(),
    })
    saved.value = true
  } catch (err) {
    error.value = localizeError(err, 'errors.saveOrg')
  } finally {
    busy.value = false
  }
}

async function removeOrg() {
  if (!confirm(i18n.global.t('orgs.confirmDelete', { slug: props.slug }))) {
    return
  }
  busy.value = true
  try {
    await orgsApi.remove(props.slug)
    await router.push({ name: 'orgs' })
  } catch (err) {
    error.value = localizeError(err, 'errors.deleteOrg')
  } finally {
    busy.value = false
  }
}

watch(tab, (t) => {
  if (t === 'activity' && org.value) {
    loadActivity()
  }
})

watch(() => props.slug, () => {
  tab.value = 'repos'
  load()
})
onMounted(load)
onUnmounted(clearBreadcrumbs)
</script>

<template>
  <div>
    <PageState :loading="loading" :error="error && !org ? error : ''">
      <template v-if="org">
        <div class="border-b border-line bg-white">
          <div class="mx-auto max-w-7xl px-4 pt-4">
            <div class="flex flex-wrap items-start justify-between gap-4">
              <div class="flex min-w-0 items-start gap-3">
                <span
                  class="flex h-12 w-12 shrink-0 items-center justify-center rounded-md bg-accent-soft text-lg font-semibold text-accent"
                >
                  {{ initial }}
                </span>
                <div class="min-w-0">
                  <div class="flex flex-wrap items-center gap-2">
                    <h1 class="truncate text-xl font-semibold">{{ org.name }}</h1>
                    <span
                      v-if="org.role"
                      class="rounded-full border border-line bg-paper px-2 py-0.5 text-xs text-ink-muted"
                    >
                      {{ roleLabel(org.role) }}
                    </span>
                  </div>
                  <p class="mt-0.5 font-mono text-sm text-ink-muted">/{{ org.slug }}</p>
                  <p v-if="org.description" class="mt-2 text-sm text-ink-muted">{{ org.description }}</p>
                </div>
              </div>
              <div v-if="canAdmin && tab === 'repos'" class="flex flex-wrap gap-2">
                <button type="button" class="btn-ghost" @click="showFolderModal = true">
                  {{ $t('home.newFolder') }}
                </button>
                <button type="button" class="btn-primary" @click="showRepoModal = true">
                  {{ $t('home.newRepository') }}
                </button>
              </div>
            </div>

            <nav class="mt-4 flex gap-1 overflow-x-auto text-sm" :aria-label="$t('orgs.tabs')">
              <button
                type="button"
                class="inline-flex items-center gap-2 border-b-2 px-3 py-2"
                :class="tab === 'repos' ? 'border-[#fd8c73] font-semibold' : 'border-transparent text-ink hover:border-line'"
                @click="tab = 'repos'"
              >
                <Icon name="repo" />
                {{ $t('orgs.repos') }}
                <span class="text-ink-muted">{{ rootRepos.length + folders.length }}</span>
              </button>
              <button
                type="button"
                class="inline-flex items-center gap-2 border-b-2 px-3 py-2"
                :class="tab === 'activity' ? 'border-[#fd8c73] font-semibold' : 'border-transparent text-ink hover:border-line'"
                @click="tab = 'activity'"
              >
                {{ $t('orgs.activity') }}
              </button>
              <button
                v-if="org.role"
                type="button"
                class="inline-flex items-center gap-2 border-b-2 px-3 py-2"
                :class="tab === 'members' ? 'border-[#fd8c73] font-semibold' : 'border-transparent text-ink hover:border-line'"
                @click="tab = 'members'"
              >
                {{ $t('orgs.members') }}
                <span class="text-ink-muted">{{ members.length }}</span>
              </button>
              <button
                v-if="canAdmin"
                type="button"
                class="inline-flex items-center gap-2 border-b-2 px-3 py-2"
                :class="tab === 'settings' ? 'border-[#fd8c73] font-semibold' : 'border-transparent text-ink hover:border-line'"
                @click="tab = 'settings'"
              >
                <Icon name="gear" />
                {{ $t('orgs.settings') }}
              </button>
            </nav>
          </div>
        </div>

        <div class="mx-auto max-w-7xl px-4 py-4">
          <p v-if="error" class="mb-3 text-sm text-warn">{{ error }}</p>

          <div v-if="tab === 'repos'">
            <ContentsList
              v-if="folders.length || rootRepos.length"
              :folders="folders"
              :repos="rootRepos"
              :repo-owner="slug"
              :organization="slug"
            />
            <div
              v-else
              class="rounded-md border border-dashed border-line bg-white px-4 py-12 text-center"
            >
              <p class="text-sm text-ink-muted">{{ $t('orgs.noRepos') }}</p>
              <div v-if="canAdmin" class="mt-4 flex flex-wrap justify-center gap-2">
                <button type="button" class="btn-ghost" @click="showFolderModal = true">
                  {{ $t('home.newFolder') }}
                </button>
                <button type="button" class="btn-primary" @click="showRepoModal = true">
                  {{ $t('home.newRepository') }}
                </button>
              </div>
            </div>
          </div>

          <div v-else-if="tab === 'activity'">
            <ActivityFeed
              :events="activity"
              :loading="activityLoading"
              :error="activityError"
            />
          </div>

          <div v-else-if="tab === 'members'" class="space-y-4">
            <form
              v-if="canAdmin"
              class="flex flex-wrap items-end gap-3 rounded-md border border-line bg-white p-4"
              @submit.prevent="addMember"
            >
              <label class="block min-w-48 flex-1 text-sm">
                <span class="mb-1 block font-semibold">{{ $t('orgs.username') }}</span>
                <input
                  v-model="memberUsername"
                  class="input"
                  required
                  :placeholder="$t('orgs.usernamePlaceholder')"
                >
              </label>
              <label class="block text-sm">
                <span class="mb-1 block font-semibold">{{ $t('orgs.role') }}</span>
                <select v-model="memberRole" class="input min-w-32">
                  <option value="member">{{ $t('orgs.role_member') }}</option>
                  <option value="admin">{{ $t('orgs.role_admin') }}</option>
                  <option v-if="canOwner" value="owner">{{ $t('orgs.role_owner') }}</option>
                </select>
              </label>
              <button type="submit" class="btn-primary" :disabled="busy">
                {{ $t('orgs.addMember') }}
              </button>
            </form>

            <section class="overflow-hidden rounded-md border border-line bg-white">
              <h2 class="border-b border-line bg-paper px-4 py-2 text-sm font-semibold">
                {{ $t('orgs.members') }}
              </h2>
              <ul v-if="members.length" class="divide-y divide-line">
                <li
                  v-for="m in members"
                  :key="m.user_id"
                  class="flex items-center justify-between gap-3 px-4 py-3"
                >
                  <div class="flex min-w-0 items-center gap-3">
                    <span
                      class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-paper-2 text-xs font-semibold text-ink-muted"
                    >
                      {{ m.username.slice(0, 1).toUpperCase() }}
                    </span>
                    <span class="truncate font-semibold">{{ m.username }}</span>
                  </div>
                  <span class="rounded-full border border-line bg-paper px-2 py-0.5 text-xs text-ink-muted">
                    {{ roleLabel(m.role) }}
                  </span>
                </li>
              </ul>
              <p v-else class="px-4 py-8 text-center text-sm text-ink-muted">
                {{ $t('orgs.noMembers') }}
              </p>
            </section>
          </div>

          <div v-else-if="tab === 'settings'" class="max-w-2xl space-y-6">
            <form
              class="space-y-4 rounded-md border border-line bg-white p-4"
              @submit.prevent="saveSettings"
            >
              <h2 class="text-lg font-semibold">{{ $t('orgs.settings') }}</h2>
              <label class="block text-sm">
                <span class="mb-1 block font-semibold">{{ $t('orgs.name') }}</span>
                <input v-model="name" class="input" required maxlength="100">
              </label>
              <label class="block text-sm">
                <span class="mb-1 block font-semibold">{{ $t('orgs.description') }}</span>
                <input v-model="description" class="input" :placeholder="$t('modal.optional')">
              </label>
              <p v-if="saved" class="text-sm text-moss">{{ $t('orgs.saved') }}</p>
              <button type="submit" class="btn-primary" :disabled="busy">
                {{ busy ? $t('settings.saving') : $t('orgs.save') }}
              </button>
            </form>

            <div v-if="canOwner" class="overflow-hidden rounded-md border border-warn">
              <h2 class="border-b border-warn px-4 py-2 text-sm font-semibold">
                {{ $t('settings.danger') }}
              </h2>
              <div class="flex flex-wrap items-center justify-between gap-3 px-4 py-3">
                <div>
                  <p class="font-semibold">{{ $t('orgs.delete') }}</p>
                  <p class="text-sm text-ink-muted">{{ $t('orgs.danger') }}</p>
                </div>
                <button
                  type="button"
                  class="btn-danger"
                  :disabled="busy"
                  @click="removeOrg"
                >
                  {{ $t('orgs.delete') }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </template>
    </PageState>

    <CreateFolderModal
      v-if="showFolderModal"
      :organization="slug"
      @close="showFolderModal = false"
      @created="load"
    />
    <CreateRepoModal
      v-if="showRepoModal"
      :organization="slug"
      @close="showRepoModal = false"
      @created="load"
    />
  </div>
</template>
