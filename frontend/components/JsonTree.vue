<template>
  <div class="select-text font-mono text-xs">
    <!-- Row Element -->
    <div 
      class="flex items-center hover:bg-white/4 py-1.5 px-2 rounded-lg transition-colors group relative border border-transparent hover:border-white/6"
      :style="{ paddingLeft: `${depth * 18}px` }"
    >
      <!-- Expand Arrow -->
      <span 
        v-if="isExpandable || isRoot" 
        class="w-5 h-5 flex items-center justify-center cursor-pointer text-zinc-400 hover:text-white mr-1.5 transform transition-transform duration-150 text-xs shrink-0 rounded hover:bg-white/10"
        :class="isOpen ? 'rotate-90 text-amber-400' : ''"
        @click.stop="toggleExpand"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M7.21 14.77a.75.75 0 01.02-1.06L11.168 10 7.23 6.29a.75.75 0 111.04-1.08l4.5 4.25a.75.75 0 010 1.08l-4.5 4.25a.75.75 0 01-1.06-.02z" clip-rule="evenodd" />
        </svg>
      </span>
      <span v-else class="w-5 mr-1.5 shrink-0"></span>

      <!-- Key Name -->
      <span 
        class="font-semibold text-amber-400 mr-2 shrink-0 select-all cursor-pointer hover:underline"
        :class="{ 'text-amber-300 font-bold text-sm tracking-wide': isRoot }"
        @click.stop="isExpandable ? toggleExpand() : startEdit()"
        :title="label"
      >
        {{ label }}:
      </span>

      <!-- Type Badge -->
      <span 
        class="text-[9px] uppercase tracking-wider font-semibold px-1.5 py-0.5 rounded mr-2.5 shrink-0 border"
        :class="badgeClass"
      >
        {{ isRoot ? 'document' : displayType }}
      </span>

      <!-- Value Display / Editor -->
      <div class="flex-1 flex items-center min-w-0">
        <!-- Edit Mode -->
        <div v-if="isEditing" class="flex items-center gap-1.5 w-full max-w-lg" @click.stop>
          <select 
            v-model="editType" 
            @change="onEditTypeChange"
            class="bg-[#090a0f] border border-white/15 text-zinc-300 rounded px-2 py-1 text-[11px] font-mono outline-none cursor-pointer focus:border-amber-500"
          >
            <option value="string">string</option>
            <option value="number">number</option>
            <option value="boolean">boolean</option>
            <option value="timestamp">timestamp</option>
            <option value="reference">reference</option>
            <option value="geopoint">geopoint</option>
            <option value="bytes">bytes</option>
            <option value="null">null</option>
          </select>

          <input 
            v-if="editType === 'boolean'"
            type="checkbox"
            v-model="editValue"
            class="accent-amber-500 w-4 h-4 cursor-pointer"
          />
          <input 
            v-else-if="editType === 'number'"
            type="number"
            v-model="editValue"
            class="bg-[#090a0f] border border-amber-500/80 text-white rounded px-2 py-0.5 text-xs font-mono outline-none w-32 focus:ring-1 focus:ring-amber-500"
            @keydown.enter.prevent.stop="saveEdit"
            @keydown.esc.prevent.stop="cancelEdit"
            ref="inputRef"
          />
          <input 
            v-else
            v-model="editValue" 
            class="bg-[#090a0f] border border-amber-500/80 text-white rounded px-2 py-0.5 text-xs font-mono outline-none flex-1 focus:ring-1 focus:ring-amber-500"
            @keydown.enter.prevent.stop="saveEdit"
            @keydown.esc.prevent.stop="cancelEdit"
            ref="inputRef"
          />
          <UButton size="xs" color="success" icon="i-heroicons-check" type="button" @click.stop.prevent="saveEdit" />
          <UButton size="xs" color="neutral" variant="ghost" icon="i-heroicons-x-mark" type="button" @click.stop.prevent="cancelEdit" />
        </div>

        <!-- View Mode -->
        <div v-else class="flex items-center group/val gap-1.5 min-w-0 flex-1">
          <!-- Primitive Value -->
          <span 
            v-if="!isExpandable && !isRoot"
            class="cursor-pointer hover:underline truncate select-all font-medium transition-colors"
            :class="valueColorClass"
            :title="String(value)"
            @dblclick.stop="startEdit"
          >
            {{ formatValue(value) }}
          </span>

          <!-- Expandable Summary -->
          <span 
            v-else 
            class="text-zinc-500 select-none italic text-[11px] truncate flex-1 min-w-0 pr-2 block"
            :title="summaryTooltip"
          >
            {{ isRoot ? (subcollectionsFetched && subcollections.length > 0 ? `{${subcollections.length} subcolls}` : valueSummary) : valueSummary }}
          </span>

          <!-- Copy Value Helper -->
          <button
            v-if="!isExpandable && !isRoot"
            type="button"
            class="opacity-0 group-hover/val:opacity-100 text-zinc-500 hover:text-amber-400 p-0.5 rounded transition-opacity shrink-0"
            @click.stop="copyValueToClipboard"
            title="Copy value"
          >
            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
          </button>

          <!-- Inline Edit Button (only for primitives) -->
          <UButton 
            v-if="!isExpandable && !isRoot"
            icon="i-heroicons-pencil"
            size="xs"
            color="neutral"
            variant="ghost"
            class="opacity-0 group-hover/val:opacity-100 p-0.5 scale-90 shrink-0"
            type="button"
            @click.stop.prevent="startEdit"
            title="Edit value"
          />
        </div>
      </div>

      <!-- Action Toolbar -->
      <div class="ml-auto opacity-0 group-hover:opacity-100 flex items-center gap-1 transition-opacity shrink-0 z-10 pl-2 bg-[#141721] rounded-md shadow-sm border border-white/5">
        <!-- Add Field Button (For root or map objects) -->
        <UButton 
          v-if="isRoot || displayType === 'map'"
          size="xs" 
          color="warning" 
          variant="subtle" 
          icon="i-heroicons-plus"
          type="button" 
          @click.stop="startAddField"
          title="Add new field"
          class="text-[10px] px-1.5 py-0.5 font-medium"
        >
          Field
        </UButton>

        <!-- Root Toolbar Actions -->
        <template v-if="isRoot">
          <UButton 
            v-if="subcollectionsFetched && subcollections.length > 0"
            size="xs" 
            color="warning" 
            variant="subtle" 
            type="button" 
            @click.stop="toggleExpand"
            class="text-[10px] px-1.5 py-0.5"
            title="Toggle subcollections"
          >
            📁 {{ subcollections.length }}
          </UButton>
          <UButton 
            size="xs" 
            color="neutral" 
            variant="subtle" 
            icon="i-heroicons-arrow-top-right-on-square"
            type="button" 
            @click.stop="$emit('openTab')"
            title="Open document in dedicated tab"
            class="text-[10px] px-1.5 py-0.5"
          >
            Tab
          </UButton>
          <UButton 
            size="xs" 
            color="neutral" 
            variant="subtle" 
            icon="i-heroicons-document-duplicate"
            type="button" 
            @click.stop="$emit('duplicateDoc')"
            title="Duplicate document"
            class="text-[10px] px-1.5 py-0.5"
          >
            Dup
          </UButton>
          <UButton 
            size="xs" 
            color="neutral" 
            variant="subtle" 
            icon="i-heroicons-code-bracket"
            type="button" 
            @click.stop="$emit('editJson')"
            title="Edit raw JSON"
            class="text-[10px] px-1.5 py-0.5"
          >
            JSON
          </UButton>
          <UButton 
            size="xs" 
            color="danger" 
            variant="subtle" 
            icon="i-heroicons-trash" 
            type="button" 
            @click.stop="$emit('deleteDoc')"
            title="Delete document"
            class="p-1"
          />
        </template>

        <!-- Non-root Delete Field Action -->
        <template v-else>
          <UButton 
            size="xs" 
            color="danger" 
            variant="ghost" 
            icon="i-heroicons-trash"
            type="button" 
            @click.stop="confirmDeleteField"
            title="Delete field"
            class="p-1"
          />
        </template>
      </div>
    </div>

    <!-- Add Field Inline Form -->
    <div 
      v-if="isAddingField" 
      class="mt-1 mb-2 p-3 bg-[#141721] border border-amber-500/30 rounded-lg flex flex-wrap items-center gap-2 shadow-lg"
      :style="{ marginLeft: `${(depth + 1) * 18}px` }"
      @click.stop
    >
      <span class="text-xs font-semibold text-amber-400 shrink-0 flex items-center gap-1">
        <span>+</span> Add Field:
      </span>
      
      <input 
        v-model="newFieldName" 
        placeholder="Field Name" 
        class="bg-[#090a0f] border border-white/15 text-white rounded px-2 py-1 text-xs font-mono outline-none w-36 focus:border-amber-500"
        @keydown.enter.prevent.stop="saveAddField"
        @keydown.esc.prevent.stop="cancelAddField"
        ref="newFieldNameRef"
      />

      <select 
        v-model="newFieldType" 
        @change="onNewFieldTypeChange"
        class="bg-[#090a0f] border border-white/15 text-zinc-300 rounded px-2 py-1 text-xs font-mono outline-none cursor-pointer focus:border-amber-500"
      >
        <option value="string">string</option>
        <option value="number">number</option>
        <option value="boolean">boolean</option>
        <option value="map">map</option>
        <option value="array">array</option>
        <option value="null">null</option>
        <option value="timestamp">timestamp</option>
        <option value="reference">reference</option>
        <option value="geopoint">geopoint</option>
      </select>

      <input 
        v-if="newFieldType === 'string' || newFieldType === 'reference'"
        v-model="newFieldValue"
        :placeholder="newFieldType === 'reference' ? 'collection/docId' : 'Value'"
        class="bg-[#090a0f] border border-white/15 text-white rounded px-2 py-1 text-xs font-mono outline-none flex-1 min-w-30 focus:border-amber-500"
        @keydown.enter.prevent.stop="saveAddField"
        @keydown.esc.prevent.stop="cancelAddField"
      />
      <input 
        v-else-if="newFieldType === 'geopoint'"
        v-model="newFieldValue"
        placeholder="lat, lng (e.g. 24.7136, 46.6753)"
        class="bg-[#090a0f] border border-white/15 text-white rounded px-2 py-1 text-xs font-mono outline-none flex-1 min-w-36 focus:border-amber-500"
        @keydown.enter.prevent.stop="saveAddField"
        @keydown.esc.prevent.stop="cancelAddField"
      />
      <input 
        v-else-if="newFieldType === 'number'"
        type="number"
        v-model="newFieldValue"
        placeholder="0"
        class="bg-[#090a0f] border border-white/15 text-white rounded px-2 py-1 text-xs font-mono outline-none w-28 focus:border-amber-500"
        @keydown.enter.prevent.stop="saveAddField"
        @keydown.esc.prevent.stop="cancelAddField"
      />
      <select 
        v-else-if="newFieldType === 'boolean'"
        v-model="newFieldValue"
        class="bg-[#090a0f] border border-white/15 text-white rounded px-2 py-1 text-xs font-mono outline-none focus:border-amber-500"
      >
        <option :value="true">true</option>
        <option :value="false">false</option>
      </select>
      <input 
        v-else-if="newFieldType === 'timestamp'"
        v-model="newFieldValue"
        placeholder="ISO Timestamp"
        class="bg-[#090a0f] border border-white/15 text-white rounded px-2 py-1 text-xs font-mono outline-none flex-1 min-w-40 focus:border-amber-500"
        @keydown.enter.prevent.stop="saveAddField"
        @keydown.esc.prevent.stop="cancelAddField"
      />
      <span v-else-if="newFieldType === 'map'" class="text-[11px] text-zinc-500 italic">{ } (empty map)</span>
      <span v-else-if="newFieldType === 'array'" class="text-[11px] text-zinc-500 italic">[ ] (empty array)</span>
      <span v-else-if="newFieldType === 'null'" class="text-[11px] text-zinc-500 italic">null</span>

      <div class="flex items-center gap-1.5 ml-auto">
        <UButton size="xs" color="warning" type="button" @click.stop.prevent="saveAddField">Add</UButton>
        <UButton size="xs" color="neutral" variant="ghost" type="button" @click.stop.prevent="cancelAddField">Cancel</UButton>
      </div>
    </div>

    <!-- Recursive Children & Subcollections -->
    <div v-if="(isExpandable || isRoot) && isOpen" class="mt-0.5">
      <!-- Object/Map fields -->
      <template v-if="displayType === 'map'">
        <JsonTree 
          v-for="(val, key) in value" 
          :key="key"
          :label="key"
          :value="val"
          :depth="depth + 1"
          @change="(path, newVal) => $emit('change', isRoot ? path : `${label}/${path}`, newVal)"
          @add-field="(path, fieldName, fieldType, initialVal) => $emit('addField', isRoot ? path : `${label}/${path}`, fieldName, fieldType, initialVal)"
          @delete-field="(path) => $emit('deleteField', isRoot ? path : `${label}/${path}`)"
        />
      </template>

      <!-- Array items -->
      <template v-else-if="displayType === 'array'">
        <JsonTree 
          v-for="(val, index) in value" 
          :key="index"
          :label="String(index)"
          :value="val"
          :depth="depth + 1"
          @change="(path, newVal) => $emit('change', isRoot ? path : `${label}/${path}`, newVal)"
          @add-field="(path, fieldName, fieldType, initialVal) => $emit('addField', isRoot ? path : `${label}/${path}`, fieldName, fieldType, initialVal)"
          @delete-field="(path) => $emit('deleteField', isRoot ? path : `${label}/${path}`)"
        />
      </template>

      <!-- Subcollections Section for Document -->
      <div v-if="isRoot" class="mt-2 pt-2 border-t border-white/8">
        <div 
          class="flex items-center justify-between hover:bg-white/4 py-1 px-2 rounded transition-colors text-xs font-semibold text-amber-400 select-none"
          :style="{ paddingLeft: `${(depth + 1) * 18}px` }"
        >
          <div class="flex items-center gap-1.5">
            <span class="text-amber-400">📁</span>
            <span>Subcollections</span>
            <span v-if="subcollectionsFetched" class="text-[10px] text-amber-300 bg-amber-500/15 px-1.5 py-0.5 rounded font-mono border border-amber-500/25">
              {{ subcollections.length }}
            </span>
            <span v-if="loadingSubcollections" class="text-[10px] text-amber-400/80 animate-pulse italic ml-1">
              Loading...
            </span>
          </div>
          <UButton 
            size="xs" 
            color="neutral" 
            variant="ghost" 
            icon="i-heroicons-arrow-path" 
            type="button"
            class="p-0.5"
            @click.stop="$emit('fetchSubcollections')" 
            title="Refresh subcollections" 
          />
        </div>

        <div v-if="subcollectionsFetched && subcollections.length === 0" 
             class="text-[11px] text-zinc-500 italic py-1 px-2"
             :style="{ paddingLeft: `${(depth + 2) * 18}px` }">
          No subcollections found
        </div>

        <div v-else-if="subcollections && subcollections.length > 0" class="space-y-1 mt-1">
          <div 
            v-for="subCol in subcollections" 
            :key="subCol"
            class="flex items-center justify-between hover:bg-white/4 py-1.5 px-2 rounded-lg transition-colors group text-xs cursor-pointer border border-transparent hover:border-amber-500/20"
            :style="{ paddingLeft: `${(depth + 2) * 18}px` }"
            @click.stop="$emit('openSubcollection', subCol)"
          >
            <div class="flex items-center gap-2 min-w-0">
              <span class="text-amber-400">📁</span>
              <span class="font-semibold text-amber-300 font-mono truncate">{{ subCol }}</span>
              <span class="text-[9px] px-1.5 py-0.5 rounded bg-amber-500/10 text-amber-400 border border-amber-500/20">subcollection</span>
            </div>
            <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity shrink-0">
              <UButton 
                size="xs" 
                color="warning" 
                variant="ghost" 
                type="button" 
                @click.stop="$emit('openSubcollection', subCol)"
              >
                Browse
              </UButton>
              <UButton 
                size="xs" 
                color="neutral" 
                variant="ghost" 
                type="button" 
                @click.stop="$emit('openSubcollectionInNewTab', subCol)"
              >
                New Tab
              </UButton>
            </div>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, nextTick } from 'vue'

