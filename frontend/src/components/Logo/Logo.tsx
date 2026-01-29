import { Text } from '@mantine/core';
import { RouteLink } from '../RouteLink/RouteLink';
import { LogoIcon } from './LogoIcon';
import classes from './Logo.module.css';

export function Logo() {
	return (
		<RouteLink to="/" underline="never" c="black">
			<Text size="xl" fw={700}>
				<LogoIcon className={classes.logoIcon} size="xl" />
				&nbsp;Neskuchka
			</Text>
		</RouteLink>
	);
}
