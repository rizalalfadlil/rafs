<template>
  <aside class="w-72 border-e border-white/10 p-4 select-none flex flex-col justify-between">
    <div>
      <p class="text-lg mb-4 cursor-pointer hover:opacity-80" @click="$emit('back')">Databases</p>
      <ul class="space-y-2 max-h-[calc(100vh-200px)] overflow-y-auto p-1">
        <li v-for="db in databases" :key="db.name" @click="$emit('select-db', db)"
            :class="['p-2 cursor-pointer hover:outline outline-white/5 rounded-md flex items-center justify-between transition-all duration-200', activeDatabase?.name === db.name ? 'bg-white/5 outline outline-white/10' : '']">
          <span class="flex items-center gap-2">
            <i class="pi pi-database text-zinc-400"></i>
            {{ db.name }}
          </span>
        </li>
      </ul>
    </div>
    <div class="pt-4 border-t border-white/5">
      <Button class="w-full flex items-center justify-center gap-2" @click="$emit('create-db-click')">
        <i class="pi pi-plus"></i>
        <span>Database Baru</span>
      </Button>
    </div>
  </aside>
</template>

<script setup>
import Button from 'primevue/button';

defineProps({
  databases: {
    type: Array,
    required: true
  },
  activeDatabase: {
    type: Object,
    default: null
  }
});

defineEmits(['select-db', 'create-db-click', 'back']);
</script>
