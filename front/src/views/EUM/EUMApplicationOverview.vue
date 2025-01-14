<template>
    <v-container>
        <v-card>
            <v-card-title>Overview: {{ applicationName }}</v-card-title>

            <v-tabs v-model="activeTab">
                <v-tab>Page Performance</v-tab>
                <v-tab>Errors</v-tab>
                <v-tab>Logs</v-tab>
            </v-tabs>
            <v-tabs-items v-model="activeTab">
                <v-tab-item>
                    <PagePerformance v-if="pagePerformance" :data="pagePerformance" />
                </v-tab-item>
                <v-tab-item>
                    <Errors v-if="errors" :data="errors" />
                </v-tab-item>
                <v-tab-item>
                    <Logs v-if="logs" :data="logs" />
                </v-tab-item>
            </v-tabs-items>
        </v-card>
    </v-container>
</template>

<script>
import { getApplicationData } from './api/EUMapi';
import PagePerformance from './PagePerformance.vue';
import Errors from './Errors.vue';
import Logs from './Logs.vue';

export default {
    name: 'EUMApplicationOverview',
    components: {
        PagePerformance,
        Errors,
        Logs,
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
