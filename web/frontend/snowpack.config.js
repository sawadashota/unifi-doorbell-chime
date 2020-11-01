module.exports = {
    mount: {
        public: '/',
        src: '/_dist_',
    },
    plugins: [
        '@snowpack/plugin-optimize',
        '@snowpack/plugin-dotenv',
        '@prefresh/snowpack',
        '@snowpack/plugin-typescript',
    ],
    install: [
        /* ... */
    ],
    installOptions: {
        /* ... */
    },
    devOptions: {
        /* ... */
    },
    buildOptions: {
        out: 'static'
    },
    proxy: {
        /* ... */
    },

    alias: {
        /* ... */
    },
};
