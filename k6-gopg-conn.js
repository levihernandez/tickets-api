import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  scenarios: {
    // Define a scenario that ramps up and down repeatedly
    interrupted_load: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '10s', target: 20 }, // Ramp-up VUs for 10 seconds
        { duration: '20s', target: 20 }, // Stay at target VUs for 20 seconds
        { duration: '10s', target: 0 },  // Ramp-down to 0 VUs for 10 seconds
      ],
      gracefulRampDown: '0s',
      exec: 'userScenario', // The function to execute for this scenario
    },
  },
  thresholds: {
    // Define any thresholds
  },
};

// Main user scenario logic
export function userScenario() {
  const BASE_URL = 'http://localhost:3001';

  // ... (rest of your user scenario code remains the same)

  sleep(1); // Control iteration pacing here if necessary
}

export function teardown(data) {
  // Optional: Any teardown logic if necessary when the test ends
}
