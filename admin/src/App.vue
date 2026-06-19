<template>
  <div class="flex h-screen w-screen">
    <aside class="w-72 border-e border-white/10">
      <h1 class="text-2xl font-bold p-4 py-6">DBAdmin</h1>
      <div class="flex flex-col gap-2 mt-8 p-2">
        <div @click="activeMenu = item.label" v-for="item in menuItems" :key="item.label" :class="['flex items-center gap-2 p-4 rounded select-none cursor-pointer transition-all duration-200 hover:outline outline-white/10', activeMenu === item.label ? 'bg-white/10' : '']">
          <i :class="item.icon"></i>
          <span>{{ item.label }}</span>
        </div>
      </div>  
    </aside>
    <div class="flex-1 flex flex-col">
      <header class="h-14 p-8 border-b border-white/10 flex items-center">
        <span class="font-medium">{{ activeMenu }}</span>
      </header>
      <main class="flex-1">
        <Dashboard v-if="activeMenu === 'Dashboard'" @databases="activeMenu = 'Databases'" @sites="activeMenu = 'Sites'"></Dashboard>
        <Databases v-if="activeMenu === 'Databases'"></Databases>
        <Sites v-if="activeMenu === 'Sites'"></Sites>
        <Storage v-if="activeMenu === 'Storage'"></Storage>
      </main>
    </div>
  </div>
</template>

<script setup> 
import { ref } from 'vue'
import 'primeicons/primeicons.css';
import Dashboard from './pages/dashboard.vue';
import Databases from './pages/database.vue';
import Sites from './pages/sites.vue';
import Storage from './pages/storage.vue';

let activeMenu = ref('Dashboard');

const menuItems = ref([
  { label: 'Dashboard', icon: 'pi pi-home' },
  { label: 'Sites', icon: 'pi pi-cloud' },
  { label: 'Databases', icon: 'pi pi-database' },
  { label: 'Storage', icon: 'pi pi-box' },
])


</script>