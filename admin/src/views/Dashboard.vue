<template>
  <div class="flex gap-12 h-full">
    <!-- Left Column: Databases List & Search -->
    <DbNavigator 
      :databases="databases" 
      :selectedDb="selectedDb" 
      @select="selectDatabase" 
      @create-click="showCreateDbModal" 
    />

    <!-- Right Column: Context Content -->
    <div class="flex-1 min-w-0">
      <!-- Welcome Panel -->
      <WelcomePanel v-if="!selectedDb" :dbCount="databases.length" />

      <!-- Database Details Panel -->
      <div v-else class="flex flex-col gap-6 h-full">
        <div class="flex justify-between items-center border-b border-slate-800 pb-4">
          <div class="flex items-center gap-4">
            <i class="pi pi-folder-open text-3xl text-indigo-500" />
            <div>
              <h2 class="text-2xl font-bold m-0">{{ selectedDb }}</h2>
              <span class="text-xs text-slate-500">Skema Publik dan Manajemen Tabel</span>
            </div>
          </div>
          <div class="flex gap-2.5">
            <Button 
              icon="pi pi-pencil" 
              label="Ubah Nama" 
              class="p-button-outlined p-button-sm"
              @click="showRenameDbModal"
            />
            <Button 
              icon="pi pi-trash" 
              label="Hapus DB" 
              class="p-button-danger p-button-sm"
              @click="confirmDeleteDatabase"
            />
          </div>
        </div>

        <!-- Tables Subsection -->
        <div class="tables-section">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-bold">Daftar Tabel</h3>
            <Button 
              icon="pi pi-plus" 
              label="Buat Tabel Baru" 
              class="p-button-sm p-button-primary"
              @click="showCreateTableModal"
            />
          </div>

          <!-- Tables Grid -->
          <div class="grid grid-cols-[repeat(auto-fill,minmax(240px,1fr))] gap-4" v-if="tables.length > 0">
            <TableCard 
              v-for="table in tables" 
              :key="table" 
              :name="table" 
              @rename="showRenameTableModal"
              @delete="confirmDeleteTable"
            />
          </div>

          <!-- Empty State for Tables -->
          <div v-else class="bg-[#121826] border border-dashed border-slate-800 rounded-lg p-12 text-center text-slate-500">
            <i class="pi pi-file-excel text-3xl mb-3" />
            <p>Database ini tidak memiliki tabel.</p>
            <Button 
              icon="pi pi-plus" 
              label="Buat Tabel Pertama" 
              class="p-button-outlined p-button-sm mt-3"
              @click="showCreateTableModal"
            />
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- MODALS -->
  <CreateDbDialog 
    v-model:visible="createDbVisible" 
    :loading="actionLoading" 
    @submit="handleCreateDatabase" 
  />

  <RenameDbDialog 
    v-model:visible="renameDbVisible" 
    :currentName="selectedDb || ''" 
    :loading="actionLoading" 
    @submit="handleRenameDatabase" 
  />

  <CreateTableDialog 
    v-model:visible="createTableVisible" 
    :loading="actionLoading" 
    @submit="handleCreateTable" 
  />

  <RenameTableDialog 
    v-model:visible="renameTableVisible" 
    :currentName="renameTable.old_name" 
    :loading="actionLoading" 
    @submit="handleRenameTable" 
  />
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import Button from 'primevue/button'

// Reusable Components
import DbNavigator from '../components/DbNavigator.vue'
import WelcomePanel from '../components/WelcomePanel.vue'
import TableCard from '../components/TableCard.vue'

// Dialog Components
import CreateDbDialog from '../components/CreateDbDialog.vue'
import RenameDbDialog from '../components/RenameDbDialog.vue'
import CreateTableDialog from '../components/CreateTableDialog.vue'
import RenameTableDialog from '../components/RenameTableDialog.vue'

const toast = useToast()

// States
const databases = ref([])
const selectedDb = ref(null)
const tables = ref([])
const actionLoading = ref(false)

// Modals Visibility
const createDbVisible = ref(false)
const renameDbVisible = ref(false)
const createTableVisible = ref(false)
const renameTableVisible = ref(false)

// Form states needing external reference
const renameTable = ref({ old_name: '', new_name: '' })

// Lifecycle
onMounted(() => {
  fetchDatabases()
})

// Helper Notify
const notifySuccess = (msg) => toast.add({ severity: 'success', summary: 'Sukses', detail: msg, life: 3000 })
const notifyError = (msg) => toast.add({ severity: 'error', summary: 'Gagal', detail: msg, life: 3000 })

// Fetch Databases
const fetchDatabases = async () => {
  try {
    const res = await fetch('/api/databases')
    const data = await res.json()
    if (data.status === 'sukses') {
      databases.value = data.databases || []
    } else {
      notifyError(data.message || 'Gagal memuat database')
    }
  } catch (err) {
    notifyError('Gagal menghubungi server')
  }
}

// Fetch Tables in selected DB
const fetchTables = async (dbName) => {
  try {
    const res = await fetch(`/api/tables?db_name=${dbName}`)
    const data = await res.json()
    if (data.status === 'sukses') {
      tables.value = data.tables || []
    } else {
      notifyError(data.message || 'Gagal memuat daftar tabel')
    }
  } catch (err) {
    notifyError('Gagal menghubungi server')
  }
}

