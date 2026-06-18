<template>
    <div class="p-2 bg-white/5 border-b border-white/5 px-4 flex gap-4 items-center text-sm text-white/50 select-none">
        <i class="pi pi-server" @click="back"></i>
        <i class="pi pi-chevron-right text-sm"></i>
        <div v-if="activeDatabase" class="flex items-center gap-4" @click="activeTable = null">
            <i class="pi pi-database"></i>
            <span>{{ activeDatabase?.name }}</span>
        </div>
        <div v-if="activeTable" class="flex items-center gap-4">
            <i class="pi pi-chevron-right text-sm"></i>
            <i class="pi pi-table"></i>
            <span>{{ activeTable.name }}</span>
        </div>
    </div>
    <div class="flex">
        <aside class="w-72 border-e border-white/10 p-4 h-screen select-none">
            <p class="text-lg mb-4" @click="activeDatabase = null">Databases</p>
            <ul class="space-y-2">
                <li v-for="db in databases" :key="db" @click="activeDatabase = db; activeTable = null"
                    :class="['p-2 cursor-pointer hover:outline outline-white/5 rounded-md', activeDatabase === db ? 'bg-white/5' : '']">
                    {{ db.name }}</li>
                <li>
                    <Button class="w-full">
                        <i class="pi pi-plus"></i>
                        <span>New Database</span>
                    </Button>
                </li>
            </ul>
        </aside>
        <div v-if="activeDatabase" class="p-4 w-full overflow-y-auto h-[calc(100vh-200px)]">
            <div class="w-full">
                <div class="flex mt-4 gap-4">
                    <div class="w-full space-y-4">
                        <div v-if="activeDatabase && !activeTable" class="space-y-4">
                            <p class="text-lg">Details</p>
                            <IftaLabel class="w-full">
                                <InputText class="w-full" v-model="activeDatabase.name" />
                                <label for="db_name">Database Name</label>
                            </IftaLabel>
                            <div class="flex gap-2 *:w-[300px]">
                                <Button :variant="outlined">Save</Button>
                                <Button severity="danger" variant="outlined">Drop Database</Button>
                            </div>
                        </div>
                        <div v-else class="p-2">
                            <p class="text-lg">{{ activeDatabase.name }}</p>
                        </div>
                        <p class="text-lg">Tables</p>
                        <ul class="mt-2 flex gap-2 border border-white/5 rounded">
                            <li v-for="table in activeDatabase.tables" :key="table"
                                :class="['p-2 px-8 cursor-pointer hover:outline outline-white/5 rounded-md items-center justify-center text-center', activeTable === table ? 'bg-white/5' : '']"
                                @click="activeTable = table">{{ table.name }}
                            </li>
                            <li class="flex-grow"></li>
                            <li>
                                <Button class="px-8">
                                    <i class="pi pi-plus"></i>
                                    <span>New</span>
                                </Button>
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
            <div v-if="activeTable" class="mt-4 space-y-4 rounded-b-md">
                <IftaLabel class="w-full">
                    <InputText class="w-full" v-model="activeTable.name" />
                    <label for="table_name">Table Name</label>
                </IftaLabel>
                <div class="flex gap-2 *:w-[300px]">
                    <Button :variant="outlined">Save</Button>
                    <Button severity="danger" variant="outlined">Drop Table</Button>
                </div>
                <p class="text-lg">Columns</p>
                <table class="w-full text-xs text-white/50  **:border **:border-white/10">
                    <tr class="text-white">
                        <th class="p-4 w-[2rem]"><i class="pi pi-pencil border-none"></i></th>
                        <th class="p-4">Name</th>
                        <th class="p-4">Type</th>
                    </tr>
                    <tr v-for="column in activeTable.columns" class="*:p-4">
                        <td>
                            <Checkbox v-model="column.checked" />
                        </td>
                        <td>{{ column.name }}</td>
                        <td>{{ column.type }}</td>
                    </tr>
                </table>
                <div class="flex gap-2">
                    <Button :variant="outlined">
                        <i class="pi pi-plus"></i>
                        <span>Add Column</span>
                    </Button>
                    <Button severity="danger" variant="outlined">
                        <i class="pi pi-trash"></i>
                        <span>Delete Selected</span>
                    </Button>
                </div>
                <p class="text-lg">Rows</p>
                <div class="overflow-x-auto">
                    <table class="w-full text-xs text-white/50  **:border **:border-white/10">
                        <tr class="text-white">
                            <th class="p-4 w-[2rem]"><i class="pi pi-pencil border-none"></i></th>
                            <th v-for="col in activeTable.columns" class="p-4">{{ col.name }}</th>
                        </tr>
                        <tr v-for="row in activeTable.rows" class="*:p-4">
                            <td>
                                <Checkbox v-model="row.checked" />
                            </td>
                            <td v-for="col in row">{{ col }}</td>
                        </tr>
                    </table>
                </div>
                <div class="flex gap-2">
                    <Button :variant="outlined">
                        <i class="pi pi-plus"></i>
                        <span>Add Row</span>
                    </Button>
                    <Button severity="danger" variant="outlined">
                        <i class="pi pi-trash"></i>
                        <span>Delete Selected</span>
                    </Button>
                </div>
            </div>
            <div v-else class="flex items-center justify-center w-full">
                <p class="text-lg opacity-50 mt-auto">Pilih Tabel</p>
            </div>
        </div>
        <div v-else class="flex items-center justify-center w-full">
            <p class="text-lg opacity-50">Pilih Database</p>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import IftaLabel from 'primevue/iftalabel';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import Checkbox from 'primevue/checkbox';

