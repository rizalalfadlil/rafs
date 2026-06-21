<script setup lang="ts">
import { ref, onMounted } from 'vue';
import Button from 'primevue/button';
import Badge from 'primevue/badge';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';

const toast = useToast();

const sitesList = ref<{ name: string; active: boolean }[]>([]);
const isSubmitting = ref(false);

// Dialog visible controls
const showCloneDialog = ref(false);
const showUploadDialog = ref(false);

// Form refs
const cloneForm = ref({ repo_url: '', site_name: '' });
const uploadForm = ref({ site_name: '' });
const selectedFile = ref<File | null>(null);

// Fetch sites list from API
const fetchSites = async () => {
    try {
        const res = await fetch('/api/sites');
        const data = await res.json();
        if (data.status === 'sukses') {
            sitesList.value = data.sites;
        } else {
            toast.add({ severity: 'error', summary: 'Gagal', detail: data.message, life: 3000 });
        }
    } catch (err: any) {
        toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000 });
    }
};

// Handle file input selection
const onFileChange = (e: Event) => {
    const target = e.target as HTMLInputElement;
    if (target.files && target.files.length > 0) {
        selectedFile.value = target.files[0];
    }
};

// Action: Clone from GitHub
const handleClone = async () => {
    if (!cloneForm.value.repo_url || !cloneForm.value.site_name) {
        toast.add({ severity: 'warn', summary: 'Peringatan', detail: 'Semua input wajib diisi', life: 3000 });
        return;
    }
    isSubmitting.value = true;
    try {
        const res = await fetch('/api/sites/clone', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(cloneForm.value)
        });
        const data = await res.json();
        if (data.status === 'sukses') {
            toast.add({ severity: 'success', summary: 'Sukses', detail: data.message, life: 3000 });
            showCloneDialog.value = false;
            cloneForm.value = { repo_url: '', site_name: '' };
            await fetchSites();
        } else {
            toast.add({ severity: 'error', summary: 'Gagal Clone', detail: data.message, life: 4000 });
        }
    } catch (err: any) {
        toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000 });
    } finally {
        isSubmitting.value = false;
    }
};

// Action: Upload ZIP Folder
const handleUpload = async () => {
    if (!uploadForm.value.site_name || !selectedFile.value) {
        toast.add({ severity: 'warn', summary: 'Peringatan', detail: 'Isi nama site dan pilih berkas zip', life: 3000 });
        return;
    }
    isSubmitting.value = true;
    try {
        const formData = new FormData();
        formData.append('site_name', uploadForm.value.site_name);
        formData.append('file', selectedFile.value);

        const res = await fetch('/api/sites/upload', {
            method: 'POST',
            body: formData
        });
        const data = await res.json();
        if (data.status === 'sukses') {
            toast.add({ severity: 'success', summary: 'Sukses', detail: data.message, life: 3000 });
            showUploadDialog.value = false;
            uploadForm.value = { site_name: '' };
            selectedFile.value = null;
            await fetchSites();
        } else {
            toast.add({ severity: 'error', summary: 'Gagal Upload', detail: data.message, life: 4000 });
        }
    } catch (err: any) {
        toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000 });
    } finally {
        isSubmitting.value = false;
    }
};

// Action: Delete Website
const deleteSite = async (siteName: string) => {
    if (!confirm(`Apakah Anda yakin ingin menghapus website '${siteName}'? Folder dan semua isinya akan terhapus permanent.`)) {
        return;
    }
    try {
        const res = await fetch(`/api/sites?site_name=${encodeURIComponent(siteName)}`, {
            method: 'DELETE'
        });
        const data = await res.json();
        if (data.status === 'sukses') {
            toast.add({ severity: 'success', summary: 'Sukses', detail: data.message, life: 3000 });
            await fetchSites();
        } else {
            toast.add({ severity: 'error', summary: 'Gagal Hapus', detail: data.message, life: 3000 });
        }
    } catch (err: any) {
        toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000 });
    }
};

// Action: Open site link
const openSite = (site: { name: string; active: boolean }) => {
    if (site.active) {
        window.open(`/sites/${site.name}/`, '_blank');
    } else {
        toast.add({
            severity: 'info',
            summary: 'Website Tidak Aktif',
            detail: 'Tidak ada file index.html di dalam root direktori website ini.',
            life: 4000
        });
    }
};

onMounted(() => {
    fetchSites();
});
</script>

