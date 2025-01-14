<template>
    <v-container class="my-10">
        <v-card>
            <v-card-title>Errors</v-card-title>
            <v-data-table :headers="headers" :items="data.errors" item-key="error" class="elevation-1">
                <template v-slot:[`item.error`]="{ item }">
                    <span>{{ item.error }}</span>
                </template>
            </v-data-table>
        </v-card>
    </v-container>
</template>

<script>
export default {
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
            tableItems: [],
        };
    },
    watch: {
        errors: {
            immediate: true,
            handler(newVal) {
                if (newVal && newVal.errors) {
                    this.tableItems = newVal.errors;
                }
            },
        },
    },
};
</script>

<style scoped></style>
