<template>
  <Dialog v-model:visible="visibleState" header="Tambah Baris Data" :modal="true" class="w-full max-w-lg bg-zinc-900 border border-white/10 rounded-lg p-6 text-white">
    <div class="space-y-4 max-h-[400px] overflow-y-auto pr-2">
      <template v-for="col in columns" :key="col.name">
        <div v-if="!isAutoIncrement(col)" class="flex flex-col gap-2">
          <label class="text-sm font-semibold">{{ col.name }} <span class="text-xs text-zinc-500">({{ col.type }})</span></label>
          
          <!-- If Boolean -->
          <select v-if="col.type.toLowerCase() === 'boolean'" v-model="form[col.name]" class="w-full bg-zinc-800 border border-white/10 text-white rounded p-2 text-sm focus:outline-none focus:border-zinc-500">
            <option value="">(NULL / Default)</option>
            <option value="true">True</option>
            <option value="false">False</option>
          </select>
          
          <!-- If Numeric -->
          <input v-else-if="col.type.toLowerCase().includes('int') || col.type.toLowerCase().includes('real') || col.type.toLowerCase().includes('numeric') || col.type.toLowerCase().includes('double') || col.type.toLowerCase().includes('precision')"
                 type="number" step="any" v-model="form[col.name]" placeholder="0" class="w-full bg-white/5 border border-white/10 rounded p-2 text-zinc-100 text-sm focus:outline-none focus:border-zinc-500" />
          
          <!-- If Date -->
          <input v-else-if="col.type.toLowerCase() === 'date'" type="date" v-model="form[col.name]" class="w-full bg-white/5 border border-white/10 rounded p-2 text-zinc-100 text-sm focus:outline-none focus:border-zinc-500" />
          
          <!-- If Timestamp / Datetime -->
          <input v-else-if="col.type.toLowerCase().includes('time')" type="datetime-local" v-model="form[col.name]" class="w-full bg-white/5 border border-white/10 rounded p-2 text-zinc-100 text-sm focus:outline-none focus:border-zinc-500" />
          
          <!-- If Text -->
          <textarea v-else-if="col.type.toLowerCase() === 'text'" v-model="form[col.name]" rows="3" placeholder="Isi teks..." class="w-full bg-white/5 border border-white/10 rounded p-2 text-zinc-100 text-sm focus:outline-none focus:border-zinc-500"></textarea>
          
          <!-- Default (Varchar/etc) -->
          <InputText v-else v-model="form[col.name]" :placeholder="'Isi ' + col.name" class="w-full bg-white/5 border border-white/10 text-white" />
        </div>
      </template>
      <div class="flex gap-2 justify-end mt-4">
        <Button label="Batal" severity="secondary" variant="outlined" @click="visibleState = false" />
        <Button label="Simpan Data" @click="submit" />
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
  },
  columns: {
    type: Array,
    default: () => []
  }
});

const emit = defineEmits(['update:visible', 'save']);

const visibleState = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val)
});

const form = ref({});

const isAutoIncrement = (col) => {
  if (!col) return false;
  const type = (col.type || '').toLowerCase();
  const def = (col.default || '').toLowerCase();
  return type.includes('serial') || def.includes('nextval') || (col.name === 'id' && type.includes('int') && def.includes('nextval'));
};

watch(() => props.visible, (newVal) => {
  if (newVal) {
    const data = {};
    props.columns.forEach(col => {
      if (!isAutoIncrement(col)) {
        data[col.name] = '';
      }
    });
    form.value = data;
  }
});

const submit = () => {
  emit('save', { ...form.value });
};
</script>
