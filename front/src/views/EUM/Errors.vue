<template>
    <div class="my-10 mx-5">
        <CustomTable :headers="headers" :items="data.errors" item-key="error" class="elevation-1">
            <template v-slot:[`item.error`]="{ item }">
                <router-link
                    class="clickable"
                    :to="{
                        name: 'overview',
                        params: { view: 'EUM', id: $route.params.id },
                        query: { ...$utils.contextQuery(), error: encodeURIComponent(item.error) },
                    }"
                    @click.native.prevent="handleErrorClick(item.error)"
                >
                    <span>{{ item.error }}</span>
                </router-link>
            </template>
        </CustomTable>
    </div>
</template>

<script>
import CustomTable from '@/components/CustomTable.vue';

export default {
    components: {
        CustomTable,
    },
    name: 'Errors',
    props: {
        data: {
            type: Object,
            required: true,
        },
    },
    data() {
        return {
            headers: [
                { text: 'Error', value: 'error' },
                { text: 'Number of Events', value: 'numberOfEvents' },
                { text: 'Users Impacted', value: 'usersImpacted' },
                { text: 'Last Reported Time', value: 'lastReportedTime' },
                { text: 'Category', value: 'category' },
            ],
        };
    },
    watch: {
        '$route.query.error': {
            immediate: true,
            handler(newError) {
                this.$emit('update:error', newError);
            },
        },
        '$route.query.eventId': {
            immediate: true,
            handler(newEventId) {
                this.$emit('update:eventId', newEventId);
            },
        },
    },
    methods: {
        handleErrorClick(error) {
            this.$emit('error-clicked', error);
        },
    },
};
</script>
<style scoped>
.clickable {
    cursor: pointer;
    color: var(--status-ok);
    text-decoration: none !important;
}
</style>
