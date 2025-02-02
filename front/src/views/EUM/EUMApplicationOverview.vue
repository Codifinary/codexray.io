<template>
    <div class="eum-container my-10 mx-5">
        <Navigation class="my-3" :id="id" :error="selectedError" :eventId="eventId" @update:error="updateError" @update:eventId="updateEventId" />

        <v-tabs height="40" slider-color="success" show-arrows slider-size="2" v-model="activeTab" @change="updateUrl">
            <v-tab>Page Performance</v-tab>
            <v-tab>Errors</v-tab>
            <v-tab>Logs</v-tab>
            <v-tab>Traces</v-tab>
        </v-tabs>
        <v-tabs-items v-model="activeTab">
            <v-tab-item>
                <PagePerformance v-if="activeTab === 0" :data="pagePerformance" :id="id" />
            </v-tab-item>
            <v-tab-item>
                <div v-if="!selectedError">
                    <Errors v-if="errors" :data="errors" @error-clicked="handleErrorClicked" />
                </div>
                <div v-else>
                    <Error :error="selectedError" :id="id" @event-clicked="handleEventClicked" />
                </div>
            </v-tab-item>
            <v-tab-item>
                <Logs v-if="activeTab === 2" />
            </v-tab-item>
            <v-tab-item>
                <EUMTraces v-if="activeTab === 3" />
            </v-tab-item>
        </v-tabs-items>
    </div>
</template>

<script>
import { getApplicationData } from './api/EUMapi';

import PagePerformance from './PagePerformance.vue';
import Errors from './Errors.vue';
import Logs from './Logs.vue';
import Error from './Error.vue';
import Navigation from './Navigation.vue';
import EUMTraces from './EUMTraces.vue';

export default {
    name: 'EUMApplicationOverview',
    components: {
        PagePerformance,
        Errors,
        Logs,
        Error,
        EUMTraces,
        Navigation,
    },
    props: {
        id: {
            type: String,
            required: true,
        },
        report: {
            type: String,
            required: false,
        },
    },
    data() {
        return {
            activeTab: 0,
            pagePerformance: null,
            errors: null,
            logs: null,
            error: null,
            selectedError: null,
            eventId: null,
            reports: [{ name: 'page-performance' }, { name: 'errors' }, { name: 'logs' }, { name: 'traces' }],
        };
    },
    watch: {
        id: {
            immediate: true,
            handler(newId) {
                this.fetchApplicationData(newId);
            },
        },
        report: {
            immediate: true,
            handler(newReport) {
                this.setActiveTab(newReport);
            },
        },
        activeTab: {
            immediate: true,
            handler(newTab) {
                this.updateUrl(newTab);
            },
        },
        '$route.query.eventId': {
            immediate: true,
            handler(newEventId) {
                this.eventId = newEventId;
            },
        },
    },
    created() {
        this.updateUrl(this.activeTab);
    },
    methods: {
        fetchApplicationData(id) {
            const data = getApplicationData(id);
            this.pagePerformance = data.pagePerformance;
            this.errors = data.errors;
            // const appData = getAppOverview(id);
            // this.pagePerformance = appData.overviews;
        },
        updateUrl(tabIndex) {
            if (tabIndex < 0 || tabIndex >= this.reports.length) {
                console.error(`Invalid tab index: ${tabIndex}`);
                return;
            }
            const report = this.reports[tabIndex].name;
            const currentRoute = this.$route;
            const targetRoute = { name: 'overview', params: { view: 'EUM', id: this.id, report }, query: this.$utils.contextQuery() };

            if (currentRoute.name !== targetRoute.name || currentRoute.params.report !== targetRoute.params.report) {
                this.$router.push(targetRoute).catch((err) => {
                    if (err.name !== 'NavigationDuplicated') {
                        console.error(err);
                    }
                });
            }
        },
        setActiveTab(report) {
            if (!report) {
                console.error('Report is undefined');
                return;
            }
            const decodedReport = decodeURIComponent(report);
            const tabIndex = this.reports.findIndex((r) => decodedReport.startsWith(r.name));
            if (tabIndex !== -1) {
                this.activeTab = tabIndex;
                if (decodedReport.startsWith('errors/')) {
                    this.selectedError = decodedReport.split('/')[1];
                } else {
                    this.selectedError = null;
                }
            } else {
                console.error(`Invalid report: ${report}`);
            }
        },
        handleErrorClicked(error) {
            this.selectedError = error;
        },
        handleEventClicked(eventId) {
            this.$router.push({
                name: 'overview',
                params: { view: 'EUM', id: this.id },
                query: { ...this.$utils.contextQuery(), error: encodeURIComponent(this.selectedError), eventId },
            });
        },
        updateError(newError) {
            this.selectedError = newError;
        },
        updateEventId(newEventId) {
            this.eventId = newEventId;
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
    margin-left: 20px !important;
    margin-right: 20px !important;
    /* box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1) !important; */
}
</style>
