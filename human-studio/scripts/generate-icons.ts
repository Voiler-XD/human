/**
 * Generate platform-specific icons from the Human Studio logo SVG.
 *
 * This script creates:
 * - resources/icon.png (512x512 for Linux)
 * - resources/icon-16.png
 * - resources/icon-32.png
 * - resources/icon-256.png
 *
 * For .ico and .icns, use a tool like `electron-icon-builder` or
 * manually convert using ImageMagick:
 *   convert icon.png -define icon:auto-resize=256,48,32,16 icon.ico
 *   iconutil -c icns icon.iconset  (macOS)
 *
 * Prerequisites: npm install sharp
 */

// Placeholder: In production, use `sharp` or `@electron/rebuild` to convert
// the SVG to raster formats at required sizes.
//
// For now, the SVG is in resources/ and electron-builder can use PNG directly.

console.log('Icon generation script — use sharp or electron-icon-builder')
console.log('Required outputs:')
console.log('  resources/icon.ico    — Windows (16,32,48,256)')
console.log('  resources/icon.icns   — macOS (16-1024 @1x and @2x)')
console.log('  resources/icon.png    — Linux (512x512)')
console.log('')
console.log('For development, place the logo SVG as resources/icon.svg')
console.log('and use `npx electron-icon-builder --input=resources/icon.svg --output=resources/`')
