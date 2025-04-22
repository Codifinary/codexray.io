<template>
    <div class="card" :style="{ '--bottom-color': '#013912' }">
      <div class="metrics-container">
        <div class="main-metric">
          <div class="metric-label" :style="{ color: 'primary' }">{{ cardData.primaryLabel }}</div>
          <div class="metric-value">{{cardData.primaryValue}} </div>
        </div>
        <div class="secondary-metric">
          <div class="metric-label">{{ cardData.secondaryLabel}}</div>
          <div class="metric-value">{{cardData.secondaryValue}}</div>
        </div>
    </div>
      <div class="percentage-indicator">
        <div v-if="cardData.percentageChange !== undefined && cardData.percentageChange !== null" class="percentage" :class="{ 'positive': (cardData.trendColor === '#66BB6A'), 'negative': (cardData.trendColor === '#EF5350') }">
          {{ (cardData.percentageChange) > 0 ? '+' : '' }}{{ (cardData.percentageChange || 0).toFixed(2) }}%
        </div>
        <BaseIcon 
          v-if="cardData.icon"
          :name="cardData.icon" 
          :iconColor="cardData.iconColor || 'primary'" 
          :class="['card-icon', cardData.background || '']" 
          style="border-radius: 30%" 
        />
      </div>
      <div class="bottom-border" :class="$vuetify.theme.dark ? 'theme--dark' : 'theme--light'"></div>
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
        required: true,
        default: () => ({
          primaryLabel: '',
          primaryValue: 0,
          secondaryLabel: '',
          secondaryValue: 0,
          percentageChange: 0,
          icon: '',
          iconColor: '',
          trendColor: '',
        })
      }
    }
  };
  </script>
  
  <style scoped>
  .card {
    background: #FFFFFF;
    border-radius: 8px;
    padding: 1.5rem;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    width: 25%;
    position: relative;
  }

  .card-count {
    font-weight: 700;
    font-size: 26px;
    padding: 0;
    line-height: 33px;
    color: #013912;
}
  
  .metrics-container {
    display: flex;
    flex-direction: column;
    height: 100%;
    gap: 1rem;
    justify-content: center;
  }
  
  .main-metric, .secondary-metric {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }
  
  .metric-label {
    color: #013912;
    font-size: 0.875rem;
    font-weight: 400;
  }
  
  .main-metric .metric-value {
    color: #013912;
    font-size: 2rem;
    font-weight: 600;
    line-height: 1;
  }
  
  .secondary-metric .metric-value {
    color: #013912;
    font-size: 1.5rem;
    font-weight: 600;
    line-height: 1;
  }
  
  .percentage-indicator {
    position: absolute;
    top: 1.5rem;
    right: 1.5rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 3.75rem;
  }
  
  .percentage {
    font-size: 0.875rem;
    font-weight: 500;
  }
  
  .percentage.positive {
    color: #34D399;
  }
  
  .percentage.negative {
    color: #EF4444;
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
    width: 18vw;
    padding: 2%;
  }

  .main-metric .metric-value {
    font-size: 2.5vw;
  }

  .secondary-metric .metric-value {
    font-size: 1.8vw;
  }
}
  </style>