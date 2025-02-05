<template>
    <div v-if="!selectedError" class="my-10 mx-5">
        <CustomTable :headers="headers" :items="errors" item-key="error_name" class="elevation-1">
            <template v-slot:[`item.error_name`]="{ item }">
                <router-link
                    class="clickable"
                    :to="{
                        name: 'overview',
                        params: { view: 'EUM', id: id, report: `errors/${item.error_name}` },
                        query: { ...$utils.contextQuery() },
                    }"
                    @click.native.prevent="handleErrorClicked(item.error_name)"
                >
                    <span>{{ item.error_name }}</span>
                </router-link>
            </template>
            <template #item.last_reported="{ item }">
                {{ $format.date(item.last_reported, '{MMM} {DD}, {HH}:{mm}:{ss}') }}
                ({{ $format.timeSinceNow(new Date(item.last_reported).getTime()) }} ago)
            </template>
        </CustomTable>
    </div>
    <div v-else>
        <Error :error="selectedError" :id="id" :report="'errors'" @event-clicked="handleEventClicked" />
    </div>
</template>

<script>
import CustomTable from '@/components/CustomTable.vue';
import Error from './Error.vue';

export default {
    components: { CustomTable, Error },
    name: 'Errors',
    props: { id: { type: String, required: true } },
    data() {
        return {
            headers: [
                { text: 'Error', value: 'error_name' },
                { text: 'Number of Events', value: 'event_count' },
                { text: 'Users Impacted', value: 'user_impacted' },
                { text: 'Last Reported Time', value: 'last_reported' },
                { text: 'Category', value: 'category' },
            ],
            errors: [],
            selectedError: null,
        };
    },
    watch: {
        '$route.query.error': {
            immediate: true,
            handler(newError) {
                this.selectedError = newError;
            },
        },
    },
    methods: {
        get() {
            this.$api.getEUMApplicationErrors(this.id, (data, error) => {
                if (error) {
                    console.error('Error fetching errors:', error);
                    return;
                }
                this.errors = data.errors || [];
            });
        },
        handleErrorClicked(error) {
            this.selectedError = error;
            this.$router
                .push({
                    name: 'overview',
                    params: { view: 'EUM', id: this.id, report: `errors` },
                    query: { ...this.$utils.contextQuery(), error: this.selectedError },
                })
                .catch((err) => {
                    if (err.name !== 'NavigationDuplicated') console.error(err);
                });
        },
        handleEventClicked(eventId) {
            this.$router.push({
                name: 'overview',
                params: { view: 'EUM', id: this.id, report: `errors` },
                query: { ...this.$utils.contextQuery(), error: this.selectedError, eventId },
            });
        },
    },
    mounted() {
        this.get();
        this.$events.watch(this, this.get, 'refresh');
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
