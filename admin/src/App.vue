<template>
  <div class="flex h-screen w-screen">
    <aside class="border-e transition-all duration-300 hidden sm:block border-white/10"
      :class="!collapse ? 'w-72' : 'w-16'">
      <div v-if="!collapse" class="p-4 pt-6">
        <h1 class="text-2xl font-bold">RFPS</h1>
        <p class="text-sm text-green-300">RizalAlfadlil Personal Server</p>
      </div>
      <div class="flex flex-col gap-2 mt-8 p-2 h-full" :class="collapse ? 'justify-center' : 'justify-start'">
        <router-link v-for="item in menuItems" :key="item.label" :to="item.path" custom v-slot="{ isActive, navigate }">
          <div @click="navigate"
            :class="['flex items-center gap-2 p-4 rounded select-none cursor-pointer transition-all duration-200 hover:outline outline-white/10', isActive ? 'bg-white/10' : '']">
            <i :class="item.icon"></i>
            <span :class="!collapse ? 'opacity-100 translate-x-0' : 'opacity-0 -translate-x-100'"
              class="transition-all duration-300">{{ item.label }}</span>
          </div>
        </router-link>
      </div>
    </aside>
    <div class="flex-1 flex flex-col">
      <header class="h-14 p-8 border-b border-white/10 flex items-center">
        <Button :icon="collapse?'pi pi-bars':'pi pi-times'" class="me-4" variant="text" @click="collapse = !collapse" />
        <span class="font-medium"> {{ activeMenu }}</span>
      </header>
      <main class="flex-1 relative overflow-y-auto">
        <router-view></router-view>
      </main>
    </div>
    <Dock :collapse="collapse" />
  </div>
</template>

<script setup>
import 'primeicons/primeicons.css';
import Button from 'primevue/button';
import { useNavigation } from './composables/useNavigation';
import { ref } from 'vue';
import Dock from './components/Dock.vue';
const { activeMenu, menuItems } = useNavigation();

let collapse = ref(false)
</script>