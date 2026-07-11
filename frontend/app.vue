<template>
  <UApp>
    <div class="h-screen w-screen flex overflow-hidden bg-[#0c0d0e] text-[#e3e6e8] font-sans antialiased select-none">
    
    <!-- Sidebar -->
    <aside class="w-80 border-r border-[#2a2d31] bg-[#121315] flex flex-col flex-shrink-0">
      <!-- Sidebar Header -->
      <div class="p-4 border-b border-[#2a2d31] flex items-center justify-between">
        <div class="flex items-center gap-2 cursor-pointer select-none" @click="reopenWelcomeTab" title="Show Welcome Screen">
          <!-- Logo -->
          <img src="/logo.svg" class="w-6 h-6 object-contain" alt="Logo" />
          <h1 class="font-bold text-lg font-display tracking-wide text-white">Mikhaz Fire</h1>
        </div>
        <!-- Add Service Account Button -->
        <UButton
          icon="i-heroicons-plus"
          size="xs"
          color="neutral"
          variant="ghost"
          @click="showAddSaModal = true"
          title="Add Service Account"
        />
      </div>

      <!-- Accounts Scan Status -->
      <div class="px-4 py-2 bg-[#16181a] text-xs text-[#9da5ac] flex items-center justify-between border-b border-[#2a2d31]">
        <span>Accounts Discovered: {{ accounts.length }}</span>
        <UButton
          size="xs"
          variant="link"
          color="neutral"
          class="p-0 text-amber-500 hover:text-amber-400"
          @click="refreshAccounts"
          :loading="scanning"
        >
          Rescan
        </UButton>
      </div>

      <!-- Accounts / Projects List -->
      <div class="flex-1 overflow-y-auto p-2 space-y-1">
        <div v-if="accounts.length === 0" class="text-xs text-[#606870] p-4 text-center">
          No accounts found. Use gcloud or firebase login, or add a service account.
        </div>
        
        <div v-for="acc in accounts" :key="acc.email" class="space-y-1">
          <!-- Account Item Header -->
          <div 
            class="flex items-center justify-between p-2 rounded-md cursor-pointer text-sm transition-colors"
            :class="expandedAccount === acc.email ? 'bg-[#1d1f22] text-white' : 'hover:bg-[#1d1f22] text-[#9da5ac]'"
            @click="toggleAccount(acc.email)"
          >
            <div class="flex items-center gap-2 truncate">
              <span class="text-sm transform transition-transform" :class="expandedAccount === acc.email ? 'rotate-90' : ''">&#9656;</span>
              <span class="truncate font-medium" :title="acc.email">{{ acc.email }}</span>
            </div>
            <UBadge size="sm" color="warning" variant="subtle" class="text-[10px] px-1.5 py-0.5">
              {{ acc.is_service_account ? 'SA' : acc.source }}
            </UBadge>
          </div>

          <!-- Projects List for Expanded Account -->
          <div v-if="expandedAccount === acc.email" class="pl-4 pr-1 space-y-1">
            <div v-if="loadingProjects" class="text-xs text-[#606870] py-2 pl-3">Loading projects...</div>
            <div v-else-if="projects.length === 0" class="text-xs text-[#606870] py-2 pl-3">No projects found.</div>
            
            <div 
              v-for="proj in projects" 
              :key="proj.projectId" 
              class="group/proj p-2 rounded-md hover:bg-[#1d1f22] text-xs text-[#9da5ac] flex flex-col gap-1 cursor-pointer transition-colors"
            >
              <div class="flex items-center justify-between">
                <span class="font-semibold text-white truncate" :title="proj.displayName">{{ proj.displayName || proj.projectId }}</span>
              </div>
              <span class="text-[10px] text-[#606870] truncate">{{ proj.projectId }}</span>
              
              <!-- Quick Action Options on Hover -->
              <div class="flex gap-2 mt-1 opacity-0 group-hover/proj:opacity-100 transition-opacity">
                <UButton
                  size="xs"
                  color="warning"
                  variant="subtle"
                  class="text-[10px] py-0.5"
                  @click.stop="openTab('auth', acc.email, proj.projectId)"
                >
                  Auth
                </UButton>
                <UButton
                  size="xs"
                  color="warning"
                  variant="subtle"
                  class="text-[10px] py-0.5"
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
      <div class="p-4 bg-[#0c0d0e] border-t border-[#2a2d31] flex items-center justify-between text-xs text-[#606870]">
        <div class="flex items-center gap-1.5">
          <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
          <span>Desktop Mode</span>
        </div>
        <span>v1.0.0</span>
      </div>
    </aside>

    <!-- Main Content Area -->
    <main class="flex-1 flex flex-col overflow-hidden">
      <!-- Tabs Bar -->
      <div class="h-11 bg-[#121315] border-b border-[#2a2d31] flex items-center overflow-x-auto select-none scrollbar-none">
        <!-- Welcome Tab -->
        <div 
          v-if="showWelcomeTab"
          class="px-4 h-full flex items-center gap-2 text-xs border-r border-[#2a2d31] cursor-pointer transition-colors group"
          :class="activeTabId === 'welcome' ? 'bg-[#0c0d0e] text-white font-medium border-t-2 border-t-amber-500' : 'hover:bg-[#1d1f22] text-[#9da5ac]'"
          @click="activeTabId = 'welcome'"
        >
          <span>Welcome</span>
          <span 
            class="text-[#606870] hover:text-white rounded-full p-0.5 leading-none transition-colors"
            @click.stop="closeWelcomeTab"
          >
            &times;
          </span>
        </div>

        <!-- Dynamic Tabs -->
        <div 
          v-for="tab in openTabs" 
          :key="tab.id"
          class="px-3 h-full flex items-center gap-2 text-xs border-r border-[#2a2d31] cursor-pointer transition-colors group"
          :class="activeTabId === tab.id ? 'bg-[#0c0d0e] text-white font-medium border-t-2 border-t-amber-500' : 'hover:bg-[#1d1f22] text-[#9da5ac]'"
          @click="activeTabId = tab.id"
        >
          <span class="font-semibold text-amber-500 truncate max-w-[160px]" :title="getTabTitle(tab)">{{ getTabTitle(tab) }}</span>
          <span 
            class="text-[#606870] hover:text-white rounded-full p-0.5 leading-none transition-colors"
            @click.stop="closeTab(tab.id)"
          >
            &times;
          </span>
        </div>
      </div>

      <!-- Tab Contents Container -->
      <div class="flex-1 overflow-hidden relative">
        <!-- Welcome Screen -->
        <div v-if="activeTabId === 'welcome'" class="h-full overflow-y-auto p-8 flex flex-col justify-center items-center text-center max-w-3xl mx-auto">
          <img src="/logo.svg" class="w-20 h-20 mb-6 drop-shadow-[0_0_15px_rgba(245,158,11,0.2)] object-contain" alt="Logo" />
          <h2 class="text-3xl font-bold font-display text-white mb-2">Welcome to Mikhaz Fire</h2>
          <p class="text-[#9da5ac] mb-8 text-sm">
            A fast, local, and premium Firebase Auth and Firestore controller. 
            Connect your credentials locally and administer database structures and accounts with ease.
          </p>
          
          <div class="grid grid-cols-2 gap-4 w-full">
            <div class="p-5 bg-[#121315] border border-[#2a2d31] rounded-lg text-left transition-all hover:border-amber-500/50">
              <h3 class="font-bold text-white mb-1 flex items-center gap-2">
                <span class="text-amber-500">&#10024;</span> Local Credential Scan
              </h3>
              <p class="text-xs text-[#9da5ac]">
                Automatically reads from standard gcloud ADC and Firebase CLI profiles stored in AppData. Secure and offline.
              </p>
            </div>
            <div class="p-5 bg-[#121315] border border-[#2a2d31] rounded-lg text-left transition-all hover:border-amber-500/50">
              <h3 class="font-bold text-white mb-1 flex items-center gap-2">
                <span class="text-amber-500">&#128274;</span> Auth Controller
              </h3>
              <p class="text-xs text-[#9da5ac]">
                Create, update, delete auth accounts, update display names, passwords, and custom user claims in real-time.
              </p>
            </div>
            <div class="p-5 bg-[#121315] border border-[#2a2d31] rounded-lg text-left transition-all hover:border-amber-500/50">
              <h3 class="font-bold text-white mb-1 flex items-center gap-2">
                <span class="text-amber-500">&#128452;</span> Firestore Explorer
              </h3>
              <p class="text-xs text-[#9da5ac]">
                Navigate root collections, document graphs, inspect fields in table formats, and perform edits directly.
              </p>
            </div>
            <div class="p-5 bg-[#121315] border border-[#2a2d31] rounded-lg text-left transition-all hover:border-amber-500/50">
              <h3 class="font-bold text-white mb-1 flex items-center gap-2">
                <span class="text-amber-500">&#9889;</span> JS Sandbox Queries
              </h3>
              <p class="text-xs text-[#9da5ac]">
                Run custom JavaScript scripts in an embedded JS runtime with access to mocked Firestore admin APIs to batch modify.
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
          <div v-if="tab.type === 'auth'" class="h-full flex flex-col overflow-hidden p-6 space-y-4">
            <div class="flex items-center justify-between">
              <div>
                <h2 class="text-xl font-bold text-white">Auth User Manager</h2>
                <p class="text-xs text-[#9da5ac]">{{ tab.projectId }} &middot; {{ tab.email }} &middot; {{ tab.users ? tab.users.length : 0 }} loaded users</p>
              </div>
              <UButton
                color="warning"
                size="sm"
                @click="openCreateUserModal(tab)"
              >
                Create User
              </UButton>
            </div>

            <!-- Search and Control Bar -->
            <div class="flex gap-4">
              <UInput
                v-model="tab.searchQuery"
                placeholder="Search by UID, email, display name..."
                class="flex-1"
                icon="i-heroicons-magnifying-glass"
                @keydown.enter="fetchAuthUsers(tab)"
              />
              <UButton color="neutral" size="sm" @click="fetchAuthUsers(tab)">
                Search
              </UButton>
            </div>

            <!-- Users Table -->
            <div class="flex-1 border border-[#2a2d31] bg-[#121315] rounded-lg overflow-hidden flex flex-col">
              <div class="flex-1 overflow-auto">
                <table class="w-full text-left border-collapse text-xs">
                  <thead>
                    <tr class="bg-[#1d1f22] text-[#e3e6e8] border-b border-[#2a2d31] font-semibold">
                      <th class="p-3">UID</th>
                      <th class="p-3">Email</th>
                      <th class="p-3">Display Name</th>
                      <th class="p-3">Disabled</th>
                      <th class="p-3">Created</th>
                      <th class="p-3">Last Login</th>
                      <th class="p-3 text-right">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-if="tab.loading" class="text-center">
                      <td colspan="7" class="p-8 text-[#9da5ac]">Loading users...</td>
                    </tr>
                    <tr v-else-if="!tab.users || tab.users.length === 0" class="text-center">
                      <td colspan="7" class="p-8 text-[#606870]">No users found.</td>
                    </tr>
                    <tr v-else v-for="user in tab.users" :key="user.uid" class="border-b border-[#2a2d31] hover:bg-[#1d1f22] transition-colors">
                      <td class="p-3 font-mono text-[11px] text-amber-500 truncate max-w-[150px]" :title="user.uid">{{ user.uid }}</td>
                      <td class="p-3 truncate max-w-[150px]" :title="user.email">{{ user.email || 'Anonymous' }}</td>
                      <td class="p-3 truncate max-w-[150px]" :title="user.displayName">{{ user.displayName || '-' }}</td>
                      <td class="p-3">
                        <UBadge :color="user.disabled ? 'danger' : 'success'" variant="subtle" size="sm">
                          {{ user.disabled ? 'Yes' : 'No' }}
                        </UBadge>
                      </td>
                      <td class="p-3 text-[#606870]">{{ formatDate(user.createdAt) }}</td>
                      <td class="p-3 text-[#606870]">{{ formatDate(user.lastLoginAt) }}</td>
                      <td class="p-3 text-right space-x-1.5">
                        <UButton size="xs" variant="ghost" color="neutral" @click="editUserClaims(tab, user)">Claims</UButton>
                        <UButton size="xs" variant="ghost" color="neutral" @click="editUserBasic(tab, user)">Edit</UButton>
                        <UButton size="xs" variant="ghost" color="danger" @click="deleteUserConfirm(tab, user)">Delete</UButton>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
              <!-- Pagination / Next Page -->
              <div v-if="tab.nextPageToken" class="p-3 bg-[#1d1f22] border-t border-[#2a2d31] flex justify-end">
                <UButton size="xs" color="warning" @click="fetchAuthUsersNext(tab)">
                  Load More Users
                </UButton>
              </div>
            </div>
          </div>

          <!-- Firestore Component -->
          <div v-else-if="tab.type === 'firestore'" class="h-full flex overflow-hidden">
            <!-- Left Pane: Collections / Documents Explorer -->
            <div class="w-80 border-r border-[#2a2d31] bg-[#121315] flex flex-col">
              <div class="p-4 border-b border-[#2a2d31] flex items-center justify-between">
                <h3 class="font-bold text-sm text-white">Collections</h3>
                <div class="flex items-center gap-1">
                  <UButton size="xs" color="neutral" variant="ghost" icon="i-heroicons-arrow-path" @click="refreshCollections(tab)" />
                </div>
              </div>
              <div class="flex-1 overflow-y-auto p-2 space-y-1">
                <div v-if="tab.loadingCollections" class="text-xs text-[#606870] p-4 text-center">Loading root collections...</div>
                <div v-else-if="tab.collections.length === 0" class="text-xs text-[#606870] p-4 text-center">No collections found.</div>
                
                <div v-for="col in tab.collections" :key="col" class="space-y-1">
                  <div 
                    class="p-2 rounded-md hover:bg-[#1d1f22] text-xs font-semibold cursor-pointer text-[#e3e6e8] flex justify-between items-center"
                    :class="tab.activeCollection === col ? 'bg-[#1d1f22] border-l-2 border-amber-500' : ''"
                    @click="openCollectionInNewTab(tab, col)"
                  >
                    <span class="truncate">{{ col }}</span>
                    <span class="text-[10px] text-amber-500">Col</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Right Pane: Document Viewer / Editor / JS Queries -->
            <div class="flex-1 flex flex-col overflow-hidden bg-[#0c0d0e]">
              <!-- Header -->
              <div class="p-4 border-b border-[#2a2d31] bg-[#121315] flex justify-between items-center">
                <div>
                  <h3 class="font-bold text-white flex items-center gap-2">
                    <span class="capitalize text-amber-500">{{ tab.firestoreView || 'Explorer' }}</span>
                    <span v-if="tab.firestoreView === 'explorer' && tab.activeCollection" class="text-xs text-[#9da5ac]">&middot; {{ tab.activeCollection }}</span>
                  </h3>
                </div>
                <div class="flex gap-2">
                  <UButton size="xs" color="neutral" variant="subtle" @click="openCopyWizard(tab)">
                    Copy Wizard
                  </UButton>
                  <UButton v-if="tab.firestoreView === 'explorer' && tab.activeCollection" size="xs" color="warning" @click="openAddDocModal(tab)">
                    Add Document
                  </UButton>
                </div>
              </div>

              <!-- Main Viewport -->
              <div class="flex-1 overflow-hidden">
                <!-- Explorer View -->
                <div v-if="tab.firestoreView === 'explorer'" class="h-full flex flex-col p-6 space-y-4 overflow-hidden">
                  <div v-if="!tab.activeCollection" class="h-full flex justify-center items-center text-[#606870] text-sm">
                    Select a collection on the left sidebar to browse documents.
                  </div>
                  <div v-else class="flex-1 flex flex-col space-y-4 min-h-0">
                    
                    <!-- Query Builder Panel (Above the Results) -->
                    <div class="border border-[#2a2d31] bg-[#121315] rounded-lg p-4 space-y-3 flex flex-col">
                      <div class="flex items-center justify-between border-b border-[#2a2d31] pb-2">
                        <div class="flex items-center gap-2">
                          <span class="text-xs font-bold text-white uppercase">Query Builder</span>
                          <span class="text-[10px] bg-[#2a2d31] text-[#9da5ac] px-1.5 py-0.5 rounded font-mono">{{ tab.activeCollection }}</span>
                        </div>
                        <div class="flex items-center gap-1 bg-[#0c0d0e] p-0.5 rounded border border-[#2a2d31]">
                          <UButton 
                            size="xs" 
                            :color="tab.queryMode === 'simple' ? 'warning' : 'neutral'" 
                            variant="ghost" 
                            @click="setQueryMode(tab, 'simple')"
                          >
                            Simple
                          </UButton>
                          <UButton 
                            size="xs" 
                            :color="tab.queryMode === 'js' ? 'warning' : 'neutral'" 
                            variant="ghost" 
                            @click="setQueryMode(tab, 'js')"
                          >
                            JS Query
                          </UButton>
                        </div>
                      </div>

                      <!-- Simple Query Form -->
                      <div v-if="tab.queryMode === 'simple'" class="grid grid-cols-4 gap-3 text-xs">
                        <div class="flex flex-col gap-1">
                          <label class="font-semibold text-[#9da5ac]">Document ID (UID)</label>
                          <UInput v-model="tab.queryDocId" placeholder="Optional ID" class="w-full bg-[#0c0d0e]" />
                        </div>
                        <div class="flex flex-col gap-1">
                          <label class="font-semibold text-[#9da5ac]">Filter Field</label>
                          <UInput v-model="tab.queryFilterField" placeholder="e.g. role" class="w-full bg-[#0c0d0e]" />
                        </div>
                        <div class="flex flex-col gap-1">
                          <label class="font-semibold text-[#9da5ac]">Op</label>
                          <USelect 
                            v-model="tab.queryFilterOp" 
                            :items="['==', '>', '>=', '<', '<=', '!=', 'array-contains']" 
                            class="w-full bg-[#0c0d0e]" 
                          />
                        </div>
                        <div class="flex flex-col gap-1">
                          <label class="font-semibold text-[#9da5ac]">Value</label>
                          <UInput v-model="tab.queryFilterVal" placeholder="value" class="w-full bg-[#0c0d0e]" />
                        </div>
                        <div class="flex flex-col gap-1">
                          <label class="font-semibold text-[#9da5ac]">Sort Field</label>
                          <UInput v-model="tab.querySortField" placeholder="e.g. created_at" class="w-full bg-[#0c0d0e]" />
                        </div>
                        <div class="flex flex-col gap-1">
                          <label class="font-semibold text-[#9da5ac]">Sort Direction</label>
                          <USelect 
                            v-model="tab.querySortDir" 
                            :items="[{label:'Ascending',value:'asc'}, {label:'Descending',value:'desc'}]" 
                            class="w-full bg-[#0c0d0e]" 
                          />
                        </div>
                        <div class="flex flex-col gap-1">
                          <label class="font-semibold text-[#9da5ac]">Limit (Quantity)</label>
                          <UInput v-model.number="tab.queryLimit" type="number" placeholder="50" class="w-full bg-[#0c0d0e]" />
                        </div>
                        <div class="flex items-end justify-end">
                          <UButton color="warning" @click="runSimpleQueryAction(tab)" :loading="tab.loadingDocuments" class="w-full justify-center">
                            Run Simple Query
                          </UButton>
                        </div>
                      </div>

                      <!-- JS Query Form -->
                      <div v-else class="space-y-3 flex flex-col">
                        <div class="flex justify-between items-center text-xs">
                          <div class="flex items-center gap-2">
                            <label class="font-semibold text-[#9da5ac]">Preset Template:</label>
                            <USelect 
                              v-model="tab.queryPreset" 
                              :items="['Basic Fetch', 'Filter & Order', 'Projection', 'Aggregate Sum', 'Subcollection Fetch']" 
                              class="bg-[#0c0d0e] shrink-0" 
                              @update:model-value="loadPresetQueryScript(tab)"
                            />
                          </div>
                          <span class="text-[10px] text-[#606870]">Runs local mock JS engine against DB</span>
                        </div>
                        <div class="grid grid-cols-2 gap-4">
                          <!-- Editor -->
                          <div class="flex flex-col space-y-2">
                            <UTextarea 
                              v-model="tab.queryScript" 
                              class="font-mono text-xs bg-[#0c0d0e] flex-1" 
                              rows="10"
                              placeholder="// JavaScript query script..." 
                            />
                            <div class="flex justify-end">
                              <UButton color="warning" @click="runQueryScript(tab)" :loading="tab.runningQuery">
                                Run JS Script
                              </UButton>
                            </div>
                          </div>
                          <!-- Console Output -->
                          <div class="flex flex-col space-y-1">
                            <label class="text-[10px] font-semibold text-[#9da5ac]">Query Output Logs</label>
                            <pre class="flex-1 font-mono text-[10px] text-green-400 overflow-auto bg-[#0c0d0e] p-3 rounded border border-[#2a2d31] whitespace-pre-wrap min-h-[200px]">{{ tab.queryOutput || 'No output yet.' }}</pre>
                          </div>
                        </div>
                      </div>
                    </div>

                    <!-- Results Section (Table / Tree / JSON) -->
                    <div class="flex-1 flex flex-col space-y-3 min-h-0">
                      <!-- Results Navigation Tabs and Selection Toolbar -->
                      <div class="flex items-center justify-between border-b border-[#2a2d31] pb-2">
                        <div class="flex items-center gap-1 bg-[#121315] p-1 rounded-md border border-[#2a2d31]">
                          <UButton 
                            size="xs" 
                            :color="tab.resultsView === 'tree' ? 'warning' : 'neutral'" 
                            variant="ghost" 
                            @click="tab.resultsView = 'tree'"
                          >
                            Tree
                          </UButton>
                          <UButton 
                            size="xs" 
                            :color="tab.resultsView === 'table' ? 'warning' : 'neutral'" 
                            variant="ghost" 
                            @click="tab.resultsView = 'table'"
                          >
                            Table
                          </UButton>
                          <UButton 
                            size="xs" 
                            :color="tab.resultsView === 'json' ? 'warning' : 'neutral'" 
                            variant="ghost" 
                            @click="tab.resultsView = 'json'"
                          >
                            JSON
                          </UButton>
                        </div>

                        <!-- Selection Action Toolbar -->
                        <div class="flex items-center justify-between gap-3">
                          <span class="text-xs text-[#9da5ac] font-medium" v-if="tab.selectedDocNames && tab.selectedDocNames.length > 0">
                            {{ tab.selectedDocNames.length }} selected
                          </span>
                          <div v-if="tab.selectedDocNames && tab.selectedDocNames.length > 0" class="flex gap-2">
                            <UButton 
                              size="xs" 
                              color="danger" 
                              @click="deleteSelectedDocuments(tab)"
                            >
                              Delete Selected
                            </UButton>
                            <UButton 
                              size="xs" 
                              color="neutral" 
                              variant="subtle"
                              @click="triggerExportSelected(tab)"
                            >
                              Export Selected
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
                      <div v-if="tab.resultsView === 'tree'" class="flex-1 border border-[#2a2d31] bg-[#121315] rounded-lg p-4 overflow-y-auto flex flex-col gap-2 min-h-0">
                        <div class="flex items-center gap-2 mb-2 pb-2 border-b border-[#2a2d31]">
                          <input 
                            type="checkbox" 
                            :checked="tab.selectedDocNames && tab.selectedDocNames.length === tab.documents.length && tab.documents.length > 0"
                            @change="toggleSelectAllDocs(tab)"
                            class="accent-amber-500 rounded"
                          />
                          <span class="text-xs text-[#9da5ac]">Select All Documents</span>
                        </div>
                        
                        <div v-if="tab.loadingDocuments" class="text-xs text-[#9da5ac] text-center p-8">Loading documents...</div>
                        <div v-else-if="!tab.documents || tab.documents.length === 0" class="text-xs text-[#606870] text-center p-8">No documents found.</div>
                        <div v-else class="space-y-3">
                          <div v-for="doc in tab.documents" :key="doc.name" class="border border-[#1d1f22] hover:border-[#2a2d31] rounded-md p-2 bg-[#0c0d0e]/50 flex items-start gap-3">
                            <input 
                              type="checkbox" 
                              :checked="tab.selectedDocNames && tab.selectedDocNames.includes(doc.name)"
                              @change="toggleSelectDoc(tab, doc.name)"
                              class="accent-amber-500 mt-1"
                            />
                            <div class="flex-1 min-w-0">
                              <JsonTree 
                                :label="getDocId(doc.name)"
                                :value="cleanDocFields(doc.fields)"
                                :is-root="true"
                                @change="(fieldPath, val) => saveDocFieldChange(tab, doc, fieldPath, val)"
                                @editJson="editDocumentJSON(tab, doc)"
                                @deleteDoc="deleteDocumentConfirm(tab, doc)"
                                @openTab="openDocInNewTab(tab, doc)"
                              />
                            </div>
                          </div>
                        </div>
                      </div>

                      <!-- Table Viewer -->
                      <div v-else-if="tab.resultsView === 'table'" class="flex-1 border border-[#2a2d31] bg-[#121315] rounded-lg overflow-hidden flex flex-col min-h-0">
                        <div class="flex-1 overflow-auto">
                          <table class="w-full text-left border-collapse text-xs">
                            <thead>
                              <tr class="bg-[#1d1f22] text-[#e3e6e8] border-b border-[#2a2d31] font-semibold select-none">
                                <th class="p-3 w-10 text-center">
                                  <input 
                                    type="checkbox" 
                                    :checked="tab.selectedDocNames && tab.selectedDocNames.length === tab.documents.length && tab.documents.length > 0"
                                    @change="toggleSelectAllDocs(tab)"
                                    class="accent-amber-500 rounded"
                                  />
                                </th>
                                <th class="p-3 w-1/4">Document ID</th>
                                <th v-for="col in getTableColumns(tab.documents)" :key="col" class="p-3">{{ col }}</th>
                                <th class="p-3 text-right">Actions</th>
                              </tr>
                            </thead>
                            <tbody>
                              <tr v-if="tab.loadingDocuments" class="text-center">
                                <td :colspan="getTableColumns(tab.documents).length + 3" class="p-8 text-[#9da5ac]">Loading documents...</td>
                              </tr>
                              <tr v-else-if="!tab.documents || tab.documents.length === 0" class="text-center">
                                <td :colspan="getTableColumns(tab.documents).length + 3" class="p-8 text-[#606870]">No documents found.</td>
                              </tr>
                              <tr 
                                v-else 
                                v-for="doc in tab.documents" 
                                :key="doc.name" 
                                class="border-b border-[#2a2d31] hover:bg-[#1d1f22] transition-colors cursor-pointer"
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
                                <td class="p-3 font-mono font-semibold text-amber-500 truncate" :title="getDocId(doc.name)">
                                  {{ getDocId(doc.name) }}
                                </td>
                                <td v-for="col in getTableColumns(tab.documents)" :key="col" class="p-3 font-mono text-[10px] text-[#9da5ac] truncate max-w-[150px]">
                                  {{ getColumnValue(doc, col) }}
                                </td>
                                <td class="p-3 text-right space-x-2" @click.stop>
                                  <UButton size="xs" variant="ghost" color="neutral" @click="openDocInNewTab(tab, doc)">Open Tab</UButton>
                                  <UButton size="xs" variant="ghost" color="neutral" @click="editDocumentJSON(tab, doc)">Edit JSON</UButton>
                                  <UButton size="xs" variant="ghost" color="danger" @click="deleteDocumentConfirm(tab, doc)">Delete</UButton>
                                </td>
                              </tr>
                            </tbody>
                          </table>
                        </div>
                      </div>

                      <!-- JSON Viewer -->
                      <div v-else-if="tab.resultsView === 'json'" class="flex-1 border border-[#2a2d31] bg-[#121315] rounded-lg p-4 overflow-hidden flex flex-col min-h-0">
                        <textarea 
                          class="flex-grow flex-1 w-full h-full bg-[#0c0d0e] text-[#9da5ac] font-mono text-xs p-3 rounded border border-[#2a2d31] resize-none outline-none select-text min-h-0" 
                          readonly
                        >{{ formatAllDocsJson(tab.documents) }}</textarea>
                      </div>

                      <!-- Load More Button -->
                      <div v-if="tab.nextDocPageToken" class="p-3 bg-[#1d1f22] border-t border-[#2a2d31] flex justify-end rounded-b-lg">
                        <UButton size="xs" color="warning" @click="fetchDocumentsNext(tab)">
                          Load More
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
        <UCard :ui="{ divide: 'y-divide-[#2a2d31]', background: 'bg-[#121315]', border: 'border-[#2a2d31]' }">
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-base font-semibold leading-6 text-white">Add Google Service Account</h3>
              <UButton color="neutral" variant="ghost" icon="i-heroicons-x-mark" class="-my-1" @click="showAddSaModal = false" />
            </div>
          </template>

          <div class="space-y-4 p-2 text-xs">
            <p class="text-[#9da5ac]">
              Paste the contents of your Google Service Account private key JSON file below. 
              This file will be loaded <strong>only in-memory</strong> and will not be saved to disk.
            </p>
            <UTextarea 
              v-model="newSaJson" 
              placeholder='{ "type": "service_account", "project_id": ... }' 
              rows="10" 
              class="font-mono text-xs bg-[#0c0d0e]" 
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
        <UCard :ui="{ background: 'bg-[#121315]', border: 'border-[#2a2d31]' }">
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-base font-semibold leading-6 text-white">Edit Custom Claims</h3>
              <UButton color="neutral" variant="ghost" icon="i-heroicons-x-mark" class="-my-1" @click="showClaimsModal = false" />
            </div>
          </template>

          <div class="space-y-4 p-2 text-xs">
            <p class="text-[#9da5ac]">
              Manage user claims for <strong>{{ activeUser?.email || activeUser?.uid }}</strong>. Custom claims must be valid JSON object or empty.
            </p>
            <UTextarea 
              v-model="claimsJson" 
              placeholder='{ "admin": true, "role": "editor" }' 
              rows="6" 
              class="font-mono text-xs bg-[#0c0d0e]" 
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
        <UCard :ui="{ background: 'bg-[#121315]', border: 'border-[#2a2d31]' }">
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-base font-semibold leading-6 text-white">Create Firebase Auth User</h3>
              <UButton color="neutral" variant="ghost" icon="i-heroicons-x-mark" class="-my-1" @click="showCreateUserModal = false" />
            </div>
          </template>

          <div class="space-y-4 p-2 text-xs">
            <div>
              <label class="block text-xs font-semibold text-[#9da5ac] mb-1">Email Address</label>
              <UInput v-model="newUser.email" placeholder="user@example.com" class="w-full bg-[#0c0d0e]" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-[#9da5ac] mb-1">Display Name</label>
              <UInput v-model="newUser.displayName" placeholder="John Doe" class="w-full bg-[#0c0d0e]" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-[#9da5ac] mb-1">Password</label>
              <UInput v-model="newUser.password" type="password" placeholder="••••••••" class="w-full bg-[#0c0d0e]" />
            </div>
            <div class="flex items-center gap-2">
              <UCheckbox v-model="newUser.emailVerified" label="Email Verified" />
              <UCheckbox v-model="newUser.disabled" label="Disabled" />
            </div>
          </div>

          <template #footer>
            <div class="flex justify-end gap-2">
              <UButton color="neutral" variant="ghost" @click="showCreateUserModal = false">Cancel</UButton>
              <UButton color="warning" @click="saveNewUser">Create</UButton>
            </div>
          </template>
        </UCard>
      </template>
    </UModal>

    <!-- Modal: Edit Auth User Basic -->
    <UModal v-model:open="showEditUserModal" :dismissible="false">
      <template #content>
        <UCard :ui="{ background: 'bg-[#121315]', border: 'border-[#2a2d31]' }">
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-base font-semibold leading-6 text-white">Edit User Profile</h3>
              <UButton color="neutral" variant="ghost" icon="i-heroicons-x-mark" class="-my-1" @click="showEditUserModal = false" />
            </div>
          </template>

          <div class="space-y-4 p-2 text-xs">
            <div>
              <label class="block text-xs font-semibold text-[#9da5ac] mb-1">Email Address</label>
              <UInput v-model="editUser.email" placeholder="user@example.com" class="w-full bg-[#0c0d0e]" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-[#9da5ac] mb-1">Display Name</label>
              <UInput v-model="editUser.displayName" placeholder="John Doe" class="w-full bg-[#0c0d0e]" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-[#9da5ac] mb-1">Password (Leave empty to keep unchanged)</label>
              <UInput v-model="editUser.password" type="password" placeholder="••••••••" class="w-full bg-[#0c0d0e]" />
            </div>
            <div class="flex items-center gap-2">
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

    <!-- Modal: Document Details JSON Viewer/Editor -->
    <UModal v-model:open="showDocModal" :dismissible="false" fullscreen :ui="{ content: 'h-full flex flex-col bg-[#0c0d0e]' }">
      <template #content>
        <UCard :ui="{ root: 'h-full flex flex-col bg-[#0c0d0e] border border-[#2a2d31]', body: 'flex-grow overflow-hidden flex flex-col h-full min-h-0' }">
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-base font-semibold leading-6 text-white">Document Editor JSON</h3>
              <UButton color="neutral" variant="ghost" icon="i-heroicons-x-mark" class="-my-1" @click="showDocModal = false" />
            </div>
          </template>

          <div class="space-y-4 p-4 text-xs flex flex-col flex-grow flex-1 h-full min-h-0 bg-[#0c0d0e]">
            <div class="bg-[#0c0d0e] p-2 rounded text-[11px] border border-[#2a2d31] text-[#9da5ac]">
              Path: <strong class="text-amber-500 font-mono">{{ activeDocPath }}</strong>
            </div>
            <label class="block text-xs font-semibold text-[#9da5ac]">Fields Data JSON (Raw Key-Value pairs)</label>
            <textarea 
              v-model="activeDocJson" 
              placeholder='{ "title": "First Post", "likes": 42 }' 
              class="flex-1 w-full bg-[#0c0d0e] text-[#e3e6e8] font-mono text-xs p-4 rounded border border-[#2a2d31] resize-none outline-none select-text min-h-0 h-full"
            ></textarea>
          </div>

          <template #footer>
            <div class="flex justify-end gap-2 bg-[#0c0d0e]">
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
        <UCard :ui="{ background: 'bg-[#121315]', border: 'border-[#2a2d31]' }">
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-base font-semibold leading-6 text-white">Create New Document</h3>
              <UButton color="neutral" variant="ghost" icon="i-heroicons-x-mark" class="-my-1" @click="showAddDocModal = false" />
            </div>
          </template>

          <div class="space-y-4 p-2 text-xs">
            <div>
              <label class="block text-xs font-semibold text-[#9da5ac] mb-1">Document ID (optional, blank auto-generates ID)</label>
              <UInput v-model="newDocId" placeholder="my-custom-id" class="w-full bg-[#0c0d0e]" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-[#9da5ac] mb-1">Fields Initial JSON</label>
              <UTextarea 
                v-model="newDocFieldsJson" 
                placeholder='{ "createdAt": "2026-07-02T10:00:00Z" }' 
                rows="6" 
                class="font-mono text-xs bg-[#0c0d0e]" 
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
        <UCard :ui="{ background: 'bg-[#121315]', border: 'border-[#2a2d31]' }">
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-base font-semibold leading-6 text-white">Firestore Copy Wizard</h3>
              <UButton color="neutral" variant="ghost" icon="i-heroicons-x-mark" class="-my-1" @click="showCopyModal = false" />
            </div>
          </template>

          <div class="space-y-4 p-2 text-xs">
            <p class="text-[#9da5ac]">
              Copy collections or documents from this project to another project safely. Checks for conflicts before execution.
            </p>
            <div>
              <label class="block text-xs font-semibold text-[#9da5ac] mb-1">Source Account & Project</label>
              <div class="bg-[#0c0d0e] p-2 rounded text-[#e3e6e8] border border-[#2a2d31] font-mono">
                {{ activeCopySource }}
              </div>
            </div>
            <div>
              <label class="block text-xs font-semibold text-[#9da5ac] mb-1">Source Path (Collection or Doc)</label>
              <UInput v-model="copyWizard.path" placeholder="posts" class="w-full bg-[#0c0d0e] font-mono" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-[#9da5ac] mb-1">Copy Node Type</label>
              <USelect 
                v-model="copyWizard.type" 
                :items="['collection', 'document']" 
                class="w-full bg-[#0c0d0e]"
              />
            </div>
            
            <div class="border-t border-[#2a2d31] pt-3">
              <label class="block text-xs font-semibold text-white mb-2">Destination Configuration</label>
              <div class="space-y-2">
                <div>
                  <label class="block text-[10px] text-[#9da5ac] mb-1">Dest Account Email</label>
                  <USelect 
                    v-model="copyWizard.destEmail" 
                    :items="accounts.map(a => a.email)" 
                    class="w-full bg-[#0c0d0e]"
                    @change="loadDestProjects"
                  />
                </div>
                <div>
                  <label class="block text-[10px] text-[#9da5ac] mb-1">Dest Project ID</label>
                  <USelect 
                    v-model="copyWizard.destProjectId" 
                    :items="destProjects.map(p => p.projectId)" 
                    class="w-full bg-[#0c0d0e]"
                  />
                </div>
              </div>
            </div>

            <div class="flex items-center gap-4 border-t border-[#2a2d31] pt-3">
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
        <UCard :ui="{ background: 'bg-[#121315]', border: 'border-[#2a2d31]' }">
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-base font-semibold leading-6 text-white">Export Results</h3>
              <UButton color="neutral" variant="ghost" icon="i-heroicons-x-mark" class="-my-1" @click="showExportModal = false" />
            </div>
          </template>

          <div class="space-y-4 p-2 text-xs">
            <p class="text-[#9da5ac]">
              Export the current list of loaded documents. You can configure the format and options below.
            </p>
            
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-xs font-semibold text-[#9da5ac] mb-1">Export Format</label>
                <USelect v-model="exportConfig.format" :items="['JSON', 'CSV']" class="w-full bg-[#0c0d0e]" @change="updateExportPreview" />
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
                class="w-full bg-[#0c0d0e] text-[#9da5ac] border border-[#2a2d31] font-mono text-[10px] p-2 rounded resize-none select-text"
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

    </div>
  </UApp>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'

// Import Wails runtime calls
import {
  GetAccounts,
  AddServiceAccount,
  ListProjects,
  ListAuthUsers,
  CreateAuthUser,
  UpdateAuthUser,
  DeleteAuthUser,
  ListRootCollections,
  ListDocuments,
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

// Auth Modal States
const showClaimsModal = ref(false)
const showCreateUserModal = ref(false)
const showEditUserModal = ref(false)
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
  refreshAccounts()
})

// Load accounts
async function refreshAccounts() {
  scanning.value = true
  try {
    const res = await GetAccounts()
    accounts.value = res || []
  } catch (err) {
    console.error('Failed to load accounts:', err)
  } finally {
    scanning.value = false
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
    const email = await AddServiceAccount(newSaJson.value.trim())
    showAddSaModal.value = false
    newSaJson.value = ''
    refreshAccounts()
  } catch (err) {
    alert('Failed to add service account: ' + err)
  } finally {
    savingSa.value = false
  }
}

function getTabTitle(tab) {
  if (tab.type === 'auth') {
    return `Auth: ${tab.projectId}`;
  }
  if (tab.type === 'firestore') {
    if (tab.activeDocPath) {
      const parts = tab.activeDocPath.split('/');
      return `Doc: ${parts[parts.length - 1]}`;
    }
    if (tab.activeCollection) {
      return `Col: ${tab.activeCollection}`;
    }
    return `Firestore: ${tab.projectId}`;
  }
  return tab.type;
}

function openCollectionInNewTab(parentTab, colName) {
  const tabId = `firestore-col-${parentTab.email}-${parentTab.projectId}-${colName}`;
  const existingTab = openTabs.value.find(t => t.id === tabId);
  
  if (existingTab) {
    activeTabId.value = tabId;
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
      activeCollection: colName,
      activeDocPath: '',
      documents: [],
      nextDocPageToken: '',
      loadingCollections: false,
      loadingDocuments: false,
      firestoreView: 'explorer',
      resultsView: 'table',
      selectedDocNames: [],
      queryMode: 'simple',
      queryDocId: '',
      queryFilterField: '',
      queryFilterOp: '==',
      queryFilterVal: '',
      querySortField: '',
      querySortDir: 'asc',
      queryLimit: 50,
      queryPreset: 'Basic Fetch',
      queryScript: `// Query with JavaScript using the Firebase Admin SDK\n// See examples at https://firefoo.app/go/firestore-js-query\nasync function run() {\n  const query = await db.collection("${colName}")\n    .limit(50)\n    .get();\n  return query;\n}`,
      queryOutput: '',
      runningQuery: false
    };
    openTabs.value.push(newTab);
    const reactiveTab = openTabs.value[openTabs.value.length - 1];
    activeTabId.value = tabId;
    fetchDocuments(reactiveTab);
  }
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
      loading: false,
      users: [],
      nextPageToken: '',
      collections: [],
      activeCollection: '',
      documents: [],
      nextDocPageToken: '',
      loadingCollections: false,
      loadingDocuments: false,
      firestoreView: 'explorer', // 'explorer' or 'query'
      resultsView: 'tree', // 'tree' | 'table' | 'json'
      selectedDocNames: [],
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
      runningQuery: false
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

// Auth Actions
async function fetchAuthUsers(tab) {
  tab.loading = true
  tab.users = []
  tab.nextPageToken = ''
  try {
    const res = await ListAuthUsers(tab.email, tab.projectId, 50, '', tab.searchQuery)
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
    const res = await ListAuthUsers(tab.email, tab.projectId, 50, tab.nextPageToken, tab.searchQuery)
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
    fetchAuthUsers(tab)
  } catch (err) {
    alert('Failed to save claims: ' + err)
  }
}

async function deleteUserConfirm(tab, user) {
  if (confirm(`Are you sure you want to delete user ${user.email || user.uid}?`)) {
    try {
      await DeleteAuthUser(tab.email, tab.projectId, user.uid)
      fetchAuthUsers(tab)
    } catch (err) {
      alert('Failed to delete user: ' + err)
    }
  }
}

// Firestore Actions
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
  tab.activeCollection = colName
  tab.firestoreView = 'explorer'
  tab.selectedDocNames = []
  
  // Reset simple query filters
  tab.queryMode = 'simple'
  tab.queryDocId = ''
  tab.queryFilterField = ''
  tab.queryFilterOp = '=='
  tab.queryFilterVal = ''
  tab.querySortField = ''
  tab.querySortDir = 'asc'
  tab.queryLimit = 50
  
  // Set default JS script
  tab.queryPreset = 'Basic Fetch'
  tab.queryScript = `// Query with JavaScript using the Firebase Admin SDK\n// See examples at https://firefoo.app/go/firestore-js-query\nasync function run() {\n  const query = await db.collection("${colName}")\n    .limit(50)\n    .get();\n  return query;\n}`

  fetchDocuments(tab)
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
        structuredQuery.order = [{
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
          if (r.document) {
            docs.push(r.document)
          }
        })
      }
      tab.documents = docs
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
      collections: [],
      activeCollection: parentTab.activeCollection,
      activeDocPath: docPath,
      documents: [doc],
      nextDocPageToken: '',
      loadingCollections: false,
      loadingDocuments: false,
      firestoreView: 'explorer',
      resultsView: 'tree',
      selectedDocNames: [],
      queryMode: 'js',
      queryScript: `// Query with JavaScript using the Firebase Admin SDK\n// See examples at https://firefoo.app/go/firestore-js-query\nasync function run() {\n  const doc = await db.doc("${docPath}").get();\n  return doc;\n}`,
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
    alert('Copied to clipboard!')
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
}

// Convert typed Firestore fields to clean JS Object recursively
function cleanDocFields(fields) {
  const res = {}
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

// Convert clean JS Object back to typed Firestore fields recursively
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
    if (cur[part] === undefined) {
      cur[part] = {}
    }
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
    fetchDocuments(tab)
  } catch (err) {
    alert('Failed to save field change: ' + err)
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

function openAddDocModal(tab) {
  activeTabContext.value = tab
  newDocId.value = ''
  newDocFieldsJson.value = '{\n  \n}'
  showAddDocModal.value = true
}

async function saveNewDoc() {
  const tab = activeTabContext.value
  if (!tab) return
  
  let fields = {}
  try {
    fields = JSON.parse(newDocFieldsJson.value.trim() || '{}')
  } catch (e) {
    alert('Invalid JSON fields: ' + e)
    return
  }

  // Path is collection + id
  const path = tab.activeCollection + '/' + (newDocId.value.trim() || generateId())
  try {
    await SaveDocument(tab.email, tab.projectId, '', path, fields, true)
    showAddDocModal.value = false
    fetchDocuments(tab)
  } catch (err) {
    alert('Failed to save document: ' + err)
  }
}

function editDocumentJSON(tab, doc) {
  activeTabContext.value = tab
  const docId = getDocId(doc.name)
  activeDocPath.value = tab.activeCollection + '/' + docId
  activeDocJson.value = JSON.stringify(doc.fields || {}, null, 2)
  showDocModal.value = true
}

async function saveDocFromJSON() {
  const tab = activeTabContext.value
  if (!tab) return
  let fields = {}
  try {
    fields = JSON.parse(activeDocJson.value.trim())
  } catch (err) {
    alert('Invalid JSON data format')
    return
  }

  try {
    await SaveDocument(tab.email, tab.projectId, '', activeDocPath.value, fields, false)
    showDocModal.value = false
    fetchDocuments(tab)
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
      fetchDocuments(tab)
    } catch (err) {
      alert('Failed to delete document: ' + err)
    }
  }
}


function loadQueryExample(tab) {
  tab.queryScript = `// Example script to batch update docs\nconst snapshot = await db.collection("${tab.activeCollection || 'users'}").limit(10).get();\nconsole.log(\`Fetched \${snapshot.size} documents\`);\n\nsnapshot.forEach(doc => {\n  console.log(\`Doc ID: \${doc.id}\`);\n});`
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
    // 1. Conflict Check
    const check = await CheckConflict(copyWizard.destEmail, copyWizard.destProjectId, '', copyWizard.path, copyWizard.type)
    if (check.conflict && !copyWizard.overwrite) {
      const proceed = confirm(`${check.message} Do you want to proceed and overwrite?`)
      if (!proceed) {
        copyingNode.value = false
        return
      }
      copyWizard.overwrite = true
    }

    // 2. Execute Copy
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

function getFieldsPreview(fields) {
  if (!fields) return '{}'
  const keys = Object.keys(fields)
  if (keys.length === 0) return '{}'
  const summary = keys.map(k => {
    let val = fields[k]
    if (typeof val === 'object' && val !== null) {
      val = JSON.stringify(val)
    }
    return `${k}: ${val}`
  }).join(', ')
  return `{ ${summary.substring(0, 100)}${summary.length > 100 ? '...' : ''} }`
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
</script>

<style>
/* Custom styled scrollbars */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}
::-webkit-scrollbar-track {
  background: #0c0d0e;
}
::-webkit-scrollbar-thumb {
  background: #2a2d31;
  border-radius: 3px;
}
::-webkit-scrollbar-thumb:hover {
  background: #3a3f45;
}
</style>
