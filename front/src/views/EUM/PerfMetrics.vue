<template>
    <v-card class="performance-card">
        <v-card-title class="title">
            Performance metrics
            <span class="info-icon">ⓘ</span>
        </v-card-title>

        <v-card-text>
            <div class="metrics-grid">
                <div v-for="(key, index) in displayOrder" :key="key" class="metric-item" :class="{ 'no-border': index === displayOrder.length - 1 }">
                    <div class="metric-label">{{ labels[key] }}</div>
                    <div class="metric-value">
                        {{
                            key === 'users' || key === 'load'
                                ? `${$format.shortenNumber(data[key]).value}${$format.shortenNumber(data[key]).unit}`
                                : `${$format.formatUnits(data[key], 'ms')}ms `
                        }}
                    </div>
                </div>
            </div>
        </v-card-text>
    </v-card>
</template>

<script>
export default {
    props: {
        data: {
            type: Object,
            required: true,
        },
    },
    data() {
        return {
            labels: {
                medLoadTime: 'Med load time',
                p90LoadTime: 'p90 load time',
                avgLoadTime: 'Avg. load time',
                users: 'Users',
                load: 'Load',
            },
            displayOrder: ['medLoadTime', 'p90LoadTime', 'avgLoadTime', 'users', 'load'],
        };
    },
};
</script>

<style scoped>
.v-sheet.v-card:not(.v-sheet--outlined) {
    box-shadow: 0px 1px 5px 0px rgba(0, 0, 0, 0.12) !important;
}

.performance-card {
    border-radius: 12px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.5);
}

.title {
    font-weight: 600;
    font-size: 1.25rem;
    color: #14532d;
    display: flex;
    align-items: center;
}

.info-icon {
    margin-left: 0.5rem;
    color: #22c55e;
    font-size: 0.85rem;
}

.metrics-grid {
    display: grid;
    grid-template-columns: repeat(5, 1fr);
    text-align: center;
    gap: 1rem;
}

.metric-item {
    border-right: 2px solid #dde2e4;
    padding-right: 1rem;
}

.metric-item.no-border {
    border-right: none;
}

.metric-label {
    font-size: 0.875rem;
    color: #6b7280;
    margin-bottom: 0.25rem;
}

.metric-value {
    font-size: 1.5rem;
    font-weight: bold;
    color: #14532d;
}
</style>