const props = defineProps({
  label: { type: String, required: true },
  value: { type: [String, Number, Boolean, Object, Array, null], default: null },
  depth: { type: Number, default: 0 },
  isRoot: { type: Boolean, default: false },
  subcollections: { type: Array, default: () => [] },
  loadingSubcollections: { type: Boolean, default: false },
  subcollectionsFetched: { type: Boolean, default: false }
})

const emit = defineEmits([
  'change', 'addField', 'deleteField', 'editJson', 'deleteDoc', 'duplicateDoc', 'openTab',
  'fetchSubcollections', 'openSubcollection', 'openSubcollectionInNewTab'
])

const isOpen = ref(false)
const isEditing = ref(false)

function toggleExpand() {
  isOpen.value = !isOpen.value
  if (props.isRoot && isOpen.value && !props.subcollectionsFetched && !props.loadingSubcollections) {
    emit('fetchSubcollections')
  }
}
const editType = ref('string')
const editValue = ref('')
const inputRef = ref(null)

const isAddingField = ref(false)
const newFieldName = ref('')
const newFieldType = ref('string')
const newFieldValue = ref('')
const newFieldNameRef = ref(null)

const displayType = computed(() => {
  if (props.value === null) return 'null'
  if (Array.isArray(props.value)) return 'array'
  if (typeof props.value === 'object') {
    if (props.value && '__time__' in props.value) return 'timestamp'
    if (props.value && '__geopoint__' in props.value) return 'geopoint'
    if (props.value && '__ref__' in props.value) return 'reference'
    if (props.value && '__bytes__' in props.value) return 'bytes'
    return 'map'
  }
  return typeof props.value
})

