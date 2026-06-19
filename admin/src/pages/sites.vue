<script setup lang="ts">
import Button from 'primevue/button';
import Badge from 'primevue/badge';

const sitesList = [
    {
        name: "test-site",
        active: true,
    },
    {
        name: "homepage-1",
        active: false,
    }
]
</script>

<template>
    <main class="p-8 space-y-8 *:space-y-2 max-w-5xl mx-auto">
        <h1 class="text-2xl font-bold">Sites</h1>
        <section>
            <p class="text-xl">Create a new site</p>
            <div class="flex gap-2">
                <Button label="Upload Folder" icon="pi pi-upload" />
                <Button label="Clone From Github" icon="pi pi-github" severity="secondary" />
            </div>
        </section>
        <section>
            <p class="text-xl">Guides</p>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-2">
                <Button variant="outlined" severity="secondary">
                    <img src="https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/html5/html5-original.svg"
                        class="size-6" />
                    native html
                </Button>
            </div>
        </section>
        <section>
            <p class="text-xl">All Sites</p>
            <div class="grid gap-2 grid-cols-2">
                <div v-for="site in sitesList" class="border rounded p-4 border-white/10 flex justify-between">
                    <div class="space-y-1">
                        <p>{{ site.name }}</p>
                        <p class="text-xs text-muted text-white/50">
                            <Badge :severity="site.active ? 'success' : 'secondary'" /> {{ site.active ? "active"
                            : "need configuration" }}
                        </p>
                    </div>
                    <Button variant="text" :severity="site.active ? 'success' : 'warn'" :icon="site.active ? 'pi pi-arrow-right' : 'pi pi-info-circle'" />
                </div>
            </div>
        </section>
    </main>
</template>