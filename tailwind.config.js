/** @type {import('tailwindcss').Config} */
module.exports = {
    darkMode: '[data-mode="dark"]',
    content: ["./resources/views/**/*.templ", "./resources/views/**/*.go"],
    theme: {
        extend: {},
    },
    plugins: [
        require("@tailwindcss/forms"),
        require("daisyui")
    ],
    daisyui: {
        themes: ["dark"]
    }
}