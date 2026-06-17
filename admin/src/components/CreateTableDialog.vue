<template>
  <Dialog 
    :visible="visible" 
    header="Buat Tabel Baru" 
    :modal="true" 
    :style="{ width: '600px' }"
    @update:visible="$emit('update:visible', $event)"
  >
    <div class="flex flex-col gap-4 pt-2">
      <div class="flex flex-col gap-1.5">
        <label for="table_name" class="text-xs font-bold text-slate-400">Nama Tabel</label>
        <InputText id="table_name" v-model="form.table_name" placeholder="contoh: produk" required class="!bg-[#0c0f19] !border-slate-800 !text-slate-100" />
      </div>
      
      <div class="flex justify-between items-center border-b border-slate-800 pb-2 mt-2">
        <span class="text-sm font-bold text-slate-100">Definisi Kolom</span>
        <Button 
          icon="pi pi-plus" 
          label="Tambah Kolom" 
          class="p-button-outlined p-button-sm" 
          @click="addColumnField" 
        />
      </div>

      <div class="max-h-[250px] overflow-y-auto flex flex-col gap-2.5 pr-1.5">
        <div v-for="(col, index) in form.columns" :key="index" class="flex gap-2.5 items-center">
          <InputText 
            v-model="col.name" 
            placeholder="Nama Kolom" 
            class="!bg-[#0c0f19] !border-slate-800 !text-slate-100 flex-[3]"
            required 
            :disabled="index === 0" 
          />
          
          <Select 
            v-model="col.selectedType" 
            :options="columnTypeOptions" 
            optionLabel="label" 
            optionValue="value" 
            class="!bg-[#0c0f19] !border-slate-800 !text-slate-100 flex-[3]" 
            :disabled="index === 0"
            @change="onColumnTypeChange(index)"
          />
          
          <InputText 
            v-if="col.selectedType === 'custom'" 
            v-model="col.customType" 
            placeholder="Tipe kustom..." 
            class="!bg-[#0c0f19] !border-slate-800 !text-slate-100 flex-[4]"
            required 
          />

          <Button 
            icon="pi pi-trash" 
            class="p-button-danger p-button-text shrink-0" 
            :disabled="index === 0" 
            @click="removeColumnField(index)" 
          />
        </div>
      </div>
    </div>
    <template #footer>
      <Button label="Batal" class="p-button-text p-button-secondary" @click="$emit('update:visible', false)" />
      <Button label="Buat Tabel" class="p-button-primary" @click="submit" :loading="loading" />
    </template>
  </Dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import Button from 'primevue/button'

const props = defineProps({
  visible: {
    type: Boolean,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:visible', 'submit'])

const form = ref({
  table_name: '',
  columns: []
})

const columnTypeOptions = [
  { label: 'SERIAL PRIMARY KEY', value: 'SERIAL PRIMARY KEY' },
  { label: 'INTEGER', value: 'INTEGER' },
  { label: 'VARCHAR(255)', value: 'VARCHAR(255)' },
  { label: 'TEXT', value: 'TEXT' },
  { label: 'BOOLEAN', value: 'BOOLEAN' },
  { label: 'TIMESTAMP', value: 'TIMESTAMP' },
  { label: 'Tipe Kustom...', value: 'custom' }
]

watch(() => props.visible, (newVal) => {
  if (newVal) {
    form.value = {
      table_name: '',
      columns: [
        { name: 'id', selectedType: 'SERIAL PRIMARY KEY', customType: '' }
      ]
    }
  }
})

const addColumnField = () => {
  form.value.columns.push({ name: '', selectedType: 'VARCHAR(255)', customType: '' })
}

const removeColumnField = (index) => {
  if (index === 0) return
  form.value.columns.splice(index, 1)
}

const onColumnTypeChange = (index) => {
  const col = form.value.columns[index]
  if (col.selectedType !== 'custom') {
    col.customType = ''
  }
}

const submit = () => {
  emit('submit', { ...form.value })
}
</script>
