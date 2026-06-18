<template>
  <Dialog v-model:visible="visibleState" header="Ubah Baris Data" :modal="true" class="w-full max-w-lg bg-zinc-900 border border-white/10 rounded-lg p-6 text-white">
    <div class="space-y-4 max-h-[400px] overflow-y-auto pr-2">
      <div v-for="col in columns" :key="col.name" class="flex flex-col gap-2">
        <label class="text-sm font-semibold">{{ col.name }} <span class="text-xs text-zinc-500">({{ col.type }})</span></label>
        
        <!-- If Auto Increment or ID (Disabled/Read Only) -->
        <InputText v-if="isAutoIncrement(col) || col.name === 'id'" :value="form[col.name]" :disabled="true" class="w-full bg-white/5 border border-white/10 opacity-60 text-white" />
        
        <!-- If Boolean -->
        <select v-else-if="col.type.toLowerCase() === 'boolean'" v-model="form[col.name]" class="w-full bg-zinc-800 border border-white/10 text-white rounded p-2 text-sm focus:outline-none focus:border-zinc-500">
          <option value="">(NULL)</option>
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
      <div class="flex gap-2 justify-end mt-4">
        <Button label="Batal" severity="secondary" variant="outlined" @click="visibleState = false" />
        <Button label="Simpan Perubahan" @click="submit" />
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
  },
  row: {
    type: Object,
    default: () => ({})
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
  if (newVal && props.row) {
    const data = {};
    props.columns.forEach(col => {
      const val = props.row[col.name];
      if (col.type.toLowerCase() === 'boolean' && val !== null && val !== undefined) {
        data[col.name] = String(val);
      } else {
        data[col.name] = val !== null && val !== undefined ? val : '';
      }
    });
    form.value = data;
  }
});

const submit = () => {
  emit('save', { ...form.value });
};
</script>
