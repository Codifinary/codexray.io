<template>
    <div>
        <CustomTable :headers="headers" :items="specificErrors" class="mt-10" />
    </div>
</template>

<script>
import CustomTable from '@/components/CustomTable.vue';
import { getSpecificErrors } from './api/EUMapi';

export default {
    components: {
        CustomTable,
    },
    data() {
        return {
            specificErrors: null,
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
    created() {
        this.specificErrors = getSpecificErrors();
    },
};
</script>
