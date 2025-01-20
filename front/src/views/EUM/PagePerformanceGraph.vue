<template>
    <div class="eum-application-overview my-10 mx-5">
        <Navigation :id="id" :pagePath="pagePath" />
        <div class="d-flex my-10 justify-space-between">
            <Chart :chart="loadChartConfig" />
            <Chart :chart="responseChartConfig" />
            <Chart :chart="errorChartConfig" />
        </div>
        <div class="d-flex">
            <Chart :chart="firstUserPaintChartConfig" />
            <div class="ml-10"></div>
            <Chart :chart="userImpactedChartConfig" :selection="{}" @select="zoomChart" />
        </div>
    </div>
</template>

<script>
import Navigation from './Navigation.vue';
import Chart from '@/components/Chart.vue';

export default {
    props: {
        id: String,
        pagePath: String,
    },
    components: { Navigation, Chart },
    data() {
        return {
            errorChartConfig: {
                title: 'Errors', //area chart
                ctx: { from: 0, to: 50, step: 1 },
                series: [
                    { name: 'JS Errors', data: this.generateSineWaveData(10, 0.1, 0, 50), color: 'red', fill: true },
                    { name: 'AJAX Errors', data: this.generateSineWaveData(8, 0.1, Math.PI / 2, 50), color: 'blue', fill: true },
                ],
            },
            userImpactedChartConfig: {
                title: 'User Impacted', // Bar Chart
                ctx: { from: 0, to: 50, step: 1 },
                series: [
                    {
                        name: 'Users Impacted',
                        data: this.generateRandomData(50, 100, 50),
                        color: 'green',
                    },
                ],
            },
            loadChartConfig: {
                title: 'Page Load Time', // Line Chart
                ctx: { from: 0, to: 50, step: 1 },
                series: [
                    {
                        name: 'Load Time (ms)',
                        data: this.generateRandomData(20, 500, 100),
                        color: 'orange',
                    },
                ],
            },
            responseChartConfig: {
                title: 'Response Time ', //Stacked Bar Chart
                ctx: { from: 0, to: 50, step: 1 },
                series: [
                    {
                        name: 'Response Time (ms)',
                        data: this.generateRandomData(15, 200, 50),
                        color: 'purple',
                        bars: true,
                        stacked: true,
                    },
                ],
            },
            firstUserPaintChartConfig: {
                title: 'First Meaningful Paint ', //Scatter Plot
                ctx: { from: 0, to: 50, step: 1 },
                series: [
                    {
                        name: 'First Meaningful Paint(ms)',
                        data: this.generateRandomData(30, 300, 100),
                        color: 'teal',
                        points: { show: true },
                        width: 0,
                    },
                ],
            },
        };
    },
    methods: {
        generateSineWaveData(amplitude, frequency, phase, length) {
            const data = [];
            for (let i = 0; i < length; i++) {
                data.push(amplitude * Math.sin(frequency * i + phase));
            }
            return data;
        },
        generateRandomData(length, max, min) {
            const data = [];
            for (let i = 0; i < length; i++) {
                data.push(Math.random() * (max - min) + min);
            }
            return data;
        },
    },
};
</script>

<style scoped>
.chart {
    width: 30%;
    height: 300px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1) !important;
    margin: 10px;
}
.eum-container {
    padding-bottom: 70px;
    margin-left: 20px !important;
    margin-right: 20px !important;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1) !important;
    margin-top: 10px;
}
</style>
