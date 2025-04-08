<template>
    <div>
        <v-progress-linear indeterminate v-if="loading" color="green" />

        <v-alert v-if="error" color="red" icon="mdi-alert-octagon-outline" outlined text>
            {{ error }}
        </v-alert>

        <CustomTable 
            v-if="sessions.summary" 
            :headers="headers" 
            :items="sessions.summary" 
            class="table"
        >
            <template #item.serviceName="{ item: { serviceName } }">
                <div class="service-name clickable" @click="goToService(serviceName)">
                    {{ serviceName }}
                </div>
            </template>
            <template #item.totalUsers="{ item: { totalUsers } }">
                <div class="total-users">
                    {{ totalUsers }}
                </div>
            </template>
            <template #item.totalRequests="{ item: { totalRequests } }">
                <div class="total-requests">
                    {{ totalRequests }}
                </div>
            </template>
            <template #item.totalErrors="{ item: { totalErrors } }">
                <div class="total-errors">
                    {{ totalErrors }}
                </div>
            </template>
            <template #item.totalSessions="{ item: { totalSessions } }">
                <div class="total-sessions">
                    {{ totalSessions }}
                </div>
            </template>
        </CustomTable>
    </div>
</template>

<script>
import CustomTable from '@/components/CustomTable.vue';
import mockData from './overview.json';

export default {

    props: {
        id: String,
        projectId: String
    },
    components: {
        CustomTable
    },

    data() {
        return {
            error: null,
            loading: false,
            sessions: {},
            headers: [
                { text: 'Session Name', value: 'serviceName', width: '20%', sortable: false },
                { text: 'No. of users', value: 'totalUsers', width: '20%', sortable: true },
                { text: 'No. of Requests', value: 'totalRequests', width: '20%', sortable: true },
                { text: 'No. of Errors', value: 'totalErrors', width: '20%', sortable: true },
                { text: 'Session Count', value: 'totalSessions', width: '20%', sortable: true }
            ]
        }
    },

    mounted() {
        this.get();
        this.$events.watch(this, this.get, 'refresh');
    },

    methods: {
        get() {
            this.loading = true;
            this.error = '';
            this.sessions = mockData.data;
            this.loading = false;
            // this.$api.getMRUMOverview((data, error) => {
            //     this.loading = false;
            //     if (error) {
            //         this.error = error;
            //         return;
            //     }
            //     this.sessions = data.sessions;
            // });
        },
        goToService(serviceName) {
            this.$router.push(this.link(serviceName));
        },
        link(serviceName) {
            return {
                name: 'overview',
                params: {
                    projectId: this.projectId,
                    view: 'MRUM',
                    id: serviceName,
                    tab: 'sessions'
                }
            };
        }
    },
};
</script>

<style scoped>
.table {
    margin-top: 20px;
    /* box-shadow: 1px 1px 5px 0 rgba(0, 0, 0, 0.1); */
}

.table th {
    font-weight: bold;
}

.table tr {
    color: #013912;
}

.clickable {
    cursor: pointer;
    text-decoration: none;
}
.clickable:hover {
    text-decoration: underline;
}
</style>
