<template>
  <div
    ref="refElement"
    class="glare-card-container"
    @pointermove="handlePointerMove"
    @pointerenter="handlePointerEnter"
    @pointerleave="handlePointerLeave"
    @click="handleClick"
  >
    <div class="glare-card-inner">
      <div class="glare-card-content">
        <div :class="cn('glare-card-slot', props.class)">
          <slot />
        </div>
      </div>
      <div class="glare-overlay"></div>
      <div class="brand-overlay"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { cn } from "@/lib/utils";
import { useTimeoutFn } from "@vueuse/core";
import { ref } from "vue";

interface GlareCardProps {
  class?: string;
}

const props = defineProps<GlareCardProps>();
const emit = defineEmits<{
  click: [event: MouseEvent]
}>();

const isPointerInside = ref(false);
const refElement = ref<HTMLElement | null>(null);

const state = ref({
  glare: { x: 50, y: 50 },
  background: { x: 50, y: 50 },
  rotate: { x: 0, y: 0 },
});

function handlePointerMove(event: PointerEvent) {
  const rotateFactor = 0.4;
  const rect = refElement.value?.getBoundingClientRect();
  if (rect) {
    const position = {
      x: event.clientX - rect.left,
      y: event.clientY - rect.top,
    };
    const percentage = {
      x: (100 / rect.width) * position.x,
      y: (100 / rect.height) * position.y,
    };
    const delta = {
      x: percentage.x - 50,
      y: percentage.y - 50,
    };
    state.value.background.x = 50 + percentage.x / 4 - 12.5;
    state.value.background.y = 50 + percentage.y / 3 - 16.67;
    state.value.rotate.x = -(delta.x / 3.5) * rotateFactor;
    state.value.rotate.y = (delta.y / 2) * rotateFactor;
    state.value.glare.x = percentage.x;
    state.value.glare.y = percentage.y;
  }
}

function handlePointerEnter() {
  isPointerInside.value = true;
  useTimeoutFn(() => {
    if (isPointerInside.value && refElement.value) {
      refElement.value.style.setProperty("--duration", "0s");
    }
  }, 300);
}

function handlePointerLeave() {
  isPointerInside.value = false;
  if (refElement.value) {
    refElement.value.style.removeProperty("--duration");
    state.value.rotate = { x: 0, y: 0 };
  }
}

function handleClick(event: MouseEvent) {
  console.log('GlareCard clicked!', event)
  emit('click', event)
}
</script>

<style scoped>
.glare-card-container {
  width: 320px;
  height: 420px;
  perspective: 600px;
  position: relative;
  cursor: pointer;
  --m-x: v-bind(state.glare.x + "%");
  --m-y: v-bind(state.glare.y + "%");
  --r-x: v-bind(state.rotate.x + "deg");
  --r-y: v-bind(state.rotate.y + "deg");
  --bg-x: v-bind(state.background.x + "%");
  --bg-y: v-bind(state.background.y + "%");
  --duration: 300ms;
  --opacity: 0;
}

.glare-card-inner {
  width: 100%;
  height: 100%;
  position: relative;
  transform-style: preserve-3d;
  transform: rotateY(var(--r-x)) rotateX(var(--r-y));
  transition: transform var(--duration) ease;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid #1e293b;
}

.glare-card-content {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: #0f172a;
  z-index: 1;
}

.glare-card-slot {
  width: 100%;
  height: 100%;
}

.glare-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: radial-gradient(
    circle at var(--m-x) var(--m-y),
    rgba(255, 255, 255, 0.8) 10%,
    rgba(255, 255, 255, 0.65) 20%,
    rgba(255, 255, 255, 0) 90%
  );
  opacity: var(--opacity);
  transition: opacity var(--duration) ease;
  z-index: 2;
  pointer-events: none;
}

.brand-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    135deg,
    var(--color-accent-primary) 0%,
    var(--color-accent-secondary) 25%,
    var(--color-accent-tertiary) 50%,
    var(--color-accent-secondary) 75%,
    var(--color-accent-primary) 100%
  );
  background-position: 0% var(--bg-y);
  background-size: 200% 200%;
  opacity: var(--opacity);
  transition: opacity var(--duration) ease;
  z-index: 3;
  pointer-events: none;
  mix-blend-mode: color-dodge;
}

.glare-card-container:hover {
  --duration: 200ms;
  --opacity: 0.6;
}
</style>
