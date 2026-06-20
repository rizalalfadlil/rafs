<script setup>
import { ref, onMounted, computed } from 'vue';
import ProgressBar from 'primevue/progressbar';
import Badge from 'primevue/badge';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';

const toast = useToast();

// State
const directory = ref([]);
const currentPath = ref("");
const selectedItems = ref(new Set());
const selectMode = ref(false);
const freeMode = ref(false);
const isSubmitting = ref(false);

const space = ref({
    used: 0,
    total: 1024 * 1024 * 1024 // 1GB
});

// Dialogs state
const showNewFolderDialog = ref(false);
const newFolderName = ref("");
const showUploadDialog = ref(false);
const uploadFilesInput = ref(null);
const showPublicLinksDialog = ref(false);
const publicLinks = ref([]);

// Form upload
const selectedFiles = ref([]);

// Computed values
const spacePercent = computed(() => {
    if (!space.value.total) return 0;
    return Math.round((space.value.used / space.value.total) * 100);
});

const usedDisplay = computed(() => {
    return formatBytes(space.value.used);
});

const freeDisplay = computed(() => {
    return formatBytes(space.value.total - space.value.used);
});

const totalDisplay = computed(() => {
    return formatBytes(space.value.total);
});

const pathSegments = computed(() => {
    if (!currentPath.value) return [];
    return currentPath.value.split("/").filter(s => s);
});

// Helper: format bytes
function formatBytes(bytes) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// Get icon for files based on extension
const getFileIcon = (name) => {
    const ext = name.split('.').pop().toLowerCase();
    if (['jpg', 'jpeg', 'png', 'gif', 'svg', 'webp', 'bmp'].includes(ext)) return 'pi-image';
    if (['zip', 'rar', 'tar', 'gz', '7z'].includes(ext)) return 'pi-folder-open';
    if (['pdf'].includes(ext)) return 'pi-file-pdf';
    if (['doc', 'docx'].includes(ext)) return 'pi-file-word';
    if (['xls', 'xlsx'].includes(ext)) return 'pi-file-excel';
    if (['ppt', 'pptx'].includes(ext)) return 'pi-file';
    if (['txt', 'md', 'json', 'yml', 'yaml', 'xml', 'html', 'css', 'js', 'ts'].includes(ext)) return 'pi-file-edit';
    return 'pi-file';
};

// Fetch storage directory content
const fetchDirectory = async () => {
    try {
        const res = await fetch(`/api/storage?path=${encodeURIComponent(currentPath.value)}`);
        const data = await res.json();
        if (data.status === 'sukses') {
            directory.value = data.contents;
            space.value = data.space;
        } else {
            toast.add({ severity: 'error', summary: 'Gagal Memuat', detail: data.message, life: 3000 });
        }
    } catch (err) {
        toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000 });
    }
};

// Handle clicks on grid items
const handleItemClick = (item) => {
    if (selectMode.value) {
        toggleSelect(item.name);
    } else {
        if (item.type === 'folder') {
            currentPath.value = currentPath.value ? `${currentPath.value}/${item.name}` : item.name;
            selectedItems.value.clear();
            fetchDirectory();
        } else {
            downloadFile(item.name);
        }
    }
};

// Navigation
const navigateUp = () => {
    const segments = [...pathSegments.value];
    if (segments.length > 0) {
        segments.pop();
        currentPath.value = segments.join("/");
        selectedItems.value.clear();
        fetchDirectory();
    }
};

const navigateToSegment = (index) => {
    const segments = pathSegments.value.slice(0, index + 1);
    currentPath.value = segments.join("/");
    selectedItems.value.clear();
    fetchDirectory();
};

const navigateToRoot = () => {
    currentPath.value = "";
    selectedItems.value.clear();
    fetchDirectory();
};

// Item selection
const toggleSelect = (name) => {
    if (selectedItems.value.has(name)) {
        selectedItems.value.delete(name);
    } else {
        selectedItems.value.add(name);
    }
};

const isSelected = (name) => {
    return selectedItems.value.has(name);
};

const toggleSelectMode = () => {
    selectMode.value = !selectMode.value;
    selectedItems.value.clear();
};

// Actions
const handleCreateFolder = async () => {
    if (!newFolderName.value.trim()) {
        toast.add({ severity: 'warn', summary: 'Peringatan', detail: 'Nama folder tidak boleh kosong', life: 3000 });
        return;
    }
    isSubmitting.value = true;
    try {
        const res = await fetch('/api/storage/folder', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                path: currentPath.value,
                name: newFolderName.value
            })
        });
        const data = await res.json();
        if (data.status === 'sukses') {
            toast.add({ severity: 'success', summary: 'Sukses', detail: data.message, life: 3000 });
            showNewFolderDialog.value = false;
            newFolderName.value = "";
            fetchDirectory();
        } else {
            toast.add({ severity: 'error', summary: 'Gagal', detail: data.message, life: 3000 });
        }
    } catch (err) {
        toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000 });
    } finally {
        isSubmitting.value = false;
    }
};

