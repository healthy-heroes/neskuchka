import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

/**
 * Сборка дизайн-системы для claude.ai/design.
 *
 * Отдельный конфиг, а не режим основного: приложению нужен html-вход, роутер-плагин
 * и прокси, дизайн-системе — один ESM-модуль с именованными экспортами. Конвертер
 * design-sync грузит его как `window.NeskuchkaDS` и резолвит в него импорты стори.
 *
 * React внешний: в бандле должна быть ровно одна его копия, иначе хуки и контексты
 * в превью и в компонентах окажутся из разных реестров.
 *
 * Шрифты инлайнятся в CSS: страница на claude.ai/design не отдаёт наши woff2,
 * а без Oswald дизайн-система приезжает с чужой типографикой.
 */
export default defineConfig({
	plugins: [react()],

	resolve: {
		tsconfigPaths: true,
	},

	define: {
		__API_URL__: JSON.stringify(process.env.VITE_BACKEND_PORT || 8080),
	},

	build: {
		outDir: 'dist-ds',
		emptyOutDir: true,
		cssCodeSplit: false,
		assetsInlineLimit: 512 * 1024,
		lib: {
			entry: 'src/design/index.ts',
			formats: ['es'],
			fileName: () => 'ds.js',
			cssFileName: 'ds',
		},
		rollupOptions: {
			external: ['react', 'react-dom', 'react/jsx-runtime', 'react-dom/client'],
			// Один файл вместо графа чанков: конвертер грузит бандл как единый модуль
			output: { inlineDynamicImports: true },
		},
	},
});
