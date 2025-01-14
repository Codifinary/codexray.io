<template>
    <v-container>
        <v-card>
            <v-card-title>Page Performance</v-card-title>
            <v-data-table :headers="headers" :items="data.pages" item-key="pagePath" class="elevation-1">
                <template v-slot:[`item.pagePath`]="{ item }">
                    <span>{{ item.pagePath }}</span>
                </template>
            </v-data-table>
        </v-card>
    </v-container>
</template>

<script>
export default {
    name: 'PagePerformance',
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
