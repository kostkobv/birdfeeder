/** K6 stress loading test **/

import http from "k6/http";

const numbers = [31123123123, 32123123123, 33123123123];
const params =  { headers: { "Content-Type": "application/json" } };
const url = "http://localhost:8081/message";

export default function() {
    for (let i = 0; i < numbers.length; i++) {
        const payload = `{
            "recipient": ${numbers[i]},
            "originator": "hey",
            "message": "Message!"
        }`;

        http.post(url, payload, params);
    }
};
