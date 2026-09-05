<template>
  <UApp>
    <div class="h-screen w-screen flex overflow-hidden bg-[#090a0f] text-[#e2e8f0] font-sans antialiased select-none">
    
    <!-- Sidebar -->
    <aside class="w-80 border-r border-white/8 bg-[#0f1118] flex flex-col shrink-0 z-20 shadow-xl">
      <!-- Sidebar Header (Full Row Stacking, No Squashed Words) -->
      <div class="p-3.5 border-b border-white/8 bg-[#141721]/60 backdrop-blur-md flex flex-col gap-3">
        <!-- Row 1: Brand & Status -->
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2.5">
            <div class="w-8 h-8 rounded-xl flame-gradient flex items-center justify-center shadow-lg shadow-amber-500/20 text-white font-bold text-base shrink-0">
              🔥
            </div>
            <div>
              <div class="flex items-center gap-1.5">
                <h1 class="text-sm font-bold text-white tracking-wide">Mikhaz Fire</h1>
                <span class="text-[9px] px-1.5 py-0.5 rounded-full bg-amber-500/15 text-amber-300 font-mono border border-amber-500/30 font-semibold">v1.1</span>
              </div>
              <p class="text-[10px] text-zinc-400 font-medium">Auth &amp; Firestore Controller</p>
            </div>
          </div>

          <!-- Status Indicator Dot -->
          <div class="flex items-center gap-1.5 px-2 py-0.5 rounded-full bg-white/4 border border-white/6" title="Engine Ready">
            <span class="relative flex h-1.5 w-1.5">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-1.5 w-1.5 bg-emerald-500"></span>
            </span>
            <span class="text-[10px] font-mono text-zinc-400">Ready</span>
          </div>
        </div>

        <!-- Row 2: Action Buttons (Proper Full Row Stacking, No Squashed Words) -->
        <div class="flex items-center gap-2">
          <UButton
            icon="i-heroicons-user-plus"
            size="xs"
            color="warning"
            variant="subtle"
            :loading="loggingInGoogle"
            @click="triggerGoogleLogin"
            title="Sign in with Google Account in Browser"
            class="flex-1 justify-center font-semibold text-xs py-1.5 whitespace-nowrap"
          >
            Google Sign In
          </UButton>
          <UButton
            icon="i-heroicons-key"
            size="xs"
            color="neutral"
            variant="subtle"
            @click="showAddSaModal = true"
            title="Add Service Account Key JSON"
            class="font-medium text-xs py-1.5 px-2.5 whitespace-nowrap bg-white/4 hover:bg-white/8 text-zinc-300 border border-white/10"
          >
            Add Key
          </UButton>
        </div>
      </div>

      <!-- Search & Filter Projects Bar -->
      <div class="p-2.5 border-b border-white/8 bg-[#0b0d13] space-y-2">
        <div class="relative">
          <UInput
            v-model="sidebarProjectFilter"
            placeholder="Search projects..."
            icon="i-heroicons-magnifying-glass"
            size="xs"
            class="w-full bg-[#141721] text-xs border-white/10"
            :ui="{ icon: { trailing: { pointer: '' } } }"
          >
            <template #trailing v-if="sidebarProjectFilter">
              <UButton
                color="neutral"
                variant="link"
                size="2xs"
                icon="i-heroicons-x-mark"
                class="text-zinc-400 hover:text-white"
                @click="sidebarProjectFilter = ''"
              />
            </template>
          </UInput>
        </div>

        <!-- Accounts & Rescan Bar -->
        <div class="flex items-center justify-between text-[11px] text-zinc-400 px-1">
          <span class="flex items-center gap-1">
            <span class="w-1.5 h-1.5 rounded-full bg-amber-400"></span>
            <span>Accounts: {{ accounts.length }}</span>
          </span>
          <UButton
            size="2xs"
            variant="ghost"
            color="warning"
            class="text-[11px] px-1 py-0.5 text-amber-400 hover:text-amber-300 font-medium"
            @click="refreshAccounts"
            :loading="scanning"
          >
            ↻ Rescan
          </UButton>
        </div>
      </div>

      <!-- Accounts / Projects List -->
      <div class="flex-1 overflow-y-auto p-2 space-y-2">
        <!-- Pinned / Favorite Projects Section -->
        <div v-if="pinnedProjects.length > 0 && !sidebarProjectFilter" class="space-y-1 mb-3">
          <div class="px-2 py-1 text-[10px] font-semibold uppercase tracking-wider text-amber-400/80 flex items-center justify-between">
            <span class="flex items-center gap-1">⭐ Pinned Projects</span>
            <span class="text-[9px] text-zinc-500 font-mono">{{ pinnedProjects.length }}</span>
          </div>
          <div 
            v-for="pin in pinnedProjects" 
            :key="pin.projectId"
            class="group/pin p-2 rounded-lg bg-[#141721] border border-amber-500/20 hover:border-amber-500/40 text-xs transition-all shadow-sm"
          >
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-1.5 truncate">
                <span class="text-amber-400 text-xs">🔥</span>
                <span class="font-semibold text-white truncate" :title="pin.displayName">{{ pin.displayName || pin.projectId }}</span>
              </div>
              <button 
                @click.stop="togglePinProject(pin)" 
                class="text-amber-400 hover:text-amber-300 text-xs p-0.5" 
                title="Unpin project"
              >
                ★
              </button>
            </div>
            <div class="text-[10px] text-zinc-400 font-mono truncate mt-0.5">{{ pin.projectId }}</div>
            <div class="flex gap-1.5 mt-2">
              <UButton
                size="2xs"
                color="warning"
                variant="subtle"
                class="text-[10px] py-0.5 flex-1 justify-center font-medium"
                icon="i-heroicons-user-group"
                @click.stop="openTab('auth', pin.email, pin.projectId)"
              >
                Auth
              </UButton>
              <UButton
                size="2xs"
                color="warning"
                variant="subtle"
                class="text-[10px] py-0.5 flex-1 justify-center font-medium"
                icon="i-heroicons-circle-stack"
                @click.stop="openTab('firestore', pin.email, pin.projectId)"
              >
                Firestore
              </UButton>
            </div>
          </div>
        </div>

        <div v-if="accounts.length === 0" class="text-xs text-zinc-400 p-6 text-center space-y-3 bg-[#141721]/50 rounded-xl border border-white/5 m-2">
          <div class="w-10 h-10 rounded-full bg-amber-500/10 text-amber-400 mx-auto flex items-center justify-center text-lg">
            🔑
          </div>
          <p class="font-medium text-white text-sm">No Accounts Connected</p>
          <p class="text-[11px] text-zinc-400 leading-relaxed">Connect a Google account or upload a Service Account JSON to inspect your Firebase databases.</p>
          <UButton size="xs" color="warning" icon="i-heroicons-user-plus" class="mx-auto font-medium" @click="triggerGoogleLogin" :loading="loggingInGoogle">
            Sign In with Google
          </UButton>
        </div>
        
        <div v-for="acc in filteredAccounts" :key="acc.email" class="space-y-1">
          <!-- Account Item Header -->
          <div 
            class="flex items-center justify-between p-2 rounded-lg cursor-pointer text-xs transition-all border"
            :class="expandedAccount === acc.email 
              ? 'bg-[#181b26] border-amber-500/30 text-white shadow-sm' 
              : 'hover:bg-[#141721] border-transparent text-zinc-300 hover:text-white'"
            @click="toggleAccount(acc.email)"
          >
            <div class="flex items-center gap-2 truncate">
              <span class="text-xs transform transition-transform duration-150 text-zinc-400" :class="expandedAccount === acc.email ? 'rotate-90 text-amber-400' : ''">
                ▸
              </span>
              <span class="truncate font-semibold tracking-tight" :title="acc.email">{{ acc.email }}</span>
            </div>
            <UBadge size="xs" color="warning" variant="subtle" class="text-[9px] font-mono px-1.5 py-0.5 uppercase tracking-wider shrink-0">
              {{ acc.is_service_account ? 'SA' : acc.source }}
            </UBadge>
          </div>

          <!-- Projects List for Expanded Account -->
          <div v-if="expandedAccount === acc.email" class="pl-3 pr-1 py-1 space-y-1 border-l-2 border-amber-500/30 ml-2">
            <div v-if="loadingProjects" class="text-xs text-zinc-400 py-3 pl-2 flex items-center gap-2">
              <UIcon name="i-heroicons-arrow-path" class="w-3.5 h-3.5 animate-spin text-amber-400" />
              <span>Loading projects...</span>
            </div>
            <div v-else-if="filteredProjects.length === 0" class="text-xs text-zinc-400 py-3 pl-2 italic">
              {{ sidebarProjectFilter ? 'No matching projects.' : 'No projects found.' }}
            </div>
            
            <div 
              v-for="proj in filteredProjects" 
              :key="proj.projectId" 
              class="group/proj p-2 rounded-lg hover:bg-[#181b26] text-xs text-zinc-300 flex flex-col gap-1 cursor-pointer transition-all border border-transparent hover:border-white/10 hover:shadow-md"
            >
              <div class="flex items-center justify-between">
                <span class="font-semibold text-white truncate text-xs" :title="proj.displayName">
                  {{ proj.displayName || proj.projectId }}
                </span>
                <button
                  @click.stop="togglePinProject({ ...proj, email: acc.email })"
                  class="text-zinc-600 group-hover/proj:text-zinc-400 hover:text-amber-400! text-xs transition-colors p-0.5"
                  :class="isPinned(proj.projectId) ? 'text-amber-400!' : ''"
                  :title="isPinned(proj.projectId) ? 'Unpin project' : 'Pin project to top'"
                >
                  ★
                </button>
              </div>
              <span class="text-[10px] text-zinc-400 font-mono truncate">{{ proj.projectId }}</span>
              
              <!-- Quick Action Options on Hover -->
              <div class="flex gap-1.5 mt-1 opacity-0 group-hover/proj:opacity-100 transition-opacity">
                <UButton
                  size="2xs"
                  color="warning"
                  variant="subtle"
                  class="text-[10px] py-0.5 flex-1 justify-center font-medium"
                  icon="i-heroicons-user-group"
                  @click.stop="openTab('auth', acc.email, proj.projectId)"
                >
                  Auth
                </UButton>
                <UButton
                  size="2xs"
                  color="warning"
                  variant="subtle"
                  class="text-[10px] py-0.5 flex-1 justify-center font-medium"
                  icon="i-heroicons-circle-stack"
                  @click.stop="openTab('firestore', acc.email, proj.projectId)"
                >
                  Firestore
                </UButton>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Sidebar Footer -->
      <div class="p-3 bg-[#0b0d13] border-t border-white/8 flex items-center justify-between text-xs text-zinc-400">
        <div class="flex items-center gap-2">
          <span class="relative flex h-2 w-2">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
          </span>
          <span class="text-[11px] font-medium text-zinc-300">Local Engine Ready</span>
        </div>
        <button 
          class="text-[10px] font-mono text-zinc-400 hover:text-amber-400 transition-colors"
          @click="reopenWelcomeTab"
          title="Open Welcome Dashboard"
        >
          Overview ↗
        </button>
      </div>
    </aside>

    <!-- Main Content Area -->
    <main class="flex-1 flex flex-col overflow-hidden bg-[#090a0f]">
      <!-- Tabs Bar -->
      <div class="h-11 bg-[#0f1118] border-b border-white/8 flex items-center overflow-x-auto select-none scrollbar-none px-2 gap-1.5 z-10">
        <!-- Welcome Tab -->
        <div 
          v-if="showWelcomeTab"
          class="px-3 py-1.5 h-8 flex items-center gap-2 text-xs rounded-lg cursor-pointer transition-all border"
          :class="activeTabId === 'welcome' 
            ? 'bg-[#181b26] text-white font-semibold border-amber-500/30 shadow-sm' 
            : 'hover:bg-white/4 text-zinc-400 border-transparent'"
          @click="activeTabId = 'welcome'"
        >
          <span class="text-amber-400 text-xs">🏠</span>
          <span>Welcome</span>
          <span 
            class="text-zinc-500 hover:text-white rounded p-0.5 leading-none transition-colors text-xs"
            @click.stop="closeWelcomeTab"
          >
            ✕
          </span>
        </div>

        <!-- Dynamic Tabs -->
        <div 
          v-for="tab in openTabs" 
          :key="tab.id"
          class="px-3 py-1.5 h-8 flex items-center gap-2 text-xs rounded-lg cursor-pointer transition-all border max-w-60 shrink-0"
          :class="activeTabId === tab.id 
            ? 'bg-[#181b26] text-white font-semibold border-amber-500/40 shadow-sm' 
            : 'hover:bg-white/4 text-zinc-400 border-transparent'"
          @click="activeTabId = tab.id"
        >
          <span class="text-xs shrink-0">
            {{ tab.type === 'auth' ? '👥' : (tab.activeDocPath ? '📄' : (tab.activeCollection ? '📁' : '🗄️')) }}
          </span>
          <span class="truncate font-medium text-amber-400" :title="getTabFullTitle(tab)">
            {{ getTabTitle(tab) }}
          </span>
          <span 
            class="text-zinc-500 hover:text-white rounded p-0.5 leading-none transition-colors text-xs ml-auto shrink-0"
            @click.stop="closeTab(tab.id)"
          >
            ✕
          </span>
        </div>
      </div>

      <!-- Tab Contents Container -->
      <div class="flex-1 overflow-hidden relative">
        <!-- Welcome Screen -->
        <div v-if="activeTabId === 'welcome'" class="h-full overflow-y-auto p-8 flex flex-col justify-center items-center text-center max-w-4xl mx-auto space-y-6">
          <div class="relative">
            <div class="absolute -inset-4 bg-linear-to-r from-amber-500 to-orange-600 rounded-full blur-xl opacity-20 animate-pulse"></div>
            <img src="/logo.svg" class="w-20 h-20 relative drop-shadow-[0_0_25px_rgba(245,158,11,0.3)] object-contain" alt="Logo" />
          </div>

          <div>
            <h2 class="text-3xl font-extrabold text-white tracking-tight flame-gradient-text">
              Mikhaz Fire Dashboard
            </h2>
            <p class="text-zinc-400 mt-2 text-sm max-w-xl mx-auto leading-relaxed">
              Fast, offline-first Firebase Auth & Firestore administration. Inspect document structures, execute custom JavaScript queries, and migrate user accounts with embedded password hashes.
            </p>
          </div>
          
          <div class="grid grid-cols-2 gap-4 w-full text-left">
            <div class="p-5 glass-card rounded-xl transition-all">
              <div class="w-9 h-9 rounded-lg bg-amber-500/15 border border-amber-500/30 flex items-center justify-center text-amber-400 text-lg mb-3">
                ✨
              </div>
              <h3 class="font-bold text-white mb-1">Local Credential Scan</h3>
              <p class="text-xs text-zinc-400 leading-relaxed">
                Automatically reads from gcloud ADC and Firebase CLI profiles stored on your machine with zero remote key transmission.
              </p>
            </div>

            <div class="p-5 glass-card rounded-xl transition-all">
              <div class="w-9 h-9 rounded-lg bg-emerald-500/15 border border-emerald-500/30 flex items-center justify-center text-emerald-400 text-lg mb-3">
                👥
              </div>
              <h3 class="font-bold text-white mb-1">Auth Controller & Claims</h3>
              <p class="text-xs text-zinc-400 leading-relaxed">
                Search, update profiles, toggle disabled status, and modify custom user claims (JSON metadata) in real-time.
              </p>
            </div>

            <div class="p-5 glass-card rounded-xl transition-all">
              <div class="w-9 h-9 rounded-lg bg-sky-500/15 border border-sky-500/30 flex items-center justify-center text-sky-400 text-lg mb-3">
                🗄️
              </div>
              <h3 class="font-bold text-white mb-1">Firestore Explorer</h3>
              <p class="text-xs text-zinc-400 leading-relaxed">
                Navigate document graphs, tree views, structured tables, subcollections, and perform nested recursive deletions.
              </p>
            </div>

            <div class="p-5 glass-card rounded-xl transition-all">
              <div class="w-9 h-9 rounded-lg bg-purple-500/15 border border-purple-500/30 flex items-center justify-center text-purple-400 text-lg mb-3">
                ⚡
              </div>
              <h3 class="font-bold text-white mb-1">JS Query Engine & Sync</h3>
              <p class="text-xs text-zinc-400 leading-relaxed">
                Run JavaScript queries locally against Firestore using embedded Goja and copy collections directly between projects.
              </p>
            </div>
          </div>
        </div>

        <!-- Active Dynamic Tab View -->
        <div 
          v-for="tab in openTabs" 
          :key="tab.id" 
          v-show="activeTabId === tab.id"
          class="h-full w-full flex flex-col"
        >
          <!-- Auth Component -->
          <div v-if="tab.type === 'auth'" class="h-full flex flex-col overflow-hidden p-5 space-y-4">
            <!-- Header & Action Bar -->
            <div class="flex items-center justify-between bg-[#0f1118] p-4 rounded-xl border border-white/8 shadow-sm">
              <div class="space-y-0.5">
                <div class="flex items-center gap-2">
                  <span class="text-lg">👥</span>
                  <h2 class="text-lg font-bold text-white tracking-tight">Auth User Manager</h2>
                  <UBadge color="warning" variant="subtle" size="xs" class="font-mono">{{ tab.projectId }}</UBadge>
                </div>
                <p class="text-[11px] text-zinc-400 font-mono">
                  Account: <span class="text-zinc-300">{{ tab.email }}</span>
                </p>
              </div>

              <div class="flex items-center gap-2">
                <UButton
                  color="info"
                  variant="subtle"
                  size="xs"
                  icon="i-heroicons-arrow-up-tray"
                  class="font-medium"
                  @click="openExportAuthModal(tab)"
                >
                  Export Users
                </UButton>
                <UButton
                  color="success"
                  variant="subtle"
                  size="xs"
                  icon="i-heroicons-arrow-down-tray"
                  class="font-medium"
                  @click="openImportAuthModal(tab)"
                >
                  Import Users
                </UButton>
                <UButton
                  color="warning"
                  size="xs"
                  icon="i-heroicons-user-plus"
                  class="font-medium shadow-sm shadow-amber-500/20"
                  @click="openCreateUserModal(tab)"
                >
                  Create User
                </UButton>
              </div>
            </div>

            <!-- Metric Summary Stats Cards -->
            <div class="grid grid-cols-4 gap-3">
              <div class="p-3 bg-[#141721] border border-white/8 rounded-xl flex items-center justify-between">
                <div>
                  <div class="text-[10px] uppercase font-semibold text-zinc-400">Total Loaded</div>
                  <div class="text-xl font-extrabold text-white mt-0.5">{{ tab.users ? tab.users.length : 0 }}</div>
                </div>
                <div class="w-8 h-8 rounded-lg bg-amber-500/10 text-amber-400 flex items-center justify-center font-bold">👥</div>
              </div>

              <div class="p-3 bg-[#141721] border border-white/8 rounded-xl flex items-center justify-between">
                <div>
                  <div class="text-[10px] uppercase font-semibold text-zinc-400">Active Accounts</div>
                  <div class="text-xl font-extrabold text-emerald-400 mt-0.5">
                    {{ tab.users ? tab.users.filter(u => !u.disabled).length : 0 }}
                  </div>
                </div>
                <div class="w-8 h-8 rounded-lg bg-emerald-500/10 text-emerald-400 flex items-center justify-center font-bold">✓</div>
              </div>

              <div class="p-3 bg-[#141721] border border-white/8 rounded-xl flex items-center justify-between">
                <div>
                  <div class="text-[10px] uppercase font-semibold text-zinc-400">Disabled</div>
                  <div class="text-xl font-extrabold text-rose-400 mt-0.5">
                    {{ tab.users ? tab.users.filter(u => u.disabled).length : 0 }}
                  </div>
                </div>
                <div class="w-8 h-8 rounded-lg bg-rose-500/10 text-rose-400 flex items-center justify-center font-bold">🚫</div>
              </div>

              <div class="p-3 bg-[#141721] border border-white/8 rounded-xl flex items-center justify-between">
                <div>
                  <div class="text-[10px] uppercase font-semibold text-zinc-400">Custom Claims</div>
                  <div class="text-xl font-extrabold text-purple-400 mt-0.5">
                    {{ countUsersWithClaims(tab.users) }}
                  </div>
                </div>
                <div class="w-8 h-8 rounded-lg bg-purple-500/10 text-purple-400 flex items-center justify-center font-bold">⚡</div>
              </div>
            </div>

            <!-- Unified Search and Filter Bar -->
            <div class="flex items-center gap-3 bg-[#141721] p-2.5 rounded-xl border border-white/8">
              <div class="flex-1 relative">
                <UInput
                  v-model="tab.searchQuery"
                  placeholder="Instant filter by UID, email, display name..."
                  size="sm"
                  class="w-full bg-[#090a0f] border-white/10"
                  icon="i-heroicons-magnifying-glass"
                  @keydown.enter="fetchAuthUsers(tab)"
                />
              </div>
              <UButton color="neutral" variant="subtle" size="sm" @click="fetchAuthUsers(tab)" :loading="tab.loading">
                Search API
              </UButton>

              <!-- Quick Status Filter Pills -->
              <div class="flex items-center gap-1 bg-[#090a0f] p-1 rounded-lg border border-white/10">
                <button
                  v-for="filterOpt in ['all', 'active', 'disabled', 'claims']"
                  :key="filterOpt"
                  type="button"
                  class="px-2.5 py-1 text-[10px] font-semibold uppercase rounded-md transition-all"
                  :class="(tab.statusFilter || 'all') === filterOpt 
                    ? 'bg-amber-500 text-black font-bold shadow' 
                    : 'text-zinc-400 hover:text-white'"
                  @click="tab.statusFilter = filterOpt"
                >
                  {{ filterOpt }}
                </button>
              </div>
            </div>

            <!-- Users Table -->
            <div class="flex-1 border border-white/8 bg-[#0f1118] rounded-xl overflow-hidden flex flex-col shadow-md min-h-0">
              <div class="flex-1 overflow-auto">
                <table class="w-full text-left border-collapse text-xs">
                  <thead>
                    <tr class="bg-[#141721] text-zinc-300 border-b border-white/8 font-semibold select-none sticky top-0 z-10 backdrop-blur">
                      <th class="p-3 w-12 text-center">#</th>
                      <th class="p-3 cursor-pointer hover:text-amber-400 transition-colors" @click="toggleAuthSort(tab, 'uid')">
                        UID
                        <span v-if="tab.authSortField === 'uid'" class="ml-1 text-[10px] text-amber-400">
                          {{ tab.authSortDir === 'asc' ? '▲' : '▼' }}
                        </span>
                      </th>
                      <th class="p-3 cursor-pointer hover:text-amber-400 transition-colors" @click="toggleAuthSort(tab, 'email')">
                        User / Email
                        <span v-if="tab.authSortField === 'email'" class="ml-1 text-[10px] text-amber-400">
                          {{ tab.authSortDir === 'asc' ? '▲' : '▼' }}
                        </span>
                      </th>
                      <th class="p-3 cursor-pointer hover:text-amber-400 transition-colors" @click="toggleAuthSort(tab, 'displayName')">
                        Display Name
                        <span v-if="tab.authSortField === 'displayName'" class="ml-1 text-[10px] text-amber-400">
                          {{ tab.authSortDir === 'asc' ? '▲' : '▼' }}
                        </span>
                      </th>
                      <th class="p-3">Status</th>
                      <th class="p-3">Claims</th>
                      <th class="p-3 cursor-pointer hover:text-amber-400 transition-colors" @click="toggleAuthSort(tab, 'createdAt')">
                        Created
                        <span v-if="tab.authSortField === 'createdAt'" class="ml-1 text-[10px] text-amber-400">
                          {{ tab.authSortDir === 'asc' ? '▲' : '▼' }}
                        </span>
                      </th>
                      <th class="p-3 text-right">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-if="tab.loading" class="text-center">
                      <td colspan="8" class="p-12 text-zinc-400">
                        <UIcon name="i-heroicons-arrow-path" class="w-6 h-6 animate-spin mx-auto text-amber-400 mb-2" />
                        <span>Loading users from project...</span>
                      </td>
                    </tr>
                    <tr v-else-if="!tab.users || tab.users.length === 0" class="text-center">
                      <td colspan="8" class="p-12 text-zinc-500">No users found in this project.</td>
                    </tr>
                    <tr v-else-if="getFilteredSortedUsers(tab).length === 0" class="text-center">
                      <td colspan="8" class="p-12 text-zinc-500">No users match the search filter.</td>
                    </tr>
                    <tr 
                      v-else 
                      v-for="(user, idx) in getFilteredSortedUsers(tab)" 
                      :key="user.uid" 
                      class="border-b border-white/4 hover:bg-white/3 transition-colors group/row"
                    >
                      <!-- Index & Avatar -->
                      <td class="p-3 text-center text-zinc-500 text-[11px] font-mono">
                        {{ idx + 1 }}
                      </td>

                      <!-- UID with 1-click Copy -->
                      <td class="p-3 font-mono text-[11px] text-amber-400">
                        <div class="flex items-center gap-1.5">
                          <span class="truncate max-w-32.5" :title="user.uid">{{ user.uid }}</span>
                          <button
                            type="button"
                            class="opacity-0 group-hover/row:opacity-100 text-zinc-500 hover:text-amber-400 p-0.5 rounded transition-opacity"
                            @click="copyText(user.uid, 'UID copied to clipboard!')"
                            title="Copy UID"
                          >
                            📋
                          </button>
                        </div>
                      </td>

                      <!-- Email with Avatar -->
                      <td class="p-3">
                        <div class="flex items-center gap-2 min-w-0">
                          <div class="w-6 h-6 rounded-full bg-linear-to-tr from-amber-500 to-orange-500 flex items-center justify-center text-white text-[10px] font-bold shrink-0">
                            {{ getUserInitial(user) }}
                          </div>
                          <div class="truncate max-w-45">
                            <span class="font-medium text-white truncate" :title="user.email">{{ user.email || 'Anonymous' }}</span>
                          </div>
                          <button
                            v-if="user.email"
                            type="button"
                            class="opacity-0 group-hover/row:opacity-100 text-zinc-500 hover:text-amber-400 p-0.5 rounded transition-opacity shrink-0"
                            @click="copyText(user.email, 'Email copied to clipboard!')"
                            title="Copy Email"
                          >
                            📋
                          </button>
                        </div>
                      </td>

                      <!-- Display Name -->
                      <td class="p-3 text-zinc-300 truncate max-w-35" :title="user.displayName">
                        {{ user.displayName || '-' }}
                      </td>

                      <!-- Status Badge -->
                      <td class="p-3">
                        <UBadge 
                          :color="user.disabled ? 'danger' : 'success'" 
                          variant="subtle" 
                          size="xs"
                          class="font-medium"
                        >
                          {{ user.disabled ? 'Disabled' : 'Active' }}
                        </UBadge>
                      </td>

                      <!-- Claims Badge -->
                      <td class="p-3">
                        <span 
                          v-if="hasCustomClaims(user)"
                          class="text-[10px] font-mono px-1.5 py-0.5 rounded bg-purple-500/15 text-purple-300 border border-purple-500/30 cursor-pointer hover:border-purple-400"
                          @click="editUserClaims(tab, user)"
                          title="Click to view/edit claims"
                        >
                          ⚡ Has Claims
                        </span>
                        <span v-else class="text-zinc-600 text-[10px]">-</span>
                      </td>

                      <!-- Created Date -->
                      <td class="p-3 text-zinc-400 text-[11px] font-mono">{{ formatDate(user.createdAt) }}</td>

                      <!-- Actions -->
                      <td class="p-3 text-right space-x-1 whitespace-nowrap">
                        <UButton size="xs" variant="ghost" color="neutral" @click="editUserClaims(tab, user)" title="Edit Claims">
                          Claims
                        </UButton>
                        <UButton size="xs" variant="ghost" color="neutral" @click="editUserBasic(tab, user)" title="Edit User">
                          Edit
                        </UButton>
                        <UButton size="xs" variant="ghost" color="danger" @click="deleteUserConfirm(tab, user)" title="Delete User">
                          Delete
                        </UButton>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <!-- Pagination / Next Page -->
              <div v-if="tab.nextPageToken" class="p-3 bg-[#141721] border-t border-white/8 flex justify-between items-center text-xs">
                <span class="text-zinc-400 font-mono text-[11px]">More users available on server</span>
                <UButton size="xs" color="warning" @click="fetchAuthUsersNext(tab)">
                  Load More Users
                </UButton>
              </div>
            </div>
          </div>

          <!-- Firestore Component -->
          <div v-else-if="tab.type === 'firestore'" class="h-full flex overflow-hidden">
            <!-- Left Pane: Collections Explorer with Search Filter -->
            <!-- Left Pane: Collections Explorer with Search Filter -->
            <div 
              class="border-r border-white/8 bg-[#0f1118] flex flex-col shrink-0 transition-all duration-200"
              :class="tab.collectionsCollapsed ? 'w-11' : 'w-64'"
            >
              <!-- Collections Header -->
              <div class="p-2.5 border-b border-white/8 flex items-center justify-between bg-[#141721]/60">
                <div v-if="!tab.collectionsCollapsed" class="flex items-center gap-1.5 min-w-0">
                  <span class="text-amber-400 text-sm shrink-0">🗄️</span>
                  <h3 class="font-bold text-xs uppercase tracking-wider text-white truncate">Collections</h3>
                  <span class="text-[10px] text-zinc-400 font-mono shrink-0">({{ tab.collections ? tab.collections.length : 0 }})</span>
                </div>
                <div class="flex items-center gap-0.5" :class="{ 'mx-auto': tab.collectionsCollapsed }">
                  <UButton 
                    v-if="!tab.collectionsCollapsed" 
                    size="xs" 
                    color="neutral" 
                    variant="ghost" 
                    icon="i-heroicons-arrow-path" 
                    @click="refreshCollections(tab)" 
                    title="Refresh Collections" 
                  />
                  <UButton 
                    size="xs" 
                    color="neutral" 
                    variant="ghost" 
                    :icon="tab.collectionsCollapsed ? 'i-heroicons-chevron-double-right' : 'i-heroicons-chevron-double-left'" 
                    @click="tab.collectionsCollapsed = !tab.collectionsCollapsed" 
                    :title="tab.collectionsCollapsed ? 'Expand Collections Panel' : 'Collapse Collections Panel'" 
                  />
                </div>
              </div>

              <!-- Collapsed quick expand trigger -->
              <div 
                v-if="tab.collectionsCollapsed" 
                class="flex-1 flex flex-col items-center py-4 gap-3 text-zinc-500 hover:text-amber-400 cursor-pointer transition-colors"
                @click="tab.collectionsCollapsed = false"
                title="Click to expand collections"
              >
                <span class="text-base">📁</span>
                <span class="text-[10px] font-mono [writing-mode:vertical-lr] tracking-widest uppercase rotate-180 select-none">
                  Colls ({{ tab.collections ? tab.collections.length : 0 }})
                </span>
              </div>

              <template v-else>
                <!-- Collection Search Input -->
                <div class="p-2 border-b border-white/8 bg-[#0b0d13]">
                  <UInput
                    v-model="tab.collectionSearch"
                    placeholder="Filter collections..."
                    icon="i-heroicons-magnifying-glass"
                    size="xs"
                    class="w-full bg-[#141721] border-white/10 text-xs"
                  />
                </div>

                <!-- Collections List -->
                <div class="flex-1 overflow-y-auto p-2 space-y-1">
                  <div v-if="tab.loadingCollections" class="text-xs text-zinc-400 p-6 text-center space-y-2">
                    <UIcon name="i-heroicons-arrow-path" class="w-5 h-5 animate-spin mx-auto text-amber-400" />
                    <span>Loading collections...</span>
                  </div>
                  <div v-else-if="!tab.collections || tab.collections.length === 0" class="text-xs text-zinc-500 p-6 text-center">
                    No collections found.
                  </div>
                  
                  <div 
                    v-for="col in getFilteredCollections(tab)" 
                    :key="col" 
                    class="p-2 rounded-lg text-xs font-semibold cursor-pointer text-zinc-300 flex justify-between items-center transition-all border"
                    :class="tab.activeCollection === col 
                      ? 'bg-[#191c28] border-amber-500/40 text-amber-300 shadow-sm' 
                      : 'hover:bg-white/4 border-transparent hover:text-white'"
                    @click="selectCollection(tab, col)"
                  >
                    <div class="flex items-center gap-2 truncate">
                      <span class="text-xs">📁</span>
                      <span class="truncate font-mono">{{ col }}</span>
                    </div>
                    <button 
                      @click.stop="openCollectionInNewTab(tab, col)"
                      class="text-[10px] text-zinc-500 hover:text-amber-400 px-1 py-0.5 rounded transition-colors"
                      title="Open in new tab"
                    >
                      ↗
                    </button>
                  </div>
                </div>
              </template>
            </div>

            <!-- Right Pane: Document Viewer / Editor / JS Queries -->
            <div class="flex-1 flex flex-col overflow-hidden bg-[#090a0f]">
              <!-- Breadcrumb Header -->
              <div class="p-3.5 border-b border-white/8 bg-[#0f1118] flex justify-between items-center shadow-sm">
                <div class="flex items-center gap-2 flex-wrap min-w-0">
                  <span class="text-xs font-bold uppercase tracking-wider text-amber-400 flex items-center gap-1">
                    <span>⚡</span> {{ tab.firestoreView || 'Explorer' }}
                  </span>
                  
                  <div v-if="tab.activeCollection" class="flex items-center gap-1.5 text-xs text-zinc-400">
                    <span class="text-zinc-600">/</span>
                    <template v-for="(seg, idx) in getCollectionBreadcrumbs(tab.activeCollection)" :key="idx">
                      <span v-if="idx > 0" class="text-zinc-600">&gt;</span>
                      <button 
                        v-if="seg.isCollection" 
                        class="text-amber-400 hover:underline font-mono text-xs cursor-pointer font-bold px-1.5 py-0.5 rounded bg-amber-500/10 border border-amber-500/20"
                        @click="selectCollection(tab, seg.path)"
                        title="Navigate to collection"
                      >
                        {{ seg.name }}
                      </button>
                      <span v-else class="text-zinc-200 font-mono text-xs font-semibold px-1 py-0.5">
                        {{ seg.name }}
                      </span>
                    </template>
                  </div>
                </div>

                <div class="flex items-center gap-2">
                  <UButton size="xs" color="neutral" variant="subtle" icon="i-heroicons-arrow-path-rounded-square" @click="openCopyWizard(tab)">
                    Copy Wizard
                  </UButton>
                  <UButton v-if="tab.activeCollection" size="xs" color="warning" icon="i-heroicons-plus" @click="openAddDocModal(tab)">
                    Add Document
                  </UButton>
                </div>
              </div>

              <!-- Main Viewport -->
              <div class="flex-1 overflow-hidden">
                <!-- Explorer View -->
                <div v-if="tab.firestoreView === 'explorer'" class="h-full flex flex-col p-5 space-y-3 overflow-hidden">
                  <div v-if="!tab.activeCollection" class="h-full flex flex-col justify-center items-center text-zinc-500 text-sm space-y-2">
                    <div class="text-3xl">📁</div>
                    <div class="font-medium text-zinc-300">Select a collection on the left sidebar</div>
                    <div class="text-xs text-zinc-500">Choose any collection to browse and query documents.</div>
                  </div>
                  <div v-else class="flex-1 flex flex-col space-y-3 min-h-0">
                    
                    <!-- Collapsible Query Builder Panel -->
                    <div class="border border-white/8 bg-[#0f1118] rounded-xl p-3.5 space-y-3 flex flex-col shadow-sm">
                      <div class="flex items-center justify-between cursor-pointer select-none" @click="tab.queryExpanded = !tab.queryExpanded">
                        <div class="flex items-center gap-2">
                          <span class="text-xs font-bold text-white uppercase tracking-wider flex items-center gap-1.5">
                            <span class="text-amber-400">🔍</span> Query Builder
                          </span>
                          <span class="text-[10px] bg-amber-500/15 text-amber-300 px-2 py-0.5 rounded font-mono border border-amber-500/25">
                            {{ tab.activeCollection }}
                          </span>
                        </div>
                        <div class="flex items-center gap-2" @click.stop>
                          <div class="flex items-center gap-1 bg-[#090a0f] p-0.5 rounded-lg border border-white/10">
                            <button 
                              type="button"
                              class="px-2.5 py-1 text-[11px] font-semibold rounded-md transition-all"
                              :class="tab.queryMode === 'simple' ? 'bg-amber-500 text-black font-bold' : 'text-zinc-400 hover:text-white'"
                              @click="setQueryMode(tab, 'simple')"
                            >
                              Simple Query
                            </button>
                            <button 
                              type="button"
                              class="px-2.5 py-1 text-[11px] font-semibold rounded-md transition-all"
                              :class="tab.queryMode === 'js' ? 'bg-amber-500 text-black font-bold' : 'text-zinc-400 hover:text-white'"
                              @click="setQueryMode(tab, 'js')"
                            >
                              JS Sandbox
                            </button>
                          </div>
                          <span class="text-zinc-400 hover:text-white text-xs p-1" @click="tab.queryExpanded = !tab.queryExpanded">
                            {{ tab.queryExpanded ? '▲' : '▼' }}
                          </span>
                        </div>
                      </div>

                      <!-- Collapsible Body -->
                      <div v-show="tab.queryExpanded !== false">
                        <!-- Simple Query Form -->
                        <div v-if="tab.queryMode === 'simple'" class="grid grid-cols-4 gap-3 text-xs pt-1">
                          <div class="flex flex-col gap-1">
                            <label class="font-semibold text-zinc-400 text-[11px]">Document ID</label>
                            <UInput v-model="tab.queryDocId" placeholder="Optional ID" size="xs" class="w-full bg-[#090a0f]" />
                          </div>
                          <div class="flex flex-col gap-1">
                            <label class="font-semibold text-zinc-400 text-[11px]">Filter Field</label>
                            <UInput v-model="tab.queryFilterField" placeholder="e.g. role" size="xs" class="w-full bg-[#090a0f]" :list="'filter-fields-list-' + tab.id" />
                            <datalist :id="'filter-fields-list-' + tab.id">
                              <option v-for="field in getAvailableFields(tab)" :key="field" :value="field">{{ field }}</option>
                            </datalist>
                          </div>
                          <div class="flex flex-col gap-1">
                            <label class="font-semibold text-zinc-400 text-[11px]">Operator</label>
                            <USelect 
                              v-model="tab.queryFilterOp" 
                              :items="['==', '>', '>=', '<', '<=', '!=', 'array-contains']" 
                              size="xs"
                              class="w-full bg-[#090a0f]" 
                            />
                          </div>
                          <div class="flex flex-col gap-1">
                            <label class="font-semibold text-zinc-400 text-[11px]">Value</label>
                            <UInput v-model="tab.queryFilterVal" placeholder="Filter value" size="xs" class="w-full bg-[#090a0f]" />
                          </div>
                          <div class="flex flex-col gap-1">
                            <label class="font-semibold text-zinc-400 text-[11px]">Sort Field</label>
                            <UInput v-model="tab.querySortField" placeholder="e.g. createdAt" size="xs" class="w-full bg-[#090a0f]" :list="'sort-fields-list-' + tab.id" />
                            <datalist :id="'sort-fields-list-' + tab.id">
                              <option v-for="field in getAvailableFields(tab)" :key="field" :value="field">{{ field }}</option>
                            </datalist>
                          </div>
                          <div class="flex flex-col gap-1">
                            <label class="font-semibold text-zinc-400 text-[11px]">Direction</label>
                            <USelect 
                              v-model="tab.querySortDir" 
                              :items="[{label:'Ascending',value:'asc'}, {label:'Descending',value:'desc'}]" 
                              size="xs"
                              class="w-full bg-[#090a0f]" 
                            />
                          </div>
                          <div class="flex flex-col gap-1">
                            <label class="font-semibold text-zinc-400 text-[11px]">Limit</label>
                            <UInput v-model.number="tab.queryLimit" type="number" placeholder="50" size="xs" class="w-full bg-[#090a0f]" />
                          </div>
                          <div class="flex items-end justify-end">
                            <UButton color="warning" size="xs" @click="runSimpleQueryAction(tab)" :loading="tab.loadingDocuments" class="w-full justify-center font-semibold">
                              Run Query
                            </UButton>
                          </div>
                        </div>

                        <!-- JS Query Form -->
                        <div v-else class="space-y-3 flex flex-col pt-1">
                          <div class="flex justify-between items-center text-xs">
                            <div class="flex items-center gap-2">
                              <label class="font-semibold text-zinc-400">Preset Template:</label>
                              <USelect 
                                v-model="tab.queryPreset" 
                                :items="['Basic Fetch', 'Filter & Order', 'Projection', 'Aggregate Sum', 'Subcollection Fetch']" 
                                size="xs"
                                class="bg-[#090a0f] shrink-0" 
                                @update:model-value="loadPresetQueryScript(tab)"
                              />
                            </div>
                            <span class="text-[10px] text-zinc-500 font-mono">Runs local mock JS engine against DB</span>
                          </div>
                          <div class="grid grid-cols-2 gap-4">
                            <!-- Editor -->
                            <div class="flex flex-col space-y-2">
                              <UTextarea 
                                v-model="tab.queryScript" 
                                class="font-mono text-xs bg-[#090a0f] flex-1 border-white/10" 
                                rows="8"
                                placeholder="// JavaScript query script..." 
                              />
                              <div class="flex justify-end">
                                <UButton color="warning" size="xs" @click="runQueryScript(tab)" :loading="tab.runningQuery">
                                  Run Script
                                </UButton>
                              </div>
                            </div>
                            <!-- Console Output -->
                            <div class="flex flex-col space-y-1">
                              <label class="text-[10px] font-semibold text-zinc-400">Query Output Logs</label>
                              <pre class="flex-1 font-mono text-[10px] text-emerald-400 overflow-auto bg-[#090a0f] p-3 rounded-lg border border-white/10 whitespace-pre-wrap min-h-40">{{ tab.queryOutput || 'No output yet.' }}</pre>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>

                    <!-- Results Section (Table / Tree / JSON) -->
                    <div class="flex-1 flex flex-col space-y-2 min-h-0">
                      <!-- Results Navigation Tabs and Selection Toolbar -->
                      <div class="flex items-center justify-between border-b border-white/8 pb-2">
                        <div class="flex items-center gap-1 bg-[#141721] p-1 rounded-lg border border-white/8">
                          <button 
                            type="button"
                            class="px-3 py-1 text-xs font-semibold rounded-md transition-all flex items-center gap-1"
                            :class="tab.resultsView === 'tree' ? 'bg-amber-500 text-black font-bold' : 'text-zinc-400 hover:text-white'"
                            @click="tab.resultsView = 'tree'"
                          >
                            <span>🌲</span> Tree
                          </button>
                          <button 
                            type="button"
                            class="px-3 py-1 text-xs font-semibold rounded-md transition-all flex items-center gap-1"
                            :class="tab.resultsView === 'table' ? 'bg-amber-500 text-black font-bold' : 'text-zinc-400 hover:text-white'"
                            @click="tab.resultsView = 'table'"
                          >
                            <span>📊</span> Table
                          </button>
                          <button 
                            type="button"
                            class="px-3 py-1 text-xs font-semibold rounded-md transition-all flex items-center gap-1"
                            :class="tab.resultsView === 'json' ? 'bg-amber-500 text-black font-bold' : 'text-zinc-400 hover:text-white'"
                            @click="tab.resultsView = 'json'"
                          >
                            <span>{ }</span> JSON
                          </button>
                        </div>

                        <!-- Document Count & Selection Action Toolbar -->
                        <div class="flex items-center gap-3">
                          <span class="text-xs text-zinc-400 font-mono">
                            {{ tab.documents ? tab.documents.length : 0 }} document(s)
                          </span>

                          <span class="text-xs text-amber-400 font-semibold" v-if="tab.selectedDocNames && tab.selectedDocNames.length > 0">
                            {{ tab.selectedDocNames.length }} selected
                          </span>
                          
                          <div v-if="tab.selectedDocNames && tab.selectedDocNames.length > 0" class="flex gap-2">
                            <UButton 
                              size="xs" 
                              color="danger" 
                              @click="deleteSelectedDocuments(tab)"
                            >
                              Delete
                            </UButton>
                            <UButton 
                              size="xs" 
                              color="neutral" 
                              variant="subtle"
                              @click="triggerExportSelected(tab)"
                            >
                              Export
                            </UButton>
                          </div>
                          <UButton 
                            size="xs" 
                            color="neutral" 
                            variant="ghost"
                            icon="i-heroicons-arrow-down-tray"
                            @click="triggerExportAll(tab)"
                            title="Export All Documents"
                          />
                        </div>
                      </div>

                      <!-- Tree Viewer -->
                      <div v-if="tab.resultsView === 'tree'" class="flex-1 border border-white/8 bg-[#0f1118] rounded-xl p-4 overflow-y-auto flex flex-col gap-2 min-h-0 shadow-inner">
                        <div class="flex items-center gap-2 mb-2 pb-2 border-b border-white/8">
                          <input 
                            type="checkbox" 
                            :checked="tab.selectedDocNames && tab.selectedDocNames.length === tab.documents.length && tab.documents.length > 0"
                            @change="toggleSelectAllDocs(tab)"
                            class="accent-amber-500 rounded"
                          />
                          <span class="text-xs text-zinc-400">Select All Documents</span>
                        </div>
                        
                        <div v-if="tab.loadingDocuments" class="text-xs text-zinc-400 text-center p-8 space-y-2">
                          <UIcon name="i-heroicons-arrow-path" class="w-6 h-6 animate-spin mx-auto text-amber-400" />
                          <span>Loading documents...</span>
                        </div>
                        <div v-else-if="!tab.documents || tab.documents.length === 0" class="text-xs text-zinc-500 text-center p-8">
                          No documents found in this collection.
                        </div>
                        <div v-else class="space-y-2">
                          <div v-for="doc in tab.documents" :key="doc.name" class="border border-white/6 hover:border-amber-500/30 rounded-xl p-2.5 bg-[#141721]/60 flex items-start gap-3 transition-all">
                            <input 
                              type="checkbox" 
                              :checked="tab.selectedDocNames && tab.selectedDocNames.includes(doc.name)"
                              @change="toggleSelectDoc(tab, doc.name)"
                              class="accent-amber-500 mt-1.5"
                            />
                            <div class="flex-1 min-w-0">
                              <JsonTree 
                                :label="getDocId(doc.name)"
                                :value="cleanDocFields(doc.fields)"
                                :is-root="true"
                                :subcollections="getDocSubcollections(tab, doc)"
                                :loading-subcollections="getDocSubcollectionsLoading(tab, doc)"
                                :subcollections-fetched="getDocSubcollectionsFetched(tab, doc)"
                                @fetch-subcollections="fetchSubcollectionsForDoc(tab, doc)"
                                @open-subcollection="(subCol) => openSubcollection(tab, doc, subCol)"
                                @open-subcollection-in-new-tab="(subCol) => openSubcollectionInNewTab(tab, doc, subCol)"
                                @change="(fieldPath, val) => saveDocFieldChange(tab, doc, fieldPath, val)"
                                @add-field="(parentPath, fieldName, fieldType, initialVal) => addDocField(tab, doc, parentPath, fieldName, fieldType, initialVal)"
                                @delete-field="(fieldPath) => deleteDocField(tab, doc, fieldPath)"
                                @editJson="editDocumentJSON(tab, doc)"
                                @deleteDoc="deleteDocumentConfirm(tab, doc)"
                                @openTab="openDocInNewTab(tab, doc)"
                              />
                            </div>
                          </div>
                        </div>
                      </div>

                      <!-- Table Viewer -->
                      <div v-else-if="tab.resultsView === 'table'" class="flex-1 border border-white/8 bg-[#0f1118] rounded-xl overflow-hidden flex flex-col min-h-0 shadow-md">
                        <div class="flex-1 overflow-auto">
                          <table class="w-full text-left border-collapse text-xs">
                            <thead>
                              <tr class="bg-[#141721] text-zinc-300 border-b border-white/8 font-semibold select-none sticky top-0 z-10 backdrop-blur">
                                <th class="p-3 w-10 text-center">
                                  <input 
                                    type="checkbox" 
                                    :checked="tab.selectedDocNames && tab.selectedDocNames.length === tab.documents.length && tab.documents.length > 0"
                                    @change="toggleSelectAllDocs(tab)"
                                    class="accent-amber-500 rounded"
                                  />
                                </th>
                                <th class="p-3 w-1/4">Document ID</th>
                                <th v-for="col in getTableColumns(tab.documents)" :key="col" class="p-3 font-mono text-zinc-300">{{ col }}</th>
                                <th class="p-3 text-right">Actions</th>
                              </tr>
                            </thead>
                            <tbody>
                              <tr v-if="tab.loadingDocuments" class="text-center">
                                <td :colspan="getTableColumns(tab.documents).length + 3" class="p-12 text-zinc-400">
                                  <UIcon name="i-heroicons-arrow-path" class="w-6 h-6 animate-spin mx-auto text-amber-400 mb-2" />
                                  <span>Loading documents...</span>
                                </td>
                              </tr>
                              <tr v-else-if="!tab.documents || tab.documents.length === 0" class="text-center">
                                <td :colspan="getTableColumns(tab.documents).length + 3" class="p-12 text-zinc-500">No documents found.</td>
                              </tr>
                              <tr 
                                v-else 
                                v-for="doc in tab.documents" 
                                :key="doc.name" 
                                class="border-b border-white/4 hover:bg-white/3 transition-colors cursor-pointer group/row"
                                @click="toggleSelectDoc(tab, doc.name)"
                              >
                                <td class="p-3 text-center" @click.stop>
                                  <input 
                                    type="checkbox" 
                                    :checked="tab.selectedDocNames && tab.selectedDocNames.includes(doc.name)"
                                    @change="toggleSelectDoc(tab, doc.name)"
                                    class="accent-amber-500 rounded"
                                  />
                                </td>
                                <td class="p-3 font-mono font-semibold text-amber-400">
                                  <div class="flex items-center gap-1.5">
                                    <span class="truncate max-w-35 cursor-pointer hover:underline" :title="getDocId(doc.name)" @click.stop="openDocInNewTab(tab, doc)">
                                      {{ getDocId(doc.name) }}
                                    </span>
                                    <button
                                      type="button"
                                      class="opacity-0 group-hover/row:opacity-100 text-zinc-500 hover:text-amber-400 p-0.5 rounded transition-opacity"
                                      @click.stop="copyText(getDocId(doc.name), 'Document ID copied!')"
                                      title="Copy Document ID"
                                    >
                                      📋
                                    </button>
                                  </div>
                                </td>
                                <td v-for="col in getTableColumns(tab.documents)" :key="col" class="p-3 font-mono text-[11px] text-zinc-400 truncate max-w-37.5">
                                  {{ getColumnValue(doc, col) }}
                                </td>
                                <td class="p-3 text-right space-x-1 whitespace-nowrap" @click.stop>
                                  <UPopover v-if="getDocSubcollectionsFetched(tab, doc) || getDocSubcollectionsLoading(tab, doc)">
                                    <UButton size="xs" variant="subtle" color="warning" class="text-[10px]">
                                      📁 {{ getDocSubcollections(tab, doc).length }} Subcoll
                                    </UButton>
                                    <template #content>
                                      <div class="p-3 space-y-1 bg-[#141721] border border-white/10 rounded-xl shadow-2xl text-xs max-w-xs text-left">
                                        <div class="font-bold text-amber-400 border-b border-white/10 pb-1.5 flex justify-between items-center">
                                          <span>Subcollections</span>
                                          <UButton size="2xs" color="neutral" variant="ghost" icon="i-heroicons-arrow-path" @click="fetchSubcollectionsForDoc(tab, doc)" />
                                        </div>
                                        <div v-if="getDocSubcollectionsLoading(tab, doc)" class="text-xs text-zinc-400 p-2">Loading...</div>
                                        <div v-else-if="getDocSubcollections(tab, doc).length === 0" class="text-xs text-zinc-500 p-2">No subcollections found</div>
                                        <div v-else v-for="subCol in getDocSubcollections(tab, doc)" :key="subCol" class="flex justify-between items-center gap-2 p-1.5 hover:bg-white/5 rounded-lg">
                                          <span class="font-mono text-amber-300 text-xs">{{ subCol }}</span>
                                          <div class="flex gap-1">
                                            <UButton size="2xs" color="warning" variant="ghost" @click="openSubcollection(tab, doc, subCol)">Open</UButton>
                                            <UButton size="2xs" color="neutral" variant="ghost" @click="openSubcollectionInNewTab(tab, doc, subCol)">Tab</UButton>
                                          </div>
                                        </div>
                                      </div>
                                    </template>
                                  </UPopover>
                                  <UButton v-else size="xs" variant="ghost" color="neutral" icon="i-heroicons-folder" title="Fetch Subcollections" @click="fetchSubcollectionsForDoc(tab, doc)">
                                    Subcolls
                                  </UButton>

                                  <UButton size="xs" variant="ghost" color="neutral" @click="openDocInNewTab(tab, doc)">Tab</UButton>
                                  <UButton size="xs" variant="ghost" color="neutral" @click="editDocumentJSON(tab, doc)">JSON</UButton>
                                  <UButton size="xs" variant="ghost" color="danger" @click="deleteDocumentConfirm(tab, doc)">Delete</UButton>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </div>
                      </div>

                      <!-- JSON Viewer -->
                      <div v-else-if="tab.resultsView === 'json'" class="flex-1 border border-white/8 bg-[#0f1118] rounded-xl p-4 overflow-hidden flex flex-col min-h-0 shadow-inner">
                        <textarea 
                          class="grow flex-1 w-full h-full bg-[#090a0f] text-zinc-300 font-mono text-xs p-3 rounded-lg border border-white/10 resize-none outline-none select-text min-h-0" 
                          readonly
                        >{{ formatAllDocsJson(tab.documents) }}</textarea>
                      </div>

                      <!-- Load More Button -->
                      <div v-if="tab.nextDocPageToken" class="p-3 bg-[#141721] border-t border-white/8 flex justify-end rounded-b-xl">
                        <UButton size="xs" color="warning" @click="fetchDocumentsNext(tab)">
                          Load More Documents
                        </UButton>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- Modal: Add Service Account -->
    <UModal v-model:open="showAddSaModal" :dismissible="false">
      <template #content>
        <UCard :ui="{ background: 'bg-[#141721]', border: 'border-white/10' }">
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <span class="text-amber-400 text-lg">🔑</span>
                <h3 class="text-base font-bold text-white">Add Google Service Account</h3>
              </div>
              <UButton color="neutral" variant="ghost" icon="i-heroicons-x-mark" class="-my-1" @click="showAddSaModal = false" />
            </div>
          </template>

          <div class="space-y-4 p-2 text-xs">
            <p class="text-zinc-400 leading-relaxed">
              Paste the contents of your Google Service Account private key JSON file below. 
              This file will be loaded <strong class="text-white">only in-memory</strong> and will not be saved to disk.
            </p>
            <UTextarea 
              v-model="newSaJson" 
              placeholder='{ "type": "service_account", "project_id": ... }' 
              rows="10" 
              class="font-mono text-xs bg-[#090a0f] border-white/10" 
            />
          </div>

          <template #footer>
            <div class="flex justify-end gap-2">
              <UButton color="neutral" variant="ghost" @click="showAddSaModal = false">Cancel</UButton>
              <UButton color="warning" @click="saveServiceAccount" :loading="savingSa">Add Account</UButton>
            </div>
          </template>
        </UCard>
      </template>
    </UModal>

    <!-- Modal: Edit Claims -->
    <UModal v-model:open="showClaimsModal" :dismissible="false">
      <template #content>
        <UCard :ui="{ background: 'bg-[#141721]', border: 'border-white/10' }">
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <span class="text-purple-400 text-lg">⚡</span>
                <h3 class="text-base font-bold text-white">Edit Custom Claims</h3>
              </div>
              <UButton color="neutral" variant="ghost" icon="i-heroicons-x-mark" class="-my-1" @click="showClaimsModal = false" />
            </div>
          </template>

          <div class="space-y-4 p-2 text-xs">
            <p class="text-zinc-400">
              User: <strong class="text-amber-400 font-mono">{{ activeUser?.email || activeUser?.uid }}</strong>. Custom claims must be valid JSON object or empty `{}`.
            </p>
            <UTextarea 
              v-model="claimsJson" 
              placeholder='{ "admin": true, "role": "editor" }' 
              rows="7" 
              class="font-mono text-xs bg-[#090a0f] border-white/10" 
            />
          </div>

          <template #footer>
            <div class="flex justify-end gap-2">
              <UButton color="neutral" variant="ghost" @click="showClaimsModal = false">Cancel</UButton>
              <UButton color="warning" @click="saveCustomClaims">Save Claims</UButton>
            </div>
          </template>
        </UCard>
      </template>
    </UModal>

    <!-- Modal: Create Auth User -->
    <UModal v-model:open="showCreateUserModal" :dismissible="false">
      <template #content>
        <UCard :ui="{ background: 'bg-[#141721]', border: 'border-white/10' }">
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <span class="text-amber-400 text-lg">👤</span>
                <h3 class="text-base font-bold text-white">Create Firebase Auth User</h3>
              </div>
              <UButton color="neutral" variant="ghost" icon="i-heroicons-x-mark" class="-my-1" @click="showCreateUserModal = false" />
            </div>
          </template>

          <div class="space-y-4 p-2 text-xs">
            <div>
              <label class="block text-xs font-semibold text-zinc-300 mb-1">Email Address</label>
              <UInput v-model="newUser.email" placeholder="user@example.com" class="w-full bg-[#090a0f]" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-zinc-300 mb-1">Display Name</label>
              <UInput v-model="newUser.displayName" placeholder="John Doe" class="w-full bg-[#090a0f]" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-zinc-300 mb-1">Password</label>
              <UInput v-model="newUser.password" type="password" placeholder="••••••••" class="w-full bg-[#090a0f]" />
            </div>
            <div class="flex items-center gap-4 pt-1">
              <UCheckbox v-model="newUser.emailVerified" label="Email Verified" />
              <UCheckbox v-model="newUser.disabled" label="Disabled" />
            </div>
          </div>

          <template #footer>
            <div class="flex justify-end gap-2">
              <UButton color="neutral" variant="ghost" @click="showCreateUserModal = false">Cancel</UButton>
              <UButton color="warning" @click="saveNewUser">Create User</UButton>
            </div>
          </template>
        </UCard>
      </template>
    </UModal>

    <!-- Modal: Edit Auth User Basic -->
    <UModal v-model:open="showEditUserModal" :dismissible="false">
      <template #content>
        <UCard :ui="{ background: 'bg-[#141721]', border: 'border-white/10' }">
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <span class="text-amber-400 text-lg">✏️</span>
                <h3 class="text-base font-bold text-white">Edit User Profile</h3>
              </div>
              <UButton color="neutral" variant="ghost" icon="i-heroicons-x-mark" class="-my-1" @click="showEditUserModal = false" />
            </div>
          </template>

          <div class="space-y-4 p-2 text-xs">
            <div>
              <label class="block text-xs font-semibold text-zinc-300 mb-1">Email Address</label>
              <UInput v-model="editUser.email" placeholder="user@example.com" class="w-full bg-[#090a0f]" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-zinc-300 mb-1">Display Name</label>
              <UInput v-model="editUser.displayName" placeholder="John Doe" class="w-full bg-[#090a0f]" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-zinc-300 mb-1">Password (Leave empty to keep unchanged)</label>
              <UInput v-model="editUser.password" type="password" placeholder="••••••••" class="w-full bg-[#090a0f]" />
            </div>
            <div class="flex items-center gap-4 pt-1">
              <UCheckbox v-model="editUser.emailVerified" label="Email Verified" />
              <UCheckbox v-model="editUser.disabled" label="Disabled" />
            </div>
          </div>

          <template #footer>
            <div class="flex justify-end gap-2">
              <UButton color="neutral" variant="ghost" @click="showEditUserModal = false">Cancel</UButton>
              <UButton color="warning" @click="saveUserEdit">Save Changes</UButton>
            </div>
          </template>
        </UCard>
      </template>
    </UModal>

    <!-- Modal: Export Auth Users -->
    <UModal v-model:open="showExportAuthModal" :dismissible="false" :ui="{ content: 'sm:max-w-xl' }">
      <template #content>
        <UCard :ui="{ background: 'bg-[#141721]', border: 'border-white/10' }">
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <span class="text-amber-400 text-lg">📤</span>
                <h3 class="text-base font-bold text-white">Export Auth Users & Custom Claims</h3>
              </div>
              <UButton color="neutral" variant="ghost" icon="i-heroicons-x-mark" class="-my-1" @click="showExportAuthModal = false" />
            </div>
          </template>

          <div class="space-y-4 p-2 text-xs">
            <div v-if="exportAuthLoading" class="py-8 text-center text-zinc-400 space-y-2">
              <UIcon name="i-heroicons-arrow-path" class="w-6 h-6 animate-spin mx-auto text-amber-400" />
              <p>Exporting all Auth users and custom claims from project...</p>
            </div>

            <div v-else-if="exportAuthData" class="space-y-4">
              <!-- Summary Card -->
              <div class="p-3 bg-[#090a0f] border border-white/10 rounded-xl flex items-center justify-between">
                <div>
                  <div class="text-sm font-semibold text-white">
                    {{ exportAuthData.count }} Users Ready to Export
                  </div>
                  <div class="text-[11px] text-zinc-400">
                    Includes UIDs, emails, providers, and {{ countUsersWithClaims(exportAuthData.users) }} user(s) with custom claims.
                  </div>
                </div>
                <UBadge color="warning" variant="subtle" size="sm">JSON Format</UBadge>
              </div>

              <!-- Password Hash Parameters Section -->
              <div class="p-3 bg-[#090a0f] border border-white/10 rounded-xl space-y-3">
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-1.5 font-semibold text-amber-400">
                    <span>🔑</span>
                    <span>Password Hash Parameters</span>
                  </div>
                  <span class="text-[10px] text-zinc-400">Recommended for password migration</span>
                </div>
                <p class="text-zinc-400 leading-relaxed text-[11px]">
                  Firebase encrypts passwords using modified SCRYPT. To allow users to sign in without password resets in destination projects, configure hash parameters:
                </p>

                <div class="grid grid-cols-2 gap-x-4 gap-y-3 pt-1">
                  <div>
                    <div class="flex items-center justify-between mb-1 text-[11px] font-medium text-zinc-300">
                      <span>Signer Key</span>
                      <code class="text-amber-400 text-[10px]">--hash-key</code>
                    </div>
                    <UInput v-model="exportAuthHashParams.signerKey" placeholder="Base64 signer key" class="w-full bg-[#141721]" />
                  </div>
                  <div>
                    <div class="flex items-center justify-between mb-1 text-[11px] font-medium text-zinc-300">
                      <span>Salt Separator</span>
                      <code class="text-amber-400 text-[10px]">--salt-separator</code>
                    </div>
                    <UInput v-model="exportAuthHashParams.saltSeparator" placeholder="Bw==" class="w-full bg-[#141721]" />
                  </div>
                  <div>
                    <div class="flex items-center justify-between mb-1 text-[11px] font-medium text-zinc-300">
                      <span>Rounds</span>
                      <code class="text-amber-400 text-[10px]">--rounds</code>
                    </div>
                    <UInput v-model.number="exportAuthHashParams.rounds" type="number" placeholder="8" class="w-full bg-[#141721]" />
                  </div>
                  <div>
                    <div class="flex items-center justify-between mb-1 text-[11px] font-medium text-zinc-300">
                      <span>Memory Cost</span>
                      <code class="text-amber-400 text-[10px]">--mem-cost</code>
                    </div>
                    <UInput v-model.number="exportAuthHashParams.memoryCost" type="number" placeholder="14" class="w-full bg-[#141721]" />
                  </div>
                </div>
              </div>
            </div>
          </div>

          <template #footer>
            <div class="flex justify-between items-center w-full">
              <div class="text-[11px] text-zinc-500 truncate max-w-50">
                Custom claims & hashes preserved.
              </div>
              <div class="flex items-center gap-2 whitespace-nowrap">
                <UButton color="neutral" variant="ghost" size="sm" @click="showExportAuthModal = false">
                  Cancel
                </UButton>
                <UButton color="neutral" variant="subtle" size="sm" icon="i-heroicons-clipboard-document" @click="copyExportAuthJson" :disabled="exportAuthLoading || !exportAuthData">
                  Copy JSON
                </UButton>
                <UButton color="warning" size="sm" icon="i-heroicons-arrow-down-tray" class="font-semibold" @click="downloadExportAuthJson" :disabled="exportAuthLoading || !exportAuthData">
                  Download JSON
                </UButton>
              </div>
            </div>
          </template>
        </UCard>
      </template>
    </UModal>

    <!-- Modal: Import Auth Users -->
    <UModal v-model:open="showImportAuthModal" :dismissible="false" :ui="{ content: 'sm:max-w-xl' }">
      <template #content>
        <UCard :ui="{ background: 'bg-[#141721]', border: 'border-white/10' }">
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <span class="text-amber-400 text-lg">📥</span>
                <h3 class="text-base font-bold text-white">Import Auth Users & Custom Claims</h3>
              </div>
              <UButton color="neutral" variant="ghost" icon="i-heroicons-x-mark" class="-my-1" @click="showImportAuthModal = false" />
            </div>
          </template>

          <div class="space-y-4 p-2 text-xs">
            <div>
              <div class="flex justify-between items-center mb-1.5">
                <label class="font-semibold text-zinc-300">1. Load Users Data (JSON)</label>
                <input type="file" id="auth-import-file-input" accept=".json" class="hidden" @change="onImportFileSelected" />
                <UButton size="xs" color="neutral" variant="subtle" icon="i-heroicons-folder-open" @click="triggerImportFilePicker">
                  Upload JSON File
                </UButton>
              </div>
              <UTextarea
                v-model="importAuthJsonInput"
                placeholder='Paste JSON here or upload file. Supports Mikhaz Fire and Firebase CLI formats.'
                rows="5"
                class="w-full font-mono text-[11px] bg-[#090a0f] border-white/10"
                @input="parseImportJsonInput"
              />
            </div>

            <!-- JSON Parsing Status Summary -->
            <div v-if="parsedImportSummary" class="p-3 bg-[#090a0f] border border-white/10 rounded-xl flex items-center justify-between">
              <div>
                <div class="text-sm font-semibold text-white">
                  {{ parsedImportSummary.userCount }} User(s) Detected
                </div>
                <div class="text-[11px] text-zinc-400">
                  {{ parsedImportSummary.claimsCount }} with claims &middot; {{ parsedImportSummary.hashCount }} with password hashes.
                </div>
              </div>
              <UBadge :color="parsedImportSummary.userCount > 0 ? 'success' : 'danger'" variant="subtle" size="sm">
                {{ parsedImportSummary.userCount > 0 ? 'Valid JSON' : 'Invalid Data' }}
              </UBadge>
            </div>

            <!-- Hash Parameters Section -->
            <div class="p-3 bg-[#090a0f] border border-white/10 rounded-xl space-y-3">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-1.5 font-semibold text-amber-400">
                  <span>🔐</span>
                  <span>2. Password Hash Parameters</span>
                </div>
                <span v-if="parsedImportSummary && parsedImportSummary.hasEmbeddedParams" class="text-[10px] text-emerald-400 font-semibold">
                  ✓ Auto-filled from JSON
                </span>
              </div>

              <div class="grid grid-cols-2 gap-x-4 gap-y-3 pt-1">
                <div>
                  <div class="mb-1 text-[11px] font-medium text-zinc-300">Hash Algorithm</div>
                  <select v-model="importAuthHashParams.hashAlgorithm" class="w-full bg-[#141721] border border-white/15 rounded-lg px-2.5 py-1.5 text-xs text-white focus:outline-none focus:border-amber-500">
                    <option value="SCRYPT">SCRYPT (Firebase Modified)</option>
                    <option value="STANDARD_SCRYPT">STANDARD_SCRYPT</option>
                    <option value="BCRYPT">BCRYPT</option>
                    <option value="PBKDF2_SHA256">PBKDF2_SHA256</option>
                    <option value="SHA256">SHA256</option>
                  </select>
                </div>
                <div>
                  <div class="mb-1 text-[11px] font-medium text-zinc-300">Signer Key</div>
                  <UInput v-model="importAuthHashParams.signerKey" placeholder="Base64 signer key" class="w-full bg-[#141721]" />
                </div>
                <div>
                  <div class="mb-1 text-[11px] font-medium text-zinc-300">Salt Separator</div>
                  <UInput v-model="importAuthHashParams.saltSeparator" placeholder="Bw==" class="w-full bg-[#141721]" />
                </div>
                <div class="grid grid-cols-2 gap-2">
                  <div>
                    <div class="mb-1 text-[11px] font-medium text-zinc-300">Rounds</div>
                    <UInput v-model.number="importAuthHashParams.rounds" type="number" placeholder="8" class="w-full bg-[#141721]" />
                  </div>
                  <div>
                    <div class="mb-1 text-[11px] font-medium text-zinc-300">Memory Cost</div>
                    <UInput v-model.number="importAuthHashParams.memoryCost" type="number" placeholder="14" class="w-full bg-[#141721]" />
                  </div>
                </div>
              </div>
            </div>

            <!-- Import Results Card -->
            <div v-if="importAuthResult" class="p-3 border rounded-xl space-y-2" :class="importAuthResult.success ? 'bg-emerald-950/20 border-emerald-500/30' : 'bg-rose-950/20 border-rose-500/30'">
              <div class="flex items-center justify-between">
                <div class="font-semibold text-white flex items-center gap-2">
                  <span :class="importAuthResult.success ? 'text-emerald-400' : 'text-amber-400'">
                    {{ importAuthResult.success ? '✓ Import Completed' : '⚠️ Import Completed with Issues' }}
                  </span>
                </div>
                <div class="text-[11px] text-zinc-400">
                  {{ importAuthResult.importedCount }} imported &middot; {{ importAuthResult.failedCount }} failed
                </div>
              </div>
            </div>
          </div>

          <template #footer>
            <div class="flex justify-between items-center w-full">
              <div class="text-[11px] text-zinc-500 truncate max-w-50">
                Target: {{ activeTabContext?.projectId || '' }}
              </div>
              <div class="flex items-center gap-2 whitespace-nowrap">
                <UButton color="neutral" variant="ghost" size="sm" @click="showImportAuthModal = false">
                  Cancel
                </UButton>
                <UButton
                  color="warning"
                  size="sm"
                  icon="i-heroicons-arrow-down-tray"
                  class="font-semibold"
                  :loading="importAuthLoading"
                  :disabled="importAuthLoading || !parsedImportSummary || parsedImportSummary.userCount === 0"
                  @click="runImportAuthUsers(activeTabContext)"
                >
                  Import {{ parsedImportSummary ? parsedImportSummary.userCount : '' }} Users
                </UButton>
              </div>
            </div>
          </template>
        </UCard>
      </template>
    </UModal>

    <!-- Modal: Document Details JSON Viewer/Editor -->
    <UModal v-model:open="showDocModal" :dismissible="false" fullscreen :ui="{ content: 'h-full flex flex-col bg-[#090a0f]' }">
      <template #content>
        <UCard :ui="{ root: 'h-full flex flex-col bg-[#090a0f] border border-white/10', body: 'grow overflow-hidden flex flex-col h-full min-h-0' }">
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <span class="text-amber-400 text-lg">📄</span>
                <h3 class="text-base font-bold text-white">Document JSON Editor</h3>
              </div>
              <UButton color="neutral" variant="ghost" icon="i-heroicons-x-mark" class="-my-1" @click="showDocModal = false" />
            </div>
          </template>

          <div class="space-y-3 p-4 text-xs flex flex-col grow flex-1 h-full min-h-0 bg-[#090a0f]">
            <div class="bg-[#141721] p-2.5 rounded-lg border border-white/10 text-zinc-300 flex items-center justify-between">
              <div>Path: <strong class="text-amber-400 font-mono">{{ activeDocPath }}</strong></div>
              <UButton size="2xs" color="neutral" variant="subtle" @click="formatActiveDocJson">Beautify JSON</UButton>
            </div>
            <textarea 
              v-model="activeDocJson" 
              placeholder='{ "title": "First Post", "likes": 42 }' 
              class="flex-1 w-full bg-[#141721] text-zinc-200 font-mono text-xs p-4 rounded-xl border border-white/10 resize-none outline-none select-text min-h-0 h-full focus:border-amber-500"
            ></textarea>
          </div>

          <template #footer>
            <div class="flex justify-end gap-2 bg-[#090a0f]">
              <UButton color="neutral" variant="ghost" @click="showDocModal = false">Cancel</UButton>
              <UButton color="warning" @click="saveDocFromJSON">Save Document</UButton>
            </div>
          </template>
        </UCard>
      </template>
    </UModal>

    <!-- Modal: Add Document -->
    <UModal v-model:open="showAddDocModal" :dismissible="false">
      <template #content>
        <UCard :ui="{ background: 'bg-[#141721]', border: 'border-white/10' }">
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <span class="text-amber-400 text-lg">+</span>
                <h3 class="text-base font-bold text-white">Create New Document</h3>
              </div>
              <UButton color="neutral" variant="ghost" icon="i-heroicons-x-mark" class="-my-1" @click="showAddDocModal = false" />
            </div>
          </template>

          <div class="space-y-4 p-2 text-xs">
            <div>
              <label class="block text-xs font-semibold text-zinc-300 mb-1">Document ID (optional, blank auto-generates ID)</label>
              <UInput v-model="newDocId" placeholder="my-custom-doc-id" class="w-full bg-[#090a0f]" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-zinc-300 mb-1">Fields Initial JSON</label>
              <UTextarea 
                v-model="newDocFieldsJson" 
                placeholder='{ "title": "New Doc", "createdAt": "2026-09-05T00:00:00Z" }' 
                rows="6" 
                class="font-mono text-xs bg-[#090a0f] border-white/10" 
              />
            </div>
          </div>

          <template #footer>
            <div class="flex justify-end gap-2">
              <UButton color="neutral" variant="ghost" @click="showAddDocModal = false">Cancel</UButton>
              <UButton color="warning" @click="saveNewDoc">Create Document</UButton>
            </div>
          </template>
        </UCard>
      </template>
    </UModal>

    <!-- Modal: Copy Node Wizard -->
    <UModal v-model:open="showCopyModal" :dismissible="false">
      <template #content>
        <UCard :ui="{ background: 'bg-[#141721]', border: 'border-white/10' }">
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <span class="text-amber-400 text-lg">🔄</span>
                <h3 class="text-base font-bold text-white">Firestore Copy Wizard</h3>
              </div>
              <UButton color="neutral" variant="ghost" icon="i-heroicons-x-mark" class="-my-1" @click="showCopyModal = false" />
            </div>
          </template>

          <div class="space-y-4 p-2 text-xs">
            <p class="text-zinc-400">
              Copy collections or documents from this project to another project safely. Checks for conflicts before execution.
            </p>
            <div>
              <label class="block text-xs font-semibold text-zinc-300 mb-1">Source Account & Project</label>
              <div class="bg-[#090a0f] p-2 rounded-lg text-white border border-white/10 font-mono">
                {{ activeCopySource }}
              </div>
            </div>
            <div>
              <label class="block text-xs font-semibold text-zinc-300 mb-1">Source Path (Collection or Doc)</label>
              <UInput v-model="copyWizard.path" placeholder="posts" class="w-full bg-[#090a0f] font-mono" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-zinc-300 mb-1">Copy Node Type</label>
              <USelect 
                v-model="copyWizard.type" 
                :items="['collection', 'document']" 
                class="w-full bg-[#090a0f]"
              />
            </div>
            
            <div class="border-t border-white/10 pt-3">
              <label class="block text-xs font-semibold text-white mb-2">Destination Configuration</label>
              <div class="space-y-2">
                <div>
                  <label class="block text-[10px] text-zinc-400 mb-1">Dest Account Email</label>
                  <USelect 
                    v-model="copyWizard.destEmail" 
                    :items="accounts.map(a => a.email)" 
                    class="w-full bg-[#090a0f]"
                    @change="loadDestProjects"
                  />
                </div>
                <div>
                  <label class="block text-[10px] text-zinc-400 mb-1">Dest Project ID</label>
                  <USelect 
                    v-model="copyWizard.destProjectId" 
                    :items="destProjects.map(p => p.projectId)" 
                    class="w-full bg-[#090a0f]"
                  />
                </div>
              </div>
            </div>

            <div class="flex items-center gap-4 border-t border-white/10 pt-3">
              <UCheckbox v-model="copyWizard.overwrite" label="Overwrite Conflicts" />
              <UCheckbox v-model="copyWizard.recursive" label="Copy Subcollections" />
            </div>
          </div>

          <template #footer>
            <div class="flex justify-end gap-2">
              <UButton color="neutral" variant="ghost" @click="showCopyModal = false">Cancel</UButton>
              <UButton color="warning" @click="runCopyConflictCheck" :loading="copyingNode">Proceed Copy</UButton>
            </div>
          </template>
        </UCard>
      </template>
    </UModal>

    <!-- Modal: Export Results -->
    <UModal v-model:open="showExportModal" :dismissible="false" size="md">
      <template #content>
        <UCard :ui="{ background: 'bg-[#141721]', border: 'border-white/10' }">
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <span class="text-amber-400 text-lg">📤</span>
                <h3 class="text-base font-bold text-white">Export Results</h3>
              </div>
              <UButton color="neutral" variant="ghost" icon="i-heroicons-x-mark" class="-my-1" @click="showExportModal = false" />
            </div>
          </template>

          <div class="space-y-4 p-2 text-xs">
            <p class="text-zinc-400">
              Export the current list of loaded documents. You can configure the format and options below.
            </p>
            
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-xs font-semibold text-zinc-300 mb-1">Export Format</label>
                <USelect v-model="exportConfig.format" :items="['JSON', 'CSV']" class="w-full bg-[#090a0f]" @change="updateExportPreview" />
              </div>
              <div class="flex flex-col gap-2 justify-end pb-1">
                <UCheckbox v-model="exportConfig.includeIds" label="Include Document IDs" @change="updateExportPreview" />
                <UCheckbox v-model="exportConfig.includePaths" label="Include Document Paths" @change="updateExportPreview" />
              </div>
            </div>

            <div class="pt-2" v-if="exportConfig.hasSelection">
              <UCheckbox v-model="exportConfig.onlySelected" :label="`Export only selected documents (${exportConfig.selectedCount} selected)`" @change="updateExportPreview" />
            </div>

            <div class="flex flex-col space-y-1 mt-4">
              <label class="block text-xs font-semibold text-white">Preview (First 5,000 characters)</label>
              <textarea 
                v-model="exportPreviewText"
                class="w-full bg-[#090a0f] text-zinc-300 border border-white/10 font-mono text-[10px] p-2 rounded-lg resize-none select-text"
                rows="8"
                readonly
              />
            </div>
          </div>

          <template #footer>
            <div class="flex justify-end gap-2">
              <UButton color="neutral" variant="ghost" @click="showExportModal = false">Cancel</UButton>
              <UButton color="neutral" variant="subtle" @click="copyExportToClipboard">Copy to Clipboard</UButton>
              <UButton color="warning" @click="downloadExportFile">Download File</UButton>
            </div>
          </template>
        </UCard>
      </template>
    </UModal>

    <!-- Modal: Notification / Error Details -->
    <UModal v-model:open="showNotificationModal" size="md">
      <template #content>
        <UCard :ui="{ background: 'bg-[#141721]', border: notificationModalType === 'error' ? 'border-rose-500/30' : 'border-emerald-500/30' }">
          <template #header>
            <div class="flex items-center justify-between">
              <h3 
                class="text-base font-bold leading-6 flex items-center gap-2"
                :class="notificationModalType === 'error' ? 'text-rose-400' : 'text-emerald-400'"
              >
                <span>{{ notificationModalType === 'error' ? '⚠️' : '✓' }}</span>
                {{ notificationModalTitle || (notificationModalType === 'error' ? 'Error' : 'Notification') }}
              </h3>
              <UButton color="neutral" variant="ghost" icon="i-heroicons-x-mark" class="-my-1" @click="showNotificationModal = false" />
            </div>
          </template>

          <div class="space-y-3 p-2 text-xs">
            <div 
              class="bg-[#090a0f] text-zinc-200 border border-white/10 font-mono text-[11px] p-4 rounded-xl overflow-auto select-text max-h-87.5 whitespace-pre-wrap break-all"
            >
              {{ notificationModalMessage }}
            </div>
          </div>

          <template #footer>
            <div class="flex justify-end">
              <UButton 
                color="warning"
                @click="showNotificationModal = false"
              >
                Dismiss
              </UButton>
            </div>
          </template>
        </UCard>
      </template>
    </UModal>

    <!-- Floating Toast Notification -->
    <div 
      v-if="toastMessage" 
      class="fixed bottom-5 right-5 z-50 px-4 py-2.5 rounded-xl bg-[#181b26] border border-amber-500/40 text-white text-xs shadow-2xl flex items-center gap-2 animate-fade-in"
    >
      <span class="text-amber-400 font-bold">✓</span>
      <span>{{ toastMessage }}</span>
    </div>

    </div>
  </UApp>
