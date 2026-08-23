import { readdirSync, readFileSync, statSync, writeFileSync } from 'node:fs';
import { dirname, join, relative, resolve } from 'node:path';

/**
 * Переписывает алиасы `@/…` в выпущенных декларациях на относительные пути.
 *
 * tsc сохраняет спецификатор импорта как написано, а конвертер design-sync
 * читает `.d.ts` через ts-morph с захардкоженными compilerOptions — без `paths`.
 * С алиасами он не находит ни одного экспорта и выбрасывает все компоненты.
 *
 * Дерево деклараций повторяет `src/`, поэтому `@/x` — это просто `<root>/x`.
 * Побочные импорты стилей выкидываем: типов в них нет, а файлов рядом с
 * декларациями не существует.
 */
const ROOT = resolve('dist-ds/types');

function* walk(dir) {
	for (const name of readdirSync(dir)) {
		const path = join(dir, name);
		if (statSync(path).isDirectory()) {
			yield* walk(path);
		} else if (path.endsWith('.d.ts')) {
			yield path;
		}
	}
}

let rewritten = 0;

for (const file of walk(ROOT)) {
	const text = readFileSync(file, 'utf8');

	const next = text
		.replace(/^import ['"][^'"]+\.css['"];?\n/gm, '')
		.replace(/(['"])@\/([^'"]+)\1/g, (_match, quote, subpath) => {
			const rel = relative(dirname(file), join(ROOT, subpath)).split('\\').join('/');

			return `${quote}${rel.startsWith('.') ? rel : `./${rel}`}${quote}`;
		});

	if (next !== text) {
		writeFileSync(file, next);
		rewritten += 1;
	}
}

console.error(`ds-dts-paths: rewrote ${rewritten} declaration file(s)`);
