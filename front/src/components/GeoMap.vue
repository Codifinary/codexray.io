<template>
  <div class="map-container">
    <div class="title">{{ title }}</div>
    <div class="content">
      <div class="chart" ref="chartContainer"></div>
      <div class="color-legend">
        <div class="legend-container">
          <ul>
            <li v-for="(color, country) in colorMapping" :key="country">
              <span :style="{ backgroundColor: color }" class="color-box"></span>
              <span class="country-name">{{ country }}</span>
            </li>
          </ul>
        </div>
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
      chart: null
    };
  },
  mounted() {
    window.addEventListener('resize', this.handleResize);
    this.$nextTick(() => {
      setTimeout(() => {
        this.initChart();
      }, 100);
    });
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.handleResize);
    if (this.chart) {
      this.chart.dispose();
      this.chart = null;
    }
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
      if (this.chart) {
        this.chart.dispose();
      }
      
      await this.$nextTick();
      
      const container = this.$refs.chartContainer;
      
      if (!container || container.clientWidth === 0 || container.clientHeight === 0) {
        setTimeout(() => this.initChart(), 100);
        return;
      }
      
      const width = container.clientWidth;
      const height = container.clientHeight;
      
      this.chart = echarts.init(container, null, {
        renderer: 'canvas',
        useDirtyRect: false,
        width: width,
        height: height
      });
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
          backgroundColor: '#fff',
          series: [
            {
              name: 'World Map',
              type: 'map',
              map: 'world',
              roam: false,
              zoom: 1,
              center: [0, 0],
              aspectScale: 1,
              top: 0,
              left: 0,
              right: 0,
              bottom: 0,
              boundingCoords: [[-180, 90], [180, -90]],
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
        this.chart.resize();
      } catch (error) {
        console.error('Error initializing chart:', error);
        this.chart.hideLoading();
      }
    }
  }
};
</script>

<style scoped>
.map-container {
  border: 1px solid #d6d6d6;
  width: 100%;
  height: 800px;
  position: relative;
  overflow: hidden;
  box-sizing: border-box;
}

.title {
  padding: 10px;
  border-bottom: 1px solid #d6d6d6;
  color: #111618;
  font-size: 10px;
}

.content {
  display: flex;
  padding-top: 50px !important;
  padding-left: 5% !important;
  width: 95% !important;
  position: relative;
  height: calc(100% - 35px);
  /* height: 100% !important; */
  padding: 5px;
}

.chart {
  height: 100%;
  width: 100%;
  padding: 0px;
  min-height: 700px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.color-legend {
  position: absolute;
  top: 300px;
  left: 50px;
  width: 180px;
  height: auto;
  background-color: #ffffff;
  border-radius: 4px;
  padding: 10px;
}

.legend-container {
  overflow-y: visible;
}

.color-legend ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.color-legend li {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.country-name {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 130px;
}

.color-box {
  min-width: 12px;
  height: 12px;
  display: inline-block;
  margin-right: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}
</style>