</template>

<script setup>
import { ref, onMounted, reactive, computed } from 'vue'

// Import Wails runtime calls
import {
  GetAccounts,
  AddServiceAccount,
  LoginGoogleAccount,
  ListProjects,
  ListAuthUsers,
  CreateAuthUser,
  UpdateAuthUser,
  DeleteAuthUser,
  ExportAuthUsers,
  ImportAuthUsers,
  ListRootCollections,
  ListDocuments,
  ListDocumentSubcollections,
  GetDocument,
  SaveDocument,
  DeleteDocument,
  RunJSScript,
  RunSimpleQuery,
  CheckConflict,
  CopyNode
} from '~/wailsjs/go/main/App'

// Sidebar State
const accounts = ref([])
const projects = ref([])
const expandedAccount = ref(null)
const scanning = ref(false)
const loadingProjects = ref(false)
const loggingInGoogle = ref(false)
const sidebarProjectFilter = ref('')
const pinnedProjects = ref([])

// Toast state
const toastMessage = ref('')
let toastTimer = null

function showToast(msg) {
  toastMessage.value = msg
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => {
    toastMessage.value = ''
  }, 2500)
}

function copyText(text, label = 'Copied to clipboard!') {
  if (!text) return
  navigator.clipboard.writeText(String(text))
  showToast(label)
}

