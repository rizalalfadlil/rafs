<template>
  <Dialog v-model:visible="visibleState" header="Buat Database Baru" :modal="true" class="w-full max-w-md bg-zinc-900 border border-white/10 rounded-lg p-6 text-white">
    <div class="space-y-4">
      <div class="flex flex-col gap-2">
        <label class="text-sm font-semibold">Nama Database</label>
        <InputText v-model="form.db_name" placeholder="contoh: db_toko" class="w-full bg-white/5 border border-white/10 text-white" />
      </div>
      <div class="flex flex-col gap-2">
        <label class="text-sm font-semibold">Username PostgreSQL</label>
        <InputText v-model="form.username" placeholder="contoh: toko_admin" class="w-full bg-white/5 border border-white/10 text-white" />
      </div>
      <div class="flex flex-col gap-2">
        <label class="text-sm font-semibold">Password</label>
        <InputText v-model="form.password" type="password" placeholder="Password user baru" class="w-full bg-white/5 border border-white/10 text-white" />
      </div>
      <div class="flex gap-2 justify-end mt-4">
        <Button label="Batal" severity="secondary" variant="outlined" @click="visibleState = false" />
        <Button label="Buat Database" @click="submit" />
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

const form = ref({ db_name: '', username: '', password: '' });

watch(() => props.visible, (newVal) => {
  if (newVal) {
    form.value = { db_name: '', username: '', password: '' };
  }
});

const submit = () => {
  emit('save', { ...form.value });
};
</script>
