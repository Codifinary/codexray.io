<template>
    <div class="my-10 mx-5">
        <CustomTable :headers="headers" :items="tableItems" item-key="serviceName" class="elevation-1">
            <template v-slot:item.serviceName="{ item }">
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
                            params: { view: 'EUM', id: item.serviceName },
                        }"
                    >
                        {{ item.serviceName }}
                    </router-link>
                </div>
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
    name: 'EUM',
    components: {
        CustomTable,
    },
    data() {
        return {
            headers: [
                { text: 'Application', value: 'serviceName' },
                { text: 'Number of Pages', value: 'pages' },
                { text: 'Load (Total Requests)', value: 'requests' },
                { text: 'Response time', value: 'avgLoadPageTime' },
                { text: 'JS Error%', value: 'jsErrorPercentage' },
                { text: 'API Error%', value: 'apiErrorPercentage' },
                { text: 'Users Impacted', value: 'impactedUsers' },
            ],
            tableItems: [],
            selectedApplication: null,
            loading: false,
            error: '',
        };
    },
    mounted() {
        this.get();
        this.$events.watch(this, this.get, 'refresh');
    },
    methods: {
        get() {
            this.loading = true;
            this.error = '';
            this.$api.getEUMApplications((data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.tableItems = data.eumapps.overviews || [];
                this.tableItems.forEach((item) => {
                    item.applicationType = item.browser === '' ? 'mobileapp' : 'webapp';
                });
            });
        },
        format(duration, unit) {
            return `${duration.toFixed(2)} ${unit}`;
        },
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
