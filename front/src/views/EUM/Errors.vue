<template>
    <div class="my-10 mx-5">
        <CustomTable :headers="headers" :items="data" item-key="error_name" class="elevation-1">
            <template v-slot:[`item.error_name`]="{ item }">
                <router-link
                    class="clickable"
                    :to="{
                        name: 'overview',
                        params: { view: 'EUM', id: $route.params.id, report: 'errors' },
                        query: { ...$utils.contextQuery(), error: encodeURIComponent(item.error_name) },
                    }"
                    @click.native.prevent="handleErrorClick(item.error)"
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
            type: Array,
            required: true,
        },
    },
    data() {
        return {
            headers: [
                { text: 'Error', value: 'error_name' },
                { text: 'Number of Events', value: 'event_count' },
                { text: 'Users Impacted', value: 'user_impacted' },
                { text: 'Last Reported Time', value: 'last_reported' },
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
    },
    methods: {
        handleErrorClick(error) {
            this.$emit('error-clicked', error);
            this.$router.push({
                name: 'overview',
                params: { view: 'EUM', id: this.$route.params.id },
                query: { ...this.$utils.contextQuery(), error: encodeURIComponent(error) },
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
