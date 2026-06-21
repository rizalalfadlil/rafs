<template>
  <Toast />

  <div class="absolute inset-0 flex flex-col">
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
    <div class="flex grow min-h-0">
      <!-- Sidebar Databases List -->
      <DbSidebar
        :databases="databases"
        :activeDatabase="activeDatabase"
        @select-db="selectDatabase"
        @create-db-click="showCreateDbDialog = true"
        @back="back"
      />

      <!-- Right Panels -->
      <div v-if="activeDatabase" class="p-6 flex-1 min-w-0 overflow-y-auto space-y-6">
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
        <div v-else class="flex flex-col items-center justify-center w-full h-[300px] border border-dashed border-white/10 rounded-lg">
          <i class="pi pi-table text-4xl opacity-20 mb-2"></i>
          <p class="text-base opacity-40">Pilih tabel di atas atau buat tabel baru</p>
        </div>
      </div>

      <!-- Database Not Selected State -->
      <div v-else class="hidden sm:flex flex-col items-center justify-center flex-1 min-w-0 border border-dashed border-white/10 rounded-lg m-6 bg-white/5 select-none">
        <i class="pi pi-database text-6xl opacity-20 mb-4 animate-pulse"></i>
        <p class="text-lg opacity-40">Silakan pilih database di menu sidebar untuk mengelola data</p>
      </div>
    </div>

    <!-- SQL Command Terminal Footer -->
    <footer class="p-4 border-t border-white/10 space-y-4 backdrop-blur-md">
      <div class="flex justify-between items-center">
        <p class="font-semibold text-sm flex items-center gap-2">
          <Button size="small" class="p-0 h-6 w-6" variant="text" :icon="expand ? 'pi pi-chevron-down' : 'pi pi-chevron-up'" @click="expand = !expand" />
          <i class="pi pi-code"></i>
          <span>SQL Query Terminal</span>
        </p>
        <div v-if="activeDatabase" class="text-xs text-white/50 flex items-center gap-1.5 select-none">
          <span class="w-1.5 h-1.5 rounded-full bg-green-500 animate-pulse"></span>
          <span>Koneksi Aktif: <strong>{{ activeDatabase.name }}</strong> (User: {{ activeDbUser }})</span>
          <Button size="small" icon="pi pi-sign-out" variant="text" severity="danger" class="h-6 w-6 p-0 ml-2" title="Keluar dari Database" @click="logoutDatabase" />
        </div>
      </div>

      <!-- SQL Command Input and Output -->
      <div v-if="activeDatabase" class="space-y-4">
        <div class="grid sm:flex gap-4 items-end">
          <Textarea 
            v-model="sqlQuery" 
            placeholder="Tulis perintah SQL di sini... (Tekan Ctrl + Enter untuk menjalankan)" 
            style="font-family: monospace;"  
            class="grow resize-none bg-white/5 border border-white/10 text-white p-3 rounded" 
            :rows="expand ? 6 : 1" 
            @keydown.ctrl.enter="runSqlCommand"
            :disabled="queryRunning"
          />
          <Button 
            :label="queryRunning ? 'Running...' : 'Run'" 
            :icon="queryRunning ? 'pi pi-spin pi-spinner' : 'pi pi-play'" 
            class="h-fit py-3 px-6" 
            @click="runSqlCommand" 
            :disabled="queryRunning"
          />
        </div>

        <!-- SQL Output Panel -->
        <div v-if="queryResult && expand" class="border border-white/10 rounded bg-black/40 p-4 max-h-60 overflow-y-auto space-y-2">
          <div class="flex justify-between items-center border-b border-white/5 pb-2 mb-2 select-none">
            <span class="text-xs font-bold text-zinc-500 uppercase tracking-wider">Hasil Eksekusi</span>
            <Button icon="pi pi-times" variant="text" class="h-5 w-5 p-0" @click="queryResult = null" />
          </div>

          <!-- DDL/DML Exec message -->
          <div v-if="queryResult.status === 'sukses' && queryResult.type === 'exec'" class="text-green-400 text-sm font-mono flex items-center gap-2">
            <i class="pi pi-check-circle"></i>
            <span>{{ queryResult.message }} ({{ queryResult.rows_affected }} baris terpengaruh)</span>
          </div>

          <!-- SELECT Results Table -->
          <div v-else-if="queryResult.status === 'sukses' && queryResult.type === 'select'" class="overflow-x-auto">
            <table v-if="queryResult.rows.length > 0" class="w-full border-collapse text-left font-mono text-xs text-zinc-300">
              <thead>
                <tr class="bg-white/5 text-white border-b border-white/10">
                  <th v-for="col in queryResult.columns" :key="col" class="p-2 border-r border-white/5">{{ col }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(row, rIdx) in queryResult.rows" :key="rIdx" class="border-b border-white/5 hover:bg-white/5">
                  <td v-for="col in queryResult.columns" :key="col" class="p-2 border-r border-white/5 max-w-[200px] truncate" :title="row[col]">
                    {{ row[col] === null ? 'NULL' : row[col] }}
                  </td>
                </tr>
              </tbody>
            </table>
            <p v-else class="text-zinc-500 italic text-sm font-mono p-2">Query berhasil dijalankan, tetapi tidak ada data yang dikembalikan.</p>
          </div>

          <!-- Query Error -->
          <div v-else-if="queryResult.status === 'error'" class="text-red-400 text-sm font-mono p-2 flex items-start gap-2 bg-red-500/10 rounded border border-red-500/20">
            <i class="pi pi-times-circle mt-0.5"></i>
            <div>
              <div class="font-bold">Error Eksekusi:</div>
              <div class="mt-1 whitespace-pre-wrap">{{ queryResult.message }}</div>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="text-center text-xs text-white/30 py-2 select-none">
        Pilih database terlebih dahulu untuk membuka SQL Query Terminal.
      </div>
    </footer>
  </div>

  <!-- Modals & Dialogs -->
  <CreateDbDialog
    v-model:visible="showCreateDbDialog"
    @save="createDatabase"
  />

  <LoginDbDialog
    v-model:visible="showLoginDbDialog"
    :dbName="dbToLogin?.name"
    @login="loginDatabase"
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
import LoginDbDialog from '../components/LoginDbDialog.vue';
import CreateTableDialog from '../components/CreateTableDialog.vue';
import AddColumnDialog from '../components/AddColumnDialog.vue';
import EditColumnDialog from '../components/EditColumnDialog.vue';
import AddRowDialog from '../components/AddRowDialog.vue';
import EditRowDialog from '../components/EditRowDialog.vue';
import Textarea from 'primevue/textarea';
import Button from 'primevue/button';
import { ref } from 'vue';

let expand = ref(false);

const {
  databases,
  activeDatabase,
  activeTable,
  activeDbUser,
  showLoginDbDialog,
  dbToLogin,
  sqlQuery,
  queryResult,
  queryRunning,
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
  loginDatabase,
  logoutDatabase,
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
  deleteSelectedRows,
  runSqlCommand
} = useDatabase();
</script>
