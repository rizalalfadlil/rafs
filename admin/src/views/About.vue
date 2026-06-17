<template>
  <div class="flex justify-center items-start h-full">
    <Card class="w-full max-w-[750px] !bg-[#121826] !border-slate-800 !rounded-xl">
      <template #title>
        <div class="flex items-center gap-3 text-lg font-bold text-slate-100">
          <InfoIcon class="text-indigo-500 w-6 h-6" />
          <span>Informasi Sistem Server</span>
        </div>
      </template>
      <template #content>
        <div class="flex flex-col gap-6">
          <p class="text-slate-400 text-sm leading-relaxed">
            Sistem ini adalah server pengembangan lokal (local development server) mandiri yang diorkestrasikan dengan Docker untuk membantu proses belajar infrastruktur, hosting tugas kuliah, serta penyediaan basis data dinamis secara terisolasi.
          </p>

          <div class="bg-[#0c0f19] border border-slate-800 rounded-lg p-4 flex flex-col gap-3">
            <div class="flex justify-between text-sm border-b border-slate-800/50 pb-2.5 last:border-0 last:pb-0">
              <span class="text-slate-500 font-medium">Sistem Operasi Host</span>
              <span class="text-slate-200 font-semibold">Windows (Local Host)</span>
            </div>
            <div class="flex justify-between text-sm border-b border-slate-800/50 pb-2.5 last:border-0 last:pb-0">
              <span class="text-slate-500 font-medium">Bahasa Backend</span>
              <span class="text-slate-200 font-semibold">Golang v1.26.4 (Alpine Linux Container)</span>
            </div>
            <div class="flex justify-between text-sm border-b border-slate-800/50 pb-2.5 last:border-0 last:pb-0">
              <span class="text-slate-500 font-medium">Database Induk</span>
              <span class="text-slate-200 font-semibold">PostgreSQL v16.3 (Alpine Container)</span>
            </div>
            <div class="flex justify-between text-sm border-b border-slate-800/50 pb-2.5 last:border-0 last:pb-0">
              <span class="text-slate-500 font-medium">Frontend Framework</span>
              <span class="text-slate-200 font-semibold">Vue 3 + Vite + PrimeVue</span>
            </div>
            <div class="flex justify-between text-sm border-b border-slate-800/50 pb-2.5 last:border-0 last:pb-0">
              <span class="text-slate-500 font-medium">Status Container</span>
              <span class="text-emerald-500 font-semibold">Orchestrated (Docker Compose)</span>
            </div>
          </div>

          <div class="system-features">
            <h3 class="text-base font-bold mb-4 text-slate-100">Fitur & Keunggulan Platform</h3>
            <div class="grid grid-cols-2 gap-5">
              <div class="flex gap-3.5">
                <i class="pi pi-shield text-indigo-500 text-lg mt-0.5" />
                <div>
                  <h4 class="text-sm font-bold mb-1.5 text-slate-200">Keamanan & Isolasi</h4>
                  <p class="text-xs text-slate-400 leading-relaxed">Setiap database yang dibuat otomatis terhubung ke user PostgreSQL baru dengan kontrol penuh, mencegah interferensi data antar proyek.</p>
                </div>
              </div>
              <div class="flex gap-3.5">
                <i class="pi pi-box text-indigo-500 text-lg mt-0.5" />
                <div>
                  <h4 class="text-sm font-bold mb-1.5 text-slate-200">Penyimpanan Persisten</h4>
                  <p class="text-xs text-slate-400 leading-relaxed">Semua berkas database tersimpan di volume lokal (`pgdata`) agar data tidak hilang ketika kontainer dimatikan.</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </Card>
  </div>
</template>

<script setup>
import Card from 'primevue/card'
import { Info as InfoIcon } from 'lucide-vue-next'
</script>
