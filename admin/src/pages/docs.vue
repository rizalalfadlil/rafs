<script setup>
import { ref, onMounted } from "vue";
import Button from "primevue/button";
import { marked } from "marked";

// Configure marked to render code blocks and links safely
marked.setOptions({
  gfm: true,
  breaks: true,
});

const activeDoc = ref(null); // 'static_site', 'database', 'storage', 'api'
const docContent = ref("");
const loading = ref(false);
const errorMsg = ref("");
const copiedIndex = ref(null);

const docsList = [
  {
    id: "static_site",
    label: "Static Site",
    icon: "pi pi-globe",
    file: "static_site.md",
    color: "from-blue-500/20 to-cyan-500/10 border-blue-500/30 text-blue-400",
    desc: "Pelajari cara melakukan hosting halaman HTML statis, clone repositori Git publik, dan upload berkas ZIP serta sistem proteksinya.",
  },
  {
    id: "database",
    label: "Database Management",
    icon: "pi pi-database",
    file: "database.md",
    color:
      "from-purple-500/20 to-indigo-500/10 border-purple-500/30 text-purple-400",
    desc: "Panduan lengkap pengelolaan PostgreSQL, pembuatan tabel/skema kolom, serta manipulasi data baris (CRUD) secara GUI.",
  },
  {
    id: "storage",
    label: "Cloud Storage",
    icon: "pi pi-box",
    file: "storage.md",
    color:
      "from-emerald-500/20 to-teal-500/10 border-emerald-500/30 text-emerald-400",
    desc: "Kelola berkas secara pribadi dengan batas kuota 1 GB, fitur multi-upload berkas, pratinjau native, serta publikasi tautan.",
  },
  {
    id: "api",
    label: "API Reference",
    icon: "pi pi-book",
    file: "api.md",
    color:
      "from-amber-500/20 to-orange-500/10 border-amber-500/30 text-amber-400",
    desc: "Spesifikasi teknis endpoint REST API server RAFS untuk integrasi aplikasi eksternal, lengkap dengan payload JSON dan respons.",
  },
];

// Preprocessor for Github-style markdown alerts
function preprocessMarkdown(md) {
  if (!md) return "";
  const lines = md.split("\n");
  let inAlert = false;
  let alertType = "";
  let alertContent = [];
  const processedLines = [];

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    const alertMatch = line.match(
      /^>\s+\[!(NOTE|TIP|IMPORTANT|WARNING|CAUTION)\]\s*$/i,
    );

    if (alertMatch) {
      if (inAlert) {
        processedLines.push(
          `<div class="markdown-alert markdown-alert-${alertType}">`,
        );
        processedLines.push(marked.parse(alertContent.join("\n")));
        processedLines.push(`</div>`);
      }
      inAlert = true;
      alertType = alertMatch[1].toLowerCase();
      alertContent = [];
      continue;
    }

    if (inAlert) {
      if (line.startsWith(">")) {
        alertContent.push(line.replace(/^>\s?/, ""));
      } else {
        processedLines.push(
          `<div class="markdown-alert markdown-alert-${alertType}">`,
        );
        processedLines.push(marked.parse(alertContent.join("\n")));
        processedLines.push(`</div>`);
        inAlert = false;
        processedLines.push(line);
      }
    } else {
      processedLines.push(line);
    }
  }

  if (inAlert) {
    processedLines.push(
      `<div class="markdown-alert markdown-alert-${alertType}">`,
    );
    processedLines.push(marked.parse(alertContent.join("\n")));
    processedLines.push(`</div>`);
  }

  return processedLines.join("\n");
}

async function selectDoc(docId) {
  activeDoc.value = docId;
  const doc = docsList.find((d) => d.id === docId);
  if (!doc) return;

  loading.value = true;
  errorMsg.value = "";
  docContent.value = "";

  try {
    const res = await fetch(`/admin/docs/${doc.file}`);
    if (!res.ok) {
      throw new Error(`Gagal memuat dokumen: ${res.statusText}`);
    }
    const text = await res.text();

    // Preprocess alerts and then parse using marked
    const preprocessedText = preprocessMarkdown(text);
    docContent.value = marked.parse(preprocessedText);

    // Add copy button helper listener after DOM updates
    setTimeout(attachCopyListeners, 100);
  } catch (err) {
    errorMsg.value = err.message || "Terjadi kesalahan saat memuat dokumen.";
  } finally {
    loading.value = false;
  }
}

