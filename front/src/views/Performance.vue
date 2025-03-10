<template>
    <div class="cards my-4 d-flex">
        <Card2
            name="Active Sessions"
            count="6158"
            sub="25%"
            background="success"
            icon="alert"
        />
        <Card2
            name="Users"
            count="514"
            sub="25%"
            background="info"
            icon="person"
        />
        <Card2
            name="Sessions"
            count="6158"
            sub="25%"
            background="warning"
            icon="clock"
        />

        <CustomTable :items="items" :headers="headers" class="table">
            <template #item.incident="{ item }">
                <div class="incident" :class="{ 'grey--text': item.resolved_at }">
                    <div class="status" :style="{ backgroundColor: item.color }" />
                    <router-link
                        :to="{
                            name: 'overview',
                            params: { view: 'incidents' },
                            query: { ...$utils.contextQuery(), incident: item.key },
                        }"
                    >
                        <span class="key" style="font-family: monospace">i-{{ item.key }}</span>
                    </router-link>
                </div>
            </template>

            <template #item.opened_at="{ item }">
                <div class="d-flex text-no-wrap" :class="{ 'grey--text': item.resolved_at }">
                    {{ $format.date(item.opened_at, '{MMM} {DD}, {HH}:{mm}:{ss}') }}
                    ({{ $format.timeSinceNow(item.opened_at) }} ago)
                </div>
            </template>

            <template #item.duration="{ item }">
                <div class="d-flex text-no-wrap" :class="{ 'grey--text': item.resolved_at }">
                    {{ $format.durationPretty(item.duration) }}
                </div>
            </template>

            <template #item.application="{ item }">
                <div class="d-flex">
                    <span :class="{ 'grey--text': item.resolved_at }">
                        {{ $utils.appId(item.application_id).name }}
                    </span>
                </div>
            </template>

            <template #item.latency="{ item }">
                <span v-if="item.latency_slo" :class="item.latency_slo.violated ? 'fired' : undefined">
                    {{ item.latency_slo.compliance }}
                </span>
            </template>

            <template #item.availability="{ item }">
                <span v-if="item.availability_slo" :class="item.availability_slo.violated ? 'fired' : undefined">
                    {{ item.availability_slo.compliance }}
                </span>
            </template>

            <template #item.affected_request_percent="{ item }">
                <div class="progress-line">
                    <v-progress-linear :value="item.affected_request_percent" color="green lighten-1" height="16px"> </v-progress-linear>
                    <div class="percent-text">{{ $format.percent(item.affected_request_percent) }} %</div>
                </div>
            </template>

            <template #item.error_budget_consumed_percent="{ item }">
                <div class="progress-line">
                    <v-progress-linear :value="item.error_budget_consumed_percent" color="purple lighten-1" height="16px"> </v-progress-linear>
                    <div class="percent-text">{{ $format.percent(item.error_budget_consumed_percent) }} %</div>
                </div>
            </template>
            <template #item.actions="{ item }">
                <v-menu offset-y>
                    <template v-slot:activator="{ attrs, on }">
                        <v-btn icon x-small class="ml-1" v-bind="attrs" v-on="on">
                            <v-icon small>mdi-dots-vertical</v-icon>
                        </v-btn>
                    </template>

                    <v-list dense>
                        <v-list-item @click="edit(item.application_id, 'SLOAvailability', 'Availability')">
                            <v-icon small class="mr-1">mdi-check-circle-outline</v-icon> Adjust Availability SLO
                        </v-list-item>
                        <v-list-item @click="edit(item.application_id, 'SLOLatency', 'Latency')">
                            <v-icon small class="mr-1">mdi-timer-outline</v-icon> Adjust Latency SLO
                        </v-list-item>
                        <v-list-item
                            :to="{
                                name: 'overview',
                                params: { view: 'incidents' },
                                query: { incident: item.key, view: 'rca' },
                            }"
                        >
                            <v-icon small class="mr-1">mdi-creation</v-icon> Investigate with AI
                        </v-list-item>
                    </v-list>
                </v-menu>
            </template>
        </CustomTable>
        <!-- <Card2 v-for="card in cards" 
              :key="card.name" 
              :name="card.name" 
              :count="card.count" 
              :background="card.background" 
              :icon="card.icon" 
              :sub="card.sub" /> -->
    </div>
</template>

<script>
import Card2 from '@/components/Card2.vue';
import CustomTable from '@/components/CustomTable.vue';

export default {
    props: {
        projectId: String,
        tab: String,
    },
    data() {
        return {
            cards: [
                { name: 'Active Sessions', count: 6158, sub: '25%', background: 'success', icon: 'alert' },
                { name: 'Users', count: 514, sub: '25%', background: 'info', icon: 'person' },
                { name: 'Sessions', count: 6158, sub: '25%', background: 'warning', icon: 'clock' },
            ],
        };
    },
    components: {
        Card2,
        CustomTable
    },
};
</script>

<style scoped>
.cards {
    display: flex;
    gap: 40px;
    align-items: center;
    justify-content: start;
    margin-left: 40px;
}

.table .incident {
    gap: 4px;
    display: flex;
}
.table .incident .status {
    height: 6px;
    width: 6px;
    border-radius: 50%;
    align-self: center;
}

.table .incident .key {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

/* .cards {
    display: flex;
    gap: 20px;
    align-items: center;
} */


.fired {
    color: #ef5350 !important;
    background-color: unset !important;
}

.progress-line {
    display: flex;
    align-items: center;
    gap: 4px;
}

.percent-text {
    font-size: 10px;
    color: rgba(0, 0, 0, 0.5);
    white-space: nowrap;
}
</style>
