import UPlot from 'uplot';

export const SECOND = 1000;
export const MINUTE = SECOND * 60;
export const HOUR = MINUTE * 60;
export const DAY = HOUR * 24;

export function duration(ms, precision) {
    let milliseconds = ms;
    const days = Math.floor(milliseconds / DAY);
    milliseconds %= DAY;
    const hours = Math.floor(milliseconds / HOUR);
    milliseconds %= HOUR;
    const minutes = Math.floor(milliseconds / MINUTE);
    milliseconds %= MINUTE;
    const seconds = Math.floor(milliseconds / SECOND);
    milliseconds %= SECOND;

    const names = {
        d: days,
        h: hours,
        m: minutes,
        s: seconds,
        ms: milliseconds,
    };

    let res = '';
    let stop = false;
    for (const n in names) {
        if (n === precision) {
            stop = true;
        }
        const v = names[n];
        if (v) {
            res += v + n;
            if (stop) {
                break;
            }
        }
    }
    return res.trimEnd();
}

export function durationPretty(ms) {
    if (ms > 5 * DAY) {
        return duration(ms, 'd');
    }
    if (ms > DAY) {
        return duration(ms, 'h');
    }
    if (ms > HOUR) {
        return duration(ms, 'm');
    }
    if (ms > MINUTE) {
        return duration(ms, 's');
    }
    return duration(ms, 'ms');
}

export function date(ms, format) {
    return UPlot.fmtDate(format)(new Date(ms));
}

export function timeSinceNow(ms) {
    return durationPretty(Date.now() - ms);
}

export function percent(p) {
    if (p > 10) {
        return p.toFixed(0);
    }
    if (p > 1) {
        return p.toFixed(1);
    }
    return p.toFixed(2);
}

export function float(f) {
    if (f === 0) {
        return '0';
    }
    if (f >= 1) {
        return f.toFixed(0);
    }
    if (f >= 0.1) {
        return f.toFixed(1);
    }
    if (f >= 0.01) {
        return f.toFixed(2);
    }
    return f.toFixed(3);
}

export function formatUnits(v, unit) {
    if (unit === 'ts') {
        return this.$format.date(v, '{MMM} {DD}, {HH}:{mm}');
    }
    if (unit === 'dur') {
        if (!v) {
            return '0';
        }
        if (v === 'inf' || v === 'err') {
            return 'Inf';
        }
        if (v >= 1) {
            return v + 's';
        }
        return v * 1000 + 'ms';
    }
    if (unit === '%') {
        v *= 100;
        if (v < 1) {
            return '<1';
        }
        let d = 1;
        if (v >= 10) {
            d = 0;
        }
        return v.toFixed(d);
    }
    if (unit === 'ms') {
        let d = 0;
        if (v < 10) {
            d = 1;
        }
        return v.toFixed(d);
    }
    let m = '';
    if (v > 1e3) {
        v /= 1000;
        m = 'K';
    }
    if (v > 1e6) {
        v /= 1000;
        m = 'M';
    }
    if (v > 1e9) {
        v /= 1000;
        m = 'G';
    }
    return v.toFixed(1) + m;
}

export function convertLatency(latency) {
    if (latency < 1000) {
        return { value: parseFloat(latency), unit: 'ms' };
    } else if (latency < 60000) {
        return { value: parseFloat(latency / 1000), unit: 's' };
    } else if (latency < 3600000) {
        return { value: parseFloat(latency / 60000), unit: 'min' };
    } else {
        return { value: parseFloat(latency / 3600000), unit: 'hr' };
    }
}

export function shortenNumber(value) {
    if (value >= 1e9) {
        return { value: parseFloat(value / 1e9), unit: 'G' };
    } else if (value >= 1e6) {
        return { value: parseFloat(value / 1e6), unit: 'M' };
    } else if (value >= 1e3) {
        return { value: parseFloat(value / 1e3), unit: 'K' };
    } else {
        return { value: parseFloat(value), unit: '' };
    }
}

export function copyToClipboard(text) {
    if (!text) {
      console.warn('copyToClipboard: No text provided to copy.');
      return;
    }
  
    if (navigator.clipboard && navigator.clipboard.writeText) {
      // Use the modern Clipboard API
      navigator.clipboard.writeText(text).catch((err) => {
        console.error('copyToClipboard: Failed with Clipboard API, falling back.', err);
        fallbackCopy(text);
      });
      return true;
    } else {
      // Fallback immediately if Clipboard API is not available
      return fallbackCopy(text);
    }
  
    function fallbackCopy(text) {
      const textarea = document.createElement('textarea');
      textarea.value = text;
  
      // Make textarea invisible and off-screen
      textarea.style.position = 'fixed';
      textarea.style.opacity = '0';
      document.body.appendChild(textarea);
  
      textarea.focus();
      textarea.select();
      textarea.setSelectionRange(0, 99999); // For mobile devices
  
      try {
        const successful = document.execCommand('copy');
        if (!successful) {
          console.error('copyToClipboard: Fallback copy failed.');
          return false;
        }
      } catch (err) {
        console.error('copyToClipboard: Error using fallback method.', err);
        return false;
      }
  
      document.body.removeChild(textarea);
      return true;
    }
  }
  