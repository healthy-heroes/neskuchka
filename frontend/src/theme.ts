import {
	Badge,
	Button,
	Card,
	createTheme,
	rem,
	Title,
	type MantineColorsTuple,
} from '@mantine/core';

/**
 * Neskuchka theme.
 * Палитра — Copper Aquamarine Dream: copper (действия, «сегодня»),
 * slate (тёмные плашки, статус «выполнено»), тёплые нейтральные вместо серых Mantine.
 * Шрифты подключаются в App.tsx через @fontsource-variable.
 */

const copper: MantineColorsTuple = [
	'#fdf3ed',
	'#f8e0d2',
	'#f0c3a9',
	'#e6a37e',
	'#dc8b5f',
	'#d7794d',
	'#c35727',
	'#a5471e',
	'#843818',
	'#632911',
];

const slate: MantineColorsTuple = [
	'#eef4f6',
	'#dbe7ea',
	'#b7ced4',
	'#8fb2bb',
	'#6d9aa5',
	'#4b848d',
	'#3d6d76',
	'#30525c',
	'#26424a',
	'#1b3037',
];

/** Тёплые нейтральные — переопределяют мантиновский gray */
const warm: MantineColorsTuple = [
	'#faf9f8',
	'#f4f2f1',
	'#eceae8',
	'#e0dddb',
	'#cfcbc8',
	'#bfb9b6',
	'#948d89',
	'#6b6561',
	'#4a4542',
	'#2b2725',
];

export const theme = createTheme({
	colors: { copper, slate, gray: warm },
	primaryColor: 'copper',
	primaryShade: { light: 6, dark: 4 },
	black: '#2b2725',

	fontFamily:
		'Commissioner Variable, -apple-system, BlinkMacSystemFont, Segoe UI, Roboto, Helvetica, Arial, sans-serif',
	fontFamilyMonospace: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace',

	headings: {
		fontFamily: 'Oswald Variable, Commissioner Variable, sans-serif',
		fontWeight: '600',
		sizes: {
			h1: { fontSize: rem(52), lineHeight: '1' },
			h2: { fontSize: rem(40), lineHeight: '1.05' },
			h3: { fontSize: rem(22), lineHeight: '1.2' },
			h4: { fontSize: rem(18), lineHeight: '1.3' },
			h5: { fontSize: rem(16), lineHeight: '1.4' },
			h6: { fontSize: rem(14), lineHeight: '1.4' },
		},
	},

	defaultRadius: 'md',
	radius: { xs: rem(4), sm: rem(6), md: rem(8), lg: rem(14), xl: rem(24) },

	components: {
		Title: Title.extend({
			styles: { root: { textTransform: 'uppercase', letterSpacing: '0.03em' } },
		}),
		Button: Button.extend({
			defaultProps: { radius: 'md' },
			styles: { root: { fontWeight: 600 } },
		}),
		Card: Card.extend({
			// 12px — радиус карточек на обоих экранах; в шкале такого шага нет намеренно
			defaultProps: { radius: 12, withBorder: true },
		}),
		Badge: Badge.extend({
			defaultProps: { radius: 'xl' },
			styles: { label: { fontWeight: 700, letterSpacing: '0.06em' } },
		}),
	},
});
