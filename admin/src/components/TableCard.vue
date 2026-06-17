<template>
  <Card class="!bg-[#121826] !border-slate-800 !rounded-lg hover:-translate-y-0.5 hover:!border-indigo-500 hover:shadow-lg transition-all duration-200">
    <template #content>
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2.5 font-semibold text-sm text-slate-100 truncate min-w-0">
          <i class="pi pi-table text-sky-400 text-sm shrink-0" />
          <span>{{ name }}</span>
        </div>
        <div class="flex gap-0.5 shrink-0">
          <Button 
            icon="pi pi-pencil" 
            class="p-button-text p-button-secondary p-button-sm" 
            title="Ubah Nama Tabel"
            @click="$emit('rename', name)"
          />
          <Button 
            icon="pi pi-trash" 
            class="p-button-text p-button-danger p-button-sm" 
            title="Hapus Tabel"
            @click="$emit('delete', name)"
          />
        </div>
      </div>
    </template>
  </Card>
</template>

<script setup>
import Card from 'primevue/card'
import Button from 'primevue/button'

defineProps({
  name: {
    type: String,
    required: true
  }
})

defineEmits(['rename', 'delete'])
</script>
