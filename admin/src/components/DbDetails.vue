<template>
  <div class="w-full">
    <!-- Database renaming panel -->
    <div v-if="!activeTable" class="space-y-4 max-w-2xl bg-zinc-900/50 p-6 rounded-lg border border-white/5 backdrop-blur-sm">
      <p class="text-lg font-semibold">Database Details</p>
      <IftaLabel class="w-full">
        <InputText class="w-full" v-model="activeDatabase.newNameInput" id="db_name_input" />
        <label for="db_name_input">Database Name</label>
      </IftaLabel>
      <div class="flex gap-4">
        <Button class="flex-1" :disabled="activeDatabase.newNameInput === activeDatabase.name" :variant="activeDatabase.newNameInput !== activeDatabase.name ? 'default' : 'outlined'" @click="$emit('rename-db')">Simpan Perubahan</Button>
        <Button severity="danger" variant="outlined" class="flex-1" @click="$emit('drop-db')">Drop Database</Button>
      </div>
    </div>

    <!-- Active DB path trace -->
    <div v-if="activeTable" class="p-2 mb-4 bg-zinc-900/30 rounded border border-white/5">
      <p class="text-sm text-zinc-400 flex items-center gap-2">
        <i class="pi pi-database"></i> Database: <strong class="text-white">{{ activeDatabase.name }}</strong>
      </p>
    </div>

    <!-- Tables navigation -->
    <div class="space-y-3">
      <div class="flex items-center justify-between">
        <p class="text-lg font-semibold">Tables</p>
        <Button size="small" class="flex items-center gap-1" @click="$emit('create-table-click')">
          <i class="pi pi-plus"></i>
          <span>Tabel Baru</span>
        </Button>
      </div>
      
      <ul class="flex flex-wrap gap-2 p-2 border border-white/5 rounded-md bg-zinc-900/20 min-h-[50px] items-center">
        <li v-if="!activeDatabase.tables || activeDatabase.tables.length === 0" class="text-sm text-zinc-500 w-full text-center py-2">
          Belum ada tabel di database ini.
        </li>
        <li v-for="table in activeDatabase.tables" :key="table.name"
            :class="['p-2 px-6 cursor-pointer hover:outline outline-white/5 rounded-md flex items-center gap-2 transition-all duration-200', activeTable?.name === table.name ? 'bg-white/5 outline outline-white/10 font-semibold' : '']"
            @click="$emit('select-table', table)">
          <i class="pi pi-table text-zinc-400"></i>
          {{ table.name }}
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import IftaLabel from 'primevue/iftalabel';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';

defineProps({
  activeDatabase: {
    type: Object,
    required: true
  },
  activeTable: {
    type: Object,
    default: null
  }
});

defineEmits(['rename-db', 'drop-db', 'select-table', 'create-table-click']);
</script>
