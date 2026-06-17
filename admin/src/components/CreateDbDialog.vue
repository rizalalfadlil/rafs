<template>
  <Dialog 
    :visible="visible" 
    header="Buat Database & User Baru" 
    :modal="true" 
    :style="{ width: '450px' }"
    @update:visible="$emit('update:visible', $event)"
  >
    <div class="flex flex-col gap-4 pt-2">
      <div class="flex flex-col gap-1.5">
        <label for="db_name" class="text-xs font-bold text-slate-400">Nama Database</label>
        <InputText id="db_name" v-model="form.db_name" placeholder="contoh: toko_online" required class="!bg-[#0c0f19] !border-slate-800 !text-slate-100" />
        <small class="text-[10px] text-slate-500">Hanya alfanumerik dan underscore. Harus dimulai dengan huruf.</small>
      </div>
      <div class="flex flex-col gap-1.5">
        <label for="username" class="text-xs font-bold text-slate-400">Username Owner</label>
        <InputText id="username" v-model="form.username" placeholder="contoh: admin_toko" required class="!bg-[#0c0f19] !border-slate-800 !text-slate-100" />
      </div>
      <div class="flex flex-col gap-1.5">
        <label for="password" class="text-xs font-bold text-slate-400">Password Owner</label>
        <InputText id="password" type="password" v-model="form.password" placeholder="Password aman..." required class="!bg-[#0c0f19] !border-slate-800 !text-slate-100" />
      </div>
    </div>
    <template #footer>
      <Button label="Batal" class="p-button-text p-button-secondary" @click="$emit('update:visible', false)" />
      <Button label="Buat Database" class="p-button-primary" @click="submit" :loading="loading" />
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
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:visible', 'submit'])

const form = ref({ db_name: '', username: '', password: '' })

watch(() => props.visible, (newVal) => {
  if (newVal) {
    form.value = { db_name: '', username: '', password: '' }
  }
})

const submit = () => {
  emit('submit', { ...form.value })
}
</script>
