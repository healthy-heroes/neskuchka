import '@mantine/core/styles.css';
import '@mantine/tiptap/styles.css';
import '@mantine/dates/styles.css';
import './App.css';

import dayjs from 'dayjs';
import { StrictMode } from 'react';
import { MantineProvider } from '@mantine/core';
import { ApiProvider } from './api/provider';
import ApiService from './api/service';
import { API_URL } from './config';
import { Router } from './Router';
import { theme } from './theme';

import 'dayjs/locale/ru';

import { DatesProvider } from '@mantine/dates';

// todo: get locale from DatesProvider
dayjs.locale('ru');

const apiConfig = {
	apiUrl: API_URL,
};

export default function App() {
	return (
		<StrictMode>
			<DatesProvider settings={{ locale: 'ru', firstDayOfWeek: 0, weekendDays: [0] }}>
				<ApiProvider apiService={new ApiService(apiConfig)}>
					<MantineProvider theme={theme}>
						<Router />
					</MantineProvider>
				</ApiProvider>
			</DatesProvider>
		</StrictMode>
	);
}
