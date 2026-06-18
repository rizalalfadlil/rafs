import { ref, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'

export function useDatabase() {
  const toast = useToast();

  const databases = ref([]);
  const activeDatabase = ref(null);
  const activeTable = ref(null);

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
  };

  // Check if a column has an auto-incrementing property
  const isAutoIncrement = (col) => {
    if (!col) return false;
    const type = (col.type || '').toLowerCase();
    const def = (col.default || '').toLowerCase();
    return type.includes('serial') || def.includes('nextval') || (col.name === 'id' && type.includes('int') && def.includes('nextval'));
  };

  // Fetch databases list from API
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

  // Select DB and Fetch Tables
  const selectDatabase = async (db) => {
    activeDatabase.value = db;
    activeDatabase.value.newNameInput = db.name;
    activeTable.value = null;
    await fetchTables(db.name);
  };

  const fetchTables = async (dbName) => {
    try {
      const res = await fetch(`/api/tables?db_name=${encodeURIComponent(dbName)}`);
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
      const colRes = await fetch(`/api/columns?db_name=${encodeURIComponent(dbName)}&table_name=${encodeURIComponent(tableName)}`);
      const colData = await colRes.json();
      
      const rowRes = await fetch(`/api/rows?db_name=${encodeURIComponent(dbName)}&table_name=${encodeURIComponent(tableName)}`);
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

  // Database: Rename
  const saveDatabaseName = async () => {
    if (!activeDatabase.value) return;
    const oldName = activeDatabase.value.name;
    const newName = activeDatabase.value.newNameInput;
    if (!newName || oldName === newName) return;

    try {
      const res = await fetch('/api/databases', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ old_name: oldName, new_name: newName })
      });
      const data = await res.json();
      if (data.status === 'sukses') {
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

  // Database: Drop
  const dropDatabase = async () => {
    if (!activeDatabase.value) return;
    if (!confirm(`Apakah Anda yakin ingin menghapus database '${activeDatabase.value.name}' beserta seluruh isinya?`)) {
      return;
    }
    try {
      const res = await fetch(`/api/databases?db_name=${encodeURIComponent(activeDatabase.value.name)}`, {
        method: 'DELETE'
      });
      const data = await res.json();
      if (data.status === 'sukses') {
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
        headers: { 'Content-Type': 'application/json' },
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
        headers: { 'Content-Type': 'application/json' },
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
        method: 'DELETE'
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
        headers: { 'Content-Type': 'application/json' },
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
        headers: { 'Content-Type': 'application/json' },
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
          method: 'DELETE'
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
        headers: { 'Content-Type': 'application/json' },
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
        headers: { 'Content-Type': 'application/json' },
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
          headers: { 'Content-Type': 'application/json' },
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

  onMounted(() => {
    fetchDatabases();
  });

  return {
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
    originalRowData,
    back,
    isAutoIncrement,
    fetchDatabases,
    selectDatabase,
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
    deleteSelectedRows
  };
}
