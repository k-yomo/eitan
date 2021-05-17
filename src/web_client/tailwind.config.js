const colors = require('tailwindcss/colors')

module.exports = {
    purge: ['./src/pages/**/*.{ts,tsx}', './src/components/**/*.{ts,tsx}'],
    darkMode: false, // or 'media' or 'class'
    theme: {
        extend: {
            gradientColorStops: theme => ({
                'blue-gray': colors.blueGray,
                'fuchsia': colors.fuchsia,
                'light-blue': colors.lightBlue,
                'orange': colors.orange,
                'rose': colors.rose,
            })
        },
    },
    variants: {
        extend: {},
    },
    plugins: [],
}
