<template>
    <div class="mt-10">
        <Navigation :id="id" :pagePath="pagePath" />

        <div v-if="mappedExperiences.length">
            <PageExp :experiences="mappedExperiences" />
        </div>

        <div class="my-5">
            <PerfMetrics
                :data="{
                    medLoadTime: performanceMetrics.MedianLoadTime,
                    p90LoadTime: performanceMetrics.P90LoadTime,
                    avgLoadTime: performanceMetrics.AvgLoadTime,
                    users: performanceMetrics.UniqueUsers,
                    load: performanceMetrics.LoadCount,
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

const STATUS_STYLES = {
    Good: {
        color: '#66BB6A',
        background: '#EBFFF6',
    },
    Moderate: {
        color: '#FFA726',
        background: '#FFF5E7',
    },
    Poor: {
        color: '#EF5350',
        background: '#FFEDED',
    },
};

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
            performanceMetrics: {},
            pageExperience: {},
            loading: false,
            error: null,
        };
    },
    computed: {
        mappedExperiences() {
            const experienceMap = {
                PageLoadingClass: 'Page Loading',
                PageInteractivityClass: 'Page Interactivity',
                PageRenderingClass: 'Page Rendering',
                ResourceLoadingClass: 'Resource Loading',
                ServerPerformanceClass: 'Server Response',
            };

            return Object.entries(this.pageExperience).map(([key, status]) => {
                const title = experienceMap[key] || key;
                const style = STATUS_STYLES[status];
                return {
                    title,
                    status,
                    color: style.color,
                    background: style.background,
                };
            });
        },
    },
    mounted() {
        this.getPerformanceData();
        this.$events.watch(this, this.getPerformanceData, 'refresh');
    },
    methods: {
        getPerformanceData() {
            this.loading = true;
            this.error = null;
            this.$api.getPagePerformanceGraphs(this.id, this.pagePath, (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    console.error('Error fetching performance data:', error);
                    return;
                }
                this.performanceData = data.performanceCharts || {};
                this.performanceMetrics = data.performanceMetrics || {};
                this.pageExperience = data.experienceScores || {};
            });
        },
    },
};
</script>
