<template>
  <Toast />

<div class="flex flex-col h-full">
    <!-- Breadcrumb Header -->
  <div class="p-2 bg-white/5 border-b border-white/5 px-4 flex gap-4 items-center text-sm text-white/50 select-none">
    <i class="pi pi-server cursor-pointer hover:text-white" @click="back"></i>
    <i class="pi pi-chevron-right text-sm"></i>
    <div v-if="activeDatabase" class="flex items-center gap-2 cursor-pointer hover:text-white" @click="activeTable = null">
      <i class="pi pi-database"></i>
      <span>{{ activeDatabase?.name }}</span>
    </div>
    <div v-if="activeTable" class="flex items-center gap-2">
      <i class="pi pi-chevron-right text-sm"></i>
      <i class="pi pi-table"></i>
      <span>{{ activeTable.name }}</span>
    </div>
  </div>

  <!-- Main Layout -->
  <div class="flex grow">
    <!-- Sidebar Databases List -->
    <DbSidebar
      :databases="databases"
      :activeDatabase="activeDatabase"
      @select-db="selectDatabase"
      @create-db-click="showCreateDbDialog = true"
      @back="back"
    />

    <!-- Right Panels -->
    <div v-if="activeDatabase" class="p-6 w-full overflow-y-auto h-[calc(100vh-100px)] space-y-6">
      <DbDetails
        :activeDatabase="activeDatabase"
        :activeTable="activeTable"
        @rename-db="saveDatabaseName"
        @drop-db="dropDatabase"
        @select-table="selectTable"
        @create-table-click="showCreateTableDialog = true"
      />

      <TableDetails
        v-if="activeTable"
        :activeDatabase="activeDatabase"
        :activeTable="activeTable"
        @rename-table="saveTableName"
        @drop-table="dropTable"
        @add-column-click="showAddColDialog = true"
        @edit-column-click="openEditColumn"
        @delete-columns="deleteSelectedColumns"
        @add-row-click="openAddRow"
        @edit-row-click="openEditRow"
        @delete-rows="deleteSelectedRows"
      />

      <!-- Table Not Selected State -->
      <div v-else class="flex flex-col items-center justify-center w-full h-[300px] border border-dashed border-white/10 rounded-lg bg-zinc-900/10">
        <i class="pi pi-table text-4xl opacity-20 mb-2"></i>
        <p class="text-base opacity-40">Pilih tabel di atas atau buat tabel baru</p>
      </div>
    </div>

    <!-- Database Not Selected State -->
    <div v-else class="flex flex-col items-center justify-center w-full h-[500px] border border-dashed border-white/10 rounded-lg m-6 bg-zinc-900/10 select-none">
      <i class="pi pi-database text-6xl opacity-20 mb-4 animate-pulse"></i>
      <p class="text-lg opacity-40">Silakan pilih database di menu sidebar untuk mengelola data</p>
    </div>
  </div>
</div>

  <!-- Modals & Dialogs -->
  <CreateDbDialog
    v-model:visible="showCreateDbDialog"
    @save="createDatabase"
  />

  <CreateTableDialog
    v-model:visible="showCreateTableDialog"
    @save="createTable"
  />

  <AddColumnDialog
    v-model:visible="showAddColDialog"
    @save="addColumn"
  />

  <EditColumnDialog
    v-model:visible="showEditColDialog"
    :column="editingCol"
    @save="updateColumn"
  />

  <AddRowDialog
    v-model:visible="showAddRowDialog"
    :columns="activeTable?.columns"
    @save="addRow"
  />

  <EditRowDialog
    v-model:visible="showEditRowDialog"
    :columns="activeTable?.columns"
    :row="editingRowData"
    @save="updateRow"
  />
</template>

<script setup>
import Toast from 'primevue/toast';
import { useDatabase } from '../composables/useDatabase.js';

// Components
import DbSidebar from '../components/DbSidebar.vue';
import DbDetails from '../components/DbDetails.vue';
import TableDetails from '../components/TableDetails.vue';
import CreateDbDialog from '../components/CreateDbDialog.vue';
import CreateTableDialog from '../components/CreateTableDialog.vue';
import AddColumnDialog from '../components/AddColumnDialog.vue';
import EditColumnDialog from '../components/EditColumnDialog.vue';
import AddRowDialog from '../components/AddRowDialog.vue';
import EditRowDialog from '../components/EditRowDialog.vue';

const {
  databases,
  activeDatabase,
  activeTable,
  showCreateDbDialog,
  showCreateTableDialog,
  showAddColDialog,
  showEditColDialog,
  showAddRowDialog,
  showEditRowDialog,
  editingCol,
  editingRowData,
  back,
  selectDatabase,
  saveDatabaseName,
  dropDatabase,
  selectTable,
  createTable,
  saveTableName,
  dropTable,
  createDatabase,
  addColumn,
  openEditColumn,
  updateColumn,
  deleteSelectedColumns,
  openAddRow,
  addRow,
  openEditRow,
  updateRow,
  deleteSelectedRows
} = useDatabase();
</script>