// Local Storage for Pinned Projects
function loadPinnedProjects() {
  try {
    const saved = localStorage.getItem('mikhaz_pinned_projects')
    if (saved) pinnedProjects.value = JSON.parse(saved)
  } catch (e) {}
}

function togglePinProject(project) {
  const idx = pinnedProjects.value.findIndex(p => p.projectId === project.projectId)
  if (idx >= 0) {
    pinnedProjects.value.splice(idx, 1)
    showToast(`Unpinned ${project.displayName || project.projectId}`)
  } else {
    pinnedProjects.value.push({
      projectId: project.projectId,
      displayName: project.displayName || project.projectId,
      email: project.email || expandedAccount.value || ''
    })
    showToast(`Pinned ${project.displayName || project.projectId}!`)
  }
  try {
    localStorage.setItem('mikhaz_pinned_projects', JSON.stringify(pinnedProjects.value))
  } catch (e) {}
}

function isPinned(projectId) {
  return pinnedProjects.value.some(p => p.projectId === projectId)
}

const filteredAccounts = computed(() => {
  return accounts.value
})

const filteredProjects = computed(() => {
  if (!sidebarProjectFilter.value.trim()) return projects.value
  const q = sidebarProjectFilter.value.trim().toLowerCase()
  return projects.value.filter(p => 
    p.projectId.toLowerCase().includes(q) || 
    (p.displayName && p.displayName.toLowerCase().includes(q))
  )
})

