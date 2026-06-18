<template>
    <div class="space-y-6 p-8 w-full max-w-5xl mx-auto">
        <p class="text-2xl">Server Info</p>
        <div class="grid grid-cols-3 gap-6">
            <div class="rounded-md p-4 border border-white/10">
                <span class="text-sm opacity-70 block pb-2">Uptime</span>
                <span class="text-green-300 block text-2xl font-bold"> 1 Min 30 Sec</span>
            </div>
            <div class="rounded-md p-4 border border-white/10">
                <span class="text-sm opacity-70 block pb-2">Memory Usage</span>
                <span class="text-green-300 block text-2xl font-bold">500 MB / 4 GB</span>
            </div>
            <div class="rounded-md p-4 border border-white/10">
                <span class="text-sm opacity-70 block pb-2">CPU Usage</span>
                <span class="text-green-300 block text-2xl font-bold">15%</span>
            </div>
            <div class="rounded-md p-4 border border-white/10 col-span-3">
                <span class="text-sm opacity-70 block pb-2">Other Info</span>
                <div class="grid grid-cols-2">
                    <ul>
                        <li>Operating Systems</li>
                        <li>Language Runtimes</li>
                        <li>Database</li>
                        <li></li>
                    </ul>
                    <ul class="text-green-300">
                        <li>Ubuntu 24.04.0 LTS</li>
                        <li>Go 1.24</li>
                        <li>PostgreSQL 16</li>
                    </ul>
                </div>
            </div>
        </div>
        <p class="text-2xl mt-20">Quick Links</p>
        <ul class="*:p-4 space-y-4 *:border *:border-white/5 *:rounded-md *:hover:bg-white/10 *:cursor-pointer">
            <li @click="$emit('databases')" class="flex items-center gap-2">
                <i class="pi pi-database mr-2"></i>
                <span class="text-lg">Databases</span>
                <i class="pi pi-chevron-right ml-auto"></i>
            </li>
            <li @click="$emit('sites')" class="flex items-center gap-2">
                <i class="pi pi-cloud mr-2"></i>
                <span class="text-lg">Sites</span>
                <i class="pi pi-chevron-right ml-auto"></i>
            </li>
        </ul>
    </div>
</template>
