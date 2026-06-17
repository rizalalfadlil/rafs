<template>
  <div class="flex h-screen w-screen bg-[#0c0f19] text-slate-100">
    <!-- PrimeVue Toast Component -->
    <Toast />

    <!-- Sidebar -->
    <aside class="w-[280px] bg-[#121826] border-r border-slate-800 flex flex-col shrink-0">
      <div class="p-6 border-b border-slate-800">
        <div class="flex items-center gap-3">
          <DatabaseIcon class="text-indigo-500 w-7 h-7 drop-shadow-[0_0_8px_rgba(99,102,241,0.4)]" />
          <span class="text-2xl font-bold tracking-wide bg-gradient-to-r from-slate-100 to-indigo-500 bg-clip-text text-transparent">DBAdmin</span>
        </div>
      </div>

      <nav class="flex-1 py-6 px-4 flex flex-col gap-2">
        <button 
          :class="[
            'flex items-center gap-3 py-3 px-4 border rounded-lg text-sm font-medium cursor-pointer transition-all duration-200 w-full text-left', 
            currentView === 'dashboard' ? 'bg-indigo-500/10 border-indigo-500 text-indigo-400' : 'bg-transparent border-transparent text-slate-400 hover:bg-slate-850 hover:text-slate-100'
          ]"
          @click="currentView = 'dashboard'"
        >
          <LayoutDashboardIcon class="w-[18px] h-[18px]" :class="currentView === 'dashboard' ? 'text-indigo-500' : 'text-slate-400'" />
          <span>Dashboard</span>
        </button>
        <button 
          :class="[
            'flex items-center gap-3 py-3 px-4 border rounded-lg text-sm font-medium cursor-pointer transition-all duration-200 w-full text-left', 
            currentView === 'about' ? 'bg-indigo-500/10 border-indigo-500 text-indigo-400' : 'bg-transparent border-transparent text-slate-400 hover:bg-slate-850 hover:text-slate-100'
          ]"
          @click="currentView = 'about'"
        >
          <InfoIcon class="w-[18px] h-[18px]" :class="currentView === 'about' ? 'text-indigo-500' : 'text-slate-400'" />
          <span>Tentang Server</span>
        </button>
      </nav>

      <div class="p-5 border-t border-slate-800">
        <div class="flex items-center gap-2 text-xs font-semibold text-slate-400 bg-[#0c0f19] py-2 px-3.5 rounded-full border border-slate-800 w-fit">
          <span class="w-2 h-2 rounded-full bg-emerald-500 shadow-[0_0_8px_#10b981]"></span>
          <span>PostgreSQL Active</span>
        </div>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 flex flex-col overflow-hidden">
      <header class="h-[72px] bg-[#121826] border-b border-slate-800 px-8 flex items-center justify-between">
        <div class="flex items-center gap-2 text-sm font-medium">
          <span class="text-slate-400">Admin Portal</span>
          <span class="text-slate-600">/</span>
          <span class="text-slate-100">{{ viewTitle }}</span>
        </div>
        <div class="text-xs text-slate-400 border border-slate-800 py-1 px-2.5 rounded">
          <span>Local server</span>
        </div>
      </header>

      <!-- Panel View -->
      <div class="flex-1 overflow-y-auto p-8">
        <DashboardView v-slot v-if="currentView === 'dashboard'" />
        <AboutView v-else-if="currentView === 'about'" />
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import Toast from 'primevue/toast'
import { 
  Database as DatabaseIcon, 
  LayoutDashboard as LayoutDashboardIcon, 
  Info as InfoIcon 
} from 'lucide-vue-next'

// Views
import DashboardView from './views/Dashboard.vue'
import AboutView from './views/About.vue'

const currentView = ref('dashboard')

const viewTitle = computed(() => {
  return currentView.value === 'dashboard' ? 'Dashboard' : 'Tentang Server'
})
</script>
