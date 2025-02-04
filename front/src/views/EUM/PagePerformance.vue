<template>
    <div>
        <CustomTable :headers="headers" :items="pagePerformance" item-key="pagePath" class="elevation-1 mt-10">
            <template v-slot:[`item.pagePath`]="{ item }">
                <router-link
                    :to="{
                        name: 'overview',
                        params: { view: 'EUM', id: id, report: 'page-performance' },
                        query: { ...$route.query, pagePath: item.pagePath },
                    }"
                    class="clickable"
                >
                    {{ item.pagePath }}
                </router-link>
            </template>
            <template v-slot:item.avgLoadPageTime="{ item }">
                {{ format(item.avgLoadPageTime, 'ms') }}
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
        id: {
            type: String,
            required: true,
        },
    },
    data() {
        return {
            pagePerformance: [],
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
    methods: {
        get(id) {
            this.loading = true;
            this.error = '';
            this.$api.getPagePerformance(id, (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.pagePerformance = data.overviews || [];
            });
        },
        format(duration, unit) {
            return `${duration.toFixed(2)} ${unit}`;
        },
    },
    created() {
        this.get(this.id);
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