const isExpandable = computed(() => {
  return displayType.value === 'array' || displayType.value === 'map'
})

const badgeClass = computed(() => {
  if (props.isRoot) return 'bg-amber-500/15 text-amber-300 border-amber-500/30'
  switch (displayType.value) {
    case 'string': return 'bg-emerald-500/10 text-emerald-400 border-emerald-500/25'
    case 'number': return 'bg-amber-500/10 text-amber-400 border-amber-500/25'
    case 'boolean': return 'bg-purple-500/10 text-purple-400 border-purple-500/25'
    case 'timestamp': return 'bg-cyan-500/10 text-cyan-400 border-cyan-500/25'
    case 'geopoint': return 'bg-teal-500/10 text-teal-400 border-teal-500/25'
    case 'reference': return 'bg-pink-500/10 text-pink-400 border-pink-500/25'
    case 'bytes': return 'bg-slate-500/10 text-slate-300 border-slate-500/25'
    case 'map': return 'bg-sky-500/10 text-sky-400 border-sky-500/25'
    case 'array': return 'bg-indigo-500/10 text-indigo-400 border-indigo-500/25'
    case 'null': return 'bg-zinc-500/10 text-zinc-400 border-zinc-500/25'
    default: return 'bg-zinc-800 text-zinc-400 border-zinc-700'
  }
})

