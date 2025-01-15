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
