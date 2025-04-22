<template>
    <div>
        <!-- Requests Card -->
        <v-card class="card-body green-card">
            <div class="card-item">
                <div class="card-content">
                    <div class="card-info">
                        <div class="card-name">Total Requests</div>
                        <v-card-text class="card-count2">{{ data.totalRequests }}</v-card-text>
                        <div class="card-name">Requests/sec</div>
                        <v-card-text class="card-sub-count">{{ rps }}</v-card-text>
                    </div>
                    <div class="trend-info" v-if="!isFromNowQuery">
                        <span class="trend-percentage">{{ parseFloat(data.requestTrend).toFixed(2) }}%</span>
                        <img
                            :src="`${$codexray.base_path}static/img/tech-icons/upArrowGreen.svg`"
                            class="card-icon"
                            :class="requestTrend"
                            alt="Trend Icon"
                        />
                    </div>
                </div>
            </div>
        </v-card>

        <!-- Errors Card -->
        <v-card class="card-body red-card">
            <div class="card-item">
                <div class="card-content">
                    <div class="card-info">
                        <div class="card-name">Total Errors</div>
                        <v-card-text class="card-count2">{{ data.totalErrors }}</v-card-text>
                        <div class="card-name">Errors/sec</div>
                        <v-card-text class="card-sub-count">{{ rps }}</v-card-text>
                    </div>
                    <div class="trend-info" v-if="!isFromNowQuery">
                        <span class="trend-percentage red">{{ parseFloat(data.errorTrend).toFixed(2) }}%</span>
                        <img
                            :src="`${$codexray.base_path}static/img/tech-icons/upArrowRed.svg`"
                            class="card-icon"
                            :class="errorTrend"
                            alt="Trend Icon"
                        />
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
        rps() {
            const rps = this.data.requestsPerSecond || 0;
            return this.$format.float(rps);
        },
        requestTrend() {
            return this.data.requestTrend >= 0 ? 'upTrend' : 'downTrend';
        },
        errorTrend() {
            return this.data.errorTrend >= 0 ? 'upTrend' : 'downTrend';
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
    margin-top: 10px;
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
    padding: 5px 10px;
    border-radius: 10px;
    background-color: rgba(0, 128, 0, 0.616);
    color: white;
    text-align: center;
}
.trend-percentage.red {
    background-color: rgba(255, 0, 0, 0.1);
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
    margin-top: 5px;
}

.card-count2 {
    font-weight: 600;
    font-size: 30px;
    color: #013912;
}

.card-sub-count {
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
    transition: transform 0.3s ease;
}

.downTrend {
    transform: scaleY(-1);
}

.green-card {
    background-color: #edffed !important;
}

.red-card {
    background-color: #ffeded !important;
}
</style>