function goBack() {
  activeDoc.value = null;
  docContent.value = "";
  errorMsg.value = "";
}

// Logic to attach copy buttons to code blocks dynamically
function attachCopyListeners() {
  const preBlocks = document.querySelectorAll(".markdown-body pre");
  preBlocks.forEach((pre, index) => {
    // Check if copy button is already added
    if (pre.querySelector(".copy-btn")) return;

    pre.style.position = "relative";
    const button = document.createElement("button");
    button.className =
      "copy-btn absolute top-2 right-2 bg-white/10 hover:bg-white/20 text-white/70 hover:text-white px-2 py-1 rounded text-xs transition-all duration-200 border border-white/10";
    button.innerHTML = '<i class="pi pi-copy mr-1 text-[10px]"></i> Copy';

    button.addEventListener("click", async () => {
      const codeElement = pre.querySelector("code");
      if (codeElement) {
        try {
          await navigator.clipboard.writeText(codeElement.innerText);
          button.innerHTML =
            '<i class="pi pi-check mr-1 text-[10px] text-green-400"></i> Copied!';
          button.classList.add("bg-green-500/20", "border-green-500/40");

          setTimeout(() => {
            button.innerHTML =
              '<i class="pi pi-copy mr-1 text-[10px]"></i> Copy';
            button.classList.remove("bg-green-500/20", "border-green-500/40");
          }, 2000);
        } catch (e) {
          console.error("Failed to copy", e);
        }
      }
    });

    pre.appendChild(button);
  });
}
</script>

