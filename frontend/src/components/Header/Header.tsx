import { Box, Button, Group } from '@mantine/core';
import { useAuth } from '@/auth/hooks';
import { Logo } from '../Logo/Logo';
import { RouteLink } from '../RouteLink/RouteLink';
import classes from './Header.module.css';

export function Header() {
	const { user, isAuthenticated, isLoading } = useAuth();

	function userInfo() {
		if (isLoading) {
			return;
		}

		console.log(isAuthenticated, user);

		if (isAuthenticated && user) {
			return <div>{user.Name}</div>;
		}

		return (
			<Button component={RouteLink} to="/login">
				Войти
			</Button>
		);
	}

	return (
		<Box>
			<header className={classes.header}>
				<Group justify="space-between" h="100%">
					<Logo />

					<Group>{userInfo()}</Group>
				</Group>
			</header>
		</Box>
	);
}
