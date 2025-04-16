<template>
    <div>
        <v-card class="card-body">
            <div class="card-item">
                <div class="card-content">
                    <div class="card-info">
                        <div class="card-name">Total Applications</div>
                        <v-card-text class="card-count">{{ data.totalApplications }}</v-card-text>
                        <div class="card-name">Total pages</div>
                        <v-card-text class="card-sub-count">{{ data.totalPages }}</v-card-text>
                    </div>
                    <div class="trend-info">
                        <div class="card-name">Avg latency</div>
                        <v-card-text class="card-latency-count"
                            >{{ avgLatency.value }} <span class="card-name"> {{ avgLatency.unit }}</span></v-card-text
                        >
                    </div>
                </div>
                <hr class="separator mb-3" />
                <div class="card-content">
                    <div class="card-info">
                        <div class="card-name">Total errors</div>
                        <v-card-text class="card-count2">{{ data.totalError }}</v-card-text>
                        <div class="card-name">Errors/sec</div>
                        <v-card-text class="card-sub-count">{{ data.errorPerSec }}</v-card-text>
                    </div>
                    <div class="trend-info" v-if="!isFromNowQuery">
                        <span :style="{ color: trend === 'upTrend' ? 'green' : 'red' }" class="trend-percentage">
                            {{ parseFloat(data.errorTrend).toFixed(2) }}%
                        </span>
                        <img :src="`${$codexray.base_path}static/img/tech-icons/${trend}.svg`" class="card-icon" alt="Trend Icon" />
                    </div>
                </div>
            </div>
        </v-card>
    </div>
</template>

<script>
export default {
    props: {
        data: {
            type: Object,
            required: true,
        },
    },
    computed: {
        trend() {
            return this.data.errorTrend > 0 ? 'downTrend' : 'upTrend';
        },
        avgLatency() {
            const count = this.$format.convertLatency(this.data.avgLatency);
            return {
                value: Number.isInteger(count.value) ? count.value : parseFloat(count.value).toFixed(2),
                unit: count.unit,
            };
        },
        isFromNowQuery() {
            return this.$route.query.from && this.$route.query.from.startsWith('now-');
        },
    },
};
</script>

<style scoped>
.card-body {
    display: flex;
    flex-direction: column;
    width: 300px;
    padding: 7px 20px;
}

.card-item {
    display: flex;
    flex-direction: column;
    width: 100%;
}

.card-content {
    display: flex;
    align-items: center;
    justify-content: space-between;
}
.trend-percentage {
    opacity: 0.8;
    font-weight: 300;
    font-size: 14px;
}
.card-info {
    align-items: center;
}

.card-name {
    font-weight: 300;
    font-size: 12px;
    color: black;
    opacity: 0.5;
    margin-left: 10px;
}

.card-count {
    font-weight: 700;
    font-size: 30px;
    color: #013912;
}
.card-count2 {
    font-weight: 600;
    font-size: 26px;
    color: #013912;
}
.card-sub-count {
    font-weight: 600;
    font-size: 24px;
    color: #013912;
}
.card-latency-count {
    font-weight: 600;
    font-size: 20px;
    color: #013912;
}

.trend-info {
    display: flex;
    flex-direction: column;
    align-items: center;
}

.card-icon {
    width: 80px;
    height: 70px;
    margin-left: 10px;
    margin-top: 10px;
}

.separator {
    border: 1px solid grey;
    opacity: 0.5;
    width: 90%;
}
</style>