const onFileChange = (e) => {
    const target = e.target;
    if (target.files && target.files.length > 0) {
        selectedFiles.value = Array.from(target.files);
    }
};

const handleUploadFiles = async () => {
    if (selectedFiles.value.length === 0) {
        toast.add({ severity: 'warn', summary: 'Peringatan', detail: 'Pilih file terlebih dahulu', life: 3000 });
        return;
    }
    isSubmitting.value = true;
    try {
        const formData = new FormData();
        formData.append('path', currentPath.value);
        selectedFiles.value.forEach(file => {
            formData.append('files', file);
        });

        const res = await fetch('/api/storage/upload', {
            method: 'POST',
            body: formData
        });
        const data = await res.json();
        if (data.status === 'sukses') {
            toast.add({ severity: 'success', summary: 'Sukses', detail: data.message, life: 3000 });
            showUploadDialog.value = false;
            selectedFiles.value = [];
            if (uploadFilesInput.value) uploadFilesInput.value.value = "";
            fetchDirectory();
        } else {
            toast.add({ severity: 'error', summary: 'Gagal', detail: data.message, life: 3000 });
        }
    } catch (err) {
        toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000 });
    } finally {
        isSubmitting.value = false;
    }
};

const handleDeleteSelected = async () => {
    if (selectedItems.value.size === 0) {
        toast.add({ severity: 'warn', summary: 'Peringatan', detail: 'Pilih item yang ingin dihapus', life: 3000 });
        return;
    }
    if (!confirm(`Apakah Anda yakin ingin menghapus ${selectedItems.value.size} item terpilih? Tindakan ini permanen.`)) {
        return;
    }

    const pathsToDelete = Array.from(selectedItems.value).map(name => {
        return currentPath.value ? `${currentPath.value}/${name}` : name;
    });

    isSubmitting.value = true;
    try {
        const res = await fetch('/api/storage/delete', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ paths: pathsToDelete })
        });
        const data = await res.json();
        if (data.status === 'sukses') {
            toast.add({ severity: 'success', summary: 'Sukses', detail: data.message, life: 3000 });
            selectedItems.value.clear();
            selectMode.value = false;
            fetchDirectory();
        } else {
            toast.add({ severity: 'error', summary: 'Gagal', detail: data.message, life: 3000 });
        }
    } catch (err) {
        toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000 });
    } finally {
        isSubmitting.value = false;
    }
};

const downloadFile = (fileName) => {
    const fileRelativePath = currentPath.value ? `${currentPath.value}/${fileName}` : fileName;
    const downloadURL = `/api/storage/download?path=${encodeURIComponent(fileRelativePath)}`;
    window.open(downloadURL, '_blank');
};

const handleSetPublic = async () => {
    if (selectedItems.value.size === 0) {
        toast.add({ severity: 'warn', summary: 'Peringatan', detail: 'Pilih file/folder yang ingin dipublikasikan', life: 3000 });
        return;
    }

    const pathsToPublic = Array.from(selectedItems.value).map(name => {
        return currentPath.value ? `${currentPath.value}/${name}` : name;
    });

    isSubmitting.value = true;
    try {
        const res = await fetch('/api/storage/public', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ paths: pathsToPublic })
        });
        const data = await res.json();
        if (data.status === 'sukses') {
            toast.add({ severity: 'success', summary: 'Sukses', detail: data.message, life: 3000 });
            
            // Generate full URLs
            const host = window.location.origin;
            publicLinks.value = data.links.map(link => `${host}${link}`);
            showPublicLinksDialog.value = true;

            selectedItems.value.clear();
            selectMode.value = false;
        } else {
            toast.add({ severity: 'error', summary: 'Gagal Publikasikan', detail: data.message, life: 3000 });
        }
    } catch (err) {
        toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000 });
    } finally {
        isSubmitting.value = false;
    }
};

const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);
    toast.add({ severity: 'success', summary: 'Disalin', detail: 'Link disalin ke clipboard!', life: 2000 });
};

onMounted(() => {
    fetchDirectory();
});
</script>