<template>
    <Toast />

    <main class="p-8 space-y-8 max-w-5xl mx-auto">
        <div class="flex items-center justify-between border-b border-white/10 pb-4">
            <h1 class="text-3xl font-bold text-white flex items-center gap-2">
                Static Sites
            </h1>
        </div>

        <section class="space-y-4">
            <p class="text-xl font-semibold text-zinc-300">Buat website baru</p>
            <div class="grid sm:flex gap-4">
                <Button label="Upload Folder (ZIP)" icon="pi pi-upload" @click="showUploadDialog = true" class="px-6" />
                <Button label="Clone dari Github" icon="pi pi-github" severity="secondary" @click="showCloneDialog = true" class="px-6" />
            </div>
        </section>

        <section class="space-y-4">
            <p class="text-xl font-semibold text-zinc-300">Panduan</p>
            <div class="grid sm:grid-cols-2 md:grid-cols-4 gap-4">
                <Button variant="outlined" severity="secondary" @click="$router.push('/admin/docs')" class="flex gap-2 items-center justify-center p-4">
                    <img src="https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/html5/html5-original.svg"
                        class="size-6" />
                    <span>native html</span>
                </Button>
            </div>
        </section>

        <section class="space-y-4">
            <p class="text-xl font-semibold text-zinc-300">Semua Website</p>
            <div class="grid gap-4 grid-cols-1 md:grid-cols-2">
                <div v-if="sitesList.length === 0" class="col-span-full border border-dashed border-white/10 p-8 rounded-lg text-center text-zinc-500">
                    Belum ada static website yang di-host. Silakan clone repo atau upload zip di atas.
                </div>
                <div v-for="site in sitesList" :key="site.name" 
                     class="border rounded-lg p-5 border-white/10 bg-zinc-950/20 hover:bg-zinc-950/40 hover:border-white/20 transition-all duration-200 flex justify-between items-center">
                    <div class="space-y-2">
                        <p class="text-lg font-semibold text-zinc-100 font-mono">{{ site.name }}</p>
                        <div class="flex items-center gap-2 text-xs">
                            <Badge :severity="site.active ? 'success' : 'secondary'" /> 
                            <span :class="site.active ? 'text-green-400' : 'text-zinc-400'">
                                {{ site.active ? "active" : "need configuration (index.html missing)" }}
                            </span>
                        </div>
                    </div>
                    <div class="flex gap-2">
                        <Button variant="text" 
                                :severity="site.active ? 'success' : 'warn'" 
                                :icon="site.active ? 'pi pi-external-link' : 'pi-info-circle'" 
                                @click="openSite(site)"
                                class="p-button-rounded" />
                        <Button variant="text" 
                                severity="danger" 
                                icon="pi pi-trash" 
                                @click="deleteSite(site.name)"
                                class="p-button-rounded" />
                    </div>
                </div>
            </div>
        </section>
    </main>

    <!-- Dialog: Clone Repository -->
    <Dialog v-model:visible="showCloneDialog" header="Clone Public Repository" :modal="true" class="w-full max-w-md bg-zinc-900 border border-white/10 rounded-lg p-6 text-white">
        <div class="space-y-4">
            <div class="flex flex-col gap-2">
                <label class="text-sm font-semibold">Nama Website</label>
                <InputText v-model="cloneForm.site_name" placeholder="contoh: portofolio-saya" class="w-full bg-white/5 border border-white/10 text-white" />
            </div>
            <div class="flex flex-col gap-2">
                <label class="text-sm font-semibold">GitHub Repository URL (Publik)</label>
                <InputText v-model="cloneForm.repo_url" placeholder="https://github.com/username/repository" class="w-full bg-white/5 border border-white/10 text-white" />
                <span class="text-[10px] text-zinc-500">Repository harus publik agar bisa diclone tanpa memerlukan autentikasi login.</span>
            </div>
            <div class="flex gap-2 justify-end mt-6">
                <Button label="Batal" severity="secondary" variant="outlined" @click="showCloneDialog = false" :disabled="isSubmitting" />
                <Button :label="isSubmitting ? 'Mengkoneksi...' : 'Clone & Host'" @click="handleClone" :disabled="isSubmitting" />
            </div>
        </div>
    </Dialog>

    <!-- Dialog: Upload ZIP Folder -->
    <Dialog v-model:visible="showUploadDialog" header="Upload ZIP Website" :modal="true" class="w-full max-w-md bg-zinc-900 border border-white/10 rounded-lg p-6 text-white">
        <div class="space-y-4">
            <div class="flex flex-col gap-2">
                <label class="text-sm font-semibold">Nama Website</label>
                <InputText v-model="uploadForm.site_name" placeholder="contoh: tugas-kelas" class="w-full bg-white/5 border border-white/10 text-white" />
            </div>
            <div class="flex flex-col gap-2">
                <label class="text-sm font-semibold">Pilih Berkas ZIP</label>
                <input type="file" accept=".zip" @change="onFileChange" class="w-full bg-white/5 border border-white/10 rounded-md p-2 text-white text-sm focus:outline-none" />
                <span class="text-[10px] text-zinc-500">Unggah berkas kompresi .zip yang berisi file static seperti index.html pada root folder.</span>
            </div>
            <div class="flex gap-2 justify-end mt-6">
                <Button label="Batal" severity="secondary" variant="outlined" @click="showUploadDialog = false" :disabled="isSubmitting" />
                <Button :label="isSubmitting ? 'Mengunggah...' : 'Upload & Host'" @click="handleUpload" :disabled="isSubmitting" />
            </div>
        </div>
    </Dialog>
</template>