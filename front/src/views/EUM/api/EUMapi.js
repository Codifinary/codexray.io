import data from '../data/data.json';

export const getApplications = () => {
    return data.applications;
};

export const getApplicationData = (applicationName) => {
    const { pagePerformance, errorTab } = data;
    const appData = {
        pagePerformance: pagePerformance?.applications?.find((app) => app.applicationName === applicationName) || null,
        errors: errorTab?.applications?.find((app) => app.applicationName === applicationName) || null,
    };
    return appData;
};

export const getErrorDetails = () => {
    return data.errorDetails || null;
};

export const getSpecificErrors = (applicationName, error) => {
    const appData = data.specificErrors?.applications?.find((app) => app.applicationName === applicationName) || null;
    if (appData) {
        const errorData = appData.errors.find((err) => err.error === error);
        return errorData ? errorData.eventLogs : [];
    }
    return [];
};
export const getBreadcrumbsByType=(category)=> {
    const tableData= data.errorDetails.breadcrumb.filter(breadcrumb => breadcrumb.category === category) || [];
    console.log(tableData);
    return tableData;
    }

import logsData from '../data/logsData.json';

export const getEventLogs = () => {
    console.log(logsData);
    const { entries, chart } = logsData.data;
    return { entries, chart };
};

import tracesData from '../data/tracesData.json';

export const getHeatmapData = () => {
  try {
    return tracesData.data.traces.heatmap || null;
  } catch (error) {
    console.error('Error fetching heatmap data:', error);
    return null;
  }
};

export const getApplicationTraces = () => {
  try {
    return tracesData.data.traces.traces || [];
  } catch (error) {
    console.error('Error fetching application traces:', error);
    return [];
  }
};
