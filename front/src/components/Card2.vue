<template>
  <v-card  class="card-container" :style="{ '--bottom-color': '#013912' }">
    <div class="outer-container d-flex">

    <div class="title-container">
      <div>
        <v-card-text class="card-title">{{ cardData.primaryLabel }}</v-card-text>
        <v-card-subtitle class="card-subtitle pri">{{ cardData.primaryValue }}</v-card-subtitle>
      </div>
    <div v-if="cardData.secondaryLabel">
      <v-card-text class="card-title">{{ cardData.secondaryLabel }}</v-card-text>
      <v-card-subtitle class="card-subtitle sec">{{ cardData.secondaryValue }}</v-card-subtitle>
    </div>
  </div>
  <div v-if="cardData.percentageChange !== undefined && cardData.percentageChange !== null" class="icon-container">
    <v-card-text  class="card-title percentage" :class="{ 'positive': (cardData.trendColor === '#66BB6A'), 'negative': (cardData.trendColor === '#EF5350') }">
        {{ (cardData.percentageChange) > 0 ? '+' : '' }}{{ (cardData.percentageChange || 0).toFixed(2) }}%
    </v-card-text>
    <BaseIcon 
        v-if="cardData.icon"
        :name="cardData.icon" 
        :iconColor="cardData.iconColor" 
        :class="['card-icon']" 
        style="border-radius: 30%" 
      />
  </div>
  <div v-else class="icon-container placeholder-icon">
    <BaseIcon :name="'users'" :iconColor="cardData.iconColor" :class="['card-icon', cardData.background]" style="border-radius: 30% ; width: 3rem" />

  </div>
</div>
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
        unit: '',
        background: 'orange lighten-4'
      })
    }
  }
};
</script>

<style scoped>

.placeholder-icon{
  margin: auto;
  height: 2rem;

}

.card-container{
  position: relative;
}

.percentage.positive {
  color: #34D399;
}

.percentage.negative {
  color: #EF4444;
}

.outer-container{
  height: 100%;
  padding: 20px;
  display: flex;
  gap: 7rem;
}

.card-icon {
  border-radius: 30%;
}

.title-container{
  display: flex;
  flex-direction: column;
  gap: 1rem;
  justify-content: center;
}

.icon-container{
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 4rem;
}

.percentage{
  text-align: end;
  padding: 0;
  font-weight: 500;
}

.card-title{
  margin: 0;
  padding: 0;
  font-size: 14px;
  font-weight: 400;
  color: #013912;
}

.card-subtitle{
  padding-left: 0;
}

.card-subtitle.pri{
  font-size: 3rem;
  color: #013912;
  font-weight: 700;
}

.card-subtitle.sec{
  color: #013912;
  font-size: 2rem;
  font-weight: 600;
  padding-top: 8px;
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
