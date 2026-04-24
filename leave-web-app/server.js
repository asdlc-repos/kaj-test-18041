const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');
const path = require('path');

const app = express();
const PORT = process.env.PORT || 8080;

const LEAVE_API_URL = (process.env.LEAVE_API_URL || 'http://localhost:8081').replace(/\/$/, '');
const MANAGER_API_URL = (process.env.MANAGER_API_URL || 'http://localhost:8082').replace(/\/$/, '');

console.log(`Starting leave-web-app server on port ${PORT}`);
console.log(`LEAVE_API_URL: ${LEAVE_API_URL}`);
console.log(`MANAGER_API_URL: ${MANAGER_API_URL}`);

// Health endpoint
app.get('/health', (_req, res) => {
  res.status(200).json({ status: 'ok', service: 'leave-web-app' });
});

// Proxy /api/leave/* → leave-api
app.use('/api/leave', createProxyMiddleware({
  target: LEAVE_API_URL,
  changeOrigin: true,
  pathRewrite: { '^/api/leave': '' },
  on: {
    error: (err, _req, res) => {
      console.error('Leave API proxy error:', err.message);
      res.status(502).json({ error: 'Leave API unavailable' });
    },
  },
}));

// Proxy /api/manager/* → manager-api
app.use('/api/manager', createProxyMiddleware({
  target: MANAGER_API_URL,
  changeOrigin: true,
  pathRewrite: { '^/api/manager': '' },
  on: {
    error: (err, _req, res) => {
      console.error('Manager API proxy error:', err.message);
      res.status(502).json({ error: 'Manager API unavailable' });
    },
  },
}));

// Serve static React build
app.use(express.static(path.join(__dirname, 'dist')));

// SPA fallback — serve index.html for all unmatched routes
app.get('*', (_req, res) => {
  res.sendFile(path.join(__dirname, 'dist', 'index.html'));
});

app.listen(PORT, () => {
  console.log(`Server is running at http://localhost:${PORT}`);
});
