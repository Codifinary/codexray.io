// parentfolder/api/EUMapi.js
import data from '../data/data.json';

export const getApplications = () => {
  return data.applications;
};

export const getApplicationData = (applicationName) => {
  const { pagePerformance, errorTab, errorLogs } = data;
  const appData = {
    pagePerformance: pagePerformance?.applications?.find(
      (app) => app.applicationName === applicationName
    ) || null,
    errors: errorTab?.applications?.find(
      (app) => app.applicationName === applicationName
    ) || null,
    logs: errorLogs?.applications?.find(
      (app) => app.applicationName === applicationName
    ) || null,
  };
  return appData;
};
