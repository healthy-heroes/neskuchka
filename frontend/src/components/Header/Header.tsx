import { Box, Group } from '@mantine/core';
import { Logo } from '../Logo/Logo';
import { UserMenu } from './UserMenu';
import classes from './Header.module.css';

export function Header() {
	return (
		<Box>
			<header className={classes.header}>
				<Group justify="space-between" h="100%">
					<Logo />
					<UserMenu />
				</Group>
			</header>
		</Box>
	);
}
