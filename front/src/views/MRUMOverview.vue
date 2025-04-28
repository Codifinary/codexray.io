<template>
    <div>
        <v-alert v-if="error" color="red" icon="mdi-alert-octagon-outline" outlined text>
            {{ error }}
        </v-alert>

        <CustomTable 
            v-if="sessions.summary" 
            :headers="headers" 
            :items="sessions.summary" 
            class="table"
        >
            <template #item.service="{ item: { service } }">
                <div class="service-name clickable">
                    <router-link
                            :to="{
                                name: 'overview',
                                params: { view: 'MRUM', id: service },
                                query: $route.query,
                            }"
                        >
                            {{ service }}
                        </router-link>
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
export default {

    props: {
        id: String,
        report: String
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
                { text: 'Service Name', value: 'service', width: '20%', sortable: false },
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
            this.$api.getMRUMOverview((data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.sessions = data;
            });
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