// Tabs State
const openTabs = ref([])
const activeTabId = ref('welcome')
const showWelcomeTab = ref(true)

function closeWelcomeTab() {
  showWelcomeTab.value = false
  if (activeTabId.value === 'welcome') {
    activeTabId.value = openTabs.value.length > 0 ? openTabs.value[openTabs.value.length - 1].id : ''
  }
}

function reopenWelcomeTab() {
  showWelcomeTab.value = true
  activeTabId.value = 'welcome'
}

// Add Service Account State
const showAddSaModal = ref(false)
const newSaJson = ref('')
const savingSa = ref(false)

// Custom notification modal state
const showNotificationModal = ref(false)
const notificationModalTitle = ref('')
const notificationModalMessage = ref('')
const notificationModalType = ref('error')

function showNotification(title, message, type = 'error') {
  notificationModalTitle.value = title
  let msg = ''
  if (message && typeof message === 'object') {
    msg = message.message || JSON.stringify(message, null, 2)
  } else if (message) {
    msg = String(message)
  } else {
    msg = 'Unknown error'
  }
  notificationModalMessage.value = msg
  notificationModalType.value = type
  showNotificationModal.value = true
}

function alert(message) {
  const msgStr = String(message)
  const isSuccess = msgStr.toLowerCase().includes('success') || 
                    msgStr.toLowerCase().includes('copied') || 
                    msgStr.includes('Successfully')
  showNotification(isSuccess ? 'Success' : 'Error', message, isSuccess ? 'success' : 'error')
}

