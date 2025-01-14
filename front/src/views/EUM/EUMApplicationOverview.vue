<template>
    <div>
        <div class="mb-5">{{ applicationName }}</div>

        <v-tabs v-model="activeTab">
            <v-tab>Page Performance</v-tab>
            <v-tab>Errors</v-tab>
            <v-tab>Second Layer Error</v-tab>
            <v-tab>Error Details</v-tab>
            <v-tab>Logs</v-tab>
            <v-tab>Traces</v-tab>
        </v-tabs>
        <v-tabs-items v-model="activeTab">
            <v-tab-item>
                <PagePerformance v-if="pagePerformance" :data="pagePerformance" />
            </v-tab-item>
            <v-tab-item>
                <Errors v-if="errors" :data="errors" />
            </v-tab-item>
            <v-tab-item>
                <Error />
            </v-tab-item>
            <v-tab-item>
                <ErrorDetails />
            </v-tab-item>
            <v-tab-item>
                <Logs v-if="logs" :data="logs" />
            </v-tab-item>
        </v-tabs-items>
    </div>
</template>

<script>
import { getApplicationData } from './api/EUMapi';
import PagePerformance from './PagePerformance.vue';
import Errors from './Errors.vue';
import Logs from './Logs.vue';
import ErrorDetails from './ErrorDetail.vue';
import Error from './Error.vue';

export default {
    name: 'EUMApplicationOverview',
    components: {
        PagePerformance,
        Errors,
        Logs,
        ErrorDetails,
        Error,
    },
    props: {
        applicationName: {
            type: String,
            required: true,
        },
    },
    data() {
        return {
            activeTab: 0,
            pagePerformance: null,
            errors: null,
            logs: null,
        };
    },
    watch: {
        applicationName: {
            immediate: true,
            handler(newName) {
                this.fetchApplicationData(newName);
            },
        },
    },
    methods: {
        fetchApplicationData(applicationName) {
            const data = getApplicationData(applicationName);
            this.pagePerformance = data.pagePerformance;
            this.errors = data.errors;
            this.logs = data.logs.errors.eventLogs;
        },
    },
};
</script>

<style scoped></style>
