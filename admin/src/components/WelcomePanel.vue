<template>
  <div class="flex justify-center items-center h-full">
    <Card class="w-full max-w-[700px] !bg-[#121826] !border-slate-800 !rounded-xl">
      <template #title>
        <div class="flex items-center gap-3 font-bold text-xl text-slate-100">
          <i class="pi pi-desktop text-indigo-500" />
          <span>DBAdmin Control Center</span>
        </div>
      </template>
      <template #content>
        <p class="text-slate-400 text-sm leading-relaxed">
          Kelola dan rancang skema database PostgreSQL Anda secara visual dengan antarmuka modern yang rapi.
        </p>
        <div class="grid grid-cols-2 gap-5 my-6">
          <div class="bg-[#0c0f19] border border-slate-800 rounded-lg p-4 flex flex-col items-center gap-1">
            <span class="text-3xl font-bold text-slate-100">{{ dbCount }}</span>
            <span class="text-xs text-slate-500 uppercase font-semibold tracking-wider">Database Dibuat</span>
          </div>
          <div class="bg-[#0c0f19] border border-slate-800 rounded-lg p-4 flex flex-col items-center gap-1">
            <span class="text-3xl font-bold text-emerald-500">Active</span>
            <span class="text-xs text-slate-500 uppercase font-semibold tracking-wider">Koneksi PostgreSQL</span>
          </div>
        </div>
        
        <div>
          <h3 class="text-sm font-bold text-slate-200 mb-2.5">Petunjuk Penggunaan</h3>
          <ul class="pl-5 text-slate-400 text-xs flex flex-col gap-2 list-disc leading-relaxed">
            <li>Klik tombol <strong>Database Baru</strong> di sebelah kiri untuk meluncurkan database terisolasi baru.</li>
            <li>Klik salah satu database di navigasi kiri untuk melihat, membuat, mengganti nama, atau menghapus tabel.</li>
            <li>Data tersimpan secara persisten di Docker volume di komputer lokal Anda.</li>
          </ul>
        </div>
      </template>
    </Card>
  </div>
</template>

<script setup>
import Card from 'primevue/card'

defineProps({
  dbCount: {
    type: Number,
    required: true
  }
})
</script>
