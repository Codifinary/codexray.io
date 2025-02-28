<template>
    <div>
        <div class="card-container">
            <EUMCard :data="cardData" />
            <div class="chart-container">
                <EChart v-for="(config, index) in chartData" :key="index" :chartOptions="config" class="chart-box" />
            </div>
        </div>
        <div class="graph-container">
            <EChart v-for="i in 4" :key="i" :chartOptions="chartData?.lineChart || {}" class="graph-box" />
        </div>
    </div>
</template>

<script>
import EUMCard from '@/components/EUMCard.vue';
import EChart from '@/components/EChart.vue';
import chartData from './ChartData.json';
import { fetchCardData } from './EUMapi.js';

export default {
    components: {
        EUMCard,
        EChart,
    },

    data() {
        return {
            cardData: [],
            chartData: {},
        };
    },
    async mounted() {
        try {
            this.cardData = await fetchCardData();
            this.chartData = chartData;
        } catch (error) {
            console.error('Error fetching card data:', error);
        }
    },
};
</script>

<style scoped>
.card-container {
    display: flex;
    flex-direction: row;
    width: 100%;
    justify-content: space-between;
}
.chart-container {
    display: flex;
    flex-direction: row;
    gap: 20px;
    justify-content: center;
    padding-left: 20px;
}
.chart-box {
    width: 350px;
    height: 410px;
    box-shadow: 0px 4px 4px rgba(0, 0, 0, 0.3);
}
.graph-container {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    padding: 20px;
    margin: 20px 0;
    gap: 20px;
}
.graph-box {
    width: 48%;
    height: 400px;
    box-shadow: 0px 4px 4px rgba(0, 0, 0, 0.3);
}
</style>
