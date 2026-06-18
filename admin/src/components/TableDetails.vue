<template>
  <div class="space-y-6">
    <div class="space-y-4 max-w-2xl bg-zinc-900/50 p-6 rounded-lg border border-white/5">
      <p class="text-lg font-semibold">Table Details</p>
      <IftaLabel class="w-full">
        <InputText class="w-full" v-model="activeTable.newNameInput" id="table_name_input" />
        <label for="table_name_input">Table Name</label>
      </IftaLabel>
      <div class="flex gap-4">
        <Button class="flex-1" :disabled="activeTable.newNameInput === activeTable.name" :variant="activeTable.newNameInput !== activeTable.name ? 'default' : 'outlined'" @click="$emit('rename-table')">Rename Tabel</Button>
        <Button severity="danger" variant="outlined" class="flex-1" @click="$emit('drop-table')">Drop Tabel</Button>
      </div>
    </div>

    <!-- Columns -->
    <div class="space-y-4">
      <div class="flex items-center justify-between">
        <p class="text-lg font-semibold flex items-center gap-2">
          <span>Columns</span>
          <span class="text-xs text-zinc-500">({{ activeTable.columns?.length || 0 }})</span>
        </p>
        <div class="flex gap-2">
          <Button size="small" :variant="'outlined'" @click="$emit('add-column-click')">
            <i class="pi pi-plus"></i>
            <span>Tambah Kolom</span>
          </Button>
          <Button size="small" severity="danger" variant="outlined"
            @click="$emit('delete-columns')">
            <i class="pi pi-trash"></i>
            <span>Hapus Terpilih</span>
          </Button>
        </div>
      </div>

      <div class="border border-white/10 rounded-md overflow-hidden bg-zinc-950/40">
        <table class="w-full text-xs text-zinc-300 border-collapse">
          <thead>
            <tr class="bg-white/5 border-b border-white/10 text-left text-white font-medium">
              <th class="p-4 w-[3rem]">Pilih</th>
              <th class="p-4">Name</th>
              <th class="p-4">Type</th>
              <th class="p-4">Nullable</th>
              <th class="p-4 w-[6rem] text-center">Edit</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!activeTable.columns || activeTable.columns.length === 0">
              <td colspan="5" class="p-8 text-center text-zinc-500">Belum ada kolom.</td>
            </tr>
            <tr v-for="column in activeTable.columns" :key="column.name"
              class="border-b border-white/5 hover:bg-white/5">
              <td class="p-4">
                <Checkbox :binary="true" v-model="column.checked" />
              </td>
              <td class="p-4 font-mono text-zinc-100">{{ column.name }}</td>
              <td class="p-4"><span class="bg-zinc-800 text-zinc-400 px-2 py-1 rounded text-[10px] uppercase">{{
                  column.type }}</span></td>
              <td class="p-4">{{ column.nullable }}</td>
              <td class="p-4 text-center">
                <Button variant="text" severity="warn" size="small" icon="pi pi-pencil"
                  @click="$emit('edit-column-click', column)" />
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Rows -->
    <div class="space-y-4">
      <div class="flex items-center justify-between">
        <p class="text-lg font-semibold flex items-center gap-2">
          <span>Rows</span>
          <span class="text-xs text-zinc-500">({{ activeTable.rows?.length || 0 }})</span>
        </p>
        <div class="flex gap-2">
          <Button size="small" :variant="'outlined'" @click="$emit('add-row-click')">
            <i class="pi pi-plus"></i>
            <span>Add Row</span>
          </Button>
          <Button size="small" severity="danger" :disabled="activeTable.rows?.filter(r => r._selected).length === 0"
            :variant="activeTable.rows?.filter(r => r._selected).length === 0 ? 'outlined' : 'default'"
            @click="$emit('delete-rows')">
            <i class="pi pi-trash"></i>
            <span>Hapus Terpilih</span>
          </Button>
        </div>
      </div>

      <div class="border border-white/10 rounded-md overflow-x-auto bg-zinc-950/40">
        <table class="w-full text-xs text-zinc-300 border-collapse">
          <thead>
            <tr class="bg-white/5 border-b border-white/10 text-left text-white font-medium">
              <th class="p-4 w-[3rem]">Pilih</th>
              <th v-for="col in activeTable.columns" :key="col.name" class="p-4">{{ col.name }}</th>
              <th class="p-4 w-[6rem] text-center">Edit</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!activeTable.rows || activeTable.rows.length === 0">
              <td :colspan="(activeTable.columns?.length || 0) + 2" class="p-8 text-center text-zinc-500">Tabel ini
                tidak memiliki data.</td>
            </tr>
            <tr v-for="(row, rowIndex) in activeTable.rows" :key="rowIndex"
              class="border-b border-white/5 hover:bg-white/5">
              <td class="p-4">
                <Checkbox :binary="true" v-model="row._selected" />
              </td>
              <td v-for="col in activeTable.columns" :key="col.name" class="p-4 font-mono">
                {{ row[col.name] !== null ? row[col.name] : 'NULL' }}
              </td>
              <td class="p-4 text-center">
                <Button variant="text" severity="warn" size="small" icon="pi pi-pencil"
                  @click="$emit('edit-row-click', row)" />
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import IftaLabel from 'primevue/iftalabel';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import Checkbox from 'primevue/checkbox';

defineProps({
  activeDatabase: {
    type: Object,
    required: true
  },
  activeTable: {
    type: Object,
    required: true
  }
});

defineEmits([
  'rename-table',
  'drop-table',
  'add-column-click',
  'edit-column-click',
  'delete-columns',
  'add-row-click',
  'edit-row-click',
  'delete-rows'
]);
</script>
