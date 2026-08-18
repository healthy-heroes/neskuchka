import { Group } from '@mantine/core';
import { Logo } from '../Logo/Logo';
import { RouteLink } from '../RouteLink/RouteLink';
import { UserMenu } from './UserMenu';
import classes from './Header.module.css';

/**
 * Header — шапка приложения.
 *
 * Без обёртки вокруг <header> намеренно: sticky отсчитывает ход внутри
 * родителя, и обёртка ровно по высоте шапки его обнуляет.
 */
export function Header() {
	return (
		<header className={classes.header}>
			<Group justify="space-between" h="100%">
				<Group gap="xl">
					<Logo />
					<RouteLink
						to="/workouts"
						className={classes.nav}
						activeProps={{ className: classes.navActive }}
						underline="never"
					>
						Тренировки
					</RouteLink>
				</Group>
				<UserMenu />
			</Group>
		</header>
	);
}
