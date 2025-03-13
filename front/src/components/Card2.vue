<template>
  <div class="card" :style="{ '--bottom-color': cardData.bottomColor || '#059669' }">
    <div class="metrics-container">
      <div class="main-metric">
        <div class="metric-label">{{ cardData.name || 'Total Requests' }}</div>
        <div class="metric-value">{{ cardData.totalRequests || 0 }}</div>
      </div>
      <div class="secondary-metric">
        <div class="metric-label">{{ cardData.secondaryLabel || 'Req/sec' }}</div>
        <div class="metric-value">{{ cardData.secondaryValue || 0 }}</div>
      </div>
    </div>
    <div class="percentage-indicator">
      <div class="percentage" :class="{ 'positive': (cardData.percentageChange || 0) > 0, 'negative': (cardData.percentageChange || 0) < 0 }">
        {{ (cardData.percentageChange || 0) > 0 ? '+' : '' }}{{ (cardData.percentageChange || 0).toFixed(2) }}%
      </div>
      <BaseIcon 
        :name="cardData.icon || 'up-arrow'" 
        :iconColor="cardData.iconColor || 'green'" 
        :class="['card-icon', cardData.background || '']" 
        style="border-radius: 30%" 
      />
    </div>
    <div class="bottom-border"></div>
  </div>
</template>

<script>
import BaseIcon from './BaseIcon.vue';

export default {
  components: {
    BaseIcon,
  },
  props: {
    cardData: {
      type: Object,
      default: () => ({
        name: 'Total Requests',
        totalRequests: 0,
        secondaryLabel: 'Req/sec',
        secondaryValue: 0,
        percentageChange: 0,
        icon: 'up-arrow',
        iconColor: 'green',
        background: '',
        bottomColor: '#059669'
      })
    }
  }
};
</script>

<style scoped>
.card {
  background: #FFFFFF;
  border-radius: 8px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  width: 280px;
  position: relative;
}

.metrics-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.main-metric, .secondary-metric {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.metric-label {
  color: #666666;
  font-size: 14px;
  font-weight: 400;
}

.main-metric .metric-value {
  color: #000000;
  font-size: 32px;
  font-weight: 600;
  line-height: 1;
}

.secondary-metric .metric-value {
  color: #000000;
  font-size: 24px;
  font-weight: 600;
  line-height: 1;
}

.percentage-indicator {
  position: absolute;
  top: 24px;
  right: 24px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.percentage {
  font-size: 14px;
  font-weight: 500;
}

.percentage.positive {
  color: #34D399;
}

.percentage.negative {
  color: #EF4444;
}

.triangle-icon {
  width: 16px;
  height: 16px;
  filter: invert(72%) sepia(40%) saturate(463%) hue-rotate(95deg) brightness(91%) contrast(91%);
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

@media (min-width: 1440px) {
  .card {
    width: 320px;
    padding: 28px;
  }

  .main-metric .metric-value {
    font-size: 36px;
  }

  .secondary-metric .metric-value {
    font-size: 28px;
  }
}
</style>