<template>
  <div class="h-full flex flex-col select-none">
    <!-- HEADER -->
    <header
      class="p-4 border-b border-white/5 bg-white/5 backdrop-blur-md sticky top-0 z-10 flex gap-4 items-center"
    >
    <Button label="Kembali" icon="pi pi-arrow-left" @click="goBack" />
      <h1>
        {{
          activeDoc
            ? docsList.find((d) => d.id === activeDoc)?.label
            : "Dokumentasi & Panduan"
        }}
      </h1> 
    </header>

    <!-- CONTENT WRAPPER -->
    <div class="flex-1 flex min-h-0 relative">
      <!-- VIEW A: SELECTION GRID (When no doc is active) -->
      <main
        v-if="!activeDoc"
        class="flex-1 p-8 max-w-6xl mx-auto overflow-y-auto w-full"
      >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mt-4">
          <div
            v-for="doc in docsList"
            :key="doc.id"
            @click="selectDoc(doc.id)"
            class="group relative flex flex-col p-6 rounded-lg border border-white/10 hover:bg-white/5 cursor-pointer transition-all duration-300 transform hover:-translate-y-1 shadow-lg overflow-hidden"
          >
            <!-- Glowing Accent Background -->
            <div
              class="absolute -right-16 -top-16 w-36 h-36 rounded-full blur-3xl opacity-0 group-hover:opacity-100 transition-opacity duration-500"
              :class="
                doc.id === 'static_site'
                  ? 'bg-blue-500/20'
                  : doc.id === 'database'
                    ? 'bg-purple-500/20'
                    : doc.id === 'storage'
                      ? 'bg-emerald-500/20'
                      : 'bg-amber-500/20'
              "
            ></div>

            <div class="flex items-center gap-4 mb-4">
              <!-- Icon Frame -->
              <div
                class="w-12 h-12 rounded-lg flex items-center justify-center border shadow-inner transition-transform duration-300 group-hover:scale-110"
                :class="doc.color"
              >
                <i :class="[doc.icon, 'text-xl']"></i>
              </div>
              <h2
                class="text-lg font-bold group-hover:text-white transition-colors duration-200"
              >
                {{ doc.label }}
              </h2>
            </div>

            <p class="text-smleading-relaxed flex-1">
              {{ doc.desc }}
            </p>

            <div
              class="mt-6 pt-4 border-t border-white/5 flex items-center justify-between text-xs font-semibold"
              :class="
                doc.id === 'static_site'
                  ? 'text-blue-400'
                  : doc.id === 'database'
                    ? 'text-purple-400'
                    : doc.id === 'storage'
                      ? 'text-emerald-400'
                      : 'text-amber-400'
              "
            >
              <span>BACA PANDUAN</span>
              <i
                class="pi pi-arrow-right transform group-hover:translate-x-1 transition-transform duration-200"
              ></i>
            </div>
          </div>
        </div>
      </main>

      <!-- VIEW B: DOCUMENT READER (When doc is active) -->
      <div v-else class="flex-1 flex min-h-0 w-full">
        <!-- Reader Sidebar Navigation -->
        <aside
          class="w-64 border-e border-white/5 p-4 hidden md:flex flex-col gap-1.5 select-none"
        >
          <span
            class="text-[10px] font-bold tracking-wider uppercase px-3 mb-2 block"
            >Daftar Panduan</span
          >
          <div
            v-for="doc in docsList"
            :key="doc.id"
            @click="selectDoc(doc.id)"
            :class="[
              'flex items-center gap-3 px-3.5 py-3 rounded cursor-pointer text-sm font-medium transition-all duration-200 border',
              activeDoc === doc.id
                ? 'bg-white/5 border-white/10 text-white'
                : 'border-transparent hover:text-white hover:bg-white/5',
            ]"
          >
            <i :class="[doc.icon]"></i>
            <span>{{ doc.label }}</span>
          </div>

          <div class="mt-auto pt-4 border-t border-white/5">
            <Button
              label="Kembali ke Menu"
              icon="pi pi-home"
              class="w-full p-button-sm"
              severity="secondary"
              variant="outlined"
              @click="goBack"
            />
          </div>
        </aside>

        <!-- Main Document Reading Area -->
        <main class="flex-1 p-6 md:p-10 overflow-y-auto flex justify-center">
          <div class="max-w-3xl w-full">
            <!-- Loading State -->
            <div v-if="loading" class="space-y-6 py-8">
              <div class="h-8 bg-zinc-900 rounded-lg animate-pulse w-3/4"></div>
              <div
                class="h-4 bg-zinc-900 rounded-lg animate-pulse w-full"
              ></div>
              <div class="h-4 bg-zinc-900 rounded-lg animate-pulse w-5/6"></div>
              <div class="h-4 bg-zinc-900 rounded-lg animate-pulse w-2/3"></div>
              <div class="space-y-3 pt-6">
                <div
                  class="h-32 bg-zinc-900 rounded-lg animate-pulse w-full"
                ></div>
              </div>
            </div>

            <!-- Error State -->
            <div
              v-else-if="errorMsg"
              class="p-6 bg-red-500/10 border border-red-500/20 rounded-lg text-center my-8"
            >
              <i
                class="pi pi-exclamation-triangle text-3xl text-red-400 mb-3 block"
              ></i>
              <p class="text-sm font-medium text-red-200 mb-4">
                {{ errorMsg }}
              </p>
              <Button
                label="Coba Lagi"
                icon="pi pi-refresh"
                severity="danger"
                size="small"
                @click="selectDoc(activeDoc)"
              />
            </div>

            <!-- Rendered Document Content -->
            <article
              v-else
              class="markdown-body py-4"
              v-html="docContent"
            ></article>
          </div>
        </main>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* CUSTOM MARKDOWN RENDER STYLING */
.markdown-body :deep(h1) {
  font-size: 2rem;
  font-weight: 800;
  color: #ffffff;
  margin-top: 2rem;
  margin-bottom: 1.25rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  padding-bottom: 0.75rem;
  letter-spacing: -0.025em;
}

.markdown-body :deep(h2) {
  font-size: 1.5rem;
  font-weight: 700;
  color: #f3f4f6;
  margin-top: 2rem;
  margin-bottom: 1rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  padding-bottom: 0.35rem;
  letter-spacing: -0.015em;
}

.markdown-body :deep(h3) {
  font-size: 1.2rem;
  font-weight: 600;
  color: #e5e7eb;
  margin-top: 1.5rem;
  margin-bottom: 0.75rem;
}

.markdown-body :deep(p) {
  color: #a1a1aa; /* zinc-400 */
  line-height: 1.7;
  margin-bottom: 1.15rem;
  font-size: 0.95rem;
}

