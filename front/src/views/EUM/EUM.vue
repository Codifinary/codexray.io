<template>
    <v-container class="my-10">
        <CustomTable :headers="headers" :items="tableItems" item-key="applicationName" class="elevation-1">
            <template v-slot:item.applicationName="{ item }">
                <div class="name">
                    <router-link
                        :to="{
                            name: 'overview',
                            params: { view: 'EUM', id: item.applicationName },
                        }"
                    >
                        {{ item.applicationName }}
                    </router-link>
                </div>
            </template>
        </CustomTable>
    </v-container>
</template>

<script>
import CustomTable from '@/components/CustomTable.vue';
import { getApplications } from './api/EUMapi';

export default {
    name: 'EUM',
    components: {
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
};
</script>
