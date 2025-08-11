import '@mantine/core/styles.css';

import dayjs from 'dayjs';

import 'dayjs/locale/ru';

import { MantineProvider } from '@mantine/core';
import { Router } from './Router';
import { theme } from './theme';

import './App.css';

dayjs.locale('ru');

export default function App() {
	return (
		<MantineProvider theme={theme}>
			<Router />
		</MantineProvider>
	);
}
