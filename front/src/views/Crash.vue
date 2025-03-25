<template>
    <div class="crash-container">
        <Card :name="name" :count="count" :background="background" :icon="icon" :iconName="iconName" :bottomColor="bottomColor"/>
        <CustomTable :headers="headers" :items="items"/>
        <Dashboard id="chart" :name="title" :widgets="chartData.widgets" />
    </div>
</template>

<script>
import Card from '@/components/Card.vue';
import CustomTable from '@/components/CustomTable.vue';

export default {
    components: {
        Card,
        CustomTable
    },
    data() {
        return {
            name: 'Crash',
            count: 0,
            background: 'light-red-bg',
            icon: 'red-icon',
            bottomColor: '#EF5350',
            iconName: 'alert',
            headers: [
                { text: 'Crash reason', value: 'CrashReason', sortable: false },
                { text: 'Total Crashes', value: 'Crashes', sortable: true },
                { text: 'Affected Users', value: 'AffectedUsers', sortable: true },
                { text: 'Last Occurrence', value: 'LastOccurrence', sortable: true },
            ],
            items: [],
        };
    },
    methods: {
        get(){
            this.$api.getCrashData((data, error) => {
                if (error) {
                    this.error = error;
                    return;
                }
                this.items = data;
            });
        }
    }
};
</script>

<style scoped>


.crash-container{
    margin: 20px;
}

.cards {
    display: flex;
    flex-wrap: wrap;
    gap: 50px;
    margin-right: 30px;  
    margin-bottom: 50px;
    margin-top: 50px;
}

.table {
    margin-bottom: 50px;
    margin-top: 50px;
}

.light-green-bg {
    background-color: rgba(5, 150, 105, 0.1);
}

.light-red-bg {
    background-color: rgba(220, 38, 38, 0.1);
}

.light-orange-bg {
    background-color: rgba(249, 115, 22, 0.1);
}

.geomap{
    margin-top: 50px;
}
</style>