// Selection
const selectDatabase = (dbName) => {
  selectedDb.value = dbName
  fetchTables(dbName)
}

// Modal Show triggers
const showCreateDbModal = () => {
  createDbVisible.value = true
}

const showRenameDbModal = () => {
  renameDbVisible.value = true
}

const showCreateTableModal = () => {
  createTableVisible.value = true
}

const showRenameTableModal = (tableName) => {
  renameTable.value = { old_name: tableName, new_name: tableName }
  renameTableVisible.value = true
}

// API OPERATIONS

// Create DB
const handleCreateDatabase = async (formData) => {
  actionLoading.value = true
  try {
    const res = await fetch('/api/databases', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(formData)
    })
    const data = await res.json()
    if (data.status === 'sukses') {
      notifySuccess('Database berhasil dibuat!')
      createDbVisible.value = false
      await fetchDatabases()
      selectDatabase(formData.db_name)
    } else {
      notifyError(data.message || 'Gagal membuat database')
    }
  } catch (err) {
    notifyError('Gagal menghubungi server')
  } finally {
    actionLoading.value = false
  }
}

// Rename DB
const handleRenameDatabase = async (newName) => {
  if (!newName) return
  actionLoading.value = true
  try {
    const res = await fetch('/api/databases', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ old_name: selectedDb.value, new_name: newName })
    })
    const data = await res.json()
    if (data.status === 'sukses') {
      notifySuccess('Nama database berhasil diubah!')
      renameDbVisible.value = false
      selectedDb.value = newName
      await fetchDatabases()
      selectDatabase(newName)
    } else {
      notifyError(data.message || 'Gagal mengubah nama database')
    }
  } catch (err) {
    notifyError('Gagal menghubungi server')
  } finally {
    actionLoading.value = false
  }
}

// Delete DB
const confirmDeleteDatabase = () => {
  if (confirm(`Apakah Anda yakin ingin menghapus database "${selectedDb.value}" beserta isinya?`)) {
    handleDeleteDatabase()
  }
}

const handleDeleteDatabase = async () => {
  const targetDb = selectedDb.value
  try {
    const res = await fetch(`/api/databases?db_name=${targetDb}`, { method: 'DELETE' })
    const data = await res.json()
    if (data.status === 'sukses') {
      notifySuccess('Database berhasil dihapus!')
      selectedDb.value = null
      tables.value = []
      await fetchDatabases()
    } else {
      notifyError(data.message || 'Gagal menghapus database')
    }
  } catch (err) {
    notifyError('Gagal menghubungi server')
  }
}

// Create Table
const handleCreateTable = async (formData) => {
  if (!formData.table_name || formData.columns.length === 0) {
    notifyError('Nama tabel dan kolom wajib diisi!')
    return
  }

  // Format columns list to payload
  const formattedColumns = formData.columns.map(col => {
    let type = col.selectedType
    if (type === 'custom') {
      type = col.customType
    }
    return { name: col.name, type: type }
  })

  actionLoading.value = true
  try {
    const res = await fetch('/api/tables', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        db_name: selectedDb.value,
        table_name: formData.table_name,
        columns: formattedColumns
      })
    })
    const data = await res.json()
    if (data.status === 'sukses') {
      notifySuccess('Tabel berhasil dibuat!')
      createTableVisible.value = false
      fetchTables(selectedDb.value)
    } else {
      notifyError(data.message || 'Gagal membuat tabel')
    }
  } catch (err) {
    notifyError('Gagal menghubungi server')
  } finally {
    actionLoading.value = false
  }
}

// Rename Table
const handleRenameTable = async (newName) => {
  if (!newName) return
  actionLoading.value = true
  try {
    const res = await fetch('/api/tables', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        db_name: selectedDb.value,
        old_name: renameTable.value.old_name,
        new_name: newName
      })
    })
    const data = await res.json()
    if (data.status === 'sukses') {
      notifySuccess('Nama tabel berhasil diubah!')
      renameTableVisible.value = false
      fetchTables(selectedDb.value)
    } else {
      notifyError(data.message || 'Gagal mengubah nama tabel')
    }
  } catch (err) {
    notifyError('Gagal menghubungi server')
  } finally {
    actionLoading.value = false
  }
}

// Delete Table
const confirmDeleteTable = (tableName) => {
  if (confirm(`Apakah Anda yakin ingin menghapus tabel "${tableName}"?`)) {
    handleDeleteTable(tableName)
  }
}

const handleDeleteTable = async (tableName) => {
  try {
    const res = await fetch(`/api/tables?db_name=${selectedDb.value}&table_name=${tableName}`, { method: 'DELETE' })
    const data = await res.json()
    if (data.status === 'sukses') {
      notifySuccess('Tabel berhasil dihapus!')
      fetchTables(selectedDb.value)
    } else {
      notifyError(data.message || 'Gagal menghapus tabel')
    }
  } catch (err) {
    notifyError('Gagal menghubungi server')
  }
}
</script>

