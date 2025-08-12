<template>
  <div class="dash-wrap" :class="{ 'drawer-open': sidebarOpen }">
    <!-- Top bar -->
    <header class="topbar">
      <div class="brand">
        <span class="brand-mark">✒️</span>
        <span class="brand-text">turbo ai</span>
      </div>
      <button v-if="!sidebarOpen" class="burger" @click="toggleSidebar" aria-label="Toggle menu">
        <span class="burger-box"><span class="burger-lines"></span></span>
      </button>
    </header>

    <!-- Single expanding side pane -->
    <aside class="navpane" :class="{ open: sidebarOpen }">
      <div v-if="sidebarOpen" class="pane-head">
        <div class="pane-logo">SpeakApperAI</div>
        <button class="collapse-btn" @click="toggleSidebar" aria-label="Close menu">
          <svg viewBox="0 0 24 24" width="18" height="18"><path d="M15 6 9 12l6 6" fill="none" stroke="#ffffff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
        </button>
      </div>
      <button v-else class="rail-burger" @click="toggleSidebar" aria-label="Open menu">
        <svg viewBox="0 0 24 24" width="22" height="22" aria-hidden="true">
          <path d="M4 7h16M4 12h16M4 17h16" stroke="#ffffff" stroke-width="2" stroke-linecap="round"/>
        </svg>
      </button>

      <nav class="menu">
        <button class="menu-item active" @click="$router.push('/dashboard')">
          <span class="mi-ico">
            <!-- Home icon (outline for collapsed look) -->
            <svg viewBox="0 0 24 24" aria-hidden="true" class="ico-outline">
              <path d="M3 10l9-7 9 7v8a2 2 0 0 1-2 2h-4v-5H9v5H5a2 2 0 0 1-2-2v-8z" fill="none" stroke="#ffffff" stroke-width="2" stroke-linejoin="round"/>
            </svg>
          </span>
          <span class="mi-text">Dashboard</span>
        </button>

        <button class="menu-item" @click="goSettings">
          <span class="mi-ico">
            <!-- Settings icon (outline for collapsed look) -->
            <svg viewBox="0 0 24 24" aria-hidden="true" class="ico-outline">
              <path d="M12 15.5a3.5 3.5 0 1 0 0-7 3.5 3.5 0 0 0 0 7z" fill="none" stroke="#ffffff" stroke-width="2"/>
              <path d="M19.4 15a8 8 0 0 0 .06-6l2.04-1.59-1.92-3.32-2.39.96a7.8 7.8 0 0 0-2.27-1.3L14.6 1h-5.2L8.98 3.75a7.8 7.8 0 0 0-2.27 1.3l-2.39-.96L2.4 7.41 4.44 9a8 8 0 0 0 0 6l-2.04 1.59 1.92 3.32 2.39-.96c.7.53 1.45.95 2.27 1.3L9.4 23h5.2l.22-1.75c.82-.35 1.57-.77 2.27-1.3l2.39.96 1.92-3.32L19.4 15z" fill="none" stroke="#ffffff" stroke-width="2" stroke-linejoin="round" stroke-linecap="round"/>
            </svg>
          </span>
          <span class="mi-text">Settings</span>
        </button>

        <transition name="upg">
          <button v-show="sidebarOpen" class="upgrade" @click="upgrade">
            <span class="spark">✨</span>
            <span class="upg-text">Upgrade to Premium</span>
          </button>
        </transition>
      </nav>
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
      sidebarOpen: false,
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
    toggleSidebar(){ this.sidebarOpen = !this.sidebarOpen },
    goSettings(){ this.toast('Settings') },
    upgrade(){ this.toast('Upgrade') },

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
:root { --bg:#0f0f12; --panel:#171722; --line:rgba(255,255,255,.12); --muted:#b0b0b0; --text:#fff; --accent:#7c3aed; }
.dash-wrap { display:grid; grid-template-columns: 72px 1fr; min-height:100vh; background:rgba(10,10,10,.9); color:var(--text); transition: grid-template-columns .35s cubic-bezier(.22,.61,.36,1); }
.dash-wrap.drawer-open { grid-template-columns: 300px 1fr; }

/* Topbar */
.topbar { position:sticky; top:0; z-index:30; grid-column: 1 / -1; display:flex; align-items:center; justify-content:space-between; padding:14px 16px; background:rgba(10,10,10,.85); backdrop-filter: blur(10px); border-bottom:1px solid var(--line); }
.brand { display:flex; align-items:center; gap:10px; font-weight:800; font-size:22px; }
.brand-text { letter-spacing:.5px; }
.burger { height:36px; width:44px; border-radius:10px; border:1px solid var(--line); background:rgba(255,255,255,.04); color:var(--text); display:flex; align-items:center; justify-content:center; }
.burger-box { position:relative; width:18px; height:14px; }
.burger-lines, .burger-lines:before, .burger-lines:after { content:""; position:absolute; left:0; right:0; height:2px; background:#fff; border-radius:2px; }
.burger-lines{ top:6px; }
.burger-lines:before{ top:-6px; }
.burger-lines:after{ top:6px; transform: translateY(6px); }

/* Expanding side pane */
.navpane { position:sticky; top:56px; height:calc(100vh - 56px); padding:12px 8px; border-right:1px solid var(--line); background:#0f0f12; width:72px; transition: width .35s cubic-bezier(.22,.61,.36,1); display:flex; flex-direction:column; gap:18px; align-items:center; }
.navpane.open { width:300px; align-items:stretch; }
.pane-head { display:flex; align-items:center; justify-content:space-between; padding:12px 16px; width:100%; }
.pane-logo { font-weight:800; font-size:22px; letter-spacing:.5px; }
.collapse-btn { height:36px; width:36px; border-radius:10px; border:1px solid var(--line); background:rgba(255,255,255,.04); color:var(--text); display:flex; align-items:center; justify-content:center; }
.collapse-btn svg { width:18px; height:18px; }
.rail-burger { height:54px; width:54px; border-radius:14px; background:rgba(255,255,255,.02); border:2px solid #1f3e83; box-shadow: inset 0 0 0 2px rgba(91,131,255,.25); display:grid; place-items:center; color:#fff; align-self:center; }
.menu { display:flex; flex-direction:column; gap:16px; width:100%; }
.menu-item { display:grid; grid-template-columns:28px 1fr; align-items:center; column-gap:12px; height:56px; padding:0 14px; border-radius:14px; border:1px solid var(--line); background:#1b1b23; color:#fff; text-align:left; width:56px; margin:0 auto; transition: width .35s cubic-bezier(.22,.61,.36,1), background .2s ease; }
.menu-item.active { box-shadow:0 0 0 2px rgba(124,58,237,.25) inset; }
.navpane.open .menu-item { width:100%; }
.mi-ico { width:28px; height:28px; display:grid; place-items:center; }
.mi-ico svg { width:22px; height:22px; fill:#fff; }
.mi-text { white-space:nowrap; overflow:hidden; opacity:0; width:0; will-change: opacity, width; transition: opacity .25s ease .1s, width .35s cubic-bezier(.22,.61,.36,1); font-size:18px; font-weight:700; }
.navpane.open .mi-text { opacity:1; width:auto; }
.upgrade { display:flex; align-items:center; justify-content:center; gap:10px; height:56px; border-radius:16px; border:1px solid rgba(124,58,237,.35); background:#6227e9; color:#fff; font-weight:800; width:56px; margin:8px auto 0; transition: width .35s cubic-bezier(.22,.61,.36,1); }
.navpane.open .upgrade { width:100%; }
.upg-text { white-space:nowrap; overflow:hidden; opacity:0; width:0; transition: opacity .25s ease .1s, width .35s cubic-bezier(.22,.61,.36,1); }
.navpane.open .upg-text { opacity:1; width:auto; }
.spark { font-size:18px; }

/* Upgrade appear from bottom */
.upg-enter-from, .upg-leave-to { opacity:0; transform: translateY(12px); }
.upg-enter-active, .upg-leave-active { transition: all .35s cubic-bezier(.22,.61,.36,1); }
.upg-enter-to, .upg-leave-from { opacity:1; transform: translateY(0); }

/* Main */
.main { padding:20px 16px 60px; }
@media (min-width: 980px){
  .topbar{ display:none; }
  .navpane{ top:0; height:100vh; }
}
@media (max-width: 979px){
  .main{ grid-column: 1 / -1; }
}

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
@media (max-width: 640px){ .quick{ grid-template-columns: 1fr; } .h1{font-size:32px;} }
</style> 