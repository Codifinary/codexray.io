<template>
    <div class="ml-5 container">
        <div class="mb-5">
            <span class="heading">Summary of all browser applications</span>
        </div>
        <div class="card-container">
            <EUMCard :data="cardData" />

            <!-- <div class="chart-container"> -->
            <div v-for="(config, index) in chartData" :key="index" class="chart-wrapper">
                <EChart :chartOptions="config" class="chart-box" :style="getChartStyle(index)" />
            </div>
            <!-- </div> -->
        </div>
    </div>
</template>

<script>
import EUMCard from '@/components/EUMCard.vue';
import EChart from '@/components/EChart.vue';

export default {
    props: {
        cardData: {
            type: Object,
            required: true,
        },
        chartData: {
            type: Array,
            required: true,
        },
    },
    components: {
        EUMCard,
        EChart,
    },
    methods: {
        getChartStyle(index) {
            if (index < 2) {
                // First two charts
                return { width: '18vw', height: '30vh' };
            } else {
                // Last chart
                return { width: '25vw', height: '30vh' };
            }
        },
    },
};
</script>

<style scoped>
.container {
    max-width: 100%;
    margin-left: 0;
}
.card-container {
    display: flex;
    flex-direction: row;
    max-width: 100%;
    gap: 20px;
}
.chart-wrapper {
    box-shadow:
        0px 3px 1px -2px rgba(0, 0, 0, 0.2),
        0px 2px 2px 0px rgba(0, 0, 0, 0.14),
        0px 1px 5px 0px rgba(0, 0, 0, 0.12);
}
.chart-box {
    transform: scale(0.9);
    transform-origin: center;
}
.heading {
    color: var(--status-ok) !important;
    font-size: 14px !important;
    font-weight: 600 !important;
}
</style>