<template>
    <Toast />
    
    <div class="p-6 space-y-6">
        <!-- Top Stats and Actions Section -->
        <section class="grid grid-cols-1 md:grid-cols-3 gap-6 items-center border border-white/5 bg-zinc-950/20 p-5 rounded-lg">
            <!-- Quota Progress Card -->
            <div class="md:col-span-2 space-y-2">
                <div class="flex justify-between items-center text-sm font-semibold">
                    <span class="text-zinc-400">Penyimpanan Digunakan</span>
                    <span class="text-zinc-100 font-mono">{{ spacePercent }}%</span>
                </div>
                <ProgressBar :value="spacePercent" :showValue="false" class="h-2 bg-white/5" />
                <div class="flex gap-4 select-none pt-1">
                    <button class="flex items-center text-xs text-zinc-400 hover:text-white transition-colors" @click="freeMode = !freeMode">
                        <Badge class="me-1.5 size-2" :severity="freeMode ? 'secondary' : 'primary'" />
                        <span>{{ freeMode ? freeDisplay : usedDisplay }} {{ freeMode ? 'tersisa' : 'digunakan' }}</span>
                    </button>
                    <div class="flex items-center text-xs text-zinc-400">
                        <Badge class="me-1.5 size-2" severity="contrast" />
                        <span>Total Kuota: {{ totalDisplay }}</span>
                    </div>
                </div>
            </div>
            
            <!-- Global Actions -->
            <div class="flex gap-2 justify-end items-center">
                <template v-if="!selectMode">
                    <Button label="Upload" icon="pi pi-upload" @click="showUploadDialog = true" class="px-4" size="small" />
                    <Button label="Folder Baru" icon="pi pi-plus" severity="secondary" @click="showNewFolderDialog = true" class="px-4" size="small" />
                    <Button label="Pilih" severity="warn" variant="outlined" icon="pi pi-check-circle" @click="toggleSelectMode" class="px-4" size="small" />
                </template>
                <template v-else>
                    <Button label="Set Publik" severity="info" icon="pi pi-globe" @click="handleSetPublic" class="px-4" size="small" />
                    <Button label="Hapus" severity="danger" icon="pi pi-trash" @click="handleDeleteSelected" class="px-4" size="small" />
                    <Button label="Batal" icon="pi pi-times" severity="secondary" variant="outlined" @click="toggleSelectMode" class="px-4" size="small" />
                </template>
            </div>
        </section>

        <!-- Breadcrumb Navigation -->
        <div class="p-3 bg-zinc-900/60 border border-white/5 rounded-lg px-4 flex gap-2 items-center text-sm text-zinc-400 select-none">
            <span @click="navigateToRoot" class="cursor-pointer hover:text-white transition-colors flex items-center gap-1 font-semibold">
                <i class="pi pi-box"></i> Root
            </span>
            
            <template v-for="(segment, idx) in pathSegments" :key="idx">
                <i class="pi pi-chevron-right text-[10px] text-zinc-600"></i>
                <span @click="navigateToSegment(idx)" class="cursor-pointer hover:text-white transition-colors font-mono">
                    {{ segment }}
                </span>
            </template>
            
            <div class="ml-auto" v-if="pathSegments.length > 0 && !selectMode">
                <Button icon="pi pi-arrow-up" severity="secondary" variant="text" size="small" @click="navigateUp" class="h-7 w-7 p-0" title="Kembali ke folder atas" />
            </div>
        </div>

        <!-- Directory Contents Grid -->
        <main class="grid grid-cols-2 sm:grid-cols-4 md:grid-cols-6 lg:grid-cols-8 gap-4 min-h-[300px]">
            <!-- Empty Folder State -->
            <div v-if="directory.length === 0" class="col-span-full border border-dashed border-white/10 p-12 rounded-lg text-center flex flex-col items-center justify-center gap-3 text-zinc-500">
                <i class="pi pi-folder-open text-4xl text-zinc-600"></i>
                <p>Folder ini kosong. Gunakan tombol 'Upload' atau 'Folder Baru' di atas.</p>
            </div>
            
            <!-- Items -->
            <div v-for="item in directory" :key="item.name"
                 @click="handleItemClick(item)"
                 :class="[
                     'relative border p-4 rounded-lg bg-zinc-950/20 flex flex-col items-center justify-between gap-3 text-center aspect-square transition-all duration-200 select-none cursor-pointer',
                     isSelected(item.name) ? 'border-amber-500/60 bg-amber-500/5 ring-1 ring-amber-500/30' : 'border-white/5 hover:border-zinc-500 hover:-translate-y-0.5 hover:bg-zinc-950/40'
                 ]">
                <!-- Select Checkbox Overlay -->
                <div v-if="selectMode" class="absolute top-2 right-2">
                    <i :class="['pi text-xs rounded-full p-1 border', isSelected(item.name) ? 'pi-check bg-amber-500 border-amber-500 text-black font-bold' : 'border-white/20 text-transparent']" style="font-size: 8px;"></i>
                </div>

                <div class="flex-1 flex items-center justify-center">
                    <i v-if="item.type === 'folder'" class="pi pi-folder text-amber-500/80" style="font-size: 3.5rem"></i>
                    <i v-else class="pi text-blue-400/80" :class="getFileIcon(item.name)" style="font-size: 3.5rem"></i>
                </div>
                
                <div class="w-full">
                    <p class="text-xs font-medium truncate text-zinc-200 font-mono w-full" :title="item.name">{{ item.name }}</p>
                    <p v-if="item.type === 'file'" class="text-[10px] text-zinc-500 mt-0.5">{{ formatBytes(item.size) }}</p>
                </div>
            </div>
        </main>
    </div>

    <!-- Dialog: New Folder -->
    <Dialog v-model:visible="showNewFolderDialog" header="Buat Folder Baru" :modal="true" class="w-full max-w-sm bg-zinc-900 border border-white/10 rounded-lg p-6 text-white">
        <div class="space-y-4">
            <div class="flex flex-col gap-2">
                <label class="text-sm font-semibold text-zinc-400">Nama Folder</label>
                <InputText v-model="newFolderName" placeholder="Nama folder..." class="w-full bg-white/5 border border-white/10 text-white" autofocus />
            </div>
            <div class="flex gap-2 justify-end mt-6">
                <Button label="Batal" severity="secondary" variant="outlined" @click="showNewFolderDialog = false" :disabled="isSubmitting" />
                <Button :label="isSubmitting ? 'Memproses...' : 'Buat Folder'" @click="handleCreateFolder" :disabled="isSubmitting" />
            </div>
        </div>
    </Dialog>

    <!-- Dialog: Upload Files -->
    <Dialog v-model:visible="showUploadDialog" header="Unggah File" :modal="true" class="w-full max-w-md bg-zinc-900 border border-white/10 rounded-lg p-6 text-white">
        <div class="space-y-4">
            <div class="flex flex-col gap-2">
                <label class="text-sm font-semibold text-zinc-400">Pilih Berkas</label>
                <input ref="uploadFilesInput" type="file" multiple @change="onFileChange" class="w-full bg-white/5 border border-white/10 rounded-md p-2 text-white text-sm focus:outline-none" />
                <span class="text-[10px] text-zinc-500">Anda dapat memilih satu atau beberapa file sekaligus untuk diunggah ke direktori saat ini.</span>
            </div>
            
            <div v-if="selectedFiles.length > 0" class="border border-white/5 rounded p-2 bg-zinc-950/20 max-h-36 overflow-y-auto space-y-1">
                <p class="text-xs font-semibold text-zinc-400 border-b border-white/5 pb-1">File Terpilih ({{ selectedFiles.length }}):</p>
                <div v-for="file in selectedFiles" :key="file.name" class="flex justify-between text-[10px] text-zinc-300 font-mono">
                    <span class="truncate pr-4">{{ file.name }}</span>
                    <span>{{ formatBytes(file.size) }}</span>
                </div>
            </div>

            <div class="flex gap-2 justify-end mt-6">
                <Button label="Batal" severity="secondary" variant="outlined" @click="showUploadDialog = false" :disabled="isSubmitting" />
                <Button :label="isSubmitting ? 'Mengunggah...' : 'Upload Berkas'" @click="handleUploadFiles" :disabled="isSubmitting" />
            </div>
        </div>
    </Dialog>

    <!-- Dialog: Public Links -->
    <Dialog v-model:visible="showPublicLinksDialog" header="Link Publik Berhasil Dibuat" :modal="true" class="w-full max-w-lg bg-zinc-900 border border-white/10 rounded-lg p-6 text-white animate-fade-in">
        <div class="space-y-4">
            <p class="text-sm text-zinc-400">File Anda telah berhasil disalin ke folder publik server. Gunakan URL di bawah ini untuk mengakses atau menampilkannya di web:</p>
            
            <div class="space-y-3 max-h-60 overflow-y-auto">
                <div v-for="(link, idx) in publicLinks" :key="idx" class="flex items-center gap-2 bg-zinc-950/40 border border-white/5 rounded p-2 text-xs">
                    <span class="font-mono text-blue-400 truncate flex-1 select-all">{{ link }}</span>
                    <Button icon="pi pi-copy" severity="secondary" variant="text" size="small" @click="copyToClipboard(link)" class="h-8 w-8" title="Copy URL" />
                </div>
            </div>

            <div class="flex justify-end mt-6">
                <Button label="Tutup" @click="showPublicLinksDialog = false" />
            </div>
        </div>
    </Dialog>
</template>