const valueColorClass = computed(() => {
  switch (displayType.value) {
    case 'string': return 'text-emerald-300'
    case 'number': return 'text-amber-300'
    case 'boolean': return 'text-purple-300 font-semibold'
    case 'timestamp': return 'text-cyan-300'
    case 'geopoint': return 'text-teal-300'
    case 'reference': return 'text-pink-300'
    case 'bytes': return 'text-slate-300'
    case 'null': return 'text-zinc-500 italic'
    default: return 'text-zinc-200'
  }
})

const valueSummary = computed(() => {
  if (displayType.value === 'array') {
    return `[${props.value.length} items]`
  }
  if (displayType.value === 'map') {
    const keys = Object.keys(props.value || {})
    if (keys.length === 0) return '{ empty }'
    const preview = keys.slice(0, 16).join(', ')
    return `{ ${preview}${keys.length > 16 ? ', ...' : ''} }`
  }
  return ''
})

const summaryTooltip = computed(() => {
  if (displayType.value === 'map') {
    const keys = Object.keys(props.value || {})
    return `{ ${keys.join(', ')} } (${keys.length} fields)`
  }
  if (displayType.value === 'array') {
    return `[${props.value.length} items]`
  }
  return ''
})

function formatValue(val) {
  if (val === null) return 'null'
  if (typeof val === 'string') return `"${val}"`
  if (displayType.value === 'timestamp') {
    return new Date(val.__time__).toLocaleString()
  }
  if (displayType.value === 'geopoint') {
    const gp = val.__geopoint__ || {}
    return `(${gp.latitude ?? 0}, ${gp.longitude ?? 0})`
  }
  if (displayType.value === 'reference') {
    return `→ ${val.__ref__ || ''}`
  }
  if (displayType.value === 'bytes') {
    const b = String(val.__bytes__ || '')
    return `[bytes: ${b.length > 24 ? b.slice(0, 24) + '...' : b}]`
  }
  return String(val)
}

