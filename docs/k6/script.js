import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  // A number specifying the number of VUs to run concurrently.
  vus: 1,
  // A string specifying the total duration of the test run.
  duration: '5s',

};

export default function() {
  const paramsWithIP = {
    headers: {
      'X-Real-IP': '123.321.123',
    },
  };

  for (let id = 1; id <= 10; id++) {
    http.get('http://ratelimiter-go:8080', paramsWithIP);
    sleep(1);
  }  

  const paramsWithToken = {
    headers: {
      'Api-Key': 'token_abc',
    },
  };

  for (let id = 1; id <= 10; id++) {
    http.get('http://ratelimiter-go:8080', paramsWithToken);
    sleep(1);
  }  
}
