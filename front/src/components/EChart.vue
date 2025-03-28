<template>
    <div ref="chart" class="chart"></div>
</template>

<script>
import { ref, onMounted, watch } from 'vue';
import * as echarts from 'echarts';

export default {
    props: {
        chartOptions: Object,
    },
    setup(props) {
        const chart = ref(null);
        let myChart = null;

        const renderChart = () => {
            if (!chart.value || !props.chartOptions) return;

            if (myChart) {
                myChart.dispose();
            }

            myChart = echarts.init(chart.value);
            myChart.setOption(props.chartOptions);
        };

        onMounted(renderChart);
        watch(() => props.chartOptions, renderChart, { deep: true });

        return { chart };
    },
};
</script>

<style scoped>
.chart {
    width: 100%;
    height: 100%;
    transform: scale(0.85);
    transform-origin: center;
}
</style>