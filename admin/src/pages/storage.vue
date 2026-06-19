<script setup>
import { ref } from 'vue';
const directory = ref([
    {
        name: "My Folder",
        type: "folder",
        contents: [
            {
                name: "Hello World.txt",
                type: "file"
            }
        ]
    },
    {
        name: "Public",
        type: "folder",
        contents: []
    },
    {
        name: "Project X",
        type: "folder",
        contents: [{
            name: "index.html",
            type: "file"
        },
        {
            name: "style.css",
            type: "file"
        },
        {
            name: "script.js",
            type: "file"
        }, {
            name: "Public",
            type: "folder",
            contents: [
                {
                    name: "images", type: "folder", contents: [
                        {
                            name: "001.jpg",
                            type: "file"
                        },
                        {
                            name: "002.jpg",
                            type: "file"
                        },
                        {
                            name: "003.jpg",
                            type: "file"
                        },
                        {
                            name: "004.jpg",
                            type: "file"
                        },
                    ]
                }
            ]
        }
        ]
    },
    {
        name: "README.md",
        type: "file"
    }
]);

let currentPath = ref("");
const paths = ref([]);


</script>
<template>
    <!-- Breadcrumb -->
    <div class="p-2 bg-white/5 border-b border-white/5 px-4 flex gap-4 items-center text-sm text-white/50 select-none">
        <i class="pi pi-box cursor-pointer hover:text-white"></i>
        <i class="pi pi-chevron-right text-sm"></i>
        <p>{{ currentPath }}</p>
    </div>
    <!-- Contents -->
    <main class="p-2 flex gap-2 flex-wrap ">
        <div
            @click="if (item.type === 'folder') { currentPath = currentPath + '/' + item.name; directory = item.contents; } else console.log(item)"
            v-for="item in directory" :key="item.name"
            class="bg-white/5 border border-white/5 p-2 rounded hover:bg-white/10 cursor-pointer grid gap-4 text-center w-32 aspect-square">
            <i class="pi" style="font-size: 4rem" :class="item.type === 'folder' ? 'pi-folder' : 'pi-file'"></i>
            <p>{{ item.name }}</p>
        </div>
    </main>
</template>