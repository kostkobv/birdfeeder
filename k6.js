/** K6 stress loading test **/

import http from "k6/http";

export default function() {
    const url = "http://localhost:8081/message";
    const payload = `{
        "recipient": 123456789,
        "originator": "hey",
        "message": "Message!"
    }`;
    const params =  { headers: { "Content-Type": "application/json" } };
    http.post(url, payload, params);
};
