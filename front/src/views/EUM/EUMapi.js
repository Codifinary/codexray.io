import data from './data.json';
import ChartData from './ChartData.json';

export function fetchCardData() {
    return data;
}
export function fetchChartData() {
    console.log(ChartData);
    return ChartData;
}
