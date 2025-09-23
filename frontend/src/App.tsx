import '@mantine/core/styles.css';
import '@mantine/tiptap/styles.css';
import '@mantine/dates/styles.css';
import './App.css';

import dayjs from 'dayjs';
import { StrictMode } from 'react';
import { MantineProvider } from '@mantine/core';
import { DatesProvider } from '@mantine/dates';
import { ApiProvider } from './api/provider';
import { Router } from './Router';
import { theme } from './theme';

import 'dayjs/locale/ru';

// todo: get locale from DatesProvider
dayjs.locale('ru');

export default function App() {
	return (
		<StrictMode>
			<DatesProvider settings={{ locale: 'ru', firstDayOfWeek: 0, weekendDays: [0] }}>
				<ApiProvider>
					<MantineProvider theme={theme}>
						<Router />
					</MantineProvider>
				</ApiProvider>
			</DatesProvider>
		</StrictMode>
	);
}
