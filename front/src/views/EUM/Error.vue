<template>
    <div>
        <div v-if="selectedEventId">
            <ErrorDetail :eventId="selectedEventId" />
        </div>
        <div v-else>
            <CustomTable :headers="headers" :items="specificErrors" class="mt-10">
                <template v-slot:[`item.eventId`]="{ item }">
                    <router-link
                        :to="{
                            name: 'overview',
                            params: { view: 'EUM', id: id },
                            query: { ...$utils.contextQuery(), error: encodeURIComponent(error), eventId: item.eventId },
                        }"
                        class="clickable"
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
        '$route.query.error': {
            immediate: true,
            handler(newError) {
                if (newError !== this.error) {
                    this.$emit('update:error', newError);
                }
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
    },
    created() {
        this.fetchSpecificErrors(this.error);
    },
};
</script>

<style scoped>
.clickable {
    cursor: pointer;
    color: blue;
    text-decoration: underline;
}
</style>
