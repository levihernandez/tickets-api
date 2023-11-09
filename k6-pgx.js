import http from 'k6/http';
import { check, sleep } from 'k6';

let targetvus= 400;
export const options = {
  stages: [
    { duration: '30s', target: targetvus },
    { duration: '2m', target: targetvus },
    { duration: '2m', target: targetvus },
    { duration: '30s', target: targetvus },
    { duration: '2m', target: targetvus },
    { duration: '2m', target: targetvus },
    { duration: '30s', target: targetvus },
    { duration: '2m', target: targetvus },
    { duration: '30s', target: 0 },
  ],
};

const BASE_URL = 'http://192.168.1.27:3002';

export default function () {
  const names = [
  // UUIDs must be retrieved at random from CRDB with
    /*
WITH cn AS (
 SELECT CONCAT('''',id,'''') AS userid FROM users ORDER BY RANDOM() LIMIT 100
 )
 SELECT array_agg(userid) AS userid FROM cn;
    */
  ];

  const randomName = names[Math.floor(Math.random() * names.length)];

  const res = http.get(`${BASE_URL}/search/user/${randomName}`);
  check(res, { 'status was 200': (r) => r.status == 200 });

  // Check if the response body can be parsed as JSON
  let users;
  try {
    users = res.json();
  } catch (error) {
    console.error(`Failed to parse response as JSON: ${error}`);
    return;
  }

  // Check if the parsed object has a length property
  if (!Array.isArray(users)) {
    console.error('Response is not an array');
    return;
  }

  // Check if the array is empty
  if (users.length === 0) {
    console.log('No users found');
    return;
  }

  const userID = randomName

  const purchasesRes = http.get(`${BASE_URL}/user/${userID}/purchases`);
  check(purchasesRes, { 'status was 200': (r) => r.status == 200 });

  const cancellationsRes = http.get(`${BASE_URL}/user/${userID}/purchases/cancellations`);
  check(cancellationsRes, { 'status was 200': (r) => r.status == 200 });

  sleep(1);
}
