<template>
    <div class="my-10 mx-5">
        <div v-if="selectedEventId">
            <ErrorDetail :eventId="selectedEventId" />
        </div>
        <div v-else>
            <CustomTable :headers="headers" :items="specificErrors" class="mt-10">
                <template v-slot:[`item.eventId`]="{ item }">
                    <router-link
                        :to="{
                            name: 'overview',
                            params: { view: 'EUM', id: $route.params.id, report: 'errors' },
                            query: { ...$utils.contextQuery(), error: encodeURIComponent(error), eventId: item.eventId },
                        }"
                        class="clickable"
                        @click.native.prevent="handleEventClick(item.eventId)"
                    >
                        {{ item.eventId }}
                    </router-link>
                </template>
            </CustomTable>
        </div>
    </div>
</template>

<script>
import CustomTable from '@/components/CustomTable.vue';
import ErrorDetail from './ErrorDetail.vue';
import { getSpecificErrors } from './api/EUMapi';

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
                { text: 'Event ID', value: 'eventId' },
                { text: 'User ID', value: 'userId' },
                { text: 'Device', value: 'device' },
                { text: 'OS', value: 'os' },
                { text: 'Browser', value: 'browserAndVersion' },
                { text: 'Last reported time', value: 'lastReportedTime' },
            ],
        };
    },
    watch: {
        error: {
            immediate: true,
            handler(newError) {
                this.fetchSpecificErrors(newError);
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
        fetchSpecificErrors(error) {
            const allErrors = getSpecificErrors(this.id, error);
            if (allErrors) {
                this.specificErrors = allErrors;
            } else {
                console.error('No specific errors found');
            }
        },
        handleEventClick(eventId) {
            this.$router.push({
                name: 'overview',
                params: { view: 'EUM', id: this.id },
                query: { ...this.$utils.contextQuery(), error: encodeURIComponent(this.error), eventId },
            });
        },
    },
    created() {
        this.fetchSpecificErrors(this.error);
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
