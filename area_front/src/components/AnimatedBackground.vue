<template>
  <div class="animated-background">
    <div class="floating-shapes">
      <div
        v-for="(shape, index) in shapes"
        :key="index"
        class="shape"
        :class="`shape-${index + 1}`"
        :style="shape.style"
      ></div>
    </div>
    <div class="gradient-overlay"></div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Shape {
  style: Record<string, string>
}

interface Props {
  shapeCount?: number
  enableGlobe?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  shapeCount: 5,
  enableGlobe: true
})

const shapes = computed<Shape[]>(() => {
  const shapesArray: Shape[] = []

  for (let i = 0; i < props.shapeCount; i++) {
    const size = Math.random() * 200 + 100 // 100-300px
    const top = Math.random() * 80 + 10 // 10-90%
    const left = Math.random() * 80 + 10 // 10-90%
    const animationDelay = Math.random() * 8 // 0-8s

    shapesArray.push({
      style: {
        width: `${size}px`,
        height: `${size}px`,
        top: `${top}%`,
        left: `${left}%`,
        animationDelay: `${animationDelay}s`
      }
    })
  }

  return shapesArray
})
</script>

<style scoped>
.animated-background {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1;
  pointer-events: none;
}

.floating-shapes {
  position: absolute;
  width: 100%;
  height: 100%;
}

.shape {
  position: absolute;
  border-radius: 50%;
  background: linear-gradient(45deg, rgba(87, 128, 232, 0.1), rgba(135, 81, 209, 0.1));
  filter: blur(2px);
  animation: float 8s ease-in-out infinite;
}

.gradient-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle at 30% 20%, rgba(87, 128, 232, 0.1) 0%, transparent 50%),
              radial-gradient(circle at 70% 80%, rgba(135, 81, 209, 0.1) 0%, transparent 50%);
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
  }
  50% {
    transform: translateY(-30px) rotate(10deg);
  }
}
</style>
