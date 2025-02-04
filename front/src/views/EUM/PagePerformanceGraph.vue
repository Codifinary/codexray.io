<template>
    <div class="mt-10">
        <Dashboard :name="'performance'" :widgets="performanceData.widgets" />
    </div>
</template>

<script>
import Dashboard from '@/components/Dashboard.vue';

export default {
    components: {
        Dashboard,
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
    async created() {
        try {
            this.$api.getPagePerformanceGraphs(this.id, this.pagePath, (data, error) => {
                // this.loading = false
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
};
</script>
