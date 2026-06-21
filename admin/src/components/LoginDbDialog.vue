<template>
  <Dialog v-model:visible="visibleState" header="Akses Database" :modal="true" class="w-full max-w-md bg-zinc-900 border border-white/10 rounded-lg p-6 text-white">
    <div class="space-y-4">
      <p class="text-sm text-white/60">Silakan masukkan username dan password pemilik database <strong>{{ dbName }}</strong> untuk melanjutkan.</p>
      
      <div class="flex flex-col gap-2">
        <label class="text-sm font-semibold">Username PostgreSQL</label>
        <InputText v-model="form.username" placeholder="contoh: toko_admin" class="w-full bg-white/5 border border-white/10 text-white" />
      </div>
      
      <div class="flex flex-col gap-2">
        <label class="text-sm font-semibold">Password</label>
        <InputText v-model="form.password" type="password" placeholder="Password pemilik database" class="w-full bg-white/5 border border-white/10 text-white" />
      </div>
      
      <div class="flex gap-2 justify-end mt-4">
        <Button label="Batal" severity="secondary" variant="outlined" @click="visibleState = false" />
        <Button label="Masuk & Akses" @click="submit" />
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
  dbName: {
    type: String,
    default: ''
  }
});

const emit = defineEmits(['update:visible', 'login']);

const visibleState = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val)
});

const form = ref({ username: '', password: '' });

watch(() => props.visible, (newVal) => {
  if (newVal) {
    // Ambil default dari cache jika ada
    const cached = localStorage.getItem(`rafs_db_cred_${props.dbName}`);
    if (cached) {
      try {
        const cred = JSON.parse(cached);
        form.value = { username: cred.username, password: cred.password };
      } catch (e) {
        form.value = { username: '', password: '' };
      }
    } else {
      form.value = { username: '', password: '' };
    }
  }
});

const submit = () => {
  emit('login', { ...form.value });
};
</script>
