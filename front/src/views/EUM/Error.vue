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
                <template #item.last_reported="{ item }">
                    {{ $format.date(item.last_reported, '{MMM} {DD}, {HH}:{mm}:{ss}') }}
                    ({{ $format.timeSinceNow(new Date(item.last_reported).getTime()) }} ago)
                </template>
            </CustomTable>
        </div>
    </div>
</template>

<script>
import CustomTable from '@/components/CustomTable.vue';
import ErrorDetail from './ErrorDetail.vue';

export default {
    components: {
        CustomTable,
        ErrorDetail,
    },
    props: {
        error: {
            type: String,
            required: true,
        },
        id: {
            type: String,
            required: true,
        },
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
                { text: 'Last reported time', value: 'last_reported' },
            ],
            localError: this.error,
        };
    },
    watch: {
        error: {
            immediate: true,
            handler(newError) {
                this.localError = newError;
                this.get(this.id, newError);
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
        get(id, error) {
            this.$api.getSpecificErrors(id, error, (data, Error) => {
                this.loading = false;
                if (Error) {
                    this.localError = Error;
                    return;
                }
                this.specificErrors = data.errors || [];
            });
        },
        handleEventClick(eventId) {
            this.$router.push({
                name: 'overview',
                params: { view: 'EUM', id: this.id },
                query: { ...this.$utils.contextQuery(), error: encodeURIComponent(this.error), eventId },
            });
        },
    },
    mounted() {
        this.$events.watch(this, this.get(this.id, this.error), 'refresh');
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
