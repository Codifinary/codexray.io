<template>
    <v-card class="card-body" :style="{ '--bottom-color': bottomColor }">
        <div class="card-content">
            <v-card-title class="card-title">
                <div class="card-name">{{ name }}</div>
                <v-card-text class="card-count d-flex align-center">
                    {{ formattedCount }}<span v-if="unit">{{ unit }} </span>
                    <v-icon v-if="trendIcon" color="success">mdi-arrow-up</v-icon>
                </v-card-text>
            </v-card-title>
        </div>
        <BaseIcon v-if="icon" :name="iconName || 'alert'" :iconColor="icon" :class="['card-icon', background]" style="border-radius: 30%" />
        <v-sparkline
            v-if="trend"
            :value="sparklineData"
            fill
            smooth
            line-width="3"
            padding="4"
            :color="bottomColor"
            :gradient="[bottomColor, 'rgba(255,255,255,0)']"
            auto-draw
            height="400"
            width="350"
            stroke-linecap="round"
            :min="sparklineData.length ? Math.min(...sparklineData.filter((v) => v !== null)) : 0"
            :max="sparklineData.length ? Math.max(...sparklineData) : 0"
            class="sparkline"
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
        trendIcon: Boolean,
        trend: Number,
        bottomColor: String,
    },
    computed: {
        formattedCount() {
            if (typeof this.count === 'number') {
                return this.count.toFixed(0);
            }
            return this.count;
        },
        sparklineData() {
  const length = 10;
  const direction = this.trend >= 0 ? 1 : -1;
  const base = Math.abs(this.trend);
  const step = base / length;

  let value = 0;
  return Array.from({ length }, () => {
    // Add base trend with noise
    const fluctuation = (Math.random() - 0.5) * step * 1.5; // random between -0.75x to +0.75x of step
    value += direction * step + fluctuation;

    return parseFloat(value.toFixed(2)); // optional: round to 2 decimal places
  });
}

    },
};
</script>

<style scoped>
.card-body {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 25%;
    height: 100px;
    gap: 10px;
    position: relative;
    padding: 16px;
}

.card-content {
    width: 150px; /* Adjust as needed */
    flex-shrink: 0;
}


.card-name {
    font-weight: 400;
    font-size: 12px;
    line-height: 16px;
    margin-bottom: 3px;
    color: #013912;
    white-space: pre-line;
    word-break: break-word;
    max-width: 120px;
}

.card-count {
    font-weight: 700;
    font-size: 26px;
    padding: 0;
    line-height: 33px;
    color: #013912;
}

.card-title {
    width: 150px;
    flex-shrink: 0;
    display: flex;
    align-items: center;
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

.sparkline {
    margin-top: 20px;
    margin-bottom: 10px;
}

@media (min-width: 1441px) {
    /* Styles for larger monitor screens */
    .card-body {
        width: 19%;
        height: 150px;
        padding: 24px;
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

    .card-content {
        width: 140px;
    }

    /* .card-sparkline {
    width: 100px;
  } */
}

.bottom-border {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 4px;
    background: var(--bottom-color);
    border-bottom-left-radius: 8px;
    border-bottom-right-radius: 8px;
}
</style>