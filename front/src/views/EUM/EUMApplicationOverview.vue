<template>
    <div class="eum-container my-10">
        <Navigation class="my-3" :id="id" :error="selectedError" :eventId="eventId" @update:error="updateError" @update:eventId="updateEventId" />

        <v-tabs v-model="activeTab" height="40" slider-color="success" show-arrows slider-size="2" @change="updateUrl">
            <v-tab v-for="(report, index) in reports" :key="index">{{ report.label }}</v-tab>
        </v-tabs>

        <v-tabs-items v-model="activeTab">
            <v-tab-item v-for="(report, index) in reports" :key="index">
                <component
                    :is="report.component"
                    v-if="activeTab === index"
                    :id="id"
                    :error="selectedError"
                    :eventId="eventId"
                    :report="reports[index].name"
                    @update:error="updateError"
                    @update:eventId="updateEventId"
                />
            </v-tab-item>
        </v-tabs-items>
    </div>
</template>

<script>
import PagePerformance from './PagePerformance.vue';
import Errors from './Errors.vue';
import Logs from './Logs.vue';
import Navigation from './Navigation.vue';
import EUMTraces from './EUMTraces.vue';

export default {
    name: 'EUMApplicationOverview',
    components: { PagePerformance, Errors, Logs, EUMTraces, Navigation },
    props: {
        id: { type: String, required: true },
        report: { type: String, required: false, default: '' },
    },
    data() {
        return {
            activeTab: 0,
            selectedError: null,
            eventId: null,
            reports: [
                { name: 'page-performance', label: 'Page Performance', component: 'PagePerformance' },
                { name: 'errors', label: 'Errors', component: 'Errors' },
                { name: 'logs', label: 'Logs', component: 'Logs' },
                { name: 'traces', label: 'Traces', component: 'EUMTraces' },
            ],
        };
    },
    watch: {
        report: { immediate: true, handler: 'setActiveTab' },
        '$route.params.report': { immediate: true, handler: 'setActiveTab' },
    },
    created() {
        const report = this.$route.params.report || this.report;
        if (!report) {
            this.$router
                .replace({
                    name: 'overview',
                    params: { view: 'BRUM', id: this.id, report: this.reports[0].name },
                    query: this.$utils.contextQuery(),
                })
                .catch(this.handleNavigationError);
        } else {
            this.setActiveTab(report);
        }
        this.eventId = this.$route.query.eventId || null;
    },
    methods: {
        updateUrl(tabIndex) {
            if (tabIndex < 0 || tabIndex >= this.reports.length) return;
            const report = this.reports[tabIndex].name;
            if (this.$route.params.report !== report) {
                this.$router
                    .push({
                        name: 'overview',
                        params: { view: 'BRUM', id: this.id, report },
                        query: { ...this.$utils.contextQuery(), error: this.selectedError || undefined },
                    })
                    .catch(this.handleNavigationError);
            }
        },
        setActiveTab(report) {
            if (!report) return;

            const decodedReport = decodeURIComponent(report);
            const tabIndex = this.reports.findIndex((r) => decodedReport.startsWith(r.name));

            if (tabIndex !== -1) {
                this.activeTab = tabIndex;
                if (decodedReport.startsWith('errors/')) {
                    this.selectedError = decodedReport.split('/')[1] || null;
                } else {
                    this.selectedError = null;
                }
            } else {
                this.activeTab = 0; // Fallback to first tab if not found
            }
        },
        updateError(newError) {
            this.selectedError = newError;
        },
        updateEventId(newEventId) {
            this.eventId = newEventId;
        },
        handleNavigationError(err) {
            if (err.name !== 'NavigationDuplicated') {
                console.error('Navigation error:', err);
            }
        },
    },
};
</script>

<style scoped>
.v-tab {
    color: var(--primary-green) !important;
    margin-left: 15px;
    text-decoration: none !important;
    text-transform: none !important;
    margin-top: 5px;
    font-weight: 400 !important;
}
.v-slide-group__wrapper {
    width: 100%;
}
.v-slide-group__content {
    position: static;
    border-bottom: 2px solid #0000001a !important;
}
.eum-container {
    padding-bottom: 70px;
}
</style>
