<template>
  <Dialog 
    :visible="visible" 
    header="Ubah Nama Tabel" 
    :modal="true" 
    :style="{ width: '400px' }"
    @update:visible="$emit('update:visible', $event)"
  >
    <div class="flex flex-col gap-4 pt-2">
      <div class="flex flex-col gap-1.5">
        <label for="rename_table_name" class="text-xs font-bold text-slate-400">Nama Tabel Baru</label>
        <InputText id="rename_table_name" v-model="newName" placeholder="Nama baru..." required class="!bg-[#0c0f19] !border-slate-800 !text-slate-100" />
      </div>
    </div>
    <template #footer>
      <Button label="Batal" class="p-button-text p-button-secondary" @click="$emit('update:visible', false)" />
      <Button label="Simpan Nama" class="p-button-primary" @click="submit" :loading="loading" />
    </template>
  </Dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'

const props = defineProps({
  visible: {
    type: Boolean,
    required: true
  },
  currentName: {
    type: String,
    default: ''
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:visible', 'submit'])

const newName = ref('')

watch(() => props.visible, (newVal) => {
  if (newVal) {
    newName.value = props.currentName
  }
})

const submit = () => {
  emit('submit', newName.value)
}
</script>