// Auth Modal States
const showClaimsModal = ref(false)
const showCreateUserModal = ref(false)
const showEditUserModal = ref(false)
const showExportAuthModal = ref(false)
const exportAuthLoading = ref(false)
const exportAuthData = ref(null)
const exportAuthHashParams = reactive({
  signerKey: '',
  saltSeparator: 'Bw==',
  rounds: 8,
  memoryCost: 14
})

const showImportAuthModal = ref(false)
const importAuthLoading = ref(false)
const importAuthJsonInput = ref('')
const importAuthHashParams = reactive({
  hashAlgorithm: 'SCRYPT',
  signerKey: '',
  saltSeparator: 'Bw==',
  rounds: 8,
  memoryCost: 14
})
const importAuthResult = ref(null)
const parsedImportSummary = ref(null)

const activeUser = ref(null)
const activeTabContext = ref(null)
const claimsJson = ref('{}')
const newUser = reactive({
  email: '',
  displayName: '',
  password: '',
  emailVerified: true,
  disabled: false
})
const editUser = reactive({
  uid: '',
  email: '',
  displayName: '',
  password: '',
  emailVerified: true,
  disabled: false
})

// Firestore Modal States
const showDocModal = ref(false)
const showAddDocModal = ref(false)
const activeDocPath = ref('')
const activeDocJson = ref('{}')
const newDocId = ref('')
const newDocFieldsJson = ref('{}')

