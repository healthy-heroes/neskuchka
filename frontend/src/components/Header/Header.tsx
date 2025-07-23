import { Box, Button, Group } from '@mantine/core';
import { Logo } from '../Logo/Logo';
import classes from './Header.module.css';

export function Header() {
	return (
		<Box>
			<header className={classes.header}>
				<Group justify="space-between" h="100%">
					<Logo />

					<Group>
						<Button>Войти</Button>
					</Group>
				</Group>
			</header>
		</Box>
	);
}
