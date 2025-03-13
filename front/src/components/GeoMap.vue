<template>
  <div class="map-container">
    <div class="title">{{ title }}</div>
    <div class="chart" ref="chartContainer" style="width: 100%; height: 600px;"></div>
  </div>
</template>

<script>
import * as echarts from 'echarts';
import axios from 'axios';

export default {
  props: {
    title: String,
    countrywiseOverviews: {
      type: Array,
      default: () => []
    }
  },
  name: 'WorldMap',
  mounted() {
    this.initChart();
  },
  watch: {
    countrywiseOverviews: {
      handler() {
        this.initChart();
      },
      deep: true
    }
  },
  methods: {
    async initChart() {
      const chart = echarts.init(this.$refs.chartContainer);
      chart.showLoading();
      try {
        const worldJson = await axios.get('/static/world.json');
        echarts.registerMap('world', worldJson.data);
        chart.hideLoading();

        // Transform the countrywiseOverviews for the chart - only use country name and color code
        const mapData = this.countrywiseOverviews.map(item => ({
          name: item.Country,
          itemStyle: {
            areaColor: item.GeoMapColorCode || '#D6D6D6'
          }
        }));

        // Default data if no countrywiseOverviews is provided
        const defaultData = [
          { name: 'China', itemStyle: { areaColor: '#d73027' } },
          { name: 'India', itemStyle: { areaColor: '#fdae61' } },
          { name: 'United States', itemStyle: { areaColor: '#74add1' } },
          { name: 'Indonesia', itemStyle: { areaColor: '#ffffbf' } },
          { name: 'Brazil', itemStyle: { areaColor: '#ffffbf' } },
          { name: 'Pakistan', itemStyle: { areaColor: '#74add1' } },
          { name: 'Nigeria', itemStyle: { areaColor: '#74add1' } },
          { name: 'Bangladesh', itemStyle: { areaColor: '#313695' } },
          { name: 'Russia', itemStyle: { areaColor: '#313695' } },
          { name: 'Mexico', itemStyle: { areaColor: '#313695' } }
        ];

        const option = {
          tooltip: {
            trigger: 'item',
            showDelay: 0,
            transitionDuration: 0.2,
            formatter: (params) => {
              return `${params.name}`;
            }
          },
          series: [
            {
              name: 'World Map',
              type: 'map',
              map: 'world',
              itemStyle: {
                areaColor: '#D6D6D6', // Set default land color
                borderColor: '#D6D6D6', // Set border color
                borderWidth: 0.5
              },
              emphasis: {
                disabled: false, // Enable hover effect
                itemStyle: {
                  areaColor: '#a0a0a0'
                }
              },
              data: this.countrywiseOverviews.length > 0 ? mapData : defaultData
            }
          ]
        };

        chart.setOption(option);

        // Handle window resize
        window.addEventListener('resize', () => {
          chart.resize();
        });

        // Clean up event listener when component is destroyed
        this.$once('hook:beforeDestroy', () => {
          window.removeEventListener('resize', chart.resize);
          chart.dispose();
        });

      } catch (error) {
        console.error('Error loading world map data:', error);
      }
    }
  }
};
</script>

<style scoped>
.chart{
  width: 100%;
  height: 600px;
  padding: 20px;
}

.map-container{
  border: 1px solid #d6d6d6;
  /* padding: 10px; */
}

.title{
  padding: 10px;
  border-bottom: 1px solid #d6d6d6;
}
</style>