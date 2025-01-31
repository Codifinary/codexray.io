<template>
    <div class="my-10 mx-5">
        <div v-if="selectedEventId">
            <ErrorDetail :eventId="selectedEventId" />
        </div>
        <div v-else>
            <CustomTable :headers="headers" :items="specificErrors" class="mt-10">
                <template v-slot:[`item.event_id`]="{ item }">
                    <router-link
                        :to="{
                            name: 'overview',
                            params: { view: 'EUM', id: $route.params.id, report: 'errors' },
                            query: { ...$utils.contextQuery(), error: encodeURIComponent(error), eventId: item.event_id },
                        }"
                        class="clickable"
                        @click.native.prevent="handleEventClick(item.event_id)"
                    >
                        {{ item.event_id }}
                    </router-link>
                </template>
            </CustomTable>
        </div>
    </div>
</template>

<script>
import CustomTable from '@/components/CustomTable.vue';
import ErrorDetail from './ErrorDetail.vue';

export default {
    components: { CustomTable, ErrorDetail },
    props: {
        id: { type: String, required: true },
        error: { type: String, required: true },
    },
    data() {
        return {
            specificErrors: [],
            selectedEventId: null,
            headers: [
                { text: 'Event ID', value: 'event_id' },
                { text: 'User ID', value: 'user_id' },
                { text: 'Device', value: 'device' },
                { text: 'OS', value: 'os' },
                { text: 'Browser', value: 'browser' },
                { text: 'Last Reported Time', value: 'last_reported' },
            ],
        };
    },
    created() {
        this.get(this.error);
    },
    watch: {
        error: {
            immediate: true,
            handler(newError) {
                this.get(newError);
            },
        },
        '$route.query.eventId': {
            immediate: true,
            handler(newEventId) {
                this.selectedEventId = newEventId;
            },
        },
    },
    methods: {
        get(error) {
            this.$api.getSpecificErrors(this.id, error, (data, Error) => {
                if (Error) {
                    console.error(Error);
                    return;
                }
                this.specificErrors = data.errors || [];
            });
        },
        handleEventClick(eventId) {
            this.$router.push({
                name: 'overview',
                params: { view: 'EUM', id: this.id, report: 'errors' },
                query: { ...this.$utils.contextQuery(), error: this.error, eventId },
            });
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
