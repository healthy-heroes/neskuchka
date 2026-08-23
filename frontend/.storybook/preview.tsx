import '@fontsource-variable/commissioner';
import '@fontsource-variable/oswald';
import '@mantine/core/styles.css';
import '@mantine/tiptap/styles.css';
import '@mantine/dates/styles.css';
import '../src/App.css';

import dayjs from 'dayjs';
import React, { useEffect } from 'react';
import { DARK_MODE_EVENT_NAME } from 'storybook-dark-mode';
import { addons } from 'storybook/preview-api';
import { MantineProvider, useMantineColorScheme } from '@mantine/core';
import { DatesProvider } from '@mantine/dates';
import { theme } from '../src/theme';

import 'dayjs/locale/ru';

// Как в App.tsx: без этого даты в стори приезжают по-английски
dayjs.locale('ru');

const channel = addons.getChannel();

export const parameters = {
	layout: 'fullscreen',
	options: {
		showPanel: false,
		storySort: (a, b) => {
			return a.title.localeCompare(b.title, undefined, { numeric: true });
		},
	},
};

function ColorSchemeWrapper({ children }: { children: React.ReactNode }) {
	const { setColorScheme } = useMantineColorScheme();
	const handleColorScheme = (value: boolean) => setColorScheme(value ? 'dark' : 'light');

	useEffect(() => {
		channel.on(DARK_MODE_EVENT_NAME, handleColorScheme);
		return () => channel.off(DARK_MODE_EVENT_NAME, handleColorScheme);
	}, [channel]);

	return children;
}

export const decorators = [
	(renderStory: any) => <ColorSchemeWrapper>{renderStory()}</ColorSchemeWrapper>,
	(renderStory: any) => <MantineProvider theme={theme}>{renderStory()}</MantineProvider>,
	(renderStory: any) => (
		<DatesProvider settings={{ locale: 'ru', firstDayOfWeek: 0, weekendDays: [0] }}>
			{renderStory()}
		</DatesProvider>
	),
];
