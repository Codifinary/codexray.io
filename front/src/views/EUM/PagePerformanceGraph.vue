<template>
    <div class="mt-10">
        <Navigation :id="id" :pagePath="pagePath" />
        <div>
            <PageExp
                :experiences="[
                    {
                        title: 'Page loading exp',
                        status: 'moderate',
                        color: '#FFA726', // orange
                        background: '#FFF5E7', // light orange
                    },
                    {
                        title: 'Page Interactivity',
                        status: 'good',
                        color: '#66BB6A', // green
                        background: '#EBFFF6', // light green
                    },
                    {
                        title: 'Page Rendering',
                        status: 'poor',
                        color: '#EF5350', // red
                        background: '#FFEDED', // light red
                    },
                    {
                        title: 'Resource Loading',
                        status: 'good',
                        color: '#66BB6A', // green
                        background: '#EBFFF6', // light green
                    },
                    {
                        title: 'Page loading exp',
                        status: 'moderate',
                        color: '#FFA726', // orange
                        background: '#FFF5E7', // light orange
                    },
                ]"
            />
        </div>
        <div class="my-5">
            <PerfMetrics
                :data="{
                    medLoadTime: 2.35,
                    p90LoadTime: 8.6,
                    avgLoadTime: 11.43,
                    users: 93,
                    load: 299,
                }"
            />
        </div>
        <Dashboard :name="'performance'" :widgets="performanceData.widgets" />
    </div>
</template>

<script>
import Dashboard from '@/components/Dashboard.vue';
import Navigation from './Navigation.vue';
import PerfMetrics from './PerfMetrics.vue';
import PageExp from './PageExp.vue';

export default {
    components: {
        Dashboard,
        Navigation,
        PerfMetrics,
        PageExp,
    },
    props: {
        id: {
            type: String,
            required: true,
        },
        pagePath: {
            type: String,
            required: true,
        },
    },
    data() {
        return {
            performanceData: {},
        };
    },
    mounted() {
        this.get();
        this.$events.watch(this, this.get, 'refresh');
    },
    methods: {
        get() {
            try {
                this.loading = true;
                this.$api.getPagePerformanceGraphs(this.id, this.pagePath, (data, error) => {
                    this.loading = false;
                    if (error) {
                        this.error = error;
                        return;
                    }
                    this.performanceData = data || [];
                });
            } catch (error) {
                console.error('Error fetching performance data:', error);
            }
        },
    },
};
</script>
