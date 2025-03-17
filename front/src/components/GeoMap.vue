<template>
  <div class="map-container">
    <div class="title">{{ title }}</div>
    <div class="content">
      <div class="chart" ref="chartContainer" style="width: 100%; height: 600px;"></div>
      <div class="color-legend">
        <ul>
          <li v-for="(color, country) in colorMapping" :key="country">
            <span :style="{ backgroundColor: color }" class="color-box"></span>
            {{ country }}
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script>
import * as echarts from 'echarts';
import axios from 'axios';

export default {
  name: 'GeoMap',
  props: {
    title: {
      type: String,
      default: 'Geographic Distribution'
    },
    countrywiseOverviews: {
      type: Array,
      default: () => [],
      deep: true
    }
  },
  data() {
    return {
      chart: null // Store chart instance
    };
  },
  mounted() {
    this.initChart();
    
    // Add resize event listener
    window.addEventListener('resize', this.handleResize);
  },
  beforeUnmount() {
    // Clean up chart and event listeners when component is destroyed
    if (this.chart) {
      this.chart.dispose();
      this.chart = null;
    }
    window.removeEventListener('resize', this.handleResize);
  },
  watch: {
    countrywiseOverviews: {
      handler() {
        this.initChart();
      },
      deep: true
    }
  },
  computed: {
    colorMapping() {
      const mapping = {};
      this.countrywiseOverviews.forEach(item => {
        if (item.GeoMapColorCode) {
          mapping[item.Country] = item.GeoMapColorCode;
        }
      });
      return mapping;
    }
  },
  methods: {
    handleResize() {
      if (this.chart) {
        this.chart.resize();
      }
    },
    async initChart() {
      // Dispose of existing chart instance if it exists
      if (this.chart) {
        this.chart.dispose();
      }
      
      // Initialize new chart
      this.chart = echarts.init(this.$refs.chartContainer);
      this.chart.showLoading();
      
      try {
        const response = await axios.get('/static/world.json');
        echarts.registerMap('world', response.data);
        this.chart.hideLoading();
        
        const mapData = this.countrywiseOverviews.map(item => ({
          name: item.Country,
          itemStyle: {
            areaColor: item.GeoMapColorCode || '#D6D6D6'
          }
        }));
        
        const option = {
          series: [
            {
              name: 'World Map',
              type: 'map',
              map: 'world',
              itemStyle: {
                areaColor: '#D6D6D6',
                borderColor: '#D6D6D6',
                borderWidth: 0.5
              },
              emphasis: {
                disabled: false,
                itemStyle: {
                  areaColor: '#D6D6D6',
                  borderColor: '#FFFFFF',
                  textStyle: {
                    color: '#000000'
                  }
                }
              },
              tooltip: {
                show: true
              },
              data: this.countrywiseOverviews.length > 0 ? mapData : []
            }
          ]
        };

        this.chart.setOption(option);
      } catch (error) {
        console.error('Error initializing chart:', error);
        this.chart.hideLoading();
      }
    }
  }
};
</script>

<style scoped>
.chart {
  width: 80%;
  height: 600px;
  padding: 20px;
}

.map-container {
  border: 1px solid #d6d6d6;
}

.title {
  padding: 10px;
  border-bottom: 1px solid #d6d6d6;
  color: #111618;
  font-size: 10px;
}

.content {
  display: flex;
  align-items: center;
}

.color-legend {
  margin-left: 20px;
  width: 200px;
}

.color-legend ul {
  list-style: none;
  padding: 0;
}

.color-legend li {
  display: flex;
  align-items: center;
  margin-bottom: 5px;
  font-size: 12px;
}

.color-box {
  width: 15px;
  height: 15px;
  display: inline-block;
  margin-right: 8px;
  border-radius: 50%;
}
</style>
