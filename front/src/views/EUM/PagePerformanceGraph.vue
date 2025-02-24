<template>
    <div class="mt-10">
        <Navigation :id="id" />
        <Dashboard :name="'performance'" :widgets="performanceData.widgets" />
    </div>
</template>

<script>
import Dashboard from '@/components/Dashboard.vue';
import Navigation from './Navigation.vue';

export default {
    components: {
        Dashboard,
        Navigation,
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
