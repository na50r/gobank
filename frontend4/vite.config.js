export default {
    // Vite Dev Server will not work with SSE events unless we proxy the events
    // No issue when built and previewed
    server: {
        proxy: {
            '/events': {
                target: 'http://localhost:3000',
                changeOrigin: true,
            }
        }
    }
}