import { IconLogout, IconUser } from '@tabler/icons-react';
import { ActionIcon, Button, Menu, Skeleton } from '@mantine/core';
import { useAuth } from '@/auth/hooks';
import { RouteLink } from '../RouteLink/RouteLink';

export function UserMenu() {
	const { user, isAuthenticated, isLoading, logout } = useAuth();

	if (isLoading) {
		return <Skeleton height={36} width={36} circle />;
	}

	if (isAuthenticated && user) {
		return (
			<Menu shadow="md" width={200}>
				<Menu.Target>
					<ActionIcon variant="outline" color="blue" size="lg" radius="xl" bd="2px solid">
						<IconUser size={20} />
					</ActionIcon>
				</Menu.Target>

				<Menu.Dropdown>
					<Menu.Label>{user.Name}</Menu.Label>
					<Menu.Item leftSection={<IconUser size={16} />} disabled>
						Профиль
					</Menu.Item>
					<Menu.Divider />
					<Menu.Item color="red" leftSection={<IconLogout size={16} />} onClick={() => logout()}>
						Выйти
					</Menu.Item>
				</Menu.Dropdown>
			</Menu>
		);
	}

	return (
		<Button component={RouteLink} to="/login">
			Войти
		</Button>
	);
}
