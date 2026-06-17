<template>
  <div class="w-[260px] bg-[#121826] border border-slate-800 rounded-xl flex flex-col p-4 shrink-0">
    <div class="flex flex-col gap-3 mb-4">
      <span class="text-base font-bold text-slate-100">Databases</span>
      <Button 
        icon="pi pi-plus" 
        label="Database Baru" 
        class="p-button-sm p-button-primary"
        @click="$emit('create-click')"
      />
    </div>

    <div class="relative mb-4">
      <i class="pi pi-search absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-500 text-xs z-[1]" />
      <InputText 
        v-model="searchQuery" 
        placeholder="Cari database..." 
        class="w-full pl-8! !bg-[#0c0f19] !border-slate-800 !text-slate-100 p-inputtext-sm" 
      />
    </div>

    <div class="flex-1 overflow-y-auto flex flex-col gap-1.5">
      <div 
        v-for="db in filteredDatabases" 
        :key="db" 
        :class="[
          'flex items-center gap-2.5 p-2.5 rounded-lg cursor-pointer text-sm font-medium transition-all duration-200',
          selectedDb === db ? 'bg-indigo-500/10 text-indigo-400' : 'text-slate-400 hover:bg-slate-800/60 hover:text-slate-100'
        ]"
        @click="$emit('select', db)"
      >
        <i class="pi pi-database text-sm" />
        <span class="truncate">{{ db }}</span>
      </div>

      <div v-if="filteredDatabases.length === 0" class="text-xs text-slate-500 text-center pt-5">
        Tidak ada database
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'

const props = defineProps({
  databases: {
    type: Array,
    required: true
  },
  selectedDb: {
    type: String,
    default: null
  }
})

defineEmits(['select', 'create-click'])

const searchQuery = ref('')

const filteredDatabases = computed(() => {
  if (!searchQuery.value) return props.databases
  return props.databases.filter(db => 
    db.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})
</script>
