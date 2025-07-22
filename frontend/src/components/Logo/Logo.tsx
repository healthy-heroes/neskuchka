import { Text } from '@mantine/core';
import { LogoIcon } from './LogoIcon';
import classes from './Logo.module.css';

export function Logo() {
	return (
		<Text size="xl" fw={700}>
			<LogoIcon className={classes.logoIcon} size="xl" />
			&nbsp;Neskuchka
		</Text>
	);
}
