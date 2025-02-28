<template>
    <div>
        <v-card class="card-body">
            <div v-for="(item, index) in data" :key="index" class="card-item">
                <div class="card-content">
                    <div class="card-info">
                        <div class="card-name">{{ item.name }}</div>
                        <v-card-text class="card-count">{{ item.count }}</v-card-text>
                        <div class="card-name">{{ item.subName }}</div>
                        <v-card-text class="card-sub-count">{{ item.subCount }}</v-card-text>
                    </div>
                    <div class="trend-info">
                        <span
                            :style="{ color: item.trend === 'upTrend' ? 'green' : item.trend === 'downTrend' ? 'red' : 'black' }"
                            class="trend-percentage"
                        >
                            {{ item.trendPercentage }}
                        </span>
                        <img :src="`${$codexray.base_path}static/img/tech-icons/${item.trend}.svg`" class="card-icon" alt="Trend Icon" />
                    </div>
                </div>
                <hr v-if="index < data.length - 1" class="separator" />
            </div>
        </v-card>
    </div>
</template>

<script>
export default {
    props: {
        data: {
            type: Array,
            required: true,
        },
    },
};
</script>

<style scoped>
.card-body {
    display: flex;
    flex-direction: column;
    width: 100%;
    padding: 20px;
}

.card-item {
    display: flex;
    flex-direction: column;
    width: 100%;
    padding: 20px 0;
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
.card-sub-count {
    font-weight: 600;
    font-size: 24px;
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
