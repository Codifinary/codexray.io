<template>
    <div>
        <CustomTable :headers="headers" :items="data" item-key="pagePath" class="elevation-1 mt-10">
            <template v-slot:[`item.pagePath`]="{ item }">
                <router-link
                    :to="{
                        name: 'overview',
                        params: { view: 'EUM', id: id, report: 'page-performance' },
                        query: { pagePath: item.pagePath },
                    }"
                    class="clickable"
                >
                    {{ item.pagePath }}
                </router-link>
            </template>
        </CustomTable>
    </div>
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
            type: Array,
            required: true,
        },
        id: {
            type: String,
            required: true,
        },
    },
    data() {
        return {
            headers: [
                { text: 'Page', value: 'pagePath' },
                { text: 'Load (Total Requests)', value: 'requests' },
                { text: 'Response time', value: 'avgLoadPageTime' },
                { text: 'JS Error%', value: 'jsErrorPercentage' },
                { text: 'API Error%', value: 'apiErrorPercentage' },
                { text: 'Users Impacted', value: 'impactedUsers' },
            ],
        };
    },
};
</script>

<style scoped>
.clickable {
    cursor: pointer;
    color: var(--status-ok);
    text-decoration: underline;
}
</style>