.markdown-body :deep(ul) {
  margin-bottom: 1.25rem;
  padding-left: 1.5rem;
}

.markdown-body :deep(ol) {
  margin-bottom: 1.25rem;
  padding-left: 1.5rem;
}

.markdown-body :deep(li) {
  color: #a1a1aa;
  margin-bottom: 0.4rem;
  list-style-type: disc;
  font-size: 0.95rem;
}

.markdown-body :deep(ol li) {
  list-style-type: decimal;
}

.markdown-body :deep(code) {
  font-family:
    ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono",
    "Courier New", monospace;
  background-color: rgba(255, 255, 255, 0.08);
  padding: 0.15rem 0.35rem;
  border-radius: 0.375rem;
  font-size: 0.85rem;
  color: #fbcfe8; /* pink-200 */
  border: 1px solid rgba(255, 255, 255, 0.04);
}

.markdown-body :deep(pre) {
  background-color: #09090b; /* zinc-950 */
  border: 1px solid rgba(255, 255, 255, 0.06);
  padding: 1.15rem;
  border-radius: 0.75rem;
  overflow-x: auto;
  margin-top: 1rem;
  margin-bottom: 1.25rem;
}

.markdown-body :deep(pre code) {
  background-color: transparent;
  padding: 0;
  border-radius: 0;
  color: #34d399; /* emerald-400 */
  font-size: 0.85rem;
  border: none;
}

.markdown-body :deep(blockquote) {
  border-left: 4px solid rgba(255, 255, 255, 0.15);
  padding-left: 1.15rem;
  color: #71717a; /* zinc-500 */
  font-style: italic;
  margin-top: 1rem;
  margin-bottom: 1.25rem;
}

.markdown-body :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin-top: 1.25rem;
  margin-bottom: 1.75rem;
  font-size: 0.9rem;
}

.markdown-body :deep(th) {
  background-color: rgba(255, 255, 255, 0.03);
  color: #ffffff;
  font-weight: 600;
  text-align: left;
  padding: 0.85rem 1rem;
  border-bottom: 2px solid rgba(255, 255, 255, 0.08);
}

.markdown-body :deep(td) {
  padding: 0.85rem 1rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
  color: #a1a1aa;
}

.markdown-body :deep(hr) {
  border: 0;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  margin: 2rem 0;
}

/* Alert styling matching Github-style markdown alerts */
.markdown-body :deep(.markdown-alert) {
  padding: 1rem 1.15rem;
  margin: 1.25rem 0;
  border-left: 4px solid;
  border-radius: 0.5rem;
  font-size: 0.9rem;
}

.markdown-body :deep(.markdown-alert-note) {
  background-color: rgba(59, 130, 246, 0.06);
  border-color: #3b82f6;
  color: #93c5fd;
}
.markdown-body :deep(.markdown-alert-note p) {
  color: #bfdbfe;
  margin-bottom: 0;
}

.markdown-body :deep(.markdown-alert-tip) {
  background-color: rgba(16, 185, 129, 0.06);
  border-color: #10b981;
  color: #6ee7b7;
}
.markdown-body :deep(.markdown-alert-tip p) {
  color: #a7f3d0;
  margin-bottom: 0;
}

.markdown-body :deep(.markdown-alert-important) {
  background-color: rgba(139, 92, 246, 0.06);
  border-color: #8b5cf6;
  color: #c4b5fd;
}
.markdown-body :deep(.markdown-alert-important p) {
  color: #ddd6fe;
  margin-bottom: 0;
}

.markdown-body :deep(.markdown-alert-warning) {
  background-color: rgba(245, 158, 11, 0.06);
  border-color: #f59e0b;
  color: #fde047;
}
.markdown-body :deep(.markdown-alert-warning p) {
  color: #fef08a;
  margin-bottom: 0;
}

.markdown-body :deep(.markdown-alert-caution) {
  background-color: rgba(239, 68, 68, 0.06);
  border-color: #ef4444;
  color: #fca5a5;
}
.markdown-body :deep(.markdown-alert-caution p) {
  color: #fecaca;
  margin-bottom: 0;
}

/* Custom Scrollbar */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}
::-webkit-scrollbar-track {
  background: transparent;
}
::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 9999px;
}
::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.2);
}
</style>
