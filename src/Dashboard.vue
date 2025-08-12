<template>
  <div class="dash-wrap">
    <!-- Sidebar -->
    <aside class="sidebar">
      <button class="side-btn active" @click="$router.push('/dashboard')" aria-label="Home">
        <span class="ico">🏠</span>
      </button>
      <button class="side-btn" @click="logout" aria-label="Logout">
        <span class="ico">↪︎</span>
      </button>
      <button class="side-btn" aria-label="Settings">
        <span class="ico">⚙️</span>
      </button>
    </aside>

    <!-- Main content -->
    <main class="main">
      <!-- Header -->
      <header class="page-header">
        <h1 class="h1">Dashboard</h1>
        <p class="sub">Create new notes</p>
      </header>

      <!-- Quick actions -->
      <section class="quick">
        <div class="qa" v-for="(card, i) in quickActions" :key="i" @click="card.action()">
          <div class="qa-ico" :class="[card.color, card.key]">
            <!-- Blank document -->
            <svg v-if="card.key==='blank'" viewBox="0 0 24 24" aria-hidden="true">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8l-6-6zM15 9V3.5L20.5 9H15z"/>
            </svg>
            <!-- Microphone -->
            <svg v-else-if="card.key==='audio'" viewBox="0 0 24 24" aria-hidden="true">
              <path d="M12 15a3 3 0 0 0 3-3V7a3 3 0 0 0-6 0v5a3 3 0 0 0 3 3zm5-4v1a5 5 0 0 1-10 0v-1H5v1a7 7 0 0 0 6 6v3h2v-3a7 7 0 0 0 6-6v-1h-2z"/>
            </svg>
            <!-- Document upload (DOC tag) -->
            <svg v-else-if="card.key==='doc'" viewBox="0 0 24 24" aria-hidden="true" class="doc-svg">
              <rect class="doc-paper" x="4" y="3.5" width="16" height="17" rx="3" ry="3"/>
              <polygon class="doc-fold" points="16,3.5 20,7.5 16,7.5"/>
              <text x="12" y="15" text-anchor="middle" class="doc-text">DOC</text>
            </svg>
            <!-- YouTube -->
            <svg v-else viewBox="0 0 24 24" aria-hidden="true" class="yt-svg">
              <rect x="2" y="5" width="20" height="14" rx="4" ry="4" fill="#FF0000"/>
              <polygon points="10,9 15.5,12 10,15" fill="#FFFFFF"/>
            </svg>
          </div>
          <div class="qa-body">
            <div class="qa-title">{{ card.title }}</div>
            <div class="qa-desc">{{ card.desc }}</div>
          </div>
          <div class="qa-arrow">›</div>
        </div>
      </section>

      <!-- Tabs + actions -->
      <section class="tabs-row">
        <div class="tabs">
          <button class="tab" :class="{active: activeTab==='my'}" @click="activeTab='my'">My Notes</button>
          <button class="tab" :class="{active: activeTab==='shared'}" @click="activeTab='shared'">Shared with Me</button>
        </div>
        <button class="btn outline" @click="createFolder">
          <span class="btn-ico">📁</span> Create Folder
        </button>
      </section>

      <!-- Notes list -->
      <section class="notes">
        <article class="note" v-for="note in filteredNotes" :key="note.id" @click="openNote(note)">
          <div class="note-ico" :class="note.type">
            <span v-if="note.type==='audio'">🎤</span>
            <span v-else>📄</span>
          </div>
          <div class="note-body">
            <div class="note-title">{{ note.title }}</div>
            <div class="note-meta">Last opened {{ note.lastOpened }}</div>
          </div>
          <button class="note-more" @click.stop="moreNote(note)">⋯</button>
        </article>
        <div v-if="filteredNotes.length===0" class="empty">No notes yet</div>
      </section>
    </main>
  </div>
</template>

<script>
export default {
  name: 'Dashboard',
  data() {
    return {
      activeTab: 'my',
      quickActions: [
        { key: 'blank', title: 'Blank document', desc: 'Start from scratch', color: 'c-purple', action: () => this.createBlank() },
        { key: 'audio', title: 'Record or upload audio', desc: 'Upload an audio file', color: 'c-violet', action: () => this.uploadAudio() },
        { key: 'doc',   title: 'Document upload', desc: 'Any PDF, DOC, PPT, etc', color: 'c-purple', action: () => this.uploadDoc() },
        { key: 'yt',    title: 'YouTube video', desc: 'Paste a YouTube link', color: 'c-red', action: () => this.addYoutube() },
      ],
      notes: [
        { id: '1', title: 'Basic Russian: Counting, Greeting, Stop', type: 'audio', lastOpened: 'about 6 hours ago', tab: 'my' },
        { id: '2', title: 'Untitled Document', type: 'doc', lastOpened: 'about 6 hours ago', tab: 'my' },
      ],
    }
  },
  computed: {
    filteredNotes() {
      return this.notes.filter(n => n.tab === this.activeTab)
    }
  },
  mounted() {
    // временно отключено требование токена
  },
  methods: {
    // Quick action stubs
    createBlank() { this.toast('Blank document') },
    uploadAudio() { this.toast('Audio upload') },
    uploadDoc()   { this.toast('Document upload') },
    addYoutube()  { this.toast('YouTube link') },

    createFolder() { this.toast('Create folder') },
    openNote(n) { this.toast('Open: ' + n.title) },
    moreNote() { this.toast('More…') },

    logout() {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      this.$router.push('/')
    },
    toast(msg) { console.log('[Dashboard]', msg) },
  }
}
</script>