// Copy Modal State
const showCopyModal = ref(false)
const activeCopySource = ref('')
const copyingNode = ref(false)
const destProjects = ref([])
const copyWizard = reactive({
  path: '',
  type: 'collection',
  destEmail: '',
  destProjectId: '',
  overwrite: false,
  recursive: true
})

onMounted(() => {
  loadPinnedProjects()
  refreshAccounts()
})

// Load accounts
async function refreshAccounts() {
  scanning.value = true
  try {
    const res = await GetAccounts()
    accounts.value = res || []
  } catch (err) {
    alert('Failed to scan credentials: ' + err)
  } finally {
    scanning.value = false
  }
}

async function triggerGoogleLogin() {
  loggingInGoogle.value = true
  showNotification('Google Login', 'Opening browser to complete Google Sign In...', 'info')
  try {
    const acc = await LoginGoogleAccount()
    if (acc && acc.email) {
      showNotification('Success', `Successfully signed in as ${acc.email}!`, 'success')
      await refreshAccounts()
      toggleAccount(acc.email)
    }
  } catch (err) {
    showNotification('Login Failed', 'Google authentication failed: ' + err, 'error')
  } finally {
    loggingInGoogle.value = false
  }
}

// Toggle account view
async function toggleAccount(email) {
  if (expandedAccount.value === email) {
    expandedAccount.value = null
    projects.value = []
  } else {
    expandedAccount.value = email
    loadingProjects.value = true
    projects.value = []
    try {
      const res = await ListProjects(email)
      projects.value = res || []
    } catch (err) {
      console.error('Failed to load projects:', err)
    } finally {
      loadingProjects.value = false
    }
  }
}

// Add Service Account
async function saveServiceAccount() {
  if (!newSaJson.value.trim()) return
  savingSa.value = true
  try {
    await AddServiceAccount(newSaJson.value.trim())
    showAddSaModal.value = false
    newSaJson.value = ''
    refreshAccounts()
    showToast('Service account added successfully!')
  } catch (err) {
    alert('Failed to add service account: ' + err)
  } finally {
    savingSa.value = false
  }
}

function getTabTitle(tab) {
  if (tab.type === 'auth') {
    return `Auth: ${tab.projectId}`
  }
  if (tab.type === 'firestore') {
    if (tab.activeDocPath) {
      const parts = tab.activeDocPath.split('/')
      return `Doc: ${parts[parts.length - 1]}`
    }
    if (tab.activeCollection) {
      return `Col: ${tab.activeCollection}`
    }
    return `Firestore: ${tab.projectId}`
  }
  return tab.type
}

function getTabFullTitle(tab) {
  if (tab.type === 'auth') {
    return `Auth: ${tab.projectId} (${tab.email})`
  }
  if (tab.type === 'firestore') {
    if (tab.activeDocPath) {
      return `Document: ${tab.activeDocPath} (${tab.projectId})`
    }
    if (tab.activeCollection) {
      return `Collection: ${tab.activeCollection} (${tab.projectId})`
    }
    return `Firestore: ${tab.projectId} (${tab.email})`
  }
  return tab.type
}

function openCollectionInNewTab(parentTab, colName) {
  const tabId = `firestore-col-${parentTab.email}-${parentTab.projectId}-${colName.replace(/\//g, '-')}`
  const existingTab = openTabs.value.find(t => t.id === tabId)
  
  if (existingTab) {
    activeTabId.value = tabId
  } else {
    const newTab = {
      id: tabId,
      type: 'firestore',
      email: parentTab.email,
      projectId: parentTab.projectId,
      searchQuery: '',
      loading: false,
      users: [],
      nextPageToken: '',
      collections: parentTab.collections || [],
      collectionsCollapsed: false,
      collectionSearch: '',
      activeCollection: colName,
      activeDocPath: '',
      documents: [],
      nextDocPageToken: '',
      loadingCollections: false,
      loadingDocuments: false,
      firestoreView: 'explorer',
      resultsView: parentTab.resultsView || 'table',
      selectedDocNames: [],
      queryExpanded: true,
      queryMode: 'simple',
      queryDocId: '',
      queryFilterField: '',
      queryFilterOp: '==',
      queryFilterVal: '',
      querySortField: '',
      querySortDir: 'asc',
      queryLimit: 50,
      queryPreset: 'Basic Fetch',
      queryScript: `// Query with JavaScript using the Firebase Admin SDK\nasync function run() {\n  const query = await db.collection("${colName}")\n    .limit(50)\n    .get();\n  return query;\n}`,
      queryOutput: '',
      runningQuery: false
    }
    openTabs.value.push(newTab)
    const reactiveTab = openTabs.value[openTabs.value.length - 1]
    activeTabId.value = tabId
    fetchDocuments(reactiveTab)
  }
}

function openCollectionTab(parentTab, colName) {
  openCollectionInNewTab(parentTab, colName)
}

// Tabs Operations
function openTab(type, email, projectId) {
  const tabId = `${type}-${email}-${projectId}`
  const existingTab = openTabs.value.find(t => t.id === tabId)
  
  if (existingTab) {
    activeTabId.value = tabId
  } else {
    const newTab = {
      id: tabId,
      type,
      email,
      projectId,
      searchQuery: '',
      statusFilter: 'all',
      loading: false,
      users: [],
      nextPageToken: '',
      collections: [],
      collectionsCollapsed: false,
      collectionSearch: '',
      activeCollection: '',
      activeDocPath: '',
      documents: [],
      nextDocPageToken: '',
      loadingCollections: false,
      loadingDocuments: false,
      firestoreView: 'explorer',
      resultsView: 'table',
      selectedDocNames: [],
      queryExpanded: true,
      queryMode: 'simple',
      queryDocId: '',
      queryFilterField: '',
      queryFilterOp: '==',
      queryFilterVal: '',
      querySortField: '',
      querySortDir: 'asc',
      queryLimit: 50,
      queryPreset: 'Basic Fetch',
      queryScript: '',
      queryOutput: '',
      runningQuery: false,
      authSortField: 'email',
      authSortDir: 'asc'
    }
    openTabs.value.push(newTab)
    const reactiveTab = openTabs.value[openTabs.value.length - 1]
    activeTabId.value = tabId
    
    // Initial fetch for tab
    if (type === 'auth') {
      fetchAuthUsers(reactiveTab)
    } else if (type === 'firestore') {
      refreshCollections(reactiveTab)
    }
  }
}

function closeTab(tabId) {
  openTabs.value = openTabs.value.filter(t => t.id !== tabId)
  if (activeTabId.value === tabId) {
    activeTabId.value = openTabs.value.length > 0 ? openTabs.value[openTabs.value.length - 1].id : 'welcome'
  }
}

function getUserInitial(user) {
  if (user.displayName) return user.displayName.charAt(0).toUpperCase()
  if (user.email) return user.email.charAt(0).toUpperCase()
  return 'U'
}

function hasCustomClaims(user) {
  return (user.customClaims && Object.keys(user.customClaims).length > 0) || (user.customAttributes && user.customAttributes !== '{}')
}

function getFilteredSortedUsers(tab) {
  if (!tab || !tab.users) return []
  
  let result = [...tab.users]
  
  // Status filter pills
  if (tab.statusFilter === 'active') {
    result = result.filter(u => !u.disabled)
  } else if (tab.statusFilter === 'disabled') {
    result = result.filter(u => u.disabled)
  } else if (tab.statusFilter === 'claims') {
    result = result.filter(u => hasCustomClaims(u))
  }

  // Unified instant search
  if (tab.searchQuery && tab.searchQuery.trim()) {
    const q = tab.searchQuery.trim().toLowerCase()
    result = result.filter(u => 
      (u.uid && u.uid.toLowerCase().includes(q)) ||
      (u.email && u.email.toLowerCase().includes(q)) ||
      (u.displayName && u.displayName.toLowerCase().includes(q))
    )
  }
  
  // Sorting
  const field = tab.authSortField || 'email'
  const dir = tab.authSortDir || 'asc'
  const isAsc = dir === 'asc'
  
  result.sort((a, b) => {
    let valA = a[field]
    let valB = b[field]
    
    if (field === 'createdAt' || field === 'lastLoginAt') {
      const parsedA = parseInt(valA)
      valA = !isNaN(parsedA) ? parsedA : (valA ? new Date(valA).getTime() : 0)
      const parsedB = parseInt(valB)
      valB = !isNaN(parsedB) ? parsedB : (valB ? new Date(valB).getTime() : 0)
    } else if (field === 'disabled') {
      valA = valA ? 1 : 0
      valB = valB ? 1 : 0
    } else {
      valA = valA ? String(valA).toLowerCase() : ''
      valB = valB ? String(valB).toLowerCase() : ''
    }
    
    if (valA < valB) return isAsc ? -1 : 1
    if (valA > valB) return isAsc ? 1 : -1
    return 0
  })
  
  return result
}

function toggleAuthSort(tab, field) {
  if (tab.authSortField === field) {
    tab.authSortDir = tab.authSortDir === 'asc' ? 'desc' : 'asc'
  } else {
    tab.authSortField = field
    tab.authSortDir = 'asc'
  }
}

// Auth Actions
async function fetchAuthUsers(tab) {
  tab.loading = true
  tab.users = []
  tab.nextPageToken = ''
  try {
    const res = await ListAuthUsers(tab.email, tab.projectId, 50, '', tab.searchQuery || '')
    if (res) {
      tab.users = res.users || []
      tab.nextPageToken = res.nextPageToken || ''
    }
  } catch (err) {
    alert('Failed to load users: ' + err)
  } finally {
    tab.loading = false
  }
}

async function fetchAuthUsersNext(tab) {
  if (!tab.nextPageToken) return
  tab.loading = true
  try {
    const res = await ListAuthUsers(tab.email, tab.projectId, 50, tab.nextPageToken, tab.searchQuery || '')
    if (res) {
      tab.users.push(...(res.users || []))
      tab.nextPageToken = res.nextPageToken || ''
    }
  } catch (err) {
    alert('Failed to load more users: ' + err)
  } finally {
    tab.loading = false
  }
}

function countUsersWithClaims(users) {
  if (!users) return 0
  return users.filter(u => hasCustomClaims(u)).length
}

async function openExportAuthModal(tab) {
  activeTabContext.value = tab
  showExportAuthModal.value = true
  exportAuthLoading.value = true
  exportAuthData.value = null
  exportAuthHashParams.signerKey = ''
  exportAuthHashParams.saltSeparator = 'Bw=='
  exportAuthHashParams.rounds = 8
  exportAuthHashParams.memoryCost = 14
  try {
    const res = await ExportAuthUsers(tab.email, tab.projectId)
    exportAuthData.value = res
  } catch (err) {
    showNotification('Error', 'Failed to export auth users: ' + err, 'error')
    showExportAuthModal.value = false
  } finally {
    exportAuthLoading.value = false
  }
}

function copyExportAuthJson() {
  if (!exportAuthData.value) return
  const fullExport = {
    users: exportAuthData.value.users,
    hashParameters: {
      hashAlgorithm: 'SCRYPT',
      signerKey: exportAuthHashParams.signerKey,
      saltSeparator: exportAuthHashParams.saltSeparator || 'Bw==',
      rounds: exportAuthHashParams.rounds || 8,
      memoryCost: exportAuthHashParams.memoryCost || 14
    }
  }
  const jsonStr = JSON.stringify(fullExport, null, 2)
  navigator.clipboard.writeText(jsonStr)
  showToast('Users & hash parameters copied to clipboard!')
}

