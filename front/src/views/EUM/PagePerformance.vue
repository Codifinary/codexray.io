<template>
    <v-container>
        <CustomTable :headers="headers" :items="data.pages" item-key="pagePath" class="elevation-1 mt-10">
            <template v-slot:[`item.pagePath`]="{ item }">
                <span>{{ item.pagePath }}</span>
            </template>
        </CustomTable>
    </v-container>
</template>

<script>
import CustomTable from '@/components/CustomTable.vue';
export default {
    name: 'PagePerformance',
    components: {
        CustomTable,
    },
    props: {
        data: {
            type: Object,
            required: true,
        },
    },
    data() {
        return {
            headers: [
                { text: 'Page', value: 'pagePath' },
                { text: 'Load Requests/Second', value: 'loadRequestsPerSecond' },
                { text: 'Response Time (ms)', value: 'responseTimeMs' },
                { text: 'JS Error%', value: 'jsErrorPercentage' },
                { text: 'API Error%', value: 'apiErrorPercentage' },
                { text: 'Users Impacted', value: 'usersImpacted' },
            ],
        };
    },
};
</script>

<style scoped></style>