let databases = ref([
    {
        name: "db_sekolah",
        tables: [
            {
                name: "Siswa",
                columns: [
                    { name: "id", type: "INTEGER" },
                    { name: "name", type: "VARCHAR" },
                    { name: "umur", type: "INTEGER" },
                    { name: "kelas", type: "INTEGER" }
                ],
                rows: [
                    { id: 1, name: "Budi", umur: 20, kelas: 1 },
                    { id: 2, name: "Ani", umur: 19, kelas: 2 }
                ]
            },
            {
                name: "Guru",
                columns: [
                    { name: "id", type: "INTEGER" },
                    { name: "name", type: "VARCHAR" },
                    { name: "pelajaran", type: "VARCHAR" }
                ],
                rows: [
                    { id: 1, name: "Pak Budi", pelajaran: "Matematika" },
                    { id: 2, name: "Ibu Ani", pelajaran: "Fisika" }
                ]
            },
            {
                name: "Kelas",
                columns: [
                    { name: "id", type: "INTEGER" },
                    { name: "jumlah_siswa", type: "INTEGER" },
                    { name: "reputasi", type: "REAL" }
                ],
                rows: [
                    { id: 1, jumlah_siswa: 30, reputasi: 0.80 },
                    { id: 2, jumlah_siswa: 20, reputasi: 0.70 },
                    { id: 3, jumlah_siswa: 20, reputasi: 0.50 },
                    { id: 4, jumlah_siswa: 20, reputasi: 0.40 },
                    { id: 5, jumlah_siswa: 30, reputasi: 0.70 },
                    { id: 6, jumlah_siswa: 20, reputasi: 0.60 }
                ]
            }
        ]
    },
    {
        name: "toko_online",
        tables: [
            {
                name: "produk",
                columns: [
                    { name: "id", type: "INTEGER" },
                    { name: "name", type: "VARCHAR" },
                    { name: "harga", type: "REAL" },
                    { name: "stok", type: "INTEGER" }
                ],
                rows: [
                    { id: 1, name: "Buku", harga: 10000, stok: 10 },
                    { id: 2, name: "Pensil", harga: 5000, stok: 10 },
                    { id: 3, name: "Pulpen", harga: 4000, stok: 10 }
                ]
            },
            {
                name: "toko",
                columns: [
                    { name: "id", type: "INTEGER" },
                    { name: "name", type: "VARCHAR" },
                    { name: "lokasi", type: "VARCHAR" },
                    { name: "barang", type: "ARRAY" }
                ],
                rows: [
                    { id: 1, name: "Toko 1", lokasi: "Jakarta", barang: ["Buku", "Pensil"] },
                    { id: 2, name: "Toko 2", lokasi: "Bandung", barang: ["Pensil", "Pulpen"] },
                    { id: 3, name: "Toko 3", lokasi: "Surabaya", barang: ["Pulpen", "Buku", "Pensil"] }
                ]
            },
            {
                name: "user",
                columns: [
                    { name: "id", type: "INTEGER" },
                    { name: "name", type: "VARCHAR" }
                ],
                rows: [
                    { id: 1, name: "Budi" },
                    { id: 2, name: "Ani" },
                    { id: 3, name: "Rudi" }
                ]
            }
        ]
    }
])

let activeDatabase = ref(null)
let activeTable = ref(null)
</script>
