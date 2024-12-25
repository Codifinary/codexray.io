<template>
    <span @click="emitClickEvent">
        <component
            v-if="iconComponent"
            :is="iconComponent"
            role="img"
            :class="className"
            :style="{
                width: width,
                height: height,
                fill: iconColor,
                transform: transform,
            }"
        />
    </span>
</template>

<script>
const icons = {};
const requireComponents = require.context('@/assets/icons', false, /\.svg$/);
requireComponents.keys().forEach((fileName) => {
    const iconName = fileName.replace(/^\.\/(.+)\.svg$/, '$1');
    const componentConfig = requireComponents(fileName);

    // Ensure proper handling of the imported SVG
    icons[iconName] = componentConfig.default || componentConfig;
});

export default {
    emits: ['onClick'],
    props: {
        name: {
            type: String,
            required: true,
            validator(value) {
                return Object.prototype.hasOwnProperty.call(icons, value);
            },
        },
        width: {
            type: String,
        },
        height: {
            type: String,
        },
        iconColor: {
            type: String,
        },
        className: {},
        transform: {
            type: String,
        },
    },
    computed: {
        iconComponent() {
            return icons[this.name] || null;
        },
    },
    methods: {
        emitClickEvent() {
            this.$emit('onClick');
        },
    },
};
</script>
<style lang="scss" scoped></style>
