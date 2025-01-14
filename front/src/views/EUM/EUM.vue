<template>
    <v-container class="my-10">
        <template v-if="selectedApplication">
            <EUMApplicationOverview :application-name="selectedApplication" />
        </template>

        <template v-else>
            <CustomTable :headers="headers" :items="tableItems" item-key="applicationName" class="elevation-1" @click:row="navigateToOverview">
                <template v-slot:[`item.applicationName`]="{ item }">
                    <a href="#" @click.prevent="navigateToOverview(item)">{{ item.applicationName }}</a>
                </template>
            </CustomTable>
        </template>
    </v-container>
</template>

<script>
import CustomTable from '@/components/CustomTable.vue';
import { getApplications } from './api/EUMapi';
import EUMApplicationOverview from './EUMApplicationOverview.vue';

export default {
    name: 'EUM',
    components: {
        EUMApplicationOverview,
        CustomTable,
    },
    data() {
        return {
            headers: [
                { text: 'Application Name', value: 'applicationName' },
                { text: 'Application Type', value: 'applicationType' },
                { text: 'Number of Pages', value: 'noOfPages' },
                { text: 'Load Requests/Second', value: 'loadRequestsPerSecond' },
                { text: 'Response time (ms)', value: 'responseTimeMs' },
                { text: 'JS Error%', value: 'jsErrorPercentage' },
                { text: 'API Error%', value: 'apiErrorPercentage' },
                { text: 'Users Impacted', value: 'usersImpacted' },
            ],
            tableItems: [],
            selectedApplication: null,
        };
    },
    created() {
        const applications = getApplications();
        this.tableItems = applications;
    },
    methods: {
        navigateToOverview(item) {
            this.selectedApplication = item.applicationName;
        },
    },
};
</script>

<style scoped></style>
