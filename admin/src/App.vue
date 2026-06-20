<template>
  <div class="flex h-screen w-screen">
    <aside class="w-72 border-e border-white/10">
      <h1 class="text-2xl font-bold p-4 py-6">RFS Admin</h1>
      <div class="flex flex-col gap-2 mt-8 p-2">
        <router-link
        v-for="item in menuItems"
        :key="item.label"
        :to="item.path"
        custom
        v-slot="{ isActive, navigate }"
        >
          <div @click="navigate" 
            :class="['flex items-center gap-2 p-4 rounded select-none cursor-pointer transition-all duration-200 hover:outline outline-white/10', isActive ? 'bg-white/10' : '']">
            <i :class="item.icon"></i>
            <span>{{ item.label }}</span>
          </div>
        </router-link>
      </div>  
    </aside>
    <div class="flex-1 flex flex-col">
      <header class="h-14 p-8 border-b border-white/10 flex items-center">
        <span class="font-medium">{{ activeMenu }}</span>
      </header>
      <main class="flex-1">
        <router-view></router-view>
      </main>
    </div>
  </div>
</template>

<script setup> 
import { ref } from 'vue'
import 'primeicons/primeicons.css';


const menuItems = ref([
  { label: 'Dashboard', path: '/admin', icon: 'pi pi-home' },
  { label: 'Sites', path: '/admin/sites', icon: 'pi pi-cloud' },
  { label: 'Databases', path: '/admin/databases', icon: 'pi pi-database' },
  { label: 'Storage', path: '/admin/storage', icon: 'pi pi-box' },
])


</script>