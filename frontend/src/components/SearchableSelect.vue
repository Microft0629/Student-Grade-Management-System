<!-- 可搜索下拉选择组件 -->
<script setup>
import { ref, computed } from 'vue'
import { vClickOutside } from '../composables/clickOutside'

const props = defineProps({
  modelValue: [Number, String],
  options: { type: Array, default: () => [] },
  placeholder: { type: String, default: '请选择...' },
  emptyText: { type: String, default: '无匹配结果' },
})

const emit = defineEmits(['update:modelValue'])

const keyword = ref('')
const open = ref(false)
const highlighted = ref(-1)

const filtered = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  if (!kw) return props.options
  return props.options.filter(o => {
    return o.label.toLowerCase().includes(kw) ||
           (o.searchText || '').toLowerCase().includes(kw) ||
           (o.term || '').toLowerCase().includes(kw)
  })
})

const selectedLabel = computed(() => {
  const opt = props.options.find(o => o.value === props.modelValue)
  return opt ? opt.label : ''
})

function select(opt) {
  emit('update:modelValue', opt.value)
  keyword.value = ''
  open.value = false
  highlighted.value = -1
}

function toggle() {
  open.value = !open.value
  if (!open.value) {
    keyword.value = ''
    highlighted.value = -1
  }
}

function close() {
  open.value = false
  keyword.value = ''
  highlighted.value = -1
}

function onKeydown(e) {
  if (!open.value) return
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    highlighted.value = Math.min(highlighted.value + 1, filtered.value.length - 1)
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    highlighted.value = Math.max(highlighted.value - 1, -1)
  } else if (e.key === 'Enter') {
    e.preventDefault()
    if (highlighted.value >= 0 && highlighted.value < filtered.value.length) {
      select(filtered.value[highlighted.value])
    }
  } else if (e.key === 'Escape') {
    close()
  }
}

function onFocus() {
  open.value = true
}
</script>

<template>
  <div class="ssel-wrap" v-click-outside="close">
    <div class="ssel-input-row" @click="toggle">
      <span v-if="!open && modelValue" class="ssel-selected">{{ selectedLabel }}</span>
      <span v-else-if="!open" class="ssel-placeholder">{{ placeholder }}</span>
      <input
        v-else
        ref="inputRef"
        v-model="keyword"
        class="ssel-search"
        :placeholder="'输入关键词筛选...'"
        @focus="onFocus"
        @keydown="onKeydown"
        @click.stop
      />
      <span class="ssel-arrow" :class="{ open }">▾</span>
    </div>

    <div v-if="open" class="ssel-dropdown">
      <div
        v-for="(opt, i) in filtered"
        :key="opt.value"
        class="ssel-option"
        :class="{ highlighted: i === highlighted, selected: opt.value === modelValue }"
        @click="select(opt)"
      >
        <span>{{ opt.label }}</span>
        <span v-if="opt.term" class="ssel-sub">{{ opt.term }}</span>
      </div>
      <div v-if="filtered.length === 0" class="ssel-empty">{{ emptyText }}</div>
    </div>
  </div>
</template>

<style scoped>
.ssel-wrap {
  position: relative;
  flex: 1;
  max-width: 480px;
}
.ssel-input-row {
  display: flex;
  align-items: center;
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  padding: 0 10px;
  cursor: pointer;
  background: #fff;
  min-height: 38px;
  box-sizing: border-box;
  transition: border-color 0.2s;
}
.ssel-input-row:focus-within {
  border-color: #4a90d9;
  box-shadow: 0 0 0 2px rgba(74,144,217,0.15);
}
.ssel-selected {
  flex: 1;
  font-size: 14px;
  color: #333;
}
.ssel-placeholder {
  flex: 1;
  font-size: 14px;
  color: #bbb;
}
.ssel-search {
  flex: 1;
  border: none !important;
  outline: none !important;
  box-shadow: none !important;
  padding: 0;
  font-size: 14px;
  background: transparent;
  min-width: 0;
}
.ssel-arrow {
  font-size: 12px;
  color: #bbb;
  margin-left: 8px;
  transition: transform 0.2s;
  flex-shrink: 0;
}
.ssel-arrow.open {
  transform: rotate(180deg);
}
.ssel-dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  margin-top: 4px;
  background: #fff;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  box-shadow: 0 6px 20px rgba(0,0,0,0.12);
  max-height: 240px;
  overflow-y: auto;
  z-index: 100;
}
.ssel-option {
  padding: 10px 14px;
  cursor: pointer;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 14px;
  transition: background 0.15s;
}
.ssel-option:hover,
.ssel-option.highlighted {
  background: #f0f5ff;
}
.ssel-option.selected {
  background: #ecf3fc;
  color: #4a90d9;
  font-weight: 600;
}
.ssel-sub {
  font-size: 12px;
  color: #bbb;
  margin-left: 8px;
}
.ssel-empty {
  padding: 14px;
  text-align: center;
  color: #ccc;
  font-size: 13px;
}
</style>
