const http = require('http');

const server = http.createServer((req, res) => {
  console.log("got request!!")
  res.writeHead(200, { 'Content-Type': 'text/plain' });
  res.end('Hello World\n');
});

const port = process.env.PORT || 3000;
const host = '0.0.0.0'; // Listen on all interfaces

server.listen(port, host, () => {
  console.log(`Server running at http://${host}:${port}/`);
});
