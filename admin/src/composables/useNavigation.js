import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'

export function useNavigation () {
const route = useRoute();

const menuItems = ref([
  { label: 'Dashboard', path: '/admin/', icon: 'pi pi-home' },
  { label: 'Static Sites', path: '/admin/sites', icon: 'pi pi-globe' },
  { label: 'Database Management', path: '/admin/databases', icon: 'pi pi-database' },
  { label: 'File Storage', path: '/admin/storage', icon: 'pi pi-box' },
  { label: 'Documentation', path: '/admin/docs', icon: 'pi pi-book' },
])

const activeMenu = computed(() => {
  const currentPath = route.path.replace(/\/$/, '') || '/';
  const item = menuItems.value.find(m => {
    const itemPath = m.path.replace(/\/$/, '') || '/';
    return itemPath === currentPath;
  });
  if (item) return item.label;

  const segment = route.path.split('/').filter(Boolean).pop();
  if (!segment || segment === 'admin') return 'Dashboard';
  return segment
    .split(/[-_]/)
    .map(word => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ');
});

return {
    activeMenu, menuItems
}
}