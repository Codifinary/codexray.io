<template>
    <div class="my-10 mx-5">
        <CustomTable :headers="headers" :items="tableItems" item-key="applicationName" class="elevation-1">
            <template v-slot:item.applicationName="{ item }">
                <div class="name d-flex">
                    <div class="mr-3">
                        <img
                            :src="`${$codexray.base_path}static/img/tech-icons/${item.applicationType}.svg`"
                            style="width: 16px; height: 16px"
                            alt="App Icon"
                        />
                    </div>
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
    </div>
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
                { text: 'Application', value: 'applicationName' },
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
<style scoped>
.eum-container {
    padding-bottom: 70px;
    margin-left: 20px !important;
    margin-right: 20px !important;
    /* box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1) !important; */
}
</style>
