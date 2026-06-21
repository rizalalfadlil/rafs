<script setup>
import { useNavigation } from '../composables/useNavigation';
import Button from 'primevue/button';

const { activeMenu, menuItems } = useNavigation();

defineProps({
    collapse: Boolean
})
</script>

<template>
    <div class="p-8 block sm:hidden fixed bottom-0 w-full transition-all duration-200" :class="collapse ? 'translate-y-full' : ''">
      <div class=" w-fit mx-auto bg-[#0c0f19]/80 backdrop-blur-lg border border-white/5 rounded p-4 flex gap-4 justify-around">
        <Button v-for="item in menuItems" :icon="item.icon" :variant="activeMenu.label === item.label ? 'primary' : 'text'" @click="$router.push(item.path)" />
      </div>
    </div>
</template>