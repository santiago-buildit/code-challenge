<template>
    <span :class="['badge', variantClass]">
        <slot>{{ label }}</slot>
    </span>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
    name: 'BaseBadge',
    props: {
        label: {
            type: String,
            default: '',
        },
        variant: {
            type: String as () => 'success' | 'error' | 'neutral' | 'info' | 'warning',
            default: 'neutral',
        },
    },
    computed: {
        variantClass(): string {
            return `badge-${this.variant}`
        },
    },
})
</script>

<style scoped>
.badge {
    padding: 0.25rem 0.5rem;
    border-radius: 0.5rem;
    font-size: 0.75rem;
    font-weight: bold;
    text-transform: uppercase;
    display: inline-block;
    min-width: 80px;
    text-align: center;
}

.badge-success {
    background-color: #d4edda;
    color: #155724;
}

.badge-error {
    background-color: #f8d7da;
    color: #721c24;
}

.badge-warning {
    background-color: #fff3cd;
    color: #856404;
}

.badge-info {
    background-color: #d1ecf1;
    color: #0c5460;
}

.badge-neutral {
    background-color: #e2e3e5;
    color: #383d41;
}
</style>