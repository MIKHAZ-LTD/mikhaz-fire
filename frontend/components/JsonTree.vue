<template>
  <div class="select-text font-mono text-xs">
    <!-- Row Element -->
    <div 
      class="flex items-center hover:bg-[#1d1f22] py-1 px-2 rounded transition-colors group relative"
      :style="{ paddingLeft: `${depth * 16}px` }"
    >
      <!-- Expand Arrow -->
      <span 
        v-if="isExpandable" 
        class="w-5 h-5 flex items-center justify-center cursor-pointer text-[#606870] hover:text-white mr-1 transform transition-transform text-sm shrink-0"
        :class="isOpen ? 'rotate-90' : ''"
        @click.stop="isOpen = !isOpen"
      >
        &#9656;
      </span>
      <span v-else class="w-5 mr-1 shrink-0"></span>

      <!-- Key Name -->
      <span 
        class="font-semibold text-amber-500 mr-2 shrink-0 select-all"
        :class="{ 'text-amber-400 font-bold text-sm': isRoot }"
      >
        {{ label }}:
      </span>

      <!-- Type Badge -->
      <span class="text-[9px] px-1 py-0.5 rounded bg-[#1d1f22] text-[#9da5ac] mr-3 shrink-0 scale-90 border border-[#2a2d31]">
        {{ displayType }}
      </span>

      <!-- Value Display / Editor -->
      <div class="flex-1 flex items-center min-w-0">
        <!-- Edit Mode -->
        <div v-if="isEditing" class="flex items-center gap-1.5 w-full max-w-lg" @click.stop>
          <select 
            v-model="editType" 
            @change="onEditTypeChange"
            class="bg-[#0c0d0e] border border-[#2a2d31] text-[#9da5ac] rounded px-1 py-0.5 text-[10px] font-mono outline-none cursor-pointer"
          >
            <option value="string">string</option>
            <option value="number">number</option>
            <option value="boolean">boolean</option>
            <option value="timestamp">timestamp</option>
            <option value="null">null</option>
          </select>

          <input 
            v-if="editType === 'boolean'"
            type="checkbox"
            v-model="editValue"
            class="accent-amber-500"
          />
          <input 
            v-else-if="editType === 'number'"
            type="number"
            v-model="editValue"
            class="bg-[#0c0d0e] border border-amber-500 text-white rounded px-1.5 py-0.5 text-xs font-mono outline-none w-28"
            @keydown.enter.prevent.stop="saveEdit"
            @keydown.esc.prevent.stop="cancelEdit"
            ref="inputRef"
          />
          <input 
            v-else
            v-model="editValue" 
            class="bg-[#0c0d0e] border border-amber-500 text-white rounded px-1.5 py-0.5 text-xs font-mono outline-none flex-1"
            @keydown.enter.prevent.stop="saveEdit"
            @keydown.esc.prevent.stop="cancelEdit"
            ref="inputRef"
          />
          <UButton size="xs" color="success" icon="i-heroicons-check" type="button" @click.stop.prevent="saveEdit" />
          <UButton size="xs" color="danger" icon="i-heroicons-x-mark" type="button" @click.stop.prevent="cancelEdit" />
        </div>

        <!-- View Mode -->
        <div v-else class="flex items-center group/val truncate">
          <!-- Primitive Value -->
          <span 
            v-if="!isExpandable"
            class="text-[#e3e6e8] cursor-pointer hover:underline truncate select-all"
            :title="String(value)"
            @dblclick.stop="startEdit"
          >
            {{ formatValue(value) }}
          </span>

          <!-- Expandable Summary -->
          <span v-else class="text-[#606870] select-none italic text-[11px]">
            {{ valueSummary }}
          </span>

          <!-- Inline Edit Button (only for primitives) -->
          <UButton 
            v-if="!isExpandable"
            icon="i-heroicons-pencil"
            size="xs"
            color="neutral"
            variant="ghost"
            class="opacity-0 group-hover/val:opacity-100 ml-2 scale-90 p-0"
            type="button"
            @click.stop.prevent="startEdit"
            title="Double-click or click to edit"
          />
        </div>
      </div>

      <!-- Action Toolbar -->
      <div class="ml-auto opacity-0 group-hover:opacity-100 flex items-center gap-1 transition-opacity shrink-0">
        <!-- Add Field Button (For root or map objects) -->
        <UButton 
          v-if="isRoot || displayType === 'map'"
          size="xs" 
          color="warning" 
          variant="ghost" 
          icon="i-heroicons-plus"
          type="button" 
          @click.stop="startAddField"
          title="Add new field"
        >
          + Field
        </UButton>

        <!-- Root Toolbar Actions -->
        <template v-if="isRoot">
          <UButton size="xs" color="neutral" variant="ghost" type="button" @click.stop="$emit('openTab')">Open Tab</UButton>
          <UButton size="xs" color="neutral" variant="ghost" type="button" @click.stop="$emit('editJson')">Edit JSON</UButton>
          <UButton size="xs" color="danger" variant="ghost" type="button" @click.stop="$emit('deleteDoc')">Delete</UButton>
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
          />
        </template>
      </div>
    </div>

    <!-- Add Field Inline Form -->
    <div 
      v-if="isAddingField" 
      class="mt-1 mb-2 p-2 bg-[#121315] border border-amber-500/40 rounded flex flex-wrap items-center gap-2"
      :style="{ marginLeft: `${(depth + 1) * 16}px` }"
      @click.stop
    >
      <span class="text-xs font-semibold text-amber-400 shrink-0">+ Add Field:</span>
      
      <input 
        v-model="newFieldName" 
        placeholder="Field Name" 
        class="bg-[#0c0d0e] border border-[#2a2d31] text-white rounded px-2 py-1 text-xs font-mono outline-none w-32 focus:border-amber-500"
        @keydown.enter.prevent.stop="saveAddField"
        @keydown.esc.prevent.stop="cancelAddField"
        ref="newFieldNameRef"
      />

      <select 
        v-model="newFieldType" 
        @change="onNewFieldTypeChange"
        class="bg-[#0c0d0e] border border-[#2a2d31] text-white rounded px-2 py-1 text-xs font-mono outline-none cursor-pointer focus:border-amber-500"
      >
        <option value="string">string</option>
        <option value="number">number</option>
        <option value="boolean">boolean</option>
        <option value="map">map</option>
        <option value="array">array</option>
        <option value="null">null</option>
        <option value="timestamp">timestamp</option>
      </select>

      <input 
        v-if="newFieldType === 'string'"
        v-model="newFieldValue"
        placeholder="Value"
        class="bg-[#0c0d0e] border border-[#2a2d31] text-white rounded px-2 py-1 text-xs font-mono outline-none flex-1 min-w-[120px] focus:border-amber-500"
        @keydown.enter.prevent.stop="saveAddField"
        @keydown.esc.prevent.stop="cancelAddField"
      />
      <input 
        v-else-if="newFieldType === 'number'"
        type="number"
        v-model="newFieldValue"
        placeholder="0"
        class="bg-[#0c0d0e] border border-[#2a2d31] text-white rounded px-2 py-1 text-xs font-mono outline-none w-28 focus:border-amber-500"
        @keydown.enter.prevent.stop="saveAddField"
        @keydown.esc.prevent.stop="cancelAddField"
      />
      <select 
        v-else-if="newFieldType === 'boolean'"
        v-model="newFieldValue"
        class="bg-[#0c0d0e] border border-[#2a2d31] text-white rounded px-2 py-1 text-xs font-mono outline-none focus:border-amber-500"
      >
        <option :value="true">true</option>
        <option :value="false">false</option>
      </select>
      <input 
        v-else-if="newFieldType === 'timestamp'"
        v-model="newFieldValue"
        placeholder="ISO Timestamp"
        class="bg-[#0c0d0e] border border-[#2a2d31] text-white rounded px-2 py-1 text-xs font-mono outline-none flex-1 min-w-[160px] focus:border-amber-500"
        @keydown.enter.prevent.stop="saveAddField"
        @keydown.esc.prevent.stop="cancelAddField"
      />
      <span v-else-if="newFieldType === 'map'" class="text-[11px] text-[#606870] italic">{ } (empty map)</span>
      <span v-else-if="newFieldType === 'array'" class="text-[11px] text-[#606870] italic">[ ] (empty array)</span>
      <span v-else-if="newFieldType === 'null'" class="text-[11px] text-[#606870] italic">null</span>

      <div class="flex items-center gap-1.5 ml-auto">
        <UButton size="xs" color="warning" type="button" @click.stop.prevent="saveAddField">Add</UButton>
        <UButton size="xs" color="neutral" variant="ghost" type="button" @click.stop.prevent="cancelAddField">Cancel</UButton>
      </div>
    </div>

    <!-- Recursive Children -->
    <div v-if="isExpandable && isOpen" class="mt-0.5">
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
    </div>

  </div>
