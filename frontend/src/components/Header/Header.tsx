import { Box, Group } from '@mantine/core';
import { Logo } from '../Logo/Logo';
import { RouteLink } from '../RouteLink/RouteLink';
import { UserMenu } from './UserMenu';
import classes from './Header.module.css';

export function Header() {
	return (
		<Box>
			<header className={classes.header}>
				<Group justify="space-between" h="100%">
					<Group gap="xl">
						<Logo />
						<RouteLink to="/workouts" c="black" underline="hover">
							Тренировки
						</RouteLink>
					</Group>
					<UserMenu />
				</Group>
			</header>
		</Box>
	);
}