function downloadExportAuthJson() {
  if (!exportAuthData.value || !activeTabContext.value) return
  const fullExport = {
    users: exportAuthData.value.users,
    hashParameters: {
      hashAlgorithm: 'SCRYPT',
      signerKey: exportAuthHashParams.signerKey,
      saltSeparator: exportAuthHashParams.saltSeparator || 'Bw==',
      rounds: exportAuthHashParams.rounds || 8,
      memoryCost: exportAuthHashParams.memoryCost || 14
    }
  }
  const jsonStr = JSON.stringify(fullExport, null, 2)
  const blob = new Blob([jsonStr], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `auth-users-${activeTabContext.value.projectId}.json`
  a.click()
  URL.revokeObjectURL(url)
}

function openImportAuthModal(tab) {
  activeTabContext.value = tab
  showImportAuthModal.value = true
  importAuthJsonInput.value = ''
  importAuthResult.value = null
  parsedImportSummary.value = null
  importAuthHashParams.hashAlgorithm = 'SCRYPT'
  importAuthHashParams.signerKey = ''
  importAuthHashParams.saltSeparator = 'Bw=='
  importAuthHashParams.rounds = 8
  importAuthHashParams.memoryCost = 14
}

function triggerImportFilePicker() {
  const el = document.getElementById('auth-import-file-input')
  if (el) el.click()
}

function onImportFileSelected(event) {
  const file = event.target.files?.[0]
  if (!file) return
  const reader = new FileReader()
  reader.onload = (e) => {
    importAuthJsonInput.value = e.target.result
    parseImportJsonInput()
  }
  reader.readAsText(file)
}

function parseImportJsonInput() {
  parsedImportSummary.value = null
  if (!importAuthJsonInput.value.trim()) return

  try {
    const raw = JSON.parse(importAuthJsonInput.value)
    let userArray = []
    let embeddedParams = null

    if (Array.isArray(raw)) {
      userArray = raw
    } else if (raw && typeof raw === 'object') {
      if (Array.isArray(raw.users)) userArray = raw.users
      if (raw.hashParameters) embeddedParams = raw.hashParameters
    }

    if (embeddedParams) {
      if (embeddedParams.hashAlgorithm) importAuthHashParams.hashAlgorithm = embeddedParams.hashAlgorithm
      if (embeddedParams.signerKey) importAuthHashParams.signerKey = embeddedParams.signerKey
      if (embeddedParams.saltSeparator) importAuthHashParams.saltSeparator = embeddedParams.saltSeparator
      if (embeddedParams.rounds) importAuthHashParams.rounds = embeddedParams.rounds
      if (embeddedParams.memoryCost) importAuthHashParams.memoryCost = embeddedParams.memoryCost
    }

    const userCount = userArray.length
    const claimsCount = userArray.filter(u => hasCustomClaims(u)).length
    const hashCount = userArray.filter(u => u.passwordHash).length

    parsedImportSummary.value = {
      userCount,
      claimsCount,
      hashCount,
      hasEmbeddedParams: !!embeddedParams,
      rawUsers: userArray
    }
  } catch (err) {
    parsedImportSummary.value = null
  }
}

async function runImportAuthUsers(tab) {
  if (!parsedImportSummary.value || parsedImportSummary.value.userCount === 0) return

  importAuthLoading.value = true
  importAuthResult.value = null

  const req = {
    users: parsedImportSummary.value.rawUsers,
    hashAlgorithm: importAuthHashParams.hashAlgorithm || 'SCRYPT',
    signerKey: importAuthHashParams.signerKey || '',
    saltSeparator: importAuthHashParams.saltSeparator || 'Bw==',
    rounds: Number(importAuthHashParams.rounds) || 8,
    memoryCost: Number(importAuthHashParams.memoryCost) || 14
  }

  try {
    const res = await ImportAuthUsers(tab.email, tab.projectId, req)
    importAuthResult.value = res
    if (res.success || res.importedCount > 0) {
      showToast(`Imported ${res.importedCount} user(s) successfully!`)
      fetchAuthUsers(tab)
    } else {
      showNotification('Import Failed', `Failed to import users. Check details in modal.`, 'error')
    }
  } catch (err) {
    showNotification('Error', 'Import failed: ' + err, 'error')
  } finally {
    importAuthLoading.value = false
  }
}

function openCreateUserModal(tab) {
  activeTabContext.value = tab
  newUser.email = ''
  newUser.displayName = ''
  newUser.password = ''
  newUser.emailVerified = true
  newUser.disabled = false
  showCreateUserModal.value = true
}

async function saveNewUser() {
  const tab = activeTabContext.value
  if (!tab) return
  try {
    const req = {
      email: newUser.email,
      displayName: newUser.displayName,
      password: newUser.password,
      emailVerified: newUser.emailVerified,
      disabled: newUser.disabled
    }
    await CreateAuthUser(tab.email, tab.projectId, req)
    showCreateUserModal.value = false
    showToast('User created successfully!')
    fetchAuthUsers(tab)
  } catch (err) {
    alert('Failed to create user: ' + err)
  }
}

function editUserBasic(tab, user) {
  activeTabContext.value = tab
  activeUser.value = user
  editUser.uid = user.uid
  editUser.email = user.email
  editUser.displayName = user.displayName
  editUser.password = ''
  editUser.emailVerified = user.emailVerified
  editUser.disabled = user.disabled
  showEditUserModal.value = true
}

async function saveUserEdit() {
  const tab = activeTabContext.value
  if (!tab) return
  try {
    const req = {
      localId: editUser.uid,
      email: editUser.email,
      displayName: editUser.displayName,
      password: editUser.password || undefined,
      emailVerified: editUser.emailVerified,
      disabled: editUser.disabled
    }
    await UpdateAuthUser(tab.email, tab.projectId, req)
    showEditUserModal.value = false
    showToast('User updated successfully!')
    fetchAuthUsers(tab)
  } catch (err) {
    alert('Failed to update user: ' + err)
  }
}

function editUserClaims(tab, user) {
  activeTabContext.value = tab
  activeUser.value = user
  claimsJson.value = JSON.stringify(user.customClaims || {}, null, 2)
  showClaimsModal.value = true
}

async function saveCustomClaims() {
  const tab = activeTabContext.value
  const user = activeUser.value
  if (!tab || !user) return
  
  let claims = {}
  if (claimsJson.value.trim()) {
    try {
      claims = JSON.parse(claimsJson.value.trim())
    } catch (err) {
      alert('Invalid JSON format for claims')
      return
    }
  }

  try {
    const req = {
      localId: user.uid,
      customAttributes: JSON.stringify(claims)
    }
    await UpdateAuthUser(tab.email, tab.projectId, req)
    showClaimsModal.value = false
    showToast('Claims saved successfully!')
    fetchAuthUsers(tab)
  } catch (err) {
    alert('Failed to save claims: ' + err)
  }
}

async function deleteUserConfirm(tab, user) {
  if (confirm(`Are you sure you want to delete user ${user.email || user.uid}?`)) {
    try {
      await DeleteAuthUser(tab.email, tab.projectId, user.uid)
      showToast('User deleted!')
      fetchAuthUsers(tab)
    } catch (err) {
      alert('Failed to delete user: ' + err)
    }
  }
}

// Firestore Actions
function getFilteredCollections(tab) {
  if (!tab.collections) return []
  if (!tab.collectionSearch || !tab.collectionSearch.trim()) return tab.collections
  const q = tab.collectionSearch.trim().toLowerCase()
  return tab.collections.filter(c => c.toLowerCase().includes(q))
}

async function refreshCollections(tab) {
  tab.loadingCollections = true
  tab.collections = []
  try {
    const res = await ListRootCollections(tab.email, tab.projectId, '')
    tab.collections = res || []
  } catch (err) {
    console.error('Failed to list collections:', err)
  } finally {
    tab.loadingCollections = false
  }
}

function selectCollection(tab, colName) {
  openCollectionInNewTab(tab, colName)
}

function setQueryMode(tab, mode) {
  tab.queryMode = mode
}

function mapSimpleOpToApi(op) {
  const mapping = {
    '==': 'EQUAL',
    '>': 'GREATER_THAN',
    '>=': 'GREATER_THAN_OR_EQUAL',
    '<': 'LESS_THAN',
    '<=': 'LESS_THAN_OR_EQUAL',
    '!=': 'NOT_EQUAL',
    'array-contains': 'ARRAY_CONTAINS'
  }
  return mapping[op] || 'EQUAL'
}

async function runSimpleQueryAction(tab) {
  if (!tab.activeCollection) return
  tab.loadingDocuments = true
  tab.selectedDocNames = []
  tab.documents = []
  
  try {
    if (tab.queryDocId && tab.queryDocId.trim()) {
      const docPath = tab.activeCollection + '/' + tab.queryDocId.trim()
      const res = await GetDocument(tab.email, tab.projectId, '', docPath)
      if (res) {
        tab.documents = [res]
      }
    } else {
      const structuredQuery = {
        from: [{ collectionId: tab.activeCollection }]
      }
      
      if (tab.queryFilterField && tab.queryFilterVal) {
        let valObj = {}
        const rawVal = tab.queryFilterVal.trim()
        if (rawVal === 'true' || rawVal === 'false') {
          valObj = { booleanValue: rawVal === 'true' }
        } else if (!isNaN(rawVal)) {
          if (rawVal.includes('.')) {
            valObj = { doubleValue: parseFloat(rawVal) }
          } else {
            valObj = { integerValue: rawVal }
          }
        } else {
          valObj = { stringValue: rawVal }
        }
        
        structuredQuery.where = {
          fieldFilter: {
            field: { fieldPath: tab.queryFilterField },
            op: mapSimpleOpToApi(tab.queryFilterOp),
            value: valObj
          }
        }
      }
      
      if (tab.querySortField) {
        structuredQuery.orderBy = [{
          field: { fieldPath: tab.querySortField },
          direction: tab.querySortDir === 'asc' ? 'ASCENDING' : 'DESCENDING'
        }]
      }
      
      if (tab.queryLimit && tab.queryLimit > 0) {
        structuredQuery.limit = parseInt(tab.queryLimit)
      }
      
      const res = await RunSimpleQuery(tab.email, tab.projectId, '', { structuredQuery })
      const docs = []
      if (res && Array.isArray(res)) {
        res.forEach(r => {
          if (r.document) docs.push(r.document)
        })
      }
      tab.documents = docs
      showToast(`Loaded ${docs.length} document(s)`)
    }
  } catch (err) {
    alert('Failed to run query: ' + err)
  } finally {
    tab.loadingDocuments = false
  }
}

function loadPresetQueryScript(tab) {
  const col = tab.activeCollection || 'users'
  if (tab.queryPreset === 'Basic Fetch') {
    tab.queryScript = `// Query with JavaScript using the Firebase Admin SDK\nasync function run() {\n  const query = await db.collection("${col}")\n    .limit(50)\n    .get();\n  return query;\n}`
  } else if (tab.queryPreset === 'Filter & Order') {
    tab.queryScript = `// Query with filters and sorting\nasync function run() {\n  const query = await db.collection("${col}")\n    .where("status", "==", "active")\n    .orderBy("createdAt", "desc")\n    .limit(20)\n    .get();\n  return query;\n}`
  } else if (tab.queryPreset === 'Projection') {
    tab.queryScript = `// Map documents to simple projections\nasync function run() {\n  const snapshot = await db.collection("${col}").limit(10).get();\n  const results = [];\n  snapshot.forEach(doc => {\n    const data = doc.data();\n    results.push({\n      id: doc.id,\n      name: data.name || data.displayName || 'Unnamed'\n    });\n  });\n  return results;\n}`
  } else if (tab.queryPreset === 'Aggregate Sum') {
    tab.queryScript = `// Calculate sum or averages\nasync function run() {\n  const snapshot = await db.collection("${col}").get();\n  let total = 0;\n  snapshot.forEach(doc => {\n    const data = doc.data();\n    total += Number(data.likes || data.amount || 0);\n  });\n  return { count: snapshot.size, totalSum: total };\n}`
  } else if (tab.queryPreset === 'Subcollection Fetch') {
    tab.queryScript = `// Query a subcollection from a specific document\nasync function run() {\n  const query = await db.collection("${col}").doc("some-doc-id").collection("subcollection")\n    .limit(50)\n    .get();\n  return query;\n}`
  }
}

function openDocInNewTab(parentTab, doc) {
  const docId = getDocId(doc.name)
  const docPath = parentTab.activeCollection + '/' + docId
  const tabId = `firestore-doc-${parentTab.email}-${parentTab.projectId}-${docPath.replace(/\//g, '-')}`
  const existingTab = openTabs.value.find(t => t.id === tabId)
  
  if (existingTab) {
    activeTabId.value = tabId
  } else {
    const newTab = {
      id: tabId,
      type: 'firestore',
      email: parentTab.email,
      projectId: parentTab.projectId,
      searchQuery: '',
      loading: false,
      users: [],
      nextPageToken: '',
      collections: parentTab.collections || [],
      collectionsCollapsed: false,
      activeCollection: parentTab.activeCollection,
      activeDocPath: docPath,
      documents: [doc],
      nextDocPageToken: '',
      loadingCollections: false,
      loadingDocuments: false,
      firestoreView: 'explorer',
      resultsView: 'tree',
      selectedDocNames: [],
      queryExpanded: false,
      queryMode: 'js',
      queryScript: `// Document query\nasync function run() {\n  const doc = await db.doc("${docPath}").get();\n  return doc;\n}`,
      queryOutput: '',
      runningQuery: false
    }
    openTabs.value.push(newTab)
    activeTabId.value = tabId
  }
}

// Firestore Document Selections
function toggleSelectDoc(tab, docName) {
  if (!tab.selectedDocNames) tab.selectedDocNames = []
  const index = tab.selectedDocNames.indexOf(docName)
  if (index === -1) {
    tab.selectedDocNames.push(docName)
  } else {
    tab.selectedDocNames.splice(index, 1)
  }
}

function toggleSelectAllDocs(tab) {
  if (!tab.selectedDocNames) tab.selectedDocNames = []
  if (tab.selectedDocNames.length === tab.documents.length) {
    tab.selectedDocNames = []
  } else {
    tab.selectedDocNames = tab.documents.map(d => d.name)
  }
}

async function deleteSelectedDocuments(tab) {
  const count = tab.selectedDocNames?.length || 0
  if (count === 0) return
  if (confirm(`Are you sure you want to delete the ${count} selected documents?`)) {
    tab.loadingDocuments = true
    try {
      for (const docName of tab.selectedDocNames) {
        const prefix = '/documents/'
        const docIdx = docName.indexOf(prefix)
        const relativePath = docIdx !== -1 ? docName.substring(docIdx + prefix.length) : docName
        await DeleteDocument(tab.email, tab.projectId, '', relativePath, false)
      }
      showToast(`Deleted ${count} documents!`)
      tab.selectedDocNames = []
      fetchDocuments(tab)
    } catch (err) {
      alert('Failed to delete some documents: ' + err)
    } finally {
      tab.loadingDocuments = false
    }
  }
}

// Export Results Modal & Action States
const showExportModal = ref(false)
const exportPreviewText = ref('')
const exportConfig = reactive({
  format: 'JSON',
  includeIds: true,
  includePaths: false,
  onlySelected: false,
  hasSelection: false,
  selectedCount: 0,
  activeTab: null
})

function triggerExportAll(tab) {
  exportConfig.activeTab = tab
  exportConfig.hasSelection = (tab.selectedDocNames?.length || 0) > 0
  exportConfig.selectedCount = tab.selectedDocNames?.length || 0
  exportConfig.onlySelected = false
  exportConfig.format = 'JSON'
  exportConfig.includeIds = true
  exportConfig.includePaths = false
  
  updateExportPreview()
  showExportModal.value = true
}

function triggerExportSelected(tab) {
  exportConfig.activeTab = tab
  exportConfig.hasSelection = true
  exportConfig.selectedCount = tab.selectedDocNames?.length || 0
  exportConfig.onlySelected = true
  exportConfig.format = 'JSON'
  exportConfig.includeIds = true
  exportConfig.includePaths = false
  
  updateExportPreview()
  showExportModal.value = true
}

function generateExportString() {
  const tab = exportConfig.activeTab
  if (!tab || !tab.documents) return ''
  
  let docsToExport = tab.documents
  if (exportConfig.onlySelected && tab.selectedDocNames) {
    docsToExport = docsToExport.filter(d => tab.selectedDocNames.includes(d.name))
  }

  if (exportConfig.format === 'JSON') {
    const list = docsToExport.map(d => {
      const docObj = {}
      if (exportConfig.includeIds) docObj.id = getDocId(d.name)
      if (exportConfig.includePaths) docObj.path = d.name
      docObj.fields = cleanDocFields(d.fields)
      return docObj
    })
    return JSON.stringify(list, null, 2)
  } else {
    const columns = getTableColumns(docsToExport)
    let csv = ''
    
    const headers = []
    if (exportConfig.includeIds) headers.push('id')
    if (exportConfig.includePaths) headers.push('path')
    headers.push(...columns)
    csv += headers.map(h => escapeCSV(h)).join(',') + '\n'

    docsToExport.forEach(d => {
      const row = []
      if (exportConfig.includeIds) row.push(getDocId(d.name))
      if (exportConfig.includePaths) row.push(d.name)
      
      const cleaned = cleanDocFields(d.fields)
      columns.forEach(col => {
        const val = cleaned[col]
        if (val === undefined) row.push('')
        else if (val === null) row.push('null')
        else if (typeof val === 'object') {
          if (val.__time__) row.push(val.__time__)
          else row.push(JSON.stringify(val))
        } else {
          row.push(String(val))
        }
      })
      csv += row.map(r => escapeCSV(r)).join(',') + '\n'
    })
    return csv
  }
}

function escapeCSV(val) {
  const str = String(val)
  if (str.includes(',') || str.includes('"') || str.includes('\n')) {
    return `"${str.replace(/"/g, '""')}"`
  }
  return str
}

function updateExportPreview() {
  const text = generateExportString()
  if (text.length > 5000) {
    exportPreviewText.value = text.substring(0, 5000) + '\n\n... [TRUNCATED PREVIEW] ...'
  } else {
    exportPreviewText.value = text
  }
}

async function copyExportToClipboard() {
  const text = generateExportString()
  try {
    await navigator.clipboard.writeText(text)
    showToast('Copied export data to clipboard!')
  } catch (err) {
    alert('Failed to copy to clipboard: ' + err)
  }
}

function downloadExportFile() {
  const text = generateExportString()
  const tab = exportConfig.activeTab
  const colName = tab ? tab.activeCollection : 'export'
  const ext = exportConfig.format.toLowerCase()
  
  const blob = new Blob([text], { type: exportConfig.format === 'JSON' ? 'application/json' : 'text/csv' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${colName}_export.${ext}`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
  
  showExportModal.value = false
  showToast('Download started!')
}

// Convert typed Firestore fields to clean JS Object recursively
function cleanDocFields(fields) {
  const res = {}
  if (!fields) return res
  for (const key in fields) {
    res[key] = cleanDocValue(fields[key])
  }
  return res
}

function cleanDocValue(valObj) {
  if (!valObj) return null
  if ('nullValue' in valObj) return null
  if ('booleanValue' in valObj) return valObj.booleanValue
  if ('stringValue' in valObj) return valObj.stringValue
  if ('integerValue' in valObj) return parseInt(valObj.integerValue)
  if ('doubleValue' in valObj) return parseFloat(valObj.doubleValue)
  if ('timestampValue' in valObj) return { __time__: valObj.timestampValue }
  if ('referenceValue' in valObj) return valObj.referenceValue
  if ('mapValue' in valObj) {
    return cleanDocFields(valObj.mapValue.fields || {})
  }
  if ('arrayValue' in valObj) {
    const values = valObj.arrayValue.values || []
    return values.map(v => cleanDocValue(v))
  }
  return null
}

function getTableColumns(docs) {
  if (!docs || docs.length === 0) return []
  const columns = new Set()
  docs.forEach(doc => {
    const cleaned = cleanDocFields(doc.fields)
    Object.keys(cleaned).forEach(k => columns.add(k))
  })
  return Array.from(columns).slice(0, 8)
}

function getAvailableFields(tab) {
  if (!tab || !tab.documents || tab.documents.length === 0) return []
  const columns = new Set()
  tab.documents.forEach(doc => {
    if (doc.fields) {
      const cleaned = cleanDocFields(doc.fields)
      Object.keys(cleaned).forEach(k => columns.add(k))
    }
  })
  return Array.from(columns).sort()
}

function getColumnValue(doc, col) {
  const cleaned = cleanDocFields(doc.fields)
  const val = cleaned[col]
  if (val === undefined) return '-'
  if (val === null) return 'null'
  if (typeof val === 'object') {
    if (val.__time__) return new Date(val.__time__).toLocaleString()
    return JSON.stringify(val)
  }
  return String(val)
}

function formatAllDocsJson(docs) {
  if (!docs || docs.length === 0) return '[]'
  const list = docs.map(d => {
    return {
      id: getDocId(d.name),
      path: d.name,
      fields: cleanDocFields(d.fields)
    }
  })
  return JSON.stringify(list, null, 2)
}

function convertJsonToFirestoreFields(obj) {
  const fields = {}
  for (const key in obj) {
    fields[key] = convertValueToFirestore(obj[key])
  }
  return fields
}

function convertValueToFirestore(val) {
  if (val === null) return { nullValue: null }
  if (typeof val === 'boolean') return { booleanValue: val }
  if (typeof val === 'string') return { stringValue: val }
  if (typeof val === 'number') {
    if (Number.isInteger(val)) {
      return { integerValue: val.toString() }
    }
    return { doubleValue: val }
  }
  if (Array.isArray(val)) {
    return {
      arrayValue: {
        values: val.map(item => convertValueToFirestore(item))
      }
    }
  }
  if (typeof val === 'object') {
    if (val.__time__) {
      return { timestampValue: val.__time__ }
    }
    const fields = {}
    for (const k in val) {
      fields[k] = convertValueToFirestore(val[k])
    }
    return { mapValue: { fields } }
  }
  return { stringValue: String(val) }
}

async function saveDocFieldChange(tab, doc, fieldPath, newVal) {
  const cleanFields = cleanDocFields(doc.fields)
  const pathParts = fieldPath.split('/')
  let cur = cleanFields
  for (let i = 0; i < pathParts.length - 1; i++) {
    const part = pathParts[i]
    if (cur[part] === undefined) cur[part] = {}
    cur = cur[part]
  }
  const lastPart = pathParts[pathParts.length - 1]
  cur[lastPart] = newVal

  const updatedTypedFields = convertJsonToFirestoreFields(cleanFields)
  
  try {
    const prefix = '/documents/'
    const docIdx = doc.name.indexOf(prefix)
    const relativePath = docIdx !== -1 ? doc.name.substring(docIdx + prefix.length) : doc.name

    await SaveDocument(tab.email, tab.projectId, '', relativePath, updatedTypedFields, false)
    doc.fields = updatedTypedFields
    showToast('Field saved!')
  } catch (err) {
    alert('Failed to save field change: ' + err)
  }
}

async function addDocField(tab, doc, parentPath, fieldName, fieldType, initialVal) {
  const cleanFields = cleanDocFields(doc.fields)
  const fullPath = parentPath ? `${parentPath}/${fieldName}` : fieldName
  const pathParts = fullPath.split('/')
  let cur = cleanFields
  for (let i = 0; i < pathParts.length - 1; i++) {
    const part = pathParts[i]
    if (cur[part] === undefined) cur[part] = {}
    cur = cur[part]
  }
  const lastPart = pathParts[pathParts.length - 1]
  cur[lastPart] = initialVal

  const updatedTypedFields = convertJsonToFirestoreFields(cleanFields)
  
  try {
    const prefix = '/documents/'
    const docIdx = doc.name.indexOf(prefix)
    const relativePath = docIdx !== -1 ? doc.name.substring(docIdx + prefix.length) : doc.name

    await SaveDocument(tab.email, tab.projectId, '', relativePath, updatedTypedFields, false)
    doc.fields = updatedTypedFields
    showToast(`Added field "${fieldName}"!`)
  } catch (err) {
    alert('Failed to add field: ' + err)
  }
}

async function deleteDocField(tab, doc, fieldPath) {
  if (!confirm(`Are you sure you want to delete field "${fieldPath}"?`)) return
  const cleanFields = cleanDocFields(doc.fields)
  const pathParts = fieldPath.split('/')
  let cur = cleanFields
  for (let i = 0; i < pathParts.length - 1; i++) {
    const part = pathParts[i]
    if (cur[part] === undefined) return
    cur = cur[part]
  }
  const lastPart = pathParts[pathParts.length - 1]
  if (Array.isArray(cur)) {
    const idx = parseInt(lastPart)
    if (!isNaN(idx)) cur.splice(idx, 1)
  } else {
    delete cur[lastPart]
  }

  const updatedTypedFields = convertJsonToFirestoreFields(cleanFields)
  
  try {
    const prefix = '/documents/'
    const docIdx = doc.name.indexOf(prefix)
    const relativePath = docIdx !== -1 ? doc.name.substring(docIdx + prefix.length) : doc.name

    await SaveDocument(tab.email, tab.projectId, '', relativePath, updatedTypedFields, false)
    doc.fields = updatedTypedFields
    showToast(`Deleted field "${fieldPath}"!`)
  } catch (err) {
    alert('Failed to delete field: ' + err)
  }
}

async function fetchDocuments(tab) {
  if (!tab.activeCollection) return
  tab.loadingDocuments = true
  tab.documents = []
  tab.nextDocPageToken = ''
  try {
    const res = await ListDocuments(tab.email, tab.projectId, '', tab.activeCollection, 50, '')
    if (res) {
      tab.documents = res.documents || []
      tab.nextDocPageToken = res.nextPageToken || ''
    }
  } catch (err) {
    alert('Failed to fetch documents: ' + err)
  } finally {
    tab.loadingDocuments = false
  }
}

async function fetchDocumentsNext(tab) {
  if (!tab.nextDocPageToken) return
  tab.loadingDocuments = true
  try {
    const res = await ListDocuments(tab.email, tab.projectId, '', tab.activeCollection, 50, tab.nextDocPageToken)
    if (res) {
      tab.documents.push(...(res.documents || []))
      tab.nextDocPageToken = res.nextPageToken || ''
    }
  } catch (err) {
    alert('Failed to fetch more documents: ' + err)
  } finally {
    tab.loadingDocuments = false
  }
}

function isTypedFirestoreField(val) {
  if (val && typeof val === 'object' && !Array.isArray(val)) {
    const keys = Object.keys(val)
    const firestoreTypes = ['stringValue', 'integerValue', 'doubleValue', 'booleanValue', 'mapValue', 'arrayValue', 'timestampValue', 'nullValue', 'bytesValue', 'referenceValue', 'geoPointValue']
    return keys.some(k => firestoreTypes.includes(k))
  }
  return false
}

function prepareFieldsForSaving(parsedJson) {
  if (!parsedJson || typeof parsedJson !== 'object') return {}
  const values = Object.values(parsedJson)
  if (values.length > 0 && values.every(isTypedFirestoreField)) {
    return parsedJson
  }
  return convertJsonToFirestoreFields(parsedJson)
}

function openAddDocModal(tab) {
  activeTabContext.value = tab
  newDocId.value = ''
  newDocFieldsJson.value = '{\n  \n}'
  showAddDocModal.value = true
}

async function saveNewDoc() {
  const tab = activeTabContext.value
  if (!tab) return
  
  let parsed = {}
  try {
    parsed = JSON.parse(newDocFieldsJson.value.trim() || '{}')
  } catch (e) {
    alert('Invalid JSON fields: ' + e)
    return
  }

  const fields = prepareFieldsForSaving(parsed)
  const path = tab.activeCollection + '/' + (newDocId.value.trim() || generateId())
  try {
    await SaveDocument(tab.email, tab.projectId, '', path, fields, true)
    showAddDocModal.value = false
    showToast('Document created successfully!')
    fetchDocuments(tab)
  } catch (err) {
    alert('Failed to save document: ' + err)
  }
}

function editDocumentJSON(tab, doc) {
  activeTabContext.value = tab
  const docId = getDocId(doc.name)
  activeDocPath.value = tab.activeCollection + '/' + docId
  activeDocJson.value = JSON.stringify(cleanDocFields(doc.fields), null, 2)
  showDocModal.value = true
}

function formatActiveDocJson() {
  try {
    const parsed = JSON.parse(activeDocJson.value)
    activeDocJson.value = JSON.stringify(parsed, null, 2)
    showToast('Formatted JSON!')
  } catch (e) {
    alert('Invalid JSON: ' + e)
  }
}

async function saveDocFromJSON() {
  const tab = activeTabContext.value
  if (!tab) return
  let parsed = {}
  try {
    parsed = JSON.parse(activeDocJson.value.trim() || '{}')
  } catch (err) {
    alert('Invalid JSON data format: ' + err)
    return
  }

  const fields = prepareFieldsForSaving(parsed)

  try {
    await SaveDocument(tab.email, tab.projectId, '', activeDocPath.value, fields, false)
    showDocModal.value = false
    showToast('Document saved!')
    const targetDoc = tab.documents ? tab.documents.find(d => getDocId(d.name) === getDocId(activeDocPath.value)) : null
    if (targetDoc) {
      targetDoc.fields = fields
    }
  } catch (err) {
    alert('Failed to save document JSON: ' + err)
  }
}

async function deleteDocumentConfirm(tab, doc) {
  const docId = getDocId(doc.name)
  const path = tab.activeCollection + '/' + docId
  const recursive = confirm(`Delete document ${docId}? Click OK to delete recursively (including subcollections) or Cancel to abort.`)
  if (recursive) {
    try {
      await DeleteDocument(tab.email, tab.projectId, '', path, true)
      showToast(`Document ${docId} deleted!`)
      fetchDocuments(tab)
    } catch (err) {
      alert('Failed to delete document: ' + err)
    }
  }
}

async function runQueryScript(tab) {
  if (!tab.queryScript.trim()) return
  tab.runningQuery = true
  tab.queryOutput = ''
  try {
    const res = await RunJSScript(tab.email, tab.projectId, '', tab.queryScript)
    if (res && res.logs) {
      tab.queryOutput = res.logs.join('\n')
    } else {
      tab.queryOutput = JSON.stringify(res, null, 2)
    }
    showToast('Script executed successfully!')
  } catch (err) {
    tab.queryOutput = 'Error running script:\n' + err
  } finally {
    tab.runningQuery = false
  }
}

// Copy Node Wizard
function openCopyWizard(tab) {
  activeTabContext.value = tab
  activeCopySource.value = `${tab.projectId} (${tab.email})`
  copyWizard.path = tab.activeCollection || ''
  copyWizard.type = 'collection'
  copyWizard.destEmail = accounts.value[0]?.email || ''
  copyWizard.destProjectId = ''
  copyWizard.overwrite = false
  copyWizard.recursive = true
  loadDestProjects()
  showCopyModal.value = true
}

async function loadDestProjects() {
  if (!copyWizard.destEmail) return
  destProjects.value = []
  try {
    const res = await ListProjects(copyWizard.destEmail)
    destProjects.value = res || []
    if (destProjects.value.length > 0) {
      copyWizard.destProjectId = destProjects.value[0].projectId
    }
  } catch (err) {
    console.error('Failed to load destination projects:', err)
  }
}

async function runCopyConflictCheck() {
  const tab = activeTabContext.value
  if (!tab) return
  
  copyingNode.value = true
  try {
    const check = await CheckConflict(copyWizard.destEmail, copyWizard.destProjectId, '', copyWizard.path, copyWizard.type)
    if (check.conflict && !copyWizard.overwrite) {
      const proceed = confirm(`${check.message} Do you want to proceed and overwrite?`)
      if (!proceed) {
        copyingNode.value = false
        return
      }
      copyWizard.overwrite = true
    }

    const copyRes = await CopyNode(
      tab.email, tab.projectId, '',
      copyWizard.destEmail, copyWizard.destProjectId, '',
      copyWizard.path, copyWizard.type,
      copyWizard.overwrite, copyWizard.recursive
    )

    alert(`Successfully copied! Copied count: ${copyRes.copiedCount}`)
    showCopyModal.value = false
  } catch (err) {
    alert('Copy operation failed: ' + err)
  } finally {
    copyingNode.value = false
  }
}

// Helper functions
function getDocId(fullName) {
  if (!fullName) return ''
  const parts = fullName.split('/')
  return parts[parts.length - 1]
}

function formatDate(timestamp) {
  if (!timestamp) return '-'
  try {
    const parsed = parseInt(timestamp)
    if (!isNaN(parsed)) {
      return new Date(parsed).toLocaleString()
    }
    return new Date(timestamp).toLocaleString()
  } catch (e) {
    return timestamp
  }
}

function generateId() {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  let id = ''
  for (let i = 0; i < 20; i++) {
    id += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  return id
}

// Subcollections Management
function getDocRelativePath(tab, doc) {
  if (!doc) return ''
  if (doc.name && doc.name.includes('/documents/')) {
    return doc.name.split('/documents/')[1]
  }
  const docId = getDocId(doc.name || '')
  if (tab && tab.activeCollection) {
    return tab.activeCollection + '/' + docId
  }
  return docId
}

async function fetchSubcollectionsForDoc(tab, doc) {
  if (!tab || !doc) return
  if (!tab.subcollectionsMap) {
    tab.subcollectionsMap = {}
  }
  const docPath = getDocRelativePath(tab, doc)
  if (!docPath) return

  tab.subcollectionsMap[docPath] = {
    list: tab.subcollectionsMap[docPath]?.list || [],
    loading: true,
    fetched: false
  }

  try {
    const res = await ListDocumentSubcollections(tab.email, tab.projectId, '', docPath)
    tab.subcollectionsMap[docPath] = {
      list: res || [],
      loading: false,
      fetched: true
    }
  } catch (err) {
    console.error('Failed to list subcollections for ' + docPath, err)
    tab.subcollectionsMap[docPath] = {
      list: [],
      loading: false,
      fetched: true,
      error: String(err)
    }
  }
}

function getDocSubcollections(tab, doc) {
  const docPath = getDocRelativePath(tab, doc)
  return (tab && tab.subcollectionsMap && tab.subcollectionsMap[docPath]) ? tab.subcollectionsMap[docPath].list : []
}

function getDocSubcollectionsLoading(tab, doc) {
  const docPath = getDocRelativePath(tab, doc)
  return (tab && tab.subcollectionsMap && tab.subcollectionsMap[docPath]) ? tab.subcollectionsMap[docPath].loading : false
}

function getDocSubcollectionsFetched(tab, doc) {
  const docPath = getDocRelativePath(tab, doc)
  return (tab && tab.subcollectionsMap && tab.subcollectionsMap[docPath]) ? tab.subcollectionsMap[docPath].fetched : false
}

function openSubcollection(tab, doc, subColName) {
  const docPath = getDocRelativePath(tab, doc)
  const fullSubColPath = `${docPath}/${subColName}`
  selectCollection(tab, fullSubColPath)
}

function openSubcollectionInNewTab(parentTab, doc, subColName) {
  const docPath = getDocRelativePath(parentTab, doc)
  const fullSubColPath = `${docPath}/${subColName}`
  openCollectionInNewTab(parentTab, fullSubColPath)
}

function getCollectionBreadcrumbs(collPath) {
  if (!collPath) return []
  const parts = collPath.split('/')
  const breadcrumbs = []
  let currentPath = ''
  for (let i = 0; i < parts.length; i++) {
    currentPath = currentPath ? `${currentPath}/${parts[i]}` : parts[i]
    breadcrumbs.push({
      name: parts[i],
      path: currentPath,
      isCollection: i % 2 === 0
    })
  }
  return breadcrumbs
}
</script>
