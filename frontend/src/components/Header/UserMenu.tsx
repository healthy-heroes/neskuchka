import { IconLogout, IconSettings, IconUser } from '@tabler/icons-react';
import { ActionIcon, Avatar, Button, Menu, Skeleton } from '@mantine/core';
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
					{user.Avatar ? (
						<Avatar src={user.Avatar} size="md" radius="xl" style={{ cursor: 'pointer' }} />
					) : (
						<ActionIcon variant="outline" color="blue" size="lg" radius="xl" bd="2px solid">
							<IconUser size={20} />
						</ActionIcon>
					)}
				</Menu.Target>

				<Menu.Dropdown>
					<Menu.Label>{user.Name}</Menu.Label>
					<Menu.Item
						component={RouteLink}
						to="/settings"
						leftSection={<IconSettings size={16} />}
					>
						Настройки
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