function copyValueToClipboard() {
  let str = ''
  if (displayType.value === 'timestamp') str = props.value.__time__
  else if (displayType.value === 'geopoint') {
    const gp = props.value.__geopoint__ || {}
    str = `${gp.latitude ?? 0}, ${gp.longitude ?? 0}`
  } else if (displayType.value === 'reference') str = props.value.__ref__ || ''
  else if (displayType.value === 'bytes') str = props.value.__bytes__ || ''
  else str = props.value === null ? 'null' : String(props.value)
  navigator.clipboard.writeText(str)
}

function startEdit() {
  if (isExpandable.value) return
  isEditing.value = true
  editType.value = displayType.value
  if (displayType.value === 'timestamp') {
    editValue.value = props.value.__time__
  } else if (displayType.value === 'geopoint') {
    const gp = props.value.__geopoint__ || {}
    editValue.value = `${gp.latitude ?? 0}, ${gp.longitude ?? 0}`
  } else if (displayType.value === 'reference') {
    editValue.value = props.value.__ref__ || ''
  } else if (displayType.value === 'bytes') {
    editValue.value = props.value.__bytes__ || ''
  } else {
    editValue.value = displayType.value === 'null' ? '' : props.value
  }
  nextTick(() => {
    if (inputRef.value) {
      inputRef.value.focus()
      if (typeof editValue.value === 'string') {
        inputRef.value.select()
      }
    }
  })
}