<style scoped>
:root { --bg:#0f0f12; --panel:#171722; --line:rgba(255,255,255,.12); --muted:#b0b0b0; --text:#fff; }
.dash-wrap { display:flex; min-height:100vh; background:rgba(10,10,10,.9); color:var(--text); }

/* Sidebar */
.sidebar { width:56px; padding:12px 8px; border-right:1px solid var(--line); background:rgba(10,10,10,.8); backdrop-filter: blur(10px); position:sticky; top:0; height:100vh; display:flex; flex-direction:column; gap:10px; }
.side-btn { height:42px; border-radius:12px; border:1px solid var(--line); background:rgba(255,255,255,.04); color:var(--text); display:flex; align-items:center; justify-content:center; cursor:pointer; }
.side-btn.active, .side-btn:hover { box-shadow:0 0 0 2px rgba(0,212,255,.15) inset; }
.ico { font-size:18px; }

/* Main */
.main { flex:1; padding:28px 28px 60px; }
.page-header { margin:8px 0 18px; }
.h1 { font-size:40px; font-weight:800; margin:0 0 8px; }
.sub { color:var(--muted); margin:0; }

/* Quick actions */
.quick { display:grid; grid-template-columns: repeat(4, minmax(200px, 1fr)); gap:16px; margin:18px 0 22px; }
.qa { display:grid; grid-auto-flow:column; grid-template-columns:48px 1fr 16px; align-items:center; gap:14px; padding:16px; border-radius:14px; border:1px solid var(--line); background:rgba(255,255,255,.04); cursor:pointer; transition:transform .15s ease, box-shadow .15s ease; }
.qa:hover { transform: translateY(-1px); box-shadow: 0 6px 18px rgba(0,0,0,.25), 0 0 0 1px rgba(0,212,255,.15) inset; }
.qa-ico { width:48px; height:48px; border-radius:12px; display:flex; align-items:center; justify-content:center; font-size:20px; box-shadow:0 0 0 1px rgba(255,255,255,.08) inset; color:#fff; }
.c-purple { background: linear-gradient(135deg,#7c3aed33,#7c3aed22); }
.c-violet { background: linear-gradient(135deg,#6d28d933,#6d28d922); }
.c-blue   { background: linear-gradient(135deg,#00d4ff33,#00d4ff22); }
.c-red    { background: linear-gradient(135deg,#ef444433,#ef444422); }

/* SVGs inside icons */
.qa-ico svg { width:22px; height:22px; display:block; }
.qa-ico.doc svg { width:28px; height:28px; }
.qa-ico path { fill:#fff; }
.qa-ico rect { fill:#fff; }
.qa-ico.yt rect { fill:#FF0000; }
.qa-ico.yt polygon { fill:#FFFFFF; }
.qa-ico.doc .doc-paper { fill:#5b21b6; }
.qa-ico.doc .doc-fold { fill:#8b5cf6; }
.qa-ico.doc .pill { fill: rgba(255,255,255,.2); }
.qa-ico .doc-text { fill:#fff; font-weight:900; font-size:8.8px; font-family: ui-sans-serif, -apple-system, Segoe UI, Roboto, Helvetica, Arial, "Apple Color Emoji", "Segoe UI Emoji"; letter-spacing:.3px; dominant-baseline: middle; }
.qa-ico.yt .yt-svg { filter: drop-shadow(0 0 0 rgba(0,0,0,0)); }

.qa-title { font-weight:700; }
.qa-desc  { color:var(--muted); font-size:12px; }
.qa-arrow { color:var(--muted); font-size:22px; justify-self:end; }

/* Tabs */
.tabs-row { display:flex; align-items:center; justify-content:space-between; gap:12px; margin:6px 0 14px; }
.tabs { display:flex; gap:8px; }
.tab { height:36px; padding:0 14px; border-radius:10px; border:1px solid var(--line); background:rgba(255,255,255,.04); color:var(--text); cursor:pointer; }
.tab.active, .tab:hover { box-shadow:0 0 0 2px rgba(0,212,255,.15) inset; }
.btn.outline { height:36px; padding:0 12px; border-radius:10px; border:1px solid var(--line); background:rgba(255,255,255,.04); color:var(--text); display:inline-flex; align-items:center; gap:8px; cursor:pointer; }
.btn-ico { font-size:16px; }

/* Notes */
.notes { display:flex; flex-direction:column; gap:12px; margin-top:8px; }
.note { display:grid; grid-template-columns:44px 1fr auto; align-items:center; gap:14px; padding:14px; border-radius:12px; border:1px solid var(--line); background:rgba(255,255,255,.035); }
.note:hover { box-shadow:0 0 0 1px rgba(0,212,255,.15) inset; }
.note-ico { width:44px; height:44px; border-radius:12px; display:flex; align-items:center; justify-content:center; font-size:18px; background:rgba(255,255,255,.05); }
.note-ico.audio { background: linear-gradient(135deg,#7c3aed33,#00d4ff22); }
.note-title { font-weight:700; }
.note-meta { color:var(--muted); font-size:12px; }
.note-more { height:32px; width:32px; border-radius:8px; border:1px solid var(--line); background:rgba(255,255,255,.04); color:var(--text); cursor:pointer; }
.empty { color:var(--muted); text-align:center; padding:40px 0; }

@media (max-width: 980px){ .quick{ grid-template-columns: 1fr 1fr; } }
@media (max-width: 640px){ .quick{ grid-template-columns: 1fr; } .sidebar{display:none;} .main{padding:20px;} .h1{font-size:32px;} }
</style> 