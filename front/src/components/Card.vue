<template>
    <v-card class="card-body" :style="{ '--bottom-color': bottomColor }">
        <v-card-title>
            <div class="card-name">{{ name }}</div>
            <v-card-text class="card-count">
                {{ formattedCount }}<span v-if="unit">{{ unit }}</span>
            </v-card-text>
        </v-card-title>
        <BaseIcon 
            v-if="icon" 
            :name="iconName" 
            :iconColor="iconColor || icon" 
            :class="['card-icon', background]" 
            style="border-radius: 30%" 
            width="24px"
            height="24px"
        />
        <v-sparkline
    v-else-if="trend && trend.chart"
    :value="trend.chart.map((v) => (v === null ? 0 : v))"
    fill
    smooth
    line-width="3"  
    padding="8"
    :color="bottomColor"  
    :gradient="[bottomColor, 'rgba(255,255,255,0)']"  
    auto-draw
    height="50"
    width="70"
    stroke-linecap="round"
    :min="trend.chart.length ? Math.min(...trend.chart.filter(v => v !== null)) : 0" 
    :max="trend.chart.length ? Math.max(...trend.chart) : 0"
/>


        <div class="bottom-border" :class="$vuetify.theme.dark ? 'theme--dark' : 'theme--light'"></div>

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
        bottomColor: String,
        trend: {
            type: Object,
            required: false
        },
        iconColor: String,
    },
    computed: {
        formattedCount() {
            return Number.isInteger(this.count) ? this.count : this.count.toFixed(2);
        },
    },
};
</script>

<style scoped>
.card-body {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 18%;
    height: 100px;
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