function onEditTypeChange() {
  if (editType.value === 'string' || editType.value === 'reference' || editType.value === 'bytes') {
    editValue.value = ''
  } else if (editType.value === 'geopoint') {
    editValue.value = '0, 0'
  } else if (editType.value === 'number') {
    editValue.value = 0
  } else if (editType.value === 'boolean') {
    editValue.value = true
  } else if (editType.value === 'timestamp') {
    editValue.value = new Date().toISOString()
  } else if (editType.value === 'null') {
    editValue.value = null
  }
}

function saveEdit() {
  let parsedVal = editValue.value
  if (editType.value === 'number') {
    parsedVal = Number(editValue.value)
    if (isNaN(parsedVal)) {
      alert('Must be a valid number')
      return
    }
  } else if (editType.value === 'boolean') {
    parsedVal = Boolean(editValue.value)
  } else if (editType.value === 'null') {
    parsedVal = null
  } else if (editType.value === 'timestamp') {
    parsedVal = { __time__: editValue.value }
  } else if (editType.value === 'geopoint') {
    const parts = String(editValue.value).split(',').map(s => Number(s.trim()))
    if (parts.length !== 2 || isNaN(parts[0]) || isNaN(parts[1])) {
      alert('Geopoint must be "latitude, longitude"')
      return
    }
    parsedVal = { __geopoint__: { latitude: parts[0], longitude: parts[1] } }
  } else if (editType.value === 'reference') {
    parsedVal = { __ref__: String(editValue.value).trim() }
  } else if (editType.value === 'bytes') {
    parsedVal = { __bytes__: String(editValue.value).trim() }
  } else if (editType.value === 'string') {
    parsedVal = String(editValue.value)
  }

  isEditing.value = false
  emit('change', props.label, parsedVal)
}

