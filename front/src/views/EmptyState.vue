<template>
    <div :style="{ height, width }">
        <span class="ml-8 mt-5 sub-heading">{{ heading }}</span>

        <div class="empty-state-container">
            <div class="empty-state-content">
                <!-- Icon -->
                <div class="empty-state-icon">
                    <BaseIcon :name="iconName" :width="iconWidth" :height="iconHeight" />
                </div>
                <!-- Title -->
                <div v-if="title" class="empty-state-title">{{ title }}</div>
                <!-- Description -->
                <div v-if="description" class="empty-state-description">{{ description }}</div>
                <!-- Button -->
                <button v-if="buttonText" class="empty-state-button">{{ buttonText }}</button>

                <div v-if="helpText" class="help-text">
                    <span>{{ helpText }}</span>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import BaseIcon from '@/components/BaseIcon.vue';

export default {
    name: 'EmptyState',
    props: {
        height: {
            type: String,
            default: '70vh',
        },
        width: {
            type: String,
            default: '100%',
        },
        iconWidth: {
            type: String,
            default: '15vw',
        },
        iconHeight: {
            type: String,
            default: '15vh',
        },
        heading: {
            type: String,
        },
        title: {
            type: String,
            required: true,
        },
        description: {
            type: String,
            required: true,
        },
        buttonText: {
            type: String,
        },
        iconName: {
            type: String,
            default: 'NoData',
        },
        helpText: {
            type: String,
        },
    },
    components: {
        BaseIcon,
    },
};
</script>

<style scoped>
.empty-state-container {
    display: flex;
    align-items: center;
    justify-content: center;
}
.empty-state-content {
    display: flex;
    flex-direction: column;
    align-items: center;
}
.empty-state-icon {
    margin-bottom: 0.5rem;
}
.sub-heading {
    font-size: 1.2rem;
    font-weight: 500;
    color: #333;
}
.empty-state-title {
    font-size: 1.25rem;
    font-weight: 600;
    margin-bottom: 0.25rem;
}
.empty-state-description {
    color: #666;
    font-size: 1rem;
    margin-bottom: 0.75rem;
}
.empty-state-button {
    background: #22c55e;
    color: #fff;
    border: none;
    border-radius: 0.375rem;
    padding: 0.625rem 1.5rem;
    font-size: 1rem;
    font-weight: 500;
    cursor: pointer;
    transition: background 0.2s;
}
.empty-state-button:hover {
    background: #16a34a;
}
.help-text {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4px;
    color: #013912;
    padding-top: 10px;
}
</style>
