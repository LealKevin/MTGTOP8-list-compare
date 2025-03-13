const { default: daisyui } = require("daisyui");

/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./internal/**/*.gohtml",
        "./internal/**/*.go",
        "./templates/**/*.html",
        "./static/*.html",
        "./static/*.css",
    ],
    theme: {
        extend: {},
    },
    plugins: [require("daisyui")],
    daisyui: {
        themes: ["dark"],
    },
};
