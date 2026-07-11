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
        class="w-5 h-5 flex items-center justify-center cursor-pointer text-[#606870] hover:text-white mr-1 transform transition-transform text-sm"
        :class="isOpen ? 'rotate-90' : ''"
        @click.stop="isOpen = !isOpen"
      >
        &#9656;
      </span>
      <span v-else class="w-5 mr-1"></span>

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
        <div v-if="isEditing" class="flex items-center gap-1.5 w-full max-w-md">
          <input 
            v-if="displayType === 'boolean'"
            type="checkbox"
            v-model="editValue"
            class="accent-amber-500"
          />
          <input 
            v-else
            v-model="editValue" 
            class="bg-[#0c0d0e] border border-amber-500 text-white rounded px-1.5 py-0.5 text-xs font-mono outline-none flex-1"
            @keydown.enter="saveEdit"
            @keydown.esc="cancelEdit"
            ref="inputRef"
          />
          <UButton size="xs" color="success" icon="i-heroicons-check" @click="saveEdit" />
          <UButton size="xs" color="danger" icon="i-heroicons-x-mark" @click="cancelEdit" />
        </div>

        <!-- View Mode -->
        <div v-else class="flex items-center group/val truncate">
          <!-- Primitive Value -->
          <span 
            v-if="!isExpandable"
            class="text-[#e3e6e8] cursor-pointer hover:underline truncate select-all"
            :title="String(value)"
            @dblclick="startEdit"
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
            @click="startEdit"
            title="Double-click or click to edit"
          />
        </div>
      </div>

      <!-- Root Action Toolbar (Delete/Edit JSON) -->
      <div v-if="isRoot" class="absolute right-2 top-1/2 -translate-y-1/2 opacity-0 group-hover:opacity-100 flex gap-1 transition-opacity">
        <UButton size="xs" color="neutral" variant="ghost" @click.stop="$emit('openTab')">Open Tab</UButton>
        <UButton size="xs" color="neutral" variant="ghost" @click.stop="$emit('editJson')">Edit JSON</UButton>
        <UButton size="xs" color="danger" variant="ghost" @click.stop="$emit('deleteDoc')">Delete</UButton>
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
          @change="(path, newVal) => $emit('change', `${label}/${path}`, newVal)"
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
          @change="(path, newVal) => $emit('change', `${label}/${path}`, newVal)"
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

const emit = defineEmits(['change', 'editJson', 'deleteDoc', 'openTab'])

const isOpen = ref(false)
const isEditing = ref(false)
const editValue = ref('')
const inputRef = ref(null)

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
    const keys = Object.keys(props.value)
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

function saveEdit() {
  let parsedVal = editValue.value
  if (displayType.value === 'number' || typeof props.value === 'number') {
    parsedVal = Number(editValue.value)
    if (isNaN(parsedVal)) {
      alert('Must be a valid number')
      return
    }
  } else if (displayType.value === 'boolean') {
    parsedVal = Boolean(editValue.value)
  } else if (displayType.value === 'null') {
    if (editValue.value === '' || editValue.value === 'null') {
      parsedVal = null
    }
  } else if (displayType.value === 'timestamp') {
    parsedVal = { __time__: editValue.value }
  }

  isEditing.value = false
  emit('change', props.label, parsedVal)
}

function cancelEdit() {
  isEditing.value = false
}
</script>
