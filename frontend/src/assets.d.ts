/** Изображения импортируются как URL: vite отдаёт путь, esbuild превью — data-URL. */
declare module '*.jpg' {
	const src: string;
	export default src;
}
