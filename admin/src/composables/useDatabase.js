import { ref, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'

export function useDatabase() {
  const toast = useToast();

  const databases = ref([]);
  const activeDatabase = ref(null);
  const activeTable = ref(null);

  // Database credentials & login state
  const activeDbUser = ref('');
  const activeDbPassword = ref('');
  const showLoginDbDialog = ref(false);
  const dbToLogin = ref(null);

  // SQL Query Runner state
  const sqlQuery = ref('');
  const queryResult = ref(null);
  const queryRunning = ref(false);

  // Dialog visible states
  const showCreateDbDialog = ref(false);
  const showCreateTableDialog = ref(false);
  const showAddColDialog = ref(false);
  const showEditColDialog = ref(false);
  const showAddRowDialog = ref(false);
  const showEditRowDialog = ref(false);

  // Refs for data transfer
  const editingCol = ref({});
  const editingRowData = ref({});
  const originalRowData = ref({});

  // Reset view states
  const back = () => {
    activeDatabase.value = null;
    activeTable.value = null;
    activeDbUser.value = '';
    activeDbPassword.value = '';
  };

  // Helper to generate auth headers
  const getAuthHeaders = () => {
    return {
      'X-Database-User': activeDbUser.value || '',
      'X-Database-Password': activeDbPassword.value || ''
    };
  };

  // Check if a column has an auto-incrementing property
  const isAutoIncrement = (col) => {
    if (!col) return false;
    const type = (col.type || '').toLowerCase();
    const def = (col.default || '').toLowerCase();
    return type.includes('serial') || def.includes('nextval') || (col.name === 'id' && type.includes('int') && def.includes('nextval'));
  };

  // Fetch databases list from API (uses superadmin credentials inside Go server by default)
  const fetchDatabases = async () => {
    try {
      const res = await fetch('/api/databases');
      const data = await res.json();
      if (data.status === 'sukses') {
        databases.value = data.databases.map(dbName => ({
          name: dbName,
          tables: [],
          newNameInput: dbName
        }));
      } else {
        toast.add({ severity: 'error', summary: 'Gagal', detail: data.message, life: 3000 });
      }
    } catch (err) {
      toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000 });
    }
  };

  // Select DB and verify or prompt credentials
  const selectDatabase = async (db) => {
    const cached = localStorage.getItem(`rafs_db_cred_${db.name}`);
    if (cached) {
      try {
        const cred = JSON.parse(cached);
        activeDbUser.value = cred.username;
        activeDbPassword.value = cred.password;
        activeDatabase.value = db;
        activeDatabase.value.newNameInput = db.name;
        activeTable.value = null;
        
        // Coba memuat tabel untuk validasi
        await fetchTables(db.name);
        return;
      } catch (e) {
        localStorage.removeItem(`rafs_db_cred_${db.name}`);
      }
    }
    
    // Kredensial tidak ditemukan, tampilkan dialog login
    dbToLogin.value = db;
    showLoginDbDialog.value = true;
  };

  // Login Database
  const loginDatabase = async (credentials) => {
    try {
      const res = await fetch(`/api/tables?db_name=${encodeURIComponent(dbToLogin.value.name)}`, {
        headers: {
          'X-Database-User': credentials.username,
          'X-Database-Password': credentials.password
        }
      });
      const data = await res.json();
      if (res.ok && data.status === 'sukses') {
        localStorage.setItem(`rafs_db_cred_${dbToLogin.value.name}`, JSON.stringify(credentials));
        activeDbUser.value = credentials.username;
        activeDbPassword.value = credentials.password;
        
        activeDatabase.value = dbToLogin.value;
        activeDatabase.value.newNameInput = dbToLogin.value.name;
        activeTable.value = null;
        
        const dbObj = databases.value.find(d => d.name === dbToLogin.value.name);
        if (dbObj) {
          dbObj.tables = data.tables.map(tableName => ({
            name: tableName,
            columns: [],
            rows: [],
            newNameInput: tableName
          }));
        }
        
        showLoginDbDialog.value = false;
        dbToLogin.value = null;
        toast.add({ severity: 'success', summary: 'Akses Diterima', detail: 'Berhasil masuk ke database!', life: 3000 });
      } else {
        toast.add({ severity: 'error', summary: 'Akses Ditolak', detail: data.message || 'Username atau password salah', life: 3000 });
      }
    } catch (err) {
      toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal terhubung ke database: ' + err.message, life: 3000 });
    }
  };

  const logoutDatabase = () => {
    if (activeDatabase.value) {
      localStorage.removeItem(`rafs_db_cred_${activeDatabase.value.name}`);
    }
    activeDatabase.value = null;
    activeTable.value = null;
    activeDbUser.value = '';
    activeDbPassword.value = '';
  };

  // Fetch Tables
  const fetchTables = async (dbName) => {
    try {
      const res = await fetch(`/api/tables?db_name=${encodeURIComponent(dbName)}`, {
        headers: getAuthHeaders()
      });
      const data = await res.json();
      if (data.status === 'sukses') {
        const db = databases.value.find(d => d.name === dbName);
        if (db) {
          db.tables = data.tables.map(tableName => ({
            name: tableName,
            columns: [],
            rows: [],
            newNameInput: tableName
          }));
        }
      } else {
        toast.add({ severity: 'error', summary: 'Gagal', detail: data.message, life: 3000 });
      }
    } catch (err) {
      toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000 });
    }
  };

  // Select Table and Fetch details
  const selectTable = async (table) => {
    activeTable.value = table;
    activeTable.value.newNameInput = table.name;
    await fetchColumnsAndRows(activeDatabase.value.name, table.name);
  };

  const fetchColumnsAndRows = async (dbName, tableName) => {
    try {
      const colRes = await fetch(`/api/columns?db_name=${encodeURIComponent(dbName)}&table_name=${encodeURIComponent(tableName)}`, {
        headers: getAuthHeaders()
      });
      const colData = await colRes.json();
      
      const rowRes = await fetch(`/api/rows?db_name=${encodeURIComponent(dbName)}&table_name=${encodeURIComponent(tableName)}`, {
        headers: getAuthHeaders()
      });
      const rowData = await rowRes.json();

      if (colData.status === 'sukses' && rowData.status === 'sukses') {
        const db = databases.value.find(d => d.name === dbName);
        if (db) {
          const table = db.tables.find(t => t.name === tableName);
          if (table) {
            table.columns = colData.columns.map(c => ({ ...c, checked: false }));
            table.rows = rowData.rows.map(r => ({ ...r, _selected: false }));
            activeTable.value = table;
          }
        }
      } else {
        const msg = colData.message || rowData.message;
        toast.add({ severity: 'error', summary: 'Gagal Load', detail: msg, life: 3000 });
      }
    } catch (err) {
      toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000 });
    }
  };

  // Database: Create
  const createDatabase = async (payload) => {
    try {
      const res = await fetch('/api/databases', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });
      const data = await res.json();
      if (data.status === 'sukses') {
        // Cache credentials so they can access it immediately
        localStorage.setItem(`rafs_db_cred_${payload.db_name}`, JSON.stringify({
          username: payload.username,
          password: payload.password
        }));

        toast.add({ severity: 'success', summary: 'Sukses', detail: data.message, life: 3000 });
        showCreateDbDialog.value = false;
        await fetchDatabases();
      } else {
        toast.add({ severity: 'error', summary: 'Gagal', detail: data.message, life: 3000 });
      }
    } catch (err) {
      toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000 });
    }
  };

  // Database: Rename (uses provided credentials to authorize first)
  const saveDatabaseName = async () => {
    if (!activeDatabase.value) return;
    const oldName = activeDatabase.value.name;
    const newName = activeDatabase.value.newNameInput;
    if (!newName || oldName === newName) return;

    try {
      const res = await fetch('/api/databases', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          ...getAuthHeaders()
        },
        body: JSON.stringify({ old_name: oldName, new_name: newName })
      });
      const data = await res.json();
      if (data.status === 'sukses') {
        // Update credentials key
        const cred = localStorage.getItem(`rafs_db_cred_${oldName}`);
        if (cred) {
          localStorage.setItem(`rafs_db_cred_${newName}`, cred);
          localStorage.removeItem(`rafs_db_cred_${oldName}`);
        }

        toast.add({ severity: 'success', summary: 'Sukses', detail: data.message, life: 3000 });
        await fetchDatabases();
        const updatedDb = databases.value.find(d => d.name === newName);
        if (updatedDb) {
          await selectDatabase(updatedDb);
        }
      } else {
        toast.add({ severity: 'error', summary: 'Gagal', detail: data.message, life: 3000 });
      }
    } catch (err) {
      toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000 });
    }
  };

  // Database: Drop (uses provided credentials to authorize first)
  const dropDatabase = async () => {
    if (!activeDatabase.value) return;
    if (!confirm(`Apakah Anda yakin ingin menghapus database '${activeDatabase.value.name}' beserta seluruh isinya?`)) {
      return;
    }
    try {
      const res = await fetch(`/api/databases?db_name=${encodeURIComponent(activeDatabase.value.name)}`, {
        method: 'DELETE',
        headers: getAuthHeaders()
      });
      const data = await res.json();
      if (data.status === 'sukses') {
        localStorage.removeItem(`rafs_db_cred_${activeDatabase.value.name}`);
        activeDbUser.value = '';
        activeDbPassword.value = '';

        toast.add({ severity: 'success', summary: 'Sukses', detail: data.message, life: 3000 });
        activeDatabase.value = null;
        activeTable.value = null;
        await fetchDatabases();
      } else {
        toast.add({ severity: 'error', summary: 'Gagal', detail: data.message, life: 3000 });
      }
    } catch (err) {
      toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000 });
    }
  };

  // Table: Create
  const createTable = async (payload) => {
    try {
      const res = await fetch('/api/tables', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...getAuthHeaders()
        },
        body: JSON.stringify({
          db_name: activeDatabase.value.name,
          table_name: payload.table_name,
          columns: payload.columns
        })
      });
      const data = await res.json();
      if (data.status === 'sukses') {
        toast.add({ severity: 'success', summary: 'Sukses', detail: data.message, life: 3000 });
        showCreateTableDialog.value = false;
        await fetchTables(activeDatabase.value.name);
      } else {
        toast.add({ severity: 'error', summary: 'Gagal', detail: data.message, life: 3000 });
      }
    } catch (err) {
      toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000 });
    }
  };

  // Table: Rename
  const saveTableName = async () => {
    if (!activeDatabase.value || !activeTable.value) return;
    const oldName = activeTable.value.name;
    const newName = activeTable.value.newNameInput;
    if (!newName || oldName === newName) return;

    try {
      const res = await fetch('/api/tables', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          ...getAuthHeaders()
        },
        body: JSON.stringify({
          db_name: activeDatabase.value.name,
          old_name: oldName,
          new_name: newName
        })
      });
      const data = await res.json();
      if (data.status === 'sukses') {
        toast.add({ severity: 'success', summary: 'Sukses', detail: data.message, life: 3000 });
        await fetchTables(activeDatabase.value.name);
        const db = databases.value.find(d => d.name === activeDatabase.value.name);
        if (db) {
          const table = db.tables.find(t => t.name === newName);
          if (table) {
            await selectTable(table);
          }
        }
      } else {
        toast.add({ severity: 'error', summary: 'Gagal', detail: data.message, life: 3000 });
      }
    } catch (err) {
      toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000 });
    }
  };

  // Table: Drop
  const dropTable = async () => {
    if (!activeDatabase.value || !activeTable.value) return;
    if (!confirm(`Apakah Anda yakin ingin menghapus tabel '${activeTable.value.name}'? Semua data di dalamnya akan hilang.`)) {
      return;
    }
    try {
      const res = await fetch(`/api/tables?db_name=${encodeURIComponent(activeDatabase.value.name)}&table_name=${encodeURIComponent(activeTable.value.name)}`, {
        method: 'DELETE',
        headers: getAuthHeaders()
      });
      const data = await res.json();
      if (data.status === 'sukses') {
        toast.add({ severity: 'success', summary: 'Sukses', detail: data.message, life: 3000 });
        activeTable.value = null;
        await fetchTables(activeDatabase.value.name);
      } else {
        toast.add({ severity: 'error', summary: 'Gagal', detail: data.message, life: 3000 });
      }
    } catch (err) {
      toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000 });
    }
  };

  // Column: Add
  const addColumn = async (payload) => {
    try {
      const res = await fetch('/api/columns', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...getAuthHeaders()
        },
        body: JSON.stringify({
          db_name: activeDatabase.value.name,
          table_name: activeTable.value.name,
          column_name: payload.name,
          column_type: payload.type
        })
      });
      const data = await res.json();
      if (data.status === 'sukses') {
        toast.add({ severity: 'success', summary: 'Sukses', detail: data.message, life: 3000 });
        showAddColDialog.value = false;
        await fetchColumnsAndRows(activeDatabase.value.name, activeTable.value.name);
      } else {
        toast.add({ severity: 'error', summary: 'Gagal', detail: data.message, life: 3000 });
      }
    } catch (err) {
      toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000 });
    }
  };

  // Column: Edit
  const openEditColumn = (column) => {
    editingCol.value = { ...column };
    showEditColDialog.value = true;
  };

  const updateColumn = async (payload) => {
    try {
      const res = await fetch('/api/columns', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          ...getAuthHeaders()
        },
        body: JSON.stringify({
          db_name: activeDatabase.value.name,
          table_name: activeTable.value.name,
          old_name: payload.old_name,
          new_name: payload.new_name,
          column_type: payload.type
        })
      });
      const data = await res.json();
      if (data.status === 'sukses') {
        toast.add({ severity: 'success', summary: 'Sukses', detail: data.message, life: 3000 });
        showEditColDialog.value = false;
        await fetchColumnsAndRows(activeDatabase.value.name, activeTable.value.name);
      } else {
        toast.add({ severity: 'error', summary: 'Gagal', detail: data.message, life: 3000 });
      }
    } catch (err) {
      toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000 });
    }
  };

  // Column: Drop Selected
  const deleteSelectedColumns = async () => {
    const checkedCols = activeTable.value.columns.filter(c => c.checked).map(c => c.name);
    if (checkedCols.length === 0) {
      toast.add({ severity: 'warn', summary: 'Peringatan', detail: 'Pilih kolom yang ingin dihapus', life: 3000 });
      return;
    }
    if (!confirm(`Apakah Anda yakin ingin menghapus kolom: ${checkedCols.join(', ')}?`)) {
      return;
    }

    try {
      let successCount = 0;
      for (const colName of checkedCols) {
        const res = await fetch(`/api/columns?db_name=${encodeURIComponent(activeDatabase.value.name)}&table_name=${encodeURIComponent(activeTable.value.name)}&column_name=${encodeURIComponent(colName)}`, {
          method: 'DELETE',
          headers: getAuthHeaders()
        });
        const data = await res.json();
        if (data.status === 'sukses') {
          successCount++;
        }
      }
      toast.add({ severity: 'success', summary: 'Sukses', detail: `${successCount} kolom berhasil dihapus`, life: 3000 });
      await fetchColumnsAndRows(activeDatabase.value.name, activeTable.value.name);
    } catch (err) {
      toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000 });
    }
  };

  // Row: Add
  const openAddRow = () => {
    showAddRowDialog.value = true;
  };

  const addRow = async (payload) => {
    try {
      const formattedData = {};
      for (const [key, val] of Object.entries(payload)) {
        const col = activeTable.value.columns.find(c => c.name === key);
        if (!col || isAutoIncrement(col)) {
          continue;
        }
        
        if (col.type.toLowerCase() === 'boolean') {
          if (val === 'true' || val === true) formattedData[key] = true;
          else if (val === 'false' || val === false) formattedData[key] = false;
          else formattedData[key] = null;
        } else if (col.type.toLowerCase().includes('int') || col.type.toLowerCase().includes('real') || col.type.toLowerCase().includes('numeric') || col.type.toLowerCase().includes('double') || col.type.toLowerCase().includes('precision')) {
          formattedData[key] = val === '' ? null : Number(val);
        } else {
          formattedData[key] = val === '' ? null : val;
        }
      }

      const res = await fetch('/api/rows', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...getAuthHeaders()
        },
        body: JSON.stringify({
          db_name: activeDatabase.value.name,
          table_name: activeTable.value.name,
          row: formattedData
        })
      });
      const data = await res.json();
      if (data.status === 'sukses') {
        toast.add({ severity: 'success', summary: 'Sukses', detail: data.message, life: 3000 });
        showAddRowDialog.value = false;
        await fetchColumnsAndRows(activeDatabase.value.name, activeTable.value.name);
      } else {
        toast.add({ severity: 'error', summary: 'Gagal', detail: data.message, life: 3000 });
      }
    } catch (err) {
      toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000 });
    }
  };

  // Row: Edit
  const openEditRow = (row) => {
    const orig = {};
    activeTable.value.columns.forEach(col => {
      orig[col.name] = row[col.name];
    });
    originalRowData.value = orig;
    editingRowData.value = { ...row };
    showEditRowDialog.value = true;
  };

  const updateRow = async (payload) => {
    try {
      const formattedData = {};
      for (const [key, val] of Object.entries(payload)) {
        const col = activeTable.value.columns.find(c => c.name === key);
        if (!col || isAutoIncrement(col) || col.name === 'id') {
          continue;
        }
        
        if (col.type.toLowerCase() === 'boolean') {
          if (val === 'true' || val === true) formattedData[key] = true;
          else if (val === 'false' || val === false) formattedData[key] = false;
          else formattedData[key] = null;
        } else if (col.type.toLowerCase().includes('int') || col.type.toLowerCase().includes('real') || col.type.toLowerCase().includes('numeric') || col.type.toLowerCase().includes('double') || col.type.toLowerCase().includes('precision')) {
          formattedData[key] = val === '' ? null : Number(val);
        } else {
          formattedData[key] = val === '' ? null : val;
        }
      }

      const res = await fetch('/api/rows', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          ...getAuthHeaders()
        },
        body: JSON.stringify({
          db_name: activeDatabase.value.name,
          table_name: activeTable.value.name,
          where: originalRowData.value,
          row: formattedData
        })
      });
      const data = await res.json();
      if (data.status === 'sukses') {
        toast.add({ severity: 'success', summary: 'Sukses', detail: data.message, life: 3000 });
        showEditRowDialog.value = false;
        await fetchColumnsAndRows(activeDatabase.value.name, activeTable.value.name);
      } else {
        toast.add({ severity: 'error', summary: 'Gagal', detail: data.message, life: 3000 });
      }
    } catch (err) {
      toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000 });
    }
  };

  // Row: Delete Selected
  const deleteSelectedRows = async () => {
    const selected = activeTable.value.rows.filter(r => r._selected);
    if (selected.length === 0) {
      toast.add({ severity: 'warn', summary: 'Peringatan', detail: 'Pilih baris yang ingin dihapus', life: 3000 });
      return;
    }
    if (!confirm(`Apakah Anda yakin ingin menghapus ${selected.length} baris data?`)) {
      return;
    }

    try {
      let successCount = 0;
      for (const row of selected) {
        const whereClause = {};
        activeTable.value.columns.forEach(col => {
          whereClause[col.name] = row[col.name];
        });

        const res = await fetch('/api/rows', {
          method: 'DELETE',
          headers: {
            'Content-Type': 'application/json',
            ...getAuthHeaders()
          },
          body: JSON.stringify({
            db_name: activeDatabase.value.name,
            table_name: activeTable.value.name,
            where: whereClause
          })
        });
        const data = await res.json();
        if (data.status === 'sukses') {
          successCount++;
        }
      }
      toast.add({ severity: 'success', summary: 'Sukses', detail: `${successCount} baris berhasil dihapus`, life: 3000 });
      await fetchColumnsAndRows(activeDatabase.value.name, activeTable.value.name);
    } catch (err) {
      toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 3000 });
    }
  };

  // SQL Runner: Execute custom query
  const runSqlCommand = async () => {
    if (!activeDatabase.value) {
      toast.add({ severity: 'warn', summary: 'Peringatan', detail: 'Pilih database terlebih dahulu', life: 3000 });
      return;
    }
    if (!sqlQuery.value.trim()) {
      toast.add({ severity: 'warn', summary: 'Peringatan', detail: 'Perintah SQL kosong', life: 3000 });
      return;
    }

    queryRunning.value = true;
    queryResult.value = null;

    try {
      const res = await fetch('/api/query', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...getAuthHeaders()
        },
        body: JSON.stringify({
          db_name: activeDatabase.value.name,
          query: sqlQuery.value
        })
      });
      const data = await res.json();
      if (data.status === 'sukses') {
        queryResult.value = data;
        toast.add({ severity: 'success', summary: 'Sukses', detail: 'Query berhasil dieksekusi!', life: 3000 });
        
        // Auto-refresh layout jika query mengubah DDL/DML
        const lowerQuery = sqlQuery.value.toLowerCase();
        if (
          lowerQuery.includes('table') || 
          lowerQuery.includes('column') || 
          lowerQuery.includes('insert') || 
          lowerQuery.includes('update') || 
          lowerQuery.includes('delete') || 
          lowerQuery.includes('create') || 
          lowerQuery.includes('drop')
        ) {
          await fetchTables(activeDatabase.value.name);
          if (activeTable.value) {
            await fetchColumnsAndRows(activeDatabase.value.name, activeTable.value.name);
          }
        }
      } else {
        queryResult.value = { status: 'error', message: data.message };
        toast.add({ severity: 'error', summary: 'Query Gagal', detail: data.message, life: 4000 });
      }
    } catch (err) {
      queryResult.value = { status: 'error', message: err.message };
      toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 4000 });
    } finally {
      queryRunning.value = false;
      sqlQuery.value = '';
    }
  };

  onMounted(() => {
    fetchDatabases();
  });

  return {
    databases,
    activeDatabase,
    activeTable,
    activeDbUser,
    activeDbPassword,
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
    originalRowData,
    back,
    isAutoIncrement,
    fetchDatabases,
    selectDatabase,
    loginDatabase,
    logoutDatabase,
    fetchTables,
    selectTable,
    fetchColumnsAndRows,
    createDatabase,
    saveDatabaseName,
    dropDatabase,
    createTable,
    saveTableName,
    dropTable,
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
  };
}