</template>

<script setup>
import { ref, computed, nextTick } from 'vue'

const props = defineProps({
  label: { type: String, required: true },
  value: { type: [String, Number, Boolean, Object, Array, null], default: null },
  depth: { type: Number, default: 0 },
  isRoot: { type: Boolean, default: false }
})

const emit = defineEmits(['change', 'addField', 'deleteField', 'editJson', 'deleteDoc', 'openTab'])

const isOpen = ref(false)
const isEditing = ref(false)
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
    if (props.value && props.value.__time__) return 'timestamp'
    return 'map'
  }
  return typeof props.value
})

const isExpandable = computed(() => {
  return displayType.value === 'array' || displayType.value === 'map'
})

const valueSummary = computed(() => {
  if (displayType.value === 'array') {
    return `[${props.value.length} items]`
  }
  if (displayType.value === 'map') {
    const keys = Object.keys(props.value || {})
    return `{ ${keys.slice(0, 3).join(', ')}${keys.length > 3 ? '...' : ''} }`
  }
  return ''
})

function formatValue(val) {
  if (val === null) return 'null'
  if (typeof val === 'string') return `"${val}"`
  if (displayType.value === 'timestamp') {
    return new Date(val.__time__).toLocaleString()
  }
  return String(val)
}

function startEdit() {
  if (isExpandable.value) return
  isEditing.value = true
  editType.value = displayType.value
  if (displayType.value === 'timestamp') {
    editValue.value = props.value.__time__
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
  if (editType.value === 'string') {
    editValue.value = ''
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
  if (newFieldType.value === 'string') {
    newFieldValue.value = ''
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
