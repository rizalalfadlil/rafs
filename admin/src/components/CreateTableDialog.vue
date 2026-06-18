<template>
  <Dialog v-model:visible="visibleState" header="Buat Tabel Baru" :modal="true" class="w-full max-w-xl bg-zinc-900 border border-white/10 rounded-lg p-6 text-white">
    <div class="space-y-4">
      <div class="flex flex-col gap-2">
        <label class="text-sm font-semibold">Nama Tabel</label>
        <InputText v-model="form.table_name" placeholder="contoh: users" class="w-full bg-white/5 border border-white/10 text-white" />
      </div>
      <div class="space-y-2">
        <div class="flex justify-between items-center">
          <label class="text-sm font-semibold">Kolom Tabel</label>
        </div>
        <div class="space-y-2 max-h-[200px] overflow-y-auto pr-1">
          <div v-for="(col, index) in form.columns" :key="index" class="flex gap-2 items-center">
            <InputText v-model="col.name" placeholder="Nama kolom" class="flex-1 bg-white/5 border border-white/10 text-white" />
            <select v-model="col.type" class="bg-zinc-800 border border-white/10 text-white rounded p-2 text-sm max-w-[180px]">
              <option value="SERIAL PRIMARY KEY">SERIAL PRIMARY KEY</option>
              <option value="INTEGER">INTEGER</option>
              <option value="VARCHAR(255)">VARCHAR(255)</option>
              <option value="TEXT">TEXT</option>
              <option value="REAL">REAL</option>
              <option value="BOOLEAN">BOOLEAN</option>
              <option value="TIMESTAMP">TIMESTAMP</option>
            </select>
            <Button icon="pi pi-trash" severity="danger" variant="text" @click="removeColumn(index)" v-if="form.columns.length > 1" />
          </div>
        </div>
        <Button label="Tambah Kolom" icon="pi pi-plus" severity="secondary" variant="outlined" @click="addColumn" class="w-full mt-2" />
      </div>
      <div class="flex gap-2 justify-end mt-4">
        <Button label="Batal" severity="secondary" variant="outlined" @click="visibleState = false" />
        <Button label="Buat Tabel" @click="submit" />
      </div>
    </div>
  </Dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';

const props = defineProps({
  visible: {
    type: Boolean,
    required: true
  }
});

const emit = defineEmits(['update:visible', 'save']);

const visibleState = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val)
});

const form = ref({
  table_name: '',
  columns: [{ name: 'id', type: 'SERIAL PRIMARY KEY' }]
});

watch(() => props.visible, (newVal) => {
  if (newVal) {
    form.value = {
      table_name: '',
      columns: [{ name: 'id', type: 'SERIAL PRIMARY KEY' }]
    };
  }
});

const addColumn = () => {
  form.value.columns.push({ name: '', type: 'VARCHAR(255)' });
};

const removeColumn = (index) => {
  form.value.columns.splice(index, 1);
};

const submit = () => {
  emit('save', { ...form.value });
};
</script>