function cancelEdit() {
  isEditing.value = false
}

function startAddField() {
  isOpen.value = true
  isAddingField.value = true
  newFieldName.value = ''
  newFieldType.value = 'string'
  newFieldValue.value = ''
  nextTick(() => {
    if (newFieldNameRef.value) {
      newFieldNameRef.value.focus()
    }
  })
}

function onNewFieldTypeChange() {
  if (newFieldType.value === 'string' || newFieldType.value === 'reference') {
    newFieldValue.value = ''
  } else if (newFieldType.value === 'geopoint') {
    newFieldValue.value = '0, 0'
  } else if (newFieldType.value === 'number') {
    newFieldValue.value = 0
  } else if (newFieldType.value === 'boolean') {
    newFieldValue.value = true
  } else if (newFieldType.value === 'timestamp') {
    newFieldValue.value = new Date().toISOString()
  } else {
    newFieldValue.value = ''
  }
}

function saveAddField() {
  const name = newFieldName.value.trim()
  if (!name) {
    alert('Field name is required')
    return
  }

  let finalVal = null
  if (newFieldType.value === 'string') {
    finalVal = String(newFieldValue.value)
  } else if (newFieldType.value === 'number') {
    finalVal = Number(newFieldValue.value)
    if (isNaN(finalVal)) {
      alert('Must be a valid number')
      return
    }
  } else if (newFieldType.value === 'boolean') {
    finalVal = Boolean(newFieldValue.value)
  } else if (newFieldType.value === 'map') {
    finalVal = {}
  } else if (newFieldType.value === 'array') {
    finalVal = []
  } else if (newFieldType.value === 'null') {
    finalVal = null
  } else if (newFieldType.value === 'timestamp') {
    finalVal = { __time__: newFieldValue.value }
  } else if (newFieldType.value === 'reference') {
    finalVal = { __ref__: String(newFieldValue.value).trim() }
  } else if (newFieldType.value === 'geopoint') {
    const parts = String(newFieldValue.value).split(',').map(s => Number(s.trim()))
    if (parts.length !== 2 || isNaN(parts[0]) || isNaN(parts[1])) {
      alert('Geopoint must be "latitude, longitude"')
      return
    }
    finalVal = { __geopoint__: { latitude: parts[0], longitude: parts[1] } }
  }

  isAddingField.value = false
  emit('addField', props.isRoot ? '' : props.label, name, newFieldType.value, finalVal)
}

function cancelAddField() {
  isAddingField.value = false
}

function confirmDeleteField() {
  emit('deleteField', props.label)
}
</script>
