<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { foldersApi, reposApi } from '@/api'
import type { FolderTreeNode } from '@/api/types'
import ModalForm from '@/components/ModalForm.vue'
import { localizeError } from '@/i18n'

const props = defineProps<{
  folderId?: number | null
  organization?: string
}>()

const emit = defineEmits<{
  close: []
  created: []
}>()

const name = ref('')
const description = ref('')
const isPrivate = ref(false)
const defaultBranch = ref('main')
const error = ref('')
const busy = ref(false)
const folderKey = ref('')
const folderOptions = ref<{ id: string; label: string }[]>([])

function flattenFolders(nodes: FolderTreeNode[], acc: { id: number; label: string }[] = []) {
  for (const node of nodes) {
    acc.push({ id: node.id, label: node.path })
    if (node.children?.length) {
      flattenFolders(node.children, acc)
    }
  }

  return acc
}

async function loadFolders() {
  if (!props.organization) {
    folderOptions.value = []
    return
  }
  
  try {
    const tree = await foldersApi.tree(props.organization)
    folderOptions.value = flattenFolders(tree.tree ?? []).map((f) => ({
      id: String(f.id),
      label: f.label,
    }))
  } catch {
    folderOptions.value = []
  }
}

watch(
  () => props.folderId,
  (id) => {
    folderKey.value = id == null ? '' : String(id)
  },
  { immediate: true },
)

onMounted(loadFolders)

async function submit() {
  error.value = ''
  busy.value = true
  try {
    const pickedFolder = props.organization
      ? (folderKey.value === '' ? null : Number(folderKey.value))
      : (props.folderId ?? null)
    await reposApi.create({
      name: name.value.trim(),
      description: description.value.trim(),
      folder_id: pickedFolder,
      organization: props.organization || undefined,
      private: isPrivate.value,
      default_branch: defaultBranch.value.trim() || 'main',
    })
    emit('created')
    emit('close')
  } catch (err) {
    error.value = localizeError(err, 'errors.createRepository')
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <ModalForm
    :title="$t('modal.newRepository')"
    :submit-label="busy ? $t('modal.creating') : $t('modal.createRepository')"
    @close="emit('close')"
    @submit="submit"
  >
    <label class="block text-sm">
      <span class="mb-1 block font-medium">{{ $t('modal.name') }}</span>
      <input
        v-model="name"
        class="input"
        required
        maxlength="100"
        placeholder="my-repo"
        autofocus
      />
    </label>
    <label class="block text-sm">
      <span class="mb-1 block font-medium">{{ $t('modal.description') }}</span>
      <input
        v-model="description"
        class="input"
        :placeholder="$t('modal.optional')"
      />
    </label>
    <label class="block text-sm">
      <span class="mb-1 block font-medium">{{ $t('modal.defaultBranch') }}</span>
      <input
        v-model="defaultBranch"
        class="input font-mono"
      />
    </label>
    <label v-if="organization && folderOptions.length" class="block text-sm">
      <span class="mb-1 block font-medium">{{ $t('modal.folder') }}</span>
      <select v-model="folderKey" class="input">
        <option value="">{{ $t('nav.root') }}</option>
        <option v-for="opt in folderOptions" :key="opt.id" :value="opt.id">
          {{ opt.label }}
        </option>
      </select>
    </label>
    <label class="flex items-center gap-2 text-sm">
      <input
        v-model="isPrivate"
        type="checkbox"
        class="accent-moss"
      />
      {{ $t('modal.private') }}
    </label>
    <p v-if="error" class="text-sm text-warn">{{ error }}</p>
  </ModalForm>
</template>
