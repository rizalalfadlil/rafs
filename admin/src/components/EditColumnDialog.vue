<template>
  <Dialog v-model:visible="visibleState" header="Ubah Kolom" :modal="true" class="w-full max-w-md bg-zinc-900 border border-white/10 rounded-lg p-6 text-white">
    <div class="space-y-4">
      <div class="flex flex-col gap-2">
        <label class="text-sm font-semibold">Nama Kolom Baru</label>
        <InputText v-model="form.new_name" placeholder="contoh: alamat" class="w-full bg-white/5 border border-white/10 text-white" />
      </div>
      <div class="flex flex-col gap-2">
        <label class="text-sm font-semibold">Ubah Tipe Data (Opsional)</label>
        <select v-model="form.type" class="bg-zinc-800 border border-white/10 text-white rounded p-2 text-sm w-full">
          <option value="">(Tidak diubah)</option>
          <option value="INTEGER">INTEGER</option>
          <option value="VARCHAR(255)">VARCHAR(255)</option>
          <option value="TEXT">TEXT</option>
          <option value="REAL">REAL</option>
          <option value="BOOLEAN">BOOLEAN</option>
          <option value="TIMESTAMP">TIMESTAMP</option>
        </select>
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
  column: {
    type: Object,
    required: true
  }
});

const emit = defineEmits(['update:visible', 'save']);

const visibleState = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val)
});

const form = ref({ old_name: '', new_name: '', type: '' });

watch(() => props.visible, (newVal) => {
  if (newVal && props.column) {
    form.value = {
      old_name: props.column.old_name || props.column.name,
      new_name: props.column.new_name || props.column.name,
      type: ''
    };
  }
});

const submit = () => {
  emit('save', { ...form.value });
};
</script>
