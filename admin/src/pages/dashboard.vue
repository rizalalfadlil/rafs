<template>
  <div class="space-y-6 p-8 w-full max-w-5xl mx-auto">
    <p class="text-2xl font-semibold">Server Info</p>
    <div v-if="errorMsg" class="bg-red-500/10 border border-red-500/20 text-red-400 p-4 rounded-md text-sm">
      Gagal memuat status server: {{ errorMsg }}
    </div>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="rounded-md p-4 border border-white/10 bg-zinc-900/40 backdrop-blur-sm">
        <span class="text-sm opacity-70 block pb-2 text-zinc-400">Uptime</span>
        <span class="text-green-300 block text-2xl font-bold font-mono">{{ info.uptime }}</span>
      </div>
      <div class="rounded-md p-4 border border-white/10 bg-zinc-900/40 backdrop-blur-sm">
        <span class="text-sm opacity-70 block pb-2 text-zinc-400">Memory Usage</span>
        <span class="text-green-300 block text-2xl font-bold font-mono">{{ info.memory_used }} / {{ info.memory_total
          }}</span>
        <span class="text-xs text-zinc-500 block pt-1">Rasio: {{ info.memory_percentage }}</span>
      </div>
      <div class="rounded-md p-4 border border-white/10 bg-zinc-900/40 backdrop-blur-sm">
        <span class="text-sm opacity-70 block pb-2 text-zinc-400">CPU Usage</span>
        <span class="text-green-300 block text-2xl font-bold font-mono">{{ info.cpu_usage }}</span>
      </div>
      <div class="rounded-md p-4 border border-white/10 bg-zinc-900/40 backdrop-blur-sm col-span-1 md:col-span-3">
        <span class="text-sm opacity-70 block pb-2 text-zinc-400">Other Info</span>
        <div class="grid grid-cols-2 max-w-xl">
          <ul class="space-y-2 text-zinc-400 text-sm">
            <li>Operating System</li>
            <li>Language Runtime</li>
            <li>Database</li>
          </ul>
          <ul class="text-green-300 font-mono space-y-2 text-sm">
            <li>{{ info.os }}</li>
            <li>{{ info.go_version }}</li>
            <li>{{ info.postgres_version }}</li>
          </ul>
        </div>
      </div>
    </div>

    <p class="text-2xl mt-16 font-semibold">Quick Links</p>
    <ul
      class="*:p-4 space-y-4 *:border *:border-white/5 *:rounded-md *:hover:bg-white/10 *:cursor-pointer *:transition-all *:duration-200 bg-zinc-900/10">
      <li @click="router.push(menu.path)" v-for="menu in menuItems" :key="menu.path" class="flex items-center gap-2">
        <i :class="menu.icon + ' mr-2 text-zinc-400'"></i>
        <span class="text-lg">{{ menu.label }}</span>
        <i class="pi pi-chevron-right ml-auto text-zinc-500"></i>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useNavigation } from '../composables/useNavigation';

const router = useRouter()
const { menuItems } = useNavigation()

const info = ref({
  uptime: 'Loading...',
  memory_used: '0 MB',
  memory_total: '0 MB',
  memory_percentage: '0%',
  cpu_usage: '0%',
  os: 'Loading...',
  go_version: 'Loading...',
  postgres_version: 'Loading...'
});

const errorMsg = ref('');

const fetchInfo = async () => {
  try {
    const res = await fetch('/api/server-info');
    if (!res.ok) {
      throw new Error(`HTTP error! status: ${res.status}`);
    }
    const data = await res.json();
    if (data.status === 'sukses') {
      info.value = data.info;
      errorMsg.value = '';
    } else {
      errorMsg.value = data.message;
    }
  } catch (err) {
    errorMsg.value = err.message;
  }
};

let intervalId = null;
onMounted(() => {
  fetchInfo();
  // Poll every 3 seconds to keep metrics updated
  intervalId = setInterval(fetchInfo, 3000);
});

onUnmounted(() => {
  if (intervalId) {
    clearInterval(intervalId);
  }
});
</script>
