<template>
    <v-card class="card-body" :style="{ '--line-color': lineColor }">
        <v-card-title class="card-title">
            <div class="card-name">{{ name }}</div>
            <v-card-text class="card-count">
                {{ formattedCount }}<span v-if="unit">{{ unit }}</span>
                <v-icon v-if="trendIcon" :color="trendIcon ? 'success' : 'error'">{{ trendIcon ? 'mdi-arrow-up' : 'mdi-arrow-down' }}</v-icon>
            </v-card-text>
        </v-card-title>
        <v-sparkline
            v-if="trend && trend !== 0"
            :value="sparklineData"
            fill
            smooth
            line-width="3"
            padding="0"
            :color="lineColor"
            :gradient="[lineColor, 'rgba(255,255,255,0)']"
            auto-draw
            height="500"
            width="400"
            stroke-linecap="round"
            :min="sparklineData.length ? Math.min(...sparklineData.filter((v) => v !== null)) : 0"
            :max="sparklineData.length ? Math.max(...sparklineData) : 0"
            class="sparkline"
        />
        <BaseIcon
            v-if="trend === undefined || trend === null"
            :name="iconName || 'alert'"
            :iconColor="icon"
            :class="['card-icon', background]"
            style="border-radius: 30%"
        />

        <div class="bottom-border"></div>
    </v-card>
</template>

<script>
import BaseIcon from './BaseIcon.vue';

export default {
    components: {
        BaseIcon,
    },
    props: {
        name: String,
        count: Number,
        background: String,
        icon: String,
        iconName: String,
        unit: String,
        lineColor: String,
        trend: Number,
        trendIcon: Boolean,
    },
    computed: {
        formattedCount() {
            return Number.isInteger(this.count) ? this.count : this.count.toFixed(2);
        },
        sparklineData() {
            const length = 10;
            const direction = this.trend >= 0 ? 1 : -1;
            const base = Math.abs(this.trend);
            const step = base / length;

            let value = 0;
            return Array.from({ length }, () => {
                const fluctuation = (Math.random() - 0.5) * step * 1.5;
                value += direction * step + fluctuation;
                return parseFloat(value.toFixed(2));
            });
        },
    },
};
</script>

<style scoped>
.sparkline {
    height: 100%;
    box-sizing: border-box;
    width: 30%;
    padding-right: 10px;
}

.bottom-border {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 4px;
    background: var(--line-color);
    border-bottom-left-radius: 8px;
    border-bottom-right-radius: 8px;
}

.card-body {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 20%;
    height: 100px;
    gap: 20px;
}

.card-name {
    font-weight: 400;
    font-size: 12px;
    line-height: 16px;
    margin-bottom: 3px;
    color: #013912;
}

.card-title {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    /* width: fit-content; */
    padding: 16px;
    padding-left: 16px;
    padding-right: 0;
}

.card-count {
    font-weight: 700;
    font-size: 24px;
    padding: 0;
    line-height: 33px;
    color: #013912;
}

.card-icon {
    margin-right: 10px;
    padding: 8px 10px 6px 10px;
}

.card-count span {
    font-size: 18px;
    color: gray;
    font-weight: 500;
}

@media (min-width: 1441px) {
    /* Styles for larger monitor screens */
    .card-body {
        width: 19%;
        height: 150px;
        padding-left: 20px;
        gap: 20px;
    }

    .card-name {
        font-size: 12px;
        line-height: 18px;
    }

    .card-count {
        font-weight: 800;
        font-size: 36px;
        line-height: 38px;
    }

    .v-icon {
        font-size: 30px;
    }
}
</style>
