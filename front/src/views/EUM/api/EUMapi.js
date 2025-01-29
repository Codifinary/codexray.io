import data from '../data/data.json';

// export const getApplications = () => {
//     return data.applications;
// };

// export const getApplicationData = () => {
//     const { pagePerformance, errorTab } = data;
//     const appData = {
//         pagePerformance: pagePerformance || null,
//         errors: errorTab || null,
//     };
//     return appData;
// };

export const getErrorDetails = () => {
    return data.errorDetails || null;
};

export const getSpecificErrors = () => {
    const appData = data.specificErrors || null;
    if (appData) {
        return appData.eventLogs;
    }
    return [];
};
export const getBreadcrumbsByType = (type) => {
    const tableData = data.errorDetails.breadcrumb.filter((breadcrumb) => breadcrumb.type === type) || [];
    console.log(tableData);
    return tableData;
};

import logsData from '../data/logsData.json';

export const getEventLogs = () => {
    const { entries, chart } = logsData.data;
    return { entries, chart };
};

export const getFilteredEventLogs = (severity = [], search = '') => {
    console.log('Filtering logs data:', { severity, search });
    const { entries } = logsData.data;
    const searchLower = search.toLowerCase();
    const filteredEntries = entries.filter((entry) => {
        const filterSeverity = severity.length === 0 || severity.includes(entry.severity);
        const filterSearch = search === '' || entry.message.toLowerCase().includes(searchLower);
        return filterSeverity && filterSearch;
    });

    return filteredEntries;
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
